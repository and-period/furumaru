<script setup lang="ts">
import { FmProductItem } from '@furumaru/shared';

import liff from '@line/liff';
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
</script>

<template>
  <div>
    <h2 class="mt-6 font-semibold font-inter text-center w-full">
      商品一覧
    </h2>
    <div class="grid grid-cols-2">
      <fm-product-item
        name="商品1"
        :price="3000"
        :stoke="15"
        thumbnail-url="https://assets.furumaru.and-period.co.jp/products/media/image/sUJfjmMoMvGwjht9QgdWZK.jpg"
      />
    </div>
  </div>
</template>
