<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h2 class="text-2xl font-bold text-white">Управление новостями</h2>
      <button
        @click="openCreateModal"
        class="rounded-xl bg-blue-600 px-6 py-2.5 font-semibold text-white shadow-lg shadow-blue-600/20 transition-all hover:bg-blue-500"
      >
        + Добавить новость
      </button>
    </div>

    <div v-if="errorMessage" class="rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-300">
      {{ errorMessage }}
    </div>

    <div class="overflow-hidden rounded-2xl border border-white/5 bg-slate-900 shadow-xl">
      <table class="w-full border-collapse text-left">
        <thead>
          <tr class="bg-white/5 text-[10px] font-bold uppercase tracking-wider text-slate-500">
            <th class="px-6 py-4">Заголовок</th>
            <th class="px-6 py-4">Категория</th>
            <th class="px-6 py-4">Создано</th>
            <th class="px-6 py-4 text-right">Действия</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-white/5">
          <tr v-for="item in news" :key="item.id" class="text-sm transition-colors hover:bg-white/[0.02]">
            <td class="px-6 py-4">
              <div class="font-medium text-white">{{ item.title }}</div>
              <div class="max-w-xs truncate text-xs text-slate-500">{{ item.description }}</div>
            </td>
            <td class="px-6 py-4">
              <span class="rounded-md bg-blue-500/10 px-2 py-1 text-[10px] font-bold uppercase text-blue-500">
                {{ categoryLabel(item.category) }}
              </span>
            </td>
            <td class="px-6 py-4 font-mono text-xs text-slate-400">
              {{ formatDate(item.created_at) }}
            </td>
            <td class="space-x-2 px-6 py-4 text-right">
              <button
                @click="deleteNews(item.id)"
                class="rounded-lg px-3 py-2 text-xs font-semibold text-red-500 transition-all hover:bg-red-500/10"
              >
                Удалить
              </button>
            </td>
          </tr>
        </tbody>
      </table>

      <div v-if="loading" class="p-12 text-center text-slate-500 animate-pulse">Загрузка новостей...</div>
      <div v-else-if="news.length === 0" class="p-12 text-center text-slate-500">Новостей пока нет.</div>
    </div>

    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/80 p-4 backdrop-blur-sm">
      <div class="w-full max-w-md rounded-2xl border border-white/10 bg-slate-900 p-8 shadow-2xl">
        <h2 class="mb-6 text-2xl font-bold text-white">Новая новость</h2>
        <div class="space-y-4">
          <div>
            <label class="mb-2 block text-sm font-medium text-slate-400">Заголовок</label>
            <input v-model="form.title" type="text" class="w-full rounded-lg border-none bg-slate-800 px-4 py-3 text-white focus:ring-2 focus:ring-blue-500/50" />
          </div>
          <div>
            <label class="mb-2 block text-sm font-medium text-slate-400">Категория</label>
            <select v-model="form.category" class="w-full rounded-lg border-none bg-slate-800 px-4 py-3 text-white focus:ring-2 focus:ring-blue-500/50">
              <option value="news">Новость</option>
              <option value="announcement">Объявление</option>
              <option value="event">Событие</option>
              <option value="deadline">Дедлайн</option>
            </select>
          </div>
          <div>
            <label class="mb-2 block text-sm font-medium text-slate-400">Описание</label>
            <textarea v-model="form.description" rows="4" class="w-full resize-none rounded-lg border-none bg-slate-800 px-4 py-3 text-white focus:ring-2 focus:ring-blue-500/50"></textarea>
          </div>
        </div>
        <div class="mt-8 flex space-x-4">
          <button @click="showModal = false" class="flex-grow rounded-xl bg-white/5 py-3 text-white transition-all hover:bg-white/10">Отмена</button>
          <button @click="saveNews" class="flex-grow rounded-xl bg-blue-600 py-3 font-semibold text-white shadow-lg shadow-blue-600/20 transition-all hover:bg-blue-500">Опубликовать</button>
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
const news = ref([])
const loading = ref(false)
const showModal = ref(false)
const errorMessage = ref('')
const form = reactive({ title: '', description: '', category: 'news' })

const formatDate = (value) => {
  if (!value) return 'дата неизвестна'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return 'дата неизвестна'
  return date.toLocaleDateString()
}

const categoryLabel = (category) => {
  const labels = {
    news: 'Новость',
    announcement: 'Объявление',
    event: 'Событие',
    deadline: 'Дедлайн'
  }
  return labels[category] || 'Новость'
}

const fetchNews = async () => {
  loading.value = true
  errorMessage.value = ''
  try {
    const data = await api('/news')
    news.value = data || []
  } catch (e) {
    const status = e?.status || e?.response?.status
    const message = e?.data?.error || e?.message || 'Не удалось загрузить новости'
    errorMessage.value = `Не удалось загрузить новости (${status || 'без статуса'}): ${message}`
    news.value = []
  } finally {
    loading.value = false
  }
}

const openCreateModal = () => {
  form.title = ''
  form.description = ''
  form.category = 'news'
  showModal.value = true
}

const saveNews = async () => {
  try {
    await api('/admin/news', {
      method: 'POST',
      body: { ...form }
    })
    showModal.value = false
    await fetchNews()
  } catch (e) {
    alert('Ошибка сохранения: ' + (e?.data?.error || e?.message || 'неизвестная ошибка'))
  }
}

const deleteNews = async (id) => {
  if (!confirm('Удалить эту новость?')) return
  try {
    await api(`/admin/news/${id}`, { method: 'DELETE' })
    await fetchNews()
  } catch (e) {
    alert('Ошибка удаления: ' + (e?.data?.error || e?.message || 'неизвестная ошибка'))
  }
}

onMounted(fetchNews)
</script>
