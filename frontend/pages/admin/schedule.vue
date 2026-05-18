<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h2 class="text-2xl font-bold text-white">Управление расписанием</h2>
      <button
        @click="openModal()"
        class="bg-blue-600 hover:bg-blue-500 text-white px-6 py-2.5 rounded-xl font-semibold transition-all shadow-lg shadow-blue-600/20"
      >
        + Добавить пару
      </button>
    </div>

    <div v-if="errorMessage" class="bg-red-500/10 border border-red-500/30 text-red-300 rounded-xl px-4 py-3 text-sm">
      {{ errorMessage }}
    </div>

    <div v-if="loading" class="space-y-4">
      <div v-for="i in 4" :key="i" class="h-24 bg-slate-900 border border-white/5 rounded-2xl animate-pulse"></div>
    </div>

    <div v-else class="space-y-4">
      <div v-if="schedule.length === 0" class="bg-slate-900 border border-dashed border-white/10 rounded-2xl p-12 text-center text-slate-500">
        Расписание пока пустое.
      </div>

      <div
        v-for="item in schedule"
        :key="item.id"
        class="bg-slate-900 border border-white/5 rounded-2xl p-5 flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4"
      >
        <div>
          <div class="text-white font-bold">{{ item.subject }}</div>
          <div class="text-sm text-slate-400 mt-1">
            {{ dayName(item.day_of_week) }}, Пара {{ item.pair_number }} • {{ item.starts_at }}{{ item.ends_at ? `-${item.ends_at}` : '' }}
          </div>
          <div class="text-xs text-slate-500 mt-2">
            Группа: {{ item.group_name }} • Преподаватель: {{ item.teacher_name || item.teacher_id }}
            <span v-if="item.room"> • Аудитория {{ item.room }}</span>
          </div>
        </div>

        <div class="flex gap-2">
          <button @click="openModal(item)" class="px-4 py-2 bg-blue-600/10 hover:bg-blue-600 text-blue-400 hover:text-white rounded-lg text-sm transition-all border border-blue-500/20">
            Изменить
          </button>
          <button @click="deleteItem(item.id)" class="px-4 py-2 bg-red-500/10 hover:bg-red-500 text-red-400 hover:text-white rounded-lg text-sm transition-all border border-red-500/20">
            Удалить
          </button>
        </div>
      </div>
    </div>

    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm">
      <div class="bg-slate-900 border border-white/10 rounded-2xl p-8 max-w-2xl w-full shadow-2xl">
        <h3 class="mb-6 text-2xl font-bold text-white">{{ isEditing ? 'Редактирование пары' : 'Новая пара' }}</h3>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="mb-2 block text-sm text-slate-400">Предмет</label>
            <input v-model="form.subject" type="text" class="w-full bg-slate-800 rounded-lg px-4 py-3 text-white" />
          </div>

          <div>
            <label class="mb-2 block text-sm text-slate-400">Группа</label>
            <input v-model="form.group_name" type="text" placeholder="P-21" class="w-full bg-slate-800 rounded-lg px-4 py-3 text-white" />
          </div>

          <div>
            <label class="mb-2 block text-sm text-slate-400">Преподаватель</label>
            <select v-model="form.teacher_id" class="w-full bg-slate-800 rounded-lg px-4 py-3 text-white">
              <option value="" disabled>Выберите преподавателя</option>
              <option v-for="teacher in teachers" :key="teacher.uid" :value="teacher.uid">
                {{ teacher.display_name || teacher.displayName || teacher.email || teacher.uid }}
              </option>
            </select>
          </div>

          <div>
            <label class="mb-2 block text-sm text-slate-400">День недели</label>
            <select v-model.number="form.day_of_week" class="w-full bg-slate-800 rounded-lg px-4 py-3 text-white">
              <option :value="1">Понедельник</option>
              <option :value="2">Вторник</option>
              <option :value="3">Среда</option>
              <option :value="4">Четверг</option>
              <option :value="5">Пятница</option>
              <option :value="6">Суббота</option>
              <option :value="7">Воскресенье</option>
            </select>
          </div>

          <div>
            <label class="mb-2 block text-sm text-slate-400">Номер пары</label>
            <input v-model.number="form.pair_number" type="number" min="1" max="10" class="w-full bg-slate-800 rounded-lg px-4 py-3 text-white" />
          </div>

          <div>
            <label class="mb-2 block text-sm text-slate-400">Аудитория</label>
            <input v-model="form.room" type="text" class="w-full bg-slate-800 rounded-lg px-4 py-3 text-white" />
          </div>

          <div>
            <label class="mb-2 block text-sm text-slate-400">Начало</label>
            <input v-model="form.starts_at" type="time" class="w-full bg-slate-800 rounded-lg px-4 py-3 text-white" />
          </div>

          <div>
            <label class="mb-2 block text-sm text-slate-400">Конец</label>
            <input v-model="form.ends_at" type="time" class="w-full bg-slate-800 rounded-lg px-4 py-3 text-white" />
          </div>
        </div>

        <div class="flex gap-3 mt-8">
          <button @click="showModal = false" class="flex-grow rounded-xl bg-white/5 py-3 text-white transition-all hover:bg-white/10">Отмена</button>
          <button @click="saveItem" class="flex-grow rounded-xl bg-blue-600 py-3 font-semibold text-white transition-all hover:bg-blue-500">
            {{ isEditing ? 'Сохранить' : 'Создать' }}
          </button>
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

