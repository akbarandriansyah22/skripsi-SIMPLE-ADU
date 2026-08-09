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
            :class="selectedDisposisiId === d.id ? 'border-amber-400 ring-1 ring-amber-300' : 'border-slate-200'"
            class="rounded-lg bg-slate-50 p-5 transition hover:bg-slate-100"
            @click="selectDisposisi(d)"
          >
            <p class="font-semibold text-slate-950">
              {{ d.pengaduan?.kode_tiket || d.kode_tiket || `ADU-${d.pengaduan_id}` }}
            </p>
            <p class="mt-2 text-sm text-slate-600">
              {{ d.catatan || d.message || d.summary || "" }}
            </p>
            <p class="mt-3 text-xs text-slate-500">
              {{ d.tanggal || d.created_at || "" }}
            </p>
            <button type="button" class="mt-3 text-xs font-semibold text-blue-950 hover:underline" @click.stop="selectDisposisi(d)">
              {{ selectedDisposisiId === d.id ? 'Percakapan dipilih' : 'Lihat koordinasi' }}
            </button>
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
    <section class="mt-6 rounded-xl border border-amber-200 bg-amber-50 p-6">
      <p class="text-sm font-bold text-slate-950">Koordinasi Internal Urgensi Tinggi</p>
      <p class="mt-1 text-xs text-slate-600">Percakapan privat Pimpinan dan Kasubag untuk pengaduan yang dipilih.</p>
      <p v-if="selectedDisposisi" class="mt-3 text-xs font-semibold text-amber-900">Pengaduan: {{ selectedDisposisi.pengaduan?.kode_tiket || `ADU-${coordinationId}` }}</p>
      <p v-else class="mt-4 text-xs text-slate-600">Pilih pengaduan dari daftar disposisi untuk membuka koordinasi.</p>
      <p v-if="coordinationError" class="mt-3 rounded-lg border border-red-200 bg-red-50 p-3 text-xs text-red-700">{{ coordinationError }}</p>
      <p v-else-if="coordinationLoading" class="mt-3 text-xs text-slate-600">Memuat percakapan koordinasi...</p>
      <div v-else-if="coordinationMessages.length" class="mt-4 space-y-2">
        <div v-for="message in coordinationMessages" :key="message.id" class="rounded-lg bg-white p-3 text-xs">
          <p class="font-semibold">{{ message.sender_name }} · {{ message.sender_role }}</p>
          <p class="mt-1 whitespace-pre-line text-slate-700">{{ message.pesan }}</p>
        </div>
      </div>
      <p v-else-if="selectedDisposisi" class="mt-4 text-xs text-slate-600">Belum ada pesan koordinasi.</p>
      <form v-if="selectedDisposisi" class="mt-4 flex gap-3" @submit.prevent="sendCoordination">
        <input v-model="coordinationMessage" type="text" placeholder="Pesan koordinasi" class="min-w-0 flex-1 rounded-lg border border-amber-200 px-3 py-2 text-xs" />
        <button type="submit" :disabled="sendingCoordination || !coordinationMessage.trim()" class="rounded-lg bg-amber-600 px-4 py-2 text-xs font-semibold text-white disabled:opacity-50">{{ sendingCoordination ? 'Mengirim...' : 'Kirim' }}</button>
      </form>
      <p v-if="sendCoordinationError" class="mt-3 text-xs text-red-700">{{ sendCoordinationError }}</p>
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
import { computed, ref, onMounted } from "vue";
import PimpinanLayout from "../../layouts/PimpinanLayout.vue";
import Modal from "../../components/Modal.vue";
import pimpinanService from "../../services/pimpinan.service";
import { useToastStore } from "../../stores/toast";
import { errorMessage, responseData } from "../../utils/api";

const disposisi = ref([]);
const loading = ref(true);
const error = ref("");
const toast = useToastStore();

const form = ref({ pengaduanId: "", unitId: "", message: "" });
const showModal = ref(false);
const modalError = ref("");
const units = ref([]);
const submitting = ref(false);
const selectedDisposisiId = ref(null);
const coordinationMessage = ref('');
const coordinationMessages = ref([]);
const coordinationLoading = ref(false);
const coordinationError = ref('');
const sendCoordinationError = ref('');
const sendingCoordination = ref(false);
const selectedDisposisi = computed(() => disposisi.value.find((item) => item.id === selectedDisposisiId.value) || null);
const coordinationId = computed(() => Number(selectedDisposisi.value?.pengaduan_id) || 0);

const load = async () => {
  loading.value = true;
  error.value = "";
  try {
    const res = await pimpinanService.myDisposisi();
    disposisi.value = res.data?.data || res.data || [];
    const selected = disposisi.value.find((item) => item.id === selectedDisposisiId.value) || disposisi.value[0];
    if (selected) await selectDisposisi(selected);
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

async function selectDisposisi(disposition) {
  selectedDisposisiId.value = disposition?.id || null;
  coordinationMessages.value = [];
  coordinationError.value = '';
  sendCoordinationError.value = '';
  if (!coordinationId.value) {
    coordinationError.value = 'ID pengaduan pada disposisi tidak tersedia.';
    return;
  }
  coordinationLoading.value = true;
  try {
    coordinationMessages.value = responseData(await pimpinanService.coordination(coordinationId.value), []) || [];
  } catch (err) {
    coordinationError.value = errorMessage(err, 'Percakapan koordinasi tidak dapat dimuat.');
  } finally {
    coordinationLoading.value = false;
  }
}

async function sendCoordination() {
  sendCoordinationError.value = '';
  if (!coordinationId.value || !coordinationMessage.value.trim() || sendingCoordination.value) return;
  sendingCoordination.value = true;
  try {
    await pimpinanService.sendCoordination(coordinationId.value, { pesan: coordinationMessage.value.trim() });
    coordinationMessage.value = '';
    await selectDisposisi(selectedDisposisi.value);
    toast.add({ type: 'success', message: 'Koordinasi berhasil dikirim.' });
  } catch (err) {
    sendCoordinationError.value = errorMessage(err, 'Koordinasi gagal dikirim.');
  } finally {
    sendingCoordination.value = false;
  }
}
</script>
