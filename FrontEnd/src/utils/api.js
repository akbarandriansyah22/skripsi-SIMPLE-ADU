export function responseData(response, fallback = null) {
  if (response == null) return fallback
  return response.data?.data ?? response.data ?? fallback
}

export function responseList(response) {
  const value = responseData(response, [])
  if (Array.isArray(value)) return value
  if (Array.isArray(value?.items)) return value.items
  if (Array.isArray(value?.data)) return value.data
  return []
}

export function errorMessage(error, fallback = 'Terjadi kesalahan pada server.') {
  return error?.response?.data?.message || error?.response?.data?.error || error?.message || fallback
}

export function statusLabel(status) {
  return status || 'Belum tersedia'
}

export function dateLabel(value, withTime = false) {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '-'
  return date.toLocaleString('id-ID', withTime
    ? { dateStyle: 'medium', timeStyle: 'short' }
    : { dateStyle: 'medium' })
}

export function statusClass(status) {
  const value = String(status || '').toLowerCase()
  if (value.includes('selesai')) return 'bg-emerald-100 text-emerald-700'
  if (value.includes('tolak')) return 'bg-red-100 text-red-700'
  if (value.includes('proses') || value.includes('distribusi') || value.includes('unit')) return 'bg-amber-100 text-amber-700'
  if (value.includes('verifikasi') || value.includes('ajukan') || value.includes('tunggu')) return 'bg-blue-100 text-blue-700'
  return 'bg-slate-100 text-slate-700'
}

export function urgencyClass(urgency) {
  const value = String(urgency || '').toLowerCase()
  if (value.includes('tinggi') || value.includes('high')) return 'bg-red-100 text-red-700'
  if (value.includes('sedang') || value.includes('medium')) return 'bg-amber-100 text-amber-700'
  return 'bg-emerald-100 text-emerald-700'
}

export function sentimentClass(sentiment) {
  const value = String(sentiment || '').toLowerCase()
  if (value.includes('negatif') || value.includes('negative')) return 'bg-red-100 text-red-700'
  if (value.includes('positif') || value.includes('positive')) return 'bg-emerald-100 text-emerald-700'
  return 'bg-slate-100 text-slate-700'
}

export function validateEvidence(file) {
  if (!file) return ''
  const allowed = ['image/jpeg', 'image/png', 'application/pdf']
  if (!allowed.includes(file.type)) return 'Lampiran hanya dapat berupa JPG, JPEG, PNG, atau PDF.'
  if (file.size <= 0 || file.size > 10 * 1024 * 1024) return 'Ukuran lampiran maksimal 10 MB.'
  return ''
}
