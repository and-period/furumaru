<script lang="ts" setup>
import type { ProductMediaInner } from '~/types/api'

interface Props {
  id: string
  name: string
  price: number
  inventory: number
  thumbnail: ProductMediaInner | undefined
}


const props = defineProps<Props>()

const priceString = computed<string>(() => {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: 'JPY',
  }).format(props.price)
})


</script>

<template>
  <div class="mx-auto flex flex-col text-main">
    <div class="relative w-[144px]">
      <div
        v-if="inventory === 0"
        class="absolute inset-0 flex items-center justify-center bg-black/50"
      >
        <p class="text-[14px] font-semibold text-white">在庫なし</p>
      </div>
      <picture
        v-if="thumbnail"
        class="w-full"
      >
        <img
          :src="thumbnail.url"
          :alt="`${name}のサムネイル画像`"
          class="aspect-square w-full"
        />
      </picture>
    </div>

    <p
      class="mt-2 line-clamp-3 grow text-[14px] tracking-[1.4px] md:text-[14px] md:tracking-[1.6px]"
    >
      {{ name }}
    </p>

    <p
      class="my-2 text-[16px] font-bold tracking-[1.6px] after:ml-2 after:text-[16px] after:content-['(税込)'] md:text-[16px] md:tracking-[2.0px] mb-4"
    >
      {{ priceString }}
    </p>
  </div>
</template>
