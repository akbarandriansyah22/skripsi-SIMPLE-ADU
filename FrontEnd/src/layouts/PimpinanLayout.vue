<template>
  <div class="min-h-screen bg-slate-50 text-slate-900">
    <div class="flex min-h-screen">
      <aside
        class="hidden w-80 flex-col border-r border-slate-200 bg-white py-8 px-6 lg:flex shadow-soft"
      >
        <div class="mb-10">
          <p
            class="text-xs uppercase tracking-widest text-blue-600 font-semibold"
          >
            SIMPEL-ADU
          </p>
          <h2 class="mt-3 text-2xl font-bold text-slate-950">Pimpinan</h2>
          <p class="mt-2 text-sm text-slate-600 leading-relaxed">
            Tinjau validasi dan disposisi pengaduan penting.
          </p>
        </div>
        <nav class="space-y-1">
          <router-link
            class="block rounded-lg px-4 py-3 text-sm font-medium text-slate-700 transition hover:bg-slate-100 hover:text-slate-900"
            active-class="bg-blue-950 text-white hover:bg-blue-950 hover:text-white"
            to="/pimpinan/dashboard"
            >Dashboard</router-link
          >
          <router-link
            class="block rounded-lg px-4 py-3 text-sm font-medium text-slate-700 transition hover:bg-slate-100 hover:text-slate-900"
            active-class="bg-blue-950 text-white hover:bg-blue-950 hover:text-white"
            to="/pimpinan/validasi"
            >Validasi Aduan</router-link
          >
          <router-link
            class="block rounded-lg px-4 py-3 text-sm font-medium text-slate-700 transition hover:bg-slate-100 hover:text-slate-900"
            active-class="bg-blue-950 text-white hover:bg-blue-950 hover:text-white"
            to="/pimpinan/disposisi"
            >Disposisi</router-link
          >
          <button
            type="button"
            class="block w-full rounded-lg px-4 py-3 text-left text-sm font-medium text-red-600 transition hover:bg-red-50 hover:text-red-700"
            @click="handleLogout"
          >
            Keluar
          </button>
          >
        </nav>
      </aside>
      <main class="flex-1 p-6 lg:p-10">
        <div
          class="mb-6 flex items-center justify-between rounded-xl bg-white px-5 py-4 shadow-soft ring-1 ring-slate-200/50"
        >
          <div>
            <p
              class="text-xs uppercase tracking-wider text-slate-600 font-medium"
            >
              Rangkuman terbaru untuk keputusan Anda
            </p>
            <h1 class="text-xl font-bold text-slate-950">Dashboard Pimpinan</h1>
          </div>
          <router-link
            class="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-blue-700 active:bg-blue-800 shadow-soft"
            to="/pimpinan/validasi"
          >
            Validasi
          </router-link>
        </div>
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from "vue-router";
import { useAuthStore } from "../stores/auth";

const router = useRouter();
const auth = useAuthStore();

const handleLogout = () => {
  auth.logout();
  router.replace("/auth/login");
};
</script>
