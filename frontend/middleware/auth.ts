export default defineNuxtRouteMiddleware((to, from) => {
  const userStore = useUserStore()

  // Skip check while loading auth state on initial load
  if (userStore.loading) {
    console.log('[Middleware] Auth is loading, skipping check...')
    return
  }

  console.log('[Middleware] Checking auth. LoggedIn:', userStore.isLoggedIn)
  if (!userStore.isLoggedIn && to.path !== '/login' && to.path !== '/register') {
    return navigateTo('/login')
  }
})
