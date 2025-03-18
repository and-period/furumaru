<script setup lang="ts">
interface Props {
  rate: number
  id?: string
}

withDefaults(defineProps<Props>(), {
  id: Math.random().toString(36).substring(2, 15),
})

/**
 * 1～5の星（i）に対して、どのくらい塗りつぶすかを 0～1 で返す関数
 */
const getStarFill = (rate: number, i: number): number => {
  const fill = rate - (i - 1)
  return Math.min(Math.max(fill, 0), 1) // 0～1の間にクリップ
}
</script>

<template>
  <div class="flex items-center">
    <template
      v-for="i in 5"
      :key="i"
    >
      <svg
        class="w-4 h-4 ms-1"
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 22 20"
      >
        <!-- グラデーション定義 -->
        <defs>
          <linearGradient
            :id="`${id}starGradient${i}`"
            x1="0%"
            y1="0%"
            x2="100%"
            y2="0%"
          >
            <!--
              「offset=◯%」までは黄色, そこから先は灰色
              例) offset=50% のとき、左半分は黄色、右半分は灰色
            -->
            <stop
              :offset="(getStarFill(rate, i) * 100) + '%'"
              stop-color="#FFF176"
            />
            <stop
              :offset="(getStarFill(rate, i) * 100) + '%'"
              stop-color="#d1d5dc"
            />
          </linearGradient>
        </defs>
        <!-- パス全体に上記のグラデーションを適用する -->
        <path
          :fill="`url(#${id}starGradient${i})`"
          d="M20.924 7.625a1.523 1.523 0 0 0-1.238-1.044l-5.051-.734-2.259-4.577a1.534 1.534 0 0 0-2.752 0L7.365 5.847l-5.051.734A1.535 1.535 0 0 0 1.463 9.2l3.656 3.563-.863 5.031a1.532 1.532 0 0 0 2.226 1.616L11 17.033l4.518 2.375a1.534 1.534 0 0 0 2.226-1.617l-.863-5.03L20.537 9.2a1.523 1.523 0 0 0 .387-1.575Z"
        />
      </svg>
    </template>
  </div>
</template>
