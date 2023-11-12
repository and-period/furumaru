<script setup lang="ts">
import { Coordinator } from '~/types/api'
import { CartItem, ShoppingCart } from '~/types/store'

interface Props {
  cartNumber: number
  coordinator: Coordinator
  cart: ShoppingCart
  items: CartItem[]
}

interface Emits {
  (e: 'click:buyButton'): void
}

const props = defineProps<Props>()
const emits = defineEmits<Emits>()

const totalPrice = computed<number>(() => {
  return props.items
    .map((item) => item.product.price * item.quantity)
    .reduce((sum, price) => sum + price)
})

const priceStringFormatter = (price: number): string => {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: 'JPY',
  }).format(price)
}

const handleBuyButton = () => {
  emits('click:buyButton')
}
</script>

<template>
  <div class="grid grid-flow-col gap-8 text-main">
    <div class="flex flex-col">
      <div class="text-[20px] font-bold tracking-[2px]">
        {{ coordinator.marcheName }}
      </div>

      <div class="my-9">
        <div>カゴの数：{{}}</div>
        <div>発送地：{{ coordinator.username }}</div>
        <div>発送元：{{}}</div>
      </div>

      <div class="flex items-center justify-between font-bold">
        <div class="text-[14px]">商品合計（税込み）</div>
        <div class="text-[20px]">{{ priceStringFormatter(totalPrice) }}</div>
      </div>

      <hr class="my-5 border-main" />

      <div>※送料はご購入手続き画面で加算されます。</div>

      <button
        class="mt-8 bg-main p-[14px] text-[16px] text-white"
        @click="handleBuyButton"
      >
        ご購入手続きへ
      </button>
    </div>

    <div class="relative col-span-10 rounded-3xl bg-white px-16 py-12">
      <div class="absolute -left-4 top-12 h-8 w-8 rotate-45 bg-white"></div>

      <div class="flex flex-col gap-y-16">
        <div
          class="mb-7 rounded-2xl bg-base p-2 text-center font-bold tracking-[1.2px] text-main"
        >
          買い物カゴ #{{ cartNumber }} -{{ coordinator.marcheName }}-
        </div>

        <div class="flex w-full flex-col text-main">
          <div
            class="grid grid-cols-5 items-center border-b py-2 text-[12px] tracking-[1.2px]"
          >
            <div class="col-span-2">商品</div>
            <div>価格（税込み）</div>
            <div>数量</div>
            <div>小計（税込み）</div>
          </div>

          <div
            v-for="(item, j) in items"
            :key="j"
            class="grid grid-cols-5 items-center border-b py-2"
          >
            <div class="col-span-2 flex gap-4">
              <img :src="item.product.thumbnail.url" class="block h-16 w-16" />
              <div>
                {{ item.product.name }}
              </div>
            </div>
            <div>{{ priceStringFormatter(item.product.price) }}</div>
            <div class="inline-flex text-[14px]">
              {{ item.quantity }}
              <button class="ml-2 text-[12px] underline">削除</button>
            </div>
            <div>
              {{ priceStringFormatter(item.product.price * item.quantity) }}
            </div>
          </div>
        </div>

        <div class="mt-12 flex gap-20">
          <the-cardboard :box-size="cart.size" :use-rate="0" />

          <div class="flex flex-col gap-4">
            <div class="flex flex-col gap-2 text-center">
              <div class="rounded-3xl bg-base py-[3px] text-[12px]">
                箱タイプ
              </div>
              <div class="text-[14px]">{{ cart.type }}</div>
            </div>

            <div class="flex flex-col gap-2 text-center">
              <div class="rounded-3xl bg-base py-[3px] text-[12px]">
                箱サイズ
              </div>
              <div class="text-[14px]">{{ cart.size }}</div>
            </div>

            <div class="flex flex-col gap-2 text-center">
              <div class="rounded-3xl bg-base py-[3px] text-[12px]">占有率</div>
              <div class="text-[14px]">{{}}</div>
            </div>
          </div>

          <div class="flex flex-col gap-7">
            <div>
              <div class="flex items-center justify-between">
                <div class="text-[14px]">小計（税込み）</div>
                <div class="text-[20px]">
                  {{ priceStringFormatter(totalPrice) }}
                </div>
              </div>
              <hr class="mt-4 border-main" />
            </div>
            <div>※送料はご購入手続き画面で加算されます。</div>

            <button class="bg-main p-[14px] text-[16px] text-white">
              このカゴのご購入手続きへ
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
