<template>
  <DashboardLayout
    role-title="Mahasiswa"
    role-description="Lihat status aduan dan kirim pengaduan baru."
    :links="links"
    header-title="Ringkasan Pengaduan"
    header-action-to="/mahasiswa/kirim"
    header-action-label="Buat Aduan"
  >
    <slot />
  </DashboardLayout>
</template>

<script setup>
import DashboardLayout from '../components/DashboardLayout.vue'

const links = [
  { to: '/mahasiswa/dashboard', label: 'Dashboard' },
  { to: '/mahasiswa/pengaduan', label: 'Pengaduan Saya' },
  { to: '/mahasiswa/kirim', label: 'Kirim Aduan' },
  { to: '/mahasiswa/notifikasi', label: 'Notifikasi' },
  { to: '/mahasiswa/profil', label: 'Profil' },
]
</script>
