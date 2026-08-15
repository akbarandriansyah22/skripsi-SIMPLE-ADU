import api from './axios'

const dashboard = () => api.get('/api/kasubag/dashboard')
const getAllPengaduan = () => api.get('/api/kasubag/pengaduan')
const getPengaduanById = (id) => api.get(`/api/kasubag/pengaduan/${id}`)
const startProcess = (id) => api.patch(`/api/kasubag/pengaduan/${id}/proses`)
const addResponse = (id, payload) => api.post(`/api/kasubag/pengaduan/${id}/respon`, payload)
const finish = (id) => api.patch(`/api/kasubag/pengaduan/${id}/selesai`)
const returnToAdmin = (id, alasan) => api.patch(`/api/kasubag/pengaduan/${id}/kembalikan`, { alasan })
const getAttachment = async (url) => (await api.get(url, { responseType: 'blob' })).data
const coordination = (id) => api.get(`/api/pengaduan/${id}/koordinasi`)
const sendCoordination = (id, payload) => api.post(`/api/pengaduan/${id}/koordinasi`, payload)

export default { dashboard, getAllPengaduan, getPengaduanById, startProcess, addResponse, finish, returnToAdmin, getAttachment, coordination, sendCoordination }
