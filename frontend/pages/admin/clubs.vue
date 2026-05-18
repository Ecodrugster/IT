<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h2 class="text-2xl font-bold text-white">Управление клубами</h2>
      <button
        @click="openModal()"
        class="rounded-xl bg-blue-600 px-6 py-2.5 font-semibold text-white shadow-lg shadow-blue-600/20 transition-all hover:bg-blue-500"
      >
        + Создать клуб
      </button>
    </div>

    <div v-if="errorMessage" class="rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-300">
      {{ errorMessage }}
    </div>

    <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
      <div v-for="club in clubs" :key="club.id" class="group flex gap-6 rounded-2xl border border-white/5 bg-slate-900 p-6 shadow-xl">
        <div class="flex h-24 w-24 items-center justify-center rounded-xl text-4xl" :class="club.color || 'bg-blue-600/10'">
          {{ club.icon || 'C' }}
        </div>
        <div class="flex-grow">
          <div class="mb-2 flex items-start justify-between">
            <div>
              <h3 class="font-bold text-white">{{ club.name }}</h3>
              <span
                class="mt-1 inline-block rounded-md px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider"
                :class="statusClass(club.status)"
              >
                {{ statusLabel(club.status) }}
              </span>
            </div>
            <div class="flex gap-2">
              <button @click="openModal(club)" class="text-sm text-slate-500 transition-all hover:text-white">Изменить</button>
              <button @click="deleteClub(club.id)" class="text-sm text-slate-500 transition-all hover:text-red-500">Удалить</button>
            </div>
          </div>
          <p class="mb-4 line-clamp-2 text-xs text-slate-400">{{ club.description }}</p>
          <div class="text-[10px] font-bold uppercase tracking-widest text-slate-500">
            {{ club.members?.length || 0 }} участников
          </div>
        </div>
      </div>
    </div>

    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/80 p-4 backdrop-blur-sm">
      <div class="w-full max-w-md rounded-2xl border border-white/10 bg-slate-900 p-8 shadow-2xl">
        <h2 class="mb-6 text-2xl font-bold text-white">{{ isEditing ? 'Редактирование клуба' : 'Новый клуб' }}</h2>
        <div class="space-y-4">
          <div>
            <label class="mb-2 block text-sm font-medium text-slate-400">Название</label>
            <input v-model="form.name" type="text" class="w-full rounded-lg border-none bg-slate-800 px-4 py-3 text-white" />
          </div>
          <div>
            <label class="mb-2 block text-sm font-medium text-slate-400">Описание</label>
            <textarea v-model="form.description" rows="3" class="w-full resize-none rounded-lg border-none bg-slate-800 px-4 py-3 text-white"></textarea>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="mb-2 block text-sm font-medium text-slate-400">Иконка</label>
              <input v-model="form.icon" type="text" placeholder="C" class="w-full rounded-lg border-none bg-slate-800 px-4 py-3 text-center text-white" />
            </div>
            <div>
              <label class="mb-2 block text-sm font-medium text-slate-400">Цвет</label>
              <input v-model="form.color" type="text" placeholder="bg-blue-600/20" class="w-full rounded-lg border-none bg-slate-800 px-4 py-3 text-xs text-white" />
            </div>
          </div>

          <div v-if="isEditing">
            <label class="mb-2 block text-sm font-medium text-slate-400">Статус</label>
            <select v-model="form.status" class="w-full rounded-lg border-none bg-slate-800 px-4 py-3 text-white">
              <option value="approved">одобрен</option>
              <option value="pending">на рассмотрении</option>
              <option value="rejected">отклонен</option>
            </select>
          </div>
        </div>

        <div class="mt-8 flex space-x-4">
          <button @click="showModal = false" class="flex-grow rounded-xl bg-white/5 py-3 text-white transition-all hover:bg-white/10">Отмена</button>
          <button @click="saveClub" class="flex-grow rounded-xl bg-blue-600 py-3 font-semibold text-white transition-all hover:bg-blue-500">Сохранить</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'admin',
  middleware: 'admin'
})

const { fetchApi: api } = useApi()
const clubs = ref([])
const errorMessage = ref('')
const showModal = ref(false)
const isEditing = ref(false)
const currentId = ref('')

const form = reactive({
  name: '',
  description: '',
  icon: 'C',
  color: 'bg-blue-600/20',
  status: 'approved'
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

const statusLabel = (status) => {
  if (status === 'pending') return 'на рассмотрении'
  if (status === 'rejected') return 'отклонен'
  return 'одобрен'
}

const fetchClubs = async () => {
  errorMessage.value = ''
  try {
    const data = await api('/admin/clubs')
    clubs.value = data || []
  } catch (e) {
    const status = e?.status || e?.response?.status
    const message = e?.data?.error || e?.message || 'Не удалось загрузить клубы'
    errorMessage.value = `Не удалось загрузить клубы (${status || 'без статуса'}): ${message}`
    clubs.value = []
  }
}

const openModal = (club = null) => {
  if (club) {
    isEditing.value = true
    currentId.value = club.id
    Object.assign(form, {
      name: club.name || '',
      description: club.description || '',
      icon: club.icon || 'C',
      color: club.color || 'bg-blue-600/20',
      status: club.status || 'approved'
    })
  } else {
    isEditing.value = false
    currentId.value = ''
    Object.assign(form, {
      name: '',
      description: '',
      icon: 'C',
      color: 'bg-blue-600/20',
      status: 'approved'
    })
  }
  showModal.value = true
}

const saveClub = async () => {
  try {
    const method = isEditing.value ? 'PUT' : 'POST'
    const url = isEditing.value ? `/admin/clubs/${currentId.value}` : '/admin/clubs'

    await api(url, {
      method,
      body: {
        name: form.name,
        description: form.description,
        icon: form.icon,
        color: form.color,
        ...(isEditing.value ? { status: form.status } : {})
      }
    })

    showModal.value = false
    await fetchClubs()
  } catch (e) {
    alert('Ошибка: ' + (e?.data?.error || e.message))
  }
}

const deleteClub = async (id) => {
  if (!confirm('Вы уверены? Клуб будет удален для всех участников.')) return
  try {
    await api(`/admin/clubs/${id}`, { method: 'DELETE' })
    await fetchClubs()
  } catch (e) {
    alert('Ошибка удаления: ' + (e?.data?.error || e.message))
  }
}

onMounted(fetchClubs)
</script>
