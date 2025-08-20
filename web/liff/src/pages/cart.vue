<script setup lang="ts">
import { storeToRefs } from 'pinia';
import { useShoppingCartStore } from '~/stores/shopping';

// ショッピングカートストアを使用
const shoppingCartStore = useShoppingCartStore();
const { shoppingCart, cartIsEmpty, totalPrice, totalQuantity } = storeToRefs(shoppingCartStore);

// コンポーネントマウント時にカート情報を取得
onMounted(async () => {
  await shoppingCartStore.getCart();
});

// 価格のフォーマット
const formatPrice = (price: number) => {
  return price.toLocaleString('ja-JP');
};
</script>

<template>
  <div class="container mx-auto px-4 py-6">
    <h1 class="text-2xl font-bold text-center mb-6">
      買い物カゴ
    </h1>

    <!-- カートが空の場合 -->
    <div
      v-if="cartIsEmpty"
      class="text-center py-12"
    >
      <p class="text-gray-500 text-lg mb-4">
        買い物カゴは空です
      </p>
      <NuxtLink
        to="/"
        class="text-blue-600 hover:text-blue-800 underline"
      >
        商品一覧に戻る
      </NuxtLink>
    </div>

    <!-- カートに商品がある場合 -->
    <div v-else>
      <!-- 各カートごとに表示 -->
      <div
        v-for="cart in shoppingCart.carts"
        :key="cart.number"
        class="mb-8"
      >
        <!-- コーディネーター情報 -->
        <div
          v-if="cart.coordinator"
          class="mb-4 p-4 bg-gray-50 rounded-lg"
        >
          <h2 class="text-lg font-semibold">
            {{ cart.coordinator.username }}
          </h2>
          <p class="text-sm text-gray-600">
            {{ cart.coordinator.prefecture }}{{ cart.coordinator.city }}
          </p>
        </div>

        <!-- 商品リスト -->
        <div class="space-y-4">
          <div
            v-for="item in cart.items"
            :key="item.productId"
            class="border border-gray-200 rounded-lg p-4"
          >
            <div class="flex gap-4">
              <!-- サムネイル画像 -->
              <div class="flex-shrink-0">
                <img
                  v-if="item.product?.thumbnail"
                  :src="item.product.thumbnail.url"
                  :alt="item.product.name"
                  class="w-20 h-20 object-cover rounded-lg"
                >
                <div
                  v-else
                  class="w-20 h-20 bg-gray-200 rounded-lg flex items-center justify-center"
                >
                  <span class="text-gray-400 text-xs">画像なし</span>
                </div>
              </div>

              <!-- 商品情報 -->
              <div class="flex-grow">
                <h3 class="font-semibold text-lg mb-2">
                  {{ item.product?.name || '商品名不明' }}
                </h3>

                <div class="flex items-center justify-between">
                  <div class="text-gray-600">
                    <p class="text-sm">
                      単価: ¥{{ formatPrice(item.product?.price || 0) }}
                    </p>
                    <p class="text-sm">
                      数量: {{ item.quantity }}{{ item.product?.itemUnit }}
                    </p>
                  </div>

                  <div class="text-right">
                    <p class="text-lg font-bold text-blue-600">
                      ¥{{ formatPrice((item.product?.price || 0) * item.quantity) }}
                    </p>
                  </div>
                </div>

                <!-- 商品説明（省略版） -->
                <p
                  v-if="item.product?.itemDescription"
                  class="text-sm text-gray-500 mt-2"
                >
                  {{ item.product.itemDescription }}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 合計金額 -->
      <div class="border-t border-gray-300 pt-4 mt-6">
        <div class="flex justify-between items-center mb-2">
          <span class="text-lg">合計商品数:</span>
          <span class="text-lg font-semibold">{{ totalQuantity }}点</span>
        </div>
        <div class="flex justify-between items-center">
          <span class="text-xl font-bold">合計金額:</span>
          <span class="text-2xl font-bold text-blue-600">¥{{ formatPrice(totalPrice) }}</span>
        </div>
      </div>

      <!-- アクション -->
      <div class="mt-6 space-y-2">
        <button class="w-full bg-blue-600 text-white py-3 px-4 rounded-lg font-semibold hover:bg-blue-700 transition-colors">
          レジに進む
        </button>
        <NuxtLink
          to="/"
          class="block w-full text-center py-3 px-4 border border-gray-300 rounded-lg font-semibold text-gray-700 hover:bg-gray-50 transition-colors"
        >
          買い物を続ける
        </NuxtLink>
      </div>
    </div>
  </div>
</template>
