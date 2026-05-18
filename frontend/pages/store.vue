<template>
  <div class="mx-auto min-h-[calc(100vh-120px)] max-w-6xl px-4 py-8 relative">
    <!-- Header -->
    <div class="mb-8 flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
      <div>
        <h1 class="text-3xl font-bold bg-gradient-to-r from-white to-slate-400 bg-clip-text text-transparent flex items-center gap-2">
          <span>🛍️</span> Магазин ITSTEP
        </h1>
        <p class="text-xs text-slate-500 mt-1">
          Обменивайте заработанные за посещаемость монеты и звезды на фирменный мерч и технику!
        </p>
      </div>

      <!-- Navigation Tabs -->
      <div class="flex items-center bg-slate-950 p-1 rounded-full border border-white/5 self-start">
        <button 
          @click="activeTab = 'shop'"
          class="px-4 py-2 text-xs font-medium rounded-full transition-all"
          :class="activeTab === 'shop' ? 'bg-slate-900 border border-white/10 text-white shadow-md' : 'text-slate-400 hover:text-white'"
        >
          Витрина
        </button>
        <button 
          v-if="userStore.role === 'student'"
          @click="activeTab = 'purchases'"
          class="px-4 py-2 text-xs font-medium rounded-full transition-all"
          :class="activeTab === 'purchases' ? 'bg-slate-900 border border-white/10 text-white shadow-md' : 'text-slate-400 hover:text-white'"
        >
          Мои покупки
        </button>
        <button 
          v-if="userStore.role === 'admin'"
          @click="activeTab = 'admin_dashboard'"
          class="px-4 py-2 text-xs font-medium rounded-full transition-all"
          :class="activeTab === 'admin_dashboard' ? 'bg-red-500/10 border border-red-500/20 text-red-300 shadow-md' : 'text-slate-400 hover:text-white'"
        >
          ⚙️ Управление заказами
        </button>
      </div>
    </div>

    <!-- Student Balance Widget / Admin Banner -->
    <ClientOnly>
      <!-- Student view -->
      <div v-if="userStore.role === 'student'" class="mb-8 p-6 bg-gradient-to-r from-slate-900 to-slate-950 border border-white/5 rounded-2xl flex flex-col sm:flex-row items-center justify-between gap-6 relative overflow-hidden shadow-2xl">
        <!-- Glow effect -->
        <div class="absolute -right-16 -top-16 w-32 h-32 bg-yellow-500/10 rounded-full blur-3xl"></div>
        <div class="absolute -left-16 -bottom-16 w-32 h-32 bg-blue-500/10 rounded-full blur-3xl"></div>

        <div class="flex items-center gap-4">
          <div class="w-12 h-12 rounded-full bg-blue-500/10 flex items-center justify-center text-2xl border border-blue-500/20">
            🎓
          </div>
          <div>
            <h3 class="text-white font-semibold text-lg">Ваша активность</h3>
            <p class="text-xs text-slate-400 mt-0.5">
              Серия посещений без пропусков: <span class="text-emerald-400 font-bold">{{ userStore.profile?.streak || 0 }} / 5</span> 
              <span class="text-slate-500 ml-1">(при серии 5 вы получите звезду 🌟)</span>
            </p>
          </div>
        </div>

        <div class="flex items-center gap-6">
          <!-- Coins Box -->
          <div class="flex items-center gap-3 bg-slate-900/60 border border-white/5 px-6 py-3 rounded-2xl shadow-inner">
            <span class="text-3xl filter drop-shadow">🪙</span>
            <div>
              <div class="text-[10px] text-slate-500 uppercase tracking-wider">Монеты</div>
              <div class="text-xl font-bold text-yellow-400">{{ userStore.profile?.coins || 0 }}</div>
            </div>
          </div>

          <!-- Stars Box -->
          <div class="flex items-center gap-3 bg-slate-900/60 border border-white/5 px-6 py-3 rounded-2xl shadow-inner">
            <span class="text-3xl filter drop-shadow">🌟</span>
            <div>
              <div class="text-[10px] text-slate-500 uppercase tracking-wider">Звезды</div>
              <div class="text-xl font-bold text-amber-500">{{ userStore.profile?.stars || 0 }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Admin view -->
      <div v-else-if="userStore.role === 'admin'" class="mb-8 p-6 bg-gradient-to-r from-red-950/20 to-slate-950 border border-red-500/10 rounded-2xl flex flex-col sm:flex-row items-center justify-between gap-6 relative overflow-hidden shadow-2xl">
        <!-- Glow effect -->
        <div class="absolute -right-16 -top-16 w-32 h-32 bg-red-500/10 rounded-full blur-3xl"></div>
        <div class="flex items-center gap-4">
          <div class="w-12 h-12 rounded-full bg-red-500/10 flex items-center justify-center text-2xl border border-red-500/20">
            ⚙️
          </div>
          <div>
            <h3 class="text-red-400 font-semibold text-lg">Панель администратора магазина</h3>
            <p class="text-xs text-slate-400 mt-0.5 leading-relaxed">
              Вы можете полноценно просматривать витрину товаров, искать выписанные промокоды студентов и подтверждать выдачу мерча.
            </p>
          </div>
        </div>
      </div>

      <!-- Teacher Warning Card -->
      <div v-else class="mb-8 p-6 bg-blue-900/10 border border-blue-500/20 rounded-2xl flex items-start gap-4">
        <div class="text-3xl">ℹ️</div>
        <div>
          <h3 class="text-blue-400 font-semibold">Режим демонстрации</h3>
          <p class="text-xs text-slate-400 mt-1 leading-relaxed">
            Вы вошли под ролью «{{ userStore.roleLabel }}». Магазин и накопление наград доступны только студентам. Вы можете просматривать витрину, но покупки заблокированы.
          </p>
        </div>
      </div>
    </ClientOnly>

    <!-- SHOP TAB -->
    <div v-if="activeTab === 'shop'">
      <div v-if="loading" class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div v-for="i in 6" :key="i" class="bg-slate-900 border border-white/5 rounded-2xl p-6 space-y-4 animate-pulse">
          <div class="h-32 bg-slate-800 rounded-xl"></div>
          <div class="h-6 bg-slate-800 rounded w-2/3"></div>
          <div class="h-4 bg-slate-800 rounded w-full"></div>
          <div class="h-10 bg-slate-800 rounded-lg"></div>
        </div>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div 
          v-for="item in items" 
          :key="item.id" 
          class="bg-slate-900/60 backdrop-blur-xl border border-white/5 hover:border-blue-500/20 rounded-2xl p-6 flex flex-col justify-between group hover:scale-[1.01] transition-all duration-300 relative shadow-lg"
        >
          <div>
            <!-- Item Icon / Banner -->
            <div class="w-full h-36 bg-slate-950 rounded-xl mb-4 flex items-center justify-center group-hover:scale-105 transition-transform duration-300 select-none relative overflow-hidden border border-white/5">
              <!-- Radial background glow -->
              <div class="absolute inset-0 bg-radial-glow group-hover:opacity-100 transition-opacity duration-300"></div>
              
              <!-- Product image element -->
              <img 
                v-if="isImageFile(item.image)" 
                :src="item.image" 
                class="w-full h-full object-cover z-10" 
                alt="Product Image"
              />
              <!-- Emoji fallback -->
              <span v-else class="z-10 filter drop-shadow-md text-6xl">{{ item.image }}</span>
            </div>

            <div class="flex justify-between items-start gap-2 mb-2">
              <h3 class="text-white font-bold text-lg group-hover:text-blue-400 transition-colors">{{ item.name }}</h3>
              <span class="text-[10px] text-slate-500 uppercase tracking-widest px-2 py-0.5 bg-slate-950 border border-white/5 rounded-full">
                осталось: {{ item.stock }}
              </span>
            </div>

            <p class="text-xs text-slate-400 leading-relaxed mb-6">
              {{ item.description }}
            </p>
          </div>

          <!-- Price & Buy Button -->
          <div class="space-y-4 pt-4 border-t border-white/5">
            <div class="flex items-center justify-between">
              <span class="text-xs text-slate-500">Стоимость:</span>
              <div class="flex items-center gap-3">
                <div class="flex items-center gap-1">
                  <span>🪙</span>
                  <span class="font-bold text-sm" :class="canAffordCoins(item) ? 'text-yellow-400' : 'text-red-500/80'">
                    {{ item.coins_price }}
                  </span>
                </div>
                <div v-if="item.stars_price > 0" class="flex items-center gap-1">
                  <span>🌟</span>
                  <span class="font-bold text-sm" :class="canAffordStars(item) ? 'text-amber-500' : 'text-red-500/80'">
                    {{ item.stars_price }}
                  </span>
                </div>
              </div>
            </div>

            <div v-if="userStore.role === 'student'" class="flex gap-2 w-full">
              <!-- Add to Cart -->
              <button
                @click="addToCart(item)"
                class="flex-grow py-3 bg-slate-950 border border-white/5 text-slate-300 hover:text-white hover:bg-slate-900 rounded-xl font-semibold text-xs transition-all flex items-center justify-center gap-1.5 shadow-inner"
              >
                <span>🛒</span> В корзину
              </button>
              
              <!-- Quick Exchange -->
              <button
                @click="confirmPurchase(item)"
                :disabled="!canAfford(item) || purchaseLoading"
                class="flex-grow py-3 rounded-xl font-semibold text-xs transition-all flex items-center justify-center gap-1.5"
                :class="[
                  !canAfford(item)
                    ? 'bg-slate-950 border border-white/5 text-slate-500 cursor-not-allowed'
                    : 'bg-blue-600 hover:bg-blue-500 text-white shadow-lg shadow-blue-600/10 hover:shadow-blue-500/20'
                ]"
              >
                <span>⚡</span> Обменять
              </button>
            </div>

            <button
              v-else
              disabled
              class="w-full py-3 bg-slate-950 border border-white/5 text-slate-600 rounded-xl font-semibold text-sm cursor-default"
            >
              Доступно студентам
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- PURCHASES HISTORY TAB (STUDENT ONLY) -->
    <div v-if="activeTab === 'purchases' && userStore.role === 'student'">
      <div v-if="historyLoading" class="space-y-4">
        <div v-for="i in 3" :key="i" class="h-20 bg-slate-900 border border-white/5 rounded-xl animate-pulse"></div>
      </div>

      <div v-else-if="purchases.length === 0" class="text-center py-20 bg-slate-900/20 rounded-2xl border border-white/5">
        <div class="text-5xl mb-4">🎁</div>
        <h3 class="text-lg text-white font-medium mb-1">История покупок пуста</h3>
        <p class="text-xs text-slate-500 max-w-md mx-auto">
          Вы еще ничего не приобретали в магазине. Зарабатывайте монеты за посещение занятий и выбирайте товары на витрине!
        </p>
      </div>

      <div v-else class="space-y-4">
        <div 
          v-for="purchase in purchases" 
          :key="purchase.id"
          class="bg-slate-900/60 border border-white/5 hover:border-white/10 rounded-2xl p-5 flex flex-col md:flex-row md:items-center justify-between gap-4 transition-all"
        >
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 bg-slate-950 rounded-xl flex items-center justify-center border border-white/5 shadow-inner select-none overflow-hidden flex-shrink-0">
              <img 
                v-if="isImageFile(purchase.item_icon)" 
                :src="purchase.item_icon" 
                class="w-full h-full object-cover" 
                alt="Product Icon"
              />
              <span v-else class="text-3xl">{{ purchase.item_icon }}</span>
            </div>
            <div>
              <h4 class="text-white font-bold text-sm">{{ purchase.item_name }}</h4>
              <p class="text-[10px] text-slate-500 mt-0.5">
                Дата покупки: {{ formatDate(purchase.created_at) }}
              </p>
            </div>
          </div>

          <!-- Promo Code Card -->
          <div class="flex flex-col md:flex-row items-start md:items-center gap-4">
            <div class="bg-slate-950 border border-white/5 rounded-xl px-4 py-2.5 flex items-center gap-3">
              <div>
                <div class="text-[8px] text-slate-500 uppercase tracking-widest">Код для получения</div>
                <div class="text-xs font-mono font-bold text-slate-300 select-all tracking-wider">
                  {{ purchase.promo_code }}
                </div>
              </div>
              <button 
                @click="copyCode(purchase.promo_code)"
                class="p-1.5 bg-slate-900 hover:bg-slate-800 text-slate-400 hover:text-white rounded-lg transition-colors text-[10px]"
                title="Копировать промокод"
              >
                📋
              </button>
            </div>

            <!-- Status Tag -->
            <div class="text-right">
              <span 
                class="inline-block text-[10px] font-bold px-3 py-1 rounded-full border"
                :class="[
                  purchase.status === 'active'
                    ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/10'
                    : 'bg-slate-800 text-slate-500 border-transparent'
                ]"
              >
                {{ purchase.status === 'active' ? 'Активен (Покажите в деканате)' : 'Получен' }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ADMIN DASHBOARD TAB -->
    <div v-if="activeTab === 'admin_dashboard' && userStore.role === 'admin'">
      <div class="mb-6 flex flex-col md:flex-row gap-4 items-center justify-between">
        <!-- Search bar -->
        <div class="relative w-full md:w-96">
          <span class="absolute left-3.5 top-3 text-slate-500">🔎</span>
          <input 
            v-model="adminSearch"
            type="text"
            placeholder="Поиск по промокоду, имени или товару..."
            class="w-full pl-10 pr-4 py-2.5 bg-slate-950 border border-white/5 rounded-xl text-slate-200 text-xs focus:outline-none focus:border-red-500/30 transition-all shadow-inner"
          />
        </div>

        <div class="text-xs text-slate-500 flex items-center gap-2">
          <span>📊 Всего заказов: {{ adminPurchases.length }}</span>
          <span class="w-1.5 h-1.5 rounded-full bg-slate-700"></span>
          <span class="text-emerald-400">Активных: {{ activePurchasesCount }}</span>
        </div>
      </div>

      <div v-if="adminLoading" class="space-y-4">
        <div v-for="i in 4" :key="i" class="h-24 bg-slate-900 border border-white/5 rounded-xl animate-pulse"></div>
      </div>

      <div v-else-if="filteredAdminPurchases.length === 0" class="text-center py-20 bg-slate-900/20 rounded-2xl border border-white/5">
        <div class="text-5xl mb-4">🔎</div>
        <h3 class="text-lg text-white font-medium mb-1">Заказы не найдены</h3>
        <p class="text-xs text-slate-500 max-w-md mx-auto">
          По вашему запросу не найдено ни одной покупки. Проверьте правильность промокода или имени студента.
        </p>
      </div>

      <!-- Purchase Table/Grid -->
      <div v-else class="space-y-4">
        <div 
          v-for="p in filteredAdminPurchases" 
          :key="p.id"
          class="bg-slate-900/60 border border-white/5 hover:border-white/10 rounded-2xl p-5 flex flex-col md:flex-row md:items-center justify-between gap-4 transition-all"
        >
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 bg-slate-950 rounded-xl flex items-center justify-center border border-white/5 shadow-inner select-none overflow-hidden flex-shrink-0">
              <img 
                v-if="isImageFile(p.item_icon)" 
                :src="p.item_icon" 
                class="w-full h-full object-cover" 
                alt="Product Icon"
              />
              <span v-else class="text-3xl">{{ p.item_icon }}</span>
            </div>
            <div>
              <div class="flex items-center gap-2">
                <h4 class="text-white font-bold text-sm">{{ p.item_name }}</h4>
                <span class="text-[9px] px-2 py-0.5 bg-slate-950 border border-white/5 rounded-full text-slate-400 font-mono tracking-wider">
                  {{ p.promo_code }}
                </span>
              </div>
              <p class="text-xs text-slate-400 mt-1">
                Покупатель: <strong class="text-blue-400">{{ p.student_name }}</strong> 
                <span class="text-slate-500 ml-1">({{ p.student_group || 'Без группы' }})</span>
              </p>
              <p class="text-[10px] text-slate-500 mt-0.5">
                Дата покупки: {{ formatDate(p.created_at) }}
              </p>
            </div>
          </div>

          <div class="flex items-center gap-4 justify-between md:justify-end">
            <div class="text-right hidden md:block mr-4">
              <div class="text-[10px] text-slate-500">Списано:</div>
              <div class="flex items-center gap-1.5 text-xs font-semibold mt-0.5 justify-end">
                <span class="text-yellow-500">{{ p.coins_spent }} 🪙</span>
                <span v-if="p.stars_spent > 0" class="text-amber-500">{{ p.stars_spent }} 🌟</span>
              </div>
            </div>

            <!-- Claim Button / Status -->
            <div>
              <button 
                v-if="p.status === 'active'"
                @click="openClaimModal(p)"
                class="px-4 py-2.5 bg-emerald-600 hover:bg-emerald-500 text-white rounded-xl text-xs font-semibold shadow-lg shadow-emerald-600/10 hover:shadow-emerald-500/20 transition-all flex items-center gap-1.5"
              >
                <span>📦</span> Выдать товар
              </button>
              <div v-else class="text-right">
                <span class="inline-block text-[10px] font-bold px-3 py-1 bg-slate-800 text-slate-500 border border-transparent rounded-full select-none">
                  Выдан (Получено)
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Floating Cart Button -->
    <div 
      v-if="userStore.role === 'student'" 
      class="fixed bottom-6 right-6 z-40"
    >
      <button 
        @click="showCartModal = true"
        class="w-14 h-14 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-500 hover:to-indigo-500 text-white rounded-full flex items-center justify-center text-2xl shadow-2xl hover:scale-105 active:scale-95 transition-all relative border border-white/10 group"
      >
        <span>🛒</span>
        <!-- Badge -->
        <span 
          v-if="cartItemCount > 0"
          class="absolute -top-1 -right-1 bg-red-500 text-white text-[10px] font-bold w-5 h-5 rounded-full flex items-center justify-center border-2 border-slate-900 animate-pulse"
        >
          {{ cartItemCount }}
        </span>
      </button>
    </div>

    <!-- SHOPPING CART MODAL / DRAWER -->
    <div v-if="showCartModal" class="fixed inset-0 z-50 flex items-center justify-end">
      <!-- Backdrop -->
      <div @click="showCartModal = false" class="absolute inset-0 bg-black/60 backdrop-blur-md"></div>
      
      <!-- Drawer Content -->
      <div class="bg-slate-900 border-l border-white/10 w-full max-w-md h-full z-10 relative flex flex-col justify-between shadow-2xl">
        <!-- Drawer Header -->
        <div class="p-6 border-b border-white/5 flex items-center justify-between">
          <div class="flex items-center gap-2">
            <span class="text-2xl">🛒</span>
            <h3 class="text-lg font-bold text-white">Корзина товаров</h3>
            <span class="text-xs px-2 py-0.5 bg-slate-950 border border-white/5 rounded-full text-slate-400 font-mono">
              {{ cartItemCount }}
            </span>
          </div>
          <button 
            @click="showCartModal = false"
            class="text-slate-400 hover:text-white text-lg p-1"
          >
            ❌
          </button>
        </div>

        <!-- Drawer Body -->
        <div class="flex-1 overflow-y-auto p-6 space-y-4">
          <div v-if="cart.length === 0" class="text-center py-20">
            <div class="text-5xl mb-4">🛒</div>
            <h4 class="text-white font-medium mb-1">Корзина пуста</h4>
            <p class="text-xs text-slate-500 max-w-xs mx-auto">
              Складывайте сюда понравившиеся товары, чтобы рассчитать стоимость или обменять их все разом!
            </p>
          </div>

          <div v-else class="space-y-3">
            <div 
              v-for="cartItem in cart" 
              :key="cartItem.item.id"
              class="bg-slate-950 border border-white/5 rounded-xl p-4 flex items-center justify-between gap-4"
            >
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 bg-slate-950 rounded-lg flex items-center justify-center border border-white/5 shadow-inner select-none overflow-hidden flex-shrink-0">
                  <img 
                    v-if="isImageFile(cartItem.item.image)" 
                    :src="cartItem.item.image" 
                    class="w-full h-full object-cover" 
                    alt="Product Icon"
                  />
                  <span v-else class="text-2xl">{{ cartItem.item.image }}</span>
                </div>
                <div>
                  <h5 class="text-white font-bold text-xs">{{ cartItem.item.name }}</h5>
                  <div class="flex items-center gap-2 text-[10px] text-slate-500 mt-1">
                    <span>🪙 {{ cartItem.item.coins_price }}</span>
                    <span v-if="cartItem.item.stars_price > 0">🌟 {{ cartItem.item.stars_price }}</span>
                  </div>
                </div>
              </div>

              <!-- Quantity selector -->
              <div class="flex items-center gap-2">
                <button 
                  @click="decreaseQuantity(cartItem.item.id)"
                  class="w-6 h-6 bg-slate-900 border border-white/5 rounded flex items-center justify-center text-xs text-slate-400 hover:text-white"
                >
                  -
                </button>
                <span class="text-xs font-bold text-white px-1">{{ cartItem.quantity }}</span>
                <button 
                  @click="increaseQuantity(cartItem.item.id)"
                  class="w-6 h-6 bg-slate-900 border border-white/5 rounded flex items-center justify-center text-xs text-slate-400 hover:text-white"
                >
                  +
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Drawer Footer -->
        <div v-if="cart.length > 0" class="p-6 border-t border-white/5 bg-slate-950 space-y-4">
          <!-- Total Cost -->
          <div class="space-y-2">
            <div class="flex justify-between items-center text-xs text-slate-400">
              <span>Итого стоимость:</span>
              <div class="flex items-center gap-2 text-white font-bold">
                <span class="text-yellow-400">{{ cartTotalCoins }} 🪙</span>
                <span v-if="cartTotalStars > 0" class="text-amber-500">{{ cartTotalStars }} 🌟</span>
              </div>
            </div>

            <!-- Current Balance vs Required -->
            <div class="flex justify-between items-center text-xs text-slate-500">
              <span>Ваш баланс:</span>
              <div class="flex items-center gap-2 font-semibold">
                <span class="text-yellow-500/80">{{ userStore.profile?.coins || 0 }} 🪙</span>
                <span class="text-amber-500/80">{{ userStore.profile?.stars || 0 }} 🌟</span>
              </div>
            </div>
          </div>

          <!-- Progress / Missing amount card -->
          <div 
            v-if="!canAffordCart"
            class="p-4 bg-red-950/20 border border-red-500/20 rounded-xl space-y-2"
          >
            <div class="text-[10px] text-red-400 font-bold uppercase tracking-wider flex items-center gap-1.5">
              <span>⚠️</span> Недостаточно средств
            </div>
            <p class="text-[10px] text-slate-400 leading-relaxed">
              Для оформления заказа вам не хватает:
              <strong v-if="missingCoins > 0" class="text-yellow-400 ml-1">{{ missingCoins }} 🪙</strong>
              <strong v-if="missingStars > 0" class="text-amber-500 ml-1">{{ missingStars }} 🌟</strong>.
            </p>
            <div class="h-1.5 w-full bg-slate-900 rounded-full overflow-hidden mt-1">
              <div 
                class="h-full bg-gradient-to-r from-red-500 to-amber-500 transition-all duration-500"
                :style="{ width: `${affordProgress}%` }"
              ></div>
            </div>
            <p class="text-[9px] text-slate-500 italic leading-snug mt-1">
              💡 За каждое посещение занятия вы получаете +5 монет. За каждые 5 занятий без пропусков подряд — +1 звезда!
            </p>
          </div>

          <!-- Success Affirmation -->
          <div 
            v-else
            class="p-3 bg-emerald-950/20 border border-emerald-500/20 rounded-xl text-[10px] text-emerald-400 leading-relaxed"
          >
            🎉 Отлично! Вашего баланса активности полностью хватает для обмена всех товаров в корзине.
          </div>

          <!-- Checkout Button -->
          <button
            @click="executeCartCheckout"
            :disabled="!canAffordCart || purchaseLoading"
            class="w-full py-3.5 rounded-xl font-bold text-xs transition-all flex items-center justify-center gap-2"
            :class="[
              !canAffordCart
                ? 'bg-slate-900 border border-white/5 text-slate-500 cursor-not-allowed'
                : 'bg-emerald-600 hover:bg-emerald-500 text-white shadow-lg shadow-emerald-600/10 hover:shadow-emerald-500/20'
            ]"
          >
            <span v-if="purchaseLoading" class="inline-block animate-spin mr-1">⟳</span>
            <span>{{ canAffordCart ? 'Оформить заказ и получить коды' : 'Недостаточно средств для заказа' }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- PURCHASE CONFIRMATION MODAL -->
    <div v-if="showConfirmModal && selectedItem" class="fixed inset-0 z-50 flex items-center justify-center px-4">
      <!-- Backdrop -->
      <div @click="showConfirmModal = false" class="absolute inset-0 bg-black/60 backdrop-blur-md"></div>
      
      <!-- Modal Content -->
      <div class="bg-slate-900 border border-white/10 rounded-2xl max-w-md w-full p-6 z-10 relative shadow-2xl">
        <h3 class="text-xl font-bold text-white mb-2 flex items-center gap-2">
          <span>🛒</span> Подтверждение покупки
        </h3>
        <p class="text-xs text-slate-400 mb-6">
          Вы собираетесь обменять накопленные баллы активности на следующий товар:
        </p>

        <!-- Product Preview Card -->
        <div class="bg-slate-950 border border-white/5 rounded-xl p-4 flex items-center gap-4 mb-6">
          <div class="w-16 h-16 bg-slate-900 rounded-lg flex items-center justify-center border border-white/5 shadow-inner select-none overflow-hidden flex-shrink-0">
            <img 
              v-if="isImageFile(selectedItem.image)" 
              :src="selectedItem.image" 
              class="w-full h-full object-cover" 
              alt="Product Preview"
            />
            <span v-else class="text-4xl">{{ selectedItem.image }}</span>
          </div>
          <div>
            <h4 class="text-white font-bold text-sm">{{ selectedItem.name }}</h4>
            <p class="text-[10px] text-slate-500 mt-0.5 leading-relaxed">{{ selectedItem.description }}</p>
          </div>
        </div>

        <div class="flex items-center justify-between py-3 border-y border-white/5 mb-6 text-xs text-slate-400">
          <span>С вашего баланса спишется:</span>
          <div class="flex items-center gap-3">
            <div class="flex items-center gap-1 text-white">
              <span>🪙</span> <span class="font-bold text-yellow-400">{{ selectedItem.coins_price }}</span>
            </div>
            <div v-if="selectedItem.stars_price > 0" class="flex items-center gap-1 text-white">
              <span>🌟</span> <span class="font-bold text-amber-500">{{ selectedItem.stars_price }}</span>
            </div>
          </div>
        </div>

        <div class="flex gap-4">
          <button 
            @click="showConfirmModal = false"
            class="flex-1 py-3 bg-slate-950 hover:bg-slate-800 text-slate-400 rounded-xl text-xs font-semibold border border-white/5 transition-colors"
          >
            Отмена
          </button>
          <button 
            @click="executePurchase"
            :disabled="purchaseLoading"
            class="flex-1 py-3 bg-blue-600 hover:bg-blue-500 text-white rounded-xl text-xs font-semibold shadow-lg shadow-blue-600/20 transition-all flex items-center justify-center gap-2"
          >
            <span v-if="purchaseLoading" class="inline-block animate-spin mr-1">⟳</span>
            <span>Да, обменять!</span>
          </button>
        </div>
      </div>
    </div>

    <!-- ADMIN CONFIRM CLAIM MODAL -->
    <div v-if="showClaimModal && selectedPurchase" class="fixed inset-0 z-50 flex items-center justify-center px-4">
      <div @click="showClaimModal = false" class="absolute inset-0 bg-black/60 backdrop-blur-md"></div>
      
      <div class="bg-slate-900 border border-white/10 rounded-2xl max-w-md w-full p-6 z-10 relative shadow-2xl">
        <h3 class="text-xl font-bold text-white mb-2 flex items-center gap-2">
          <span>📦</span> Выдача товара студенту
        </h3>
        <p class="text-xs text-slate-400 mb-6">
          Подтвердите выдачу товара студенту. Это активирует промокод и спишет его из списка активных:
        </p>

        <!-- Product Preview Card -->
        <div class="bg-slate-950 border border-white/5 rounded-xl p-4 flex items-center gap-4 mb-6">
          <div class="w-16 h-16 bg-slate-900 rounded-lg flex items-center justify-center border border-white/5 shadow-inner select-none overflow-hidden flex-shrink-0">
            <img 
              v-if="isImageFile(selectedPurchase.item_icon)" 
              :src="selectedPurchase.item_icon" 
              class="w-full h-full object-cover" 
              alt="Product Preview"
            />
            <span v-else class="text-4xl">{{ selectedPurchase.item_icon }}</span>
          </div>
          <div>
            <h4 class="text-white font-bold text-sm">{{ selectedPurchase.item_name }}</h4>
            <p class="text-xs text-slate-400 mt-1">
              Получатель: <strong class="text-white">{{ selectedPurchase.student_name }}</strong>
            </p>
            <p class="text-[10px] text-slate-500 font-mono mt-0.5">Код: {{ selectedPurchase.promo_code }}</p>
          </div>
        </div>

        <div class="flex gap-4">
          <button 
            @click="showClaimModal = false"
            class="flex-1 py-3 bg-slate-950 hover:bg-slate-800 text-slate-400 rounded-xl text-xs font-semibold border border-white/5 transition-colors"
          >
            Отмена
          </button>
          <button 
            @click="executeClaim"
            :disabled="purchaseLoading"
            class="flex-1 py-3 bg-emerald-600 hover:bg-emerald-500 text-white rounded-xl text-xs font-semibold shadow-lg shadow-emerald-600/20 transition-all flex items-center justify-center gap-2"
          >
            <span v-if="purchaseLoading" class="inline-block animate-spin mr-1">⟳</span>
            <span>Выдать товар</span>
          </button>
        </div>
      </div>
    </div>

    <!-- PURCHASE SUCCESS MODAL -->
    <div v-if="showSuccessModal && successData" class="fixed inset-0 z-50 flex items-center justify-center px-4">
      <div @click="showSuccessModal = false" class="absolute inset-0 bg-black/70 backdrop-blur-md"></div>
      
      <div class="bg-slate-900 border border-white/10 rounded-2xl max-w-md w-full p-6 z-10 text-center relative shadow-2xl overflow-hidden">
        <!-- Shine background effect -->
        <div class="absolute -right-24 -top-24 w-48 h-48 bg-emerald-500/10 rounded-full blur-3xl"></div>

        <div class="w-16 h-16 bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 text-4xl rounded-full flex items-center justify-center mx-auto mb-4 animate-bounce">
          🎉
        </div>

        <h3 class="text-xl font-bold text-white mb-1">Поздравляем с покупкой!</h3>
        <p class="text-xs text-slate-400 mb-6">
          Вы успешно обменяли свои баллы на: <strong class="text-white">{{ successData.item_name }}</strong> 
          <span v-if="!isImageFile(successData.item_icon)">{{ successData.item_icon }}</span>
        </p>

        <!-- Promo Code Card -->
        <div class="bg-slate-950 border border-white/5 rounded-2xl p-5 mb-6 relative">
          <div class="text-[9px] text-slate-500 uppercase tracking-widest mb-1.5">Покажите этот код в деканате</div>
          <div class="text-2xl font-mono font-bold text-emerald-400 select-all tracking-widest">
            {{ successData.promo_code }}
          </div>
          <button 
            @click="copyCode(successData.promo_code)"
            class="absolute right-4 bottom-4 p-1.5 bg-slate-900 hover:bg-slate-800 text-slate-400 hover:text-white rounded-lg transition-colors text-xs"
            title="Копировать код"
          >
            📋
          </button>
        </div>

        <p class="text-[11px] text-slate-500 leading-relaxed mb-6">
          Код сохранен в разделе «Мои покупки» и в вашем Центре уведомлений. Вы можете забрать ваш товар в деканате колледжа в любое рабочее время.
        </p>

        <button 
          @click="showSuccessModal = false"
          class="w-full py-3 bg-emerald-600 hover:bg-emerald-500 text-white rounded-xl text-xs font-semibold shadow-lg shadow-emerald-600/20 transition-all"
        >
          Отлично!
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useUserStore } from '~/stores/user'

