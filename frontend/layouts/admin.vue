<template>
  <div class="flex min-h-screen bg-slate-950 text-slate-200">
    <AdminSidebar />

    <main class="flex h-screen flex-grow flex-col overflow-hidden">
      <header class="flex h-16 items-center justify-between border-b border-white/5 bg-slate-950/60 px-6 backdrop-blur-md md:px-8">
        <h3 class="text-sm font-medium text-slate-300">
          {{ currentPageTitle }}
        </h3>
        <div class="flex items-center gap-3">
          <NuxtLink
            to="/"
            class="hidden rounded-lg bg-white/5 px-3 py-1.5 text-xs text-slate-300 transition-all hover:bg-white/10 hover:text-white sm:inline-flex"
          >
            Главная
          </NuxtLink>
          <div class="text-right">
            <div class="text-xs font-bold text-white">{{ userDisplayName }}</div>
            <div class="text-[10px] font-mono uppercase tracking-wider text-blue-500">Админ</div>
          </div>
          <img
            v-if="userStore.user?.photoURL"
            :src="userStore.user.photoURL"
            class="h-8 w-8 rounded-lg border border-white/10"
            alt="avatar"
          />
          <div v-else class="flex h-8 w-8 items-center justify-center rounded-lg bg-slate-800 text-xs font-bold">
            A
          </div>
        </div>
      </header>

      <div class="custom-scrollbar flex-grow overflow-y-auto p-6 md:p-8">
        <slot />
      </div>
    </main>
  </div>
</template>

<script setup>
const userStore = useUserStore()
const route = useRoute()

const userDisplayName = computed(() => {
  return userStore.user?.displayName || userStore.profile?.display_name || userStore.profile?.displayName || userStore.profile?.email || 'Администратор'
})

const currentPageTitle = computed(() => {
  const titles = {
    '/admin': 'Сводка системы',
    '/admin/users': 'Управление пользователями',
    '/admin/schedule': 'Управление расписанием',
    '/admin/news': 'Управление новостями',
    '/admin/clubs': 'Управление клубами',
    '/admin/moderation': 'Модерация контента'
  }
  return titles[route.path] || 'Панель администратора'
})
</script>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.08);
  border-radius: 10px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.14);
}
</style>
