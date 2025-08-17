<script setup lang="ts">
const isExpand = ref<boolean>(false);

const toggleExpand = () => {
  isExpand.value = !isExpand.value;
};
</script>

<template>
  <div>
    <div class="mb-14">
      <slot />
    </div>
    <div
      class="fixed p-4 w-full bottom-0 bg-white border-t border-gray-200 shadow-sm rounded-2xl flex flex-col transition-all gap-4"
      :class="{ 'h-svh': isExpand, 'h-[56px]': !isExpand }"
    >
      <div class="text-center">
        <button @click="toggleExpand">
          カゴの中身を見る
        </button>
      </div>
      <div
        v-if="isExpand"
        class="w-full border-t border-gray-300 py-4 px-2"
      >
        コンテンツ
      </div>
    </div>
  </div>
</template>
