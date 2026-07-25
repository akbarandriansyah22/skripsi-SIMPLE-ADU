<template>
  <KasubagLayout>
    <div v-if="error" class="mb-5 flex items-center justify-between rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-xs text-red-700"><span>{{ error }}</span><button class="font-semibold underline" @click="load">Muat ulang</button></div>
    <section class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4"><div v-for="card in cards" :key="card.label" class="rounded-xl border border-slate-200 bg-white p-5 shadow-sm"><p class="text-[11px] font-semibold uppercase tracking-wider text-blue-950">{{ card.label }}</p><p class="mt-3 text-3xl font-bold text-slate-950">{{ loading ? '—' : card.value }}</p></div></section>
    <section class="mt-5 rounded-xl border border-slate-200 bg-white p-5 shadow-sm sm:p-6"><div class="flex items-center justify-between"><div><h2 class="text-sm font-bold text-slate-950">Pengaduan Terbaru</h2><p class="mt-1 text-xs text-slate-500">Pengaduan yang telah didistribusikan ke unit Anda.</p></div><router-link to="/kasubag/pengaduan" class="text-xs font-semibold text-blue-950 hover:underline">Lihat semua</router-link></div><div v-if="loading" class="mt-5 h-24 animate-pulse rounded-lg bg-slate-100"></div><div v-else-if="!items.length" class="mt-5 rounded-lg border border-dashed border-slate-200 px-4 py-8 text-center text-xs text-slate-500">Belum ada pengaduan yang ditugaskan ke unit ini.</div><div v-else class="mt-5 overflow-x-auto rounded-lg border border-slate-200"><table class="min-w-full text-left text-xs"><thead class="bg-slate-50 text-slate-500"><tr><th class="px-4 py-3">Nomor Tiket</th><th class="px-4 py-3">Judul Aduan</th><th class="px-4 py-3">Nama Mahasiswa</th><th class="px-4 py-3">Urgensi</th><th class="px-4 py-3">Status</th><th class="px-4 py-3">Aksi</th></tr></thead><tbody class="divide-y divide-slate-200"><tr v-for="item in items.slice(0, 5)" :key="item.id"><td class="px-4 py-3 font-semibold text-slate-950">{{ ticket(item) }}</td><td class="min-w-48 px-4 py-3">{{ item.judul }}</td><td class="px-4 py-3">{{ item.user?.nama_lengkap || '-' }}</td><td class="px-4 py-3"><StatusBadge kind="urgency" :label="item.urgensi" /></td><td class="px-4 py-3"><StatusBadge :label="item.status" /></td><td class="px-4 py-3"><router-link :to="{ name: 'KasubagDetail', params: { id: item.id } }" class="font-semibold text-blue-950 hover:underline">Detail</router-link></td></tr></tbody></table></div></section>
  </KasubagLayout>
</template>
<script setup>
import { computed, onMounted, ref } from 'vue'
import KasubagLayout from '../../layouts/KasubagLayout.vue'
import StatusBadge from '../../components/StatusBadge.vue'
import service from '../../services/kasubag.service'
import { errorMessage, responseData, responseList } from '../../utils/api'
const items = ref([]); const stats = ref({}); const loading = ref(true); const error = ref('')
const cards = computed(() => [{ label: 'Pengaduan Masuk', value: stats.value.diteruskan || 0 }, { label: 'Sedang Diproses', value: stats.value.diproses || 0 }, { label: 'Menunggu Pembaruan', value: items.value.filter((i) => String(i.status || '').toLowerCase().includes('tunggu')).length }, { label: 'Selesai', value: stats.value.selesai || 0 }])
function ticket(item) { return item.kode_tiket || `ADU-${item.id}` }
async function load() { loading.value = true; error.value = ''; try { const [dashboard, list] = await Promise.all([service.dashboard(), service.getAllPengaduan()]); stats.value = responseData(dashboard, {}); items.value = responseList(list) } catch (err) { error.value = errorMessage(err, 'Dashboard unit tidak dapat dimuat.') } finally { loading.value = false } }
onMounted(load)
</script>
