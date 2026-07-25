<template>
  <MahasiswaLayout>
    <section class="rounded-xl border border-slate-200 bg-white p-5 shadow-sm sm:p-6">
      <div class="mb-6">
        <h2 class="text-sm font-bold text-slate-950">Kirim Aduan</h2>
        <p class="mt-1 text-xs text-slate-500">Isi detail aduan Anda dengan jelas dan lengkap.</p>
      </div>

      <form class="space-y-6" novalidate @submit.prevent="handleSubmit">
        <div class="grid gap-5 lg:grid-cols-2">
          <div class="space-y-2">
            <label for="judul" class="text-xs font-semibold text-slate-900">Judul Aduan</label>
            <input id="judul" v-model.trim="form.judul" type="text" maxlength="150" placeholder="Masukkan judul aduan" class="form-input" />
            <p v-if="errors.judul" class="text-xs text-red-600">{{ errors.judul }}</p>
          </div>
          <div class="space-y-2">
            <label for="kategori" class="text-xs font-semibold text-slate-900">Kategori</label>
            <select id="kategori" v-model.number="form.kategori_id" class="form-input">
              <option :value="0" disabled>Pilih kategori</option>
              <option v-for="item in kategoriOptions" :key="item.id" :value="item.id">{{ item.nama }}</option>
            </select>
            <p v-if="errors.kategori_id" class="text-xs text-red-600">{{ errors.kategori_id }}</p>
          </div>
        </div>

        <div class="space-y-2">
          <label for="deskripsi" class="text-xs font-semibold text-slate-900">Deskripsi Aduan</label>
          <textarea id="deskripsi" v-model.trim="form.deskripsi" rows="6" maxlength="5000" placeholder="Jelaskan masalah secara rinci" class="form-input resize-y"></textarea>
          <div class="flex justify-between text-[11px] text-slate-400"><span v-if="errors.deskripsi" class="text-red-600">{{ errors.deskripsi }}</span><span class="ml-auto">{{ form.deskripsi.length }}/5000</span></div>
        </div>

        <div class="grid gap-5 lg:grid-cols-2 lg:items-end">
          <div class="space-y-2">
            <label for="lampiran" class="text-xs font-semibold text-slate-900">Lampiran Bukti</label>
            <input id="lampiran" ref="fileInput" type="file" accept=".jpg,.jpeg,.png,.pdf,image/jpeg,image/png,application/pdf" class="form-input file:mr-3 file:rounded file:border-0 file:bg-slate-100 file:px-2 file:py-1 file:text-xs" @change="onFileChange" />
            <div v-if="form.lampiran" class="flex items-center gap-3 text-xs text-slate-500">
              <span class="truncate">{{ form.lampiran.name }}</span><button type="button" class="text-red-600 hover:underline" @click="clearFile">Hapus</button>
            </div>
            <p v-if="errors.lampiran" class="text-xs text-red-600">{{ errors.lampiran }}</p>
            <img v-if="previewUrl" :src="previewUrl" alt="Pratinjau lampiran" class="mt-2 h-20 w-20 rounded-lg border border-slate-200 object-cover" />
          </div>
          <div class="rounded-2xl border border-blue-100 bg-blue-50 px-4 py-3 text-xs leading-relaxed text-blue-700">
            Setelah dikirim, sistem AI akan menganalisis deskripsi untuk skor sentimen, sentimen, urgensi, dan status analisis.
          </div>
        </div>

        <div class="flex items-center gap-4">
          <button type="submit" :disabled="submitting" class="rounded-full bg-emerald-600 px-6 py-3 text-xs font-semibold text-white shadow-sm transition hover:bg-emerald-700 disabled:cursor-not-allowed disabled:opacity-60">
            {{ submitting ? 'Mengirim...' : 'Kirim Aduan' }}
          </button>
          <p v-if="submitError" class="text-xs text-red-600">{{ submitError }}</p>
        </div>
      </form>
    </section>
  </MahasiswaLayout>
</template>

<script setup>
import { computed, onBeforeUnmount, ref } from 'vue'
import { useRouter } from 'vue-router'
import MahasiswaLayout from '../../layouts/MahasiswaLayout.vue'
import pengaduanService from '../../services/pengaduan.service'
import { errorMessage, validateEvidence } from '../../utils/api'
import { useToastStore } from '../../stores/toast'

const router = useRouter()
const toast = useToastStore()
const fileInput = ref(null)
const submitting = ref(false)
const submitError = ref('')
const previewUrl = ref('')
const errors = ref({})
const kategoriOptions = [
  { id: 1, nama: 'Akademik' },
  { id: 2, nama: 'Fasilitas' },
  { id: 3, nama: 'Kemahasiswaan' },
]
const form = ref({ judul: '', kategori_id: 0, deskripsi: '', lampiran: null })

function validate() {
  const next = {}
  if (!form.value.judul) next.judul = 'Judul aduan wajib diisi.'
  if (!form.value.kategori_id) next.kategori_id = 'Kategori wajib dipilih.'
  if (!form.value.deskripsi) next.deskripsi = 'Deskripsi aduan wajib diisi.'
  if (form.value.lampiran) next.lampiran = validateEvidence(form.value.lampiran)
  errors.value = next
  return !Object.values(next).some(Boolean)
}

function onFileChange(event) {
  const file = event.target.files?.[0] || null
  const validationError = validateEvidence(file)
  if (validationError) {
    errors.value = { ...errors.value, lampiran: validationError }
    form.value.lampiran = null
    if (fileInput.value) fileInput.value.value = ''
    return
  }
  errors.value = { ...errors.value, lampiran: '' }
  form.value.lampiran = file
  if (previewUrl.value) URL.revokeObjectURL(previewUrl.value)
  previewUrl.value = file?.type.startsWith('image/') ? URL.createObjectURL(file) : ''
}

function clearFile() {
  form.value.lampiran = null
  if (fileInput.value) fileInput.value.value = ''
  if (previewUrl.value) URL.revokeObjectURL(previewUrl.value)
  previewUrl.value = ''
}

async function handleSubmit() {
  submitError.value = ''
  if (!validate()) return
  submitting.value = true
  try {
    const payload = new FormData()
    payload.append('kategori_id', String(form.value.kategori_id))
    payload.append('judul', form.value.judul)
    payload.append('deskripsi', form.value.deskripsi)
    if (form.value.lampiran) payload.append('lampiran', form.value.lampiran)
    await pengaduanService.create(payload)
    toast.add({ type: 'success', message: 'Aduan berhasil dikirim.' })
    router.replace({ name: 'MahasiswaPengaduan', query: { success: '1' } })
  } catch (error) {
    submitError.value = errorMessage(error, 'Gagal mengirim aduan.')
    toast.add({ type: 'danger', message: submitError.value })
  } finally {
    submitting.value = false
  }
}

onBeforeUnmount(() => { if (previewUrl.value) URL.revokeObjectURL(previewUrl.value) })
</script>

<style scoped>
.form-input {
  @apply w-full rounded-2xl border border-slate-200 bg-white px-4 py-3 text-xs text-slate-900 shadow-sm transition placeholder:text-slate-400 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/10;
}
</style>
