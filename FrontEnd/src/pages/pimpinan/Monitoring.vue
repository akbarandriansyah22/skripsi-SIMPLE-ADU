<template>
  <PimpinanLayout>
    <section class="space-y-5">
      <div>
        <p class="text-[11px] font-semibold uppercase tracking-wider text-blue-950">Pimpinan Fakultas</p>
        <h2 class="mt-1 text-xl font-bold text-slate-950">Monitoring Pengaduan</h2>
        <p class="mt-1 text-xs text-slate-500">Pantau status terkini, disposisi, dan progres aduan setelah diteruskan ke unit sampai selesai.</p>
      </div>

      <section class="grid gap-3 rounded-xl border border-slate-200 bg-white p-4 shadow-sm sm:grid-cols-2 lg:grid-cols-4">
        <select v-model="filters.status" class="rounded-xl border border-slate-200 px-3 py-2.5 text-xs">
          <option value="">Semua status</option>
          <option v-for="status in statuses" :key="status" :value="status">{{ status }}</option>
        </select>
        <select v-model="filters.unit" class="rounded-xl border border-slate-200 px-3 py-2.5 text-xs">
          <option value="">Semua unit</option>
          <option v-for="unit in units" :key="unit" :value="unit">{{ unit }}</option>
        </select>
        <input v-model.trim="filters.q" type="search" placeholder="Cari tiket, judul, mahasiswa" class="rounded-xl border border-slate-200 px-3 py-2.5 text-xs lg:col-span-2" />
      </section>

      <div v-if="error" class="rounded-xl border border-red-200 bg-red-50 p-4 text-xs text-red-700">
        {{ error }}
        <button type="button" class="ml-2 font-semibold underline" @click="load">Coba lagi</button>
      </div>

      <section class="overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm">
        <div class="overflow-x-auto">
          <table class="min-w-[1100px] w-full text-left text-xs">
            <thead class="bg-slate-50 text-[11px] uppercase tracking-wider text-slate-500">
              <tr>
                <th class="px-4 py-3">Nomor Tiket</th>
                <th class="px-4 py-3">Judul</th>
                <th class="px-4 py-3">Mahasiswa</th>
                <th class="px-4 py-3">Unit</th>
                <th class="px-4 py-3">Status terkini</th>
                <th class="px-4 py-3">Perubahan terakhir</th>
                <th class="px-4 py-3">Diperbarui</th>
                <th class="px-4 py-3">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-200">
              <tr v-if="loading"><td colspan="8" class="px-4 py-10 text-center text-slate-500">Memuat monitoring...</td></tr>
              <tr v-else-if="!filtered.length"><td colspan="8" class="px-4 py-10 text-center text-slate-500">Belum ada pengaduan yang dapat dipantau.</td></tr>
              <tr v-for="item in filtered" :key="item.id" class="hover:bg-slate-50/70">
                <td class="px-4 py-3 font-semibold text-blue-950">{{ item.kode_tiket || `ADU-${item.id}` }}</td>
                <td class="max-w-[220px] px-4 py-3 font-semibold text-slate-900">{{ item.judul }}</td>
                <td class="px-4 py-3">
                  <p class="font-semibold">{{ item.user?.nama_lengkap || '-' }}</p>
                  <p class="text-[11px] text-slate-500">{{ item.user?.mahasiswa?.nim || '-' }}</p>
                </td>
                <td class="px-4 py-3">{{ item.unit?.nama_unit || item.disposisi?.unit?.nama_unit || '-' }}</td>
                <td class="px-4 py-3"><StatusBadge :label="item.status" /></td>
                <td class="max-w-[240px] px-4 py-3 text-slate-600">{{ lastHistory(item) }}</td>
                <td class="px-4 py-3 whitespace-nowrap text-slate-500">{{ dateLabel(item.updated_at, true) }}</td>
                <td class="px-4 py-3">
                  <router-link :to="{ name: 'PimpinanMonitoringDetail', params: { id: item.id } }" class="font-semibold text-blue-950 hover:underline">Detail progres</router-link>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </section>
  </PimpinanLayout>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import PimpinanLayout from '../../layouts/PimpinanLayout.vue'
import StatusBadge from '../../components/StatusBadge.vue'
import service from '../../services/pimpinan.service'
import { dateLabel, errorMessage, responseData } from '../../utils/api'

const items = ref([])
const loading = ref(true)
const error = ref('')
const statuses = ['Menunggu Disposisi', 'Diteruskan ke Unit', 'Diproses', 'Selesai']
const filters = reactive({ q: '', status: '', unit: '' })

const units = computed(() => {
  const names = new Set()
  items.value.forEach((item) => {
    const name = item.unit?.nama_unit || item.disposisi?.unit?.nama_unit
    if (name) names.add(name)
  })
  return [...names]
})

const filtered = computed(() => {
  const q = filters.q.toLowerCase()
  return items.value.filter((item) => {
    if (filters.status && item.status !== filters.status) return false
    const unit = item.unit?.nama_unit || item.disposisi?.unit?.nama_unit || ''
    if (filters.unit && unit !== filters.unit) return false
    if (!q) return true
    const haystack = [item.kode_tiket, item.judul, item.user?.nama_lengkap, item.user?.mahasiswa?.nim].join(' ').toLowerCase()
    return haystack.includes(q)
  })
})

function lastHistory(item) {
  const history = [...(item.riwayat_status_pengaduan || [])].sort((a, b) => new Date(a.created_at || 0) - new Date(b.created_at || 0))
  const last = history[history.length - 1]
  if (!last) return '-'
  return `${last.status_lama || 'Awal'} → ${last.status_baru}${last.catatan ? ` (${last.catatan})` : ''}`
}

async function load() {
  loading.value = true
  error.value = ''
  try {
    items.value = responseData(await service.monitoring(), []) || []
  } catch (err) {
    error.value = errorMessage(err, 'Data monitoring tidak dapat dimuat.')
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>
