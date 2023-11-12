<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useShoppingCartStore } from '~/store/shopping'
const router = useRouter()

const shoppingCartStore = useShoppingCartStore()
const { shoppingCart } = storeToRefs(shoppingCartStore)

const handleClickBuyButton = () => {
  router.push('/v1/purchase/address')
}

useSeoMeta({
  title: '買い物カゴ',
})
</script>

<template>
  <div class="container mx-auto">
    <div class="text-center text-[20px] font-bold tracking-[2px] text-main">
      買い物カゴ
    </div>

    <div class="my-10 border border-orange bg-white px-6 py-7 text-orange">
      <div>現在のカゴの数：{{ shoppingCart.carts.length }}</div>

      <ul class="list-disc px-6">
        <li>マルシェごとのご注文手続き・お届けとなります。</li>
        <li>
          買い物カゴごとに送料がかかります。詳しくはこちらからご確認ください。
        </li>
      </ul>
    </div>

    <div class="mt-10 flex flex-col gap-y-10">
      <the-marche-cart-item
        v-for="(cartItem, i) in shoppingCart.carts"
        :key="i"
        :cart="cartItem"
        :cart-number="cartItem.number"
        :coordinator="cartItem.coordinator"
        :items="cartItem.items"
        @click:buy-button="handleClickBuyButton"
      />
    </div>
  </div>
</template>
