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
            <label class="block text-sm font-semibold text-slate-900">Pilih unit penanganan</label>
            <select
              v-model="unitId"
              :disabled="submitting"
              class="w-full rounded-lg border border-slate-200 bg-white px-4 py-3 text-sm text-slate-900 shadow-soft focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500/20"
            >
              <option value="">Pilih unit penanganan</option>
              <option v-for="unit in units" :key="unit.id" :value="unit.id">
                {{ unit.nama_unit }}
              </option>
            </select>
            <p class="text-xs leading-relaxed text-slate-500">Pengaduan akan diteruskan kepada Kasubag sesuai unit yang dipilih.</p>
            <label class="block text-sm font-semibold text-slate-900">Catatan Disposisi</label>
            <textarea
              rows="6"
              v-model="catatan"
              :disabled="submitting"
              class="w-full rounded-lg border border-slate-200 bg-white px-4 py-4 text-sm text-slate-900 shadow-soft transition focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500/20"
              placeholder="Ketik catatan disposisi untuk unit tujuan..."
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
              @click="submitDisposisi"
              :disabled="submitting || !selected || !unitId || !catatan.trim()"
              class="mt-3 w-full rounded-lg px-4 py-3 text-sm font-semibold text-white shadow-soft transition disabled:cursor-not-allowed disabled:bg-emerald-300 disabled:text-white"
              :class="
                selected && unitId && catatan.trim() && !submitting
                  ? 'bg-emerald-600 hover:bg-emerald-700 active:bg-emerald-800'
                  : 'bg-emerald-300'
              "
            >
              {{ submitting ? "Mengirim..." : "Disposisikan ke Unit" }}
            </button>
            <p v-if="!selected" class="mt-3 text-xs text-slate-500">
              Pilih aduan urgensi tinggi terlebih dahulu.
            </p>
            <p v-else-if="!unitId" class="mt-3 text-xs text-slate-500">
              Pilih unit penanganan sebelum mengirim disposisi.
            </p>
            <p v-else-if="!catatan.trim()" class="mt-3 text-xs text-slate-500">
              Isi catatan disposisi sebelum mengirim.
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
import { errorMessage } from "../../utils/api";

const loading = ref(true);
const error = ref("");
const urgent = ref([]);
const selected = ref(null);
const catatan = ref("");
const unitId = ref("");
const units = ref([]);
const submitting = ref(false);
const toast = useToastStore();

const load = async () => {
  loading.value = true;
  error.value = "";
  try {
    const res = await pimpinanService.urgentReports();
    urgent.value = res.data?.data || res.data || [];
    selected.value = urgent.value[0] || null;
  } catch (err) {
    error.value = errorMessage(err, "Daftar aduan urgensi tinggi tidak dapat dimuat.");
    toast.add({ type: "danger", message: error.value });
  } finally {
    loading.value = false;
  }
};

function selectReport(report) {
  selected.value = report;
  unitId.value = "";
  catatan.value = "";
}

async function submitDisposisi() {
  if (!selected.value) {
    toast.add({ type: "danger", message: "Pilih aduan terlebih dahulu." });
    return;
  }
  if (!unitId.value) {
    toast.add({ type: "danger", message: "Pilih unit penanganan terlebih dahulu." });
    return;
  }
  if (!catatan.value.trim()) {
    toast.add({
      type: "danger",
      message: "Isi catatan disposisi sebelum mengirim.",
    });
    return;
  }
  if (submitting.value) return;
  submitting.value = true;
  try {
    await pimpinanService.createDisposisi(selected.value.id, {
      unit_id: Number(unitId.value),
      catatan: catatan.value.trim(),
    });
    toast.add({ type: "success", message: "Disposisi berhasil dikirim." });
    unitId.value = "";
    catatan.value = "";
    await Promise.all([load(), pimpinanService.myDisposisi()]);
  } catch (err) {
    toast.add({ type: "danger", message: errorMessage(err, "Disposisi gagal dikirim.") });
  } finally {
    submitting.value = false;
  }
}

onMounted(async () => {
  await Promise.all([load(), loadUnits()]);
});

async function loadUnits() {
  try {
    const res = await pimpinanService.getUnits();
    units.value = res.data?.data || res.data || [];
  } catch (err) {
    toast.add({ type: "danger", message: errorMessage(err, "Unit tujuan tidak dapat dimuat.") });
  }
}
</script>
