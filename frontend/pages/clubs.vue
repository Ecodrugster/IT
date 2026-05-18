<template>
  <div class="max-w-6xl mx-auto py-8 px-4">
    <div class="flex flex-col md:flex-row justify-between items-start md:items-center mb-8 gap-4">
      <div>
        <h1 class="text-3xl font-bold text-white mb-2">Студенческие клубы</h1>
        <p class="text-slate-400 text-sm">Найдите клуб по интересам или отправьте заявку на создание своего</p>
      </div>
      <button
        @click="showCreateModal = true"
        class="bg-blue-600 hover:bg-blue-500 text-white px-6 py-2.5 rounded-xl font-semibold transition-all shadow-lg shadow-blue-600/20 flex items-center gap-2"
      >
        <span>+</span> Создать клуб
      </button>
    </div>

    <div class="flex space-x-1 bg-slate-900/50 p-1 rounded-xl w-fit mb-8 border border-white/5">
      <button
        v-for="tab in ['all', 'mine']"
        :key="tab"
        @click="activeTab = tab"
        class="px-6 py-2 rounded-lg text-sm font-medium transition-all"
        :class="activeTab === tab ? 'bg-blue-600 text-white shadow-lg' : 'text-slate-400 hover:text-white'"
      >
        {{ tab === 'all' ? 'Все клубы' : 'Мои клубы' }}
      </button>
    </div>

    <div v-if="showCreateModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm">
      <div class="bg-slate-900 border border-white/10 rounded-2xl p-8 max-w-md w-full shadow-2xl">
        <h2 class="text-2xl font-bold text-white mb-6">Новый клуб</h2>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-slate-400 mb-2">Название</label>
            <input v-model="newClub.name" type="text" placeholder="Клуб робототехники" class="w-full bg-slate-800 rounded-lg px-4 py-3 text-white" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-400 mb-2">Описание</label>
            <textarea v-model="newClub.description" rows="3" placeholder="О чем ваш клуб?" class="w-full bg-slate-800 rounded-lg px-4 py-3 text-white resize-none"></textarea>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-slate-400 mb-2">Иконка</label>
              <select v-model="newClub.icon" class="w-full bg-slate-800 rounded-lg px-4 py-3 text-white">
                <option v-for="i in ['💻', '🎨', '🚀', '♟️', '⚽', '📚', '🎶', '🎮']" :key="i">{{ i }}</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-slate-400 mb-2">Цвет</label>
              <select v-model="newClub.color" class="w-full bg-slate-800 rounded-lg px-4 py-3 text-white">
                <option value="bg-blue-600/20">Синий</option>
                <option value="bg-purple-600/20">Фиолетовый</option>
                <option value="bg-red-600/20">Красный</option>
                <option value="bg-green-600/20">Зеленый</option>
                <option value="bg-orange-600/20">Оранжевый</option>
              </select>
            </div>
          </div>
        </div>
        <div class="flex space-x-4 mt-8">
          <button @click="showCreateModal = false" class="flex-grow py-3 bg-white/5 hover:bg-white/10 text-white rounded-xl transition-all">Отмена</button>
          <button @click="handleCreateClub" class="flex-grow py-3 bg-blue-600 hover:bg-blue-500 text-white rounded-xl font-semibold transition-all">Отправить</button>
        </div>
      </div>
    </div>

    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div v-for="i in 3" :key="i" class="bg-slate-900 border border-white/5 rounded-xl p-6 shadow-xl animate-pulse">
        <div class="w-full h-32 bg-slate-800 rounded-lg mb-4"></div>
        <div class="h-5 w-3/4 bg-slate-800 rounded mb-2"></div>
        <div class="h-4 w-full bg-slate-800 rounded mb-4"></div>
        <div class="w-full h-10 bg-slate-800 rounded-lg"></div>
      </div>
    </div>

    <div v-else-if="filteredClubs.length === 0" class="text-center py-20 bg-slate-900/50 rounded-xl border border-white/5">
      <div class="text-6xl mb-4">👥</div>
      <h3 class="text-xl text-white font-medium mb-2">Ничего не найдено</h3>
      <p class="text-slate-500">Попробуйте сменить вкладку или создать свой клуб.</p>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div v-for="club in filteredClubs" :key="club.id" class="bg-slate-900 border border-white/5 rounded-xl p-6 shadow-xl flex flex-col group hover:border-blue-500/30 transition-all">
        <div class="w-full h-32 rounded-lg mb-4 flex items-center justify-center text-4xl group-hover:scale-105 transition-transform" :class="club.color || 'bg-blue-600/20'">
          {{ club.icon || '👥' }}
        </div>

        <div class="flex items-center justify-between mb-2">
          <h3 class="text-white font-semibold text-lg">{{ club.name }}</h3>
          <span class="px-2 py-0.5 rounded-md text-[10px] font-bold uppercase tracking-wider" :class="statusClass(club.status)">
            {{ club.status === 'pending' ? 'На модерации' : club.status === 'rejected' ? 'Отклонен' : 'Одобрен' }}
          </span>
        </div>

        <p class="text-slate-400 text-sm mb-4 flex-grow">{{ club.description }}</p>

        <div class="flex items-center justify-between mb-4">
          <div class="flex -space-x-2">
            <div v-for="i in Math.min(club.members?.length || 0, 3)" :key="i" class="w-7 h-7 rounded-full bg-slate-800 border-2 border-slate-900 flex items-center justify-center text-[10px] text-slate-400">
              👤
            </div>
            <div v-if="(club.members?.length || 0) > 3" class="w-7 h-7 rounded-full bg-slate-800 border-2 border-slate-900 flex items-center justify-center text-[10px] text-slate-500">
              +{{ club.members.length - 3 }}
            </div>
          </div>
          <span class="text-xs text-slate-500">{{ club.members?.length || 0 }} участников</span>
        </div>

        <button
          v-if="!isMember(club) && isApproved(club)"
          @click="joinClub(club)"
          :disabled="joining === club.id"
          class="w-full py-2.5 bg-blue-600 hover:bg-blue-500 text-white rounded-lg text-sm font-medium transition-all disabled:opacity-50"
        >
          {{ joining === club.id ? 'Вступаем...' : 'Вступить в клуб' }}
        </button>

        <button
          v-else-if="!isApproved(club)"
          disabled
          class="w-full py-2.5 bg-amber-500/10 border border-amber-500/20 text-amber-400 rounded-lg text-sm font-medium"
        >
          {{ club.status === 'rejected' ? 'Заявка отклонена' : 'На модерации' }}
        </button>

        <div v-else class="flex gap-2">
          <button
            disabled
            class="flex-grow py-2.5 bg-green-500/10 border border-green-500/20 text-green-400 rounded-lg text-sm font-medium"
          >
            Вы участник ✅
          </button>

          <button
            v-if="club.created_by === myId"
            @click="openEditModal(club)"
            class="px-3 py-2.5 bg-blue-600/10 hover:bg-blue-600 text-blue-500 hover:text-white rounded-lg text-sm transition-all border border-blue-500/20"
            title="Управлять клубом"
          >
            ⚙️
          </button>

          <button
            @click="leaveClub(club)"
            class="px-3 py-2.5 bg-red-500/10 hover:bg-red-500/20 text-red-500 rounded-lg text-sm transition-all"
            title="Покинуть клуб"
          >
            🚪
          </button>
        </div>
      </div>
    </div>

    <div v-if="showEditModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm">
      <div class="bg-slate-900 border border-white/10 rounded-2xl p-8 max-w-md w-full shadow-2xl">
        <h2 class="text-2xl font-bold text-white mb-6">Настройки клуба</h2>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-slate-400 mb-2">Название</label>
            <input v-model="editClubForm.name" type="text" class="w-full bg-slate-800 rounded-lg px-4 py-3 text-white" />
          </div>
          <div>
            <label class="block text-sm font-medium text-slate-400 mb-2">Описание</label>
            <textarea v-model="editClubForm.description" rows="3" class="w-full bg-slate-800 rounded-lg px-4 py-3 text-white resize-none"></textarea>
          </div>
        </div>
        <div class="flex flex-col gap-3 mt-8">
          <div class="flex gap-3">
            <button @click="showEditModal = false" class="flex-grow py-3 bg-white/5 hover:bg-white/10 text-white rounded-xl transition-all">Отмена</button>
            <button @click="handleUpdateClub" class="flex-grow py-3 bg-blue-600 hover:bg-blue-500 text-white rounded-xl font-semibold transition-all">Сохранить</button>
          </div>
          <button @click="handleDeleteClub" class="w-full py-3 bg-red-500/10 hover:bg-red-500 text-red-500 hover:text-white rounded-xl text-sm font-bold transition-all border border-red-500/20">Удалить клуб</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  middleware: 'auth'
})

