<template>
  <div class="max-w-6xl mx-auto py-8 px-4 grid grid-cols-1 md:grid-cols-12 gap-8">
    <aside class="md:col-span-3 space-y-4">
      <div class="bg-slate-900 border border-white/5 rounded-xl p-4 sticky top-8">
        <nav class="space-y-2">
          <NuxtLink to="/" class="flex items-center space-x-3 p-3 bg-blue-600/10 text-blue-400 rounded-lg">
            <span>🏠</span>
            <span class="font-medium">Лента</span>
          </NuxtLink>
          <NuxtLink to="/clubs" class="flex items-center space-x-3 p-3 text-slate-400 hover:bg-white/5 rounded-lg transition-all">
            <span>🎭</span>
            <span>Клубы</span>
          </NuxtLink>
          <NuxtLink to="/chat" class="flex items-center space-x-3 p-3 text-slate-400 hover:bg-white/5 rounded-lg transition-all">
            <span>💬</span>
            <span>Сообщения</span>
          </NuxtLink>
          <NuxtLink to="/profile" class="flex items-center space-x-3 p-3 text-slate-400 hover:bg-white/5 rounded-lg transition-all">
            <span>👤</span>
            <span>Профиль</span>
          </NuxtLink>
        </nav>

        <div class="mt-8 pt-8 border-t border-white/5">
          <button @click="handleLogout" class="w-full flex items-center space-x-3 p-3 text-red-400 hover:bg-red-500/10 rounded-lg transition-all">
            <span>🚪</span>
            <span>Выйти</span>
          </button>
        </div>
      </div>
    </aside>

    <main class="md:col-span-6 space-y-6">
      <div class="bg-slate-900 border border-white/5 rounded-xl p-6 shadow-xl">
        <div class="flex items-start space-x-4">
          <div class="w-10 h-10 rounded-full bg-gradient-to-tr from-blue-500 to-indigo-500 flex-shrink-0"></div>
          <div class="flex-grow">
            <textarea
              v-model="newPost"
              class="w-full bg-slate-800 border-none rounded-lg p-3 text-white placeholder-slate-500 focus:ring-1 focus:ring-blue-500/50 min-h-[100px] resize-none"
              placeholder="Что нового в колледже?"
            ></textarea>
            <div class="mt-3 flex justify-end">
              <button
                class="bg-blue-600 hover:bg-blue-500 text-white px-6 py-2 rounded-lg font-medium transition-all shadow-lg shadow-blue-600/20"
                @click="createPost"
              >
                Опубликовать
              </button>
            </div>
          </div>
        </div>
      </div>

      <div class="space-y-6">
        <div v-if="loading" class="space-y-6">
          <div v-for="i in 3" :key="i" class="bg-slate-900 border border-white/5 rounded-xl p-6 shadow-xl animate-pulse">
            <div class="h-20 bg-slate-800 rounded"></div>
          </div>
        </div>

        <div v-for="post in posts" :key="post.id" class="bg-slate-900 border border-white/5 rounded-xl p-6 shadow-xl hover:border-blue-500/20 transition-all">
          <div class="flex items-center space-x-4 mb-4">
            <div class="w-10 h-10 rounded-full bg-slate-800 border border-white/10 flex items-center justify-center text-xs text-slate-500">
              {{ avatarInitial(post.author_name, post.author_id) }}
            </div>
            <div>
              <div class="text-white font-semibold text-sm">{{ post.author_name || `Пользователь #${shortId(post.author_id)}` }}</div>
              <div class="text-slate-500 text-[10px] uppercase">{{ roleLabel(post.author_role) }} • {{ formatDateTime(post.created_at) }}</div>
            </div>
          </div>

          <p class="text-slate-300 leading-relaxed">{{ post.content }}</p>

          <div class="mt-4 pt-4 border-t border-white/5 flex items-center space-x-6">
            <button
              @click="likePost(post)"
              class="flex items-center space-x-2 transition-colors"
              :class="post.liked ? 'text-blue-400' : 'text-slate-500 hover:text-blue-400'"
            >
              <span class="text-lg">💙</span>
              <span class="text-xs">{{ post.likes || 0 }}</span>
            </button>
            <button
              @click="toggleComments(post)"
              class="flex items-center space-x-2 text-slate-500 hover:text-blue-400 transition-colors"
            >
              <span class="text-lg">💬</span>
              <span class="text-xs">{{ post.comments_count || 0 }}</span>
            </button>
          </div>

          <div v-if="post.showComments" class="mt-4 pt-4 border-t border-white/5 space-y-4">
            <div v-for="comment in post.comments" :key="comment.id" class="flex items-start space-x-3">
              <div class="w-8 h-8 rounded-full bg-slate-800 flex-shrink-0 flex items-center justify-center text-[10px] text-slate-500 border border-white/5">
                {{ avatarInitial(comment.author_name, comment.author_id) }}
              </div>
              <div class="bg-slate-800/50 rounded-lg p-3 flex-grow">
                <div class="text-xs font-semibold text-white mb-1">
                  {{ comment.author_name || `Пользователь #${shortId(comment.author_id)}` }}
                  <span class="text-slate-500 font-normal">• {{ roleLabel(comment.author_role) }}</span>
                </div>
                <p class="text-slate-300 text-xs">{{ comment.text }}</p>
              </div>
            </div>

            <div class="flex items-center space-x-3 mt-2">
              <input
                v-model="post.newComment"
                type="text"
                placeholder="Напишите комментарий..."
                class="flex-grow bg-slate-800 border-none rounded-lg px-3 py-2 text-xs text-white focus:ring-1 focus:ring-blue-500/50"
                @keyup.enter="addComment(post)"
              />
              <button
                @click="addComment(post)"
                class="text-blue-500 text-sm font-medium hover:text-blue-400"
              >
                Отправить
              </button>
            </div>
          </div>
        </div>
      </div>
    </main>

    <aside class="md:col-span-3 space-y-6">
      <div class="bg-slate-900 border border-white/5 rounded-xl p-6 shadow-xl">
        <h3 class="text-white font-semibold mb-4">Новости колледжа</h3>
        <div class="space-y-4">
          <div v-if="news.length === 0" class="text-slate-400 text-sm">Пока нет новостей.</div>

          <div v-else v-for="item in news" :key="item.id" class="border-b border-white/5 pb-4 last:border-0 last:pb-0">
            <h4 class="text-blue-400 text-sm font-medium mb-1">{{ item.title }}</h4>
            <p class="text-slate-400 text-xs">{{ item.description }}</p>
            <p class="text-xs text-slate-600 mt-1">{{ formatDateTime(item.created_at) }}</p>
          </div>
        </div>
      </div>
    </aside>
  </div>
