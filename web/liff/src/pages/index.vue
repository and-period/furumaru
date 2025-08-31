<script setup lang="ts">
import { FmProductItem } from '@furumaru/shared';

import { NuxtLink } from '#components';
import liff from '@line/liff';
import { storeToRefs } from 'pinia';
import { useProductStore } from '~/stores/product';

const route = useRoute();

const facilityId = computed<string>(() => {
  return String(route.query.facilityId || '');
});

// Import runtime config for env variables
const runtimeConfig = useRuntimeConfig();
const liffId = runtimeConfig.public.LIFF_ID;

const isLogin = ref<boolean>(false);
const accessToken = ref<string>('');

// Init LIFF when DOM is mounted
// https://vuejs.org/api/composition-api-lifecycle.html#onmounted
onMounted(async () => {
  if (!liffId) {
    console.error('Please set LIFF_ID in .env file');
    return;
  };

  try {
    await liff.init({ liffId: liffId });
    console.log('LIFF init success');
    console.log('LIFF SDK version', liff.getVersion());
  }
  catch (error) {
    console.error('LIFF init failed', error);
  }

  if (!liff.isLoggedIn()) {
    // liff.login();
  }
  else {
    isLogin.value = true;
    const liffAccessToken = liff.getAccessToken();
    if (liffAccessToken) {
      accessToken.value = liffAccessToken;
      console.log('LIFF access token:', accessToken.value);
    }
  }
});

onMounted(() => {
  // Fetch products after LIFF initialization
  productStore.fetchProducts();
});

const productStore = useProductStore();
const { products, isLoading, error } = storeToRefs(productStore);
</script>

<template>
  <div>
    <h2 class="mt-6 font-semibold font-inter text-center w-full">
      商品一覧
    </h2>
    <div class="text-center">
      {{ isLogin ? 'ログイン済み' : '未ログイン' }} /
      {{ accessToken || 'アクセストークンの取得に失敗しました' }}
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
            :link-component-props="{ to: `/items/${product.id}?facilityId=${facilityId}`, class: 'block' }"
          />
        </template>
      </div>
    </div>
  </div>
</template>
