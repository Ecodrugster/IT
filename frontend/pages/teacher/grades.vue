<template>
  <div class="max-w-6xl mx-auto py-8 px-4">
    <div class="flex justify-between items-center mb-8">
      <div>
        <h1 class="text-3xl font-bold text-white">Журнал преподавателя</h1>
        <p class="text-slate-400">Оценки и посещаемость по парам расписания</p>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <div class="lg:col-span-1 space-y-4">
        <div class="bg-slate-900 border border-white/5 rounded-2xl p-6 shadow-xl">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-xl font-bold text-white">Студенты</h2>
            <select
              v-model="selectedGroup"
              class="bg-slate-800 border border-white/10 rounded-lg px-3 py-2 text-xs text-white"
            >
              <option value="">Все группы</option>
              <option v-for="g in taughtGroups" :key="g" :value="g">{{ g }}</option>
            </select>
          </div>

          <div class="space-y-2 max-h-[600px] overflow-y-auto pr-2">
            <button
              v-for="student in filteredStudents"
              :key="student.uid"
              @click="selectedStudent = student"
              class="w-full flex items-center space-x-3 p-3 rounded-xl transition-all border border-transparent"
              :class="selectedStudent?.uid === student.uid ? 'bg-blue-600/20 border-blue-600/50' : 'hover:bg-white/5'"
            >
              <div class="w-10 h-10 rounded-full bg-slate-800 flex items-center justify-center text-xs font-bold text-white">
                {{ avatarInitial(student.display_name || student.email) }}
              </div>
              <div class="text-left overflow-hidden">
                <div class="text-sm font-medium text-white truncate">{{ student.display_name || student.email }}</div>
                <div class="text-[10px] text-slate-500 uppercase">
                  {{ student.group_name || 'Без группы' }}
                </div>
              </div>
            </button>
          </div>

          <div v-if="filteredStudents.length === 0" class="text-slate-500 text-sm mt-4">
            В выбранной группе пока нет студентов.
          </div>
        </div>
      </div>

      <div class="lg:col-span-2">
        <div v-if="selectedStudent" class="bg-slate-900 border border-white/5 rounded-2xl p-8 shadow-xl">
          <div class="flex items-center space-x-4 mb-8 pb-8 border-b border-white/5">
            <div class="w-16 h-16 rounded-full bg-blue-600 flex items-center justify-center text-xl font-bold text-white shadow-lg shadow-blue-600/30">
              {{ avatarInitial(selectedStudent.display_name || selectedStudent.email) }}
            </div>
            <div>
              <h2 class="text-2xl font-bold text-white">{{ selectedStudent.display_name || selectedStudent.email }}</h2>
              <p class="text-slate-400">Группа: {{ selectedStudent.group_name || 'не указана' }}</p>
            </div>
          </div>

          <form @submit.prevent="submitGrade" class="space-y-6">
            <div>
              <label class="block text-sm font-medium text-slate-400 mb-2">Пара из расписания</label>
              <select
                v-model="form.schedule_id"
                class="w-full bg-slate-800 border-none rounded-lg px-4 py-3 text-white focus:ring-2 focus:ring-blue-500/50"
                required
              >
                <option value="" disabled>Выберите пару</option>
                <option v-for="item in availablePairs" :key="item.id" :value="item.id">
                  {{ formatPairLabel(item) }}
                </option>
              </select>
              <p v-if="availablePairs.length === 0" class="text-xs text-amber-400 mt-2">
                В вашем расписании нет пар для этой группы.
              </p>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label class="block text-sm font-medium text-slate-400 mb-2">Дата занятия</label>
                <input
                  v-model="form.lesson_date"
                  type="date"
                  class="w-full bg-slate-800 border-none rounded-lg px-4 py-3 text-white focus:ring-2 focus:ring-blue-500/50"
                  required
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-slate-400 mb-2">Оценка (1-12)</label>
                <input
                  v-model.number="form.value"
                  type="number"
                  min="1"
                  max="12"
                  class="w-full bg-slate-800 border-none rounded-lg px-4 py-3 text-white focus:ring-2 focus:ring-blue-500/50"
                  required
                />
              </div>
            </div>

            <div>
              <label class="block text-sm font-medium text-slate-400 mb-2">Комментарий к оценке</label>
              <textarea
                v-model="form.comment"
                rows="3"
                placeholder="Почему была выставлена эта оценка..."
                class="w-full bg-slate-800 border-none rounded-lg px-4 py-3 text-white focus:ring-2 focus:ring-blue-500/50 resize-none"
              ></textarea>
            </div>

            <div class="flex justify-end">
              <button
                type="submit"
                :disabled="submittingGrade || !form.schedule_id"
                class="bg-blue-600 hover:bg-blue-500 text-white px-8 py-3 rounded-xl font-bold shadow-lg shadow-blue-600/20 transition-all disabled:opacity-50"
              >
                {{ submittingGrade ? 'Сохранение...' : 'Добавить оценку' }}
              </button>
            </div>
          </form>

          <div class="mt-8 p-5 rounded-xl border border-white/10 bg-white/5">
            <h3 class="text-white font-semibold mb-4">Посещаемость</h3>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
              <select v-model="attendance.status" class="bg-slate-800 border border-white/10 rounded-lg px-4 py-3 text-white text-sm">
                <option value="present">Присутствует</option>
                <option value="late">Опоздал</option>
                <option value="absent">Отсутствует</option>
                <option value="excused">Уважительная причина</option>
              </select>
              <input
                v-model="attendance.comment"
                type="text"
                placeholder="Комментарий"
                class="bg-slate-800 border border-white/10 rounded-lg px-4 py-3 text-white text-sm md:col-span-2"
              />
            </div>
            <div class="mt-3 flex justify-end">
              <button
                @click="markAttendance"
                :disabled="submittingAttendance || !form.schedule_id"
                class="px-6 py-2.5 rounded-lg text-sm font-semibold bg-emerald-600 hover:bg-emerald-500 text-white disabled:opacity-50"
              >
                {{ submittingAttendance ? 'Сохранение...' : 'Отметить посещаемость' }}
              </button>
            </div>
          </div>

          <div class="mt-10 pt-8 border-t border-white/5">
            <h3 class="text-lg font-bold text-white mb-4">Последние оценки</h3>
            <div class="space-y-3 mb-8">
              <div v-for="g in studentGrades" :key="g.id" class="bg-white/5 p-4 rounded-xl flex justify-between items-center">
                <div>
                  <div class="text-sm font-bold text-white">{{ g.subject }}</div>
                  <div class="text-xs text-slate-400 mt-1">{{ g.lesson_date || formatDate(g.created_at) }} • Пара {{ g.pair_number }}</div>
                  <div class="text-xs text-slate-500">{{ g.comment || 'Без комментария' }}</div>
                </div>
                <div class="text-2xl font-bold text-blue-500">{{ g.value }}</div>
              </div>
              <div v-if="studentGrades.length === 0" class="text-center text-slate-500 py-4 italic text-sm">Оценок пока нет</div>
            </div>

            <h3 class="text-lg font-bold text-white mb-4">Последняя посещаемость</h3>
            <div class="space-y-3">
              <div v-for="a in attendanceHistory" :key="a.id" class="bg-white/5 p-4 rounded-xl flex justify-between items-center">
                <div>
                  <div class="text-sm font-semibold text-white">{{ a.subject }}</div>
                  <div class="text-xs text-slate-400 mt-1">{{ a.lesson_date }} • Пара {{ a.pair_number }}</div>
                  <div class="text-xs text-slate-500">{{ a.comment || 'Без комментария' }}</div>
                </div>
                <div class="text-xs font-bold uppercase" :class="attendanceStatusClass(a.status)">
                  {{ attendanceStatusLabel(a.status) }}
                </div>
              </div>
              <div v-if="attendanceHistory.length === 0" class="text-center text-slate-500 py-4 italic text-sm">Записей о посещаемости пока нет</div>
            </div>
          </div>
        </div>

        <div v-else class="h-full bg-slate-900/50 border border-dashed border-white/10 rounded-2xl flex flex-col items-center justify-center text-slate-500 p-12">
          <div class="text-6xl mb-4">📘</div>
          <p>Выберите студента слева для работы с оценками и посещаемостью</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  middleware: ['auth', 'teacher']
})

