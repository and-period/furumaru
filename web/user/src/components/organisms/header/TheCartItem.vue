<script setup lang='ts'>
interface Props {
  cartNumber: number
  marcheName: string
  boxType: string
  boxSize: number
  items: any[]
}

const props = defineProps<Props>()

const boxSizeIs60 = computed<boolean>(() => {
  return props.boxSize === 60
})

const boxSizeIs80 = computed<boolean>(() => {
  return props.boxSize === 80
})

const boxSizeIs100 = computed<boolean>(() => {
  return props.boxSize === 100
})

const useRate = computed<number>(() => {
  switch (props.boxSize) {
    case 60:
      return 30
    case 80:
      return 70
    case 100:
      return 95
    default:
      return 0
  }
})
</script>

<template>
  <div class="bg-base p-4">
    <p class="mt-2 mb-6 text-center">
      買い物カゴ #{{ cartNumber }}
    </p>

    <dl class="text-sm flex flex-col gap-y-1">
      <div class="flex">
        <dt>マルシュ:</dt>
        <dd>{{ marcheName }}</dd>
      </div>
      <div class="flex">
        <dt>箱タイプ:</dt>
        <dd>{{ boxType }}</dd>
      </div>
      <div class="flex">
        <dt>箱サイズ{{ boxSize }}:</dt>
        <dd>{{ useRate }}%使用</dd>
      </div>
    </dl>

    <div class="mt-4">
      <div class="flex items-center gap-x-2">
        <the-mandarin-orange-icon v-if="boxSizeIs60" />
        <the-apple-icon v-if="boxSizeIs80" />
        <the-melon-icon v-if="boxSizeIs100" />
        <div
          :class="{
            'bg-white w-full rounded-full h-4 border-2': true,
            'border-orange': boxSizeIs60,
            'border-apple-red': boxSizeIs80,
            'border-green': boxSizeIs100
          }"
        >
          <div
            :class="{
              'rounded-l h-3': true,
              'bg-orange border-orange border': boxSizeIs60,
              'bg-apple-red border-apple-red border': boxSizeIs80,
              'bg-green border-green border': boxSizeIs100
            }"
            :style="`width: ${useRate}%`"
          />
        </div>
      </div>

      <hr class="my-2 border-main border-dashed">

      <div
        v-for="item in items"
        :key="item.id"
      >
        <the-cart-product-item
          :name="item.name"
          :price="item.price"
          :img-src="item.imgSrc"
        />
        <hr class="my-2 border-main border-dashed">
      </div>

      <button class="py-1 bg-main text-white w-full">
        ログインして購入
      </button>
    </div>
  </div>
</template>
