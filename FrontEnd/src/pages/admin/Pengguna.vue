<template>
  <AdminLayout>
    <section class="grid gap-6 lg:grid-cols-[0.35fr_1fr]">
      <aside class="space-y-4">
        <div class="rounded-xl bg-white p-5 shadow-card ring-1 ring-slate-200/50">
          <p class="text-sm font-bold text-slate-950">Pencarian</p>
          <input v-model="query" type="search" placeholder="Cari nama atau email..." class="mt-4 w-full rounded-lg border border-slate-200 px-4 py-3 text-sm" />
        </div>
        <div class="rounded-xl bg-white p-5 shadow-card ring-1 ring-slate-200/50">
          <p class="text-sm font-bold text-slate-950">Filter Peran</p>
          <div class="mt-4 space-y-2 text-sm text-slate-700">
            <label v-for="role in roles" :key="role.value" class="flex items-center gap-2">
              <input v-model="selectedRoles" :value="role.value" type="checkbox" class="rounded border-slate-300 text-blue-950" />
              {{ role.label }}
            </label>
          </div>
        </div>
      </aside>

      <section class="rounded-xl bg-white shadow-card ring-1 ring-slate-200/50">
        <div class="overflow-x-auto">
          <table class="min-w-full text-left text-sm">
            <thead class="bg-slate-50 text-xs uppercase tracking-wider text-slate-500">
              <tr>
                <th class="px-6 py-3">Pengguna</th>
                <th class="px-6 py-3">Peran</th>
                <th class="px-6 py-3">Status</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-200">
              <tr v-if="loading"><td colspan="3" class="px-6 py-8 text-center text-slate-500">Memuat pengguna...</td></tr>
              <tr v-else-if="error"><td colspan="3" class="px-6 py-8 text-center text-amber-700">{{ error }}</td></tr>
              <tr v-else-if="!filtered.length"><td colspan="3" class="px-6 py-8 text-center text-slate-500">Belum ada data pengguna dari endpoint admin.</td></tr>
              <tr v-for="user in paged" :key="user.id">
                <td class="px-6 py-4">
                  <p class="font-semibold text-slate-950">{{ user.nama_lengkap || user.nama || user.name }}</p>
                  <p class="text-xs text-slate-500">{{ user.email }}</p>
                </td>
                <td class="px-6 py-4"><span class="rounded-full bg-slate-100 px-3 py-1 text-xs font-semibold text-slate-700">{{ user.role }}</span></td>
                <td class="px-6 py-4 text-slate-600">{{ user.is_active === false ? "Nonaktif" : "Aktif" }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div class="flex items-center justify-between border-t border-slate-200 px-6 py-4 text-sm text-slate-600">
          <span>Menampilkan {{ paged.length }} dari {{ filtered.length }} pengguna</span>
          <div class="flex gap-2">
            <button :disabled="page === 1" @click="page--" class="rounded-lg bg-slate-100 px-3 py-2 disabled:opacity-40">Sebelumnya</button>
            <button :disabled="page >= pageCount" @click="page++" class="rounded-lg bg-slate-100 px-3 py-2 disabled:opacity-40">Berikutnya</button>
          </div>
        </div>
      </section>
    </section>
  </AdminLayout>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import AdminLayout from "../../layouts/AdminLayout.vue";
import api from "../../services/axios";

const users = ref([]);
const loading = ref(true);
const error = ref("");
const query = ref("");
const page = ref(1);
const perPage = 10;
const selectedRoles = ref([]);
const roles = [
  { label: "Mahasiswa", value: "mahasiswa" },
  { label: "Admin", value: "admin" },
  { label: "Pimpinan", value: "pimpinan" },
];

const filtered = computed(() => {
  const term = query.value.toLowerCase();
  return users.value.filter((user) => {
    const text = `${user.nama_lengkap || user.nama || user.name || ""} ${user.email || ""}`.toLowerCase();
    const roleOk = !selectedRoles.value.length || selectedRoles.value.some((role) => String(user.role || "").toLowerCase().includes(role));
    return text.includes(term) && roleOk;
  });
});
const pageCount = computed(() => Math.max(1, Math.ceil(filtered.value.length / perPage)));
const paged = computed(() => filtered.value.slice((page.value - 1) * perPage, page.value * perPage));

onMounted(async () => {
  try {
    const { data } = await api.get("/api/admin/users");
    users.value = data?.data || data || [];
  } catch (err) {
    error.value = "Endpoint manajemen pengguna belum tersedia di backend.";
  } finally {
    loading.value = false;
  }
});
</script>