const { fetchApi: api } = useApi()

const allStudents = ref([])
const selectedStudent = ref(null)
const studentGrades = ref([])
const schedulePairs = ref([])
const attendanceHistory = ref([])

const selectedGroup = ref('')
const submittingGrade = ref(false)
const submittingAttendance = ref(false)

const form = reactive({
  schedule_id: '',
  lesson_date: new Date().toISOString().slice(0, 10),
  value: 12,
  comment: ''
})

const attendance = reactive({
  status: 'present',
  comment: ''
})

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
  if (!d || Number.isNaN(d.getTime())) return ''
  return d.toLocaleDateString('ru-RU')
}

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

const avatarInitial = (value) => {
  const text = typeof value === 'string' ? value.trim() : ''
  return (text[0] || 'U').toUpperCase()
}

const attendanceStatusLabel = (status) => {
  if (status === 'present') return 'Присутствует'
  if (status === 'late') return 'Опоздал'
  if (status === 'excused') return 'Уважительная причина'
  return 'Отсутствует'
}

const attendanceStatusClass = (status) => {
  if (status === 'present') return 'text-emerald-400'
  if (status === 'late') return 'text-amber-400'
  if (status === 'excused') return 'text-blue-400'
  return 'text-red-400'
}

const taughtGroups = computed(() => {
  const set = new Set()
  for (const pair of schedulePairs.value) {
    const group = (pair.group_name || '').trim()
    if (group) set.add(group)
  }
  return Array.from(set).sort()
})

