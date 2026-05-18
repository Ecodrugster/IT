import { defineStore } from 'pinia'

export interface Notification {
  id: string
  user_id: string
  type: 'message' | 'grade' | 'attendance' | 'news' | 'system'
  title: string
  message: string
  link: string
  read: boolean
  created_at: string
}

export const useNotificationStore = defineStore('notifications', {
  state: () => ({
    notifications: [] as Notification[],
    unreadCount: 0,              // Total combined unread count (chats + alerts)
    unreadNotificationsCount: 0, // Unread firestore notifications only
    unreadChatsCount: 0,         // Unread MongoDB chat messages only
    loading: false
  }),
  actions: {
    async fetchNotifications() {
      const { fetchApi: api } = useApi()
      this.loading = true
      try {
        const data = await api<Notification[]>('/notifications')
        this.notifications = data || []
      } catch (e) {
        console.error('[NotificationStore] Failed to fetch notifications:', e)
      } finally {
        this.loading = false
      }
    },
    async fetchUnreadCount() {
      const { fetchApi: api } = useApi()
      try {
        const data = await api<{ total_unread: number; notifications: number; chats: number }>('/notifications/unread-count')
        this.unreadCount = Number(data?.total_unread || 0)
        this.unreadNotificationsCount = Number(data?.notifications || 0)
        this.unreadChatsCount = Number(data?.chats || 0)
      } catch (e) {
        console.error('[NotificationStore] Failed to fetch unread count:', e)
      }
    },
    async markAsRead(id: string) {
      const { fetchApi: api } = useApi()
      try {
        await api(`/notifications/${encodeURIComponent(id)}/read`, { method: 'POST' })
        // Update local state
        const notif = this.notifications.find(n => n.id === id)
        if (notif && !notif.read) {
          notif.read = true
          if (this.unreadNotificationsCount > 0) this.unreadNotificationsCount--
          if (this.unreadCount > 0) this.unreadCount--
        }
      } catch (e) {
        console.error('[NotificationStore] Failed to mark notification as read:', e)
      }
    },
    async markAllAsRead() {
      const { fetchApi: api } = useApi()
      try {
        await api('/notifications/read-all', { method: 'POST' })
        // Update local state
        this.notifications.forEach(n => {
          n.read = true
        })
        this.unreadNotificationsCount = 0
        this.unreadCount = this.unreadChatsCount // Only chat messages remain unread
      } catch (e) {
        console.error('[NotificationStore] Failed to mark all notifications as read:', e)
      }
    },
    clear() {
      this.notifications = []
      this.unreadCount = 0
      this.unreadNotificationsCount = 0
      this.unreadChatsCount = 0
    },
    setUnreadCount(count: number) {
      this.unreadCount = count
    }
  }
})
