<script lang="ts" setup>
import { storeToRefs } from 'pinia'

import { useAlert } from '~/lib/hooks'
import { useCoordinatorStore, useCustomerStore, useOrderStore, useProductStore, usePromotionStore } from '~/store'

const route = useRoute()
const orderStore = useOrderStore()
const coordinatorStore = useCoordinatorStore()
const customerStore = useCustomerStore()
const promotionStore = usePromotionStore()
const productStore = useProductStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const orderId = route.params.id as string

const { order } = storeToRefs(orderStore)
const { coordinators } = storeToRefs(coordinatorStore)
const { customers } = storeToRefs(customerStore)
const { promotions } = storeToRefs(promotionStore)
const { products } = storeToRefs(productStore)

const loading = ref<boolean>(false)

const fetchState = useAsyncData(async (): Promise<void> => {
  await fetchOrder()
})

const fetchOrder = async (): Promise<void> => {
  try {
    await orderStore.getOrder(orderId)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const isLoading = (): boolean => {
  return fetchState?.pending?.value || loading.value
}

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-order-show
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :order="order"
    :coordinators="coordinators"
    :customers="customers"
    :promotions="promotions"
    :products="products"
  />
</template>
