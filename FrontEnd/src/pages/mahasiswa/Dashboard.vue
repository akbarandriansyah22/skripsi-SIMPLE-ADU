<template>
  <MahasiswaLayout>
    <div v-if="error" class="mb-5 flex items-center justify-between rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-xs text-red-700">
      <span>{{ error }}</span><button type="button" class="font-semibold underline" @click="load">Muat ulang</button>
    </div>
    <section class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
      <div v-for="card in cards" :key="card.label" class="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
        <p class="text-[11px] font-semibold uppercase tracking-wider text-blue-950">{{ card.label }}</p>
        <p class="mt-3 text-3xl font-bold text-slate-950">{{ loading ? '—' : card.value }}</p>
      </div>
    </section>

    <section class="mt-5 rounded-xl border border-slate-200 bg-white p-5 shadow-sm sm:p-6">
      <div class="flex items-center justify-between gap-3">
        <div><h2 class="text-sm font-bold text-slate-950">Pengaduan Terbaru</h2><p class="mt-1 text-xs text-slate-500">Status pengaduan yang Anda kirim.</p></div>
        <router-link to="/mahasiswa/pengaduan" class="text-xs font-semibold text-blue-950 hover:underline">Lihat semua</router-link>
      </div>
      <div v-if="loading" class="mt-5 space-y-3"><div v-for="n in 3" :key="n" class="h-14 animate-pulse rounded-lg bg-slate-100"></div></div>
      <div v-else-if="!recent.length" class="mt-5 rounded-lg border border-dashed border-slate-200 px-4 py-8 text-center text-xs text-slate-500">Belum ada pengaduan.</div>
      <div v-else class="mt-5 overflow-x-auto rounded-lg border border-slate-200">
        <table class="min-w-full text-left text-xs"><thead class="bg-slate-50 text-slate-500"><tr><th class="px-4 py-3 font-semibold">Nomor Tiket</th><th class="px-4 py-3 font-semibold">Judul</th><th class="px-4 py-3 font-semibold">Status</th><th class="px-4 py-3 font-semibold">Tanggal</th><th class="px-4 py-3"></th></tr></thead>
          <tbody class="divide-y divide-slate-200"><tr v-for="item in recent" :key="item.id" class="text-slate-700"><td class="whitespace-nowrap px-4 py-3 font-semibold text-slate-950">{{ ticket(item) }}</td><td class="min-w-56 px-4 py-3">{{ item.judul }}</td><td class="px-4 py-3"><StatusBadge :label="item.status" /></td><td class="whitespace-nowrap px-4 py-3">{{ dateLabel(item.created_at) }}</td><td class="px-4 py-3"><router-link class="font-semibold text-blue-950 hover:underline" :to="{ name: 'MahasiswaDetail', params: { id: item.id } }">Detail</router-link></td></tr></tbody>
        </table>
      </div>
    </section>
  </MahasiswaLayout>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import MahasiswaLayout from '../../layouts/MahasiswaLayout.vue'
import StatusBadge from '../../components/StatusBadge.vue'
import pengaduanService from '../../services/pengaduan.service'
import { dateLabel, errorMessage, responseList } from '../../utils/api'

const loading = ref(true)
const error = ref('')
const complaints = ref([])
const cards = computed(() => [
  { label: 'Total Pengaduan', value: complaints.value.length },
  { label: 'Menunggu Verifikasi', value: countBy('menunggu verifikasi') },
  { label: 'Sedang Diproses', value: countBy('proses') },
  { label: 'Selesai', value: countBy('selesai') },
])
const recent = computed(() => [...complaints.value].sort((a, b) => new Date(b.created_at || 0) - new Date(a.created_at || 0)).slice(0, 5))

function countBy(term) { return complaints.value.filter((item) => String(item.status || '').toLowerCase().includes(term)).length }
function ticket(item) { return item.kode_tiket || `ADU-${item.id}` }

async function load() {
  loading.value = true; error.value = ''
  try {
    complaints.value = responseList(await pengaduanService.myPengaduan())
  } catch (err) {
    error.value = errorMessage(err, 'Dashboard tidak dapat dimuat.')
  } finally { loading.value = false }
}

onMounted(load)
</script>
