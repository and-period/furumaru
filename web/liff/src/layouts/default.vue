<script setup lang="ts">
import { storeToRefs } from 'pinia';
import liff from '@line/liff';
import { useShoppingCartStore } from '~/stores/shopping';
import { useLiffInit } from '~/composables/useLiffInit';
import { useAuthStore } from '~/stores/auth';
import { ResponseError } from '~/types/api/facility';

const route = useRoute();
const router = useRouter();

const facilityId = computed<string>(() => String(route.params.facilityId || ''));

const isExpand = ref<boolean>(false);
const isMenuOpen = ref<boolean>(false);

// 該当ページではカートUIを非表示にする
const shouldHideCart = computed(() => {
  // /:facilityId/checkin/new のみ非表示
  const path = route.path || '';
  return /^\/[^/]+\/checkin\/new$/.test(path);
});

// ストア
const shoppingCartStore = useShoppingCartStore();
const authStore = useAuthStore();
const { shoppingCart, cartIsEmpty, totalPrice, totalQuantity } = storeToRefs(shoppingCartStore);

const toggleExpand = () => {
  isExpand.value = !isExpand.value;
};

const toggleMenu = () => {
  isMenuOpen.value = !isMenuOpen.value;
};

const onMenuSelect = () => {
  // 項目選択時はメニューを閉じる
  isMenuOpen.value = false;
};

// カゴに商品があるかどうか
const hasCartItems = computed(() => !cartIsEmpty.value);

// 削除中のプロダクトID
const removingProductId = ref<string | null>(null);

const handleRemoveItem = async (productId: string) => {
  if (removingProductId.value) {
    return;
  }

  removingProductId.value = productId;
  try {
    await shoppingCartStore.removeCartItem(productId);
  }
  catch (e) {
    console.error('Failed to remove item:', e);
  }
  finally {
    removingProductId.value = null;
  }
};

const goToCheckout = async () => {
  const id = facilityId.value;
  const path = id ? `/${id}/checkout` : '/checkout';
  router.push(path);
  toggleExpand();
};

// マウント時に認証処理とカート取得を実行
const { init: initLiff } = useLiffInit();
const runtimeConfig = useRuntimeConfig();

onMounted(async () => {
  await initLiff(runtimeConfig.public.LIFF_ID);

  // ログイン済みであれば、サインインとカート取得を実行
  if (!liff.isLoggedIn()) {
    return;
  }

  const liffIDToken = liff.getIDToken();
  if (liffIDToken) {
    try {
      await authStore.signIn(liffIDToken);
      // 認証処理の後にカート情報を取得（トークン付与のため）
      await shoppingCartStore.getCart();
    }
    catch (err) {
      if (err instanceof ResponseError) {
        if (err.response.status === 404) {
          // 404 Error の場合は新規登録へリダイレクト
          const path = facilityId.value ? `/${facilityId.value}/checkin/new` : '/checkin/new';
          await router.push(path);
          return;
        }
        if (err.response.status === 401) {
          // 401 Error の場合は再ログイン
          liff.logout();
          return;
        }
      }

      console.error('Auth signIn or redirect failed:', err);
    }
  }
});

// 価格のフォーマット
const formatPrice = (price: number) => price.toLocaleString('ja-JP');
</script>

