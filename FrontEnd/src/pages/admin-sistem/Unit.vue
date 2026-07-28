<template>
  <AdminSistemLayout>
    <section class="space-y-5">
      <div class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between"><div><p class="text-[11px] font-semibold uppercase tracking-wider text-blue-950">Konfigurasi Sistem</p><h2 class="mt-1 text-xl font-bold text-slate-950">Manajemen Unit</h2><p class="mt-1 text-xs text-slate-500">Kelola unit yang dapat dipilih untuk akun Kasubag.</p></div><button type="button" class="inline-flex w-fit items-center rounded-xl bg-blue-950 px-4 py-2.5 text-xs font-semibold text-white hover:bg-blue-900" @click="openCreate">+ Tambah Unit</button></div>
      <div v-if="error" class="flex items-start justify-between gap-4 rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-xs text-red-700"><span>{{ error }}</span><button type="button" class="shrink-0 font-semibold underline" @click="load">Muat Ulang</button></div>
      <section class="rounded-xl border border-slate-200 bg-white p-4 shadow-sm sm:p-5"><label class="block max-w-md"><span class="mb-1 block text-xs font-semibold text-slate-700">Cari unit</span><input v-model="query" type="search" placeholder="Cari nama atau email unit..." class="w-full rounded-xl border border-slate-200 px-3 py-2.5 text-xs" /></label></section>
      <section class="overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm"><div class="overflow-x-auto"><table class="min-w-[620px] w-full text-left text-xs"><thead class="bg-slate-50 text-[11px] uppercase tracking-wider text-slate-500"><tr><th class="px-4 py-3">Nama Unit</th><th class="px-4 py-3">Email Unit</th><th class="px-4 py-3">Aksi</th></tr></thead><tbody class="divide-y divide-slate-200"><tr v-if="loading"><td colspan="3" class="px-4 py-10 text-center text-slate-500">Memuat unit...</td></tr><tr v-else-if="!filteredUnits.length"><td colspan="3" class="px-4 py-10 text-center text-slate-500">{{ units.length ? 'Tidak ada unit yang sesuai pencarian.' : 'Belum ada unit.' }}</td></tr><tr v-for="unit in filteredUnits" :key="unit.id" class="hover:bg-slate-50/70"><td class="px-4 py-3 font-semibold text-slate-950">{{ unit.nama_unit || '-' }}</td><td class="px-4 py-3 text-slate-600">{{ unit.email || '-' }}</td><td class="px-4 py-3"><button type="button" class="font-semibold text-blue-950 hover:underline" @click="openEdit(unit)">Edit</button></td></tr></tbody></table></div></section>
    </section>
    <Modal :visible="modalVisible" :title="editingId ? 'Edit Unit' : 'Tambah Unit'" :subtitle="editingId ? 'Perbarui data unit tanpa mengubah ID.' : 'Tambahkan unit baru untuk konfigurasi sistem.'" @close="closeModal"><form class="space-y-4" @submit.prevent="submit"><div v-if="modalError" class="rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-xs text-red-700">{{ modalError }}</div><label><span class="form-label">Nama Unit</span><input v-model.trim="form.nama_unit" type="text" class="form-input" /></label><label><span class="form-label">Email Unit <span class="font-normal text-slate-400">(opsional)</span></span><input v-model.trim="form.email" type="email" class="form-input" /></label><div v-if="validationError" class="rounded-lg border border-amber-200 bg-amber-50 px-3 py-2 text-xs text-amber-800">{{ validationError }}</div><div class="flex flex-col-reverse gap-2 border-t border-slate-100 pt-4 sm:flex-row sm:justify-end"><button type="button" class="rounded-xl border border-slate-200 px-4 py-2.5 text-xs font-semibold text-slate-700" :disabled="saving" @click="closeModal">Batal</button><button type="submit" class="rounded-xl bg-blue-950 px-4 py-2.5 text-xs font-semibold text-white disabled:opacity-60" :disabled="saving">{{ saving ? 'Menyimpan...' : 'Simpan' }}</button></div></form></Modal>
  </AdminSistemLayout>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import AdminSistemLayout from '../../layouts/AdminSistemLayout.vue'
import Modal from '../../components/Modal.vue'
import service from '../../services/admin-sistem.service'
import { errorMessage, responseList } from '../../utils/api'
import { useToastStore } from '../../stores/toast'

const toast = useToastStore(); const units = ref([]); const loading = ref(true); const error = ref(''); const query = ref(''); const modalVisible = ref(false); const editingId = ref(null); const saving = ref(false); const modalError = ref(''); const validationError = ref(''); const form = reactive({ nama_unit: '', email: '' })
const filteredUnits = computed(() => { const term = query.value.trim().toLowerCase(); return units.value.filter((unit) => `${unit.nama_unit || ''} ${unit.email || ''}`.toLowerCase().includes(term)) })
function clearForm() { Object.assign(form, { nama_unit: '', email: '' }) }
function openCreate() { editingId.value = null; clearForm(); modalError.value = ''; validationError.value = ''; modalVisible.value = true }
function openEdit(unit) { editingId.value = unit.id; Object.assign(form, { nama_unit: unit.nama_unit || '', email: unit.email || '' }); modalError.value = ''; validationError.value = ''; modalVisible.value = true }
function closeModal() { if (!saving.value) modalVisible.value = false }
function validate() { if (!form.nama_unit.trim()) return 'Nama unit wajib diisi.'; if (form.email && !/^\S+@\S+\.\S+$/.test(form.email.trim())) return 'Format email unit tidak valid.'; return '' }
async function submit() { if (saving.value) return; validationError.value = validate(); if (validationError.value) return; saving.value = true; modalError.value = ''; const payload = { nama_unit: form.nama_unit.trim(), email: form.email.trim() }; try { if (editingId.value) await service.updateUnit(editingId.value, payload); else await service.createUnit(payload); toast.add({ type: 'success', message: editingId.value ? 'Unit berhasil diperbarui.' : 'Unit berhasil ditambahkan.' }); closeModal(); await load() } catch (err) { modalError.value = errorMessage(err, 'Unit gagal disimpan.'); toast.add({ type: 'danger', message: modalError.value }) } finally { saving.value = false } }
async function load() { loading.value = true; error.value = ''; try { units.value = responseList(await service.units()) } catch (err) { error.value = errorMessage(err, 'Data unit tidak dapat dimuat.') } finally { loading.value = false } }
onMounted(load)
</script>

<style scoped>
.form-label { display: block; margin-bottom: 0.35rem; font-size: 0.75rem; font-weight: 600; color: rgb(51 65 85); }
.form-input { width: 100%; border-radius: 0.75rem; border: 1px solid rgb(226 232 240); padding: 0.65rem 0.75rem; font-size: 0.75rem; color: rgb(15 23 42); }
</style>