definePageMeta({
  middleware: 'auth'
})

const userStore = useUserStore()
const { fetchApi: api } = useApi()

const activeTab = ref('shop')
const loading = ref(false)
const historyLoading = ref(false)
const purchaseLoading = ref(false)

const items = ref([])
const purchases = ref([])

const showConfirmModal = ref(false)
const showSuccessModal = ref(false)
const selectedItem = ref(null)
const successData = ref(null)

// Admin variables
const adminPurchases = ref([])
const adminLoading = ref(false)
const adminSearch = ref('')
const showClaimModal = ref(false)
const selectedPurchase = ref(null)

// Cart state variables
const cart = ref([])
const showCartModal = ref(false)

// Utility: check if string is a realistic image file path
const isImageFile = (val) => {
  return typeof val === 'string' && (val.includes('.') || val.startsWith('/'))
}

// Affordability helpers
const canAffordCoins = (item) => {
  return (userStore.profile?.coins || 0) >= item.coins_price
}

const canAffordStars = (item) => {
  return (userStore.profile?.stars || 0) >= item.stars_price
}

const canAfford = (item) => {
  return canAffordCoins(item) && canAffordStars(item)
}

const activePurchasesCount = computed(() => {
  return adminPurchases.value.filter(p => p.status === 'active').length
})

