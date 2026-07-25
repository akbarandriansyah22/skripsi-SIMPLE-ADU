<template>
  <div class="min-h-screen bg-slate-50 text-slate-900">
    <router-view />
    <ToastContainer />
  </div>
</template>

<script setup>
import { onBeforeUnmount, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ToastContainer from "./components/ToastContainer.vue";
import { useAuthStore } from './stores/auth'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()
const handleExpiredSession = () => {
  auth.logout()
  if (route.name !== 'Login') router.replace('/auth/login')
}
onMounted(() => window.addEventListener('simpel-auth-expired', handleExpiredSession))
onBeforeUnmount(() => window.removeEventListener('simpel-auth-expired', handleExpiredSession))
</script>
