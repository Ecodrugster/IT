<template>
  <div class="mx-auto h-[calc(100vh-120px)] max-w-6xl px-4 py-8">
    <div class="flex h-full overflow-hidden rounded-xl border border-white/5 bg-slate-900 shadow-xl">
      <!-- Список пользователей -->
      <div class="flex w-80 flex-col border-r border-white/5">
        <div class="border-b border-white/5 p-4">
          <h3 class="mb-2 font-semibold text-white">Студенты</h3>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Поиск..."
            class="w-full rounded-lg border-none bg-slate-800 px-4 py-2 text-sm text-white focus:ring-1 focus:ring-blue-500/50"
          />
        </div>

        <div class="flex-grow overflow-y-auto">
          <div
            v-for="user in filteredUsers"
            :key="user.uid"
            class="flex cursor-pointer items-center space-x-3 border-b border-white/5 p-4 transition-all last:border-0 hover:bg-white/5"
            :class="{ 'border-l-2 border-l-blue-600 bg-blue-600/10': selectedUser?.uid === user.uid }"
            @click="selectUser(user)"
          >
            <img v-if="user.photoURL" :src="user.photoURL" class="h-10 w-10 rounded-full" />
            <div
              v-else
              class="flex h-10 w-10 items-center justify-center rounded-full bg-slate-800 text-xs font-bold text-slate-500"
            >
              {{ (user.displayName || user.email || 'U')[0].toUpperCase() }}
            </div>

            <div class="min-w-0 flex-grow">
              <h4 class="truncate text-sm font-medium text-white">{{ user.displayName || user.email }}</h4>
              <p class="text-[10px] uppercase text-slate-500">Студент</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Окно переписки -->
      <div class="flex flex-grow flex-col">
        <template v-if="selectedUser">
          <div class="flex items-center justify-between border-b border-white/5 bg-slate-900/50 p-4">
            <div class="flex items-center space-x-3">
              <img v-if="selectedUser.photoURL" :src="selectedUser.photoURL" class="h-10 w-10 rounded-full" />
              <div
                v-else
                class="flex h-10 w-10 items-center justify-center rounded-full bg-slate-800 text-xs font-bold text-white"
              >
                {{ (selectedUser.displayName || selectedUser.email || 'U')[0].toUpperCase() }}
              </div>

              <div>
                <h4 class="text-sm font-medium text-white">{{ selectedUser.displayName || selectedUser.email }}</h4>
                <p class="text-[10px] font-bold uppercase text-green-500">В сети</p>
              </div>
            </div>
          </div>

          <div ref="messageContainer" class="flex-grow space-y-4 overflow-y-auto bg-slate-950/30 p-6">
            <div
              v-for="msg in messages"
              :key="msg.id"
              class="flex"
              :class="msg.senderId === userStore.user.uid ? 'justify-end' : 'justify-start'"
            >
              <div
                class="max-w-[72%] rounded-2xl shadow-lg"
                :class="[
                  msg.senderId === userStore.user.uid
                    ? 'rounded-tr-none bg-blue-600 text-white'
                    : 'rounded-tl-none bg-slate-800 text-slate-200',
                  msg.imageUrl ? 'px-2 py-2' : 'px-4 py-2'
                ]"
              >
                <a
                  v-if="msg.imageUrl"
                  :href="msg.imageUrl"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="block"
                >
                  <img
                    :src="msg.imageUrl"
                    class="max-h-80 w-auto rounded-xl border border-white/10 object-cover"
                    loading="lazy"
                    alt="Изображение в чате"
                  />
                </a>

                <p v-if="msg.text" class="break-words whitespace-pre-wrap" :class="msg.imageUrl ? 'mt-2 px-2 pb-1' : ''">
                  {{ msg.text }}
                </p>

                <div class="mt-1 text-right text-[9px] opacity-50" :class="msg.imageUrl ? 'px-2 pb-1' : ''">
                  {{ msg.createdAt ? new Date(msg.createdAt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) : '...' }}
                </div>
              </div>
            </div>
          </div>

          <div class="border-t border-white/5 p-4">
            <input
              ref="galleryInputRef"
              type="file"
              accept="image/*"
              class="hidden"
              @change="onImageSelected"
            />
            <input
              ref="cameraInputRef"
              type="file"
              accept="image/*"
              capture="environment"
              class="hidden"
              @change="onImageSelected"
            />

            <div
              v-if="selectedImagePreview"
              class="mb-3 flex items-center gap-3 rounded-lg border border-white/10 bg-slate-800/70 p-3"
            >
              <img
                :src="selectedImagePreview"
                class="h-16 w-16 rounded-lg border border-white/10 object-cover"
                alt="Превью изображения"
              />
              <div class="min-w-0 flex-grow">
                <div class="truncate text-xs text-white">{{ selectedImageFile?.name }}</div>
                <div class="text-[11px] text-slate-400">Фото готово к отправке</div>
              </div>
              <button
                type="button"
                class="rounded-md border border-white/10 px-3 py-1 text-xs text-slate-300 transition-colors hover:bg-white/10 hover:text-white"
                @click="clearSelectedImage"
              >
                Убрать
              </button>
            </div>

            <div
              v-if="sendError"
              class="mb-3 rounded-lg border border-red-500/30 bg-red-500/10 px-3 py-2 text-xs text-red-300"
            >
              {{ sendError }}
            </div>

            <form class="flex items-center gap-2" @submit.prevent="sendMessage">
              <button
                type="button"
                class="rounded-lg border border-white/10 px-3 py-3 text-xs font-medium text-slate-300 transition-colors hover:bg-white/10 hover:text-white"
                @click="openGallery"
              >
                Фото
              </button>
              
              <input
                v-model="newMessage"
                type="text"
                placeholder="Введите сообщение или добавьте фото..."
                class="flex-grow rounded-lg border-none bg-slate-800 px-4 py-3 text-white focus:ring-1 focus:ring-blue-500/50"
                @input="sendError = ''"
              />

              <button
                type="submit"
                :disabled="!canSend"
                class="rounded-lg bg-blue-600 px-4 py-3 text-sm text-white shadow-lg shadow-blue-600/20 transition-all hover:bg-blue-500 disabled:opacity-50"
              >
                {{ isSending ? 'Отправка...' : 'Отправить' }}
              </button>
            </form>

            <div class="mt-2 flex items-center justify-between text-[11px] text-slate-500">
              <span>Фото до {{ MAX_IMAGE_SIZE_MB }} МБ</span>
              <span :class="messageLength > MAX_MESSAGE_LENGTH ? 'text-red-400' : 'text-slate-500'">
                {{ messageLength }} / {{ MAX_MESSAGE_LENGTH }}
              </span>
            </div>
          </div>
        </template>

        <div v-else class="flex flex-grow flex-col items-center justify-center space-y-4 text-slate-500">
          <div class="text-6xl text-slate-800">💬</div>
          <p class="italic">Выберите студента, чтобы начать диалог</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  middleware: 'auth'
})

