<template>
  <AdminSistemLayout>
    <section class="space-y-5">
      <div class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between"><div><p class="text-[11px] font-semibold uppercase tracking-wider text-blue-950">Konfigurasi Sistem</p><h2 class="mt-1 text-xl font-bold text-slate-950">Manajemen Kategori</h2><p class="mt-1 text-xs text-slate-500">Kelola kategori yang digunakan pada pengaduan mahasiswa.</p></div><button type="button" class="inline-flex w-fit items-center rounded-xl bg-blue-950 px-4 py-2.5 text-xs font-semibold text-white hover:bg-blue-900" @click="openCreate">+ Tambah Kategori</button></div>
      <div v-if="error" class="flex items-start justify-between gap-4 rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-xs text-red-700"><span>{{ error }}</span><button type="button" class="shrink-0 font-semibold underline" @click="load">Muat Ulang</button></div>
      <section class="rounded-xl border border-slate-200 bg-white p-4 shadow-sm sm:p-5"><label class="block max-w-md"><span class="mb-1 block text-xs font-semibold text-slate-700">Cari kategori</span><input v-model="query" type="search" placeholder="Cari nama atau deskripsi..." class="w-full rounded-xl border border-slate-200 px-3 py-2.5 text-xs" /></label></section>
      <section class="overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm"><div class="overflow-x-auto"><table class="min-w-[720px] w-full text-left text-xs"><thead class="bg-slate-50 text-[11px] uppercase tracking-wider text-slate-500"><tr><th class="px-4 py-3">Nama Kategori</th><th class="px-4 py-3">Deskripsi</th><th class="px-4 py-3">Aksi</th></tr></thead><tbody class="divide-y divide-slate-200"><tr v-if="loading"><td colspan="3" class="px-4 py-10 text-center text-slate-500">Memuat kategori...</td></tr><tr v-else-if="!filteredCategories.length"><td colspan="3" class="px-4 py-10 text-center text-slate-500">{{ categories.length ? 'Tidak ada kategori yang sesuai pencarian.' : 'Belum ada kategori.' }}</td></tr><tr v-for="category in filteredCategories" :key="category.id" class="hover:bg-slate-50/70"><td class="px-4 py-3 font-semibold text-slate-950">{{ category.nama || '-' }}</td><td class="max-w-xl px-4 py-3 text-slate-600">{{ category.deskripsi || '-' }}</td><td class="px-4 py-3"><button type="button" class="font-semibold text-blue-950 hover:underline" @click="openEdit(category)">Edit</button></td></tr></tbody></table></div></section>
    </section>
    <Modal :visible="modalVisible" :title="editingId ? 'Edit Kategori' : 'Tambah Kategori'" :subtitle="editingId ? 'Perbarui kategori tanpa mengubah ID.' : 'Tambahkan kategori pengaduan baru.'" @close="closeModal"><form class="space-y-4" @submit.prevent="submit"><div v-if="modalError" class="rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-xs text-red-700">{{ modalError }}</div><label><span class="form-label">Nama Kategori</span><input v-model.trim="form.nama" type="text" class="form-input" /></label><label><span class="form-label">Deskripsi <span class="font-normal text-slate-400">(opsional)</span></span><textarea v-model.trim="form.deskripsi" rows="4" class="form-input resize-y"></textarea></label><div v-if="validationError" class="rounded-lg border border-amber-200 bg-amber-50 px-3 py-2 text-xs text-amber-800">{{ validationError }}</div><div class="flex flex-col-reverse gap-2 border-t border-slate-100 pt-4 sm:flex-row sm:justify-end"><button type="button" class="rounded-xl border border-slate-200 px-4 py-2.5 text-xs font-semibold text-slate-700" :disabled="saving" @click="closeModal">Batal</button><button type="submit" class="rounded-xl bg-blue-950 px-4 py-2.5 text-xs font-semibold text-white disabled:opacity-60" :disabled="saving">{{ saving ? 'Menyimpan...' : 'Simpan' }}</button></div></form></Modal>
  </AdminSistemLayout>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import AdminSistemLayout from '../../layouts/AdminSistemLayout.vue'
import Modal from '../../components/Modal.vue'
import service from '../../services/admin-sistem.service'
import { errorMessage, responseList } from '../../utils/api'
import { useToastStore } from '../../stores/toast'

const toast = useToastStore(); const categories = ref([]); const loading = ref(true); const error = ref(''); const query = ref(''); const modalVisible = ref(false); const editingId = ref(null); const saving = ref(false); const modalError = ref(''); const validationError = ref(''); const form = reactive({ nama: '', deskripsi: '' })
const filteredCategories = computed(() => { const term = query.value.trim().toLowerCase(); return categories.value.filter((category) => `${category.nama || ''} ${category.deskripsi || ''}`.toLowerCase().includes(term)) })
function clearForm() { Object.assign(form, { nama: '', deskripsi: '' }) }
function openCreate() { editingId.value = null; clearForm(); modalError.value = ''; validationError.value = ''; modalVisible.value = true }
function openEdit(category) { editingId.value = category.id; Object.assign(form, { nama: category.nama || '', deskripsi: category.deskripsi || '' }); modalError.value = ''; validationError.value = ''; modalVisible.value = true }
function closeModal() { if (!saving.value) modalVisible.value = false }
async function submit() { if (saving.value) return; validationError.value = form.nama.trim() ? '' : 'Nama kategori wajib diisi.'; if (validationError.value) return; saving.value = true; modalError.value = ''; const payload = { nama: form.nama.trim(), deskripsi: form.deskripsi.trim() }; try { if (editingId.value) await service.updateCategory(editingId.value, payload); else await service.createCategory(payload); toast.add({ type: 'success', message: editingId.value ? 'Kategori berhasil diperbarui.' : 'Kategori berhasil ditambahkan.' }); closeModal(); await load() } catch (err) { modalError.value = errorMessage(err, 'Kategori gagal disimpan.'); toast.add({ type: 'danger', message: modalError.value }) } finally { saving.value = false } }
async function load() { loading.value = true; error.value = ''; try { categories.value = responseList(await service.categories()) } catch (err) { error.value = errorMessage(err, 'Data kategori tidak dapat dimuat.') } finally { loading.value = false } }
onMounted(load)
</script>

<style scoped>
.form-label { display: block; margin-bottom: 0.35rem; font-size: 0.75rem; font-weight: 600; color: rgb(51 65 85); }
.form-input { width: 100%; border-radius: 0.75rem; border: 1px solid rgb(226 232 240); padding: 0.65rem 0.75rem; font-size: 0.75rem; color: rgb(15 23 42); }
</style>
