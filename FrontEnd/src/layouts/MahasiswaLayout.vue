<template>
  <div class="min-h-screen bg-slate-50 text-slate-900">
    <div class="flex min-h-screen flex-col lg:flex-row">
      <aside
        class="w-full border-b border-slate-200 bg-white py-6 px-5 lg:w-80 lg:border-b-0 lg:border-r lg:pb-0 shadow-soft lg:shadow-none"
      >
        <div class="mb-8">
          <p
            class="text-xs uppercase tracking-widest text-blue-600 font-semibold"
          >
            SIMPEL-ADU
          </p>
          <h2 class="mt-3 text-2xl font-bold text-slate-950">Mahasiswa</h2>
          <p class="mt-2 text-sm text-slate-600 leading-relaxed">
            Lihat status aduan dan kirim pengaduan baru.
          </p>
        </div>
        <nav class="space-y-1">
          <router-link
            class="block rounded-lg px-4 py-3 text-sm font-medium text-slate-700 transition hover:bg-slate-100 hover:text-slate-900"
            active-class="bg-blue-950 text-white hover:bg-blue-950 hover:text-white"
            to="/mahasiswa/dashboard"
            >Dashboard</router-link
          >
          <router-link
            class="block rounded-lg px-4 py-3 text-sm font-medium text-slate-700 transition hover:bg-slate-100 hover:text-slate-900"
            active-class="bg-blue-950 text-white hover:bg-blue-950 hover:text-white"
            to="/mahasiswa/pengaduan"
            >Pengaduan Saya</router-link
          >
          <router-link
            class="block rounded-lg px-4 py-3 text-sm font-medium text-slate-700 transition hover:bg-slate-100 hover:text-slate-900"
            active-class="bg-blue-950 text-white hover:bg-blue-950 hover:text-white"
            to="/mahasiswa/kirim"
            >Kirim Aduan</router-link
          >
          <router-link
            class="block rounded-lg px-4 py-3 text-sm font-medium text-slate-700 transition hover:bg-slate-100 hover:text-slate-900"
            active-class="bg-blue-950 text-white hover:bg-blue-950 hover:text-white"
            to="/mahasiswa/notifikasi"
            >Notifikasi</router-link
          >
          <router-link
            class="block rounded-lg px-4 py-3 text-sm font-medium text-slate-700 transition hover:bg-slate-100 hover:text-slate-900"
            active-class="bg-blue-950 text-white hover:bg-blue-950 hover:text-white"
            to="/mahasiswa/profil"
            >Profil</router-link
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
          class="mb-6 rounded-xl bg-white px-5 py-4 shadow-soft ring-1 ring-slate-200/50"
        >
          <div
            class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between"
          >
            <div>
              <p
                class="text-xs text-slate-600 uppercase tracking-wider font-medium"
              >
                Halo, selamat datang di dashboard Anda
              </p>
              <h1 class="text-xl font-bold text-slate-950">
                Ringkasan Pengaduan
              </h1>
            </div>
            <router-link
              class="inline-flex items-center rounded-lg bg-emerald-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-emerald-700 active:bg-emerald-800 shadow-soft"
              to="/mahasiswa/kirim"
              >Buat Aduan</router-link
            >
          </div>
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