const userStore = useUserStore()
const { fetchApi: api } = useApi()
const route = useRoute()

const MAX_MESSAGE_LENGTH = 10000
const MAX_IMAGE_SIZE_BYTES = 10 * 1024 * 1024
const MAX_IMAGE_SIZE_MB = Math.floor(MAX_IMAGE_SIZE_BYTES / (1024 * 1024))

const users = ref([])
const searchQuery = ref('')
const selectedUser = ref(null)
const messages = ref([])
const newMessage = ref('')
const sendError = ref('')
const isSending = ref(false)
const selectedImageFile = ref(null)
const selectedImagePreview = ref('')

const messageContainer = ref(null)
const galleryInputRef = ref(null)
const cameraInputRef = ref(null)

let pollTimer = null

const messageLength = computed(() => Array.from(newMessage.value).length)

const canSend = computed(() => {
  return !!selectedUser.value && !isSending.value && (!!newMessage.value.trim() || !!selectedImageFile.value)
})

const filteredUsers = computed(() => {
  return users.value.filter((u) => {
    return (
      u.uid !== userStore.user.uid &&
      (u.displayName?.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
        u.email?.toLowerCase().includes(searchQuery.value.toLowerCase()))
    )
  })
})

const normalizeMessage = (raw) => {
  const imageUrl = String(raw?.imageUrl || raw?.image_url || '').trim()
  const messageType = String(raw?.messageType || raw?.message_type || (imageUrl ? 'image' : 'text')).trim()

  return {
    ...raw,
    text: typeof raw?.text === 'string' ? raw.text : '',
    imageUrl,
    messageType
  }
}

const stopPolling = () => {
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = null
  }
}

const scrollToBottom = () => {
  nextTick(() => {
    if (messageContainer.value) {
      messageContainer.value.scrollTop = messageContainer.value.scrollHeight
    }
  })
}

const clearSelectedImage = () => {
  if (selectedImagePreview.value && selectedImagePreview.value.startsWith('blob:')) {
    URL.revokeObjectURL(selectedImagePreview.value)
  }
  selectedImagePreview.value = ''
  selectedImageFile.value = null
}

const openGallery = () => {
  galleryInputRef.value?.click()
}

const openCamera = () => {
  cameraInputRef.value?.click()
}

