<script lang="ts" setup>
import { storeToRefs } from 'pinia'

import { useAlert } from '~/lib/hooks'
import { useCoordinatorStore, useCustomerStore, useOrderStore, useProductStore, usePromotionStore } from '~/store'
import type { DraftOrderRequest, CompleteOrderRequest, RefundOrderRequest, UpdateOrderFulfillmentRequest, OrderFulfillment } from '~/types/api'
import type { FulfillmentInput } from '~/types/props'

const route = useRoute()
const orderStore = useOrderStore()
const coordinatorStore = useCoordinatorStore()
const customerStore = useCustomerStore()
const promotionStore = usePromotionStore()
const productStore = useProductStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const orderId = route.params.id as string

const { order } = storeToRefs(orderStore)
const { coordinator } = storeToRefs(coordinatorStore)
const { customer } = storeToRefs(customerStore)
const { promotions } = storeToRefs(promotionStore)
const { products } = storeToRefs(productStore)

const loading = ref<boolean>(false)
const cancelDialog = ref<boolean>(false)
const refundDialog = ref<boolean>(false)
const completeFormData = ref<CompleteOrderRequest>({
  shippingMessage: '',
})
const refundFormData = ref<RefundOrderRequest>({
  description: '',
})
const fulfillmentsFormData = ref<FulfillmentInput[]>([])

const fetchState = useAsyncData('orders', async (): Promise<boolean> => {
  await fetchOrder()
  return true
})

const fetchOrder = async (): Promise<void> => {
  try {
    await orderStore.getOrder(orderId)
    const inputs = order.value.fulfillments.map((fulfillment: OrderFulfillment): FulfillmentInput => ({
      fulfillmentId: fulfillment.fulfillmentId,
      shippingCarrier: fulfillment.shippingCarrier,
      trackingNumber: fulfillment.trackingNumber,
    }))
    completeFormData.value = {
      shippingMessage: order.value.shippingMessage,
    }
    refundFormData.value = {
      description: order.value.refund?.reason || '',
    }
    fulfillmentsFormData.value = inputs
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const isLoading = (): boolean => {
  return fetchState?.pending?.value || loading.value
}

const handleSubmitCapture = async (): Promise<void> => {
  try {
    loading.value = true
    await orderStore.captureOrder(orderId)
    fetchState.refresh()
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    loading.value = false
  }
}

const handleSubmitDraft = async (): Promise<void> => {
  try {
    loading.value = true
    const req: DraftOrderRequest = { ...completeFormData.value }
    await orderStore.draftOrder(orderId, req)
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    loading.value = false
  }
}

const handleSubmitComplete = async (): Promise<void> => {
  try {
    loading.value = true
    await orderStore.completeOrder(orderId, completeFormData.value)
    fetchState.refresh()
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    loading.value = false
  }
}

const handleSubmitCancel = async (): Promise<void> => {
  try {
    loading.value = true
    await orderStore.cancelOrder(orderId)
    cancelDialog.value = false
    fetchState.refresh()
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    loading.value = false
  }
}

const handleSubmitRefund = async (): Promise<void> => {
  try {
    loading.value = true
    await orderStore.refundOrder(orderId, refundFormData.value)
    refundDialog.value = false
    fetchState.refresh()
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    loading.value = false
  }
}

const handleSubmitUpdateFulfillment = async (fulfillmentId: string): Promise<void> => {
  const payload = fulfillmentsFormData.value.find((formData: FulfillmentInput): boolean => {
    return formData.fulfillmentId === fulfillmentId
  })
  if (!payload) {
    return
  }

  try {
    loading.value = true
    const req: UpdateOrderFulfillmentRequest = { ...payload }
    await orderStore.updateFulfillment(orderId, fulfillmentId, req)
    fetchState.refresh()
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    loading.value = false
  }
}

try {
  await fetchState.execute()
}
catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-order-show
    v-model:complete-form-data="completeFormData"
    v-model:refund-form-data="refundFormData"
    v-model:fulfillments-form-data="fulfillmentsFormData"
    v-model:cancel-dialog="cancelDialog"
    v-model:refund-dialog="refundDialog"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :order="order"
    :coordinator="coordinator"
    :customer="customer"
    :promotions="promotions"
    :products="products"
    @submit:capture="handleSubmitCapture"
    @submit:draft="handleSubmitDraft"
    @submit:complete="handleSubmitComplete"
    @submit:update-fulfillment="handleSubmitUpdateFulfillment"
    @submit:cancel="handleSubmitCancel"
    @submit:refund="handleSubmitRefund"
  />
</template>
