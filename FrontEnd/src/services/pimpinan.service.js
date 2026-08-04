import api from './axios'

const dashboard = () => api.get('/api/pimpinan/dashboard')
const urgentReports = () => api.get('/api/pimpinan/pengaduan/urgensi-tinggi')
const createDisposisi = (id, payload) => api.post(`/api/pimpinan/pengaduan/${id}/disposisi`, payload)
const myDisposisi = () => api.get('/api/pimpinan/disposisi')
const getUnits = () => api.get('/api/pimpinan/unit/penanganan')
const coordination = (id) => api.get(`/api/pengaduan/${id}/koordinasi`)
const sendCoordination = (id, payload) => api.post(`/api/pengaduan/${id}/koordinasi`, payload)
const categories = () => api.get('/api/pengaduan/kategori')
const history = (params = {}) => api.get('/api/pimpinan/pengaduan/riwayat', { params })
const historyDetail = (id) => api.get(`/api/pimpinan/pengaduan/riwayat/${id}`)
const getAttachment = async (url) => (await api.get(url, { responseType: 'blob' })).data

export default { dashboard, urgentReports, createDisposisi, myDisposisi, getUnits, coordination, sendCoordination, categories, history, historyDetail, getAttachment }
