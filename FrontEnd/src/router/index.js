import { createRouter, createWebHistory } from 'vue-router'
import LoginPage from '../pages/auth/Login.vue'
import RegisterPage from '../pages/auth/Register.vue'
import ForbiddenPage from '../pages/auth/Forbidden.vue'
import ChangePasswordPage from '../pages/auth/ChangePassword.vue'
import MahasiswaDashboard from '../pages/mahasiswa/Dashboard.vue'
import MahasiswaPengaduan from '../pages/mahasiswa/PengaduanSaya.vue'
import MahasiswaBuatAduan from '../pages/mahasiswa/BuatPengaduan.vue'
import MahasiswaDetail from '../pages/mahasiswa/PengaduanDetail.vue'
import MahasiswaNotifikasi from '../pages/mahasiswa/Notifikasi.vue'
import MahasiswaProfil from '../pages/mahasiswa/Profil.vue'
import AdminDashboard from '../pages/admin/Dashboard.vue'
import AdminPengaduan from '../pages/admin/PengaduanList.vue'
import AdminDetail from '../pages/admin/PengaduanDetail.vue'
import AdminPengguna from '../pages/admin/Pengguna.vue'
import AdminNotifikasi from '../pages/admin/Notifikasi.vue'
import PimpinanDashboard from '../pages/pimpinan/Dashboard.vue'
import PimpinanValidasi from '../pages/pimpinan/ValidasiAduan.vue'
import PimpinanDisposisi from '../pages/pimpinan/Disposisi.vue'
import PimpinanRiwayatTinggi from '../pages/pimpinan/RiwayatPengaduanTinggi.vue'
import PimpinanRiwayatTinggiDetail from '../pages/pimpinan/RiwayatPengaduanTinggiDetail.vue'
import KasubagDashboard from '../pages/kasubag/Dashboard.vue'
import KasubagPengaduan from '../pages/kasubag/PengaduanList.vue'
import KasubagDetail from '../pages/kasubag/PengaduanDetail.vue'
import AdminSistemDashboard from '../pages/admin-sistem/Dashboard.vue'
import AdminSistemPengguna from '../pages/admin-sistem/Pengguna.vue'
import AdminSistemImportMahasiswa from '../pages/admin-sistem/ImportMahasiswa.vue'
import AdminSistemStatus from '../pages/admin-sistem/StatusKonfigurasi.vue'
import AdminSistemUnit from '../pages/admin-sistem/Unit.vue'
import AdminSistemKategori from '../pages/admin-sistem/Kategori.vue'
import { useAuthStore } from '../stores/auth'

