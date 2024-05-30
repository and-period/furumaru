<script setup lang="ts">
import dayjs from 'dayjs'
import { useOrderStore } from '~/store/order'
import { priceFormatter } from '~/lib/price'
import {
  getOrderStatusString,
  getOperationResultFromOrderStatus,
  getPaymentMethodNameByPaymentMethodType,
} from '~/lib/order'
import type { OrderStatus } from '~/types/api'

const route = useRoute()

const orderId = computed<string>(() => {
  const id = route.params.id
  if (id instanceof Array) {
    return id[0]
  }
  return route.params.id as string
})

const orderStore = useOrderStore()
const { fetchOrderHistory } = orderStore
const { orderHistory } = storeToRefs(orderStore)

const getClassNameFromOrderStatus = (status: OrderStatus) => {
  const operationResult = getOperationResultFromOrderStatus(status)
  switch (operationResult) {
    case 'success':
      return 'bg-green'
    case 'failed':
      return 'bg-red-700'
    case 'canceled':
      return 'bg-yellow-600'
    default:
      return ''
  }
}

const { error } = useAsyncData(`orders-${orderId.value}`, () => {
  return fetchOrderHistory(orderId.value)
})

useSeoMeta({
  title: '注文履歴',
})
</script>

<template>
  <div class="container mx-auto p-4 xl:p-0">
    <template v-if="error">
      <div class="border border-orange bg-white p-4 text-main">
        指定した注文を取得できませんでした。
        {{ error.message }}
      </div>
    </template>

    <template v-if="!orderHistory" />

    <template v-else>
      <div class="rounded bg-white p-4">
        <div class="flex flex-col gap-4 sm:grid sm:grid-cols-3 md:grid-cols-4">
          <div class="flex items-center justify-center sm:aspect-square">
            <nuxt-img
              provider="cloudFront"
              class="block max-w-[80px] rounded-full"
              width="80px"
              :src="orderHistory.coordinator?.thumbnailUrl"
              :alt="`${orderHistory.coordinator?.username}のサムネイル`"
            />
          </div>
          <dl class="col-span-2 md:col-span-3">
            <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
              <dt>注文ID</dt>
              <dd class="sm:col-span-2">
                {{ orderHistory.id }}
              </dd>
            </div>
            <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
              <dt>注文日時</dt>
              <dd class="sm:col-span-2">
                {{
                  dayjs
                    .unix(orderHistory.payment.orderedAt)
                    .format('YYYY/MM/DD HH:mm')
                }}
              </dd>
            </div>
            <div
              v-if="orderHistory.coordinator"
              class="py-2 sm:grid sm:grid-cols-3 sm:gap-4"
            >
              <dt>マルシェ名</dt>
              <dd class="sm:col-span-2">
                {{ orderHistory.coordinator.marcheName }}
              </dd>
            </div>
            <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
              <dt>注文品数</dt>
              <dd>
                {{ orderHistory.items.length }}
              </dd>
            </div>

            <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
              <dt>支払い金額</dt>
              <dd>
                {{ priceFormatter(orderHistory.payment.total) }}
              </dd>
            </div>

            <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
              <dt>決済方法</dt>
              <dd class="sm:col-span-2">
                {{
                  getPaymentMethodNameByPaymentMethodType(
                    orderHistory.payment.methodType,
                  )
                }}
              </dd>
            </div>

            <div class="py-2 sm:grid sm:grid-cols-3 sm:gap-4">
              <dt>ステータス</dt>
              <dd class="sm:col-span-2">
                <span
                  :class="getClassNameFromOrderStatus(orderHistory.status)"
                  class="inline-block rounded-lg px-2 py-1 text-[14px] text-white sm:mt-1"
                >
                  {{ getOrderStatusString(orderHistory.status) }}
                </span>
              </dd>
            </div>
          </dl>
        </div>

        <hr class="my-2 border-dashed">

        <!-- 注文商品情報 -->
        <div
          v-if="orderHistory.items.length > 0"
          class="hidden grid-cols-5 items-center border-b py-2 text-[12px] tracking-[1.2px] md:grid"
        >
          <div class="col-span-2">
            商品
          </div>
          <div class="text-right">
            価格（税込み）
          </div>
          <div class="text-right">
            数量
          </div>
          <div class="text-right">
            小計（税込み）
          </div>
        </div>
        <div
          v-for="(item, i) in orderHistory.items"
          :key="i"
        >
          <div
            v-if="item.product"
            class="md-text-[16px] flex grid-cols-5 flex-col gap-2 border-b py-2 text-[14px] md:grid md:items-center"
          >
            <div class="col-span-2 flex gap-4">
              <nuxt-img
                provider="cloudFront"
                width="64px"
                height="64px"
                class="aspect-square h-16 w-16 object-contain"
                :src="item.product.thumbnailUrl"
                :alt="`${item.product.name}のサムネイル画像`"
              />
              {{ item.product.name }}
            </div>
            <div class="text-right">
              <span class="md:hidden">価格：</span>
              {{ priceFormatter(item.price) }}
            </div>
            <div class="text-right">
              <span class="md:hidden">数量：</span>
              {{ item.quantity }}
            </div>
            <div class="text-right">
              <span class="md:hidden">小計：</span>
              {{ priceFormatter(item.price * item.quantity) }}
            </div>
          </div>
        </div>

        <div class="mt-4 grid grid-cols-3 md:grid-cols-5">
          <div class="col-span-2 col-start-4 flex flex-col gap-2 text-[14px]">
            <div class="grid grid-cols-2 text-right">
              <div class="px-4">
                商品合計
              </div>
              <div>
                {{ priceFormatter(orderHistory.payment.subtotal) }}
              </div>
            </div>
            <div class="grid grid-cols-2 text-right">
              <div class="px-4">
                配送料
              </div>
              <div>
                {{ priceFormatter(orderHistory.payment.shippingFee) }}
              </div>
            </div>

            <div class="grid grid-cols-2 text-right">
              <div class="px-4">
                割引
              </div>
              <div>
                {{ priceFormatter(orderHistory.payment.discount) }}
              </div>
            </div>

            <div class="grid grid-cols-2 text-right">
              <div class="px-4">
                消費税
              </div>
              <div>
                {{
                  priceFormatter(
                    orderHistory.payment.total
                      - orderHistory.payment.subtotal
                      - orderHistory.payment.shippingFee
                      - orderHistory.payment.discount,
                  )
                }}
              </div>
            </div>

            <div class="mt-2 grid grid-cols-2 text-right text-[16px] font-bold">
              <div class="px-4">
                合計
              </div>
              <div class="underline">
                {{ priceFormatter(orderHistory.payment.total) }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>
