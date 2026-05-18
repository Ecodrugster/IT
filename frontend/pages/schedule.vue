<template>
  <div class="max-w-6xl mx-auto py-8 px-4 space-y-8">
    <div class="flex flex-col md:flex-row md:items-end md:justify-between gap-4">
      <div>
        <h1 class="text-3xl font-bold text-white">Расписание</h1>
        <p class="text-slate-400 text-sm">{{ weekLabel }}</p>
      </div>

      <form @submit.prevent="loadSchedule" class="flex gap-2">
        <input
          v-model="groupFilter"
          type="text"
          :disabled="groupLocked"
          placeholder="Группа П-21"
          class="bg-slate-900 border border-white/10 rounded-lg px-4 py-2 text-white text-sm disabled:opacity-60"
        />
        <button
          type="submit"
          class="bg-blue-600 hover:bg-blue-500 text-white px-4 py-2 rounded-lg text-sm font-semibold transition-all"
        >
          Принять
        </button>
      </form>
    </div>

    <div v-if="errorMessage" class="bg-red-500/10 border border-red-500/30 text-red-300 rounded-xl px-4 py-3 text-sm">
      {{ errorMessage }}
    </div>

    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div v-for="i in 4" :key="i" class="h-40 bg-slate-900 border border-white/5 rounded-2xl animate-pulse"></div>
    </div>

    <div v-else class="space-y-6">
      <div v-if="schedule.length === 0" class="bg-slate-900 border border-dashed border-white/10 rounded-2xl p-12 text-center text-slate-500">
        Расписание пусто.
      </div>

      <div v-for="day in orderedDays" :key="day.value" class="bg-slate-900 border border-white/5 rounded-2xl p-6">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h2 class="text-xl font-bold text-white">{{ day.label }}</h2>
            <p class="text-xs text-slate-500 mt-1">{{ formatDayDate(day.value) }}</p>
          </div>
          <span class="text-xs text-slate-500 uppercase tracking-widest">{{ dayItems(day.value).length }} пар</span>
        </div>

        <div v-if="dayItems(day.value).length === 0" class="text-sm text-slate-500 italic">Нет пар</div>

        <div v-else class="space-y-3">
          <div
            v-for="item in dayItems(day.value)"
            :key="item.id"
            class="bg-white/5 border border-white/5 rounded-xl p-4"
          >
            <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-3">
              <div>
                <div class="text-white font-semibold">{{ item.subject }}</div>
                <div class="text-sm text-slate-400 mt-1">
                  Пара {{ item.pair_number }} • {{ item.starts_at }}{{ item.ends_at ? `-${item.ends_at}` : '' }}
                </div>
              </div>

              <div class="text-sm text-slate-300">
                <div><span class="text-slate-500">Группа:</span> {{ item.group_name }}</div>
                <div><span class="text-slate-500">Преподаватель:</span> {{ item.teacher_name || item.teacher_id }}</div>
                <div v-if="item.room"><span class="text-slate-500">Аудитория:</span> {{ item.room }}</div>
              </div>
            </div>
          </div>
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

const loading = ref(true)
const schedule = ref([])
const groupFilter = ref('')
const errorMessage = ref('')

const orderedDays = [
  { value: 1, label: 'Понедельник' },
  { value: 2, label: 'Вторник' },
  { value: 3, label: 'Среда' },
  { value: 4, label: 'Четверг' },
  { value: 5, label: 'Пятница' },
  { value: 6, label: 'Суббота' },
  { value: 7, label: 'Воскресенье' }
]

const groupLocked = computed(() => {
  return userStore.role === 'student' && !!(userStore.profile?.group_name || userStore.profile?.group)
})

const weekStart = computed(() => {
  const now = new Date()
  const day = now.getDay() === 0 ? 7 : now.getDay()
  const monday = new Date(now)
  monday.setDate(now.getDate() - (day - 1))
  monday.setHours(0, 0, 0, 0)
  return monday
})

const weekLabel = computed(() => {
  const start = weekStart.value
  const end = new Date(start)
  end.setDate(start.getDate() + 6)
  return `${start.toLocaleDateString('ru-RU')} - ${end.toLocaleDateString('ru-RU')}`
})

const formatDayDate = (dayOfWeek) => {
  const date = new Date(weekStart.value)
  date.setDate(weekStart.value.getDate() + (Number(dayOfWeek) - 1))
  return date.toLocaleDateString('ru-RU', { weekday: 'long', day: 'numeric', month: 'long' })
}

const dayItems = (day) => {
  return schedule.value.filter(item => Number(item.day_of_week) === day)
}

const loadSchedule = async () => {
  loading.value = true
  errorMessage.value = ''
  try {
    const query = groupFilter.value.trim() ? `?group=${encodeURIComponent(groupFilter.value.trim())}` : ''
    const data = await api(`/schedule${query}`)
    schedule.value = data || []
  } catch (e) {
    console.error('Failed to fetch schedule:', e)
    const status = e?.status || e?.response?.status
    const message = e?.data?.error || e?.message || 'Failed to load schedule'
    errorMessage.value = `Could not load schedule (${status || 'no-status'}): ${message}`
    schedule.value = []
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  const profileGroup = userStore.profile?.group || userStore.profile?.group_name
  if (profileGroup) {
    groupFilter.value = profileGroup
  }
  await loadSchedule()
})
</script>
