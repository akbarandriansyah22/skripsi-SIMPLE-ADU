import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  headers: { 'Content-Type': 'application/json' },
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('simpelToken')
  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use((response) => {
  if (
    response.data &&
    typeof response.data === 'object' &&
    Object.prototype.hasOwnProperty.call(response.data, 'success') &&
    Object.prototype.hasOwnProperty.call(response.data, 'data')
  ) {
    return {
      ...response,
      data: response.data.data,
      message: response.data.message,
      success: response.data.success,
    }
  }

  return response
})

export default api
