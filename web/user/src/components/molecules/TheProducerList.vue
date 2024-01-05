<script setup lang="ts">
interface Props {
  id: string | undefined
  name: string | undefined
  profile: string | undefined
  imgSrc: string | undefined
  products: undefined
}

defineProps<Props>()
</script>

<template>
  <div class="mx-auto w-[384px] rounded-3xl bg-white text-main">
    <div class="relative bottom-10">
      <img
        :src="imgSrc"
        class="mx-auto block aspect-square w-[120px] rounded-full border-2 border-white"
      />
    </div>
    <p class="relative bottom-5 text-center text-[20px] underline">
      {{ name }}
    </p>
    <p class="px-6 text-center text-[16px] tracking-[1.6px]">{{ profile }}</p>
    <div class="pt-[40px]">
      <div
        class="mx-auto flex w-[320px] justify-center rounded-3xl bg-base py-[3px] text-[16px]"
      >
        この生産者の商品
      </div>
    </div>
    <div class="grid grid-cols-2 pt-10">
      <the-coorrdinator-product-list
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
