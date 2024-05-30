<script setup lang="ts">
interface Props {
  imgSrc: string
  name: string
  price: number
}

const props = defineProps<Props>()

const priceString = computed(() => {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: 'JPY',
  }).format(props.price)
})
</script>

<template>
  <div class="flex gap-x-[10px] bg-white p-2">
    <nuxt-img
      :src="imgSrc"
      class="h-20 w-20"
      :alt="`${name}のサムネイル`"
      provider="cloudFront"
      width="80px"
      height="80px"
      format="webp"
    />
    <div
      class="flex w-full flex-col justify-between text-[12px] tracking-[1.2px]"
    >
      <div>
        <div class="min-h-[28px]">
          {{ name }}
        </div>
        <div class="font-bold">
          {{ priceString }}（税込み）
        </div>
      </div>
      <button class="mt-2 bg-main px-2 py-1 text-white">
        追加する
      </button>
    </div>
  </div>
</template>
