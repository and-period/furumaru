<script lang="ts" setup>
import type { ShoppingCart } from '~/types/store/shopping'

interface Props {
  isAuthenticated: boolean
  cartIsEmpty: boolean
  cartMenuMessage: string
  cartTotalPriceText: string
  cartTotalPriceTaxIncludedText: string
  totalPrice: number
  cartItems: ShoppingCart[]
  viewMycartText: string
  numberOfCartsText: string
  shippingFeeAnnotation: string
  shippingFeeAnnotationLinkText: string
  shippingFeeAnnotationCheckText: string
}

defineProps<Props>()

interface Emits {
  (e: 'click:buyButton'): void
  (e: 'click:removeItemFromCart', cartNumber: number, id: string): void
}

const emits = defineEmits<Emits>()

interface Expose {
  open: () => void
}

const area = ref<{ open: () => void, close: () => void }>({
  open: () => {},
  close: () => {},
})

const priceStringFormatter = (price: number): string => {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: 'JPY',
  }).format(price)
}

const handleClickBuyButton = () => {
  emits('click:buyButton')
  area.value.close()
}

const handleClickRemoveItemButton = (cartNumber: number, id: string) => {
  emits('click:removeItemFromCart', cartNumber, id)
}

const handleOpen = () => {
  area.value.open()
}

defineExpose<Expose>({
  open: handleOpen,
})
</script>

<template>
  <the-dropdown-with-icon ref="area">
    <template #icon>
      <div class="relative">
        <span
          v-if="!cartIsEmpty"
          class="absolute right-[2px] top-[-2px] inline-flex h-[8px] w-[8px] animate-ping rounded-full bg-orange opacity-75"
        />
        <span
          v-if="!cartIsEmpty"
          class="absolute right-[2px] top-[-2px] inline-flex h-[8px] w-[8px] rounded-full bg-orange"
        />
        <the-cart-icon
          id="header-cart-icon"
          fill="#604C3F"
        />
      </div>
    </template>
    <template #content>
      <div
        class="flex max-h-[calc(100vh_-_150px)] flex-col gap-y-4 overflow-auto p-4 leading-8"
      >
        <p v-html="cartMenuMessage" />
        <hr class="border-main">
        <div>
          {{ cartTotalPriceText }}:
          <p class="font-bold after:ml-2 after:text-[16px]">
            {{ priceStringFormatter(totalPrice)
            }}{{ cartTotalPriceTaxIncludedText }}
          </p>
        </div>
        <button
          class="w-full bg-orange py-1 text-white"
          @click="handleClickBuyButton"
        >
          {{ viewMycartText }}
        </button>

        <div class="border border-orange p-3 text-sm text-orange">
          {{ numberOfCartsText }}: {{ cartItems.length }}
          <p>
            {{ shippingFeeAnnotation }}
            <nuxt-link
              href="/legal-notice"
              class="underline"
            >
              {{ shippingFeeAnnotationLinkText }}
            </nuxt-link>
            {{ shippingFeeAnnotationCheckText }}
          </p>
        </div>

        <the-cart-item
          v-for="(item, i) in cartItems"
          :key="i"
          :cart-number="i + 1"
          :marche-name="item.coordinator.marcheName"
          :box-type="item.boxType"
          :box-size="item.boxSize"
          :use-rate="item.useRate"
          :shopping-cart="item"
          @click:buy-button="handleClickBuyButton"
          @click:remove-item-button="handleClickRemoveItemButton"
        />
      </div>
    </template>
  </the-dropdown-with-icon>
</template>
