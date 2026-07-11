<template>
  <AdminLayout>
    <section class="grid gap-6 lg:grid-cols-3">
      <div
        class="rounded-xl bg-gradient-to-br from-slate-900 to-slate-800 p-6 text-white shadow-elevated ring-1 ring-slate-900/20"
      >
        <p
          class="text-xs uppercase tracking-widest text-blue-300 font-semibold"
        >
          Laporan Baru
        </p>
        <p class="mt-4 text-2xl font-bold">
          {{
            stats.value?.latest?.judul ||
            stats.value?.latest?.title ||
            "Tidak ada laporan terbaru"
          }}
        </p>
        <p class="mt-4 text-sm text-slate-300 leading-relaxed">
          {{
            stats.value?.latest?.deskripsi || stats.value?.latest?.summary || ""
          }}
        </p>
      </div>
      <div class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50">
        <p class="text-xs font-bold uppercase tracking-wider text-slate-600">
          Ringkasan
        </p>
        <div class="mt-5 grid gap-4">
          <div class="rounded-lg bg-slate-50 p-4 hover:bg-slate-100 transition">
            <p class="text-sm text-slate-600">Manajemen Aduan</p>
            <p class="mt-2 text-2xl font-bold text-slate-950">
              {{
                stats.value?.manajemen_count ||
                stats.value?.total_pengaduan ||
                stats.value?.total ||
                "-"
              }}
            </p>
          </div>
          <div class="rounded-lg bg-slate-50 p-4 hover:bg-slate-100 transition">
            <p class="text-sm text-slate-600">Notifikasi</p>
            <p class="mt-2 text-2xl font-bold text-slate-950">
              {{
                stats.value?.notifications_count ||
                stats.value?.notif_count ||
                "-"
              }}
            </p>
          </div>
        </div>
      </div>
      <div class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50">
        <p class="text-xs font-bold uppercase tracking-wider text-slate-600">
          Riwayat Aduan
        </p>
        <div class="mt-4 space-y-4">
          <template v-if="stats.value?.recent?.length">
            <div
              v-for="r in stats.value.recent"
              :key="r.id"
              class="rounded-lg bg-slate-50 p-4 hover:bg-slate-100 transition"
            >
              <p class="font-semibold text-slate-950">
                {{ r.title || r.kategori || r.unit || "—" }}
              </p>
              <p class="mt-2 text-sm text-slate-600">
                {{ r.summary || r.judul || r.deskripsi || "" }}
              </p>
              <p class="mt-3 text-xs text-slate-500">{{ r.time || "" }}</p>
            </div>
          </template>
          <div
            v-else
            class="rounded-lg bg-slate-50 p-4 hover:bg-slate-100 transition"
          >
            <p class="text-sm text-slate-600">Belum ada riwayat aduan.</p>
          </div>
        </div>
      </div>
    </section>
  </AdminLayout>
</template>

<script setup>
import { ref, onMounted } from "vue";
import adminService from "../../services/admin.service";
import { useToastStore } from "../../stores/toast";

const loading = ref(true);
const error = ref("");
const stats = ref({});
const toast = useToastStore();

const load = async () => {
  loading.value = true;
  error.value = "";
  try {
    const res = await adminService.dashboard();
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
