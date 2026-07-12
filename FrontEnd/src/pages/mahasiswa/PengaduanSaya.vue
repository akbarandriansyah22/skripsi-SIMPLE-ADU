<template>
  <MahasiswaLayout>
    <section class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50">
      <div
        class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between"
      >
        <div>
          <p class="text-sm font-bold text-slate-950">Pengaduan Saya</p>
          <p class="mt-1 text-sm text-slate-600">
            Kelola laporan yang sudah dibuat.
          </p>
        </div>
        <router-link
          class="inline-flex items-center rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-blue-700 active:bg-blue-800 shadow-soft"
          to="/mahasiswa/kirim"
          >Kirim Aduan Baru</router-link
        >
      </div>
      <div class="mt-6 overflow-hidden rounded-lg border border-slate-200">
        <div
          v-if="showSuccess"
          class="mb-4 rounded-lg bg-emerald-50 border border-emerald-200 p-4 text-emerald-800"
        >
          Aduan berhasil dikirim.
        </div>
        <table
          class="min-w-full divide-y divide-slate-200 text-left text-sm text-slate-700"
        >
          <thead class="bg-slate-50 text-slate-600">
            <tr>
              <th
                class="px-6 py-3 font-semibold text-xs uppercase tracking-wider"
              >
                ID Aduan
              </th>
              <th
                class="px-6 py-3 font-semibold text-xs uppercase tracking-wider"
              >
                Judul
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
            <tr v-if="loading">
              <td colspan="5" class="px-6 py-8 text-center text-slate-500">
                Memuat daftar aduan...
              </td>
            </tr>
            <tr v-else-if="error">
              <td colspan="5" class="px-6 py-8 text-center text-red-600">
                Gagal memuat data: {{ error }}
              </td>
            </tr>
            <tr v-else-if="!pengaduans.length">
              <td colspan="5" class="px-6 py-8 text-center text-slate-500">
                Anda belum mengirimkan aduan.
              </td>
            </tr>
            <tr
              v-else
              v-for="p in pengaduans"
              :key="p.id"
              class="hover:bg-slate-50 transition"
            >
              <td class="px-6 py-4 font-medium text-slate-900">
                {{ p.kode_tiket || p.kode || "ADU-" + p.id }}
              </td>
              <td class="px-6 py-4 text-slate-700">{{ p.judul }}</td>
              <td class="px-6 py-4">
                <span :class="statusBadgeClass(p.status)">{{ p.status }}</span>
              </td>
              <td class="px-6 py-4 text-slate-600">
                {{ formatDate(p.created_at) }}
              </td>
              <td class="px-6 py-4">
                <button
                  type="button"
                  @click="openDetail(p.id)"
                  class="rounded-lg bg-slate-100 px-3 py-2 text-xs font-semibold text-slate-700 transition hover:bg-slate-200"
                >
                  Detail
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>
  </MahasiswaLayout>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import MahasiswaLayout from "../../layouts/MahasiswaLayout.vue";
import pengaduanService from "../../services/pengaduan.service";
import { useToastStore } from "../../stores/toast";

const route = useRoute();
const router = useRouter();
const showSuccess = ref(false);
const pengaduans = ref([]);
const loading = ref(true);
const error = ref("");
const toast = useToastStore();

function openDetail(id) {
  if (!id) return;
  router.push({ name: "MahasiswaDetail", params: { id: String(id) } });
}

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
    const d = new Date(iso);
    return d.toLocaleDateString();
  } catch {
    return "";
  }
}

const load = async () => {
  loading.value = true;
  error.value = "";
  try {
    const res = await pengaduanService.myPengaduan();
    pengaduans.value = res.data?.data || res.data || [];
  } catch (err) {
    error.value = err?.message || "Server error";
    toast.add({ type: "danger", message: error.value });
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  if (route.query && route.query.success === "1") {
    showSuccess.value = true;
    // remove query from URL after showing
    router.replace({ name: "MahasiswaPengaduan" });
    setTimeout(() => (showSuccess.value = false), 4000);
  }
  load();
});
</script>
