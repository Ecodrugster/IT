<template>
  <aside class="flex h-full w-72 flex-col border-r border-white/5 bg-slate-900">
    <div class="border-b border-white/5 p-6">
      <div class="flex items-center justify-between gap-3">
        <div class="flex items-center gap-3">
          <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-blue-600 font-bold text-white">A</div>
          <h2 class="text-lg font-bold text-white">Панель администратора</h2>
        </div>
        <NuxtLink
          to="/"
          class="rounded-lg bg-white/5 px-3 py-1.5 text-xs text-slate-300 transition-all hover:bg-white/10 hover:text-white"
        >
          На сайт
        </NuxtLink>
      </div>
    </div>

    <nav class="flex-grow space-y-2 p-4">
      <NuxtLink
        v-for="item in menuItems"
        :key="item.path"
        :to="item.path"
        class="flex items-center gap-3 rounded-xl px-4 py-3 text-sm font-medium transition-all"
        :class="isActive(item.path) ? 'bg-blue-600 text-white shadow-lg shadow-blue-600/20' : 'text-slate-400 hover:bg-white/5 hover:text-white'"
      >
        <span class="w-5 text-center text-xs font-bold">{{ item.icon }}</span>
        {{ item.label }}
      </NuxtLink>
    </nav>

    <div class="border-t border-white/5 p-4">
      <NuxtLink to="/" class="flex items-center gap-3 rounded-xl px-4 py-3 text-sm text-slate-400 transition-all hover:bg-white/5 hover:text-white">
        <span class="w-5 text-center">&lt;</span>
        На сайт
      </NuxtLink>
    </div>
  </aside>
</template>

<script setup>
const route = useRoute()

const menuItems = [
  { label: 'Обзор', path: '/admin', icon: 'DB' },
  { label: 'Пользователи', path: '/admin/users', icon: 'US' },
  { label: 'Расписание', path: '/admin/schedule', icon: 'SC' },
  { label: 'Новости', path: '/admin/news', icon: 'NW' },
  { label: 'Клубы', path: '/admin/clubs', icon: 'CL' },
  { label: 'Модерация', path: '/admin/moderation', icon: 'MD' }
]

const isActive = (path) => {
  if (path === '/admin') {
    return route.path === '/admin'
  }
  return route.path.startsWith(path)
}
</script>
