<template>
  <DashboardLayout
    role-title="Admin Fakultas"
    role-description="Verifikasi, distribusikan, dan pantau pengaduan mahasiswa."
    :links="links"
    header-title="Ringkasan Pengaduan"
    header-action-to="/admin/pengaduan"
    header-action-label="Lihat Pengaduan"
  >
    <slot />
  </DashboardLayout>
</template>

<script setup>
import DashboardLayout from '../components/DashboardLayout.vue'

const links = [
  { to: '/admin/dashboard', label: 'Dashboard' },
  { to: '/admin/pengaduan', label: 'Daftar Pengaduan' },
  { to: '/admin/notifikasi', label: 'Notifikasi' },
]
</script>
