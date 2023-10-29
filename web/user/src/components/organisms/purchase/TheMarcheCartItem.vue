<script setup lang="ts">
import { PurchaseInnerItemMock } from '~/constants/mock'

interface Props {
  marche: string
  address: string
  sender: string
  cartItems: PurchaseInnerItemMock[]
}

interface Emits {
  (e: 'click:buyButton'): void
}

const props = defineProps<Props>()
const emits = defineEmits<Emits>()

const priceStringFormatter = (price: number): string => {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: 'JPY',
  }).format(price)
}

const itemPrices = computed(() => {
  return props.cartItems.map((i) => {
    const subTotal = i.items.map((j) => j.price)
    return subTotal.reduce((m, n) => m + n)
  })
})

const totalPrice = computed(() => {
  return itemPrices.value.reduce((a, b) => a + b)
})

const handleBuyButton = () => {
  emits('click:buyButton')
}
</script>

<template>
  <div class="grid grid-flow-col gap-8 text-main">
    <div class="flex flex-col">
      <div class="text-[20px] font-bold tracking-[2px]">{{ marche }}</div>

      <div class="my-9">
        <div>カゴの数：{{ cartItems.length }}</div>
        <div>発送地：{{ address }}</div>
        <div>発送元：{{ sender }}</div>
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
        <div v-for="(cartItem, i) in cartItems" :key="i">
          <div
            class="mb-7 rounded-2xl bg-base p-2 text-center font-bold tracking-[1.2px] text-main"
          >
            買い物カゴ #{{ i + 1 }} -{{ marche }}-
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
              v-for="(item, j) in cartItem.items"
              :key="j"
              class="grid grid-cols-5 items-center border-b py-2"
            >
              <div class="col-span-2 flex gap-4">
                <img :src="item.imgSrc" class="block h-16 w-16" />
                <div>
                  {{ item.name }}
                </div>
              </div>
              <div>{{ priceStringFormatter(item.price) }}</div>
              <div class="inline-flex text-[14px]">
                <select class="h-full border-[1px] border-main px-2">
                  <option value="0">1</option>
                </select>
                <button class="ml-2 text-[12px] underline">削除</button>
              </div>
              <div>{{ priceStringFormatter(item.price) }}</div>
            </div>
          </div>

          <div class="mt-12 flex gap-20">
            <the-cardboard
              :box-size="cartItem.boxSize"
              :use-rate="cartItem.useRate"
            />

            <div class="flex flex-col gap-4">
              <div class="flex flex-col gap-2 text-center">
                <div class="rounded-3xl bg-base py-[3px] text-[12px]">
                  箱タイプ
                </div>
                <div class="text-[14px]">{{ cartItem.boxType }}</div>
              </div>

              <div class="flex flex-col gap-2 text-center">
                <div class="rounded-3xl bg-base py-[3px] text-[12px]">
                  箱サイズ
                </div>
                <div class="text-[14px]">{{ cartItem.boxSize }}</div>
              </div>

              <div class="flex flex-col gap-2 text-center">
                <div class="rounded-3xl bg-base py-[3px] text-[12px]">
                  占有率
                </div>
                <div class="text-[14px]">{{ cartItem.useRate }}</div>
              </div>
            </div>

            <div class="flex flex-col gap-7">
              <div>
                <div class="flex items-center justify-between">
                  <div class="text-[14px]">小計（税込み）</div>
                  <div class="text-[20px]">
                    {{ priceStringFormatter(itemPrices[i]) }}
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
          <div class="mt-10 rounded-2xl bg-base p-6">
            <div>
              <p class="text-[16px] tracking-[1.6px]">
                この箱には
                <span class="text-orange">{{ 100 - cartItem.useRate }}％</span>
                の空きがあります。
              </p>
              <p class="mt-1 text-[12px] tracking-[1.2px]">
                残りのスペースにはこちらの商品が同梱できます。ご一緒にいかがですか？
              </p>
            </div>
            <div class="mt-6 grid grid-cols-3 gap-4">
              <the-recommend-item
                v-for="(recommendItem, n) in cartItem.recommendItems"
                :key="n"
                :img-src="recommendItem.imgSrc"
                :name="recommendItem.name"
                :price="recommendItem.price"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
