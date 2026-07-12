import { defineStore } from 'pinia'
import authService from '../services/auth.service'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    token: localStorage.getItem('simpelToken') || '',
    role: localStorage.getItem('simpelRole') || '',
    status: 'idle',
  }),
  getters: {
    isAuthenticated: (state) => !!state.token,
    normalizedRole: (state) => String(state.role || '').toLowerCase().replace(/\s+/g, '_'),
    isMahasiswa: (state) => String(state.role || '').toLowerCase() === 'mahasiswa',
    isAdmin: (state) => {
      const role = String(state.role || '').toLowerCase().replace(/\s+/g, '_')
      return role === 'petugas' || role.startsWith('admin')
    },
    isPimpinan: (state) => String(state.role || '').toLowerCase().replace(/\s+/g, '_').startsWith('pimpinan'),
  },
  actions: {
    async login(credentials) {
      this.status = 'loading'
      const response = await authService.login(credentials)
      this.token = response.data.token
      this.user = response.data.user
      this.role = this.user?.role || ''
      localStorage.setItem('simpelToken', this.token)
      localStorage.setItem('simpelRole', this.role)
      localStorage.setItem('simpelUser', JSON.stringify(this.user))
      this.status = 'success'
      return response
    },
    logout() {
      this.user = null
      this.token = ''
      this.role = ''
      localStorage.removeItem('simpelToken')
      localStorage.removeItem('simpelRole')
      localStorage.removeItem('simpelUser')
    },
    loadFromStorage() {
      const token = localStorage.getItem('simpelToken')
      const role = localStorage.getItem('simpelRole')
      const user = localStorage.getItem('simpelUser')
      this.token = token || ''
      this.role = role || ''
      try {
        this.user = user ? JSON.parse(user) : null
      } catch {
        this.user = null
      }
    },
  },
})
