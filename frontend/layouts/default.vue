<script setup>
const userStore = useUserStore()
const notificationStore = useNotificationStore()
const { logout } = useAuth()
const { fetchApi: api } = useApi()
const router = useRouter()
const route = useRoute()

const searchQuery = ref('')
const handleSearch = () => {
  if (!searchQuery.value.trim()) return
  router.push({ path: '/search', query: { q: searchQuery.value } })
  searchQuery.value = ''
}

const handleLogout = async () => {
  await logout()
  router.push('/login')
}

const isNavActive = (path, exact = false) => {
  if (exact) return route.path === path
  return route.path === path || route.path.startsWith(`${path}/`)
}

const navClasses = (path, options = {}) => {
  const exact = options.exact || false
  const accent = options.accent || 'blue'
  const active = isNavActive(path, exact)

  const accentLine = {
    blue: 'after:bg-blue-500',
    yellow: 'after:bg-yellow-400',
    green: 'after:bg-emerald-400',
    red: 'after:bg-red-400'
  }[accent] || 'after:bg-blue-500'

  const inactiveText = {
    blue: 'text-slate-400 hover:text-white',
    yellow: 'text-yellow-500 hover:text-yellow-400',
    green: 'text-green-500 hover:text-green-400',
    red: 'text-red-500 hover:text-red-400'
  }[accent] || 'text-slate-400 hover:text-white'

  const activeText = {
    blue: 'text-white',
    yellow: 'text-yellow-300',
    green: 'text-green-300',
    red: 'text-red-300'
  }[accent] || 'text-white'

  return [
    'relative text-sm font-medium transition-all duration-200 pb-2',
    "after:content-[''] after:absolute after:left-0 after:-bottom-[18px] after:h-[2px] after:w-full after:rounded-full after:origin-left after:transition-all after:duration-300",
    accentLine,
    active ? `after:opacity-100 after:scale-x-100 ${activeText}` : `after:opacity-0 after:scale-x-0 ${inactiveText}`
  ]
}

onMounted(async () => {
  if (!userStore.isLoggedIn) return
  try {
    const profile = await api('/profile')
    userStore.setProfile(profile)
  } catch (e) {
    console.error('[Layout] Failed to refresh profile:', e)
  }
})
</script>

<template>
  <div class="min-h-screen bg-slate-950 text-slate-200 selection:bg-blue-500/30">
    <header class="sticky top-0 z-50 bg-slate-950/80 backdrop-blur-md border-b border-white/5">
      <div class="max-w-6xl mx-auto px-4 h-16 flex items-center justify-between">
        <NuxtLink to="/" class="flex items-center space-x-4 cursor-pointer">
          <div class="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center font-bold text-white shadow-lg shadow-blue-600/30">
            S
          </div>
          <span class="text-lg font-bold bg-gradient-to-r from-white to-slate-400 bg-clip-text text-transparent">
            ITSTEP Social
          </span>
        </NuxtLink>

        <div class="hidden md:flex items-center bg-slate-900 border border-white/5 rounded-full px-4 py-1.5 w-96">
          <span class="text-slate-500 mr-2">🔎</span>
          <input
            v-model="searchQuery"
            @keyup.enter="handleSearch"
            type="text"
            placeholder="Поиск..."
            class="bg-transparent border-none focus:outline-none text-sm w-full placeholder-slate-600 text-white"
          />
        </div>

        <div class="flex items-center space-x-6">
          <NuxtLink to="/" :class="navClasses('/', { exact: true })">Лента</NuxtLink>
          <NuxtLink to="/clubs" :class="navClasses('/clubs')">Клубы</NuxtLink>
          <NuxtLink to="/schedule" :class="navClasses('/schedule')">Расписание</NuxtLink>
          <NuxtLink to="/store" :class="navClasses('/store')">Магазин</NuxtLink>
          <NuxtLink to="/guide" :class="[...navClasses('/guide', { accent: 'yellow' }), 'font-bold', 'flex items-center']">
            <span class="mr-1">🧭</span> Навигатор
          </NuxtLink>
          <NuxtLink to="/chat" :class="navClasses('/chat')">Чат</NuxtLink>

          <NuxtLink v-if="userStore.isTeacherLike" to="/teacher/grades" :class="[...navClasses('/teacher/grades', { accent: 'blue' }), 'font-bold', 'flex items-center']">
            <span class="mr-1">📖</span> Журнал
          </NuxtLink>
          <NuxtLink v-if="userStore.role === 'student'" to="/profile/grades" :class="[...navClasses('/profile/grades', { accent: 'green' }), 'font-bold', 'flex items-center']">
            <span class="mr-1">📊</span> Мои оценки
          </NuxtLink>
          <NuxtLink v-if="userStore.isAdmin" to="/admin" :class="[...navClasses('/admin', { accent: 'red' }), 'font-bold', 'flex items-center']">
            <span class="mr-1">⚙️</span> Админ
          </NuxtLink>
        </div>

        <ClientOnly>
          <!-- Баланс студента (Монеты и Звезды) -->
          <NuxtLink v-if="userStore.user && userStore.role === 'student'" to="/store" class="hidden md:flex items-center space-x-3 bg-slate-900 border border-white/5 px-3 py-1.5 rounded-full text-xs text-slate-300 hover:border-yellow-500/30 transition-all hover:bg-slate-900/80 shadow-inner shadow-black/40" title="Магазин колледжа: проверить баланс">
            <div class="flex items-center space-x-1 hover:scale-105 transition-transform">
              <span>💰</span>
              <span class="font-bold text-yellow-400">{{ userStore.profile?.coins || 0 }}</span>
            </div>
            <div class="w-px h-3 bg-white/10"></div>
            <div class="flex items-center space-x-1 hover:scale-105 transition-transform">
              <span>🌟</span>
              <span class="font-bold text-amber-500">{{ userStore.profile?.stars || 0 }}</span>
            </div>
          </NuxtLink>

          <NuxtLink to="/notifications" class="p-2 hover:bg-white/5 rounded-full transition-all relative group" title="Уведомления">
            <span class="text-xl transition-transform group-hover:scale-110 group-hover:rotate-12 inline-block">📬</span>
            <span v-if="notificationStore.unreadCount > 0" class="absolute top-1.5 right-1.5 w-2 h-2 bg-red-500 rounded-full border border-slate-950 animate-pulse"></span>
          </NuxtLink>

          <div v-if="userStore.user" class="flex items-center space-x-3">
            <div class="hidden md:block text-right">
              <div class="text-xs font-semibold text-white truncate w-24">
                {{ userStore.user.displayName || userStore.user.email }}
              </div>
              <div class="text-[10px] text-slate-500 uppercase tracking-wider">
                {{ userStore.roleLabel }}
              </div>
            </div>
            <NuxtLink to="/profile" class="relative group">
              <img
                v-if="userStore.user.photoURL"
                :src="userStore.user.photoURL"
                class="w-9 h-9 rounded-full border border-white/10 group-hover:border-blue-500/50 transition-all shadow-lg"
              />
              <div v-else class="w-9 h-9 rounded-full bg-slate-800 border border-white/10 flex items-center justify-center text-xs group-hover:border-blue-500/50">
                {{ (userStore.user.displayName || userStore.user.email || 'U')[0].toUpperCase() }}
              </div>
            </NuxtLink>
          </div>
        </ClientOnly>
      </div>
    </header>

    <main>
      <slot />
    </main>
  </div>
</template>
