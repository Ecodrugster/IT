export default defineNuxtPlugin(async () => {
  const { initAuth } = useAuth()
  
  // Initialize auth only on client side
  if (process.client) {
    try {
      initAuth()
      
      const userStore = useUserStore()
      
      const waitAuth = () => new Promise((resolve) => {
        if (!userStore.loading) return resolve(true)
        const unwatch = watch(() => userStore.loading, (loading) => {
          if (!loading) {
            unwatch()
            resolve(true)
          }
        })
        // Timeout after 5s just in case
        setTimeout(() => { unwatch(); resolve(false) }, 5000)
      })

      await waitAuth()
      
      // If after waiting we are not logged in and on a protected route, redirect
      const route = useRoute()
      if (!userStore.isLoggedIn && !['/login', '/register'].includes(route.path)) {
        await navigateTo('/login')
      }
    } catch (e) {
      console.error('Auth initialization error:', e)
    }
  }
})
