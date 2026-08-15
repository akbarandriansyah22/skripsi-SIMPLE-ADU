<template>
  <AdminSistemLayout>
    <div v-if="error" class="mb-5 flex items-center justify-between gap-4 rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-xs text-red-700"><span>{{ error }}</span><button class="shrink-0 font-semibold underline" @click="load">Muat Ulang</button></div>
    <section class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
      <div v-for="card in cards" :key="card.label" class="rounded-xl border border-slate-200 bg-white p-5 shadow-sm"><p class="text-[11px] font-semibold uppercase tracking-wider text-blue-950">{{ card.label }}</p><p class="mt-3 text-3xl font-bold text-slate-950">{{ loading ? '—' : card.value }}</p></div>
    </section>
    <section class="mt-5 rounded-xl border border-slate-200 bg-white p-5 shadow-sm sm:p-6">
      <div class="flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
        <div><h2 class="text-sm font-bold text-slate-950">Aksi Cepat</h2><p class="mt-1 text-xs leading-relaxed text-slate-500">Kelola akun dan konfigurasi sistem dari satu tempat.</p></div>
        <button type="button" class="text-xs font-semibold text-blue-950 underline" @click="load">Muat Ulang Ringkasan</button>
      </div>
      <div class="mt-5 grid gap-3 sm:grid-cols-3">
        <RouterLink v-for="action in quickActions" :key="action.to" :to="action.to" class="rounded-xl border border-slate-200 p-4 transition hover:border-blue-300 hover:bg-blue-50"><p class="text-sm font-semibold text-slate-950">{{ action.label }}</p><p class="mt-1 text-xs text-slate-500">{{ action.description }}</p></RouterLink>
      </div>
    </section>
  </AdminSistemLayout>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import AdminSistemLayout from '../../layouts/AdminSistemLayout.vue'
import service from '../../services/admin-sistem.service'
import { errorMessage, responseData } from '../../utils/api'

const loading = ref(true); const error = ref(''); const stats = ref({})
const cards = computed(() => [
  { label: 'Total Pengguna', value: stats.value.total_users ?? 0 },
  { label: 'Total Pengaduan', value: stats.value.total_pengaduan ?? 0 },
  { label: 'Total Unit', value: stats.value.total_units ?? 0 },
  { label: 'Total Kategori', value: stats.value.total_categories ?? 0 },
])
const quickActions = [
  { to: '/admin-sistem/pengguna', label: 'Tambah Pengguna', description: 'Buat dan kelola akun internal.' },
  { to: '/admin-sistem/import-mahasiswa', label: 'Import Mahasiswa', description: 'Unggah CSV/XLSX dan buat akun mahasiswa.' },
  { to: '/admin-sistem/unit', label: 'Kelola Unit', description: 'Atur unit untuk kebutuhan sistem.' },
  { to: '/admin-sistem/kategori', label: 'Kelola Kategori', description: 'Atur kategori pengaduan.' },
]
async function load() { loading.value = true; error.value = ''; try { stats.value = responseData(await service.dashboard(), {}) } catch (err) { error.value = errorMessage(err, 'Ringkasan sistem tidak dapat dimuat.') } finally { loading.value = false } }
onMounted(load)
</script>
