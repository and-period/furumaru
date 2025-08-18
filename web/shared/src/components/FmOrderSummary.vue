<script setup lang="ts">
import { computed } from 'vue'

// Types based on the API structure
export interface ProductMedia {
  url: string
  isThumbnail: boolean
}

export interface Product {
  id: string
  name: string
  price: number
  thumbnail?: ProductMedia
}

export interface CartItem {
  productId?: string
  quantity: number
  product?: Product
}

export interface Coordinator {
  id: string
  marcheName: string
  username: string
  prefecture: string
  city: string
}

export interface Cart {
  id: string
}

export interface Promotion {
  id?: string
  title?: string
  code?: string
}

export interface FmOrderSummaryProps {
  items: CartItem[]
  coordinator: Coordinator
  carts: Cart[]
  promotion?: Promotion
  subtotal: number
  discount: number
  shippingFee?: number
  total: number
  isLoading?: boolean
  // Coupon functionality
  enableCoupon?: boolean
  appliedCouponCode?: string
  invalidCoupon?: boolean
  // Text props for internationalization
  texts?: {
    title?: string
    shipFromLabel?: string
    coordinatorLabel?: string
    boxCountLabel?: string
    quantityLabel?: string
    couponAppliedMessage?: string
    couponPlaceholder?: string
    applyButtonText?: string
    couponInvalidMessage?: string
    itemTotalPriceLabel?: string
    applyCouponLabel?: string
    shippingFeeLabel?: string
    calculateNextPageMessage?: string
    totalPriceLabel?: string
  }
}

interface FmOrderSummaryEmits {
  (e: 'apply:coupon', code: string): void
  (e: 'cancel:coupon'): void
}

const props = withDefaults(defineProps<FmOrderSummaryProps>(), {
  isLoading: false,
  enableCoupon: false,
  appliedCouponCode: '',
  invalidCoupon: false,
  promotion: undefined,
  shippingFee: undefined,
  texts: () => ({
    title: '注文内容',
    shipFromLabel: '発送地：',
    coordinatorLabel: '取扱元：',
    boxCountLabel: '箱の数：',
    quantityLabel: '数量：',
    couponAppliedMessage: 'クーポンが適用されています',
    couponPlaceholder: 'クーポンコード',
    applyButtonText: '適用する',
    couponInvalidMessage: 'クーポンコードが正しくありません',
    itemTotalPriceLabel: '商品合計（税込）',
    applyCouponLabel: 'クーポン利用',
    shippingFeeLabel: '送料',
    calculateNextPageMessage: '次のページで計算されます',
    totalPriceLabel: '合計（税込）',
  })
})

const emits = defineEmits<FmOrderSummaryEmits>()

// Reactive state for coupon input
const couponCodeInput = defineModel<string>('couponCode', { default: '' })

// Computed properties
const priceFormatter = (price: number): string => {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: 'JPY',
  }).format(price)
}

const itemThumbnailAlt = (itemName: string): string => {
  return `商品画像: ${itemName}`
}

const hasCouponApplied = computed(() => {
  return props.appliedCouponCode && props.appliedCouponCode.length > 0
})

// Event handlers
const handleApplyCoupon = () => {
  if (couponCodeInput.value.trim()) {
    emits('apply:coupon', couponCodeInput.value.trim())
  }
}

const handleCancelCoupon = () => {
  couponCodeInput.value = ''
  emits('cancel:coupon')
}
</script>

