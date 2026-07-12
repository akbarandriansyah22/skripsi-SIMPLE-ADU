import { defineStore } from 'pinia'
import notifikasiService from '../services/notifikasi.service'

export const useNotifikasiStore = defineStore('notifikasi', {
  state: () => ({
    items: [],
    unreadCount: 0,
    status: 'idle',
  }),
  actions: {
    async loadNotifications() {
      this.status = 'loading'
      const response = await notifikasiService.mine()
      this.items = response.data?.data || response.data || []
      this.unreadCount = this.items.filter((item) => !item.is_read && !item.read_at).length
      this.status = 'success'
      return response
    },
    async markRead(id) {
      await notifikasiService.markAsRead(id)
      const item = this.items.find((notification) => notification.id === id)
      if (item) {
        item.is_read = true
        item.read_at = new Date().toISOString()
      }
      this.unreadCount = this.items.filter((item) => !item.is_read && !item.read_at).length
    },
  },
})
