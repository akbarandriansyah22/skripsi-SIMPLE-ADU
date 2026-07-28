<template>
  <KasubagLayout>
    <div v-if="loading" class="grid gap-5 xl:grid-cols-[minmax(0,1.4fr)_minmax(280px,0.6fr)]">
      <div class="h-[520px] animate-pulse rounded-xl border border-slate-200 bg-white"></div>
      <div class="h-[360px] animate-pulse rounded-xl border border-slate-200 bg-white"></div>
    </div>

    <div v-else-if="error" class="rounded-xl border border-red-200 bg-red-50 p-5 text-xs text-red-700">
      <p>{{ error }}</p>
      <button type="button" class="mt-3 rounded-full bg-red-100 px-4 py-2 font-semibold text-red-800 hover:bg-red-200" @click="load">Coba Lagi</button>
    </div>

    <section v-else-if="item" class="grid gap-5 xl:grid-cols-[minmax(0,1.4fr)_minmax(280px,0.6fr)]">
      <div class="space-y-5">
        <article class="rounded-xl border border-slate-200 bg-white p-5 shadow-sm sm:p-6">
          <div class="flex flex-col gap-3 border-b border-slate-100 pb-5 sm:flex-row sm:items-start sm:justify-between">
            <div class="min-w-0">
              <p class="text-[11px] font-semibold uppercase tracking-wider text-blue-950">{{ ticket(item) }}</p>
              <h2 class="mt-2 break-words text-xl font-bold text-slate-950">{{ item.judul || 'Pengaduan tanpa judul' }}</h2>
              <p class="mt-1 text-xs text-slate-500">Diajukan {{ dateLabel(item.created_at, true) }}</p>
            </div>
            <StatusBadge :label="item.status" />
          </div>

          <dl class="mt-5 grid gap-4 text-xs sm:grid-cols-2">
            <div><dt class="text-slate-500">Nama mahasiswa</dt><dd class="mt-1 font-semibold text-slate-900">{{ item.user?.nama_lengkap || item.nama_mahasiswa || '-' }}</dd></div>
            <div><dt class="text-slate-500">NIM</dt><dd class="mt-1 font-semibold text-slate-900">{{ item.user?.nim || item.nim || '-' }}</dd></div>
            <div><dt class="text-slate-500">Email</dt><dd class="mt-1 break-all font-semibold text-slate-900">{{ item.user?.email || '-' }}</dd></div>
            <div><dt class="text-slate-500">Kategori</dt><dd class="mt-1 font-semibold text-slate-900">{{ item.kategori?.nama || item.kategori || '-' }}</dd></div>
            <div><dt class="text-slate-500">Unit penanganan</dt><dd class="mt-1 font-semibold text-slate-900">{{ item.unit?.nama_unit || '-' }}</dd></div>
            <div><dt class="text-slate-500">Tanggal pengajuan</dt><dd class="mt-1 font-semibold text-slate-900">{{ dateLabel(item.created_at) }}</dd></div>
          </dl>

          <div class="mt-5 rounded-lg border border-slate-200 bg-slate-50 p-4">
            <h3 class="text-xs font-bold text-slate-950">Deskripsi Pengaduan</h3>
            <p class="mt-2 whitespace-pre-line text-sm leading-relaxed text-slate-700">{{ item.deskripsi || 'Tidak ada deskripsi.' }}</p>
          </div>

          <div class="mt-5 rounded-lg border border-slate-200 p-4">
            <div class="flex items-center justify-between gap-3">
              <h3 class="text-xs font-bold text-slate-950">Lampiran Bukti</h3>
              <span v-if="item.lampiran" class="text-[11px] text-slate-500">{{ formatBytes(item.lampiran_ukuran) }}</span>
            </div>
            <div v-if="attachment" class="mt-3">
              <img v-if="attachment.isImage" :src="attachment.url" :alt="attachment.name" class="max-h-72 max-w-full rounded-lg border border-slate-200 object-contain" />
              <a v-else-if="attachment.isPdf" :href="attachment.url" target="_blank" rel="noreferrer" class="inline-flex rounded-full bg-blue-50 px-4 py-2 text-xs font-semibold text-blue-950 hover:bg-blue-100">Buka PDF</a>
              <p class="mt-2 text-xs text-slate-600">{{ attachment.name }}</p>
            </div>
            <p v-else class="mt-3 text-xs text-slate-500">Tidak ada lampiran atau URL lampiran tidak disediakan oleh API.</p>
          </div>
        </article>

        <article class="rounded-xl border border-slate-200 bg-white p-5 shadow-sm sm:p-6">
          <div class="flex items-center justify-between gap-3"><div><h3 class="text-sm font-bold text-slate-950">Riwayat Percakapan</h3><p class="mt-1 text-xs text-slate-500">Balasan dari mahasiswa dan petugas terkait.</p></div><span class="text-[11px] text-slate-400">{{ conversation.length }} pesan</span></div>
          <div v-if="conversation.length" class="mt-5 space-y-3">
            <div v-for="messageItem in conversation" :key="messageItem.id" :class="messageItem.isMine ? 'ml-5 bg-blue-50 sm:ml-12' : 'mr-5 bg-slate-50 sm:mr-12'" class="rounded-xl border border-slate-200 p-4">
              <div class="flex flex-wrap items-center justify-between gap-2"><div><p class="text-xs font-semibold text-slate-950">{{ senderName(messageItem) }}</p><p class="mt-0.5 text-[11px] text-slate-500">{{ senderRole(messageItem) }}</p></div><time class="text-[11px] text-slate-500">{{ dateLabel(messageItem.created_at, true) }}</time></div>
              <p class="mt-3 whitespace-pre-line break-words text-sm leading-relaxed text-slate-700">{{ messageItem.pesan || messageItem.message || '-' }}</p>
              <a v-if="messageItem.lampiran && responseAttachmentUrl(messageItem)" :href="responseAttachmentUrl(messageItem)" target="_blank" rel="noreferrer" class="mt-3 inline-flex text-xs font-semibold text-blue-950 hover:underline">Buka lampiran respons</a>
            </div>
          </div>
          <p v-else class="mt-5 rounded-lg border border-dashed border-slate-200 px-4 py-8 text-center text-xs text-slate-500">Belum ada balasan pada pengaduan ini.</p>
        </article>

        <article class="rounded-xl border border-slate-200 bg-white p-5 shadow-sm sm:p-6">
          <div><h3 class="text-sm font-bold text-slate-950">Balas Pengaduan</h3><p class="mt-1 text-xs text-slate-500">Tuliskan tindak lanjut atau jawaban untuk mahasiswa.</p></div>
          <form class="mt-5 space-y-3" @submit.prevent="sendResponse">
            <textarea v-model="replyText" rows="5" maxlength="5000" :disabled="sending" placeholder="Tuliskan tindak lanjut atau jawaban untuk mahasiswa." class="w-full resize-y rounded-2xl border border-slate-200 px-4 py-3 text-xs text-slate-900 shadow-sm placeholder:text-slate-400 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/10"></textarea>
            <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"><div><input ref="replyFileInput" type="file" accept=".jpg,.jpeg,.png,.pdf,image/jpeg,image/png,application/pdf" :disabled="sending" class="block max-w-full text-xs text-slate-600 file:mr-3 file:rounded file:border-0 file:bg-slate-100 file:px-2 file:py-1 file:text-xs" @change="onReplyFileChange" /><p v-if="replyFile" class="mt-1 text-[11px] text-slate-500">{{ replyFile.name }} <button type="button" class="font-semibold text-red-600" @click="clearReplyFile">Hapus</button></p><p v-if="replyError" class="mt-1 text-xs text-red-600">{{ replyError }}</p></div><button type="submit" :disabled="sending || !canReply" class="rounded-full bg-emerald-600 px-5 py-3 text-xs font-semibold text-white shadow-sm transition hover:bg-emerald-700 disabled:cursor-not-allowed disabled:opacity-50">{{ sending ? 'Mengirim...' : 'Kirim Balasan' }}</button></div>
          </form>
          <p v-if="!canReply" class="mt-3 text-[11px] text-amber-700">Balasan dapat dikirim setelah aduan berstatus Diproses.</p>
        </article>
      </div>

      <aside class="space-y-5">
        <article class="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <h3 class="text-sm font-bold text-slate-950">Informasi Pengaduan</h3>
          <dl class="mt-4 space-y-3 text-xs"><div class="flex items-center justify-between gap-3"><dt class="text-slate-500">Status</dt><dd><StatusBadge :label="item.status" /></dd></div><div class="flex items-center justify-between gap-3"><dt class="text-slate-500">Kategori</dt><dd class="font-semibold text-right">{{ item.kategori?.nama || '-' }}</dd></div><div class="flex items-center justify-between gap-3"><dt class="text-slate-500">Unit</dt><dd class="font-semibold text-right">{{ item.unit?.nama_unit || '-' }}</dd></div><div class="flex items-center justify-between gap-3"><dt class="text-slate-500">Sentimen</dt><dd><StatusBadge kind="sentiment" :label="item.sentimen" /></dd></div><div class="flex items-center justify-between gap-3"><dt class="text-slate-500">Skor sentimen</dt><dd class="font-semibold">{{ item.skor_sentimen ?? '-' }}</dd></div><div class="flex items-center justify-between gap-3"><dt class="text-slate-500">Urgensi</dt><dd><StatusBadge kind="urgency" :label="item.urgensi" /></dd></div><div class="flex items-center justify-between gap-3"><dt class="text-slate-500">AI</dt><dd class="font-semibold">{{ item.ai_status || 'pending' }}</dd></div></dl>
          <p class="mt-4 border-t border-slate-100 pt-3 text-[11px] leading-relaxed text-slate-500">Hasil AI hanya informasi pendukung dan tidak dapat diubah oleh Kasubag.</p>
        </article>

        <article class="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <h3 class="text-sm font-bold text-slate-950">Perbarui Penanganan</h3>
          <p class="mt-1 text-xs text-slate-500">Gunakan status yang tersedia dari backend.</p>
          <form class="mt-4 space-y-3" @submit.prevent="updateStatus">
            <select v-model="nextStatus" :disabled="savingStatus || !statusOptions.length" class="w-full rounded-xl border border-slate-200 px-3 py-3 text-xs disabled:bg-slate-50"><option value="">{{ statusOptions.length ? 'Pilih status berikutnya' : 'Tidak ada aksi status tersedia' }}</option><option v-for="option in statusOptions" :key="option.value" :value="option.value">{{ option.label }}</option></select>
            <button type="submit" :disabled="savingStatus || !nextStatus" class="w-full rounded-full bg-emerald-600 px-4 py-3 text-xs font-semibold text-white disabled:cursor-not-allowed disabled:opacity-50">{{ savingStatus ? 'Menyimpan...' : 'Simpan Perubahan' }}</button>
          </form>
        </article>

        <article class="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <h3 class="text-sm font-bold text-slate-950">Riwayat Penanganan</h3>
          <div v-if="history.length" class="mt-4 space-y-4"><div v-for="(event, index) in history" :key="event.id || `${event.status_baru}-${index}`" class="relative pl-6"><span class="absolute left-0 top-1.5 h-2.5 w-2.5 rounded-full bg-blue-950"></span><span v-if="index < history.length - 1" class="absolute bottom-[-18px] left-[4px] top-4 w-px bg-slate-200"></span><p class="text-xs font-semibold text-slate-950">{{ event.status_baru || event.status || '-' }}</p><p v-if="event.catatan" class="mt-1 text-xs leading-relaxed text-slate-600">{{ event.catatan }}</p><p class="mt-1 text-[11px] text-slate-500">{{ dateLabel(event.created_at, true) }}</p></div></div><p v-else class="mt-4 rounded-lg border border-dashed border-slate-200 px-3 py-6 text-center text-xs text-slate-500">Belum ada riwayat penanganan.</p>
        </article>
      </aside>
    </section>
  </KasubagLayout>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import KasubagLayout from '../../layouts/KasubagLayout.vue'
