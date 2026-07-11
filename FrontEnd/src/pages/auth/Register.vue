<template>
  <AuthLayout>
    <div class="space-y-8">
      <div>
        <p
          class="text-xs uppercase tracking-widest text-blue-600 font-semibold"
        >
          Buat Akun
        </p>
        <h1 class="mt-4 text-3xl font-bold text-slate-950">
          Daftar SIMPEL-ADU
        </h1>
        <p class="mt-3 text-sm text-slate-600 leading-relaxed">
          Gunakan email resmi kampus untuk verifikasi dan akses dashboard.
        </p>
      </div>
      <form
        class="space-y-6 rounded-2xl border border-slate-200 bg-slate-50 p-6 shadow-card"
        @submit.prevent="handleSubmit"
      >
        <div class="grid gap-5 sm:grid-cols-2">
          <div class="space-y-2">
            <label class="text-sm font-semibold text-slate-900"
              >Nama Lengkap</label
            >
            <input
              v-model="form.nama"
              type="text"
              placeholder="Masukkan nama lengkap"
              class="w-full rounded-lg border border-slate-200 bg-white px-4 py-3 text-sm text-slate-900 shadow-soft transition focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500/20"
            />
          </div>
          <div class="space-y-2">
            <label class="text-sm font-semibold text-slate-900">NIM</label>
            <input
              v-model="form.nim"
              type="text"
              placeholder="Nomor Induk Mahasiswa"
              class="w-full rounded-lg border border-slate-200 bg-white px-4 py-3 text-sm text-slate-900 shadow-soft transition focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500/20"
            />
          </div>
        </div>
        <div class="grid gap-5 sm:grid-cols-2">
          <div class="space-y-2">
            <label class="text-sm font-semibold text-slate-900"
              >Program Studi</label
            >
            <input
              v-model="form.prodi"
              type="text"
              placeholder="Pilih Program Studi"
              class="w-full rounded-lg border border-slate-200 bg-white px-4 py-3 text-sm text-slate-900 shadow-soft transition focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500/20"
            />
          </div>
          <div class="space-y-2">
            <label class="text-sm font-semibold text-slate-900"
              >Email Kampus</label
            >
            <input
              v-model="form.email"
              type="email"
              placeholder="nim@mahasiswa.umj.ac.id"
              class="w-full rounded-lg border border-slate-200 bg-white px-4 py-3 text-sm text-slate-900 shadow-soft transition focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500/20"
            />
          </div>
        </div>
        <div class="grid gap-5 sm:grid-cols-2">
          <div class="space-y-2">
            <label class="text-sm font-semibold text-slate-900">Password</label>
            <input
              v-model="form.password"
              type="password"
              placeholder="Masukkan password"
              class="w-full rounded-lg border border-slate-200 bg-white px-4 py-3 text-sm text-slate-900 shadow-soft transition focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500/20"
            />
          </div>
          <div class="space-y-2">
            <label class="text-sm font-semibold text-slate-900"
              >Konfirmasi Password</label
            >
            <input
              v-model="form.confirmPassword"
              type="password"
              placeholder="Konfirmasi password"
              class="w-full rounded-lg border border-slate-200 bg-white px-4 py-3 text-sm text-slate-900 shadow-soft transition focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500/20"
            />
          </div>
        </div>
        <button
          type="submit"
          class="w-full rounded-lg bg-blue-600 px-4 py-3 text-sm font-semibold text-white transition hover:bg-blue-700 active:bg-blue-800 shadow-soft"
        >
          Daftar Sekarang
        </button>
      </form>
      <div
        class="grid gap-3 rounded-2xl border border-slate-200 bg-white p-6 shadow-soft"
      >
        <div class="space-y-2">
          <p class="text-sm font-semibold text-slate-950">Aman & Rahasia</p>
          <p class="text-sm text-slate-600 leading-relaxed">
            Identitas Anda dilindungi dalam setiap laporan yang disampaikan.
          </p>
        </div>
      </div>
      <p class="text-center text-sm text-slate-600">
        Sudah memiliki akun?
        <router-link
          to="/auth/login"
          class="font-semibold text-blue-600 hover:text-blue-700 transition"
          >Masuk di sini</router-link
        >
      </p>
    </div>
  </AuthLayout>
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import authService from "../../services/auth.service";
import { useToastStore } from "../../stores/toast";

const router = useRouter();
const form = ref({
  nama: "",
  nim: "",
  prodi: "",
  email: "",
  password: "",
  confirmPassword: "",
});

const handleSubmit = async () => {
  if (form.value.password !== form.value.confirmPassword) {
    const toast = useToastStore();
    toast.add({
      type: "danger",
      message: "Password dan konfirmasi tidak cocok.",
    });
    return;
  }
  try {
    await authService.register({
      nama: form.value.nama,
      nim: form.value.nim,
      prodi: form.value.prodi,
      email: form.value.email,
      password: form.value.password,
    });
    router.push("/auth/login");
  } catch (error) {
    const toast = useToastStore();
    toast.add({ type: "danger", message: "Gagal mendaftar. Coba lagi." });
  }
};
</script>