const { fetchApi: api } = useApi()
const userStore = useUserStore()

const clubs = ref([])
const loading = ref(true)
const joining = ref(null)
const activeTab = ref('all')
const showCreateModal = ref(false)
const showEditModal = ref(false)

const myId = computed(() => userStore.user?.uid || userStore.profile?.uid || '')

const editClubForm = reactive({
  id: '',
  name: '',
  description: ''
})

const newClub = reactive({
  name: '',
  description: '',
  icon: '💻',
  color: 'bg-blue-600/20'
})

const statusClass = (status) => {
  switch (status) {
    case 'pending':
      return 'bg-amber-500/10 text-amber-400'
    case 'rejected':
      return 'bg-red-500/10 text-red-400'
    default:
      return 'bg-green-500/10 text-green-400'
  }
}

const isApproved = (club) => (club.status || 'approved') === 'approved'

const openEditModal = (club) => {
  editClubForm.id = club.id
  editClubForm.name = club.name
  editClubForm.description = club.description
  showEditModal.value = true
}

const fetchClubs = async () => {
  loading.value = true
  try {
    const data = await api('/clubs')
    clubs.value = data || []
  } catch (e) {
    console.error('Failed to fetch clubs:', e)
    clubs.value = []
  } finally {
    loading.value = false
  }
}

