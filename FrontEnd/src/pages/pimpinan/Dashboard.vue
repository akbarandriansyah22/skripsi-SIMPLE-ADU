<template>
  <PimpinanLayout>
    <section class="grid gap-6 lg:grid-cols-3">
      <div class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50">
        <p
          class="text-xs uppercase tracking-widest text-blue-600 font-semibold"
        >
          Validasi Anda
        </p>
        <p class="mt-4 text-4xl font-bold text-slate-950">
          {{ stats.value?.validasi_count || stats.value?.pending || "-" }}
        </p>
        <p class="mt-3 text-sm text-slate-600 leading-relaxed">
          Pengaduan yang menunggu keputusan resmi Anda.
        </p>
      </div>
      <div class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50">
        <p
          class="text-xs uppercase tracking-widest text-blue-600 font-semibold"
        >
          Urgensi Tinggi
        </p>
        <p class="mt-4 text-4xl font-bold text-slate-950">
          {{ stats.value?.urgent_count || stats.value?.urgensi || "-" }}
        </p>
        <p class="mt-3 text-sm text-slate-600 leading-relaxed">
          Laporan prioritas yang perlu tindak lanjut cepat.
        </p>
      </div>
      <div class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50">
        <p
          class="text-xs uppercase tracking-widest text-blue-600 font-semibold"
        >
          Disposisi
        </p>
        <p class="mt-4 text-4xl font-bold text-slate-950">
          {{ stats.value?.disposisi_count || stats.value?.disposisi || "-" }}
        </p>
        <p class="mt-3 text-sm text-slate-600 leading-relaxed">
          Instruksi yang telah dikirim ke unit terkait.
        </p>
      </div>
    </section>
    <section class="mt-6 grid gap-6 lg:grid-cols-2">
      <div class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50">
        <div class="flex items-center justify-between">
          <p class="text-sm font-bold text-slate-950">
            Pengaduan Urgensi Tinggi
          </p>
          <button
            class="rounded-lg bg-slate-100 px-3 py-2 text-xs font-semibold text-slate-700 transition hover:bg-slate-200"
          >
            Lihat Semua
          </button>
        </div>
        <div class="mt-6 space-y-4">
          <template v-if="stats.value?.urgent?.length">
            <div
              v-for="r in stats.value.urgent"
              :key="r.id"
              class="rounded-lg border border-slate-200 bg-slate-50 p-4 hover:bg-slate-100 transition"
            >
              <p class="font-semibold text-slate-950">
                {{ r.kode || r.judul || r.title }}
              </p>
              <p class="mt-2 text-sm text-slate-600">
                {{ r.judul || r.deskripsi || r.summary }}
              </p>
              <p class="mt-3 text-xs text-slate-500">
                {{ r.tanggal || r.date || "" }}
              </p>
            </div>
          </template>
          <div
            v-else
            class="rounded-lg border border-slate-200 bg-slate-50 p-4 hover:bg-slate-100 transition"
          >
            <p class="text-sm text-slate-600">
              Tidak ada pengaduan urgensi tinggi.
            </p>
          </div>
        </div>
      </div>
      <div class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50">
        <p class="text-sm font-bold text-slate-950">Tugas Validasi</p>
        <div class="mt-6 space-y-4 text-sm text-slate-700 leading-relaxed">
          <p>Review jawaban resmi yang belum dipublikasikan.</p>
          <p>Pastikan keputusan sesuai prosedur fakultas.</p>
        </div>
        <button
          class="mt-6 rounded-lg bg-emerald-600 px-4 py-3 text-sm font-semibold text-white transition hover:bg-emerald-700 active:bg-emerald-800 shadow-soft w-full"
        >
          Validasi Sekarang
        </button>
      </div>
    </section>
  </PimpinanLayout>
</template>

<script setup>
import { ref, onMounted } from "vue";
import pimpinanService from "../../services/pimpinan.service";
import { useToastStore } from "../../stores/toast";

const loading = ref(true);
const error = ref("");
const stats = ref({});
const toast = useToastStore();

const load = async () => {
  loading.value = true;
  error.value = "";
  try {
    const res = await pimpinanService.dashboard();
    stats.value = res.data || {};
  } catch (err) {
    error.value = err?.message || "Server error";
    toast.add({ type: "danger", message: error.value });
  } finally {
    loading.value = false;
  }
};

onMounted(load);
</script>
