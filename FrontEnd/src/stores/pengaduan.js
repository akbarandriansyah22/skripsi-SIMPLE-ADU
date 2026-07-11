import { defineStore } from 'pinia'
import pengaduanService from '../services/pengaduan.service'

export const usePengaduanStore = defineStore('pengaduan', {
  state: () => ({
    list: [],
    detail: null,
    status: 'idle',
  }),
  actions: {
    async loadMyPengaduan() {
      this.status = 'loading'
      const response = await pengaduanService.myPengaduan()
      this.list = response.data || []
      this.status = 'success'
      return response
    },
    async loadDetail(id) {
      this.status = 'loading'
      const response = await pengaduanService.detail(id)
      this.detail = response.data || null
      this.status = 'success'
      return response
    },
  },
})
