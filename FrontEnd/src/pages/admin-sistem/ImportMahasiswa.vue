<template>
  <AdminSistemLayout>
    <section class="space-y-5">
      <div><p class="text-[11px] font-semibold uppercase tracking-wider text-blue-950">Admin Sistem</p><h2 class="mt-1 text-xl font-bold text-slate-950">Import Mahasiswa</h2><p class="mt-1 text-xs text-slate-500">CSV/XLSX dengan kolom nama lengkap, NIM, email, program studi, dan angkatan.</p></div>
      <section class="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
        <div v-if="error" class="mb-4 rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-xs text-red-700">{{ error }}</div>
        <div v-if="success" class="mb-4 rounded-lg border border-emerald-200 bg-emerald-50 px-3 py-2 text-xs text-emerald-700">{{ success }}</div>
        <div class="flex flex-wrap gap-3"><button type="button" class="rounded-xl border border-slate-200 px-4 py-2.5 text-xs font-semibold text-slate-700" :disabled="downloading" @click="downloadTemplate">{{ downloading ? 'Menyiapkan...' : 'Download Template CSV' }}</button><label class="rounded-xl bg-blue-950 px-4 py-2.5 text-xs font-semibold text-white cursor-pointer hover:bg-blue-900"><input type="file" class="hidden" accept=".csv,.xlsx" @change="selectFile" />{{ file?.name || 'Pilih CSV/XLSX' }}</label><button type="button" class="rounded-xl bg-emerald-600 px-4 py-2.5 text-xs font-semibold text-white disabled:opacity-50" :disabled="!file || uploading" @click="submit">{{ uploading ? 'Mengimpor...' : 'Mulai Import' }}</button></div>
        <p class="mt-3 text-[11px] text-slate-500">Akun berhasil dibuat dengan password sementara dan wajib menggantinya saat login pertama.</p>
      </section>
      <section v-if="result" class="overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm"><div class="grid gap-3 border-b border-slate-200 bg-slate-50 p-4 text-xs sm:grid-cols-3"><div><span class="text-slate-500">Total</span><p class="font-bold text-slate-950">{{ result.total_rows }}</p></div><div><span class="text-slate-500">Berhasil</span><p class="font-bold text-emerald-700">{{ result.success_rows }}</p></div><div><span class="text-slate-500">Gagal</span><p class="font-bold text-red-700">{{ result.failed_rows }}</p></div></div><div class="overflow-x-auto"><table class="min-w-[720px] w-full text-left text-xs"><thead class="text-[11px] uppercase tracking-wider text-slate-500"><tr><th class="px-4 py-3">Baris</th><th class="px-4 py-3">NIM</th><th class="px-4 py-3">Email</th><th class="px-4 py-3">Status</th><th class="px-4 py-3">Keterangan</th></tr></thead><tbody class="divide-y divide-slate-200"><tr v-for="row in result.rows" :key="row.row_number"><td class="px-4 py-3">{{ row.row_number }}</td><td class="px-4 py-3">{{ row.nim }}</td><td class="px-4 py-3">{{ row.email }}</td><td class="px-4 py-3"><span :class="row.status === 'berhasil' ? 'text-emerald-700' : 'text-red-700'" class="font-semibold">{{ row.status }}</span></td><td class="px-4 py-3 text-slate-600">{{ row.reason || (row.temporary_password ? `Password sementara: ${row.temporary_password}` : '-') }}</td></tr></tbody></table></div></section>
    </section>
  </AdminSistemLayout>
</template>
<script setup>
import { ref } from 'vue'
import AdminSistemLayout from '../../layouts/AdminSistemLayout.vue'
import service from '../../services/admin-sistem.service'
import { errorMessage } from '../../utils/api'
const file = ref(null); const result = ref(null); const error = ref(''); const success = ref(''); const uploading = ref(false); const downloading = ref(false)
function selectFile(event) { file.value = event.target.files?.[0] || null; error.value = ''; result.value = null; if (file.value && !/\.(csv|xlsx)$/i.test(file.value.name)) { error.value = 'Pilih file CSV atau XLSX.'; file.value = null } }
async function downloadTemplate() { downloading.value = true; try { const response = await service.downloadImportTemplate(); const url = URL.createObjectURL(response.data); const link = document.createElement('a'); link.href = url; link.download = 'template-mahasiswa.csv'; link.click(); URL.revokeObjectURL(url) } catch (err) { error.value = errorMessage(err, 'Template tidak dapat diunduh.') } finally { downloading.value = false } }
async function submit() { if (!file.value || uploading.value) return; uploading.value = true; error.value = ''; success.value = ''; const payload = new FormData(); payload.append('file', file.value); try { const response = await service.importMahasiswa(payload); result.value = response.data?.data || response.data; success.value = 'Import selesai. Simpan password sementara untuk mahasiswa yang berhasil.' } catch (err) { error.value = errorMessage(err, 'Import gagal.') } finally { uploading.value = false } }
</script>
