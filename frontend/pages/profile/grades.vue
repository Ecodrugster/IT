<template>
  <div class="max-w-4xl mx-auto py-8 px-4">
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-white">Мои оценки</h1>
      <p class="text-slate-400">Успеваемость по конкретным парам</p>
    </div>

    <div v-if="loading" class="space-y-4">
      <div v-for="i in 5" :key="i" class="h-24 bg-slate-900 animate-pulse rounded-2xl border border-white/5"></div>
    </div>

    <div v-else class="space-y-4">
      <div v-for="grade in grades" :key="grade.id" class="bg-slate-900 border border-white/5 p-6 rounded-2xl shadow-xl flex justify-between items-center group hover:border-blue-500/30 transition-all">
        <div class="flex items-start space-x-6">
          <div class="w-14 h-14 rounded-xl bg-blue-600/10 flex items-center justify-center text-2xl font-bold text-blue-500">
            {{ grade.value }}
          </div>
          <div>
            <h3 class="text-lg font-bold text-white">{{ grade.subject }}</h3>
            <div class="flex flex-wrap items-center gap-2 mt-1 text-xs text-slate-500">
              <span>{{ formatDate(grade.lesson_date || grade.created_at) }}</span>
              <span>•</span>
              <span>{{ dayName(grade.day_of_week) }}</span>
              <span>•</span>
              <span>{{ grade.pair_number }} пара</span>
              <span v-if="grade.group_name">• {{ grade.group_name }}</span>
              <span v-if="grade.room">• ауд. {{ grade.room }}</span>
            </div>
            <div class="text-xs text-slate-400 italic mt-1">{{ grade.comment || 'Без комментария' }}</div>
          </div>
        </div>

        <div class="text-right hidden md:block">
          <div class="text-[10px] text-slate-500 uppercase font-bold tracking-widest mb-1">Выставил(а)</div>
          <div class="text-xs text-slate-300 font-medium">{{ grade.teacher_name || 'Преподаватель' }}</div>
        </div>
      </div>

      <div v-if="grades.length === 0" class="text-center py-20 bg-slate-900 rounded-2xl border border-dashed border-white/10">
        <div class="text-5xl mb-4">📊</div>
        <p class="text-slate-500">У вас пока нет выставленных оценок</p>
      </div>
    </div>
  </div>
</template>

<script setup>
const { fetchApi: api } = useApi()
const grades = ref([])
const loading = ref(true)

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
  return names[Number(day)] || 'День не указан'
}

const parseDate = (value) => {
  if (!value) return null
  if (value instanceof Date) return value
  if (typeof value === 'string' || typeof value === 'number') return new Date(value)
  if (typeof value === 'object' && value.seconds) return new Date(value.seconds * 1000)
  if (typeof value === 'object' && value._seconds) return new Date(value._seconds * 1000)
  return null
}

const formatDate = (date) => {
  const d = parseDate(date)
  if (!d || Number.isNaN(d.getTime())) return 'Дата не указана'
  return d.toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric'
  })
}

const fetchGrades = async () => {
  try {
    const data = await api('/grades')
    grades.value = data || []
  } catch (e) {
    console.error('Failed to fetch grades:', e)
  } finally {
    loading.value = false
  }
}

onMounted(fetchGrades)
</script>

