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
            <select
              v-model="form.program_studi"
              class="w-full rounded-lg border border-slate-200 bg-white px-4 py-3 text-sm text-slate-900 shadow-soft transition focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500/20"
            >
              <option value="">Pilih Program Studi</option>
              <option>Teknik Informatika</option>
              <option>Teknik Industri</option>
              <option>Teknik Mesin</option>
              <option>Teknik Elektro</option>
              <option>Teknik Sipil</option>
              <option>Arsitektur</option>
            </select>
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
            <label class="text-sm font-semibold text-slate-900">Angkatan</label>
            <input
              v-model.number="form.angkatan"
              type="number"
              min="2000"
              max="2100"
              placeholder="2023"
              class="w-full rounded-lg border border-slate-200 bg-white px-4 py-3 text-sm text-slate-900 shadow-soft transition focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500/20"
            />
          </div>
          <div class="space-y-2">
            <label class="text-sm font-semibold text-slate-900">No. HP</label>
            <input
              v-model="form.no_hp"
              type="tel"
              placeholder="08xxxxxxxxxx"
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
        <p v-if="error" class="rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-xs text-red-700">{{ error }}</p>
        <button
          type="submit"
          :disabled="submitting"
          class="w-full rounded-lg bg-blue-600 px-4 py-3 text-sm font-semibold text-white transition hover:bg-blue-700 active:bg-blue-800 shadow-soft"
        >
          {{ submitting ? "Mendaftarkan..." : "Daftar Sekarang" }}
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
import AuthLayout from "../../layouts/AuthLayout.vue";
import authService from "../../services/auth.service";
import { useToastStore } from "../../stores/toast";
import { errorMessage } from "../../utils/api";

const router = useRouter();
const submitting = ref(false);
const error = ref("");
const form = ref({
  nama: "",
  nim: "",
  program_studi: "",
  angkatan: new Date().getFullYear(),
  no_hp: "",
  email: "",
  password: "",
  confirmPassword: "",
});

const handleSubmit = async () => {
  const toast = useToastStore();
  if (
    !form.value.nama.trim() ||
    !form.value.nim.trim() ||
    !form.value.program_studi.trim() ||
    !form.value.email.trim() ||
    !form.value.password.trim()
  ) {
    error.value = "Lengkapi seluruh field wajib.";
    return;
  }
  if (form.value.password !== form.value.confirmPassword) {
    error.value = "Password dan konfirmasi tidak cocok.";
    return;
  }
  if (form.value.password.length < 6) {
    error.value = "Password minimal 6 karakter.";
    return;
  }
  error.value = "";
  submitting.value = true;
  try {
    await authService.register({
      nama_lengkap: form.value.nama,
      nim: form.value.nim,
      program_studi: form.value.program_studi,
      angkatan: Number(form.value.angkatan),
      no_hp: form.value.no_hp,
      email: form.value.email,
      password: form.value.password,
    });
    toast.add({ type: "success", message: "Registrasi berhasil. Silakan masuk." });
    router.push("/auth/login");
  } catch (err) {
    error.value = errorMessage(err, "Gagal mendaftar. Coba lagi.");
    toast.add({ type: "danger", message: error.value });
  } finally {
    submitting.value = false;
  }
};
</script>
