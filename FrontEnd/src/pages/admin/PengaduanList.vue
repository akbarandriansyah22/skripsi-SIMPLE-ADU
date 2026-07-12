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
      <div class="mt-6 grid gap-3 md:grid-cols-4">
        <input
          v-model="filters.search"
          type="search"
          placeholder="Cari ID, nama, judul..."
          class="rounded-lg border border-slate-200 px-4 py-3 text-sm"
        />
        <select v-model="filters.status" class="rounded-lg border border-slate-200 px-4 py-3 text-sm">
          <option value="">Semua Status</option>
          <option v-for="status in statusOptions" :key="status" :value="status">{{ status }}</option>
        </select>
        <select v-model="filters.kategori" class="rounded-lg border border-slate-200 px-4 py-3 text-sm">
          <option value="">Semua Kategori</option>
          <option v-for="kategori in kategoriOptions" :key="kategori" :value="kategori">{{ kategori }}</option>
        </select>
        <select v-model="sortBy" class="rounded-lg border border-slate-200 px-4 py-3 text-sm">
          <option value="newest">Terbaru</option>
          <option value="oldest">Terlama</option>
          <option value="title">Judul A-Z</option>
        </select>
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
            <tr v-else-if="!filteredPengaduans.length" class="">
              <td colspan="7" class="px-6 py-8 text-center text-slate-500">
                Tidak ada aduan.
              </td>
            </tr>
            <tr
              v-else
              v-for="p in pagedPengaduans"
              :key="p.id"
              class="hover:bg-slate-50 transition"
            >
              <td class="px-6 py-4 font-medium text-slate-900">
                {{ p.kode_tiket || p.kode || "ADU-" + p.id }}
              </td>
              <td class="px-6 py-4 text-slate-700">{{ p.judul }}</td>
              <td class="px-6 py-4 text-slate-600">
                {{ p.pelapor?.nama || p.pelapor_name || p.user?.nama_lengkap || p.user?.nama || "-" }}
              </td>
              <td class="px-6 py-4 text-slate-600">{{ p.kategori?.nama || p.kategori_prediksi || p.kategori || "-" }}</td>
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
      <div class="mt-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <p class="text-sm text-slate-600">
          Menampilkan {{ pagedPengaduans.length }} dari {{ filteredPengaduans.length }} aduan
        </p>
        <div class="flex items-center gap-2">
          <button :disabled="page === 1" @click="page--" class="rounded-lg bg-slate-100 px-3 py-2 text-sm disabled:opacity-40">Sebelumnya</button>
          <span class="rounded-lg bg-blue-950 px-3 py-2 text-sm font-semibold text-white">{{ page }}</span>
          <button :disabled="page >= pageCount" @click="page++" class="rounded-lg bg-slate-100 px-3 py-2 text-sm disabled:opacity-40">Berikutnya</button>
        </div>
      </div>
    </section>
  </AdminLayout>
</template>

<script setup>
import { computed, ref, onMounted, watch } from "vue";
import AdminLayout from "../../layouts/AdminLayout.vue";
import adminService from "../../services/admin.service";

const pengaduans = ref([]);
const loading = ref(true);
const error = ref("");
const page = ref(1);
const perPage = 10;
const sortBy = ref("newest");
const filters = ref({ search: "", status: "", kategori: "" });

const statusOptions = computed(() => uniqueValues("status"));
const kategoriOptions = computed(() => uniqueValues("kategori"));

const uniqueValues = (key) =>
  [...new Set(pengaduans.value.map((item) => key === "kategori" ? (item?.kategori?.nama || item?.kategori_prediksi || item?.kategori) : item?.[key]).filter(Boolean))];

const filteredPengaduans = computed(() => {
  const search = filters.value.search.toLowerCase();
  const items = pengaduans.value.filter((item) => {
    const text = `${item.kode_tiket || ""} ${item.kode || ""} ${item.judul || ""} ${item.pelapor?.nama || ""} ${item.pelapor_name || ""} ${item.user?.nama_lengkap || ""} ${item.user?.nama || ""}`.toLowerCase();
    const statusOk = !filters.value.status || item.status === filters.value.status;
    const kategori = item.kategori?.nama || item.kategori_prediksi || item.kategori;
    const kategoriOk = !filters.value.kategori || kategori === filters.value.kategori;
    return text.includes(search) && statusOk && kategoriOk;
  });
  return [...items].sort((a, b) => {
    if (sortBy.value === "title") return String(a.judul || "").localeCompare(String(b.judul || ""));
    const aTime = new Date(a.created_at || 0).getTime();
    const bTime = new Date(b.created_at || 0).getTime();
    return sortBy.value === "oldest" ? aTime - bTime : bTime - aTime;
  });
});
const pageCount = computed(() => Math.max(1, Math.ceil(filteredPengaduans.value.length / perPage)));
const pagedPengaduans = computed(() => filteredPengaduans.value.slice((page.value - 1) * perPage, page.value * perPage));

watch([filters, sortBy], () => {
  page.value = 1;
}, { deep: true });

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
    pengaduans.value = res.data?.data || res.data || [];
  } catch (err) {
    error.value = err?.message || "Server error";
  } finally {
    loading.value = false;
  }
};

onMounted(load);
</script>