import StatusBadge from '../../components/StatusBadge.vue'
import service from '../../services/kasubag.service'
import { dateLabel, errorMessage, responseData, validateEvidence } from '../../utils/api'
import { useAuthStore } from '../../stores/auth'
import { useToastStore } from '../../stores/toast'

const route = useRoute()
const auth = useAuthStore()
const toast = useToastStore()
const item = ref(null)
const loading = ref(true)
const error = ref('')
const replyText = ref('')
const replyFile = ref(null)
const replyFileInput = ref(null)
const replyError = ref('')
const sending = ref(false)
const savingStatus = ref(false)
const nextStatus = ref('')
const attachmentUrls = ref({ complaint: '', responses: {} })

const canReply = computed(() => item.value?.status === 'Diproses' && replyText.value.trim().length > 0)
const conversation = computed(() => [...(item.value?.respon_pengaduan || [])].sort((a, b) => new Date(a.created_at || 0).getTime() - new Date(b.created_at || 0).getTime()).map((message) => ({ ...message, isMine: Number(message.user_id) === Number(auth.user?.id) })))
const history = computed(() => [...(item.value?.riwayat_status_pengaduan || item.value?.riwayat_status || [])].sort((a, b) => new Date(a.created_at || 0).getTime() - new Date(b.created_at || 0).getTime()))
const statusOptions = computed(() => {
  if (item.value?.status === 'Diteruskan ke Unit') return [{ value: 'Diproses', label: 'Diproses' }]
  if (item.value?.status === 'Diproses') return [{ value: 'Selesai', label: 'Selesai' }]
  return []
})
const attachment = computed(() => {
  const file = item.value?.lampiran
  const url = attachmentUrls.value.complaint
  if (!file || !url) return null
  const type = String(item.value?.lampiran_mime_type || file).toLowerCase()
  return { name: item.value?.lampiran_nama_asli || file, url, isImage: ['image/jpeg', 'image/png'].includes(type) || /\.(jpe?g|png)$/i.test(type), isPdf: type === 'application/pdf' || /\.pdf$/i.test(type) }
})

