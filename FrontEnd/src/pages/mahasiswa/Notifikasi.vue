<template>
  <MahasiswaLayout>
    <section
      class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50"
    >
      <div
        class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between"
      >
        <div>
          <p class="text-sm font-bold text-slate-950">Notifikasi</p>
          <p class="mt-1 text-sm text-slate-600">
            Notifikasi terbaru dari sistem.
          </p>
        </div>
        <button
          @click="markAllRead"
          class="rounded-lg bg-slate-100 px-4 py-2 text-sm font-medium text-slate-700 transition hover:bg-slate-200"
        >
          Tandai Semua Dibaca
        </button>
      </div>
      <div class="mt-6 space-y-4">
        <template v-if="loading">
          <div class="rounded-lg border border-slate-200 bg-slate-50 p-5">
            <p class="text-sm text-slate-600">Memuat notifikasi...</p>
          </div>
        </template>
        <template v-else-if="error">
          <div class="rounded-lg border border-red-200 bg-red-50 p-5">
            <p class="text-sm text-red-600">
              Gagal memuat notifikasi: {{ error }}
            </p>
          </div>
        </template>
        <template v-else-if="notifications.length">
          <article
            v-for="item in notifications"
            :key="item.id"
            @click="openNotification(item)"
            class="cursor-pointer rounded-lg border border-slate-200 bg-slate-50 p-5 shadow-soft hover:bg-slate-100 transition flex flex-col gap-2"
          >
            <div class="flex items-center justify-between gap-3">
              <p class="font-semibold text-slate-900">
                {{ item.title || item.judul || item.subject }}
              </p>
              <span
                class="text-xs uppercase tracking-widest text-slate-500 font-medium"
                >{{ timeAgo(item.created_at) }}</span
              >
            </div>
            <p class="mt-1 text-sm text-slate-600">
              {{ item.message || item.body || item.content }}
            </p>
            <div
              v-if="!item.read_at"
              class="mt-2 inline-block rounded-full bg-blue-100 text-blue-700 px-2 py-1 text-xs font-semibold"
            >
              Baru
            </div>
          </article>
        </template>
        <p v-else class="text-sm text-slate-500">Tidak ada notifikasi.</p>
      </div>
    </section>
  </MahasiswaLayout>
</template>

<script setup>
import { ref, onMounted } from "vue";
import notifikasiService from "../../services/notifikasi.service";
import { useToastStore } from "../../stores/toast";
import { useRouter } from "vue-router";

const notifications = ref([]);
const router = useRouter();
const toast = useToastStore();
const loading = ref(true);
const error = ref("");

const load = async () => {
  loading.value = true;
  error.value = "";
  try {
    const res = await notifikasiService.mine();
    notifications.value = res.data || [];
  } catch (err) {
    error.value = err?.message || "Server error";
    toast.add({ type: "danger", message: "Gagal memuat notifikasi." });
  } finally {
    loading.value = false;
  }
};

onMounted(load);

const openNotification = async (item) => {
  try {
    if (!item.read_at) {
      await notifikasiService.markAsRead(item.id);
      // optimistic update
      item.read_at = new Date().toISOString();
    }

    // Try to resolve related pengaduan id from common fields
    const pengaduanId =
      item.pengaduan_id ||
      item.reference_id ||
      (item.data && item.data.pengaduan_id) ||
      (item.meta && item.meta.pengaduan_id);
    if (pengaduanId) {
      // open the detail page (admin detail route exists and shows detail page)
      router.push({ name: "AdminDetail", params: { id: String(pengaduanId) } });
    }
  } catch (err) {
    toast.add({ type: "danger", message: "Gagal membuka notifikasi." });
  }
};

const markAllRead = async () => {
  try {
    const unread = notifications.value.filter((n) => !n.read_at);
    await Promise.all(unread.map((n) => notifikasiService.markAsRead(n.id)));
    unread.forEach((n) => (n.read_at = new Date().toISOString()));
  } catch (err) {
    toast.add({ type: "danger", message: "Gagal menandai notifikasi." });
  }
};

const timeAgo = (iso) => {
  if (!iso) return "";
  try {
    const diff = Math.floor((Date.now() - new Date(iso).getTime()) / 1000);
    if (diff < 60) return "Baru saja";
    if (diff < 3600) return Math.floor(diff / 60) + " menit lalu";
    if (diff < 86400) return Math.floor(diff / 3600) + " jam lalu";
    return Math.floor(diff / 86400) + " hari lalu";
  } catch {
    return "";
  }
};
</script>
