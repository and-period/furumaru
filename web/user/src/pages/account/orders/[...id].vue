<script setup lang="ts">
import { useOrderStore } from '~/store/order'
const route = useRoute()

const orderId = computed<string>(() => {
  const id = route.params.id
  if (id instanceof Array) {
    return id[0]
  }
  return route.params.id as string
})

const orderStore = useOrderStore()
const { orderHistory } = storeToRefs(orderStore)
const { fetchOrderHistory } = orderStore

const { error } = useAsyncData(`orders-${orderId.value}`, () => {
  return fetchOrderHistory(orderId.value)
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
    <template v-if="!orderHistory"> </template>
    <template v-else>
      {{ orderHistory }}
    </template>
  </div>
</template>
