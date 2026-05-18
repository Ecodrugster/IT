<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h2 class="text-2xl font-bold text-white">Пользователи</h2>
      <input
        v-model="search"
        type="text"
        placeholder="Поиск по имени или почте..."
        class="bg-slate-900 border border-white/10 rounded-xl px-4 py-2 text-sm text-white focus:ring-2 focus:ring-blue-500/50 w-72"
      />
    </div>

    <div v-if="errorMessage" class="bg-red-500/10 border border-red-500/30 text-red-300 rounded-xl px-4 py-3 text-sm">
      {{ errorMessage }}
    </div>

    <div class="bg-slate-900 border border-white/5 rounded-2xl overflow-hidden shadow-xl">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="bg-white/5 text-[10px] uppercase tracking-wider text-slate-500 font-bold">
            <th class="px-6 py-4">Пользователь</th>
            <th class="px-6 py-4">Почта</th>
            <th class="px-6 py-4">Роль</th>
            <th class="px-6 py-4">Группа</th>
            <th class="px-6 py-4">Действия</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-white/5">
          <tr v-for="user in filteredUsers" :key="user.uid" class="hover:bg-white/[0.02] transition-colors">
            <td class="px-6 py-4">
              <div class="flex items-center gap-3">
                <img v-if="user.photo_url" :src="user.photo_url" class="w-8 h-8 rounded-lg object-cover" />
                <div v-else class="w-8 h-8 rounded-lg bg-slate-800 flex items-center justify-center text-xs font-bold text-slate-500">
                  {{ (user.display_name || user.email || 'U')[0].toUpperCase() }}
                </div>
                <div class="text-sm font-medium text-white">{{ user.display_name || 'Без имени' }}</div>
              </div>
            </td>

            <td class="px-6 py-4 text-sm text-slate-400 font-mono">{{ user.email }}</td>

            <td class="px-6 py-4">
              <span
                class="px-2 py-1 rounded-md text-[10px] font-bold uppercase tracking-tight"
                :class="roleBadgeClass(user.role)"
              >
                {{ roleLabel(user.role) }}
              </span>
            </td>

            <td class="px-6 py-4">
              <input
                v-model="user.groupDraft"
                type="text"
                placeholder="Например: P-21"
                class="bg-slate-800 border border-white/10 rounded-lg px-3 py-2 text-xs text-white w-28"
              />
            </td>

            <td class="px-6 py-4">
              <div class="flex items-center gap-2">
                <select
                  :value="user.role || 'student'"
                  @change="updateRole(user, $event.target.value)"
                  :disabled="updatingRole === user.uid"
                  class="bg-slate-800 border-none rounded-lg text-xs font-bold text-white px-3 py-2 focus:ring-2 focus:ring-blue-500/50"
                >
                  <option value="student">Студент</option>
                  <option value="teacher">Преподаватель</option>
                  <option value="admin">Администратор</option>
                </select>

                <button
                  @click="updateGroup(user)"
                  :disabled="updatingGroup === user.uid"
                  class="px-3 py-2 rounded-lg text-xs font-semibold bg-blue-600 hover:bg-blue-500 text-white disabled:opacity-50"
                >
                  Сохранить группу
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>

      <div v-if="loading" class="p-12 text-center text-slate-500 animate-pulse">Загрузка пользователей...</div>
      <div v-else-if="filteredUsers.length === 0" class="p-12 text-center text-slate-500">Пользователи не найдены.</div>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'admin',
  middleware: 'admin'
})

const { fetchApi: api } = useApi()

const search = ref('')
const users = ref([])
const loading = ref(true)
const updatingRole = ref(null)
const updatingGroup = ref(null)
const errorMessage = ref('')

const roleLabel = (role) => {
  if (role === 'admin') return 'Администратор'
  if (role === 'teacher') return 'Преподаватель'
  return 'Студент'
}

const roleBadgeClass = (role) => {
  if (role === 'admin') return 'bg-blue-500/10 text-blue-400'
  if (role === 'teacher') return 'bg-green-500/10 text-green-400'
  return 'bg-slate-500/10 text-slate-400'
}

const prepareUser = (u) => {
  const groupName = u.group_name || u.group || ''
  return {
    ...u,
    role: u.role || 'student',
    group_name: groupName,
    group: groupName,
    groupDraft: groupName
  }
}

const fetchUsers = async () => {
  loading.value = true
  errorMessage.value = ''
  try {
    const data = await api('/admin/users?limit=200')
    const list = data.users || []
    users.value = list.map(prepareUser)
  } catch (e) {
    console.error('Не удалось получить пользователей:', e)
    const status = e?.status || e?.response?.status
    const message = e?.data?.error || e?.message || 'Не удалось загрузить пользователей'
    errorMessage.value = `Не удалось загрузить пользователей (${status || 'без статуса'}): ${message}`
  } finally {
    loading.value = false
  }
}

const filteredUsers = computed(() => {
  if (!search.value) return users.value
  const s = search.value.toLowerCase()
  return users.value.filter(u =>
    u.display_name?.toLowerCase().includes(s) ||
    u.email?.toLowerCase().includes(s)
  )
})

const updateRole = async (user, newRole) => {
  if (!confirm(`Изменить роль для ${user.display_name || user.email} на "${roleLabel(newRole)}"?`)) return

  updatingRole.value = user.uid
  try {
    await api(`/admin/users/${user.uid}/role`, {
      method: 'PUT',
      body: { role: newRole }
    })
    user.role = newRole
    alert('Роль обновлена')
  } catch (e) {
    alert('Ошибка обновления роли: ' + (e?.data?.error || e.message))
  } finally {
    updatingRole.value = null
  }
}

const updateGroup = async (user) => {
  updatingGroup.value = user.uid
  try {
    const groupName = (user.groupDraft || '').trim()
    await api(`/admin/users/${user.uid}/group`, {
      method: 'PUT',
      body: { group_name: groupName }
    })
    user.group_name = groupName
    user.group = groupName
    alert('Группа сохранена')
  } catch (e) {
    alert('Ошибка обновления группы: ' + (e?.data?.error || e.message))
  } finally {
    updatingGroup.value = null
  }
}

onMounted(fetchUsers)
</script>
