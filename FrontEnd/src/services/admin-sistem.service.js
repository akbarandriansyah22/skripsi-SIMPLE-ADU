import api from './axios'

const dashboard = () => api.get('/api/admin-sistem/dashboard')
const users = () => api.get('/api/admin-sistem/users')
const units = () => api.get('/api/admin-sistem/units')
const createUser = (payload) => api.post('/api/admin-sistem/users', payload)
const updateUser = (id, payload) => api.patch(`/api/admin-sistem/users/${id}`, payload)
const setUserStatus = (id, payload) => api.patch(`/api/admin-sistem/users/${id}/status`, payload)
const resetPassword = (id, payload) => api.patch(`/api/admin-sistem/users/${id}/reset-password`, payload)

export default { dashboard, users, units, createUser, updateUser, setUserStatus, resetPassword }
