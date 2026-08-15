<template>
  <PimpinanLayout>
    <div v-if="loading" class="rounded-xl border border-slate-200 bg-white p-6 text-sm text-slate-500 shadow-sm">Memuat detail monitoring...</div>
    <div v-else-if="error" class="rounded-xl border border-red-200 bg-red-50 p-5 text-xs text-red-700">
      <p>{{ error }}</p>
      <button type="button" class="mt-3 rounded-full bg-red-100 px-4 py-2 font-semibold text-red-800" @click="load">Coba Lagi</button>
    </div>

    <section v-else-if="item" class="space-y-5">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <router-link to="/pimpinan/monitoring" class="text-xs font-semibold text-blue-950 hover:underline">← Kembali ke monitoring</router-link>
          <p class="mt-4 text-[11px] font-semibold uppercase tracking-wider text-blue-950">{{ item.kode_tiket || `ADU-${item.id}` }}</p>
          <h2 class="mt-1 text-xl font-bold text-slate-950">{{ item.judul || 'Pengaduan tanpa judul' }}</h2>
        </div>
        <StatusBadge :label="item.status" />
      </div>

      <div class="grid gap-5 xl:grid-cols-[minmax(0,1.4fr)_minmax(280px,0.6fr)]">
        <div class="space-y-5">
          <article class="rounded-xl border border-slate-200 bg-white p-5 shadow-sm sm:p-6">
            <h3 class="text-sm font-bold text-slate-950">Progress Penanganan</h3>
            <dl class="mt-4 grid gap-4 text-xs sm:grid-cols-2">
              <div><dt class="text-slate-500">Status terkini</dt><dd class="mt-1 font-semibold text-slate-900">{{ item.status }}</dd></div>
              <div><dt class="text-slate-500">Unit penanganan</dt><dd class="mt-1 font-semibold text-slate-900">{{ item.unit?.nama_unit || item.disposisi?.unit?.nama_unit || '-' }}</dd></div>
              <div><dt class="text-slate-500">Mahasiswa</dt><dd class="mt-1 font-semibold text-slate-900">{{ item.user?.nama_lengkap || '-' }}</dd></div>
              <div><dt class="text-slate-500">NIM</dt><dd class="mt-1 font-semibold text-slate-900">{{ item.user?.mahasiswa?.nim || '-' }}</dd></div>
              <div><dt class="text-slate-500">Kategori</dt><dd class="mt-1 font-semibold text-slate-900">{{ item.kategori?.nama || '-' }}</dd></div>
              <div><dt class="text-slate-500">Pembaruan terakhir</dt><dd class="mt-1 font-semibold text-slate-900">{{ dateLabel(item.updated_at, true) }}</dd></div>
            </dl>
            <div class="mt-5 rounded-lg border border-slate-200 bg-slate-50 p-4">
              <h4 class="text-xs font-bold text-slate-950">Isi Pengaduan</h4>
              <p class="mt-2 whitespace-pre-line text-sm leading-relaxed text-slate-700">{{ item.deskripsi || 'Tidak ada deskripsi.' }}</p>
            </div>
          </article>

          <article class="rounded-xl border border-slate-200 bg-white p-5 shadow-sm sm:p-6">
            <h3 class="text-sm font-bold text-slate-950">Disposisi</h3>
            <div v-if="item.disposisi" class="mt-4 grid gap-4 text-xs sm:grid-cols-2">
              <div><dt class="text-slate-500">Unit tujuan</dt><dd class="mt-1 font-semibold text-slate-900">{{ item.disposisi.unit?.nama_unit || item.unit?.nama_unit || '-' }}</dd></div>
              <div><dt class="text-slate-500">Catatan disposisi</dt><dd class="mt-1 whitespace-pre-line text-slate-700">{{ item.disposisi.catatan || '-' }}</dd></div>
            </div>
            <p v-else class="mt-4 text-xs text-slate-500">Belum ada disposisi. Aduan mungkin diteruskan langsung oleh Admin Fakultas.</p>
          </article>

          <article class="rounded-xl border border-amber-200 bg-amber-50 p-5 shadow-sm sm:p-6">
            <h3 class="text-sm font-bold text-slate-950">Koordinasi Internal</h3>
            <div v-if="coordination.length" class="mt-4 space-y-3">
              <div v-for="message in coordination" :key="message.id" class="rounded-xl border border-amber-100 bg-white p-3">
                <div class="flex flex-wrap justify-between gap-2 text-[11px] text-slate-500">
                  <span class="font-semibold text-slate-900">{{ message.sender_name }} · {{ message.sender_role }}</span>
                  <span>{{ dateLabel(message.created_at, true) }}</span>
                </div>
                <p class="mt-2 whitespace-pre-line text-xs text-slate-700">{{ message.pesan }}</p>
              </div>
            </div>
            <p v-else class="mt-4 text-xs text-slate-500">Belum ada koordinasi internal.</p>
          </article>
        </div>

        <aside class="space-y-5">
          <article class="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <h3 class="text-sm font-bold text-slate-950">Riwayat Status</h3>
            <div v-if="timeline.length" class="mt-4 space-y-4">
              <div v-for="(event, index) in timeline" :key="event.id || `${event.status_baru}-${index}`" class="relative pl-6">
                <span class="absolute left-0 top-1.5 h-2.5 w-2.5 rounded-full bg-blue-950"></span>
                <span v-if="index < timeline.length - 1" class="absolute bottom-[-18px] left-[4px] top-4 w-px bg-slate-200"></span>
                <p class="text-xs font-semibold text-slate-950">{{ event.status_lama || 'Awal' }} → {{ event.status_baru }}</p>
                <p v-if="event.catatan" class="mt-1 whitespace-pre-line text-xs leading-relaxed text-slate-600">{{ event.catatan }}</p>
                <p class="mt-1 text-[11px] text-slate-500">{{ dateLabel(event.created_at, true) }}</p>
              </div>
            </div>
            <p v-else class="mt-4 rounded-lg border border-dashed border-slate-200 px-3 py-6 text-center text-xs text-slate-500">Belum ada riwayat status.</p>
          </article>
        </aside>
      </div>
    </section>
  </PimpinanLayout>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import PimpinanLayout from '../../layouts/PimpinanLayout.vue'
import StatusBadge from '../../components/StatusBadge.vue'
import service from '../../services/pimpinan.service'
import { dateLabel, errorMessage, responseData } from '../../utils/api'

const route = useRoute()
const item = ref(null)
const loading = ref(true)
const error = ref('')

const timeline = computed(() => [...(item.value?.riwayat_status_pengaduan || [])].sort((a, b) => new Date(a.created_at || 0) - new Date(b.created_at || 0)))
const coordination = computed(() => item.value?.koordinasi_internal || [])

async function load() {
  loading.value = true
  error.value = ''
  try {
    item.value = responseData(await service.getPengaduan(route.params.id))
  } catch (err) {
    error.value = errorMessage(err, 'Detail monitoring tidak dapat dimuat.')
  } finally {
    loading.value = false
  }
}

onMounted(load)
watch(() => route.params.id, load)
</script>
