<script setup lang="ts">
import type { Coordinator } from '~/types/api'
import type { CartItem, ShoppingCart } from '~/types/store'
import type { I18n } from '~/types/locales'

interface Props {
  cartNumber: number
  coordinator: Coordinator
  cart: ShoppingCart
  items: CartItem[]
}

interface Emits {
  (e: 'click:buyButton', coordinatorId: string): void
  (e: 'click:cartBuyButton', coordinatorId: string, cartNumber: number): void
  (e: 'click:removeItemFromCart', cartNumber: number, id: string): void
}

const i18n = useI18n()

const ct = (str: keyof I18n['purchase']['cart']) => {
  return i18n.t(`purchase.cart.${str}`)
}

const props = defineProps<Props>()
const emits = defineEmits<Emits>()

const totalPrice = computed<number>(() => {
  return props.items
    .map(item => item.product.price * item.quantity)
    .reduce((sum, price) => sum + price)
})

const priceStringFormatter = (price: number): string => {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: 'JPY',
  }).format(price)
}

const handleBuyButton = () => {
  emits('click:buyButton', props.coordinator.id)
}

const handleCartBuyButton = () => {
  emits('click:cartBuyButton', props.coordinator.id, props.cartNumber)
}

const handelClickRemoveItemButton = (id: string) => {
  emits('click:removeItemFromCart', props.cartNumber, id)
}
</script>

