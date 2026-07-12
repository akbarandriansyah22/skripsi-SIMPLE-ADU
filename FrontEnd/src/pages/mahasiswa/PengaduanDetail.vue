<template>
  <MahasiswaLayout>
    <section class="grid gap-6 lg:grid-cols-[1.3fr_0.7fr]">
      <div class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50">
        <div v-if="loading" class="text-sm text-slate-600">Memuat detail aduan...</div>
        <div v-else-if="error" class="rounded-lg bg-red-50 p-5 text-sm text-red-700">Gagal memuat detail: {{ error }}</div>
        <template v-else>
          <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
            <div>
              <p class="text-xs font-semibold uppercase tracking-wider text-blue-600">{{ detail?.kode_tiket || detail?.kode || "ADU-" + detail?.id }}</p>
              <h2 class="mt-2 text-2xl font-bold text-slate-950">{{ detail?.judul || "Detail Aduan" }}</h2>
              <p class="mt-1 text-sm text-slate-600">{{ formatDate(detail?.created_at) }}</p>
            </div>
            <span :class="statusBadgeClass(detail?.status)">{{ detail?.status || "-" }}</span>
          </div>

          <div class="mt-6 rounded-lg border border-slate-200 bg-slate-50 p-5">
            <p class="text-sm font-bold text-slate-950">Isi Laporan</p>
            <p class="mt-3 text-sm leading-relaxed text-slate-700">{{ detail?.deskripsi || "Tidak ada deskripsi." }}</p>
          </div>

          <div class="mt-6 rounded-lg border border-slate-200 bg-slate-50 p-5">
            <p class="text-sm font-bold text-slate-950">Percakapan Admin</p>
            <div class="mt-4 space-y-3">
              <article v-for="item in responses" :key="item.id" class="rounded-lg border border-slate-200 bg-white p-4">
                <p class="text-sm font-semibold text-slate-950">{{ item.user?.nama_lengkap || item.user?.nama || item.author || item.role || "Admin" }}</p>
                <p class="mt-2 text-sm text-slate-700">{{ item.pesan || item.message || item.body }}</p>
              </article>
              <p v-if="!responses.length" class="text-sm text-slate-500">Belum ada percakapan.</p>
            </div>
          </div>
        </template>
      </div>

      <aside class="space-y-6">
        <div class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50">
          <p class="text-sm font-bold text-slate-950">Status Laporan</p>
          <ol class="mt-5 space-y-4 text-sm">
            <li v-for="step in timeline" :key="step.label" class="flex gap-3">
              <span :class="step.done ? 'bg-blue-950' : 'bg-slate-300'" class="mt-1 h-3 w-3 rounded-full"></span>
              <span :class="step.done ? 'text-slate-900' : 'text-slate-500'">{{ step.label }}</span>
            </li>
          </ol>
        </div>
        <div class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50">
          <p class="text-sm font-bold text-slate-950">Ringkasan AI</p>
          <dl class="mt-5 space-y-3 text-sm">
            <div><dt class="font-semibold text-slate-900">Kategori</dt><dd class="text-slate-600">{{ detail?.kategori?.nama || detail?.kategori_prediksi || detail?.kategori || "-" }}</dd></div>
            <div><dt class="font-semibold text-slate-900">Sentimen</dt><dd class="text-slate-600">{{ detail?.sentimen || "-" }}</dd></div>
            <div><dt class="font-semibold text-slate-900">Urgensi</dt><dd class="text-slate-600">{{ detail?.urgensi || "-" }}</dd></div>
          </dl>
        </div>
      </aside>
    </section>
  </MahasiswaLayout>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import MahasiswaLayout from "../../layouts/MahasiswaLayout.vue";
import pengaduanService from "../../services/pengaduan.service";

const props = defineProps({ id: [String, Number] });
const route = useRoute();
const detail = ref(null);
const loading = ref(true);
const error = ref("");
const id = props.id || route.params.id;

const responses = computed(() => detail.value?.responses || detail.value?.respon || detail.value?.comments || []);
const timeline = computed(() => {
  const status = String(detail.value?.status || "").toLowerCase();
  return [
    { label: "Terkirim", done: true },
    { label: "Diverifikasi", done: !status.includes("terkirim") },
    { label: "Diproses", done: status.includes("proses") || status.includes("selesai") },
    { label: "Selesai", done: status.includes("selesai") },
  ];
});

function statusBadgeClass(status) {
  const s = String(status || "").toLowerCase();
  if (s.includes("selesai")) return "rounded-full bg-emerald-100 px-3 py-1 text-xs font-semibold text-emerald-700";
  if (s.includes("proses")) return "rounded-full bg-amber-100 px-3 py-1 text-xs font-semibold text-amber-700";
  return "rounded-full bg-slate-100 px-3 py-1 text-xs font-semibold text-slate-700";
}

function formatDate(iso) {
  return iso ? new Date(iso).toLocaleString("id-ID") : "";
}

onMounted(async () => {
  try {
    const { data } = await pengaduanService.detail(id);
    detail.value = data?.data || data;
  } catch (err) {
    error.value = err?.response?.data?.message || err?.message || "Server error";
  } finally {
    loading.value = false;
  }
});
</script>
