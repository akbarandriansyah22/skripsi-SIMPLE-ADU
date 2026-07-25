<template>
  <DashboardLayout
    :role-title="roleTitle"
    :role-description="roleDescription"
    :links="links"
    :header-title="headerTitle"
    header-action-to="/kasubag/pengaduan"
    header-action-label="Lihat Aduan Unit"
  >
    <slot />
  </DashboardLayout>
</template>

<script setup>
import { computed } from 'vue'
import DashboardLayout from '../components/DashboardLayout.vue'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const unit = computed(() => auth.user?.unit_name || '')
const roleTitle = computed(() => unit.value ? `Kasubag ${unit.value}` : 'Kasubag')
const roleDescription = computed(() => unit.value === 'Akademik'
  ? 'Kelola dan tindak lanjuti pengaduan bidang akademik.'
  : 'Kelola dan tindak lanjuti pengaduan sarana dan prasarana.')
const headerTitle = computed(() => unit.value
  ? `Ringkasan Pengaduan ${unit.value}`
  : 'Ringkasan Pengaduan Unit')
const links = [
  { to: '/kasubag/dashboard', label: 'Dashboard' },
  { to: '/kasubag/pengaduan', label: 'Aduan Unit' },
]
</script>
