<template>
  <div class="mx-auto min-h-[calc(100vh-120px)] max-w-4xl px-4 py-8">
    <!-- Шапка страницы -->
    <div class="mb-8 flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
      <div>
        <h1 class="text-2xl font-bold bg-gradient-to-r from-white to-slate-400 bg-clip-text text-transparent flex items-center gap-2">
          <span>📬</span> Центр уведомлений
        </h1>
        <p class="text-xs text-slate-500 mt-1">
          Здесь собраны ваши личные оповещения: сообщения, оценки, посещаемость и новости
        </p>
      </div>

      <button
        v-if="hasUnreadAlerts"
        @click="markAllAsRead"
        :disabled="actionLoading"
        class="flex items-center justify-center gap-2 rounded-full border border-blue-500/20 bg-blue-500/10 px-4 py-2 text-xs font-semibold text-blue-400 transition-all hover:bg-blue-500 hover:text-white hover:shadow-lg hover:shadow-blue-500/20 disabled:opacity-50"
      >
        <span v-if="actionLoading" class="inline-block animate-spin mr-1">⟳</span>
        <span>✓ Отметить все как прочитанные</span>
      </button>
    </div>

    <!-- Фильтры (Табы) -->
    <div class="mb-6 flex flex-wrap gap-2 border-b border-white/5 pb-4">
      <button
        v-for="tab in tabs"
        :key="tab.value"
        @click="activeTab = tab.value"
        class="relative px-4 py-2 text-xs font-medium rounded-full transition-all duration-200"
        :class="activeTab === tab.value 
          ? 'bg-slate-900 border border-white/10 text-white shadow-md' 
          : 'bg-transparent text-slate-400 border border-transparent hover:text-slate-200 hover:bg-white/5'"
      >
        <span class="mr-1">{{ tab.icon }}</span>
        <span>{{ tab.label }}</span>
        <span
          v-if="tab.count > 0"
          class="ml-1.5 inline-flex h-4 min-w-4 items-center justify-center rounded-full bg-red-500 px-1 text-[9px] font-bold text-white"
        >
          {{ tab.count }}
        </span>
      </button>
    </div>

    <!-- Список уведомлений -->
    <div class="space-y-4">
      <!-- Загрузка / Скелетоны -->
      <template v-if="notificationStore.loading">
        <div 
          v-for="i in 4" 
          :key="i"
          class="animate-pulse rounded-xl border border-white/5 bg-slate-900/30 p-5 flex items-start gap-4"
        >
          <div class="h-10 w-10 rounded-full bg-slate-800"></div>
          <div class="flex-grow space-y-2">
            <div class="h-4 w-1/4 rounded bg-slate-800"></div>
            <div class="h-3 w-3/4 rounded bg-slate-800"></div>
            <div class="h-2 w-1/6 rounded bg-slate-800"></div>
          </div>
        </div>
      </template>

      <!-- Основной контент -->
      <template v-else-if="filteredNotifications.length > 0">
        <TransitionGroup 
          name="list" 
          tag="div" 
          class="space-y-4"
        >
          <div
            v-for="item in filteredNotifications"
            :key="item.id"
            @click="handleNotificationClick(item)"
            class="group relative cursor-pointer rounded-xl border p-5 transition-all duration-200 flex flex-col sm:flex-row sm:items-center justify-between gap-4"
            :class="[
              item.read 
                ? 'bg-slate-950/20 hover:bg-slate-900/30 border-white/5 opacity-70 hover:opacity-100' 
                : 'bg-slate-900/60 hover:bg-slate-900/80 border-white/10 shadow-lg shadow-black/30',
              getTypeStyles(item.type).borderLeft
            ]"
          >
            <!-- Левая сторона: Иконка и детали -->
            <div class="flex items-start gap-4 min-w-0 flex-grow">
              <div 
                class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl transition-all group-hover:scale-105"
                :class="getTypeStyles(item.type).badge"
              >
                <span class="text-lg">{{ getTypeStyles(item.type).icon }}</span>
              </div>
              <div class="min-w-0 flex-grow">
                <div class="flex flex-wrap items-center gap-x-2 gap-y-1">
                  <h4 class="text-sm font-semibold text-white truncate max-w-md">
                    {{ item.title }}
                  </h4>
                  <span 
                    v-if="!item.read" 
                    class="rounded bg-blue-500/10 px-1.5 py-0.5 text-[9px] font-bold text-blue-400 uppercase tracking-wider animate-pulse"
                  >
                    Новое
                  </span>
                </div>
                <p class="text-xs text-slate-400 mt-1 break-words leading-relaxed">
                  {{ item.message }}
                </p>
                <span class="text-[10px] text-slate-500 mt-2 block">
                  {{ formatDate(item.created_at) }}
                </span>
              </div>
            </div>

            <!-- Правая сторона: Действия -->
            <div class="flex shrink-0 items-center gap-3 justify-end self-end sm:self-center">
              <button
                v-if="!item.read"
                @click.stop="markAsRead(item.id)"
                class="rounded-full bg-slate-800/80 border border-white/5 hover:border-blue-500/30 hover:bg-blue-950/30 px-3 py-1.5 text-[10px] font-medium text-slate-400 hover:text-blue-400 transition-all"
                title="Пометить как прочитанное"
              >
                Прочитано
              </button>
              <span class="text-slate-600 text-xs hidden sm:inline transition-transform group-hover:translate-x-1">
                ➜
              </span>
            </div>
          </div>
        </TransitionGroup>
      </template>

      <!-- Пустое состояние -->
      <template v-else>
        <div class="flex flex-col items-center justify-center rounded-2xl border border-white/5 bg-slate-900/30 px-6 py-16 text-center shadow-xl backdrop-blur-sm">
          <div class="mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-slate-800/50 border border-white/5 text-3xl shadow-inner">
            {{ getEmptyState().icon }}
          </div>
          <h3 class="mb-2 text-base font-semibold text-white">
            {{ getEmptyState().title }}
          </h3>
          <p class="max-w-md text-xs leading-relaxed text-slate-500">
            {{ getEmptyState().description }}
          </p>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  middleware: 'auth'
})