const loading = ref(true)
const errorMessage = ref('')
const showModal = ref(false)
const isEditing = ref(false)
const currentId = ref('')
const schedule = ref([])
const teachers = ref([])

const form = reactive({
  subject: '',
  group_name: '',
  teacher_id: '',
  day_of_week: 1,
  pair_number: 1,
  starts_at: '09:00',
  ends_at: '10:20',
  room: ''
})

const dayName = (day) => {
  const names = {
    1: 'Понедельник',
    2: 'Вторник',
    3: 'Среда',
    4: 'Четверг',
    5: 'Пятница',
    6: 'Суббота',
    7: 'Воскресенье'
  }
  return names[Number(day)] || 'Неизвестно'
}

const fetchSchedule = async () => {
  loading.value = true
  errorMessage.value = ''
  try {
    const data = await api('/admin/schedule')
    schedule.value = data || []
  } catch (e) {
    console.error('Не удалось получить расписание:', e)
    const status = e?.status || e?.response?.status
    const message = e?.data?.error || e?.message || 'Не удалось загрузить расписание'
    errorMessage.value = `Не удалось загрузить расписание (${status || 'без статуса'}): ${message}`
    schedule.value = []
  } finally {
    loading.value = false
  }
}

const fetchTeachers = async () => {
  try {
    const data = await api('/admin/users?limit=200')
    const users = data.users || []
    teachers.value = users.filter((u) => ['teacher', 'admin'].includes(u.role))
  } catch (e) {
    const status = e?.status || e?.response?.status
    const message = e?.data?.error || e?.message || 'Не удалось загрузить преподавателей'
    errorMessage.value = `Не удалось загрузить преподавателей (${status || 'без статуса'}): ${message}`
    teachers.value = []
  }
}

const resetForm = () => {
  Object.assign(form, {
    subject: '',
    group_name: '',
    teacher_id: teachers.value[0]?.uid || '',
    day_of_week: 1,
    pair_number: 1,
    starts_at: '09:00',
    ends_at: '10:20',
    room: ''
  })
}

const openModal = (item = null) => {
  if (item) {
    isEditing.value = true
    currentId.value = item.id
    Object.assign(form, {
      subject: item.subject || '',
      group_name: item.group_name || '',
      teacher_id: item.teacher_id || '',
      day_of_week: Number(item.day_of_week || 1),
      pair_number: Number(item.pair_number || 1),
      starts_at: item.starts_at || '09:00',
      ends_at: item.ends_at || '',
      room: item.room || ''
    })
  } else {
    isEditing.value = false
    currentId.value = ''
    resetForm()
  }

  showModal.value = true
}

const saveItem = async () => {
  try {
    const url = isEditing.value ? `/admin/schedule/${currentId.value}` : '/admin/schedule'
    const method = isEditing.value ? 'PUT' : 'POST'

    await api(url, {
      method,
      body: {
        subject: form.subject,
        group_name: form.group_name,
        teacher_id: form.teacher_id,
        day_of_week: form.day_of_week,
        pair_number: form.pair_number,
        starts_at: form.starts_at,
        ends_at: form.ends_at,
        room: form.room
      }
    })

    showModal.value = false
    await fetchSchedule()
  } catch (e) {
    const message = e?.data?.error || e.message || 'Не удалось сохранить пару'
    alert('Ошибка: ' + message)
  }
}

const deleteItem = async (id) => {
  if (!confirm('Удалить эту пару из расписания?')) return
  try {
    await api(`/admin/schedule/${id}`, { method: 'DELETE' })
    await fetchSchedule()
  } catch (e) {
    alert('Ошибка удаления: ' + (e?.data?.error || e.message))
  }
}

onMounted(async () => {
  await fetchTeachers()
  resetForm()
  await fetchSchedule()
})
</script>



