<template>
  <div class="max-w-6xl mx-auto py-12 px-4">
    <div class="text-center mb-16">
      <h1 class="text-4xl font-extrabold text-white mb-4 bg-gradient-to-r from-blue-500 to-green-500 bg-clip-text text-transparent">
        Навигатор студента ITSTEP
      </h1>
      <p class="text-slate-400 max-w-2xl mx-auto">
        Всё, что нужно знать о жизни в колледже: от расписания до поиска преподавателя и администратора.
      </p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
      <div class="bg-slate-900 border border-white/5 rounded-3xl p-8 hover:border-blue-500/30 transition-all group">
        <div class="w-14 h-14 bg-blue-600/10 rounded-2xl flex items-center justify-center text-2xl mb-6 group-hover:scale-110 transition-transform">
          📅
        </div>
        <h2 class="text-xl font-bold text-white mb-4">Расписание</h2>
        <p class="text-slate-400 text-sm mb-6">
          Актуальное расписание занятий, зачетов и экзаменов всегда доступно онлайн.
        </p>
        <NuxtLink to="/schedule" class="inline-flex items-center text-blue-500 font-bold hover:underline">
          Перейти к расписанию <span class="ml-2">→</span>
        </NuxtLink>
      </div>

      <div class="bg-slate-900 border border-white/5 rounded-3xl p-8 hover:border-green-500/30 transition-all group">
        <div class="w-14 h-14 bg-green-600/10 rounded-2xl flex items-center justify-center text-2xl mb-6 group-hover:scale-110 transition-transform">
          🎭
        </div>
        <h2 class="text-xl font-bold text-white mb-4">Кружки и активности</h2>
        <p class="text-slate-400 text-sm mb-6">
          Хочешь развиваться вне учебы? У нас есть клубы по интересам: от программирования до театра.
        </p>
        <NuxtLink to="/clubs" class="inline-flex items-center text-green-500 font-bold hover:underline">
          Посмотреть все клубы <span class="ml-2">→</span>
        </NuxtLink>
      </div>

      <div class="bg-slate-900 border border-white/5 rounded-3xl p-8 hover:border-purple-500/30 transition-all group">
        <div class="w-14 h-14 bg-purple-600/10 rounded-2xl flex items-center justify-center text-2xl mb-6 group-hover:scale-110 transition-transform">
          📞
        </div>
        <h2 class="text-xl font-bold text-white mb-4">Куда обращаться?</h2>
        <div class="space-y-3">
          <div class="flex justify-between text-xs">
            <span class="text-slate-500">Учебная часть:</span>
            <span class="text-white">+7 (XXX) XXX-XX-XX</span>
          </div>
          <div class="flex justify-between text-xs">
            <span class="text-slate-500">IT-поддержка:</span>
            <span class="text-white">support@itstep.edu</span>
          </div>
          <div class="flex justify-between text-xs">
            <span class="text-slate-500">Библиотека:</span>
            <span class="text-white">Кабинет 302</span>
          </div>
        </div>
      </div>
    </div>

    <div class="mt-16 bg-slate-900/50 border border-white/5 rounded-3xl p-10 overflow-hidden relative">
      <div class="absolute -right-20 -top-20 w-64 h-64 bg-blue-600/10 rounded-full blur-3xl"></div>
      
      <div class="relative z-10">
        <h2 class="text-2xl font-bold text-white mb-8">Как найти преподавателя или администратора?</h2>
        <p class="text-slate-400 mb-8 max-w-xl">
          Используйте поиск по пользователям. Можно фильтровать по роли «Учитель» и «Админ».
        </p>
        
        <div class="flex flex-wrap gap-4">
          <NuxtLink to="/search?role=teacher" class="bg-white/5 hover:bg-white/10 px-6 py-3 rounded-xl text-white font-medium transition-all">
            Найти преподавателя
          </NuxtLink>
          <NuxtLink to="/search?role=admin" class="bg-white/5 hover:bg-white/10 px-6 py-3 rounded-xl text-white font-medium transition-all">
            Найти администратора
          </NuxtLink>
        </div>
      </div>
    </div>

    <div class="mt-16">
      <h2 class="text-2xl font-bold text-white mb-8 text-center">Часто задаваемые вопросы</h2>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div v-for="(item, i) in faq" :key="i" class="bg-slate-900 border border-white/5 p-6 rounded-2xl">
          <h4 class="text-white font-bold mb-2">{{ item.q }}</h4>
          <p class="text-slate-500 text-sm">{{ item.a }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
const faq = [
  { q: 'Как получить справку об обучении?', a: 'Обратитесь в учебную часть (каб. 201) в приемные часы: с 10:00 до 16:00.' },
  { q: 'Где взять пароль от Wi-Fi?', a: 'Пароль можно узнать на стенде в IT-коворкинге или у системного администратора.' },
  { q: 'Как создать свой клуб?', a: 'Перейдите в раздел «Клубы» и нажмите «Создать». После одобрения админом ваш клуб станет публичным.' },
  { q: 'Что делать, если потерял студенческий?', a: 'Напишите заявление в деканате. Стоимость восстановления — 200 StepCoins.' }
]
</script>