const filteredAdminPurchases = computed(() => {
  const s = adminSearch.value.trim().toLowerCase()
  if (!s) return adminPurchases.value
  return adminPurchases.value.filter(p => {
    return (p.promo_code?.toLowerCase().includes(s)) ||
           (p.student_name?.toLowerCase().includes(s)) ||
           (p.item_name?.toLowerCase().includes(s)) ||
           (p.student_group?.toLowerCase().includes(s))
  })
})

const loadStoreItems = async () => {
  loading.value = true
  try {
    const data = await api('/store/items')
    items.value = data
  } catch (e) {
    console.error('Failed to load store items:', e)
  } finally {
    loading.value = false
  }
}

const loadPurchaseHistory = async () => {
  if (userStore.role !== 'student') return
  historyLoading.value = true
  try {
    const data = await api('/store/history')
    purchases.value = data
  } catch (e) {
    console.error('Failed to load purchase history:', e)
  } finally {
    historyLoading.value = false
  }
}

const loadAdminPurchases = async () => {
  if (userStore.role !== 'admin') return
  adminLoading.value = true
  try {
    const data = await api('/admin/store/purchases')
    adminPurchases.value = data
  } catch (e) {
    console.error('Failed to load admin purchases:', e)
  } finally {
    adminLoading.value = false
  }
}

