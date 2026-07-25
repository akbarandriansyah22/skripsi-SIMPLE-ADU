import api from './axios'

const dashboard = () => api.get('/api/pimpinan/dashboard')
const urgentReports = () => api.get('/api/pimpinan/pengaduan/urgensi-tinggi')
const createDisposisi = (id, payload) => api.post(`/api/pimpinan/pengaduan/${id}/disposisi`, payload)
const myDisposisi = () => api.get('/api/pimpinan/disposisi')
const getUnits = () => api.get('/api/admin/unit')

export default { dashboard, urgentReports, createDisposisi, myDisposisi, getUnits }
