<template>
  <div v-if="show" class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-slate-950/90 backdrop-blur-sm">
    <div class="bg-slate-900 border border-white/10 rounded-2xl p-8 max-w-sm w-full shadow-2xl text-center">
      <div class="w-16 h-16 bg-red-500/10 text-red-500 rounded-full flex items-center justify-center text-2xl mx-auto mb-4">⚠️</div>
      <h3 class="text-xl font-bold text-white mb-2">Вы уверены?</h3>
      <p class="text-slate-400 text-sm mb-8">{{ message || 'Это действие нельзя будет отменить.' }}</p>
      
      <div class="flex gap-4">
        <button @click="$emit('cancel')" class="flex-grow py-3 bg-white/5 hover:bg-white/10 text-white rounded-xl transition-all font-medium">Отмена</button>
        <button @click="$emit('confirm')" class="flex-grow py-3 bg-red-600 hover:bg-red-500 text-white rounded-xl font-bold shadow-lg shadow-red-600/20 transition-all">Удалить</button>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  show: Boolean,
  message: String
})
defineEmits(['confirm', 'cancel'])
</script>
