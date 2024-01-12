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

useAsyncData(`orders-${orderId.value}`, () => {
  return fetchOrderHistory(orderId.value)
})
</script>

<template>
  <div class="container">
    <template v-if="!orderHistory"> </template>
    <template v-else>
      {{ orderHistory }}
    </template>
  </div>
</template>
