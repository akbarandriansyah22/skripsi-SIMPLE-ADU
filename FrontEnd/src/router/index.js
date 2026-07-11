import { createRouter, createWebHistory } from 'vue-router'
import LoginPage from '../pages/auth/Login.vue'
import RegisterPage from '../pages/auth/Register.vue'
import MahasiswaDashboard from '../pages/mahasiswa/Dashboard.vue'
import MahasiswaPengaduan from '../pages/mahasiswa/PengaduanSaya.vue'
import MahasiswaBuatAduan from '../pages/mahasiswa/BuatPengaduan.vue'
import MahasiswaNotifikasi from '../pages/mahasiswa/Notifikasi.vue'
import AdminDashboard from '../pages/admin/Dashboard.vue'
import AdminPengaduan from '../pages/admin/PengaduanList.vue'
import AdminDetail from '../pages/admin/PengaduanDetail.vue'
import PimpinanDashboard from '../pages/pimpinan/Dashboard.vue'
import PimpinanValidasi from '../pages/pimpinan/ValidasiAduan.vue'
import PimpinanDisposisi from '../pages/pimpinan/Disposisi.vue'

const routes = [
  { path: '/auth/login', name: 'Login', component: LoginPage },
  { path: '/auth/register', name: 'Register', component: RegisterPage },
  { path: '/mahasiswa/dashboard', name: 'MahasiswaDashboard', component: MahasiswaDashboard },
  { path: '/mahasiswa/pengaduan', name: 'MahasiswaPengaduan', component: MahasiswaPengaduan },
  { path: '/mahasiswa/kirim', name: 'MahasiswaBuatAduan', component: MahasiswaBuatAduan },
  { path: '/mahasiswa/notifikasi', name: 'MahasiswaNotifikasi', component: MahasiswaNotifikasi },
  { path: '/admin/dashboard', name: 'AdminDashboard', component: AdminDashboard },
  { path: '/admin/pengaduan', name: 'AdminPengaduan', component: AdminPengaduan },
  { path: '/admin/pengaduan/:id', name: 'AdminDetail', component: AdminDetail, props: true },
  { path: '/pimpinan/dashboard', name: 'PimpinanDashboard', component: PimpinanDashboard },
  { path: '/pimpinan/validasi', name: 'PimpinanValidasi', component: PimpinanValidasi },
  { path: '/pimpinan/disposisi', name: 'PimpinanDisposisi', component: PimpinanDisposisi },
  { path: '/:pathMatch(.*)*', redirect: '/auth/login' },
]

export default createRouter({
  history: createWebHistory(),
  routes,
})