</template>

<script setup>
definePageMeta({
  middleware: 'auth'
})

const { logout } = useAuth()
const { fetchApi: api } = useApi()
const router = useRouter()

const newPost = ref('')
const posts = ref([])
const loading = ref(true)
const news = ref([])

const roleLabel = (role) => {
  if (role === 'admin') return 'Администратор'
  if (role === 'teacher') return 'Преподаватель'
  return 'Студент'
}

const shortId = (uid) => {
  if (!uid || typeof uid !== 'string') return '----'
  return uid.slice(0, 4)
}

const avatarInitial = (name, uid) => {
  const base = (typeof name === 'string' && name.trim()) ? name.trim() : shortId(uid)
  return (base[0] || 'U').toUpperCase()
}

const formatDateTime = (value) => {
  if (!value) return 'время не указано'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return 'время не указано'
  return date.toLocaleString('ru-RU')
}

const fetchNews = async () => {
  try {
    const data = await api('/news')
    news.value = data || []
  } catch (e) {
    console.error('Failed to fetch news:', e)
  }
}

const handleLogout = async () => {
  await logout()
  router.push('/login')
}

const fetchPosts = async () => {
  loading.value = true
  try {
    const data = await api('/posts')
    posts.value = data || []
  } catch (e) {
    console.error('Failed to fetch posts:', e)
  } finally {
    loading.value = false
  }
}

const createPost = async () => {
  if (!newPost.value.trim()) return

  try {
    await api('/posts', {
      method: 'POST',
      body: { content: newPost.value }
    })
    newPost.value = ''
    await fetchPosts()
  } catch (e) {
    alert('Ошибка при публикации: ' + (e?.message || 'unknown error'))
  }
}

const likePost = async (post) => {
  if (post.liked) return

  post.likes = (post.likes || 0) + 1
  post.liked = true

  try {
    await api(`/posts/${post.id}/like`, { method: 'POST' })
  } catch (e) {
    post.likes = Math.max((post.likes || 1) - 1, 0)
    post.liked = false
    console.error('Failed to like post:', e)
  }
}

const toggleComments = async (post) => {
  post.showComments = !post.showComments
  if (post.showComments && !post.comments) {
    try {
      const data = await api(`/posts/${post.id}/comments`)
      post.comments = data || []
    } catch (e) {
      console.error('Failed to fetch comments:', e)
    }
  }
}

const addComment = async (post) => {
  if (!post.newComment?.trim()) return

  const text = post.newComment
  post.newComment = ''

  try {
    const data = await api(`/posts/${post.id}/comments`, {
      method: 'POST',
      body: { text }
    })
    if (!post.comments) post.comments = []
    post.comments.push(data)
    post.comments_count = (post.comments_count || 0) + 1
  } catch (e) {
    alert('Ошибка при добавлении комментария')
  }
}

onMounted(() => {
  fetchPosts()
  fetchNews()
})
</script>