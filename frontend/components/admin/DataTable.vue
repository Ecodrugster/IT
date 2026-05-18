<template>
  <div class="bg-slate-900 border border-white/5 rounded-2xl overflow-hidden shadow-xl">
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="bg-white/5 text-[10px] uppercase tracking-wider text-slate-500 font-bold">
            <th v-for="col in columns" :key="col.key" class="px-6 py-4">{{ col.label }}</th>
            <th v-if="$slots.actions" class="px-6 py-4 text-right">Действия</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-white/5">
          <tr v-for="(item, index) in data" :key="index" class="hover:bg-white/[0.02] transition-colors">
            <td v-for="col in columns" :key="col.key" class="px-6 py-4 text-sm">
              <slot :name="`cell-${col.key}`" :item="item">
                <span class="text-slate-300">{{ item[col.key] }}</span>
              </slot>
            </td>
            <td v-if="$slots.actions" class="px-6 py-4 text-right">
              <slot name="actions" :item="item"></slot>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    
    <!-- Empty State -->
    <div v-if="data.length === 0 && !loading" class="p-12 text-center text-slate-500 text-sm">
      Данные не найдены
    </div>
    
    <!-- Loading State -->
    <div v-if="loading" class="p-12 text-center text-slate-500 animate-pulse text-sm">
      Загрузка...
    </div>
  </div>
</template>

<script setup>
defineProps({
  columns: Array,
  data: Array,
  loading: Boolean
})
</script>