<template>
  <div class="bg-base px-[16px] py-[24px] text-main md:p-10">
    <div class="text-[14px] font-bold tracking-[1.6px] md:text-[16px]">
      {{ texts.title }}
    </div>

    <template v-if="!isLoading">
      <!-- Coordinator Information -->
      <div class="my-[16px] text-[12px] tracking-[1.2px] md:my-6">
        <p>{{ coordinator.marcheName }}</p>
        <p>
          {{ texts.shipFromLabel }}{{ coordinator.prefecture }}{{ coordinator.city }}
        </p>
        <p>{{ texts.coordinatorLabel }}{{ coordinator.username }}</p>
        <p>{{ texts.boxCountLabel }}{{ carts.length }}</p>
      </div>

      <!-- Order Items -->
      <div>
        <div class="divide-y border-y">
          <div
            v-for="(item, i) in items"
            :key="i"
            class="grid grid-cols-5 py-2 text-[12px] tracking-[1.2px]"
          >
            <template v-if="item.product">
              <template v-if="item.product.thumbnail">
                <!-- Video Thumbnail -->
                <template v-if="item.product.thumbnail.url.endsWith('.mp4')">
                  <video
                    width="56px"
                    height="56px"
                    :src="item.product.thumbnail.url"
                    class="block aspect-square h-[56px] w-[56px]"
                    :alt="itemThumbnailAlt(item.product.name)"
                  />
                </template>
                <!-- Image Thumbnail -->
                <template v-else>
                  <img
                    :src="item.product.thumbnail.url"
                    :alt="itemThumbnailAlt(item.product.name)"
                    width="56"
                    height="56"
                    class="block aspect-square h-[56px] w-[56px] object-cover"
                  >
                </template>

                <!-- Product Info -->
                <div class="col-span-3 pl-[24px] md:pl-0">
                  <div>{{ item.product.name }}</div>
                  <div class="mt-4 md:mt-0 md:items-center md:justify-self-end md:text-right">
                    {{ texts.quantityLabel }}{{ item.quantity }}
                  </div>
                </div>

                <!-- Price -->
                <div class="flex items-center justify-self-end text-right">
                  {{ priceFormatter(item.product.price) }}
                </div>
              </template>
            </template>
          </div>
        </div>

        <!-- Coupon Section -->
        <template v-if="enableCoupon">
          <!-- Applied Coupon Display -->
          <template v-if="hasCouponApplied">
            <div class="leading-[1.4px] mt-4 flex justify-between rounded-lg border border-orange p-2 text-[12px] text-orange">
              <div class="flex items-center gap-1">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke-width="1.5"
                  stroke="currentColor"
                  class="h-4 w-4"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M9 12.75 11.25 15 15 9.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
                  />
                </svg>
                {{ texts.couponAppliedMessage }}
              </div>
              <button @click="handleCancelCoupon">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke-width="1.5"
                  stroke="currentColor"
                  class="h-4 w-4"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M6 18 18 6M6 6l12 12"
                  />
                </svg>
              </button>
            </div>
          </template>

          <!-- Coupon Input -->
          <template v-else>
            <div class="mt-4 flex gap-2">
              <div class="grow">
                <input
                  v-model="couponCodeInput"
                  type="text"
                  class="w-full border border-gray-300 bg-gray-50 p-2.5 text-[14px] md:text-[16px]"
                  :placeholder="texts.couponPlaceholder"
                >
              </div>
              <button
                class="whitespace-nowrap bg-orange p-2 text-[14px] text-white md:text-[16px]"
                @click="handleApplyCoupon"
              >
                {{ texts.applyButtonText }}
              </button>
            </div>
            <div
              v-if="invalidCoupon"
              class="mt-2 px-1 text-[12px] leading-[1.2px] text-red-600"
            >
              {{ texts.couponInvalidMessage }}
            </div>
          </template>
        </template>

        <!-- Price Breakdown -->
        <div class="mt-4 grid grid-cols-5 gap-y-4 border-y border-main py-6 text-[12px] tracking-[1.4px] md:grid-cols-2 md:text-[14px]">
          <!-- Subtotal -->
          <div class="col-span-2 md:col-span-1">
            {{ texts.itemTotalPriceLabel }}
          </div>
          <div class="col-span-3 text-right md:col-span-1">
            {{ priceFormatter(subtotal) }}
          </div>

          <!-- Discount -->
          <div class="col-span-2 md:col-span-1">
            {{ texts.applyCouponLabel }}
          </div>
          <div class="col-span-3 text-right md:col-span-1">
            {{ priceFormatter(discount) }}
          </div>

          <!-- Shipping Fee -->
          <div class="col-span-2 md:col-span-1">
            {{ texts.shippingFeeLabel }}
          </div>
          <div class="col-span-3 text-right md:col-span-1">
            <template v-if="shippingFee !== undefined">
              {{ priceFormatter(shippingFee) }}
            </template>
            <template v-else>
              {{ texts.calculateNextPageMessage }}
            </template>
          </div>
        </div>

        <!-- Total -->
        <div class="mt-6 grid grid-cols-2 text-[14px] font-bold tracking-[1.4px]">
          <div>{{ texts.totalPriceLabel }}</div>
          <div class="text-right">
            {{ priceFormatter(total) }}
          </div>
        </div>
      </div>
    </template>

    <!-- Loading State -->
    <template v-else>
      <div class="my-8 flex items-center justify-center">
        <div class="h-8 w-8 animate-spin rounded-full border-4 border-main border-t-transparent" />
      </div>
    </template>
  </div>
</template>