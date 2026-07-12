<template>
  <PimpinanLayout>
    <section
      class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50"
    >
      <div class="grid gap-8 lg:grid-cols-[1.2fr_0.8fr]">
        <div>
          <p
            class="text-xs uppercase tracking-wider text-blue-600 font-semibold"
          >
            Validasi & Tindak Lanjut Aduan
          </p>
          <h2 class="mt-4 text-2xl font-bold text-slate-950">
            Tinjau aduan yang perlu keputusan resmi Anda.
          </h2>
          <div class="mt-6">
            <template v-if="loading">
              <p class="text-sm text-slate-600">
                Memuat daftar aduan urgensi tinggi...
              </p>
            </template>
            <template v-else-if="error">
              <p class="text-sm text-red-600">{{ error }}</p>
            </template>
            <template v-else-if="urgent.length">
              <div class="flex flex-wrap gap-2">
                <button
                  v-for="item in urgent"
                  :key="item.id"
                  @click="selectReport(item)"
                  :class="[
                    'rounded-full px-4 py-2 text-sm font-semibold transition',
                    selected?.id === item.id
                      ? 'bg-blue-600 text-white'
                      : 'bg-slate-100 text-slate-700 hover:bg-slate-200',
                  ]"
                >
                  {{ item.kode_tiket || item.kode || "ADU-" + item.id }}
                </button>
              </div>
            </template>
            <template v-else>
              <p class="text-sm text-slate-600">
                Tidak ada aduan urgensi tinggi.
              </p>
            </template>
          </div>
          <div class="mt-6 rounded-lg border border-slate-200 bg-slate-50 p-5">
            <p class="text-sm font-bold text-slate-950">
              ID:
              {{
                selected?.kode_tiket || selected?.kode || (selected?.id ? "ADU-" + selected.id : "-")
              }}
            </p>
            <p class="mt-3 text-slate-700 leading-relaxed">
              {{
                selected?.deskripsi ||
                selected?.judul ||
                "Pilih pengaduan untuk ditinjau."
              }}
            </p>
            <p class="mt-3 text-xs text-slate-500">
              {{ selected?.created_at || selected?.tanggal || "" }}
            </p>
          </div>
          <div class="mt-6 space-y-4">
            <label class="block text-sm font-semibold text-slate-900"
              >Jawaban Resmi / Keputusan Tindak Lanjut</label
            >
            <textarea
              rows="6"
              v-model="answer"
              class="w-full rounded-lg border border-slate-200 bg-white px-4 py-4 text-sm text-slate-900 shadow-soft transition focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500/20"
              placeholder="Ketik jawaban resmi yang akan disampaikan ke pelapor..."
            ></textarea>
          </div>
        </div>
        <div class="space-y-4">
          <div
            class="rounded-xl bg-white p-5 shadow-soft ring-1 ring-slate-200"
          >
            <p class="text-sm font-bold text-slate-950">Pengirim</p>
            <p class="mt-3 text-sm text-slate-700">
              {{
                selected?.pelapor_name ||
                selected?.user?.nama_lengkap ||
                selected?.mahasiswa?.nama ||
                "Anonim (Mahasiswa FT)"
              }}
            </p>
            <p class="mt-2 text-xs text-slate-500">
              {{
                selected?.unit?.nama_unit || selected?.unit || selected?.kategori?.nama || selected?.kategori || "Terverifikasi Admin"
              }}
            </p>
          </div>
          <div
            class="rounded-xl bg-white p-5 shadow-soft ring-1 ring-slate-200"
          >
            <button
              @click="handleReturn()"
              :disabled="!selected || !answer.trim()"
              class="w-full rounded-lg px-4 py-3 text-sm font-semibold transition disabled:cursor-not-allowed disabled:bg-slate-200 disabled:text-slate-400"
              :class="
                selected && answer.trim()
                  ? 'bg-slate-100 text-slate-700 hover:bg-slate-200'
                  : 'bg-slate-100 text-slate-400'
              "
            >
              Kembalikan ke Admin
            </button>
            <button
              @click="handlePublish()"
              :disabled="!selected || !answer.trim()"
              class="mt-3 w-full rounded-lg px-4 py-3 text-sm font-semibold text-white shadow-soft transition disabled:cursor-not-allowed disabled:bg-emerald-300 disabled:text-white"
              :class="
                selected && answer.trim()
                  ? 'bg-emerald-600 hover:bg-emerald-700 active:bg-emerald-800'
                  : 'bg-emerald-300'
              "
            >
              Validasi & Publikasikan
            </button>
            <p v-if="!selected" class="mt-3 text-xs text-slate-500">
              Pilih aduan urgent untuk memulai proses validasi.
            </p>
            <p v-else-if="!answer.trim()" class="mt-3 text-xs text-slate-500">
              Isi jawaban resmi sebelum mengirim keputusan.
            </p>
          </div>
        </div>
      </div>
    </section>
  </PimpinanLayout>
</template>

<script setup>
import { ref, onMounted } from "vue";
import PimpinanLayout from "../../layouts/PimpinanLayout.vue";
import pimpinanService from "../../services/pimpinan.service";
import { useToastStore } from "../../stores/toast";

const loading = ref(true);
const error = ref("");
const urgent = ref([]);
const selected = ref(null);
const answer = ref("");
const toast = useToastStore();

const load = async () => {
  loading.value = true;
  error.value = "";
  try {
    const res = await pimpinanService.urgentReports();
    urgent.value = res.data?.data || res.data || [];
    selected.value = urgent.value[0] || null;
  } catch (err) {
    error.value = err?.message || "Server error";
    toast.add({ type: "danger", message: error.value });
  } finally {
    loading.value = false;
  }
};

async function handleReturn() {
  if (!selected.value) {
    toast.add({ type: "danger", message: "Pilih aduan terlebih dahulu." });
    return;
  }
  if (!answer.value.trim()) {
    toast.add({
      type: "danger",
      message: "Isi jawaban resmi sebelum mengirim keputusan.",
    });
    return;
  }
  try {
    await pimpinanService.createDisposisi(selected.value.id, {
      catatan: answer.value,
    });
    toast.add({ type: "success", message: "Dikembalikan ke admin." });
    answer.value = "";
    load();
  } catch (err) {
    toast.add({ type: "danger", message: "Gagal mengembalikan." });
  }
}

async function handlePublish() {
  if (!selected.value) {
    toast.add({ type: "danger", message: "Pilih aduan terlebih dahulu." });
    return;
  }
  if (!answer.value.trim()) {
    toast.add({
      type: "danger",
      message: "Isi jawaban resmi sebelum mengirim keputusan.",
    });
    return;
  }
  try {
    await pimpinanService.createDisposisi(selected.value.id, {
      catatan: answer.value,
    });
    toast.add({ type: "success", message: "Validasi dan publikasi berhasil." });
    answer.value = "";
    load();
  } catch (err) {
    toast.add({ type: "danger", message: "Gagal mempublikasikan." });
  }
}

onMounted(load);
</script>
