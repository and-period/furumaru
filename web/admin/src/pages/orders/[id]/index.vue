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

// Pinia storeからの参照を取得
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

// 注文データの初期取得用の非同期処理
// 問題点: API呼び出しに失敗した場合もtrueを返しているため、UIが正常に描画されてしまう
const fetchState = useAsyncData('orders', async (): Promise<boolean> => {
  // fetchOrder関数を呼び出して注文情報を取得
  await fetchOrder()
  return true // 常にtrueを返している
})

// 注文情報を取得する関数
const fetchOrder = async (): Promise<void> => {
  try {
    // orderStoreを通じて注文情報を取得
    await orderStore.getOrder(orderId)
    // 取得したデータをフォーム入力用に変換
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
    // エラーハンドリング - ここでエラーをキャッチしているが、呼び出し元のfetchStateには伝播していない
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
    // 注: ここでエラーを再スローしていないため、fetchStateは常に成功として扱われる
  }
}

// ローディング状態を返す関数
const isLoading = (): boolean => {
  return fetchState?.pending?.value || loading.value
}

// 売上確定処理のハンドラー
const handleSubmitCapture = async (): Promise<void> => {
  try {
    loading.value = true
    await orderStore.captureOrder(orderId)
    fetchState.refresh() // 処理成功後にデータを再取得
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

// 下書き保存処理のハンドラー
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

// 完了処理のハンドラー
const handleSubmitComplete = async (): Promise<void> => {
  try {
    loading.value = true
    await orderStore.completeOrder(orderId, completeFormData.value)
    fetchState.refresh() // 処理成功後にデータを再取得
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

// キャンセル処理のハンドラー
const handleSubmitCancel = async (): Promise<void> => {
  try {
    loading.value = true
    await orderStore.cancelOrder(orderId)
    cancelDialog.value = false
    fetchState.refresh() // 処理成功後にデータを再取得
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

// 返金処理のハンドラー
const handleSubmitRefund = async (): Promise<void> => {
  try {
    loading.value = true
    await orderStore.refundOrder(orderId, refundFormData.value)
    refundDialog.value = false
    fetchState.refresh() // 処理成功後にデータを再取得
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

// 配送情報更新処理のハンドラー
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
    fetchState.refresh() // 処理成功後にデータを再取得
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

// コンポーネントマウント時に非同期データ取得を実行
// 問題点: ここでのエラーはキャッチしているが、UI側に状態が反映されていない
try {
  await fetchState.execute()
}
catch (err) {
  console.log('failed to setup', err)
  // 注: このエラーがキャッチされてもUI側にエラー状態が反映されていない
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
