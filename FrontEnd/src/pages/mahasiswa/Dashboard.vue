<template>
  <MahasiswaLayout>
    <section class="grid gap-6 lg:grid-cols-3">
      <div class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50">
        <p
          class="text-xs uppercase tracking-widest text-blue-600 font-semibold"
        >
          Pengaduan Aktif
        </p>
        <p class="mt-4 text-4xl font-bold text-slate-950">12</p>
        <p class="mt-3 text-sm text-slate-600 leading-relaxed">
          Aduan Anda yang sedang diproses oleh admin dan pimpinan.
        </p>
      </div>
      <div class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50">
        <p
          class="text-xs uppercase tracking-widest text-blue-600 font-semibold"
        >
          Selesai
        </p>
        <p class="mt-4 text-4xl font-bold text-slate-950">8</p>
        <p class="mt-3 text-sm text-slate-600 leading-relaxed">
          Laporan yang telah ditindaklanjuti dan ditutup.
        </p>
      </div>
      <div class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50">
        <p
          class="text-xs uppercase tracking-widest text-blue-600 font-semibold"
        >
          Notifikasi Baru
        </p>
        <p class="mt-4 text-4xl font-bold text-slate-950">3</p>
        <p class="mt-3 text-sm text-slate-600 leading-relaxed">
          Pesan terbaru tentang aduan dan layanan.
        </p>
      </div>
    </section>
    <section class="mt-6 grid gap-6 lg:grid-cols-2">
      <div class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-bold text-slate-950">Ringkasan Aktif</p>
            <p class="mt-1 text-sm text-slate-600">
              Pemantauan aduan terbaru Anda.
            </p>
          </div>
          <span
            class="rounded-lg bg-emerald-100 px-3 py-1 text-xs font-semibold text-emerald-700"
            >Tingkat Respons Cepat</span
          >
        </div>
        <div class="mt-6 space-y-4">
          <div
            class="rounded-lg border border-slate-200 bg-slate-50 p-4 hover:bg-slate-100 transition"
          >
            <div class="flex items-center justify-between gap-4">
              <p class="font-semibold text-slate-900">ADU-2025-0001</p>
              <span
                class="rounded-full bg-slate-200 px-3 py-1 text-xs font-semibold uppercase tracking-wider text-slate-700"
                >Diproses</span
              >
            </div>
            <p class="mt-3 text-sm text-slate-600">
              Nilai UTS mata kuliah Basis Data tidak sesuai dengan absensi.
            </p>
          </div>
          <div
            class="rounded-lg border border-slate-200 bg-slate-50 p-4 hover:bg-slate-100 transition"
          >
            <div class="flex items-center justify-between gap-4">
              <p class="font-semibold text-slate-900">ADU-2024-0982</p>
              <span
                class="rounded-full bg-emerald-100 px-3 py-1 text-xs font-semibold uppercase tracking-wider text-emerald-700"
                >Selesai</span
              >
            </div>
            <p class="mt-3 text-sm text-slate-600">
              Aduan Anda telah diselesaikan oleh pihak terkait.
            </p>
          </div>
        </div>
      </div>
      <div class="rounded-xl bg-white p-6 shadow-card ring-1 ring-slate-200/50">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-bold text-slate-950">Notifikasi Terbaru</p>
            <p class="mt-1 text-sm text-slate-600">
              Perubahan status aduan Anda.
            </p>
          </div>
          <button
            class="rounded-lg bg-slate-100 px-4 py-2 text-sm font-medium text-slate-700 transition hover:bg-slate-200"
          >
            Lihat Semua
          </button>
        </div>
        <div class="mt-6 space-y-4">
          <div
            class="rounded-lg border border-slate-200 bg-slate-50 p-4 hover:bg-slate-100 transition"
          >
            <p class="font-semibold text-slate-900">Status Aduan Diperbarui</p>
            <p class="mt-2 text-sm text-slate-600">
              Aduan ADU-2025-0001 (Fasilitas Kelas) kini dalam proses.
            </p>
            <p class="mt-3 text-xs text-slate-500">Baru saja</p>
          </div>
          <div
            class="rounded-lg border border-slate-200 bg-slate-50 p-4 hover:bg-slate-100 transition"
          >
            <p class="font-semibold text-slate-900">Balasan Baru dari Admin</p>
            <p class="mt-2 text-sm text-slate-600">
              Admin Fakultas Teknik telah memberi balasan untuk laporan KRS.
            </p>
            <p class="mt-3 text-xs text-slate-500">1 jam lalu</p>
          </div>
        </div>
      </div>
    </section>
  </MahasiswaLayout>
</template>

<script setup>
import MahasiswaLayout from "../../layouts/MahasiswaLayout.vue";
</script>
