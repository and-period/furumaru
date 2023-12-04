<script setup lang="ts">
import type { ShoppingCart } from '~/types/store'

interface Props {
  cartNumber: number
  marcheName: string
  boxType: string
  boxSize: number
  useRate: number
  shoppingCart: ShoppingCart
}

const props = defineProps<Props>()

interface Emits {
  (e: 'click:buyButton'): void
  (e: 'click:removeItemButton', cartNumber: number, id: string): void
}

const emits = defineEmits<Emits>()

const boxSizeIs60 = computed<boolean>(() => {
  return props.boxSize === 60
})

const boxSizeIs80 = computed<boolean>(() => {
  return props.boxSize === 80
})

const boxSizeIs100 = computed<boolean>(() => {
  return props.boxSize === 100
})

const handleClick = () => {
  emits('click:buyButton')
}

const handleClickRemoveButton = (id: string) => {
  emits('click:removeItemButton', props.cartNumber, id)
}
</script>

<template>
  <div class="bg-base p-4">
    <p class="mb-6 mt-2 text-center">買い物カゴ #{{ cartNumber }}</p>

    <dl class="flex flex-col gap-y-1 text-sm">
      <div class="flex">
        <dt>マルシェ：</dt>
        <dd>{{ marcheName }}</dd>
      </div>
      <div class="flex">
        <dt>箱タイプ：</dt>
        <dd>{{ boxType }}</dd>
      </div>
      <div class="flex">
        <dt>箱サイズ{{ boxSize }}：</dt>
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
            'h-4 w-full rounded-full border-2 bg-white': true,
            'border-orange': boxSizeIs60,
            'border-apple-red': boxSizeIs80,
            'border-green': boxSizeIs100,
          }"
        >
          <div
            :class="{
              'h-3 rounded-l': true,
              'border border-orange bg-orange': boxSizeIs60,
              'border border-apple-red bg-apple-red': boxSizeIs80,
              'border border-green bg-green': boxSizeIs100,
            }"
            :style="`width: ${useRate}%`"
          />
        </div>
      </div>

      <hr class="my-2 border-dashed border-main" />

      <div v-for="item in shoppingCart.items" :key="item.productId">
        <the-cart-product-item
          v-if="item"
          :id="item.productId"
          :name="item.product.name"
          :price="item.product.price"
          :img-src="item.product.thumbnail.url"
          :quantity="item.quantity"
          @click:remove-button="handleClickRemoveButton"
        />
        <hr class="my-2 border-dashed border-main" />
      </div>

      <button class="w-full bg-main py-1 text-white" @click="handleClick">
        買い物カゴを見る
      </button>
    </div>
  </div>
</template>
