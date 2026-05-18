import type { Auth } from 'firebase/auth'

export const useApi = () => {
  const userStore = useUserStore()
  const config = useRuntimeConfig()
  const baseUrl = config.public.apiBaseUrl || 'http://localhost:8080/api/v1'

  const ensureToken = async (forceRefresh = false) => {
    if (!process.client) return userStore.token

    const { $auth } = useNuxtApp() as { $auth: Auth | null }
    const currentUser = $auth?.currentUser
    if (!currentUser) return userStore.token

    if (!userStore.token || forceRefresh) {
      try {
        const freshToken = await currentUser.getIdToken(forceRefresh)
        userStore.setToken(freshToken)
      } catch (e) {
        console.error('[API] Failed to refresh Firebase token:', e)
      }
    }

    return userStore.token
  }

  const request = async <T>(url: string, options: any = {}) => {
    await ensureToken(false)

    const headers: Record<string, string> = {
      ...options.headers
    }

    if (userStore.token) {
      headers['Authorization'] = `Bearer ${userStore.token}`
    }

    return $fetch<T>(`${baseUrl}${url}`, {
      ...options,
      headers
    })
  }

  const fetchApi = async <T = any>(url: string, options: any = {}) => {
    try {
      return await request<T>(url, options)
    } catch (e: any) {
      const status = e?.status || e?.response?.status
      if (status === 401 && process.client) {
        await ensureToken(true)
        return await request<T>(url, options)
      }
      throw e
    }
  }

  return {
    fetchApi
  }
}
