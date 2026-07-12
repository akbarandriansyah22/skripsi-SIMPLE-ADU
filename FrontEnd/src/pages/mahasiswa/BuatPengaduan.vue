<template>
  <MahasiswaLayout>
    <section
      class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50"
    >
      <div class="mb-6">
        <p class="text-sm font-bold text-slate-950">Kirim Aduan</p>
        <p class="mt-1 text-sm text-slate-600">
          Isi detail aduan Anda dengan jelas dan lengkap.
        </p>
      </div>
      <form class="space-y-6" @submit.prevent="handleSubmit">
        <div class="grid gap-5 lg:grid-cols-2">
          <div class="space-y-2">
            <label class="text-sm font-semibold text-slate-900"
              >Judul Aduan</label
            >
            <input
              v-model="form.judul"
              type="text"
              placeholder="Masukkan judul aduan"
              class="w-full rounded-lg border border-slate-200 bg-white px-4 py-3 text-sm text-slate-900 shadow-soft transition focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500/20"
            />
          </div>
          <div class="space-y-2">
            <label class="text-sm font-semibold text-slate-900">Kategori</label>
            <select
              v-model.number="form.kategori_id"
              class="w-full rounded-lg border border-slate-200 bg-white px-4 py-3 text-sm text-slate-900 shadow-soft transition focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500/20"
            >
              <option v-for="item in kategoriOptions" :key="item.id" :value="item.id">
                {{ item.nama }}
              </option>
            </select>
          </div>
        </div>
        <div class="space-y-2">
          <label class="text-sm font-semibold text-slate-900"
            >Deskripsi Aduan</label
          >
          <textarea
            v-model="form.deskripsi"
            rows="6"
            placeholder="Jelaskan masalah secara rinci"
            class="w-full rounded-lg border border-slate-200 bg-white px-4 py-4 text-sm text-slate-900 shadow-soft transition focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500/20"
          ></textarea>
        </div>
        <div class="grid gap-5 lg:grid-cols-2">
          <div class="space-y-2">
            <label class="text-sm font-semibold text-slate-900"
              >Lampiran Bukti</label
            >
            <input
              type="file"
              @change="onFileChange"
              class="w-full rounded-lg border border-slate-200 bg-white px-3 py-3 text-sm text-slate-700 shadow-soft transition focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500/20"
            />
          </div>
          <div class="rounded-lg border border-blue-100 bg-blue-50 px-4 py-3 text-sm text-blue-800">
            Setelah dikirim, sistem AI akan mengisi kategori prediksi, sentimen, urgensi, dan status pada detail aduan.
          </div>
        </div>
        <button
          type="submit"
          :disabled="submitting"
          class="rounded-lg bg-emerald-600 px-6 py-3 text-sm font-semibold text-white transition hover:bg-emerald-700 active:bg-emerald-800 shadow-soft disabled:opacity-60"
        >
          <span v-if="!submitting">Kirim Aduan</span>
          <span v-else>Mengirim...</span>
        </button>
      </form>
    </section>
  </MahasiswaLayout>
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import MahasiswaLayout from "../../layouts/MahasiswaLayout.vue";
import pengaduanService from "../../services/pengaduan.service";
import { useToastStore } from "../../stores/toast";

const router = useRouter();
const submitting = ref(false);
const toast = useToastStore();
const kategoriOptions = [
  { id: 1, nama: "Sarana Prasarana" },
  { id: 2, nama: "Akademik" },
  { id: 3, nama: "Keuangan" },
  { id: 4, nama: "Layanan Administrasi" },
];

const form = ref({
  judul: "",
  kategori_id: 1,
  deskripsi: "",
  lampiran: null,
});

const handleSubmit = async () => {
  try {
    // basic frontend validation
    if (!form.value.judul.trim() || !form.value.deskripsi.trim()) {
      toast.add({
        type: "danger",
        message: "Judul dan deskripsi wajib diisi.",
      });
      return;
    }
    submitting.value = true;
    const payload = new FormData();
    payload.append("kategori_id", String(form.value.kategori_id));
    payload.append("judul", form.value.judul);
    payload.append("deskripsi", form.value.deskripsi);
    if (form.value.lampiran) payload.append("lampiran", form.value.lampiran);

    await pengaduanService.create(payload);
    // Redirect ke daftar pengaduan dan tampilkan notifikasi sukses melalui query
    router.push({ name: "MahasiswaPengaduan", query: { success: "1" } });
  } catch (err) {
    toast.add({ type: "danger", message: "Gagal mengirim aduan. Coba lagi." });
  } finally {
    submitting.value = false;
  }
};

const onFileChange = (event) => {
  form.value.lampiran = event.target.files?.[0] || null;
};
</script>