const refreshProfile = async () => {
  try {
    const profile = await api('/profile')
    userStore.setProfile(profile)
  } catch (e) {
    console.error('Failed to refresh user profile:', e)
  }
}

const confirmPurchase = (item) => {
  selectedItem.value = item
  showConfirmModal.value = true
}

const executePurchase = async () => {
  if (!selectedItem.value) return
  purchaseLoading.value = true
  try {
    const data = await api('/store/purchase', {
      method: 'POST',
      body: { item_id: selectedItem.value.id }
    })
    
    // Close confirmation and open success
    showConfirmModal.value = false
    successData.value = data
    showSuccessModal.value = true

    // Sync student balance and reload lists
    await refreshProfile()
    await loadStoreItems()
    await loadPurchaseHistory()
  } catch (e) {
    alert(e.message || 'Ошибка оформления покупки')
  } finally {
    purchaseLoading.value = false
  }
}

const openClaimModal = (purchase) => {
  selectedPurchase.value = purchase
  showClaimModal.value = true
}

const executeClaim = async () => {
  if (!selectedPurchase.value) return
  purchaseLoading.value = true
  try {
    await api('/admin/store/claim', {
      method: 'POST',
      body: { purchase_id: selectedPurchase.value.id }
    })
    
    showClaimModal.value = false
    alert('Товар успешно выдан студенту! Студент получил системное уведомление.')
    
    // Reload admin purchases list
    await loadAdminPurchases()
  } catch (e) {
    alert(e.message || 'Ошибка выдачи товара')
  } finally {
    purchaseLoading.value = false
  }
}

