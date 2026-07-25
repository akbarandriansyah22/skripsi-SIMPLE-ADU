<template>
  <AdminLayout>
    <section class="grid gap-6 lg:grid-cols-[1.3fr_0.9fr]">
      <div class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50">
        <div
          class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between"
        >
          <div>
            <p
              class="text-xs uppercase tracking-wider text-blue-600 font-semibold"
            >
              Detail Tiket
            </p>
            <h2 class="mt-2 text-2xl font-bold text-slate-950">
              {{ detail?.kode_tiket || detail?.kode || "ADU-" + detail?.id || "—" }}
            </h2>
            <p class="mt-1 text-sm text-slate-600">
              {{ detail?.unit?.nama_unit || detail?.unit_name || detail?.unit || "—" }}
            </p>
          </div>
          <span
            :class="statusBadgeClass(detail?.status)"
            class="px-4 py-2 text-sm font-semibold"
            >{{ detail?.status || "—" }}</span
          >
        </div>
        <div class="mt-6 space-y-4">
          <div
            class="rounded-lg border border-slate-200 bg-slate-50 p-5 hover:bg-slate-100 transition"
          >
            <p class="text-sm font-bold text-slate-950">Isi Laporan</p>
            <p class="mt-3 text-slate-700 leading-relaxed">
              {{ detail?.deskripsi || "Tidak ada deskripsi." }}
            </p>
          </div>
          <div class="grid gap-4 sm:grid-cols-2">
            <div
              class="rounded-lg border border-slate-200 bg-slate-50 p-5 hover:bg-slate-100 transition"
            >
              <p class="text-sm font-bold text-slate-950">Skor Sentimen</p>
              <p class="mt-2 text-sm text-slate-700">
                {{ detail?.ai_status === "success" ? detail?.skor_sentimen : "—" }}
              </p>
            </div>
            <div
              class="rounded-lg border border-slate-200 bg-slate-50 p-5 hover:bg-slate-100 transition"
            >
              <p class="text-sm font-bold text-slate-950">Sentimen</p>
              <p class="mt-2 text-sm text-slate-700">
                {{ detail?.sentimen || "—" }}
              </p>
            </div>
            <div
              class="rounded-lg border border-slate-200 bg-slate-50 p-5 hover:bg-slate-100 transition"
            >
              <p class="text-sm font-bold text-slate-950">Urgensi</p>
              <p class="mt-2 text-sm text-slate-700">
                {{ detail?.urgensi || "—" }}
              </p>
            </div>
          </div>
          <div class="rounded-lg border border-slate-200 bg-slate-50 p-5">
            <p class="text-sm font-bold text-slate-950">Status AI</p>
            <p class="mt-2 text-sm text-slate-700">
              {{ detail?.ai_status || "pending" }}
            </p>
            <p v-if="detail?.ai_status !== 'success'" class="mt-2 text-sm text-slate-600">
              Analisis AI sedang menunggu pemrosesan.
            </p>
          </div>
          <div v-if="detail?.validasi" class="rounded-lg border border-blue-200 bg-blue-50 p-5">
            <p class="text-sm font-bold text-slate-950">Validasi Admin</p>
            <p class="mt-2 text-sm text-slate-700">{{ detail.validasi.status_validasi }}{{ detail.validasi.catatan ? ` — ${detail.validasi.catatan}` : "" }}</p>
          </div>
          <div class="rounded-lg border border-slate-200 bg-slate-50 p-5">
            <p class="text-sm font-bold text-slate-950">Bukti Pendukung</p>
            <div v-if="detail?.attachments?.length" class="mt-3 space-y-3">
              <div
                v-for="att in detail.attachments"
                :key="att.id"
                class="inline-flex items-center gap-3 rounded-lg bg-white px-4 py-3 shadow-soft border border-slate-200"
              >
                <div
                  class="h-10 w-10 rounded-lg bg-blue-100 text-blue-700 flex items-center justify-center font-bold text-sm"
                >
                  {{
                    (att.type || "").toUpperCase().includes("PDF")
                      ? "PDF"
                      : (att.type || "FILE").toUpperCase()
                  }}
                </div>
                <div>
                  <p class="font-semibold text-slate-950">{{ att.name }}</p>
                  <p class="text-xs text-slate-600">{{ att.size || "" }}</p>
                </div>
              </div>
            </div>
            <div v-else class="mt-3 text-sm text-slate-600">
              Tidak ada lampiran.
            </div>
          </div>
          <div class="rounded-lg border border-slate-200 bg-slate-50 p-5">
            <p class="text-sm font-bold text-slate-950">Riwayat Komunikasi</p>
            <div class="mt-4 space-y-3">
              <template v-if="responses.length">
                <div
                  v-for="c in responses"
                  :key="c.id"
                  class="rounded-lg bg-white p-4 border border-slate-200"
                >
                  <p class="text-sm font-semibold text-slate-950">
                    {{ c.author_name || c.author || c.user?.nama_lengkap || c.user?.nama || "—" }}
                  </p>
                  <p class="mt-2 text-sm text-slate-700">
                    {{ c.pesan || c.message || c.body || c.text }}
                  </p>
                  <p class="mt-2 text-xs text-slate-500">
                    {{ formatDate(c.created_at) }}
                  </p>
                </div>
              </template>
              <div
                v-else
                class="rounded-lg bg-white p-4 border border-slate-200"
              >
                <p class="text-sm text-slate-600">Belum ada komunikasi.</p>
              </div>
            </div>
          </div>
          <form class="rounded-lg border border-slate-200 bg-white p-5" @submit.prevent="sendResponse">
            <label class="text-sm font-bold text-slate-950">Jawaban Admin</label>
            <textarea
              v-model="responseText"
              rows="4"
              placeholder="Tulis balasan untuk mahasiswa..."
              class="mt-3 w-full rounded-lg border border-slate-200 px-4 py-3 text-sm text-slate-900 shadow-soft"
            ></textarea>
            <button
              type="submit"
              :disabled="sendingResponse || !responseText.trim()"
              class="mt-3 rounded-lg bg-blue-950 px-4 py-2 text-sm font-semibold text-white transition hover:bg-blue-900 disabled:opacity-50"
            >
              {{ sendingResponse ? "Mengirim..." : "Kirim Jawaban" }}
            </button>
          </form>
        </div>
      </div>
      <div class="space-y-6">
        <div
          class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50"
        >
          <p class="text-sm font-bold text-slate-950">Data Pelapor</p>
          <div class="mt-5 space-y-4 text-sm text-slate-700">
            <div class="grid gap-3 sm:grid-cols-2">
              <div>
                <p class="font-semibold text-slate-900">Nama</p>
                <p>
                  {{
                    detail?.pelapor?.nama ||
                    detail?.pelapor_name ||
                    detail?.user?.nama_lengkap ||
                    detail?.user?.nama ||
                    "-"
                  }}
                </p>
              </div>
              <div>
                <p class="font-semibold text-slate-900">NIM</p>
                <p>{{ detail?.pelapor?.nim || detail?.nim || "-" }}</p>
              </div>
            </div>
            <div class="grid gap-3 sm:grid-cols-2">
              <div>
                <p class="font-semibold text-slate-900">Program Studi</p>
                <p>{{ detail?.pelapor?.prodi || detail?.program_studi || detail?.prodi || "-" }}</p>
              </div>
              <div>
                <p class="font-semibold text-slate-900">Email</p>
                <p class="truncate">
                  {{ detail?.pelapor?.email || detail?.user?.email || "-" }}
                </p>
              </div>
            </div>
          </div>
        </div>
        <div
          class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50"
        >
          <p class="text-sm font-bold text-slate-950">Aksi Admin</p>
          <div class="mt-5 space-y-4">
			<textarea v-model="validationNote" rows="3" placeholder="Catatan validasi awal" class="w-full rounded-lg border border-slate-200 px-4 py-3 text-sm"></textarea>
			<div class="grid grid-cols-2 gap-3">
			  <button @click="validateComplaint('Diterima')" class="rounded-lg bg-emerald-600 px-4 py-3 text-sm font-semibold text-white">Terima Validasi</button>
			  <button @click="validateComplaint('Ditolak')" class="rounded-lg bg-red-600 px-4 py-3 text-sm font-semibold text-white">Tolak Validasi</button>
			</div>
            <select
              v-model.number="unitForm"
              class="w-full rounded-lg border border-slate-200 px-4 py-3 text-sm"
            >
              <option :value="0">Pilih unit tujuan</option>
              <option v-for="unit in units" :key="unit.id" :value="unit.id">
                {{ unit.nama_unit || unit.nama || unit.name }}
              </option>
            </select>
            <button
              @click="assignSelectedUnit()"
              :disabled="!unitForm"
              class="w-full rounded-lg bg-slate-100 px-4 py-3 text-sm font-semibold text-slate-700 transition hover:bg-slate-200 disabled:opacity-50"
            >
              Tetapkan Unit
            </button>
            <button
              @click="handleForward()"
              class="w-full rounded-lg bg-blue-600 px-4 py-3 text-sm font-semibold text-white transition hover:bg-blue-700 active:bg-blue-800 shadow-soft"
            >
              Teruskan ke Pimpinan
            </button>
            <button
              @click="handleReanalyze()"
              :disabled="reanalyzing"
              class="w-full rounded-lg bg-indigo-600 px-4 py-3 text-sm font-semibold text-white transition hover:bg-indigo-700 active:bg-indigo-800 shadow-soft disabled:opacity-50"
            >
              {{ reanalyzing ? "Menganalisis..." : "Analisis Ulang AI" }}
            </button>
          </div>
        </div>
      </div>
    </section>
  </AdminLayout>
