<script setup lang="ts">
import { MOCK_CART_ITEMS } from '~/constants/mock'

const cartItems = MOCK_CART_ITEMS
</script>

<template>
  <div class="container mx-auto">
    <div class="text-center text-[20px] font-bold tracking-[2px] text-main">
      買い物カゴ
    </div>

    <div class="my-10 border border-orange bg-white px-6 py-7 text-orange">
      <div>現在のカゴの数：{{ cartItems.length }}</div>

      <ul class="list-disc px-6">
        <li>マルシェごとのご注文手続き・お届けとなります。</li>
        <li>
          買い物カゴごとに送料がかかります。詳しくはこちらからご確認ください。
        </li>
      </ul>
    </div>

    <div class="mt-10 flex flex-col gap-y-10">
      <the-marche-cart-item
        v-for="(cartItem, i) in cartItems"
        :key="i"
        :marche="cartItem.marche"
        :box-type="cartItem.boxType"
        :box-size="cartItem.boxSize"
        :items="cartItem.items"
      />
    </div>
  </div>
</template>