<template>
  <div>
    <!-- クリック外閉じ用オーバーレイ -->
    <div
      v-if="isMenuOpen"
      class="fixed inset-0 z-30"
      aria-hidden="true"
      @click="isMenuOpen = false"
    />

    <!-- 画面右上：メニューボタン -->
    <div class="fixed top-1 right-3 z-40">
      <button
        id="top-right-menu-button"
        class="flex items-center justify-center w-10 h-10 rounded-full bg-white/90  hover:bg-white active:scale-[0.98] transition"
        :aria-expanded="isMenuOpen ? 'true' : 'false'"
        aria-haspopup="menu"
        aria-controls="top-right-menu-panel"
        @click="toggleMenu"
      >
        <!-- 簡易ハンバーガーアイコン -->
        <span class="sr-only">メニュー</span>
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          stroke-width="1.5"
          stroke="currentColor"
          class="size-6"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M3.75 9h16.5m-16.5 6.75h16.5"
          />
        </svg>
      </button>

      <!-- メニューのドロップダウン -->
      <div
        v-if="isMenuOpen"
        id="top-right-menu-panel"
        class="absolute right-0 mt-2 w-44 bg-white border border-gray-200 rounded-xl shadow py-2 z-30 flex flex-col text-gray-800 text-sm [&>*]:py-1 [&>*]:px-2"
        role="menu"
        aria-labelledby="top-right-menu-button"
        @click.stop
      >
        <nuxt-link
          :to="`/${facilityId}`"
          @click="onMenuSelect()"
        >
          商品一覧
        </nuxt-link>
        <nuxt-link
          :to="`/${facilityId}/orders`"
          @click="onMenuSelect()"
        >
          注文一覧
        </nuxt-link>
        <nuxt-link
          :to="`/${facilityId}/mypage`"
          @click="onMenuSelect()"
        >
          マイページ
        </nuxt-link>
        <nuxt-link
          to="/privacy"
          @click="onMenuSelect()"
        >
          プライバシーポリシー
        </nuxt-link>

        <nuxt-link
          to="/legal-notice"
          @click="onMenuSelect()"
        >
          特定商取引法に基づく表記
        </nuxt-link>
      </div>
    </div>
  </div>

  <div class="mb-14">
    <div v-if="!facilityId">
      <p class="text-center m-4 border border-yellow-800 text-yellow-800 rounded-lg bg-yellow-50 p-2">
        施設が選択されていません
      </p>
    </div>

    <slot />
  </div>
  <div
    v-if="!shouldHideCart"
    class="fixed p-4 w-full bottom-0 bg-white border-t border-gray-200 shadow-sm rounded-2xl flex flex-col transition-all gap-4 z-20"
    :class="{ 'h-9/10': isExpand, 'h-[64px]': !isExpand }"
  >
    <div class="text-center">
      <button
        class="relative flex items-center justify-center w-full border rounded-xl py-1 border-orange text-orange font-semibold"
        @click="toggleExpand"
      >
        <!-- カートに商品がある場合、右上に点滅する丸を表示 -->
        <span
          v-if="hasCartItems && !isExpand"
          class="absolute top-[-11px] right-[-2px] pointer-events-none"
          aria-hidden="true"
        >
          <span class="inline-block size-3 rounded-full bg-orange  animate-ping" />
        </span>
        <span v-if="isExpand">カゴを閉じる</span>
        <span v-else>
          カゴの中身を見る
          <span v-if="!cartIsEmpty"> ({{ totalQuantity }}点)</span>
        </span>
      </button>
    </div>
    <div
      v-if="isExpand"
      class="w-full border-t border-gray-300 py-4 px-2 overflow-y-auto"
    >
      <!-- カートが空の場合 -->
      <div
        v-if="cartIsEmpty"
        class="text-center py-8"
      >
        <p class="text-gray-500 text-lg mb-4">
          買い物カゴは空です
        </p>
      </div>

      <!-- カートに商品がある場合 -->
      <div
        v-else
        class="space-y-4"
      >
        <!-- 各カートごとに表示 -->
        <div
          v-for="cart in shoppingCart.carts"
          :key="cart.number"
          class="mb-6"
        >
          <!-- コーディネーター情報 -->
          <div
            v-if="cart.coordinator"
            class="mb-4 p-3 bg-gray-50 rounded-lg"
          >
            <h2 class="text-sm font-semibold">
              {{ cart.coordinator.username }}
            </h2>
            <p class="text-xs text-gray-600">
              {{ cart.coordinator.prefecture }}{{ cart.coordinator.city }}
            </p>
          </div>

          <!-- 商品リスト -->
          <div class="space-y-3">
            <div
              v-for="item in cart.items"
              :key="item.productId"
              class="border border-gray-200 rounded-lg p-3"
            >
              <div class="flex gap-3">
                <!-- サムネイル画像 -->
                <div class="flex-shrink-0">
                  <img
                    v-if="item.product?.thumbnail"
                    :src="item.product.thumbnail.url"
                    :alt="item.product.name"
                    class="w-16 h-16 object-cover rounded-lg"
                  >
                  <div
                    v-else
                    class="w-16 h-16 bg-gray-200 rounded-lg flex items-center justify-center"
                  >
                    <span class="text-gray-400 text-xs">画像なし</span>
                  </div>
                </div>

                <!-- 商品情報 -->
                <div class="flex-grow min-w-0">
                  <h3 class="font-semibold text-sm mb-1 truncate">
                    {{ item.product?.name || '商品名不明' }}
                  </h3>

                  <div class="space-y-1">
                    <p class="text-xs text-gray-600">
                      単価: ¥{{ formatPrice(item.product?.price || 0) }}
                    </p>
                    <p class="text-xs text-gray-600">
                      数量: {{ item.quantity }}{{ item.product?.itemUnit }}
                    </p>
                    <p class="text-sm font-bold text-main">
                      ¥{{ formatPrice((item.product?.price || 0) * item.quantity) }}
                    </p>
                  </div>

                  <div class="mt-2 text-right">
                    <button
                      class="text-xs text-red-600 border border-red-300 px-2 py-1 rounded hover:bg-red-50 disabled:opacity-50"
                      :disabled="removingProductId === item.productId"
                      @click="handleRemoveItem(item.productId)"
                    >
                      {{ removingProductId === item.productId ? '削除中…' : '削除' }}
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 合計金額 -->
        <div class="border-t border-gray-300 pt-4 mt-4">
          <div class="flex justify-between items-center mb-2">
            <span class="text-sm">合計商品数:</span>
            <span class="text-sm font-semibold">{{ totalQuantity }}点</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-lg font-bold">合計金額:</span>
            <span class="text-xl font-bold text-main">¥{{ formatPrice(totalPrice) }}</span>
          </div>
        </div>

        <!-- アクション -->
        <div class="mt-4 space-y-2">
          <button
            class="w-full bg-orange text-white py-3 px-4 rounded-lg font-semibold hover:bg-orange/[0.7] transition-colors"
            @click="goToCheckout"
          >
            購入処理に進む
          </button>
          <button
            class="w-full text-center py-3 px-4 border border-gray-300 rounded-lg font-semibold text-gray-700 hover:bg-gray-50 transition-colors"
            @click="toggleExpand"
          >
            買い物を続ける
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
