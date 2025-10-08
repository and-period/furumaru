<script setup lang="ts">
import type { Product } from '~/types/api'

interface Props {
  id: string | undefined
  name: string | undefined
  profile: string | undefined
  imgSrc: string | undefined
  products: Product[] | undefined
}

defineProps<Props>()

interface Emits {
  (e: 'click:product-item', id: string): void
}

const emits = defineEmits<Emits>()

const handleClickProductItem = (productId: string) => {
  emits('click:product-item', productId)
}
</script>

<template>
  <div class="mx-auto w-full rounded-3xl bg-white text-main">
    <div class="relative bottom-10">
      <img
        :src="imgSrc ? imgSrc : '/img/furuneko.png'"
        class="mx-auto block aspect-square w-[96px] rounded-full border-4 border-white object-cover md:w-[120px] shadow-lg"
      >
    </div>
    <p
      class="relative bottom-5 text-center text-[16px] tracking-[1.6px] underline md:text-[20px] md:tracking-[2.0px]"
    >
      {{ name }}
    </p>
    <p
      class="px-6 text-[14px] tracking-[1.4px] md:text-[16px] md:tracking-[1.6px]"
    >
      {{ profile }}
    </p>
    <div class="px-4 pt-[40px]">
      <div
        class="mx-4 flex justify-center rounded-3xl bg-base py-[3px] text-[14px] font-bold md:mx-auto md:text-[16px]"
      >
        この生産者の商品
      </div>
    </div>
    <div class="px-4 pt-10 pb-6">
      <div
        v-if="!products || products.length === 0"
        class="flex flex-col items-center justify-center py-6 px-4 bg-gradient-to-r from-amber-50 to-orange-50 rounded-xl border border-amber-100"
      >
        <div class="text-amber-500 text-2xl mb-2">
          🌱
        </div>
        <div class="text-center">
          <p class="text-amber-700 text-[14px] font-medium mb-1">
            新商品を準備中
          </p>
          <p class="text-amber-600 text-[11px] opacity-80">
            お楽しみに！
          </p>
        </div>
      </div>
      <div
        v-else
        class="grid grid-cols-2 gap-2"
      >
        <the-coordinator-product-list
          v-for="product in products"
          :id="product.id"
          :key="product.id"
          :name="product.name"
          :inventory="product.inventory"
          :price="product.price"
          :thumbnail="product.media[0]"
          @click:item="handleClickProductItem"
        />
      </div>
    </div>
  </div>
</template>
