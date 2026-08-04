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
    isAdminSistem: (state) => String(state.role || '').toLowerCase().replace(/\s+/g, '_') === 'admin_sistem',
    isAdminFakultas: (state) => {
      const role = String(state.role || '').toLowerCase().replace(/\s+/g, '_')
      return role === 'admin_fakultas' || role === 'petugas'
    },
    isAdmin: (state) => {
      const role = String(state.role || '').toLowerCase().replace(/\s+/g, '_')
      return role === 'admin_sistem' || role === 'admin_fakultas' || role === 'petugas'
    },
    isPimpinan: (state) => String(state.role || '').toLowerCase().replace(/\s+/g, '_').startsWith('pimpinan'),
    isKasubag: (state) => String(state.role || '').toLowerCase().replace(/\s+/g, '_') === 'kasubag',
  },
  actions: {
    async login(credentials) {
      this.status = 'loading'
      try {
        const response = await authService.login({
          ...credentials,
          email: String(credentials.email || '').trim(),
        })
        this.token = response.data.token
        this.user = response.data.user
        this.role = this.user?.role || ''
        localStorage.setItem('simpelToken', this.token)
        localStorage.setItem('simpelRole', this.role)
        localStorage.setItem('simpelUser', JSON.stringify(this.user))
        this.status = 'success'
        return response
      } catch (error) {
        this.status = 'error'
        throw error
      }
    },
    logout() {
      this.user = null
      this.token = ''
      this.role = ''
      localStorage.removeItem('simpelToken')
      localStorage.removeItem('simpelRole')
      localStorage.removeItem('simpelUser')
      sessionStorage.clear()
      this.status = 'idle'
    },
    markPasswordChanged() {
      if (this.user) this.user.password_must_change = false
      localStorage.setItem('simpelUser', JSON.stringify(this.user))
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
        this.logout()
      }
    },
  },
})
