<template>
  <AuthLayout>
    <div class="text-center"><p class="text-xs font-semibold uppercase tracking-widest text-red-600">Akses ditolak</p><h1 class="mt-3 text-2xl font-bold text-slate-950">Anda tidak memiliki izin</h1><p class="mt-3 text-sm leading-relaxed text-slate-600">Halaman ini hanya dapat diakses oleh role yang berwenang.</p><button class="mt-6 rounded-full bg-blue-950 px-5 py-3 text-xs font-semibold text-white" @click="goHome">Kembali ke dashboard</button></div>
  </AuthLayout>
</template>
<script setup>
import { useRouter } from 'vue-router'
import AuthLayout from '../../layouts/AuthLayout.vue'
import { useAuthStore } from '../../stores/auth'
const router = useRouter(); const auth = useAuthStore()
function goHome() { if (auth.isAdminSistem) return router.replace('/admin-sistem/dashboard'); if (auth.isAdminFakultas) return router.replace('/admin/dashboard'); if (auth.isPimpinan) return router.replace('/pimpinan/dashboard'); if (auth.isKasubag) return router.replace('/kasubag/dashboard'); return router.replace('/mahasiswa/dashboard') }
</script>
