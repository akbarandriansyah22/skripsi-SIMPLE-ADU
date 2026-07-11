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
              {{ detail?.kode || "ADU-" + detail?.id || "—" }}
            </h2>
            <p class="mt-1 text-sm text-slate-600">
              {{ detail?.unit_name || detail?.unit || "—" }}
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
              <template v-if="detail?.comments?.length">
                <div
                  v-for="c in detail.comments"
                  :key="c.id"
                  class="rounded-lg bg-white p-4 border border-slate-200"
                >
                  <p class="text-sm font-semibold text-slate-950">
                    {{ c.author_name || c.author || c.user?.nama || "—" }}
                  </p>
                  <p class="mt-2 text-sm text-slate-700">
                    {{ c.message || c.body || c.text }}
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
                <p>{{ detail?.pelapor?.prodi || detail?.prodi || "-" }}</p>
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
            <button
              @click="handleForward()"
              class="w-full rounded-lg bg-blue-600 px-4 py-3 text-sm font-semibold text-white transition hover:bg-blue-700 active:bg-blue-800 shadow-soft"
            >
              Teruskan ke Pimpinan
            </button>
            <button
              @click="markComplete()"
              class="w-full rounded-lg bg-emerald-600 px-4 py-3 text-sm font-semibold text-white transition hover:bg-emerald-700 active:bg-emerald-800 shadow-soft"
            >
              Update Status
            </button>
          </div>
        </div>
      </div>
    </section>
  </AdminLayout>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
import adminService from "../../services/admin.service";

const props = defineProps({ id: [String, Number] });
const route = useRoute();
const id = props.id || route.params.id;

const detail = ref(null);
const loading = ref(true);
const error = ref("");

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
    detail.value = res.data || null;
  } catch (err) {
    error.value = err?.message || "Server error";
    toast.add({ type: "danger", message: error.value });
  } finally {
    loading.value = false;
  }
};

onMounted(load);

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

async function markComplete() {
  if (!detail.value) return;
  try {
    await adminService.updateStatus(detail.value.id, { status: "selesai" });
    toast.add({ type: "success", message: "Status diperbarui." });
    load();
  } catch (err) {
    toast.add({ type: "danger", message: "Gagal memperbarui status." });
  }
}
</script>