const copyCode = (code) => {
  navigator.clipboard.writeText(code)
  alert('Промокод скопирован в буфер обмена!')
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleString('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// Shopping Cart Actions
const addToCart = (item) => {
  const existing = cart.value.find(c => c.item.id === item.id)
  if (existing) {
    existing.quantity++
  } else {
    cart.value.push({ item, quantity: 1 })
  }
  alert(`Товар «${item.name}» успешно добавлен в корзину!`)
}

const increaseQuantity = (itemId) => {
  const existing = cart.value.find(c => c.item.id === itemId)
  if (existing) {
    existing.quantity++
  }
}

const decreaseQuantity = (itemId) => {
  const existing = cart.value.find(c => c.item.id === itemId)
  if (existing) {
    if (existing.quantity > 1) {
      existing.quantity--
    } else {
      cart.value = cart.value.filter(c => c.item.id !== itemId)
    }
  }
}

const cartItemCount = computed(() => {
  return cart.value.reduce((acc, c) => acc + c.quantity, 0)
})

const cartTotalCoins = computed(() => {
  return cart.value.reduce((acc, c) => acc + (c.item.coins_price * c.quantity), 0)
})

const cartTotalStars = computed(() => {
  return cart.value.reduce((acc, c) => acc + (c.item.stars_price * c.quantity), 0)
})

const missingCoins = computed(() => {
  const diff = cartTotalCoins.value - (userStore.profile?.coins || 0)
  return diff > 0 ? diff : 0
})

const missingStars = computed(() => {
  const diff = cartTotalStars.value - (userStore.profile?.stars || 0)
  return diff > 0 ? diff : 0
})

const canAffordCart = computed(() => {
  return missingCoins.value === 0 && missingStars.value === 0
})

const affordProgress = computed(() => {
  const totalCoinsReq = cartTotalCoins.value || 1
  const currentCoins = userStore.profile?.coins || 0
  const coinPercent = Math.floor((currentCoins / totalCoinsReq) * 100)
  return coinPercent > 100 ? 100 : coinPercent
})

const executeCartCheckout = async () => {
  purchaseLoading.value = true
  try {
    for (const cartItem of cart.value) {
      for (let i = 0; i < cartItem.quantity; i++) {
        await api('/store/purchase', {
          method: 'POST',
          body: { item_id: cartItem.item.id }
        })
      }
    }
    
    alert('Ура! Все товары из корзины успешно обменяны! Ваши промокоды доступны в истории покупок.')
    cart.value = []
    showCartModal.value = false
    
    await refreshProfile()
    await loadStoreItems()
    if (userStore.role === 'student') {
      await loadPurchaseHistory()
    }
  } catch (e) {
    alert(e.message || 'Ошибка оформления заказа')
  } finally {
    purchaseLoading.value = false
  }
}

onMounted(() => {
  loadStoreItems()
  if (userStore.role === 'student') {
    loadPurchaseHistory()
  }
  if (userStore.role === 'admin') {
    loadAdminPurchases()
  }
})
</script>

<style scoped>
.bg-radial-glow {
  background: radial-gradient(circle, rgba(59, 130, 246, 0.08) 0%, transparent 75%);
}
</style>
