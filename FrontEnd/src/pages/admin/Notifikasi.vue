<template>
  <AdminLayout>
    <section class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50">
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <p class="text-sm font-bold text-slate-950">Notifikasi Admin</p>
          <p class="mt-1 text-sm text-slate-600">Pantau pembaruan laporan dan tindakan sistem.</p>
        </div>
        <button @click="markAllRead" class="rounded-lg bg-slate-100 px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-200">Tandai Semua Dibaca</button>
      </div>

      <div class="mt-6 space-y-3">
        <p v-if="loading" class="rounded-lg bg-slate-50 p-5 text-sm text-slate-600">Memuat notifikasi...</p>
        <p v-else-if="error" class="rounded-lg bg-red-50 p-5 text-sm text-red-700">Gagal memuat notifikasi: {{ error }}</p>
        <article v-else v-for="item in notifications" :key="item.id" class="rounded-lg border border-slate-200 bg-slate-50 p-5">
          <div class="flex items-start justify-between gap-3">
            <div>
              <p class="font-semibold text-slate-950">{{ item.judul || item.title || "Notifikasi" }}</p>
              <p class="mt-2 text-sm text-slate-600">{{ item.isi || item.message || item.body }}</p>
            </div>
            <button v-if="!item.read_at && !item.is_read" @click="markRead(item)" class="rounded-lg bg-blue-950 px-3 py-2 text-xs font-semibold text-white">Dibaca</button>
          </div>
        </article>
        <p v-if="!loading && !error && !notifications.length" class="rounded-lg bg-slate-50 p-5 text-sm text-slate-500">Belum ada notifikasi.</p>
      </div>
    </section>
  </AdminLayout>
</template>

<script setup>
import { onMounted, ref } from "vue";
import AdminLayout from "../../layouts/AdminLayout.vue";
import notifikasiService from "../../services/notifikasi.service";
import { useToastStore } from "../../stores/toast";

const toast = useToastStore();
const notifications = ref([]);
const loading = ref(true);
const error = ref("");

const load = async () => {
  loading.value = true;
  error.value = "";
  try {
    const { data } = await notifikasiService.mine();
    notifications.value = data?.data || data || [];
  } catch (err) {
    error.value = err?.response?.data?.message || err?.message || "Server error";
  } finally {
    loading.value = false;
  }
};

const markRead = async (item) => {
  await notifikasiService.markAsRead(item.id);
  item.read_at = new Date().toISOString();
  item.is_read = true;
};

const markAllRead = async () => {
  try {
    await Promise.all(notifications.value.filter((item) => !item.read_at && !item.is_read).map(markRead));
    toast.add({ type: "success", message: "Notifikasi ditandai dibaca." });
  } catch {
    toast.add({ type: "danger", message: "Gagal menandai notifikasi." });
  }
};

onMounted(load);
</script>
