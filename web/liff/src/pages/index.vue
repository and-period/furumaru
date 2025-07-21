<script setup lang="ts">
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
    <div />
  </div>
</template>
