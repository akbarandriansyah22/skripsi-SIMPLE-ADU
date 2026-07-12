import { createRouter, createWebHistory } from 'vue-router'
import LoginPage from '../pages/auth/Login.vue'
import RegisterPage from '../pages/auth/Register.vue'
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
import { useAuthStore } from '../stores/auth'

const routes = [
  { path: '/auth/login', name: 'Login', component: LoginPage, meta: { guest: true } },
  { path: '/auth/register', name: 'Register', component: RegisterPage, meta: { guest: true } },
  { path: '/mahasiswa/dashboard', name: 'MahasiswaDashboard', component: MahasiswaDashboard, meta: { requiresAuth: true, role: 'mahasiswa' } },
  { path: '/mahasiswa/pengaduan', name: 'MahasiswaPengaduan', component: MahasiswaPengaduan, meta: { requiresAuth: true, role: 'mahasiswa' } },
  { path: '/mahasiswa/pengaduan/:id', name: 'MahasiswaDetail', component: MahasiswaDetail, props: true, meta: { requiresAuth: true, role: 'mahasiswa' } },
  { path: '/mahasiswa/kirim', name: 'MahasiswaBuatAduan', component: MahasiswaBuatAduan, meta: { requiresAuth: true, role: 'mahasiswa' } },
  { path: '/mahasiswa/notifikasi', name: 'MahasiswaNotifikasi', component: MahasiswaNotifikasi, meta: { requiresAuth: true, role: 'mahasiswa' } },
  { path: '/mahasiswa/profil', name: 'MahasiswaProfil', component: MahasiswaProfil, meta: { requiresAuth: true, role: 'mahasiswa' } },
  { path: '/admin/dashboard', name: 'AdminDashboard', component: AdminDashboard, meta: { requiresAuth: true, role: 'admin' } },
  { path: '/admin/pengaduan', name: 'AdminPengaduan', component: AdminPengaduan, meta: { requiresAuth: true, role: 'admin' } },
  { path: '/admin/pengaduan/:id', name: 'AdminDetail', component: AdminDetail, props: true, meta: { requiresAuth: true, role: 'admin' } },
  { path: '/admin/pengguna', name: 'AdminPengguna', component: AdminPengguna, meta: { requiresAuth: true, role: 'admin' } },
  { path: '/admin/notifikasi', name: 'AdminNotifikasi', component: AdminNotifikasi, meta: { requiresAuth: true, role: 'admin' } },
  { path: '/pimpinan/dashboard', name: 'PimpinanDashboard', component: PimpinanDashboard, meta: { requiresAuth: true, role: 'pimpinan' } },
  { path: '/pimpinan/validasi', name: 'PimpinanValidasi', component: PimpinanValidasi, meta: { requiresAuth: true, role: 'pimpinan' } },
  { path: '/pimpinan/disposisi', name: 'PimpinanDisposisi', component: PimpinanDisposisi, meta: { requiresAuth: true, role: 'pimpinan' } },
  { path: '/:pathMatch(.*)*', redirect: '/auth/login' },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

const homeForRole = (auth) => {
  if (auth.isAdmin) return '/admin/dashboard'
  if (auth.isPimpinan) return '/pimpinan/dashboard'
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

  if (to.meta.guest && auth.isAuthenticated) {
    return homeForRole(auth)
  }

  if (to.meta.role === 'mahasiswa' && !auth.isMahasiswa) return homeForRole(auth)
  if (to.meta.role === 'admin' && !auth.isAdmin) return homeForRole(auth)
  if (to.meta.role === 'pimpinan' && !auth.isPimpinan) return homeForRole(auth)

  return true
})

export default router
