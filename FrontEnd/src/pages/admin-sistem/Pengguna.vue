<template>
  <AdminSistemLayout>
    <section class="space-y-5">
      <div class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p class="text-[11px] font-semibold uppercase tracking-wider text-blue-950">Admin Sistem</p>
          <h2 class="mt-1 text-xl font-bold text-slate-950">Manajemen Pengguna</h2>
          <p class="mt-1 text-xs text-slate-500">Kelola akun, role, unit Kasubag, status akses, dan password pengguna.</p>
        </div>
        <button type="button" class="inline-flex w-fit items-center rounded-xl bg-blue-950 px-4 py-2.5 text-xs font-semibold text-white shadow-sm transition hover:bg-blue-900" @click="openCreateUser">
          + Tambah Pengguna
        </button>
      </div>

      <div v-if="error" class="flex items-start justify-between gap-4 rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-xs text-red-700">
        <span>{{ error }}</span><button type="button" class="shrink-0 font-semibold underline" @click="load">Muat Ulang</button>
      </div>

      <section class="rounded-xl border border-slate-200 bg-white p-4 shadow-sm sm:p-5">
        <div class="grid gap-3 md:grid-cols-2 xl:grid-cols-5">
          <label class="xl:col-span-2"><span class="mb-1 block text-xs font-semibold text-slate-700">Cari nama atau email</span><input v-model="query" type="search" placeholder="Ketik nama atau email..." class="w-full rounded-xl border border-slate-200 px-3 py-2.5 text-xs" /></label>
          <label><span class="mb-1 block text-xs font-semibold text-slate-700">Role</span><select v-model="roleFilter" class="w-full rounded-xl border border-slate-200 px-3 py-2.5 text-xs"><option value="">Semua role</option><option v-for="role in roleOptions" :key="role.value" :value="role.value">{{ role.label }}</option></select></label>
          <label><span class="mb-1 block text-xs font-semibold text-slate-700">Unit</span><select v-model="unitFilter" class="w-full rounded-xl border border-slate-200 px-3 py-2.5 text-xs"><option value="">Semua unit</option><option v-for="unit in units" :key="unit.id" :value="String(unit.id)">{{ unit.nama_unit }}</option></select></label>
          <label><span class="mb-1 block text-xs font-semibold text-slate-700">Status</span><select v-model="statusFilter" class="w-full rounded-xl border border-slate-200 px-3 py-2.5 text-xs"><option value="">Semua status</option><option value="active">Aktif</option><option value="inactive">Nonaktif</option></select></label>
        </div>
      </section>

      <section class="overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm">
        <div class="overflow-x-auto">
          <table class="min-w-[980px] w-full text-left text-xs">
            <thead class="bg-slate-50 text-[11px] uppercase tracking-wider text-slate-500"><tr><th class="px-4 py-3">Nama Lengkap</th><th class="px-4 py-3">Email</th><th class="px-4 py-3">Role</th><th class="px-4 py-3">Unit</th><th class="px-4 py-3">Status</th><th class="px-4 py-3">Dibuat</th><th class="px-4 py-3">Aksi</th></tr></thead>
            <tbody class="divide-y divide-slate-200">
              <tr v-if="loading"><td colspan="7" class="px-4 py-10 text-center text-slate-500">Memuat pengguna...</td></tr>
              <tr v-else-if="!pagedUsers.length"><td colspan="7" class="px-4 py-10 text-center text-slate-500">{{ users.length ? 'Tidak ada pengguna yang sesuai filter.' : 'Belum ada pengguna.' }}</td></tr>
              <tr v-for="user in pagedUsers" :key="user.id" class="align-top hover:bg-slate-50/70">
                <td class="px-4 py-3 font-semibold text-slate-950">{{ user.nama_lengkap || '-' }}</td>
                <td class="px-4 py-3 text-slate-600">{{ user.email || '-' }}</td>
                <td class="px-4 py-3"><span :class="roleClass(user.role)" class="inline-flex rounded-full px-2.5 py-1 font-semibold">{{ roleLabel(user.role) }}</span></td>
                <td class="px-4 py-3 text-slate-600">{{ user.unit?.nama_unit || user.unit_name || '-' }}</td>
                <td class="px-4 py-3"><span :class="user.is_active === false ? 'bg-red-100 text-red-700' : 'bg-emerald-100 text-emerald-700'" class="inline-flex rounded-full px-2.5 py-1 font-semibold">{{ user.is_active === false ? 'Nonaktif' : 'Aktif' }}</span></td>
                <td class="px-4 py-3 whitespace-nowrap text-slate-500">{{ dateLabel(user.created_at || user.createdAt) }}</td>
                <td class="px-4 py-3"><div class="flex flex-wrap gap-x-3 gap-y-2"><button type="button" class="font-semibold text-blue-950 hover:underline" @click="openEditUser(user)">Edit</button><button type="button" class="font-semibold text-slate-700 hover:underline" @click="openResetPassword(user)">Reset Password</button><button type="button" :disabled="isCurrentUser(user) || statusSavingId === user.id" :title="isCurrentUser(user) ? 'Akun sedang digunakan' : ''" :class="isCurrentUser(user) ? 'cursor-not-allowed text-slate-400' : user.is_active === false ? 'text-emerald-700 hover:underline' : 'text-red-700 hover:underline'" class="font-semibold disabled:no-underline" @click="openStatusConfirm(user)">{{ isCurrentUser(user) ? 'Akun sedang digunakan' : user.is_active === false ? 'Aktifkan' : 'Nonaktifkan' }}</button></div></td>
              </tr>
            </tbody>
          </table>
        </div>
        <div class="flex flex-col gap-3 border-t border-slate-200 px-4 py-3 text-xs text-slate-500 sm:flex-row sm:items-center sm:justify-between"><span>Menampilkan {{ pagedUsers.length }} dari {{ filteredUsers.length }} pengguna</span><div class="flex items-center gap-2"><button type="button" class="rounded-lg border border-slate-200 px-3 py-2 disabled:cursor-not-allowed disabled:opacity-40" :disabled="page === 1" @click="page--">Sebelumnya</button><span>Halaman {{ page }} / {{ pageCount }}</span><button type="button" class="rounded-lg border border-slate-200 px-3 py-2 disabled:cursor-not-allowed disabled:opacity-40" :disabled="page >= pageCount" @click="page++">Berikutnya</button></div></div>
      </section>
    </section>

    <Modal :visible="userModalVisible" :title="userModalMode === 'create' ? 'Tambah Pengguna' : 'Edit Pengguna'" :subtitle="userModalMode === 'create' ? 'Buat akun internal baru dengan role yang sesuai.' : 'Perbarui data akun tanpa mengubah password.'" @close="closeUserModal">
      <div class="max-h-[65vh] overflow-y-auto">
        <div v-if="userDetailLoading" class="py-10 text-center text-xs text-slate-500">Memuat detail pengguna...</div>
        <form v-else class="space-y-4" @submit.prevent="submitUser">
          <div v-if="userModalError" class="rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-xs text-red-700">{{ userModalError }}</div>
          <div class="grid gap-4 sm:grid-cols-2"><label><span class="form-label">Nama Lengkap</span><input v-model.trim="userForm.nama_lengkap" type="text" class="form-input" autocomplete="name" /></label><label><span class="form-label">Email</span><input v-model.trim="userForm.email" type="email" class="form-input" autocomplete="email" /></label></div>
          <div v-if="userModalMode === 'create'" class="grid gap-4 sm:grid-cols-2"><label><span class="form-label">Password</span><input v-model="userForm.password" type="password" class="form-input" autocomplete="new-password" /></label><label><span class="form-label">Konfirmasi Password</span><input v-model="userForm.konfirmasi_password" type="password" class="form-input" autocomplete="new-password" /></label></div>
          <div class="grid gap-4 sm:grid-cols-2"><label><span class="form-label">Role</span><select v-model="userForm.role" class="form-input" @change="onRoleChanged"><option value="">Pilih role</option><option v-for="role in availableRoleOptions" :key="role.value" :value="role.value">{{ role.label }}</option></select></label><label v-if="userForm.role === 'kasubag'"><span class="form-label">Unit Kasubag</span><select v-model="userForm.unit_id" class="form-input"><option value="">Pilih unit</option><option v-for="unit in units" :key="unit.id" :value="String(unit.id)">{{ unit.nama_unit }}</option></select></label></div>
          <div v-if="validationError" class="rounded-lg border border-amber-200 bg-amber-50 px-3 py-2 text-xs text-amber-800">{{ validationError }}</div>
          <div class="flex flex-col-reverse gap-2 border-t border-slate-100 pt-4 sm:flex-row sm:justify-end"><button type="button" class="rounded-xl border border-slate-200 px-4 py-2.5 text-xs font-semibold text-slate-700" :disabled="savingUser" @click="closeUserModal">Batal</button><button type="submit" class="rounded-xl bg-blue-950 px-4 py-2.5 text-xs font-semibold text-white disabled:cursor-not-allowed disabled:opacity-60" :disabled="savingUser">{{ savingUser ? 'Menyimpan...' : 'Simpan' }}</button></div>
        </form>
      </div>
    </Modal>

    <Modal :visible="resetModalVisible" title="Reset Password" :subtitle="`Password akun ${resetTarget?.nama_lengkap || ''} akan diubah.`" @close="closeResetPassword">
      <form class="space-y-4" @submit.prevent="submitResetPassword"><div v-if="resetError" class="rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-xs text-red-700">{{ resetError }}</div><label><span class="form-label">Password Baru</span><input v-model="resetForm.password" type="password" class="form-input" autocomplete="new-password" /></label><label><span class="form-label">Konfirmasi Password Baru</span><input v-model="resetForm.konfirmasi_password" type="password" class="form-input" autocomplete="new-password" /></label><div v-if="resetValidationError" class="rounded-lg border border-amber-200 bg-amber-50 px-3 py-2 text-xs text-amber-800">{{ resetValidationError }}</div><div class="flex flex-col-reverse gap-2 border-t border-slate-100 pt-4 sm:flex-row sm:justify-end"><button type="button" class="rounded-xl border border-slate-200 px-4 py-2.5 text-xs font-semibold text-slate-700" :disabled="resetSaving" @click="closeResetPassword">Batal</button><button type="submit" class="rounded-xl bg-blue-950 px-4 py-2.5 text-xs font-semibold text-white disabled:cursor-not-allowed disabled:opacity-60" :disabled="resetSaving">{{ resetSaving ? 'Menyimpan...' : 'Reset Password' }}</button></div></form>
    </Modal>

    <Modal :visible="statusModalVisible" :title="statusTarget?.is_active === false ? 'Aktifkan Akun' : 'Nonaktifkan Akun'" @close="closeStatusConfirm">
      <div class="space-y-4"><p class="text-sm leading-relaxed text-slate-600">{{ statusTarget?.is_active === false ? 'Akun ini dapat login kembali setelah diaktifkan.' : 'Akun ini tidak dapat login selama berstatus nonaktif.' }}</p><p class="rounded-lg bg-slate-50 px-3 py-2 text-xs font-semibold text-slate-800">{{ statusTarget?.nama_lengkap }} — {{ statusTarget?.email }}</p><div class="flex flex-col-reverse gap-2 border-t border-slate-100 pt-4 sm:flex-row sm:justify-end"><button type="button" class="rounded-xl border border-slate-200 px-4 py-2.5 text-xs font-semibold text-slate-700" :disabled="statusSavingId !== null" @click="closeStatusConfirm">Batal</button><button type="button" class="rounded-xl bg-blue-950 px-4 py-2.5 text-xs font-semibold text-white disabled:cursor-not-allowed disabled:opacity-60" :disabled="statusSavingId !== null" @click="confirmStatusChange">{{ statusSavingId !== null ? 'Menyimpan...' : 'Konfirmasi' }}</button></div></div>
    </Modal>
  </AdminSistemLayout>
