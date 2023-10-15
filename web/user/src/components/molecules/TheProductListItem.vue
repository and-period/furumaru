<script lang="ts" setup>
interface Props {
  id: string
  name: string
  inventory: number
  price: number
  imgSrc: string
  address: string
  cnName: string
  cnImgSrc: string
}

const props = defineProps<Props>()

const priceString = computed<string>(() => {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: 'JPY',
  }).format(props.price)
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
      <div
        v-if="!hasStock"
        class="absolute inset-0 flex items-center justify-center bg-black/50"
      >
        <p class="text-lg font-semibold text-white">在庫なし</p>
      </div>
      <div v-if="thumbnail" class="w-full">
        <img :src="thumbnail" :alt="name" class="aspect-square w-full" />
      </div>
    </div>

    <p class="mt-2">
      {{ name }}
    </p>

    <p
      class="my-4 text-xl after:ml-2 after:text-[16px] after:content-['(税込)']"
    >
      {{ priceString }}
    </p>

    <div class="flex h-8 items-center gap-2 text-sm">
      <div class="inline-flex items-center">
        <label class="mr-2 block text-[10px] xl:text-[14px]">数量</label>
        <select class="h-full border-[1px] border-main px-2">
          <option value="0">0</option>
        </select>
      </div>
      <button
        class="flex h-full grow items-center justify-center bg-main p-1 text-[10px] text-white lg:px-4 lg:text-[14px]"
      >
        <the-cart-icon id="add-cart-icon" class="mr-1 h-3 w-3 lg:h-4 lg:w-4" />
        カゴに入れる
      </button>
    </div>
    <div class="mt-4 flex items-center gap-x-4 text-xs">
      <div class="grow">
        <span class="whitespace-pre-wrap">
          {{ address }}
        </span>
        <hr class="my-2 border-dashed border-main" />
        <span class="before:content-['CN：']">
          {{ cnName }}
        </span>
      </div>
      <div class="h-14 w-14">
        <img :src="cnImgSrc" class="aspect-square rounded-full" />
      </div>
    </div>
  </div>
</template>
