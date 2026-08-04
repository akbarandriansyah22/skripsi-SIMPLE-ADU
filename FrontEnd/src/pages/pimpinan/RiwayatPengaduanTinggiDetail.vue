<template>
  <PimpinanLayout>
    <div v-if="loading" class="rounded-xl border border-slate-200 bg-white p-6 text-sm text-slate-500 shadow-sm">
      Memuat detail riwayat...
    </div>

    <div v-else-if="error" class="rounded-xl border border-red-200 bg-red-50 p-5 text-xs text-red-700">
      <p>{{ error }}</p>
      <button type="button" class="mt-3 rounded-full bg-red-100 px-4 py-2 font-semibold text-red-800" @click="load">Coba Lagi</button>
    </div>

    <section v-else-if="item" class="space-y-5">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <router-link to="/pimpinan/riwayat-tinggi" class="text-xs font-semibold text-blue-950 hover:underline">← Kembali ke riwayat</router-link>
          <p class="mt-4 text-[11px] font-semibold uppercase tracking-wider text-blue-950">{{ ticket(item) }}</p>
          <h2 class="mt-1 text-xl font-bold text-slate-950">{{ item.judul || 'Pengaduan tanpa judul' }}</h2>
        </div>
        <StatusBadge :label="item.status" />
      </div>

      <div class="grid gap-5 xl:grid-cols-[minmax(0,1.4fr)_minmax(280px,0.6fr)]">
        <div class="space-y-5">
          <article class="rounded-xl border border-slate-200 bg-white p-5 shadow-sm sm:p-6">
            <h3 class="text-sm font-bold text-slate-950">Informasi Pengaduan</h3>
            <dl class="mt-4 grid gap-4 text-xs sm:grid-cols-2">
              <div><dt class="text-slate-500">Tanggal masuk ke Pimpinan</dt><dd class="mt-1 font-semibold text-slate-900">{{ dateLabel(item.tanggal_masuk_pimpinan, true) }}</dd></div>
              <div><dt class="text-slate-500">Pembaruan terakhir</dt><dd class="mt-1 font-semibold text-slate-900">{{ dateLabel(item.updated_at, true) }}</dd></div>
              <div><dt class="text-slate-500">Nama mahasiswa</dt><dd class="mt-1 font-semibold text-slate-900">{{ studentName }}</dd></div>
              <div><dt class="text-slate-500">NIM</dt><dd class="mt-1 font-semibold text-slate-900">{{ item.user?.mahasiswa?.nim || '-' }}</dd></div>
              <div><dt class="text-slate-500">Program studi</dt><dd class="mt-1 font-semibold text-slate-900">{{ item.user?.mahasiswa?.program_studi || '-' }}</dd></div>
              <div><dt class="text-slate-500">Angkatan</dt><dd class="mt-1 font-semibold text-slate-900">{{ item.user?.mahasiswa?.angkatan || '-' }}</dd></div>
              <div><dt class="text-slate-500">Kategori</dt><dd class="mt-1 font-semibold text-slate-900">{{ item.kategori?.nama || '-' }}</dd></div>
              <div><dt class="text-slate-500">Unit yang menangani</dt><dd class="mt-1 font-semibold text-slate-900">{{ unitName }}</dd></div>
            </dl>
            <div class="mt-5 rounded-lg border border-slate-200 bg-slate-50 p-4">
              <h4 class="text-xs font-bold text-slate-950">Isi Pengaduan</h4>
              <p class="mt-2 whitespace-pre-line text-sm leading-relaxed text-slate-700">{{ item.deskripsi || 'Tidak ada deskripsi.' }}</p>
            </div>
          </article>

          <article class="rounded-xl border border-slate-200 bg-white p-5 shadow-sm sm:p-6">
            <h3 class="text-sm font-bold text-slate-950">Hasil AI</h3>
            <dl class="mt-4 grid gap-4 text-xs sm:grid-cols-2">
              <div><dt class="text-slate-500">Urgensi final</dt><dd class="mt-1"><StatusBadge kind="urgency" :label="item.urgensi || 'Tinggi'" /></dd></div>
              <div><dt class="text-slate-500">Status AI</dt><dd class="mt-1 font-semibold text-slate-900">{{ item.ai_status || '-' }}</dd></div>
              <div><dt class="text-slate-500">Sentimen</dt><dd class="mt-1"><StatusBadge kind="sentiment" :label="item.sentimen || '-'" /></dd></div>
              <div><dt class="text-slate-500">Skor sentimen</dt><dd class="mt-1 font-semibold text-slate-900">{{ item.skor_sentimen ?? '-' }}</dd></div>
              <div><dt class="text-slate-500">Skor positif</dt><dd class="mt-1 font-semibold text-slate-900">{{ item.skor_positif ?? '-' }}</dd></div>
              <div><dt class="text-slate-500">Skor negatif</dt><dd class="mt-1 font-semibold text-slate-900">{{ item.skor_negatif ?? '-' }}</dd></div>
              <div><dt class="text-slate-500">Skor urgensi</dt><dd class="mt-1 font-semibold text-slate-900">{{ item.urgency_score ?? '-' }}</dd></div>
            </dl>
            <div class="mt-5 space-y-3 text-xs">
              <div><h4 class="font-bold text-slate-950">Teks hasil pembersihan</h4><p class="mt-1 whitespace-pre-line text-slate-600">{{ item.cleaned_text || '-' }}</p></div>
              <div><h4 class="font-bold text-slate-950">Alasan urgensi</h4><p class="mt-1 whitespace-pre-line text-slate-600">{{ item.urgency_reason || '-' }}</p></div>
              <div><h4 class="font-bold text-slate-950">Penjelasan sentimen</h4><p class="mt-1 whitespace-pre-line text-slate-600">{{ item.sentiment_explanation || '-' }}</p></div>
              <div v-if="matchedWords.length"><h4 class="font-bold text-slate-950">Kata yang terdeteksi</h4><div class="mt-2 flex flex-wrap gap-2"><span v-for="word in matchedWords" :key="word" class="rounded-full bg-slate-100 px-2.5 py-1 text-[11px] text-slate-700">{{ word }}</span></div></div>
            </div>
          </article>

          <article class="rounded-xl border border-slate-200 bg-white p-5 shadow-sm sm:p-6">
            <h3 class="text-sm font-bold text-slate-950">Validasi Admin Fakultas</h3>
            <div v-if="item.validasi" class="mt-4 grid gap-4 text-xs sm:grid-cols-2">
              <div><dt class="text-slate-500">Hasil validasi</dt><dd class="mt-1 font-semibold text-slate-900">{{ item.validasi.status_validasi }}</dd></div>
              <div><dt class="text-slate-500">Catatan</dt><dd class="mt-1 whitespace-pre-line text-slate-700">{{ item.validasi.catatan || '-' }}</dd></div>
            </div>
            <p v-else class="mt-4 text-xs text-slate-500">Data validasi belum tersedia.</p>
          </article>

          <article class="rounded-xl border border-slate-200 bg-white p-5 shadow-sm sm:p-6">
            <h3 class="text-sm font-bold text-slate-950">Disposisi Pimpinan</h3>
            <div v-if="item.disposisi" class="mt-4 grid gap-4 text-xs sm:grid-cols-2">
              <div><dt class="text-slate-500">Unit tujuan</dt><dd class="mt-1 font-semibold text-slate-900">{{ item.disposisi.unit?.nama_unit || unitName }}</dd></div>
              <div><dt class="text-slate-500">Catatan disposisi</dt><dd class="mt-1 whitespace-pre-line text-slate-700">{{ item.disposisi.catatan || '-' }}</dd></div>
            </div>
            <p v-else class="mt-4 text-xs text-slate-500">Belum ada disposisi.</p>
          </article>

          <article class="rounded-xl border border-amber-200 bg-amber-50 p-5 shadow-sm sm:p-6">
            <h3 class="text-sm font-bold text-slate-950">Koordinasi Internal</h3>
            <p class="mt-1 text-xs text-slate-600">Pesan internal Pimpinan dan unit yang ditugaskan.</p>
            <div v-if="item.koordinasi_internal?.length" class="mt-4 space-y-3">
              <div v-for="message in item.koordinasi_internal" :key="message.id" class="rounded-xl border border-amber-100 bg-white p-3">
                <div class="flex flex-wrap justify-between gap-2 text-[11px] text-slate-500"><span class="font-semibold text-slate-900">{{ message.sender_name }} · {{ message.sender_role }}</span><span>{{ dateLabel(message.created_at, true) }}</span></div>
                <p class="mt-2 whitespace-pre-line text-xs text-slate-700">{{ message.pesan }}</p>
                <a v-if="coordinationAttachment(message)" :href="coordinationAttachment(message)" target="_blank" rel="noreferrer" class="mt-2 inline-flex text-xs font-semibold text-blue-950 hover:underline">Buka lampiran koordinasi</a>
              </div>
            </div>
            <p v-else class="mt-4 text-xs text-slate-500">Belum ada koordinasi internal.</p>
          </article>
        </div>

        <aside class="space-y-5">
          <article class="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <h3 class="text-sm font-bold text-slate-950">Lampiran</h3>
            <div v-if="attachment" class="mt-4"><img v-if="attachment.isImage" :src="attachment.url" :alt="attachment.name" class="max-h-72 max-w-full rounded-lg border border-slate-200 object-contain" /><a v-else :href="attachment.url" target="_blank" rel="noreferrer" class="inline-flex rounded-full bg-blue-50 px-4 py-2 text-xs font-semibold text-blue-950 hover:bg-blue-100">Buka lampiran</a><p class="mt-2 break-all text-xs text-slate-600">{{ attachment.name }}</p></div>
            <p v-else class="mt-4 text-xs text-slate-500">Tidak ada lampiran.</p>
          </article>

          <article class="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <h3 class="text-sm font-bold text-slate-950">Timeline Perubahan Status</h3>
            <div v-if="timeline.length" class="mt-4 space-y-4"><div v-for="(event, index) in timeline" :key="event.id || `${event.status_baru}-${index}`" class="relative pl-6"><span class="absolute left-0 top-1.5 h-2.5 w-2.5 rounded-full bg-blue-950"></span><span v-if="index < timeline.length - 1" class="absolute bottom-[-18px] left-[4px] top-4 w-px bg-slate-200"></span><p class="text-xs font-semibold text-slate-950">{{ event.status_lama || 'Awal' }} → {{ event.status_baru }}</p><p v-if="event.catatan" class="mt-1 whitespace-pre-line text-xs leading-relaxed text-slate-600">{{ event.catatan }}</p><p class="mt-1 text-[11px] text-slate-500">{{ dateLabel(event.created_at, true) }}</p></div></div>
            <p v-else class="mt-4 rounded-lg border border-dashed border-slate-200 px-3 py-6 text-center text-xs text-slate-500">Belum ada timeline status.</p>
          </article>
        </aside>
      </div>
    </section>
  </PimpinanLayout>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import PimpinanLayout from '../../layouts/PimpinanLayout.vue'
