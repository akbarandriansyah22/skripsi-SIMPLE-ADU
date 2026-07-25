<template>
  <AdminSistemLayout>
    <section class="rounded-xl border border-slate-200 bg-white p-5 shadow-sm sm:p-6">
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"><div><h2 class="text-sm font-bold text-slate-950">Manajemen Pengguna</h2><p class="mt-1 text-xs text-slate-500">Data pengguna diambil dari endpoint Admin Sistem.</p></div><input v-model="query" type="search" placeholder="Cari nama atau email" class="w-full rounded-xl border border-slate-200 px-4 py-2.5 text-xs sm:w-64" /></div>
      <div v-if="error" class="mt-5 rounded-lg border border-red-200 bg-red-50 p-4 text-xs text-red-700">{{ error }} <button class="font-semibold underline" @click="load">Muat ulang</button></div>
      <div class="mt-5 overflow-x-auto rounded-lg border border-slate-200"><table class="min-w-full text-left text-xs"><thead class="bg-slate-50 text-slate-500"><tr><th class="px-4 py-3">Pengguna</th><th class="px-4 py-3">Role</th><th class="px-4 py-3">Unit</th><th class="px-4 py-3">Status</th><th class="px-4 py-3">Aksi</th></tr></thead><tbody class="divide-y divide-slate-200"><tr v-if="loading"><td colspan="5" class="px-4 py-8 text-center text-slate-500">Memuat pengguna...</td></tr><tr v-else-if="!filtered.length"><td colspan="5" class="px-4 py-8 text-center text-slate-500">Belum ada pengguna.</td></tr><tr v-for="user in filtered" v-else :key="user.id"><td class="px-4 py-3"><p class="font-semibold text-slate-950">{{ user.nama_lengkap }}</p><p class="text-slate-500">{{ user.email }}</p></td><td class="px-4 py-3">{{ user.role }}</td><td class="px-4 py-3">{{ user.unit?.nama_unit || '-' }}</td><td class="px-4 py-3">{{ user.is_active === false ? 'Nonaktif' : 'Aktif' }}</td><td class="px-4 py-3"><button class="font-semibold text-blue-950 hover:underline" @click="toggleStatus(user)">{{ user.is_active === false ? 'Aktifkan' : 'Nonaktifkan' }}</button></td></tr></tbody></table></div>
    </section>
  </AdminSistemLayout>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import AdminSistemLayout from '../../layouts/AdminSistemLayout.vue'
import service from '../../services/admin-sistem.service'
import { errorMessage, responseList } from '../../utils/api'
import { useToastStore } from '../../stores/toast'

const users = ref([]); const loading = ref(true); const error = ref(''); const query = ref(''); const toast = useToastStore()
const filtered = computed(() => { const term = query.value.toLowerCase(); return users.value.filter((user) => `${user.nama_lengkap || ''} ${user.email || ''} ${user.role || ''}`.toLowerCase().includes(term)) })
async function load() { loading.value = true; error.value = ''; try { users.value = responseList(await service.users()) } catch (err) { error.value = errorMessage(err, 'Data pengguna tidak dapat dimuat.') } finally { loading.value = false } }
async function toggleStatus(user) { try { await service.setUserStatus(user.id, { is_active: user.is_active === false }); user.is_active = user.is_active === false; toast.add({ type: 'success', message: 'Status pengguna diperbarui.' }) } catch (err) { toast.add({ type: 'danger', message: errorMessage(err, 'Status pengguna gagal diperbarui.') }) } }
onMounted(load)
</script>
