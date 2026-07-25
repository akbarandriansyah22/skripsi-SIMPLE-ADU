<template>
  <div class="min-h-screen bg-slate-50 text-slate-900">
    <div class="flex min-h-screen">
      <button
        type="button"
        class="fixed left-4 top-4 z-40 inline-flex h-10 w-10 items-center justify-center rounded-lg bg-blue-950 text-white shadow-lg lg:hidden"
        aria-label="Buka menu navigasi"
        @click="mobileOpen = true"
      >
        <span class="text-xl leading-none">☰</span>
      </button>

      <div
        v-if="mobileOpen"
        class="fixed inset-0 z-40 bg-slate-950/40 lg:hidden"
        aria-hidden="true"
        @click="mobileOpen = false"
      ></div>

      <aside
        :class="mobileOpen ? 'translate-x-0' : '-translate-x-full'"
        class="fixed inset-y-0 left-0 z-50 flex w-72 flex-col border-r border-slate-200 bg-white px-5 py-6 shadow-xl transition-transform lg:static lg:z-auto lg:w-[282px] lg:translate-x-0 lg:shadow-none"
      >
        <div class="flex items-start justify-between">
          <div>
            <p class="text-[11px] font-semibold uppercase tracking-[0.22em] text-blue-950">SIMPEL-ADU</p>
            <h2 class="mt-3 text-[17px] font-bold text-slate-950">{{ roleTitle }}</h2>
            <p class="mt-2 max-w-[220px] text-xs leading-relaxed text-slate-500">{{ roleDescription }}</p>
          </div>
          <button
            type="button"
            class="rounded-md p-1 text-slate-400 hover:bg-slate-100 hover:text-slate-700 lg:hidden"
            aria-label="Tutup menu navigasi"
            @click="mobileOpen = false"
          >
            ✕
          </button>
        </div>

        <nav class="mt-9 space-y-1" aria-label="Navigasi utama">
          <router-link
            v-for="link in links"
            :key="link.to"
            :to="link.to"
            class="block rounded-full px-4 py-3 text-xs font-medium text-blue-950 transition hover:bg-slate-100"
            active-class="bg-blue-950 text-white hover:bg-blue-950 hover:text-white"
            @click="mobileOpen = false"
          >
            {{ link.label }}
          </router-link>
          <button
            type="button"
            class="mt-2 block w-full rounded-full px-4 py-3 text-left text-xs font-medium text-red-600 transition hover:bg-red-50 hover:text-red-700"
            @click="logout"
          >
            Keluar
          </button>
        </nav>
      </aside>

      <main class="min-w-0 flex-1 px-4 pb-8 pt-20 sm:px-6 lg:px-9 lg:pt-1">
        <header class="mb-5 rounded-xl border border-slate-200 bg-white px-5 py-3 shadow-sm sm:px-6">
          <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <p class="text-[11px] uppercase tracking-[0.08em] text-blue-950">{{ headerEyebrow }}</p>
              <h1 class="mt-1 text-[17px] font-bold text-slate-950">{{ headerTitle }}</h1>
            </div>
            <router-link
              v-if="headerActionTo && headerActionLabel"
              :to="headerActionTo"
              class="inline-flex w-fit items-center rounded-full bg-emerald-600 px-5 py-2.5 text-xs font-semibold text-white shadow-sm transition hover:bg-emerald-700 focus:outline-none focus:ring-2 focus:ring-emerald-500/30"
            >
              {{ headerActionLabel }}
            </router-link>
          </div>
        </header>
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

defineProps({
  roleTitle: { type: String, required: true },
  roleDescription: { type: String, required: true },
  links: { type: Array, default: () => [] },
  headerEyebrow: { type: String, default: 'Halo, selamat datang di dashboard Anda' },
  headerTitle: { type: String, default: 'Ringkasan Pengaduan' },
  headerActionTo: { type: String, default: '' },
  headerActionLabel: { type: String, default: '' },
})

const mobileOpen = ref(false)
const auth = useAuthStore()
const router = useRouter()

function logout() {
  auth.logout()
  router.replace('/auth/login')
}
</script>
