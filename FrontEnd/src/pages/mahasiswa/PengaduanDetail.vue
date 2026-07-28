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

          <div v-if="detail?.lampiran" class="mt-6 rounded-lg border border-slate-200 bg-slate-50 p-5">
            <p class="text-sm font-bold text-slate-950">Lampiran</p>
            <a v-if="attachmentUrls.complaint" :href="attachmentUrls.complaint" target="_blank" rel="noreferrer" class="mt-3 inline-flex rounded-lg bg-blue-50 px-4 py-2 text-sm font-semibold text-blue-700 hover:bg-blue-100">
              Buka {{ detail?.lampiran_nama_asli || detail.lampiran }}
            </a>
            <p v-else class="mt-3 text-sm text-slate-500">Memuat lampiran...</p>
          </div>

          <div class="mt-6 rounded-lg border border-slate-200 bg-slate-50 p-5">
            <p class="text-sm font-bold text-slate-950">Percakapan Admin</p>
            <div class="mt-4 space-y-3">
              <article v-for="item in responses" :key="item.id" class="rounded-lg border border-slate-200 bg-white p-4">
                <p class="text-sm font-semibold text-slate-950">{{ item.user?.nama_lengkap || item.user?.nama || item.author || "Admin" }}</p>
                <p class="mt-1 text-xs text-slate-500">{{ responseRole(item) }}</p>
                <p class="mt-2 text-sm text-slate-700">{{ item.pesan || item.message || item.body }}</p>
                <p class="mt-2 text-xs text-slate-500">{{ formatDate(item.created_at) }}</p>
                <a v-if="attachmentUrls.responses[item.id]" :href="attachmentUrls.responses[item.id]" target="_blank" rel="noreferrer" class="mt-2 inline-flex text-xs font-semibold text-blue-700 hover:underline">Buka lampiran respons</a>
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
            <li v-for="step in timeline" :key="step.id || step.created_at" class="flex gap-3">
              <span class="mt-1 h-3 w-3 rounded-full bg-blue-950"></span>
              <span class="text-slate-900">
                <span class="font-semibold">{{ step.status_lama || "Status awal" }}{{ step.status_lama ? " → " : "" }}{{ step.status_baru || detail?.status }}</span>
                <span v-if="step.catatan" class="mt-1 block text-slate-600">{{ step.catatan }}</span>
                <span class="mt-1 block text-xs text-slate-500">{{ formatDate(step.created_at) }}</span>
              </span>
            </li>
          </ol>
        </div>
        <div class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50">
          <p class="text-sm font-bold text-slate-950">Ringkasan AI</p>
          <dl class="mt-5 space-y-3 text-sm">
            <div><dt class="font-semibold text-slate-900">Kategori</dt><dd class="text-slate-600">{{ detail?.kategori?.nama || detail?.kategori || "-" }}</dd></div>
            <div><dt class="font-semibold text-slate-900">Status AI</dt><dd class="text-slate-600">{{ detail?.ai_status || "pending" }}</dd></div>
            <template v-if="detail?.ai_status === 'success'">
              <div><dt class="font-semibold text-slate-900">Skor Sentimen</dt><dd class="text-slate-600">{{ detail?.skor_sentimen }}</dd></div>
              <div><dt class="font-semibold text-slate-900">Sentimen</dt><dd class="text-slate-600">{{ detail?.sentimen || "-" }}</dd></div>
              <div><dt class="font-semibold text-slate-900">Urgensi</dt><dd class="text-slate-600">{{ detail?.urgensi || "-" }}</dd></div>
            </template>
            <div v-else><dd class="text-slate-600">Analisis AI sedang menunggu pemrosesan.</dd></div>
          </dl>
        </div>
      </aside>
    </section>
  </MahasiswaLayout>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import MahasiswaLayout from "../../layouts/MahasiswaLayout.vue";
import pengaduanService from "../../services/pengaduan.service";

const props = defineProps({ id: [String, Number] });
const route = useRoute();
const detail = ref(null);
const loading = ref(true);
const error = ref("");
const id = props.id || route.params.id;
const attachmentUrls = ref({ complaint: "", responses: {} });

const responses = computed(() => [...(
  detail.value?.respon_pengaduan ||
  detail.value?.responses ||
  detail.value?.respon ||
  detail.value?.comments ||
  []
)].sort((a, b) => new Date(a.created_at || 0).getTime() - new Date(b.created_at || 0).getTime()));
const timeline = computed(() => {
  const history = detail.value?.riwayat_status_pengaduan || detail.value?.riwayat_status || [];
  const sorted = [...history].sort((a, b) => new Date(a.created_at || 0).getTime() - new Date(b.created_at || 0).getTime());
  return sorted.length ? sorted : detail.value ? [{ status_baru: detail.value.status, created_at: detail.value.created_at }] : [];
});

function responseRole(item) {
  const role = String(item.user?.role || item.role || "").toLowerCase();
  if (role.includes("kasubag")) return "Kasubag";
  if (role.includes("admin")) return "Admin Fakultas";
  if (role.includes("pimpinan")) return "Pimpinan Fakultas";
  if (role.includes("mahasiswa")) return "Mahasiswa";
  return item.user?.role || item.role || "Petugas";
}

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
    await loadAttachments();
  } catch (err) {
    error.value = err?.response?.data?.message || err?.message || "Server error";
  } finally {
    loading.value = false;
  }
});

async function loadAttachments() {
  const next = { complaint: "", responses: {} };
  if (detail.value?.lampiran_url) {
    try { next.complaint = URL.createObjectURL(await pengaduanService.getAttachment(detail.value.lampiran_url)); } catch {}
  }
  await Promise.all(responses.value.map(async (item) => {
    if (!item.lampiran_url) return;
    try { next.responses[item.id] = URL.createObjectURL(await pengaduanService.getAttachment(item.lampiran_url)); } catch {}
  }));
  attachmentUrls.value = next;
}

onBeforeUnmount(() => {
  if (attachmentUrls.value.complaint) URL.revokeObjectURL(attachmentUrls.value.complaint);
  Object.values(attachmentUrls.value.responses).forEach((url) => URL.revokeObjectURL(url));
});
</script>
