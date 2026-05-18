<template>
  <div class="max-w-4xl mx-auto py-8 px-4">
    <ClientOnly>
      <div v-if="userStore.user" class="bg-slate-900 border border-white/5 rounded-2xl shadow-xl overflow-hidden">
        <div class="h-48 bg-gradient-to-r from-blue-600 to-indigo-600"></div>

        <div class="px-8 pb-8">
          <div class="relative -mt-16 mb-6">
            <img
              v-if="avatarUrl"
              :src="avatarUrl"
              class="w-32 h-32 rounded-2xl bg-slate-800 border-4 border-slate-900 shadow-2xl object-cover"
            />
            <div v-else class="w-32 h-32 rounded-2xl bg-slate-800 border-4 border-slate-900 shadow-2xl flex items-center justify-center text-4xl font-bold text-slate-500">
              {{ initials }}
            </div>
          </div>

          <div class="flex flex-col md:flex-row md:items-start md:justify-between gap-4">
            <div>
              <h1 class="text-3xl font-bold text-white mb-1">{{ displayName }}</h1>
              <p class="text-slate-400">{{ userStore.user?.email }} • {{ userStore.roleLabel }}</p>
              <p v-if="currentGroup" class="text-xs text-slate-500 mt-1">Group: {{ currentGroup }}</p>
            </div>

            <button
              @click="openEditModal"
              class="px-6 py-2 bg-white/5 hover:bg-white/10 border border-white/10 rounded-lg text-sm font-medium transition-all"
            >
              Редактировать Профиль
            </button>
          </div>

          <div class="mt-8 grid grid-cols-1 md:grid-cols-2 gap-8">
            <div class="space-y-6">
              <div>
                <h3 class="text-slate-300 font-semibold mb-3">О себе</h3>
                <p class="text-slate-400 text-sm leading-relaxed">
                  Твой профиль ITSTEP. Группа указана здесь и используется для расписания и журнала.
                </p>
              </div>
            </div>

            <div class="bg-white/5 rounded-xl p-6 border border-white/5">
              <h3 class="text-slate-300 font-semibold mb-4">Статистика</h3>
              <div class="grid grid-cols-2 gap-4">
                <div class="text-center p-4 bg-slate-950/50 rounded-lg">
                  <div class="text-2xl font-bold text-white">{{ stats.posts }}</div>
                  <div class="text-xs text-slate-500 uppercase">Посты</div>
                </div>
                <div class="text-center p-4 bg-slate-950/50 rounded-lg">
                  <div class="text-2xl font-bold text-white">{{ stats.comments }}</div>
                  <div class="text-xs text-slate-500 uppercase">Комментарии</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="text-center py-20 text-slate-500">Загрузка профиля...</div>

      <div v-if="showEditModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm">
        <div class="bg-slate-900 border border-white/10 rounded-2xl p-8 max-w-md w-full shadow-2xl">
          <h2 class="text-2xl font-bold text-white mb-6">Настройки профиля</h2>
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-slate-400 mb-2">Имя</label>
              <input v-model="editForm.displayName" type="text" class="w-full bg-slate-800 border-none rounded-lg px-4 py-3 text-white focus:ring-2 focus:ring-blue-500/50" />
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-400 mb-2">Фото</label>
              <input v-model="editForm.photoURL" type="text" class="w-full bg-slate-800 border-none rounded-lg px-4 py-3 text-white focus:ring-2 focus:ring-blue-500/50" />
            </div>
          </div>
          <div class="flex space-x-4 mt-8">
            <button @click="showEditModal = false" class="flex-grow py-3 bg-white/5 hover:bg-white/10 text-white rounded-xl transition-all">Cancel</button>
            <button @click="saveProfile" class="flex-grow py-3 bg-blue-600 hover:bg-blue-500 text-white rounded-xl font-semibold shadow-lg shadow-blue-600/20 transition-all">Save</button>
          </div>
        </div>
      </div>
    </ClientOnly>
  </div>
</template>

<script setup>
import { updateProfile } from 'firebase/auth'

definePageMeta({
  middleware: 'auth'
})

const { $auth } = useNuxtApp()
const userStore = useUserStore()
const { fetchApi: api } = useApi()

const showEditModal = ref(false)
const editForm = reactive({
  displayName: '',
  photoURL: ''
})

const stats = ref({ posts: 0, comments: 0 })

const displayName = computed(() => {
  return userStore.profile?.display_name || userStore.user?.displayName || userStore.user?.email || 'Пользователь'
})

const avatarUrl = computed(() => {
  return userStore.profile?.photo_url || userStore.user?.photoURL || ''
})

const initials = computed(() => {
  const base = displayName.value || 'U'
  return base[0]?.toUpperCase() || 'U'
})

const currentGroup = computed(() => {
  return userStore.profile?.group_name || userStore.profile?.group || ''
})

const fetchStats = async () => {
  try {
    const data = await api('/profile/stats')
    stats.value = data || { posts: 0, comments: 0 }
  } catch (e) {
    console.error('Failed to fetch stats:', e)
  }
}

const openEditModal = () => {
  editForm.displayName = userStore.profile?.display_name || userStore.user?.displayName || ''
  editForm.photoURL = userStore.profile?.photo_url || userStore.user?.photoURL || ''
  showEditModal.value = true
}

const saveProfile = async () => {
  try {
    if ($auth?.currentUser) {
      await updateProfile($auth.currentUser, {
        displayName: editForm.displayName,
        photoURL: editForm.photoURL
      })
    }

    await api('/profile', {
      method: 'PUT',
      body: {
        display_name: editForm.displayName,
        photo_url: editForm.photoURL
      }
    })

    const profile = await api('/profile')
    userStore.setProfile(profile)
    if ($auth?.currentUser) {
      userStore.setUser({ ...$auth.currentUser })
    }

    showEditModal.value = false
    alert('Profile updated')
  } catch (e) {
    alert('Failed to update profile: ' + (e?.message || 'unknown error'))
  }
}

onMounted(() => {
  fetchStats()
})
</script>
