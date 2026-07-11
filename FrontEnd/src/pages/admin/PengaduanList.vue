<template>
  <AdminLayout>
    <section
      class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50"
    >
      <div
        class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between"
      >
        <div>
          <p class="text-sm font-bold text-slate-950">Manajemen Aduan</p>
          <p class="mt-1 text-sm text-slate-600">
            Cari dan tindak lanjuti laporan masuk.
          </p>
        </div>
        <button
          class="rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-blue-700 active:bg-blue-800 shadow-soft"
        >
          Ekspor Laporan
        </button>
      </div>
      <div class="mt-6 overflow-hidden rounded-lg border border-slate-200">
        <table
          class="min-w-full divide-y divide-slate-200 text-left text-sm text-slate-700"
        >
          <thead class="bg-slate-50 text-slate-600">
            <tr>
              <th
                class="px-6 py-3 font-semibold text-xs uppercase tracking-wider"
              >
                ID Laporan
              </th>
              <th
                class="px-6 py-3 font-semibold text-xs uppercase tracking-wider"
              >
                Judul
              </th>
              <th
                class="px-6 py-3 font-semibold text-xs uppercase tracking-wider"
              >
                Pelapor
              </th>
              <th
                class="px-6 py-3 font-semibold text-xs uppercase tracking-wider"
              >
                Kategori
              </th>
              <th
                class="px-6 py-3 font-semibold text-xs uppercase tracking-wider"
              >
                Status
              </th>
              <th
                class="px-6 py-3 font-semibold text-xs uppercase tracking-wider"
              >
                Tanggal
              </th>
              <th
                class="px-6 py-3 font-semibold text-xs uppercase tracking-wider"
              >
                Aksi
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-200 bg-white">
            <tr v-if="loading" class="">
              <td colspan="7" class="px-6 py-8 text-center text-slate-500">
                Memuat data...
              </td>
            </tr>
            <tr v-else-if="error" class="">
              <td colspan="7" class="px-6 py-8 text-center text-red-600">
                Gagal memuat data: {{ error }}
              </td>
            </tr>
            <tr v-else-if="!pengaduans.length" class="">
              <td colspan="7" class="px-6 py-8 text-center text-slate-500">
                Tidak ada aduan.
              </td>
            </tr>
            <tr
              v-else
              v-for="p in pengaduans"
              :key="p.id"
              class="hover:bg-slate-50 transition"
            >
              <td class="px-6 py-4 font-medium text-slate-900">
                {{ p.kode || "ADU-" + p.id }}
              </td>
              <td class="px-6 py-4 text-slate-700">{{ p.judul }}</td>
              <td class="px-6 py-4 text-slate-600">
                {{ p.pelapor?.nama || p.pelapor_name || p.user?.nama || "-" }}
              </td>
              <td class="px-6 py-4 text-slate-600">{{ p.kategori || "-" }}</td>
              <td class="px-6 py-4">
                <span :class="statusBadgeClass(p.status)">{{
                  p.status || "-"
                }}</span>
              </td>
              <td class="px-6 py-4 text-slate-600">
                {{ formatDate(p.created_at) }}
              </td>
              <td class="px-6 py-4">
                <router-link
                  class="rounded-lg bg-slate-100 px-3 py-2 text-xs font-semibold text-slate-700 transition hover:bg-slate-200 inline-block"
                  :to="{ name: 'AdminDetail', params: { id: String(p.id) } }"
                  >Tinjau</router-link
                >
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <p class="mt-4 text-sm text-slate-600">
        Menampilkan 1 - 10 dari 45 aduan
      </p>
    </section>
  </AdminLayout>
</template>

<script setup>
import { ref, onMounted } from "vue";
import adminService from "../../services/admin.service";

const pengaduans = ref([]);
const loading = ref(true);
const error = ref("");

function statusBadgeClass(status) {
  if (!status)
    return "rounded-full px-3 py-1 text-xs font-semibold inline-block bg-slate-100 text-slate-700";
  const s = String(status).toLowerCase();
  if (s.includes("selesai") || s.includes("done") || s.includes("seles"))
    return "rounded-full bg-emerald-100 px-3 py-1 text-xs font-semibold text-emerald-700 inline-block";
  if (
    s.includes("proses") ||
    s.includes("diproses") ||
    s.includes("in progress")
  )
    return "rounded-full bg-amber-100 px-3 py-1 text-xs font-semibold text-amber-700 inline-block";
  if (s.includes("ditolak") || s.includes("tolak") || s.includes("rejected"))
    return "rounded-full bg-red-100 px-3 py-1 text-xs font-semibold text-red-700 inline-block";
  return "rounded-full px-3 py-1 text-xs font-semibold inline-block bg-slate-100 text-slate-700";
}

function formatDate(iso) {
  if (!iso) return "";
  try {
    return new Date(iso).toLocaleDateString();
  } catch {
    return "";
  }
}

const load = async () => {
  loading.value = true;
  error.value = "";
  try {
    const res = await adminService.getAllPengaduan();
    pengaduans.value = res.data || [];
  } catch (err) {
    error.value = err?.message || "Server error";
  } finally {
    loading.value = false;
  }
};

onMounted(load);
</script>
