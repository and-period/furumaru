<script setup lang="ts">
import { MOCK_PURCHASE_ITEMS } from '~/constants/mock'
const router = useRouter()

const cartItems = MOCK_PURCHASE_ITEMS

const handleClickBuyButton = () => {
  router.push('/v1/purchase/address')
}
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
        :address="cartItem.address"
        :sender="cartItem.sender"
        :cart-items="cartItem.cartItems"
        @click:buy-button="handleClickBuyButton"
      />
    </div>
  </div>
</template>
