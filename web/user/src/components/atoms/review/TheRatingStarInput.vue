<script setup lang="ts">
const model = defineModel<number>({ required: true })

// マウスオーバー中の一時的な評価
const hoverRating = ref(0)

// 現在表示すべき評価（マウスオーバー中ならその値、そうでなければ実際の評価）
const currentRating = computed(() => hoverRating.value || model.value)

const setRating = (rating: number) => {
  model.value = rating
}
</script>

<template>
  <div class="inline-flex items-center gap-4">
    <!-- コンテナにマウスが離れたら hoverRating をリセット -->
    <div
      class="inline-flex items-center gap-1"
      @mouseleave="hoverRating = 0"
    >
      <template
        v-for="i in 5"
        :key="i"
      >
        <svg
          :class="['w-6 h-6 ms-1', i <= currentRating ? 'text-yellow-300' : 'text-gray-300']"
          aria-hidden="true"
          xmlns="http://www.w3.org/2000/svg"
          fill="currentColor"
          viewBox="0 0 22 20"
          @mouseenter="hoverRating = i"
          @click="setRating(i)"
        >
          <path d="M20.924 7.625a1.523 1.523 0 0 0-1.238-1.044l-5.051-.734-2.259-4.577a1.534 1.534 0 0 0-2.752 0L7.365 5.847l-5.051.734A1.535 1.535 0 0 0 1.463 9.2l3.656 3.563-.863 5.031a1.532 1.532 0 0 0 2.226 1.616L11 17.033l4.518 2.375a1.534 1.534 0 0 0 2.226-1.617l-.863-5.03L20.537 9.2a1.523 1.523 0 0 0 .387-1.575Z" />
        </svg>
      </template>
    </div>
    <div class="font-semibold">
      {{ model>0 ? model :'-' }}
    </div>
  </div>
</template>
