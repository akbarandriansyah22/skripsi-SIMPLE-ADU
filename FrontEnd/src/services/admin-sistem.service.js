import api from './axios'

const dashboard = () => api.get('/api/admin-sistem/dashboard')
const users = () => api.get('/api/admin-sistem/users')
const user = (id) => api.get(`/api/admin-sistem/users/${id}`)
const units = () => api.get('/api/admin-sistem/units')
const createUser = (payload) => api.post('/api/admin-sistem/users', payload)
const updateUser = (id, payload) => api.patch(`/api/admin-sistem/users/${id}`, payload)
const setUserStatus = (id, payload) => api.patch(`/api/admin-sistem/users/${id}/status`, payload)
const resetPassword = (id, payload) => api.patch(`/api/admin-sistem/users/${id}/reset-password`, payload)
const createUnit = (payload) => api.post('/api/admin-sistem/units', payload)
const updateUnit = (id, payload) => api.patch(`/api/admin-sistem/units/${id}`, payload)
const categories = () => api.get('/api/admin-sistem/categories')
const createCategory = (payload) => api.post('/api/admin-sistem/categories', payload)
const updateCategory = (id, payload) => api.patch(`/api/admin-sistem/categories/${id}`, payload)
const setUnitStatus = (id, payload) => api.patch(`/api/admin-sistem/units/${id}/status`, payload)
const setCategoryStatus = (id, payload) => api.patch(`/api/admin-sistem/categories/${id}/status`, payload)
const downloadImportTemplate = () => api.get('/api/admin-sistem/mahasiswa/template', { responseType: 'blob' })
const importMahasiswa = (payload) => api.post('/api/admin-sistem/mahasiswa/import', payload)

export default {
  dashboard,
  users,
  user,
  createUser,
  updateUser,
  setUserStatus,
  resetPassword,
  units,
  createUnit,
  updateUnit,
  categories,
  createCategory,
  updateCategory,
  setUnitStatus,
  setCategoryStatus,
  downloadImportTemplate,
  importMahasiswa,
}
