<script setup lang="ts">
import { storeToRefs } from 'pinia';
import liff from '@line/liff';
import { useAuthStore } from '~/stores/auth';
import { useShoppingCartStore } from '~/stores/shopping';
import { ResponseError } from '~/types/api/facility';

const route = useRoute();
const router = useRouter();
const runtimeConfig = useRuntimeConfig();

const facilityId = computed<string>(() => String(route.params.facilityId || ''));

const isExpand = ref<boolean>(false);

// ストア
const shoppingCartStore = useShoppingCartStore();
const authStore = useAuthStore();
const { shoppingCart, cartIsEmpty, totalPrice, totalQuantity } = storeToRefs(shoppingCartStore);

const toggleExpand = () => {
  isExpand.value = !isExpand.value;
};

// マウント時に認証処理とカート取得を実行
onMounted(async () => {
  const liffId = runtimeConfig.public.LIFF_ID;
  if (!liffId) {
    console.error('Please set LIFF_ID in .env file');
  }
  else {
    try {
      console.log(liffId);
      await liff.init({ liffId });
    }
    catch (error) {
      console.error('LIFF init failed', error);
    }
  }

  if (!liff.isLoggedIn()) {
    liff.login();
    return;
  }

  // サインイン（LIFFログイン済みの場合のみ）
  if (liff.isLoggedIn()) {
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
  }
});

// 価格のフォーマット
const formatPrice = (price: number) => price.toLocaleString('ja-JP');
</script>

<template>
  <div>
    <div class="mb-14">
      <div v-if="!facilityId">
        <p class="text-center m-4 border border-yellow-800 text-yellow-800 rounded-lg bg-yellow-50 p-2">
          施設が選択されていません
        </p>
      </div>

      <slot />
    </div>
    <div
      class="fixed p-4 w-full bottom-0 bg-white border-t border-gray-200 shadow-sm rounded-2xl flex flex-col transition-all gap-4"
      :class="{ 'h-svh': isExpand, 'h-[56px]': !isExpand }"
    >
      <div class="text-center">
        <button
          class="flex items-center justify-center w-full"
          @click="toggleExpand"
        >
          <span v-if="cartIsEmpty">カゴの中身を見る</span>
          <span v-else>カゴの中身を見る ({{ totalQuantity }}点)</span>
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
            <button class="w-full bg-orange text-white py-3 px-4 rounded-lg font-semibold hover:bg-orange/[0.7] transition-colors">
              レジに進む
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
  </div>
</template>