</template>

<script setup>
import { computed, ref, onMounted } from "vue";
import { useRoute } from "vue-router";
import AdminLayout from "../../layouts/AdminLayout.vue";
import adminService from "../../services/admin.service";
import pengaduanService from "../../services/pengaduan.service";

const props = defineProps({ id: [String, Number] });
const route = useRoute();
const id = props.id || route.params.id;

const detail = ref(null);
const loading = ref(true);
const error = ref("");
const units = ref([]);
const unitForm = ref(0);
const responseText = ref("");
const validationNote = ref("");
const sendingResponse = ref(false);
const reanalyzing = ref(false);
const responses = computed(() => detail.value?.respon_pengaduan || detail.value?.responses || detail.value?.respon || detail.value?.comments || []);

function statusBadgeClass(status) {
  if (!status)
    return "rounded-full px-3 py-1 text-xs font-semibold inline-block bg-slate-100 text-slate-700";
  const s = String(status).toLowerCase();
  if (s.includes("selesai") || s.includes("done") || s.includes("seles"))
    return "rounded-full bg-emerald-100 px-3 py-1 text-xs font-semibold text-emerald-700 inline-block";
  if (
    s.includes("proses") ||
    s.includes("diproses") ||
    s.includes("in progress")
  )
    return "rounded-full bg-amber-100 px-3 py-1 text-xs font-semibold text-amber-700 inline-block";
  if (s.includes("ditolak") || s.includes("tolak") || s.includes("rejected"))
    return "rounded-full bg-red-100 px-3 py-1 text-xs font-semibold text-red-700 inline-block";
  return "rounded-full px-3 py-1 text-xs font-semibold inline-block bg-slate-100 text-slate-700";
}

