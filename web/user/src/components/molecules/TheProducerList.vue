<script setup lang="ts">
import type { Product } from "~/types/api"

interface Props {
  id: string | undefined
  name: string | undefined
  profile: string | undefined
  imgSrc: string | undefined
  products: Product[] | undefined
}

defineProps<Props>()
</script>

<template>
  <div class="mx-auto rounded-3xl bg-white text-main md:w-[384px]">
    <div class="relative bottom-10">
      <img
        :src="imgSrc"
        class="mx-auto block aspect-square w-[96px] rounded-full border-2 border-white md:w-[120px]"
      />
    </div>
    <p class="relative bottom-5 text-center text-[16px] tracking-[1.6px] underline md:text-[20px] md:tracking-[2.0px]">
      {{ name }}
    </p>
    <p class="px-6 text-[14px] tracking-[1.4px] md:text-[16px] md:tracking-[1.6px]">{{ profile }}</p>
    <div class="pt-[40px]">
      <div
        class="mx-4 md:mx-auto flex w-[320px] justify-center rounded-3xl bg-base py-[3px] text-[14px] font-bold md:text-[16px]"
      >
        この生産者の商品
      </div>
    </div>
    <div class="grid grid-cols-2 pt-10">
      <the-coordinator-product-list
        v-for="product in products"
        :id="product.id"
        :key="product.id"
        :name="product.name"
        :inventory="product.inventory"
        :price="product.price"
        :thumbnail="product.media[0]"
      />
    </div>
  </div>
</template>
