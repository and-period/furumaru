<script setup lang="ts">
import { FmProductItem } from '@furumaru/shared';

import liff from '@line/liff';
import { storeToRefs } from 'pinia';
import { useProductStore } from '~/stores/product';

// Import runtime config for env variables
const runtimeConfig = useRuntimeConfig();
const liffId = runtimeConfig.public.LIFF_ID;

// Init LIFF when DOM is mounted
// https://vuejs.org/api/composition-api-lifecycle.html#onmounted
onMounted(async () => {
  if (!liffId) {
    console.error('Please set LIFF_ID in .env file');
    return;
  };

  await liff.init({ liffId: liffId });
  console.log('LIFF init success');
  console.log('LIFF SDK version', liff.getVersion());
});

const productStore = useProductStore();
const { products } = storeToRefs(productStore);
</script>

<template>
  <div>
    <h2 class="mt-6 font-semibold font-inter text-center w-full">
      商品一覧
    </h2>
    <div
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
          />
        </template>
      </div>
    </div>
  </div>
</template>
