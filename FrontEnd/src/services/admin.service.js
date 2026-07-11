import api from './axios'

const dashboard = () => api.get('/api/admin/dashboard')
const getAllPengaduan = () => api.get('/api/admin/pengaduan')
const getPengaduanById = (id) => api.get(`/api/admin/pengaduan/${id}`)
const updateStatus = (id, payload) => api.patch(`/api/admin/pengaduan/${id}/status`, payload)
const assignUnit = (id, payload) => api.patch(`/api/admin/pengaduan/${id}/unit`, payload)
const forwardToPimpinan = (id) => api.patch(`/api/admin/pengaduan/${id}/forward`)
const getUnits = () => api.get('/api/admin/unit')

export default { dashboard, getAllPengaduan, getPengaduanById, updateStatus, assignUnit, forwardToPimpinan, getUnits }
