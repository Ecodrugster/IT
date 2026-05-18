export default defineNuxtRouteMiddleware(async () => {
  const userStore = useUserStore()

  if (userStore.loading) {
    return
  }

  if (!userStore.isLoggedIn) {
    return navigateTo('/login')
  }

  if (process.client) {
    try {
      const { fetchApi } = useApi()
      const profile = await fetchApi('/profile')
      userStore.setProfile(profile)
    } catch (e) {
      console.error('[Admin Middleware] Failed to load profile:', e)
    }
  }

  const role = userStore.profile?.role
  const canAccessAdmin = role === 'admin'

  if (!canAccessAdmin) {
    console.warn('[Admin Middleware] Access denied. Admin role required.')
    return navigateTo('/')
  }
})
