<template>
  <div class="space-y-8">
    <section class="space-y-4">
      <div v-if="errorMessage" class="rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-300">
        {{ errorMessage }}
      </div>

      <div class="flex items-center justify-between">
        <h2 class="text-2xl font-bold text-white">Заявки в клубы</h2>
        <div class="text-sm text-slate-500">{{ clubRequests.length }} в очереди</div>
      </div>

      <div v-if="loadingClubs" class="space-y-3">
        <div v-for="i in 2" :key="i" class="h-24 animate-pulse rounded-2xl border border-white/5 bg-slate-900"></div>
      </div>

      <div v-else-if="clubRequests.length === 0" class="rounded-2xl border border-white/5 bg-slate-900 p-10 text-center text-slate-500">
        Нет заявок на модерацию.
      </div>

      <div v-else class="space-y-3">
        <div v-for="club in clubRequests" :key="club.id" class="flex flex-col gap-4 rounded-2xl border border-white/5 bg-slate-900 p-5 lg:flex-row lg:items-center lg:justify-between">
          <div>
            <div class="font-bold text-white">{{ club.name }}</div>
            <div class="mt-1 text-sm text-slate-400">{{ club.description }}</div>
            <div class="mt-2 text-xs text-slate-500">
              Автор: {{ club.created_by_name || club.created_by || 'неизвестно' }} • Участников: {{ club.members?.length || 0 }}
            </div>
          </div>

          <div class="flex gap-2">
            <button
              @click="approveClub(club.id)"
              class="rounded-lg border border-green-500/20 bg-green-500/10 px-4 py-2 text-sm text-green-400 transition-all hover:bg-green-500 hover:text-white"
            >
              Одобрить
            </button>
            <button
              @click="rejectClub(club.id)"
              class="rounded-lg border border-red-500/20 bg-red-500/10 px-4 py-2 text-sm text-red-400 transition-all hover:bg-red-500 hover:text-white"
            >
              Отклонить
            </button>
          </div>
        </div>
      </div>
    </section>

    <section class="space-y-4">
      <div class="flex items-center justify-between">
        <h2 class="text-2xl font-bold text-white">Модерация постов</h2>
        <div class="text-sm text-slate-500">Последние 50 постов</div>
      </div>

      <div v-if="loadingPosts" class="space-y-4">
        <div v-for="i in 3" :key="i" class="h-32 animate-pulse rounded-2xl border border-white/5 bg-slate-900 p-6"></div>
      </div>

      <div v-else-if="posts.length === 0" class="rounded-2xl border border-white/5 bg-slate-900 p-12 text-center text-slate-500">
        Нет постов для модерации.
      </div>

      <div v-else class="space-y-4">
        <div v-for="post in posts" :key="post.id" class="group flex items-start justify-between rounded-2xl border border-white/5 bg-slate-900 p-6 transition-all hover:border-red-500/30">
          <div class="flex gap-4">
            <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-slate-800 text-xs text-slate-500">P</div>
            <div>
              <div class="mb-1 flex items-center gap-2">
                <span class="text-sm font-bold text-white">
                  Автор: {{ post.author_name || post.author_id || 'неизвестно' }}
                </span>
                <span class="font-mono text-[10px] text-slate-500">{{ formatDateTime(post.created_at) }}</span>
              </div>
              <p class="text-sm leading-relaxed text-slate-300">{{ post.content }}</p>
            </div>
          </div>

          <button
            @click="deletePost(post.id)"
            class="rounded-xl border border-red-500/20 bg-red-500/10 px-4 py-2 text-xs font-bold text-red-500 opacity-0 transition-all hover:bg-red-500 hover:text-white group-hover:opacity-100"
          >
            Удалить пост
          </button>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'admin',
  middleware: 'admin'
})

const { fetchApi: api } = useApi()
const posts = ref([])
const clubRequests = ref([])
const loadingPosts = ref(true)
const loadingClubs = ref(true)
const errorMessage = ref('')

const formatDateTime = (value) => {
  if (!value) return 'время неизвестно'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return 'время неизвестно'
  return date.toLocaleString()
}

const fetchPosts = async () => {
  loadingPosts.value = true
  errorMessage.value = ''
  try {
    const data = await api('/admin/posts')
    posts.value = data || []
  } catch (e) {
    console.error('Не удалось получить посты:', e)
    const status = e?.status || e?.response?.status
    const message = e?.data?.error || e?.message || 'Не удалось загрузить посты'
    errorMessage.value = `Не удалось загрузить посты (${status || 'без статуса'}): ${message}`
  } finally {
    loadingPosts.value = false
  }
}

const fetchClubRequests = async () => {
  loadingClubs.value = true
  try {
    const data = await api('/admin/club-requests')
    clubRequests.value = data || []
  } catch (e) {
    console.error('Не удалось получить заявки клубов:', e)
    const status = e?.status || e?.response?.status
    const message = e?.data?.error || e?.message || 'Не удалось загрузить заявки клубов'
    errorMessage.value = `Не удалось загрузить заявки клубов (${status || 'без статуса'}): ${message}`
  } finally {
    loadingClubs.value = false
  }
}

const deletePost = async (id) => {
  if (!confirm('Удалить этот пост безвозвратно?')) return

  try {
    await api(`/admin/posts/${id}`, { method: 'DELETE' })
    posts.value = posts.value.filter((p) => p.id !== id)
    alert('Пост удален модератором')
  } catch (e) {
    alert('Ошибка удаления: ' + (e?.data?.error || e.message))
  }
}

const approveClub = async (id) => {
  try {
    await api(`/admin/club-requests/${id}/approve`, { method: 'POST' })
    clubRequests.value = clubRequests.value.filter((c) => c.id !== id)
  } catch (e) {
    alert('Ошибка одобрения: ' + (e?.data?.error || e.message))
  }
}

const rejectClub = async (id) => {
  const reason = prompt('Причина отклонения (необязательно):') || ''
  try {
    await api(`/admin/club-requests/${id}/reject`, {
      method: 'POST',
      body: { comment: reason }
    })
    clubRequests.value = clubRequests.value.filter((c) => c.id !== id)
  } catch (e) {
    alert('Ошибка отклонения: ' + (e?.data?.error || e.message))
  }
}

onMounted(async () => {
  await Promise.all([fetchPosts(), fetchClubRequests()])
})
</script>
