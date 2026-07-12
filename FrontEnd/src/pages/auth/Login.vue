<template>
  <AuthLayout>
    <div class="space-y-8">
      <div>
        <p
          class="text-xs uppercase tracking-widest text-blue-600 font-semibold"
        >
          Selamat Datang
        </p>
        <h1 class="mt-4 text-3xl font-bold text-slate-950">
          Masuk ke SIMPEL-ADU
        </h1>
        <p class="mt-3 text-sm text-slate-600 leading-relaxed">
          Gunakan email kampus atau NIM untuk melanjutkan ke dashboard Anda.
        </p>
      </div>
      <form
        class="space-y-5 rounded-2xl border border-slate-200 bg-slate-50 p-6 shadow-card"
        @submit.prevent="handleSubmit"
      >
        <div class="space-y-2">
          <label class="text-sm font-semibold text-slate-900"
            >Email atau NIM</label
          >
          <input
            v-model="form.email"
            type="text"
            placeholder="Masukkan Email atau NIM"
            class="w-full rounded-lg border border-slate-200 bg-white px-4 py-3 text-sm text-slate-900 shadow-soft transition focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500/20"
          />
        </div>
        <div class="space-y-2">
          <label class="text-sm font-semibold text-slate-900">Kata Sandi</label>
          <input
            v-model="form.password"
            type="password"
            placeholder="Masukkan kata sandi"
            class="w-full rounded-lg border border-slate-200 bg-white px-4 py-3 text-sm text-slate-900 shadow-soft transition focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500/20"
          />
        </div>
        <div class="flex items-center justify-between text-sm">
          <label class="inline-flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              class="h-4 w-4 rounded border-slate-200 text-blue-600 focus:ring-blue-500"
            />
            <span class="text-slate-700">Ingat Saya</span>
          </label>
          <router-link
            to="/auth/login"
            class="font-medium text-blue-600 hover:text-blue-700 transition"
            >Lupa Password?</router-link
          >
        </div>
        <button
          type="submit"
          class="w-full rounded-lg bg-blue-600 px-4 py-3 text-sm font-semibold text-white transition hover:bg-blue-700 active:bg-blue-800 shadow-soft"
        >
          Login
        </button>
      </form>
      <p class="text-center text-sm text-slate-600">
        Belum memiliki akun?
        <router-link
          to="/auth/register"
          class="font-semibold text-blue-600 hover:text-blue-700 transition"
          >Daftar Sekarang</router-link
        >
      </p>
    </div>
  </AuthLayout>
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import AuthLayout from "../../layouts/AuthLayout.vue";
import { useAuthStore } from "../../stores/auth";
import { useToastStore } from "../../stores/toast";

const router = useRouter();
const auth = useAuthStore();
const form = ref({ email: "", password: "" });

const handleSubmit = async () => {
  // basic validation
  if (!form.value.email.trim() || !form.value.password.trim()) {
    const toast = useToastStore();
    toast.add({
      type: "danger",
      message: "Email/NIM dan kata sandi wajib diisi.",
    });
    return;
  }
  try {
    await auth.login(form.value);
    if (auth.isAdmin) {
      router.push("/admin/dashboard");
    } else if (auth.isPimpinan) {
      router.push("/pimpinan/dashboard");
    } else {
      router.push("/mahasiswa/dashboard");
    }
  } catch (error) {
    const toast = useToastStore();
    toast.add({
      type: "danger",
      message: "Gagal masuk. Periksa kredensial Anda.",
    });
  }
};
</script>