const filteredClubs = computed(() => {
  if (activeTab.value === 'all') return clubs.value
  return clubs.value.filter(c => isMember(c) || c.created_by === myId.value)
})

const isMember = (club) => {
  if (!club.members || !Array.isArray(club.members)) return false
  return club.members.includes(myId.value)
}

const joinClub = async (club) => {
  if (!myId.value) {
    alert('Не удалось определить пользователя. Перезайдите в аккаунт.')
    return
  }

  joining.value = club.id
  try {
    await api(`/clubs/${club.id}/join`, { method: 'POST' })
    if (!club.members) club.members = []
    if (!club.members.includes(myId.value)) {
      club.members.push(myId.value)
    }
  } catch (e) {
    alert('Ошибка при вступлении: ' + (e?.data?.error || e.message))
  } finally {
    joining.value = null
  }
}

const leaveClub = async (club) => {
  if (!myId.value) {
    alert('Не удалось определить пользователя. Перезайдите в аккаунт.')
    return
  }

  if (!confirm(`Вы действительно хотите покинуть клуб "${club.name}"?`)) return
  try {
    await api(`/clubs/${club.id}/leave`, { method: 'POST' })
    club.members = (club.members || []).filter(uid => uid !== myId.value)
  } catch (e) {
    alert('Ошибка: ' + (e?.data?.error || e.message))
  }
}

const handleCreateClub = async () => {
  if (!newClub.name.trim() || !newClub.description.trim()) {
    alert('Заполните все поля')
    return
  }

  try {
    const result = await api('/clubs', {
      method: 'POST',
      body: newClub
    })

    showCreateModal.value = false
    newClub.name = ''
    newClub.description = ''

    alert(result?.message || 'Заявка отправлена на модерацию администраторам')
    await fetchClubs()
  } catch (e) {
    alert('Ошибка создания клуба: ' + (e?.data?.error || e.message))
  }
}

const handleUpdateClub = async () => {
  try {
    const result = await api(`/clubs/${editClubForm.id}`, {
      method: 'PUT',
      body: { name: editClubForm.name, description: editClubForm.description }
    })
    showEditModal.value = false
    alert(result?.message || 'Клуб обновлен')
    await fetchClubs()
  } catch (e) {
    alert('Ошибка обновления: ' + (e?.data?.error || e.message))
  }
}

const handleDeleteClub = async () => {
  if (!confirm('Вы уверены? Клуб будет удален навсегда.')) return
  try {
    await api(`/clubs/${editClubForm.id}`, { method: 'DELETE' })
    showEditModal.value = false
    await fetchClubs()
  } catch (e) {
    alert('Ошибка удаления: ' + (e?.data?.error || e.message))
  }
}

onMounted(fetchClubs)
</script>

