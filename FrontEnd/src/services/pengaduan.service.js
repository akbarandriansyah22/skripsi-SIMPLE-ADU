import api from './axios'

const create = (payload) => api.post('/api/pengaduan', payload)
const myPengaduan = () => api.get('/api/pengaduan')
const detail = (id) => api.get(`/api/pengaduan/${id}`)
const update = (id, payload) => api.put(`/api/pengaduan/${id}`, payload)
const addRespon = (id, payload) => api.post(`/api/pengaduan/${id}/respon`, payload)
const finish = (id) => api.patch(`/api/pengaduan/${id}/selesai`)

export default { create, myPengaduan, detail, update, addRespon, finish }
