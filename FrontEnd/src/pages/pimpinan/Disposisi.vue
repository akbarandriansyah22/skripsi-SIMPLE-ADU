<template>
  <PimpinanLayout>
    <section
      class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50"
    >
      <div
        class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between"
      >
        <div>
          <p class="text-sm font-bold text-slate-950">Disposisi</p>
          <p class="mt-1 text-sm text-slate-600">
            Lihat instruksi yang telah dikirim ke unit terkait.
          </p>
        </div>
        <button
          @click="openModal"
          class="rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-blue-700 active:bg-blue-800 shadow-soft"
        >
          Tambah Disposisi
        </button>
      </div>
      <div class="mt-6 space-y-4">
        <template v-if="loading">
          <div class="rounded-lg border border-slate-200 bg-slate-50 p-5">
            <p class="text-sm text-slate-600">Memuat disposisi...</p>
          </div>
        </template>
        <template v-else-if="error">
          <div class="rounded-lg border border-red-200 bg-red-50 p-5">
            <p class="text-sm text-red-600">
              Gagal memuat disposisi: {{ error }}
            </p>
          </div>
        </template>
        <template v-else-if="disposisi.length">
          <div
            v-for="d in disposisi"
            :key="d.id"
            class="rounded-lg border border-slate-200 bg-slate-50 p-5 hover:bg-slate-100 transition"
          >
            <p class="font-semibold text-slate-950">
              {{ d.kode_tiket || d.kode || "ADU-" + d.pengaduan_id || d.title }}
            </p>
            <p class="mt-2 text-sm text-slate-600">
              {{ d.catatan || d.message || d.summary || "" }}
            </p>
            <p class="mt-3 text-xs text-slate-500">
              {{ d.tanggal || d.created_at || "" }}
            </p>
          </div>
        </template>
        <div
          v-else
          class="rounded-lg border border-slate-200 bg-slate-50 p-5 hover:bg-slate-100 transition"
        >
          <p class="text-sm text-slate-600">Belum ada disposisi.</p>
        </div>
      </div>
    </section>
    <Modal
      v-if="showModal"
      :visible="showModal"
      title="Tambah Disposisi"
      subtitle="Kirim instruksi baru ke unit terkait."
      @close="showModal = false"
    >
      <div class="space-y-4">
        <div>
		  <label class="block text-sm font-semibold text-slate-900">Pilih unit penanganan</label>
		  <select v-model="form.unitId" class="mt-2 w-full rounded-lg border border-slate-200 bg-white px-4 py-3 text-sm">
		    <option value="">Pilih unit penanganan</option>
		    <option v-for="unit in units" :key="unit.id" :value="unit.id">{{ unit.nama_unit }}</option>
		  </select>
		  <p class="mt-2 text-xs leading-relaxed text-slate-500">Pengaduan akan diteruskan kepada Kasubag sesuai unit yang dipilih.</p>
		</div>
		<div>
          <label class="block text-sm font-semibold text-slate-900"
            >ID Pengaduan</label
          >
          <input
            v-model="form.pengaduanId"
            type="text"
            placeholder="Masukkan ID pengaduan"
            class="mt-2 w-full rounded-lg border border-slate-200 bg-white px-4 py-3 text-sm text-slate-900 shadow-soft transition focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500/20"
          />
        </div>
        <div>
          <label class="block text-sm font-semibold text-slate-900"
            >Pesan Disposisi</label
          >
          <textarea
            v-model="form.message"
            rows="5"
            placeholder="Masukkan pesan disposisi"
            class="mt-2 w-full rounded-lg border border-slate-200 bg-white px-4 py-3 text-sm text-slate-900 shadow-soft transition focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500/20"
          ></textarea>
        </div>
        <p v-if="modalError" class="text-sm text-red-600">{{ modalError }}</p>
      </div>
      <template #footer>
        <div class="flex justify-end gap-3">
          <button
            type="button"
            @click="showModal = false"
            class="rounded-lg border border-slate-200 bg-white px-4 py-2 text-sm font-semibold text-slate-700 transition hover:bg-slate-50"
          >
            Batal
          </button>
          <button
            type="button"
            @click="submitDisposisi"
            :disabled="submitting"
            class="rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-blue-700"
          >
            {{ submitting ? "Menyimpan..." : "Simpan" }}
          </button>
        </div>
      </template>
    </Modal>
  </PimpinanLayout>
</template>

<script setup>
import { ref, onMounted } from "vue";
import PimpinanLayout from "../../layouts/PimpinanLayout.vue";
import Modal from "../../components/Modal.vue";
import pimpinanService from "../../services/pimpinan.service";
import { useToastStore } from "../../stores/toast";
import { errorMessage } from "../../utils/api";

const disposisi = ref([]);
const loading = ref(true);
const error = ref("");
const toast = useToastStore();

const form = ref({ pengaduanId: "", unitId: "", message: "" });
const showModal = ref(false);
const modalError = ref("");
const units = ref([]);
const submitting = ref(false);

const load = async () => {
  loading.value = true;
  error.value = "";
  try {
    const res = await pimpinanService.myDisposisi();
    disposisi.value = res.data?.data || res.data || [];
  } catch (err) {
    error.value = err?.message || "Server error";
  } finally {
    loading.value = false;
  }
};

function openModal() {
  form.value = { pengaduanId: "", unitId: "", message: "" };
  modalError.value = "";
  showModal.value = true;
}

async function submitDisposisi() {
  modalError.value = "";
  if (submitting.value) return;
  if (!String(form.value.pengaduanId).trim()) {
    modalError.value = "ID pengaduan wajib diisi.";
    return;
  }
  if (!form.value.message.trim()) {
    modalError.value = "Pesan disposisi wajib diisi.";
    return;
  }
	if (!form.value.unitId) {
	  modalError.value = "Unit tujuan wajib dipilih.";
	  return;
	}

  submitting.value = true;
  try {
    await pimpinanService.createDisposisi(form.value.pengaduanId, {
	  unit_id: Number(form.value.unitId),
      catatan: form.value.message.trim(),
    });
    toast.add({ type: "success", message: "Disposisi ditambahkan." });
    showModal.value = false;
    await load();
    await pimpinanService.myDisposisi();
  } catch (err) {
    modalError.value = errorMessage(err, "Gagal menambahkan disposisi.");
    toast.add({ type: "danger", message: modalError.value });
  } finally {
    submitting.value = false;
  }
}

onMounted(async () => {
  await load();
  try {
    const response = await pimpinanService.getUnits();
    units.value = response.data?.data || response.data || [];
  } catch {}
});
</script>
