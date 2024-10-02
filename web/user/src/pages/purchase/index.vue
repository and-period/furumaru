<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useAuthStore } from '~/store/auth'
import { useShoppingCartStore } from '~/store/shopping'
import type { I18n } from '~/types/locales'

const i18n = useI18n()
const router = useRouter()

const ct = (str: keyof I18n['purchase']['cart']) => {
  return i18n.t(`purchase.cart.${str}`)
}

const shoppingCartStore = useShoppingCartStore()
const { removeProductFromCart } = shoppingCartStore
const { shoppingCart } = storeToRefs(shoppingCartStore)

const authStore = useAuthStore()
const { isAuthenticated } = storeToRefs(authStore)

const handleClickBuyButton = (coordinatorId: string) => {
  if (isAuthenticated.value) {
    router.push(`/v1/purchase/address?coordinatorId=${coordinatorId}`)
  }
  else {
    router.push(
      `/v1/purchase/auth?required=true&coordinatorId=${coordinatorId}`,
    )
  }
}

const handleClickCartBuyButton = (
  coordinatorId: string,
  cartNumber: number,
) => {
  if (isAuthenticated.value) {
    router.push({
      path: '/v1/purchase/address',
      query: {
        coordinatorId,
        cartNumber,
      },
    })
  }
  else {
    router.push({
      path: '/v1/purchase/auth',
      query: {
        required: true,
        coordinatorId,
        cartNumber,
      },
    })
  }
}

const handelClickRemoveItemFromCartButton = (
  cartNumber: number,
  id: string,
) => {
  removeProductFromCart(cartNumber, id)
}

useSeoMeta({
  title: '買い物カゴ',
})
</script>

<template>
  <div class="container mx-auto px-4 xl:px-0 my-10">
    <div class="text-center text-[20px] font-bold tracking-[2px] text-main">
      {{ ct('cartTitle') }}
    </div>

    <!--
    <div class="my-10 border border-orange bg-white px-6 py-7 text-orange">
      <p>{{ ct('cartCountLabel') }}{{ shoppingCart.carts.length }}</p>

      <ul class="list-disc px-6">
        <li>{{ ct('firstNotice') }}</li>
        <li>{{ ct('secondNotice') }}</li>
      </ul>
    </div>
    -->

    <div class="mt-10 flex flex-col gap-y-10">
      <the-marche-cart-item
        v-for="(cartItem, i) in shoppingCart.carts"
        :key="i"
        :cart="cartItem"
        :cart-number="cartItem.number"
        :coordinator="cartItem.coordinator"
        :items="cartItem.items"
        @click:cart-buy-button="handleClickCartBuyButton"
        @click:buy-button="handleClickBuyButton"
        @click:remove-item-from-cart="handelClickRemoveItemFromCartButton"
      />
    </div>
  </div>
</template>
