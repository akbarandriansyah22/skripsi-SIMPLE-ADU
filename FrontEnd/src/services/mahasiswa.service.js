import api from './axios'

const profile = () => api.get('/api/mahasiswa/profile')
const updateProfile = (payload) => api.put('/api/mahasiswa/profile', payload)

export default { profile, updateProfile }
