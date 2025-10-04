<script setup lang="ts">
const route = useRoute();
const { $liffReady } = useNuxtApp();

onMounted(async () => {
  await $liffReady;
});

const loading = computed(() => 'liff.state' in route.query);

definePageMeta({
  layout: 'init',
});
</script>

<template>
  <div>
    <!-- Loading screen -->
    <div
      v-if="loading"
      class="flex h-screen items-center justify-center"
    >
      <div class="flex flex-col items-center gap-3">
        <div class="h-10 w-10 animate-spin rounded-full border-4 border-gray-200 border-t-blue-500" />
        <p class="text-gray-600">
          読み込み中...
        </p>
      </div>
    </div>

    <!-- Content when not loading -->
    <div
      v-else
      class="text-center"
    >
      {{ route.query }}
    </div>
  </div>
</template>
