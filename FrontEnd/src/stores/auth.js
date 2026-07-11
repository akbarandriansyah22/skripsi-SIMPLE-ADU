import { defineStore } from 'pinia'
import { useRouter } from 'vue-router'
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
    isMahasiswa: (state) => state.role === 'mahasiswa',
    isAdmin: (state) => state.role.startsWith('admin'),
    isPimpinan: (state) => state.role.startsWith('pimpinan'),
  },
  actions: {
    async login(credentials) {
      this.status = 'loading'
      const response = await authService.login(credentials)
      this.token = response.data.token
      this.role = response.data.role
      this.user = response.data.user
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
      this.user = user ? JSON.parse(user) : null
    },
  },
})