</template>

<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import AdminSistemLayout from '../../layouts/AdminSistemLayout.vue'
import Modal from '../../components/Modal.vue'
import service from '../../services/admin-sistem.service'
import { dateLabel, errorMessage, responseData, responseList } from '../../utils/api'
import { useAuthStore } from '../../stores/auth'
import { useToastStore } from '../../stores/toast'

const auth = useAuthStore()
const toast = useToastStore()
const roleOptions = [
  { value: 'admin_sistem', label: 'Admin Sistem' },
  { value: 'admin_fakultas', label: 'Admin Fakultas' },
  { value: 'pimpinan_fakultas', label: 'Pimpinan Fakultas' },
  { value: 'kasubag', label: 'Kasubag' },
  { value: 'mahasiswa', label: 'Mahasiswa' },
]
const availableRoleOptions = computed(() => userModalMode.value === 'create' ? roleOptions.filter((role) => role.value !== 'mahasiswa') : roleOptions)
const users = ref([]); const units = ref([]); const loading = ref(true); const error = ref('')
const query = ref(''); const roleFilter = ref(''); const unitFilter = ref(''); const statusFilter = ref(''); const page = ref(1); const perPage = 10
const userModalVisible = ref(false); const userModalMode = ref('create'); const userDetailLoading = ref(false); const savingUser = ref(false); const userModalError = ref(''); const validationError = ref(''); const editingId = ref(null)
const userForm = reactive({ nama_lengkap: '', email: '', password: '', konfirmasi_password: '', role: '', unit_id: '' })
const resetModalVisible = ref(false); const resetTarget = ref(null); const resetSaving = ref(false); const resetError = ref(''); const resetValidationError = ref(''); const resetForm = reactive({ password: '', konfirmasi_password: '' })
const statusModalVisible = ref(false); const statusTarget = ref(null); const statusSavingId = ref(null)

