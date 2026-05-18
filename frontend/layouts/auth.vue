<template>
  <div class="min-h-screen bg-slate-950 flex items-center justify-center p-4">
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="absolute -top-[10%] -left-[10%] w-[40%] h-[40%] rounded-full bg-blue-500/20 blur-[120px]"></div>
      <div class="absolute -bottom-[10%] -right-[10%] w-[40%] h-[40%] rounded-full bg-purple-500/20 blur-[120px]"></div>
    </div>
    <div class="w-full max-w-md relative z-10">
      <slot />
    </div>
  </div>
</template>
