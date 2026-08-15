<template>
  <PimpinanLayout>
    <section class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50">
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <p class="text-sm font-bold text-slate-950">Notifikasi Pimpinan</p>
          <p class="mt-1 text-sm text-slate-600">Pantau pengaduan urgensi tinggi dan pembaruan disposisi.</p>
        </div>
        <button type="button" @click="markAllRead" class="rounded-lg bg-slate-100 px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-200">Tandai Semua Dibaca</button>
      </div>

      <div class="mt-6 space-y-3">
        <p v-if="loading" class="rounded-lg bg-slate-50 p-5 text-sm text-slate-600">Memuat notifikasi...</p>
        <p v-else-if="error" class="rounded-lg bg-red-50 p-5 text-sm text-red-700">Gagal memuat notifikasi: {{ error }}</p>
        <article
          v-else
          v-for="item in notifications"
          :key="item.id"
          class="cursor-pointer rounded-lg border border-slate-200 bg-slate-50 p-5 hover:bg-slate-100"
          @click="openNotification(item)"
        >
          <div class="flex items-start justify-between gap-3">
            <div>
              <p class="font-semibold text-slate-950">{{ item.judul || item.title || 'Notifikasi' }}</p>
              <p class="mt-2 text-sm text-slate-600">{{ item.isi || item.message || item.body }}</p>
              <p class="mt-2 text-xs text-slate-500">{{ timeAgo(item.created_at) }}</p>
            </div>
            <span v-if="!item.read_at && !item.is_read" class="rounded-full bg-blue-100 px-2 py-1 text-xs font-semibold text-blue-700">Baru</span>
          </div>
        </article>
        <p v-if="!loading && !error && !notifications.length" class="rounded-lg bg-slate-50 p-5 text-sm text-slate-500">Belum ada notifikasi.</p>
      </div>
    </section>
  </PimpinanLayout>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import PimpinanLayout from '../../layouts/PimpinanLayout.vue'
import notifikasiService from '../../services/notifikasi.service'
import { useToastStore } from '../../stores/toast'

const router = useRouter()
const toast = useToastStore()
const notifications = ref([])
const loading = ref(true)
const error = ref('')

async function load() {
  loading.value = true
  error.value = ''
  try {
    const { data } = await notifikasiService.mine()
    notifications.value = data?.data || data || []
  } catch (err) {
    error.value = err?.response?.data?.message || err?.message || 'Server error'
  } finally {
    loading.value = false
  }
}

async function markRead(item) {
  if (item.read_at || item.is_read) return
  await notifikasiService.markAsRead(item.id)
  item.read_at = new Date().toISOString()
  item.is_read = true
}

async function openNotification(item) {
  try {
    await markRead(item)
    const pengaduanId = item.pengaduan_id || item.reference_id
    if (pengaduanId) {
      router.push({ name: 'PimpinanMonitoringDetail', params: { id: String(pengaduanId) } })
    }
  } catch {
    toast.add({ type: 'danger', message: 'Gagal membuka notifikasi.' })
  }
}

async function markAllRead() {
  try {
    await Promise.all(notifications.value.filter((item) => !item.read_at && !item.is_read).map(markRead))
    toast.add({ type: 'success', message: 'Notifikasi ditandai dibaca.' })
  } catch {
    toast.add({ type: 'danger', message: 'Gagal menandai notifikasi.' })
  }
}

function timeAgo(iso) {
  if (!iso) return ''
  const diff = Math.floor((Date.now() - new Date(iso).getTime()) / 1000)
  if (diff < 60) return 'Baru saja'
  if (diff < 3600) return `${Math.floor(diff / 60)} menit lalu`
  if (diff < 86400) return `${Math.floor(diff / 3600)} jam lalu`
  return `${Math.floor(diff / 86400)} hari lalu`
}

onMounted(load)
</script>
