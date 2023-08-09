<script setup lang="ts">
import { CartItemMock } from '~/constants/mock'

interface Props {
  marche: string
  items: CartItemMock[]
  boxType: string
  boxSize: number
}

defineProps<Props>()

const priceStringFormatter = (price: number): string => {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: 'JPY',
  }).format(price)
}
</script>

<template>
  <div class="grid grid-flow-col gap-8 text-main">
    <div class="flex flex-col">
      <div class="text-[20px] font-bold tracking-[2px]">{{ marche }}</div>

      <div class="my-9">
        <div>カゴの数：2</div>
        <div>発送地：2</div>
        <div>発送元：2</div>
      </div>

      <div class="flex items-center justify-between font-bold">
        <div class="text-[14px]">商品合計（税込み）</div>
        <div class="text-[20px]">{{ priceStringFormatter(9000) }}</div>
      </div>

      <hr class="my-5 border-main" />

      <div>※送料はご購入手続き画面で加算されます。</div>

      <button class="mt-8 bg-main p-[14px] text-[16px] text-white">
        ご購入手続きへ
      </button>
    </div>

    <div class="relative col-span-10 rounded-3xl bg-white px-16 py-12">
      <div class="absolute -left-4 top-12 h-8 w-8 rotate-45 bg-white"></div>

      <div
        class="mb-7 rounded-2xl bg-base p-2 text-center font-bold tracking-[1.2px] text-main"
      >
        買い物カゴ #{{ 1 }} -{{ marche }}-
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
          v-for="(item, i) in items"
          :key="i"
          class="grid grid-cols-5 items-center border-b py-2"
        >
          <div class="col-span-2 flex gap-4">
            <img :src="item.imgSrc" class="block h-16 w-16" />
            <div>
              {{ item.name }}
            </div>
          </div>
          <div>{{ priceStringFormatter(item.price) }}</div>
          <div>{{ item.inventory }}</div>
          <div>{{ priceStringFormatter(item.price) }}</div>
        </div>
      </div>

      <div class="mt-12 flex gap-20">
        <the-cardboard />

        <div class="flex flex-col gap-4">
          <div class="flex flex-col gap-2 text-center">
            <div class="rounded-3xl bg-base py-[3px] text-[12px]">箱タイプ</div>
            <div class="text-[14px]">{{ boxType }}</div>
          </div>

          <div class="flex flex-col gap-2 text-center">
            <div class="rounded-3xl bg-base py-[3px] text-[12px]">箱サイズ</div>
            <div class="text-[14px]">{{ boxSize }}</div>
          </div>

          <div class="flex flex-col gap-2 text-center">
            <div class="rounded-3xl bg-base py-[3px] text-[12px]">占有率</div>
            <div class="text-[14px]">{{ boxSize }}</div>
          </div>
        </div>

        <div class="flex flex-col gap-7">
          <div>
            <div class="flex items-center justify-between">
              <div class="text-[14px]">小計（税込み）</div>
              <div class="text-[20px]">{{ priceStringFormatter(6000) }}</div>
            </div>
            <hr class="mt-4 border-main" />
          </div>
          <div>※送料はご購入手続き画面で加算されます。</div>

          <button class="bg-main p-[14px] text-[16px] text-white">
            このかごのご購入手続きへ
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
