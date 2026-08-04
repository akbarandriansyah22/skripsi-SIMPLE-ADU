import api from './axios'

const login = (credentials) => api.post('/api/auth/login', credentials)
const register = (payload) => api.post('/api/auth/register', payload)
const profile = () => api.get('/api/auth/profile')
const changePassword = (payload) => api.patch('/api/auth/change-password', payload)

export default { login, register, profile, changePassword }
