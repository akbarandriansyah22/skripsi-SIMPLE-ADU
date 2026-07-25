<template>
  <AdminSistemLayout>
    <div v-if="error" class="mb-5 flex items-center justify-between rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-xs text-red-700"><span>{{ error }}</span><button class="font-semibold underline" @click="load">Muat ulang</button></div>
    <section class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
      <div v-for="card in cards" :key="card.label" class="rounded-xl border border-slate-200 bg-white p-5 shadow-sm"><p class="text-[11px] font-semibold uppercase tracking-wider text-blue-950">{{ card.label }}</p><p class="mt-3 text-3xl font-bold text-slate-950">{{ loading ? '—' : card.value }}</p></div>
    </section>
    <section class="mt-5 rounded-xl border border-slate-200 bg-white p-5 shadow-sm"><h2 class="text-sm font-bold text-slate-950">Konfigurasi Akses</h2><p class="mt-2 text-xs leading-relaxed text-slate-500">Kelola akun internal dan pastikan setiap pengguna memiliki peran serta unit yang sesuai.</p></section>
  </AdminSistemLayout>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import AdminSistemLayout from '../../layouts/AdminSistemLayout.vue'
import service from '../../services/admin-sistem.service'
import { errorMessage, responseData } from '../../utils/api'

const loading = ref(true); const error = ref(''); const stats = ref({})
const cards = computed(() => [
  { label: 'Total Pengguna', value: stats.value.total_users || 0 },
  { label: 'Total Pengaduan', value: stats.value.total_pengaduan || 0 },
  { label: 'Total Unit', value: stats.value.total_units || 0 },
  { label: 'Total Kategori', value: stats.value.total_categories || 0 },
])
async function load() { loading.value = true; error.value = ''; try { stats.value = responseData(await service.dashboard(), {}) } catch (err) { error.value = errorMessage(err, 'Ringkasan sistem tidak dapat dimuat.') } finally { loading.value = false } }
onMounted(load)
</script>