const notificationStore = useNotificationStore()
const router = useRouter()

const activeTab = ref('all')
const actionLoading = ref(false)

const tabs = computed(() => [
  { value: 'all', label: 'Все', icon: '📬', count: notificationStore.unreadNotificationsCount },
  { value: 'message', label: 'Сообщения', icon: '✉️', count: 0 }, // Chat has its own counts in chat menu, but message alerts are highlighted
  { value: 'grade', label: 'Оценки', icon: '📊', count: 0 },
  { value: 'attendance', label: 'Пропуски', icon: '📆', count: 0 },
  { value: 'news', label: 'Новости', icon: '📰', count: 0 }
])

const hasUnreadAlerts = computed(() => {
  return notificationStore.unreadNotificationsCount > 0
})

const filteredNotifications = computed(() => {
  if (activeTab.value === 'all') {
    return notificationStore.notifications
  }
  return notificationStore.notifications.filter(n => n.type === activeTab.value)
})

const getTypeStyles = (type) => {
  switch (type) {
    case 'message':
      return {
        icon: '✉️',
        badge: 'bg-blue-500/10 text-blue-400 border border-blue-500/10',
        borderLeft: 'border-l-4 border-l-blue-500'
      }
    case 'grade':
      return {
        icon: '📊',
        badge: 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/10',
        borderLeft: 'border-l-4 border-l-emerald-500'
      }
    case 'attendance':
      return {
        icon: '📆',
        badge: 'bg-amber-500/10 text-amber-400 border border-amber-500/10',
        borderLeft: 'border-l-4 border-l-amber-500'
      }
    case 'news':
      return {
        icon: '📰',
        badge: 'bg-purple-500/10 text-purple-400 border border-purple-500/10',
        borderLeft: 'border-l-4 border-l-purple-500'
      }
    default:
      return {
        icon: '📬',
        badge: 'bg-slate-500/10 text-slate-400 border border-slate-500/10',
        borderLeft: 'border-l-4 border-l-slate-500'
      }
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('ru-RU', {
    day: '2-digit',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getEmptyState = () => {
  switch (activeTab.value) {
    case 'message':
      return {
        icon: '✉️',
        title: 'Личных сообщений нет',
        description: 'Вы прочитали все уведомления о новых сообщениях. Общаться с одногруппниками можно во вкладке «Чат».'
      }
    case 'grade':
      return {
        icon: '📊',
        title: 'Оценок пока нет',
        description: 'Преподаватели еще не выставляли новые оценки. Посмотреть общую успеваемость можно во вкладке «Мои оценки».'
      }
    case 'attendance':
      return {
        icon: '📆',
        title: 'Все пары посещены',
        description: 'Нет записей о пропусках или опозданиях. Отличная посещаемость, продолжайте в том же духе!'
      }
    case 'news':
      return {
        icon: '📰',
        title: 'Новостей пока нет',
        description: 'Информационная лента пуста. Как только администрация опубликует важные изменения, они появятся здесь.'
      }
    default:
      return {
        icon: '💫',
        title: 'Уведомлений нет',
        description: 'У вас нет новых или архивных оповещений. Все системы работают в штатном режиме!'
      }
  }
}

const markAsRead = async (id) => {
  await notificationStore.markAsRead(id)
  await notificationStore.fetchUnreadCount()
}

const markAllAsRead = async () => {
  actionLoading.value = true
  try {
    await notificationStore.markAllAsRead()
    await notificationStore.fetchUnreadCount()
  } finally {
    actionLoading.value = false
  }
}

const handleNotificationClick = async (item) => {
  if (!item.read) {
    await notificationStore.markAsRead(item.id)
  }
  if (item.link) {
    router.push(item.link)
  }
}

onMounted(async () => {
  await notificationStore.fetchNotifications()
  await notificationStore.fetchUnreadCount()
})
</script>

<style scoped>
/* Анимации списков Vue Transition */
.list-enter-active,
.list-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateY(12px);
}
</style>