const filteredUsers = computed(() => {
  const term = query.value.trim().toLowerCase()
  return users.value.filter((user) => {
    const unitId = user.unit?.id ?? user.unit_id
    const searchable = `${user.nama_lengkap || ''} ${user.email || ''} ${user.role || ''} ${user.unit?.nama_unit || ''}`.toLowerCase()
    const roleOk = !roleFilter.value || user.role === roleFilter.value
    const unitOk = !unitFilter.value || String(unitId) === String(unitFilter.value)
    const statusOk = !statusFilter.value || (statusFilter.value === 'inactive' ? user.is_active === false : user.is_active !== false)
    return searchable.includes(term) && roleOk && unitOk && statusOk
  })
})
const pageCount = computed(() => Math.max(1, Math.ceil(filteredUsers.value.length / perPage)))
const pagedUsers = computed(() => filteredUsers.value.slice((page.value - 1) * perPage, page.value * perPage))
const isCurrentUser = (user) => auth.user?.id != null && String(auth.user.id) === String(user.id)

watch([query, roleFilter, unitFilter, statusFilter], () => { page.value = 1 })
watch(pageCount, (value) => { if (page.value > value) page.value = value })

function roleLabel(role) { return roleOptions.find((item) => item.value === role)?.label || role || '-' }
function roleClass(role) {
  if (role === 'admin_sistem') return 'bg-blue-100 text-blue-800'
  if (role === 'admin_fakultas') return 'bg-violet-100 text-violet-800'
  if (role === 'pimpinan_fakultas') return 'bg-amber-100 text-amber-800'
  if (role === 'kasubag') return 'bg-cyan-100 text-cyan-800'
  return 'bg-slate-100 text-slate-700'
}
function clearUserForm() { Object.assign(userForm, { nama_lengkap: '', email: '', password: '', konfirmasi_password: '', role: '', unit_id: '' }) }
function onRoleChanged() { if (userForm.role !== 'kasubag') userForm.unit_id = '' }
function openCreateUser() { userModalMode.value = 'create'; editingId.value = null; clearUserForm(); userModalError.value = ''; validationError.value = ''; userDetailLoading.value = false; userModalVisible.value = true }
async function openEditUser(user) {
  userModalMode.value = 'edit'; editingId.value = user.id; clearUserForm(); userForm.role = user.role || ''; userModalError.value = ''; validationError.value = ''; userDetailLoading.value = true; userModalVisible.value = true
  try { const detail = responseData(await service.user(user.id), {}); Object.assign(userForm, { nama_lengkap: detail.nama_lengkap || '', email: detail.email || '', role: detail.role || '', unit_id: detail.unit_id != null ? String(detail.unit_id) : detail.unit?.id != null ? String(detail.unit.id) : '' }) } catch (err) { userModalError.value = errorMessage(err, 'Detail pengguna gagal dimuat.') } finally { userDetailLoading.value = false }
}
function closeUserModal() { if (!savingUser.value) userModalVisible.value = false }
function validateUser() {
  if (!userForm.nama_lengkap.trim()) return 'Nama lengkap wajib diisi.'
  if (!userForm.email.trim()) return 'Email wajib diisi.'
  if (!/^\S+@\S+\.\S+$/.test(userForm.email.trim())) return 'Format email tidak valid.'
  if (!userForm.role) return 'Role wajib dipilih.'
  if (userForm.role === 'kasubag' && !userForm.unit_id) return 'Unit wajib dipilih untuk Kasubag.'
  if (userModalMode.value === 'create') {
    if (!userForm.password) return 'Password wajib diisi.'
    if (userForm.password.length < 6) return 'Password minimal 6 karakter.'
    if (userForm.password !== userForm.konfirmasi_password) return 'Konfirmasi password harus sama.'
  }
  return ''
}
async function submitUser() {
  if (savingUser.value) return
  validationError.value = validateUser(); if (validationError.value) return
  savingUser.value = true; userModalError.value = ''
  const payload = { nama_lengkap: userForm.nama_lengkap.trim(), email: userForm.email.trim(), role: userForm.role, unit_id: userForm.role === 'kasubag' ? Number(userForm.unit_id) : null }
  if (userModalMode.value === 'create') payload.password = userForm.password
  try { if (userModalMode.value === 'create') await service.createUser(payload); else await service.updateUser(editingId.value, payload); toast.add({ type: 'success', message: userModalMode.value === 'create' ? 'Pengguna berhasil dibuat.' : 'Pengguna berhasil diperbarui.' }); closeUserModal(); await load() } catch (err) { userModalError.value = errorMessage(err, 'Pengguna gagal disimpan.'); toast.add({ type: 'danger', message: userModalError.value }) } finally { savingUser.value = false }
}
function openResetPassword(user) { resetTarget.value = user; resetForm.password = ''; resetForm.konfirmasi_password = ''; resetError.value = ''; resetValidationError.value = ''; resetModalVisible.value = true }
function closeResetPassword() { if (!resetSaving.value) resetModalVisible.value = false }
async function submitResetPassword() {
  if (resetSaving.value) return
  if (!resetForm.password) resetValidationError.value = 'Password baru wajib diisi.'; else if (resetForm.password.length < 6) resetValidationError.value = 'Password minimal 6 karakter.'; else if (resetForm.password !== resetForm.konfirmasi_password) resetValidationError.value = 'Konfirmasi password harus sama.'; else resetValidationError.value = ''
  if (resetValidationError.value) return
  resetSaving.value = true; resetError.value = ''
  try { await service.resetPassword(resetTarget.value.id, { password: resetForm.password }); toast.add({ type: 'success', message: 'Password berhasil direset.' }); closeResetPassword() } catch (err) { resetError.value = errorMessage(err, 'Password gagal direset.'); toast.add({ type: 'danger', message: resetError.value }) } finally { resetSaving.value = false }
}
function openStatusConfirm(user) { if (!isCurrentUser(user)) { statusTarget.value = user; statusModalVisible.value = true } }
function closeStatusConfirm() { if (statusSavingId.value === null) statusModalVisible.value = false }
async function confirmStatusChange() {
  if (!statusTarget.value || statusSavingId.value !== null) return
  const target = statusTarget.value; const nextActive = target.is_active === false; statusSavingId.value = target.id
  try { await service.setUserStatus(target.id, { is_active: nextActive }); const index = users.value.findIndex((user) => String(user.id) === String(target.id)); if (index >= 0) users.value[index] = { ...users.value[index], is_active: nextActive }; toast.add({ type: 'success', message: nextActive ? 'Akun berhasil diaktifkan.' : 'Akun berhasil dinonaktifkan.' }); statusModalVisible.value = false } catch (err) { toast.add({ type: 'danger', message: errorMessage(err, 'Status akun gagal diperbarui.') }) } finally { statusSavingId.value = null }
}
async function load() { loading.value = true; error.value = ''; try { const [userResponse, unitResponse] = await Promise.all([service.users(), service.units()]); users.value = responseList(userResponse); units.value = responseList(unitResponse) } catch (err) { error.value = errorMessage(err, 'Data pengguna tidak dapat dimuat.') } finally { loading.value = false } }
onMounted(load)
</script>

<style scoped>
.form-label { display: block; margin-bottom: 0.35rem; font-size: 0.75rem; font-weight: 600; color: rgb(51 65 85); }
.form-input { width: 100%; border-radius: 0.75rem; border: 1px solid rgb(226 232 240); padding: 0.65rem 0.75rem; font-size: 0.75rem; color: rgb(15 23 42); }
</style>
