import api from './axios'

const dashboard = () => api.get('/api/kasubag/dashboard')
const getAllPengaduan = () => api.get('/api/kasubag/pengaduan')
const getPengaduanById = (id) => api.get(`/api/kasubag/pengaduan/${id}`)
const startProcess = (id) => api.patch(`/api/kasubag/pengaduan/${id}/proses`)
const addResponse = (id, payload) => api.post(`/api/kasubag/pengaduan/${id}/respon`, payload)
const finish = (id) => api.patch(`/api/kasubag/pengaduan/${id}/selesai`)
const getAttachment = async (url) => (await api.get(url, { responseType: 'blob' })).data

export default { dashboard, getAllPengaduan, getPengaduanById, startProcess, addResponse, finish, getAttachment }
