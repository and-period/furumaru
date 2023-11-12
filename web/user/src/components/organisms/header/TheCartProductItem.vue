<script setup lang="ts">
interface Props {
  name: string
  price: number
  imgSrc: string
  quantity: number
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
  <div>
    <p>{{ name }}</p>

    <p
      class="my-4 font-bold after:ml-2 after:text-[16px] after:content-['(税込)']"
    >
      {{ priceString }}
    </p>

    <div class="flex gap-x-3 text-sm">
      <img :src="imgSrc" :alt="name" class="aspect-square h-[72px] w-[72px]" />
      <div class="flex flex-col">
        <div class="flex grow items-start">
          <div class="inline-flex">
            <div class="mr-2 block whitespace-nowrap">
              数量
              {{ quantity }}
            </div>
          </div>
        </div>
        <div class="text-right">
          <a href="#" class="underline">削除</a>
        </div>
      </div>
    </div>
  </div>
</template>
