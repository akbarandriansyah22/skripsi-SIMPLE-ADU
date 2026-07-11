import api from './axios'

const mine = () => api.get('/api/notifikasi')
const unread = () => api.get('/api/notifikasi/unread')
const markAsRead = (id) => api.patch(`/api/notifikasi/${id}/read`)

export default { mine, unread, markAsRead }
