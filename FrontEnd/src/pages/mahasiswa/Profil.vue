<template>
  <MahasiswaLayout>
    <section class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50">
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <p class="text-sm font-bold text-slate-950">Profil Pengguna</p>
          <p class="mt-1 text-sm text-slate-600">Perbarui identitas mahasiswa yang terhubung dengan akun Anda.</p>
        </div>
        <button
          type="submit"
          form="profile-form"
          :disabled="saving || loading"
          class="rounded-lg bg-blue-950 px-4 py-2 text-sm font-semibold text-white transition hover:bg-blue-900 disabled:opacity-60"
        >
          {{ saving ? "Menyimpan..." : "Simpan Perubahan" }}
        </button>
      </div>

      <div v-if="loading" class="mt-6 rounded-lg bg-slate-50 p-5 text-sm text-slate-600">Memuat profil...</div>
      <div v-else-if="error" class="mt-6 rounded-lg bg-red-50 p-5 text-sm text-red-700">Gagal memuat profil: {{ error }}</div>

      <form v-else id="profile-form" class="mt-6 grid gap-5" @submit.prevent="save">
        <div class="grid gap-5 lg:grid-cols-2">
          <div class="space-y-2">
            <label class="text-sm font-semibold text-slate-900">Nama Lengkap</label>
            <input v-model="form.nama_lengkap" type="text" class="w-full rounded-lg border border-slate-200 px-4 py-3 text-sm" />
          </div>
          <div class="space-y-2">
            <label class="text-sm font-semibold text-slate-900">NIM</label>
            <input v-model="form.nim" type="text" disabled class="w-full rounded-lg border border-slate-200 bg-slate-100 px-4 py-3 text-sm text-slate-500" />
          </div>
        </div>
        <div class="grid gap-5 lg:grid-cols-2">
          <div class="space-y-2">
            <label class="text-sm font-semibold text-slate-900">Program Studi</label>
            <input v-model="form.program_studi" type="text" disabled class="w-full rounded-lg border border-slate-200 bg-slate-100 px-4 py-3 text-sm text-slate-500" />
          </div>
          <div class="space-y-2">
            <label class="text-sm font-semibold text-slate-900">Email</label>
            <input :value="email" type="email" disabled class="w-full rounded-lg border border-slate-200 bg-slate-100 px-4 py-3 text-sm text-slate-500" />
          </div>
        </div>
        <div class="space-y-2">
          <label class="text-sm font-semibold text-slate-900">No. HP</label>
          <input v-model="form.no_hp" type="tel" class="w-full rounded-lg border border-slate-200 px-4 py-3 text-sm" />
        </div>
      </form>
    </section>
  </MahasiswaLayout>
</template>

<script setup>
import { onMounted, ref } from "vue";
import MahasiswaLayout from "../../layouts/MahasiswaLayout.vue";
import mahasiswaService from "../../services/mahasiswa.service";
import { useToastStore } from "../../stores/toast";

const toast = useToastStore();
const loading = ref(true);
const saving = ref(false);
const error = ref("");
const email = ref("");
const form = ref({ nama_lengkap: "", nim: "", program_studi: "", no_hp: "" });

const pick = (data, keys, fallback = "") => {
  for (const key of keys) {
    if (data?.[key]) return data[key];
  }
  return fallback;
};

const load = async () => {
  loading.value = true;
  error.value = "";
  try {
    const { data } = await mahasiswaService.profile();
    const profile = data?.data || data || {};
    form.value = {
      nama_lengkap: pick(profile, ["nama_lengkap", "nama", "namaLengkap"]),
      nim: pick(profile, ["nim", "NIM"]),
      program_studi: pick(profile, ["prodi", "program_studi", "programStudi"]),
      no_hp: pick(profile, ["no_hp", "noHP"]),
    };
    email.value = pick(profile, ["email"], profile?.user?.email || "");
  } catch (err) {
    error.value = err?.response?.data?.message || err?.message || "Server error";
  } finally {
    loading.value = false;
  }
};

const save = async () => {
  saving.value = true;
  try {
    await mahasiswaService.updateProfile({
      nama_lengkap: form.value.nama_lengkap,
      no_hp: form.value.no_hp,
    });
    toast.add({ type: "success", message: "Profil berhasil diperbarui." });
  } catch (err) {
    toast.add({ type: "danger", message: "Gagal menyimpan profil." });
  } finally {
    saving.value = false;
  }
};

onMounted(load);
</script>
