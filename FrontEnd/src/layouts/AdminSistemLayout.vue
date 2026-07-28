<template>
  <DashboardLayout
    role-title="Admin Sistem"
    role-description="Kelola pengguna dan konfigurasi akses sistem."
    :links="links"
    header-title="Ringkasan Sistem"
    header-action-to="/admin-sistem/pengguna"
    header-action-label="Tambah Pengguna"
  >
    <slot />
  </DashboardLayout>
</template>

<script setup>
import DashboardLayout from '../components/DashboardLayout.vue'
const links = [
  { to: '/admin-sistem/dashboard', label: 'Dashboard' },
  { to: '/admin-sistem/pengguna', label: 'Manajemen Pengguna' },
  { to: '/admin-sistem/unit', label: 'Manajemen Unit' },
  { to: '/admin-sistem/kategori', label: 'Manajemen Kategori' },
]
</script>
