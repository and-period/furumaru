<script lang="ts" setup>
interface Props {
  id: string;
  name: string;
  inventory: number;
  price: number;
  imgSrc: string
  address: string
  cnName: string
  cnImgSrc: string
}

const props = defineProps<Props>()

const priceString = computed<string>(() => {
  return new Intl.NumberFormat('ja-JP', { style: 'currency', currency: 'JPY' }).format(props.price)
})

const thumbnail = computed(() => {
  return props.imgSrc
})

const hasStock = computed(() => {
  return props.inventory > 0
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
      <div v-if="thumbnail" class="w-full">
        <img :src="thumbnail" :alt="name" class="aspect-square w-full">
      </div>
    </div>

    <p class="mt-2">
      {{ name }}
    </p>

    <p class="after:ml-2 text-xl my-4 after:content-['(税込)'] after:text-[16px]">
      {{ priceString }}
    </p>

    <div class="flex items-center gap-2 h-8 text-sm">
      <div class="inline-flex items-center">
        <label class="block mr-2">数量</label>
        <select class="border-main border-[1px] px-2 h-full ">
          <option value="0">
            0
          </option>
        </select>
      </div>
      <button
        class="grow h-full bg-main text-white py-1 px-4 flex items-center justify-center"
      >
        <the-cart-icon id="add-cart-icon" class="h-4 w-4 mr-1" />
        カゴに入れる
      </button>
    </div>
    <div class="mt-4 flex items-center gap-x-4 text-xs">
      <div class="grow">
        <span class="whitespace-pre-wrap">
          {{ address }}
        </span>
        <hr class=" border-main border-dashed my-2">
        <span class="before:content-['CN：']">
          {{ cnName }}
        </span>
      </div>
      <div class="min-h-14 min-w-14">
        <img :src="cnImgSrc" class="h-full w-full rounded-full">
      </div>
    </div>
  </div>
</template>
