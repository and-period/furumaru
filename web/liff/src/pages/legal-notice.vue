<script setup lang="ts">
const router = useRouter();

definePageMeta({
  title: '特定商取引法に基づく表記',
  layout: 'init',
});
</script>

<template>
  <div class="min-h-screen bg-white relative mb-[186px] overflow-auto">
    <!-- Navigation header for mobile -->
    <div class="top-0 z-10 bg-white border-b border-gray-200 px-4 py-3 fixed w-full">
      <div class="flex items-center">
        <button
          class="flex items-center text-gray-600 hover:text-gray-800"
          @click="router.back()"
        >
          <svg
            class="w-6 h-6 mr-2"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M15 19l-7-7 7-7"
            />
          </svg>
          戻る
        </button>
      </div>
    </div>

    <div class="mt-12 py-2 px-4">
      <h2 class="font-semibold font-inter text-center w-full">
        特定商取引法に基づく表記
      </h2>
      <!-- 特定商取引法に基づく表記本文 -->
      <div class="mt-2 text-gray-600">
        株式会社アンドピリオド
      </div>
    </div>
  </div>
</template>
