export default defineNuxtPlugin(() => {
  const userStore = useUserStore()
  const notificationStore = useNotificationStore()

  let pollTimer: ReturnType<typeof setInterval> | null = null

  const stopPolling = () => {
    if (pollTimer) {
      clearInterval(pollTimer)
      pollTimer = null
    }
  }

  const loadUnreadCount = async () => {
    if (!userStore.user) {
      notificationStore.clear()
      return
    }

    try {
      await notificationStore.fetchUnreadCount()
    } catch (e) {
      console.error('[Notifications Plugin] Failed to load unread count:', e)
    }
  }

  const startPolling = () => {
    stopPolling()
    loadUnreadCount()
    pollTimer = setInterval(loadUnreadCount, 5000)
  }

  watch(
    () => userStore.user,
    (user) => {
      if (user) {
        startPolling()
      } else {
        stopPolling()
        notificationStore.clear()
      }
    },
    { immediate: true }
  )
})