const filteredStudents = computed(() => {
  const taught = taughtGroups.value
  return allStudents.value
    .filter((u) => (u.role || 'student') === 'student')
    .filter((u) => {
      const group = (u.group_name || u.group || '').trim()
      if (!group) return false
      if (selectedGroup.value) return group.toLowerCase() === selectedGroup.value.toLowerCase()
      if (taught.length === 0) return true
      return taught.some(g => g.toLowerCase() === group.toLowerCase())
    })
})

const availablePairs = computed(() => {
  if (!selectedStudent.value) return schedulePairs.value
  const studentGroup = (selectedStudent.value.group_name || selectedStudent.value.group || '').trim().toLowerCase()
  if (!studentGroup) return []
  return schedulePairs.value.filter((item) => (item.group_name || '').trim().toLowerCase() === studentGroup)
})

const formatPairLabel = (item) => {
  return `${dayName(item.day_of_week)}, Пара ${item.pair_number} • ${item.subject} • ${item.group_name}`
}

const fetchStudents = async () => {
  const data = await api('/users?role=student&limit=2000')
  allStudents.value = (data || []).map((u) => ({
    ...u,
    group_name: u.group_name || u.group || ''
  }))
}

const fetchSchedule = async () => {
  const data = await api('/teacher/schedule')
  schedulePairs.value = data || []
}

const fetchStudentGrades = async (uid) => {
  if (!uid) {
    studentGrades.value = []
    return
  }
  const data = await api(`/grades?student_id=${uid}`)
  studentGrades.value = data || []
}

const fetchAttendance = async () => {
  if (!selectedStudent.value?.uid || !form.schedule_id) {
    attendanceHistory.value = []
    return
  }

  const params = new URLSearchParams({
    student_id: selectedStudent.value.uid,
    schedule_id: form.schedule_id
  })

  const data = await api(`/teacher/attendance?${params.toString()}`)
  attendanceHistory.value = data || []
}

watch(filteredStudents, (list) => {
  if (!selectedStudent.value && list.length > 0) {
    selectedStudent.value = list[0]
  }
})

watch(availablePairs, (pairs) => {
  if (pairs.length === 0) {
    form.schedule_id = ''
    return
  }

  const exists = pairs.some(p => p.id === form.schedule_id)
  if (!exists) {
    form.schedule_id = pairs[0].id
  }
})

watch([selectedStudent, () => form.schedule_id], async ([student]) => {
  if (student?.uid) {
    await fetchStudentGrades(student.uid)
  }
  await fetchAttendance()
})

const submitGrade = async () => {
  if (!selectedStudent.value || !form.schedule_id) return

  submittingGrade.value = true
  try {
    await api('/teacher/grades', {
      method: 'POST',
      body: {
        student_id: selectedStudent.value.uid,
        schedule_id: form.schedule_id,
        lesson_date: form.lesson_date,
        value: form.value,
        comment: form.comment
      }
    })

    form.comment = ''
    await fetchStudentGrades(selectedStudent.value.uid)
    alert('Оценка успешно добавлена')
  } catch (e) {
    alert('Error: ' + (e?.data?.error || e.message))
  } finally {
    submittingGrade.value = false
  }
}

const markAttendance = async () => {
  if (!selectedStudent.value || !form.schedule_id) return

  submittingAttendance.value = true
  try {
    await api('/teacher/attendance', {
      method: 'POST',
      body: {
        student_id: selectedStudent.value.uid,
        schedule_id: form.schedule_id,
        lesson_date: form.lesson_date,
        status: attendance.status,
        comment: attendance.comment
      }
    })

    attendance.comment = ''
    await fetchAttendance()
    alert('Посещаемость сохранена')
  } catch (e) {
    alert('Attendance error: ' + (e?.data?.error || e.message))
  } finally {
    submittingAttendance.value = false
  }
}

onMounted(async () => {
  await Promise.all([fetchStudents(), fetchSchedule()])
  if (!selectedGroup.value && taughtGroups.value.length > 0) {
    selectedGroup.value = taughtGroups.value[0]
  }
})
</script>