import StatusBadge from '../../components/StatusBadge.vue'
import service from '../../services/pimpinan.service'
import { dateLabel, errorMessage, responseData } from '../../utils/api'

const route = useRoute()
const item = ref(null)
const loading = ref(true)
const error = ref('')
const attachmentUrl = ref('')
const coordinationAttachmentUrls = ref({})

const studentName = computed(() => item.value?.user?.mahasiswa?.nama_lengkap || item.value?.user?.nama_lengkap || '-')
const unitName = computed(() => item.value?.unit?.nama_unit || item.value?.disposisi?.unit?.nama_unit || '-')
const timeline = computed(() => [...(item.value?.riwayat_status_pengaduan || [])].sort((a, b) => new Date(a.created_at || 0).getTime() - new Date(b.created_at || 0).getTime()))
const matchedWords = computed(() => {
  const value = item.value?.matched_words
  if (Array.isArray(value)) return value.map(String)
  if (typeof value === 'string') {
    try { const parsed = JSON.parse(value); return Array.isArray(parsed) ? parsed.map(String) : [] } catch { return [] }
  }
  return []
})
const attachment = computed(() => {
  if (!item.value?.lampiran || !attachmentUrl.value) return null
  const type = String(item.value.lampiran_mime_type || item.value.lampiran).toLowerCase()
  return { name: item.value.lampiran_nama_asli || item.value.lampiran, url: attachmentUrl.value, isImage: ['image/jpeg', 'image/png'].includes(type) || /\.(jpe?g|png)$/i.test(type) }
})

function ticket(value) { return value.kode_tiket || `ADU-${value.id}` }
function coordinationAttachment(message) { return coordinationAttachmentUrls.value[message.id] || '' }

async function load() {
  loading.value = true
  error.value = ''
  try {
    item.value = responseData(await service.historyDetail(route.params.id))
    if (item.value?.lampiran_url) attachmentUrl.value = URL.createObjectURL(await service.getAttachment(item.value.lampiran_url))
    await Promise.all((item.value?.koordinasi_internal || []).map(async (message) => {
      if (!message.lampiran_url) return
      try {
        coordinationAttachmentUrls.value[message.id] = URL.createObjectURL(await service.getAttachment(message.lampiran_url))
      } catch {}
    }))
  } catch (err) {
    error.value = errorMessage(err, 'Detail riwayat tidak dapat dimuat.')
  } finally {
    loading.value = false
  }
}

onMounted(load)
onBeforeUnmount(() => {
  if (attachmentUrl.value) URL.revokeObjectURL(attachmentUrl.value)
  Object.values(coordinationAttachmentUrls.value).forEach((url) => URL.revokeObjectURL(url))
})
</script>