function formatDate(iso) {
  if (!iso) return "";
  try {
    return new Date(iso).toLocaleString();
  } catch {
    return "";
  }
}

const load = async () => {
  loading.value = true;
  error.value = "";
  try {
    const res = await adminService.getPengaduanById(id);
    detail.value = res.data?.data || res.data || null;
    unitForm.value = detail.value?.unit_id || 0;
  } catch (err) {
    error.value = err?.message || "Server error";
    toast.add({ type: "danger", message: error.value });
  } finally {
    loading.value = false;
  }
};

onMounted(load);

onMounted(async () => {
  try {
    const res = await adminService.getUnits();
    units.value = res.data?.data || res.data || [];
  } catch {
    units.value = [];
  }
});

import { useToastStore } from "../../stores/toast";
const toast = useToastStore();

async function handleForward() {
  if (!detail.value) return;
  try {
    await adminService.forwardToPimpinan(detail.value.id);
    toast.add({
      type: "success",
      message: "Pengaduan diteruskan ke pimpinan.",
    });
    load();
  } catch (err) {
    toast.add({ type: "danger", message: "Gagal meneruskan pengaduan." });
  }
}

async function validateComplaint(status) {
  if (!detail.value) return;
  try {
    await adminService.validate(detail.value.id, { status_validasi: status, catatan: validationNote.value });
    toast.add({ type: "success", message: "Validasi pengaduan tersimpan." });
    load();
  } catch {
    toast.add({ type: "danger", message: "Validasi pengaduan gagal disimpan." });
  }
}

async function handleReanalyze() {
  if (!detail.value) return;
  reanalyzing.value = true;
  try {
    await adminService.reanalyzeAI(detail.value.id);
    toast.add({ type: "success", message: "Analisis AI berhasil diperbarui." });
    load();
  } catch (err) {
    toast.add({ type: "danger", message: "Gagal memperbarui analisis AI." });
  } finally {
    reanalyzing.value = false;
  }
}

async function assignSelectedUnit() {
  if (!detail.value || !unitForm.value) return;
  try {
    await adminService.assignUnit(detail.value.id, { unit_id: unitForm.value });
    toast.add({ type: "success", message: "Unit berhasil ditetapkan." });
    load();
  } catch (err) {
    toast.add({ type: "danger", message: "Gagal menetapkan unit." });
  }
}

async function sendResponse() {
  if (!detail.value || !responseText.value.trim()) return;
  sendingResponse.value = true;
  try {
    await pengaduanService.addRespon(detail.value.id, { pesan: responseText.value });
    toast.add({ type: "success", message: "Jawaban berhasil dikirim." });
    responseText.value = "";
    load();
  } catch (err) {
    toast.add({ type: "danger", message: "Gagal mengirim jawaban." });
  } finally {
    sendingResponse.value = false;
  }
}
</script>