<template>
  <div class="flex flex-col gap-8 text-main lg:grid lg:grid-flow-col">
    <div class="flex flex-col">
      <div class="text-[20px] font-bold tracking-[2px]">
        {{ coordinator.marcheName }}
      </div>

      <div class="my-9">
        <div>{{ ct('cartCountLabel') }}{{ cart.items.length }}</div>
        <div>{{ ct('shipFromLabel') }}{{ `${coordinator.prefecture}${coordinator.city}` }}</div>
        <div>{{ ct('coordinatorLabel') }}{{ coordinator.username }}</div>
      </div>

      <div class="flex items-center justify-between font-bold">
        <div class="text-[14px]">
          {{ ct('totalPriceLabel') }}
        </div>
        <div class="text-[20px]">
          {{ priceStringFormatter(totalPrice) }}
        </div>
      </div>

      <hr class="my-5 border-main">

      <div>{{ ct('shippingFeeNotice') }}</div>

      <button
        class="mt-8 bg-main p-[14px] text-[16px] text-white"
        @click="handleBuyButton"
      >
        {{ ct('checkoutButtonText') }}
      </button>
    </div>

    <div
      class="relative col-span-10 rounded-3xl bg-white px-4 py-6 lg:px-16 lg:py-12"
    >
      <div
        class="absolute -left-4 top-12 hidden h-8 w-8 rotate-45 bg-white lg:block"
      />

      <div class="flex flex-col gap-y-6 lg:gap-y-16">
        <div
          class="mb-7 rounded-2xl bg-base p-2 text-center font-bold tracking-[1.2px] text-main"
        >
          {{ ct('cartTitle') }} #{{ cartNumber }} -{{ coordinator.marcheName }}-
        </div>

        <div class="flex w-full flex-col text-main">
          <!-- PC、タブレットのみで表示する -->
          <div
            class="hidden grid-cols-5 items-center border-b py-2 text-[12px] tracking-[1.2px] md:grid"
          >
            <div class="col-span-2">
              {{ ct('productNameLabel') }}
            </div>
            <div>{{ ct('productPriceLabel') }}</div>
            <div>{{ ct('quantityLabel') }}</div>
            <div>{{ ct('subtotalLabel') }}</div>
          </div>

          <!-- PC、タブレットのみで表示する -->
          <div
            v-for="(item, j) in items"
            :key="j"
            class="hidden grid-cols-5 items-center border-b py-2 md:grid"
          >
            <div class="col-span-2 flex gap-4">
              <nuxt-img
                provider="cloudFront"
                :src="item.product.thumbnail.url"
                class="block h-16 w-16"
                height="64px"
                width="64px"
                :alt="`${item.product.name}のサムネイル画像`"
              />
              {{ item.product.name }}
            </div>
            <div>
              {{ priceStringFormatter(item.product.price) }}
            </div>
            <div class="inline-flex text-[14px]">
              {{ item.quantity }}
              <button
                class="ml-2 text-[12px] underline"
                type="button"
                @click="handelClickRemoveItemButton(item.product.id)"
              >
                {{ ct('deleteButtonText') }}
              </button>
            </div>
            <div>
              {{ priceStringFormatter(item.product.price * item.quantity) }}
            </div>
          </div>

          <!-- スマートフォンのみで表示する -->
          <div
            v-for="(item, j) in items"
            :key="j"
            class="flex gap-3 border-b py-2 md:hidden"
          >
            <nuxt-img
              provider="cloudFront"
              :src="item.product.thumbnail.url"
              class="block h-16 w-16"
              width="64px"
              height="64px"
              :alt="`${item.product.name}のサムネイル画像`"
            />
            <div class="flex grow flex-col justify-between">
              <div>
                {{ item.product.name }}
              </div>
              <div class="items-cneter flex justify-between">
                <div>
                  {{ priceStringFormatter(item.product.price) }}
                </div>

                <div class="flex items-center gap-4 text-[12px]">
                  <div>{{ ct('quantityLabel') }}: {{ item.quantity }}</div>
                  <button
                    class="text-[12px] underline"
                    type="button"
                    @click="handelClickRemoveItemButton(item.product.id)"
                  >
                    {{ ct('deleteButtonText') }}
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="flex flex-col gap-10 lg:mt-12 lg:flex-row">
          <!-- スマートフォン、タブレットのみで表示する -->
          <div class="lg:hidden">
            <div class="flex items-center justify-between">
              <div class="text-[14px]">
                {{ ct('subtotalLabel') }}
              </div>
              <div class="text-[20px]">
                {{ priceStringFormatter(totalPrice) }}
              </div>
            </div>
            <hr class="my-4 border-main">
            <div class="text-[12px]">
              {{ ct('shippingFeeNotice') }}
            </div>
          </div>

          <div class="flex grow gap-10">
            <the-cardboard
              :box-size="cart.size"
              :use-rate="cart.useRate"
            />

            <div class="flex grow flex-col gap-4 whitespace-nowrap">
              <div class="flex flex-col gap-2 text-center">
                <div class="rounded-3xl bg-base py-[3px] text-[12px]">
                  {{ ct('boxTypeLabel') }}
                </div>
                <div class="text-[14px]">
                  {{ cart.boxType }}
                </div>
              </div>

              <div class="flex flex-col gap-2 text-center">
                <div class="rounded-3xl bg-base py-[3px] text-[12px]">
                  {{ ct('boxSizeLabel') }}
                </div>
                <div class="text-[14px]">
                  {{ cart.boxSize }}
                </div>
              </div>

              <div class="flex flex-col gap-2 text-center">
                <div class="rounded-3xl bg-base py-[3px] text-[12px]">
                  {{ ct('utilizationRateLabel') }}
                </div>
                <div class="text-[14px]">
                  {{ cart.useRate }}
                </div>
              </div>
            </div>
          </div>

          <div class="flex flex-col gap-7">
            <div class="hidden lg:block">
              <div class="flex items-center justify-between">
                <div class="text-[14px]">
                  {{ ct('subtotalLabel') }}
                </div>
                <div class="text-[20px]">
                  {{ priceStringFormatter(totalPrice) }}
                </div>
              </div>
              <hr class="my-4 border-main">
              <div>{{ ct('shippingFeeNotice') }}</div>
            </div>

            <button
              class="bg-main p-[14px] text-[16px] text-white"
              @click="handleCartBuyButton"
            >
              {{ ct('checkoutButtonText') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