function ticket(value) { return value.kode_tiket || `ADU-${value.id}` }
function senderName(message) { return message.user?.nama_lengkap || (message.isMine ? auth.user?.nama_lengkap : 'Pengguna') || 'Pengguna' }
function senderRole(message) { const role = String(message.user?.role || '').toLowerCase(); if (message.isMine) return auth.user?.unit_name ? `Kasubag ${auth.user.unit_name}` : 'Kasubag'; if (role.includes('mahasiswa')) return 'Mahasiswa'; if (role.includes('admin')) return 'Admin Fakultas'; if (role.includes('pimpinan')) return 'Pimpinan Fakultas'; if (role.includes('kasubag')) return 'Kasubag'; return 'Petugas' }
function responseAttachmentUrl(message) { return attachmentUrls.value.responses[message.id] || '' }
function formatBytes(value) { if (!value) return ''; if (value < 1024) return `${value} B`; if (value < 1024 * 1024) return `${Math.round(value / 1024)} KB`; return `${(value / (1024 * 1024)).toFixed(1)} MB` }

async function load() {
  loading.value = true; error.value = ''
  try { item.value = responseData(await service.getPengaduanById(route.params.id)); nextStatus.value = ''; await loadAttachments() }
  catch (err) { error.value = errorMessage(err, 'Detail pengaduan tidak dapat dimuat.') }
  finally { loading.value = false }
}