const routes = [
  { path: '/auth/login', name: 'Login', component: LoginPage, meta: { guest: true } },
  { path: '/auth/register', name: 'Register', component: RegisterPage, meta: { guest: true } },
  { path: '/auth/change-password', name: 'ChangePassword', component: ChangePasswordPage, meta: { requiresAuth: true } },
  { path: '/403', name: 'Forbidden', component: ForbiddenPage, meta: { requiresAuth: true } },
  { path: '/mahasiswa/dashboard', name: 'MahasiswaDashboard', component: MahasiswaDashboard, meta: { requiresAuth: true, role: 'mahasiswa' } },
  { path: '/mahasiswa/pengaduan', name: 'MahasiswaPengaduan', component: MahasiswaPengaduan, meta: { requiresAuth: true, role: 'mahasiswa' } },
  { path: '/mahasiswa/pengaduan/:id', name: 'MahasiswaDetail', component: MahasiswaDetail, props: true, meta: { requiresAuth: true, role: 'mahasiswa' } },
  { path: '/mahasiswa/kirim', name: 'MahasiswaBuatAduan', component: MahasiswaBuatAduan, meta: { requiresAuth: true, role: 'mahasiswa' } },
  { path: '/mahasiswa/notifikasi', name: 'MahasiswaNotifikasi', component: MahasiswaNotifikasi, meta: { requiresAuth: true, role: 'mahasiswa' } },
  { path: '/mahasiswa/profil', name: 'MahasiswaProfil', component: MahasiswaProfil, meta: { requiresAuth: true, role: 'mahasiswa' } },
  { path: '/admin/dashboard', name: 'AdminDashboard', component: AdminDashboard, meta: { requiresAuth: true, role: 'admin_fakultas' } },
  { path: '/admin/pengaduan', name: 'AdminPengaduan', component: AdminPengaduan, meta: { requiresAuth: true, role: 'admin_fakultas' } },
  { path: '/admin/pengaduan/:id', name: 'AdminDetail', component: AdminDetail, props: true, meta: { requiresAuth: true, role: 'admin_fakultas' } },
  { path: '/admin/pengguna', name: 'AdminPengguna', component: AdminPengguna, meta: { requiresAuth: true, role: 'admin_fakultas' } },
  { path: '/admin/notifikasi', name: 'AdminNotifikasi', component: AdminNotifikasi, meta: { requiresAuth: true, role: 'admin_fakultas' } },
  { path: '/admin-sistem/dashboard', name: 'AdminSistemDashboard', component: AdminSistemDashboard, meta: { requiresAuth: true, role: 'admin_sistem' } },
  { path: '/admin-sistem/pengguna', name: 'AdminSistemPengguna', component: AdminSistemPengguna, meta: { requiresAuth: true, role: 'admin_sistem' } },
  { path: '/admin-sistem/import-mahasiswa', name: 'AdminSistemImportMahasiswa', component: AdminSistemImportMahasiswa, meta: { requiresAuth: true, role: 'admin_sistem' } },
  { path: '/admin-sistem/status', name: 'AdminSistemStatus', component: AdminSistemStatus, meta: { requiresAuth: true, role: 'admin_sistem' } },
  { path: '/admin-sistem/unit', name: 'AdminSistemUnit', component: AdminSistemUnit, meta: { requiresAuth: true, role: 'admin_sistem' } },
  { path: '/admin-sistem/kategori', name: 'AdminSistemKategori', component: AdminSistemKategori, meta: { requiresAuth: true, role: 'admin_sistem' } },
  { path: '/pimpinan/dashboard', name: 'PimpinanDashboard', component: PimpinanDashboard, meta: { requiresAuth: true, role: 'pimpinan' } },
  { path: '/pimpinan/validasi', name: 'PimpinanValidasi', component: PimpinanValidasi, meta: { requiresAuth: true, role: 'pimpinan' } },
  { path: '/pimpinan/disposisi', name: 'PimpinanDisposisi', component: PimpinanDisposisi, meta: { requiresAuth: true, role: 'pimpinan' } },
  { path: '/pimpinan/riwayat-tinggi', name: 'PimpinanRiwayatTinggi', component: PimpinanRiwayatTinggi, meta: { requiresAuth: true, role: 'pimpinan' } },
  { path: '/pimpinan/riwayat-tinggi/:id', name: 'PimpinanRiwayatTinggiDetail', component: PimpinanRiwayatTinggiDetail, props: true, meta: { requiresAuth: true, role: 'pimpinan' } },
  { path: '/kasubag/dashboard', name: 'KasubagDashboard', component: KasubagDashboard, meta: { requiresAuth: true, role: 'kasubag' } },
  { path: '/kasubag/pengaduan', name: 'KasubagPengaduan', component: KasubagPengaduan, meta: { requiresAuth: true, role: 'kasubag' } },
  { path: '/kasubag/pengaduan/:id', name: 'KasubagDetail', component: KasubagDetail, props: true, meta: { requiresAuth: true, role: 'kasubag' } },
  { path: '/:pathMatch(.*)*', redirect: '/auth/login' },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

const homeForRole = (auth) => {
  if (auth.isAdminSistem) return '/admin-sistem/dashboard'
  if (auth.isAdminFakultas) return '/admin/dashboard'
  if (auth.isPimpinan) return '/pimpinan/dashboard'
  if (auth.isKasubag) return '/kasubag/dashboard'
  return '/mahasiswa/dashboard'
}

router.beforeEach((to) => {
  const auth = useAuthStore()
  auth.loadFromStorage()

  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return '/auth/login'
  }

  if (to.meta.requiresAuth && !auth.role) {
    auth.logout()
    return '/auth/login'
  }

  if (auth.isAuthenticated && auth.user?.password_must_change && to.path !== '/auth/change-password') return '/auth/change-password'

  if (to.meta.guest && auth.isAuthenticated) {
    return homeForRole(auth)
  }

  const allowed = {
    mahasiswa: auth.isMahasiswa,
    admin_fakultas: auth.isAdminFakultas,
    admin_sistem: auth.isAdminSistem,
    pimpinan: auth.isPimpinan,
    kasubag: auth.isKasubag,
  }
  if (to.meta.role && !allowed[to.meta.role]) return '/403'

  return true
})

export default router
