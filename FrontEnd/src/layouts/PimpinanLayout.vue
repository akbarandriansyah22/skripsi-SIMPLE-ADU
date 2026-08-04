<template>
  <DashboardLayout
    role-title="Pimpinan Fakultas"
    role-description="Pantau pengaduan dan berikan disposisi untuk laporan prioritas."
    :links="links"
    header-title="Ringkasan Pengaduan"
    header-action-to="/pimpinan/validasi"
    header-action-label="Pengaduan Urgensi Tinggi"
  >
    <slot />
  </DashboardLayout>
</template>

<script setup>
import DashboardLayout from '../components/DashboardLayout.vue'

const links = [
  { to: '/pimpinan/dashboard', label: 'Dashboard' },
  { to: '/pimpinan/validasi', label: 'Pengaduan Urgensi Tinggi' },
  { to: '/pimpinan/riwayat-tinggi', label: 'Riwayat Pengaduan Tinggi' },
  { to: '/pimpinan/disposisi', label: 'Disposisi' },
]
</script>