const onImageSelected = (event) => {
  const input = event?.target
  const file = input?.files?.[0]

  if (input) {
    input.value = ''
  }
  if (!file) return

  if (!file.type.startsWith('image/')) {
    sendError.value = 'Можно отправлять только изображения.'
    return
  }
  if (file.size > MAX_IMAGE_SIZE_BYTES) {
    sendError.value = `Фото слишком большое. Максимум ${MAX_IMAGE_SIZE_MB} МБ.`
    return
  }

  clearSelectedImage()
  selectedImageFile.value = file
  selectedImagePreview.value = URL.createObjectURL(file)
  sendError.value = ''
}

const markAsRead = async (otherUserId) => {
  try {
    await api('/chat/messages/read', {
      method: 'POST',
      body: {
        user_id: otherUserId
      }
    })
  } catch (e) {
    console.error('[Chat Debug] Failed to mark messages as read:', e)
  }
}

const fetchMessages = async () => {
  if (!selectedUser.value) return

  try {
    const otherUserId = selectedUser.value.uid
    const data = await api(`/chat/messages?user_id=${encodeURIComponent(otherUserId)}&limit=200`)
    const nextMessages = (data || []).map(normalizeMessage)

    const prevLastId = messages.value.length ? messages.value[messages.value.length - 1].id : null
    const nextLastId = nextMessages.length ? nextMessages[nextMessages.length - 1].id : null

    messages.value = nextMessages

    await markAsRead(otherUserId)

    if (prevLastId !== nextLastId) {
      scrollToBottom()
    }
  } catch (e) {
    console.error('[Chat Debug] Failed to fetch messages:', e)
  }
}

const startPollingMessages = () => {
  stopPolling()
  fetchMessages()
  pollTimer = setInterval(fetchMessages, 2500)
}

const fetchUsers = async () => {
  try {
    const data = await api('/users')
    users.value = data || []

    const targetUID = typeof route.query.uid === 'string' ? route.query.uid.trim() : ''
    if (targetUID) {
      const userFromQuery = users.value.find((u) => u.uid === targetUID)
      if (userFromQuery) {
        selectUser(userFromQuery)
      }
    }
  } catch (e) {
    console.error('[Chat Debug] Failed to fetch users:', e)
  }
}

const selectUser = (user) => {
  messages.value = []
  sendError.value = ''
  selectedUser.value = user
  startPollingMessages()
}

const formatSendError = (e) => {
  const backendMessage = String(e?.data?.error || e?.data?.message || e?.message || 'Неизвестная ошибка')
  const lowered = backendMessage.toLowerCase()

  if (/too long|maximum|длин/.test(lowered)) {
    const match = backendMessage.match(/(\d{2,6})/)
    const backendLimit = match ? Number(match[1]) : null
    const safeLimit = Number.isFinite(backendLimit) ? backendLimit : MAX_MESSAGE_LENGTH
    return `Сообщение слишком длинное. Максимум ${safeLimit} символов.`
  }

  if (/image is too large|слишком больш|maximum is \d+ mb/.test(lowered)) {
    const match = backendMessage.match(/(\d{1,2})\s*mb/i)
    const maxMb = match ? Number(match[1]) : MAX_IMAGE_SIZE_MB
    return `Фото слишком большое. Максимум ${maxMb} МБ.`
  }

  if (/unsupported image format|неподдерж/.test(lowered)) {
    return 'Неподдерживаемый формат фото. Используйте JPG, PNG, WEBP, GIF или HEIC.'
  }

  return `Не удалось отправить сообщение: ${backendMessage}`
}

const sendMessage = async () => {
  if (!selectedUser.value || isSending.value) return

  const text = newMessage.value.trim()
  if (!text && !selectedImageFile.value) return

  if (messageLength.value > MAX_MESSAGE_LENGTH) {
    sendError.value = `Сообщение слишком длинное. Максимум ${MAX_MESSAGE_LENGTH} символов.`
    return
  }

  const otherUserId = selectedUser.value.uid
  isSending.value = true
  sendError.value = ''

  try {
    if (selectedImageFile.value) {
      const formData = new FormData()
      formData.append('receiver_id', otherUserId)
      formData.append('image', selectedImageFile.value)
      if (text) {
        formData.append('text', text)
      }

      await api('/chat/messages/image', {
        method: 'POST',
        body: formData
      })
    } else {
      await api('/chat/messages', {
        method: 'POST',
        body: {
          receiver_id: otherUserId,
          text
        }
      })
    }

    newMessage.value = ''
    clearSelectedImage()

    await fetchMessages()
    scrollToBottom()
  } catch (e) {
    sendError.value = formatSendError(e)
    console.error('[Chat Debug] Error sending message:', e)
  } finally {
    isSending.value = false
  }
}

onMounted(() => {
  fetchUsers()
})

onUnmounted(() => {
  stopPolling()
  clearSelectedImage()
})
</script>
