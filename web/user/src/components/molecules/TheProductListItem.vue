<script lang="ts" setup>
import { ProductItem } from '~/types/store'

interface Props {
  item: ProductItem
}

const props = defineProps<Props>()

const priceString = computed<string>(() => {
  return new Intl.NumberFormat('ja-JP', { style: 'currency', currency: 'JPY' }).format(props.item.price)
})

const thumbnail = computed(() => {
  return props.item.media.find(item => item.isThumbnail)
})

const hasStock = computed(() => {
  return props.item.inventory > 0
})
</script>

<template>
  <div class="text-main">
    <div class="relative">
      <div v-if="!hasStock" class="absolute inset-0 flex justify-center items-center bg-black bg-opacity-50">
        <p class="text-white text-lg font-semibold">
          在庫なし
        </p>
      </div>
      <div v-if="thumbnail">
        <picture>
          <img src="~/assets/img/sample.png" :alt="item.name">
        </picture>
      </div>
    </div>

    <p class="mt-2">
      {{ item.name }}
    </p>

    <p class="after:ml-2 after:content-['(税込み)'] text-xl my-4">
      {{ priceString }}
    </p>

    <div class="flex items-center justify-between gap-1">
      <p>数量</p>
      <div class="border-main border-[1px] pa-2">
        <select>
          <option value="0">
            0
          </option>
        </select>
      </div>
      <the-button
        class="bg-main text-white py-1 px-4 flex items-center"
      >
        <the-cart-icon id="add-cart-icon" class="h-4 w-4 mr-1" />
        カゴに入れる
      </the-button>
    </div>
  </div>
</template>