function onReplyFileChange(event) {
  const file = event.target.files?.[0] || null
  const validation = validateEvidence(file)
  if (validation) { replyError.value = validation; replyFile.value = null; if (replyFileInput.value) replyFileInput.value.value = ''; return }
  replyError.value = ''; replyFile.value = file
}
function clearReplyFile() { replyFile.value = null; if (replyFileInput.value) replyFileInput.value.value = '' }

async function sendResponse() {
  replyError.value = ''
  if (!replyText.value.trim()) { replyError.value = 'Balasan wajib diisi.'; return }
  if (!canReply.value) { replyError.value = 'Aduan harus berstatus Diproses sebelum dibalas.'; return }
  sending.value = true
  try {
    const payload = replyFile.value ? new FormData() : { pesan: replyText.value.trim() }
    if (replyFile.value) { payload.append('pesan', replyText.value.trim()); payload.append('lampiran', replyFile.value) }
    await service.addResponse(route.params.id, payload)
    replyText.value = ''; clearReplyFile(); toast.add({ type: 'success', message: 'Balasan berhasil dikirim.' }); await load()
  } catch (err) { replyError.value = errorMessage(err, 'Balasan gagal dikirim.'); toast.add({ type: 'danger', message: replyError.value }) }
  finally { sending.value = false }
}

async function updateStatus() {
  const option = statusOptions.value.find((entry) => entry.value === nextStatus.value)
  if (!option) return
  savingStatus.value = true
  try { if (option.value === 'Diproses') await service.startProcess(route.params.id); if (option.value === 'Selesai') await service.finish(route.params.id); toast.add({ type: 'success', message: `Status diperbarui menjadi ${option.label}.` }); await load() }
  catch (err) { toast.add({ type: 'danger', message: errorMessage(err, 'Status gagal diperbarui.') }) }
  finally { savingStatus.value = false }
}

async function loadAttachments() {
  const next = { complaint: '', responses: {} }
  if (item.value?.lampiran_url) {
    try { next.complaint = URL.createObjectURL(await service.getAttachment(item.value.lampiran_url)) } catch {}
  }
  await Promise.all(conversation.value.map(async (message) => {
    if (!message.lampiran_url) return
    try { next.responses[message.id] = URL.createObjectURL(await service.getAttachment(message.lampiran_url)) } catch {}
  }))
  attachmentUrls.value = next
}

onMounted(load)
onBeforeUnmount(() => {
  clearReplyFile()
  if (attachmentUrls.value.complaint) URL.revokeObjectURL(attachmentUrls.value.complaint)
  Object.values(attachmentUrls.value.responses).forEach((url) => URL.revokeObjectURL(url))
})
</script>
