<template>
  <div class="bg-slate-900/50 backdrop-blur-xl border border-white/10 p-8 rounded-2xl shadow-2xl">
    <div class="text-center mb-8">
      <h1 class="text-3xl font-bold text-white mb-2">ITSTEP Social</h1>
      <p class="text-slate-400">Войдите в свой аккаунт</p>
    </div>

    <form @submit.prevent="handleLogin" class="space-y-6">
      <div>
        <label class="block text-sm font-medium text-slate-300 mb-2">Электронная почта</label>
        <input 
          v-model="email"
          type="email" 
          required
          class="w-full bg-slate-800 border border-white/5 rounded-lg px-4 py-3 text-white focus:outline-none focus:ring-2 focus:ring-blue-500/50 transition-all"
          placeholder="student@itstep.org"
        />
      </div>

      <div>
        <label class="block text-sm font-medium text-slate-300 mb-2">Пароль</label>
        <input 
          v-model="password"
          type="password" 
          required
          class="w-full bg-slate-800 border border-white/5 rounded-lg px-4 py-3 text-white focus:outline-none focus:ring-2 focus:ring-blue-500/50 transition-all"
          placeholder="••••••••"
        />
      </div>

      <button 
        type="submit"
        :disabled="loading"
        class="w-full bg-blue-600 hover:bg-blue-500 text-white font-semibold py-3 rounded-lg transition-all disabled:opacity-50"
      >
        {{ loading ? 'Вход...' : 'Войти' }}
      </button>

      <div class="relative py-4">
        <div class="absolute inset-0 flex items-center">
          <div class="w-full border-t border-white/10"></div>
        </div>
        <div class="relative flex justify-center text-xs uppercase">
          <span class="bg-slate-900 px-2 text-slate-500">Или</span>
        </div>
      </div>

      <button 
        type="button"
        @click="handleGoogleLogin"
        :disabled="loading"
        class="w-full bg-white/5 hover:bg-white/10 border border-white/10 text-white font-medium py-3 rounded-lg transition-all flex items-center justify-center space-x-3"
      >
        <svg class="w-5 h-5" viewBox="0 0 24 24">
          <path fill="currentColor" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" />
          <path fill="currentColor" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" />
          <path fill="currentColor" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l3.66-2.84z" />
          <path fill="currentColor" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" />
        </svg>
        <span>Войти через Google</span>
      </button>
    </form>

    <div class="mt-6 text-center text-sm">
      <span class="text-slate-400">Нет аккаунта? </span>
      <NuxtLink to="/register" class="text-blue-400 hover:text-blue-300 font-medium">Создать аккаунт</NuxtLink>
    </div>

    <div v-if="error" class="mt-4 p-3 bg-red-500/10 border border-red-500/20 rounded-lg text-red-400 text-sm text-center">
      {{ error }}
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'auth'
})

const { login, loginWithGoogle } = useAuth()
const router = useRouter()

const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  loading.value = true
  error.value = ''
  try {
    await login(email.value, password.value)
    router.push('/')
  } catch (e) {
    error.value = 'Ошибка входа: ' + e.message
  } finally {
    loading.value = false
  }
}

const handleGoogleLogin = async () => {
  loading.value = true
  error.value = ''
  try {
    await loginWithGoogle()
    router.push('/')
  } catch (e) {
    error.value = 'Ошибка входа через Google: ' + e.message
  } finally {
    loading.value = false
  }
}
</script>
