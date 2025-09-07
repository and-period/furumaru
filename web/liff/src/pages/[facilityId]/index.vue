<script setup lang="ts">
import { FmProductItem } from '@furumaru/shared';
import { NuxtLink } from '#components';
import { storeToRefs } from 'pinia';
import { useProductStore } from '~/stores/product';
import { useAuthStore } from '~/stores/auth';
import { useShoppingCartStore } from '~/stores/shopping';

const route = useRoute();
const facilityId = computed<string>(() => String(route.params.facilityId || ''));

// 認証状態はレイアウトで処理済み。表示用にストアから参照。
const authStore = useAuthStore();
const isLogin = computed(() => authStore.isAuthenticated);
const idToken = computed(() => authStore.token?.accessToken || '');

// 商品取得
const productStore = useProductStore();
const { products, isLoading, error } = storeToRefs(productStore);
onMounted(() => {
  productStore.fetchProducts();
});

// カート追加
const shoppingCartStore = useShoppingCartStore();
async function handleAddToCart(productId: string, quantity: number) {
  try {
    await shoppingCartStore.addCartItem(productId, quantity);
  }
  catch (e) {
    console.error('Failed to add to cart:', e);
  }
}
</script>

<template>
  <div>
    <h2 class="mt-6 font-semibold font-inter text-center w-full">
      商品一覧
    </h2>
    <div class="text-center">
      {{ isLogin ? 'ログイン済み' : '未ログイン' }} /
      {{ idToken || 'IDトークンの取得に失敗しました' }}
    </div>

    <!-- Loading state -->
    <div
      v-if="isLoading"
      class="container mx-auto mt-6 text-center"
    >
      <p>商品を読み込み中...</p>
    </div>

    <!-- Error state -->
    <div
      v-else-if="error"
      class="container mx-auto mt-6 text-center text-red-600"
    >
      <p>エラー: {{ error }}</p>
    </div>

    <!-- Products grid -->
    <div
      v-else
      class="container mx-auto mt-6"
    >
      <div class="grid lg:grid-cols-5 md:grid-cols-3 grid-cols-2 gap-4 px-4">
        <template
          v-for="product in products"
          :key="product.id"
        >
          <FmProductItem
            :name="product.name"
            :price="product.price"
            :stock="product.inventory"
            :thumbnail-url="product.thumbnailUrl"
            :link-component="NuxtLink"
            :link-component-props="{ to: `/${facilityId}/items/${product.id}`, class: 'block' }"
            @click:add-cart="(q) => handleAddToCart(product.id, q)"
          />
        </template>
      </div>
    </div>
  </div>
</template>
