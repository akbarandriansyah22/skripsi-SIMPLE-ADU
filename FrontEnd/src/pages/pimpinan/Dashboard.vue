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
          {{ stats?.belum_disposisi || stats?.validasi_count || stats?.pending || 0 }}
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
          {{ stats?.total_urgensi_tinggi || stats?.urgent_count || stats?.urgensi || 0 }}
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
          {{ stats?.sudah_disposisi || stats?.disposisi_count || stats?.disposisi || 0 }}
        </p>
        <p class="mt-3 text-sm text-slate-600 leading-relaxed">
          Instruksi yang telah dikirim ke unit terkait.
        </p>
      </div>
    </section>
    <section class="mt-6 grid grid-cols-2 gap-3 lg:grid-cols-4"><div v-for="item in [{ label: 'Total Pengaduan', value: stats?.total_pengaduan }, { label: 'Belum Dikerjakan', value: stats?.belum_dikerjakan }, { label: 'Sedang Diproses', value: stats?.sedang_diproses }, { label: 'Selesai', value: stats?.selesai }, { label: 'Ditolak', value: stats?.ditolak }, { label: 'Urgensi Rendah', value: stats?.urgensi_rendah }, { label: 'Urgensi Sedang', value: stats?.urgensi_sedang }, { label: 'Urgensi Tinggi', value: stats?.total_urgensi_tinggi }]" :key="item.label" class="rounded-xl bg-white p-4 shadow-card ring-1 ring-slate-200/50"><p class="text-[11px] text-slate-500">{{ item.label }}</p><p class="mt-2 text-2xl font-bold text-slate-950">{{ item.value ?? 0 }}</p></div></section>
    <section class="mt-6 grid gap-6 lg:grid-cols-2">
      <div class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50">
        <div class="flex items-center justify-between">
          <p class="text-sm font-bold text-slate-950">
            Pengaduan Urgensi Tinggi
          </p>
          <router-link
            to="/pimpinan/validasi"
            class="rounded-lg bg-slate-100 px-3 py-2 text-xs font-semibold text-slate-700 transition hover:bg-slate-200"
          >
            Lihat Semua
          </router-link>
        </div>
        <div class="mt-6 space-y-4">
          <template v-if="urgent.length">
            <div
              v-for="r in urgent"
              :key="r.id"
              class="rounded-lg border border-slate-200 bg-slate-50 p-4 hover:bg-slate-100 transition"
            >
              <p class="font-semibold text-slate-950">
                {{ r.kode_tiket || r.kode || r.judul || r.title }}
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
        <router-link
          to="/pimpinan/validasi"
          class="mt-6 rounded-lg bg-emerald-600 px-4 py-3 text-sm font-semibold text-white transition hover:bg-emerald-700 active:bg-emerald-800 shadow-soft w-full"
        >
          Validasi Sekarang
        </router-link>
        <router-link
          to="/pimpinan/monitoring"
          class="mt-3 block rounded-lg bg-slate-100 px-4 py-3 text-center text-sm font-semibold text-slate-700 transition hover:bg-slate-200"
        >
          Pantau Progress Aduan
        </router-link>
      </div>
    </section>
  </PimpinanLayout>
</template>

<script setup>
import { ref, onMounted } from "vue";
import PimpinanLayout from "../../layouts/PimpinanLayout.vue";
import pimpinanService from "../../services/pimpinan.service";
import { useToastStore } from "../../stores/toast";

const loading = ref(true);
const error = ref("");
const stats = ref({});
const urgent = ref([]);
const toast = useToastStore();

const load = async () => {
  loading.value = true;
  error.value = "";
  try {
    const res = await pimpinanService.dashboard();
    stats.value = res.data?.data || res.data || {};
    const urgentRes = await pimpinanService.urgentReports();
    urgent.value = urgentRes.data?.data || urgentRes.data || [];
  } catch (err) {
    error.value = err?.message || "Server error";
    toast.add({ type: "danger", message: error.value });
  } finally {
    loading.value = false;
  }
};

onMounted(load);
</script>
