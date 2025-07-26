<script lang="ts" setup>
import { storeToRefs } from 'pinia'

import { useAlert } from '~/lib/hooks'
import { useCustomerStore, useOrderStore } from '~/store'
import { OrderStatus, OrderType } from '~/types/api'
import type { DraftOrderRequest, CompleteOrderRequest, RefundOrderRequest, UpdateOrderFulfillmentRequest, OrderFulfillment } from '~/types/api'
import type { FulfillmentInput } from '~/types/props'

const route = useRoute()
const orderStore = useOrderStore()
const customerStore = useCustomerStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const orderId = route.params.id as string

// Pinia storeからの参照を取得
const { customer } = storeToRefs(customerStore)

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
const { data, refresh, status, error } = useAsyncData(`order-${orderId}`, () => {
  console.log('注文詳細を取得します。')
  // fetchOrder関数を呼び出して注文情報を取得
  return orderStore.getOrder(orderId)
})

watch(data, (newData) => {
  // 注文情報が更新されたら、フォームデータを初期化
  if (newData) {
    completeFormData.value = {
      shippingMessage: newData.order.shippingMessage,
    }
    refundFormData.value = {
      description: newData.order.refund.reason || '',
    }
    fulfillmentsFormData.value = newData.order.fulfillments.map((fulfillment: OrderFulfillment): FulfillmentInput => ({
      fulfillmentId: fulfillment.fulfillmentId,
      shippingCarrier: fulfillment.shippingCarrier,
      trackingNumber: fulfillment.trackingNumber,
    }))
  }
})

// 注文情報を取得
const order = computed(() => {
  return data.value?.order || null
})

// コーディネーター情報を取得
const coordinator = computed(() => {
  return data.value?.coordinator
})

// 商品情報を取得
const products = computed(() => {
  return data.value?.products || []
})

watch(error, (err) => {
  // エラーが発生した場合、アラートを表示
  if (err) {
    show(err.message)
    console.error(err)
  }
})

// ローディング状態を返す関数
const isLoading = computed<boolean>(() => {
  // statusが'pending'の場合はローディング中と判断
  return status.value === 'pending' || loading.value
})

// 売上確定処理のハンドラー
const handleSubmitCapture = async (): Promise<void> => {
  try {
    loading.value = true
    await orderStore.captureOrder(orderId)
    refresh() // 処理成功後にデータを再取得
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
    refresh() // 処理成功後にデータを再取得
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
    refresh() // 処理成功後にデータを再取得
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
    refresh() // 処理成功後にデータを再取得
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
    refresh() // 処理成功後にデータを再取得
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

// Helper functions for button visibility (copied from OrderShow template)
const isAuthorized = (): boolean => {
  return order.value?.status === OrderStatus.WAITING
}

// 発送連絡時のメッセージ下書き保存 - 商品購入時のみ
const isPreservable = (): boolean => {
  if (!order.value || order.value.type !== OrderType.PRODUCT) {
    return false
  }
  const targets: OrderStatus[] = [
    OrderStatus.WAITING,
    OrderStatus.PREPARED,
    OrderStatus.COMPLETED,
  ]
  return targets.includes(order.value.status)
}

const isCompletable = (): boolean => {
  const targets: OrderStatus[] = [OrderStatus.PREPARED]
  return order.value ? targets.includes(order.value.status) : false
}

const isCancelable = (): boolean => {
  const targets: OrderStatus[] = [OrderStatus.WAITING, OrderStatus.PREPARED]
  return order.value ? targets.includes(order.value.status) : false
}

const isRefundable = (): boolean => {
  const targets: OrderStatus[] = [OrderStatus.COMPLETED]
  return order.value ? targets.includes(order.value.status) : false
}

// Action handlers for the buttons
const handleClickOpenCancelDialog = (): void => {
  cancelDialog.value = true
}

const handleClickOpenRefundDialog = (): void => {
  refundDialog.value = true
}
</script>

<template>
  <div>
    <v-alert
      v-show="isShow"
      :type="alertType"
      v-text="alertText"
    />

    <templates-order-show
      v-if="order"
      v-model:complete-form-data="completeFormData"
      v-model:refund-form-data="refundFormData"
      v-model:fulfillments-form-data="fulfillmentsFormData"
      v-model:cancel-dialog="cancelDialog"
      v-model:refund-dialog="refundDialog"
      :loading="isLoading"
      :is-alert="isShow"
      :alert-type="alertType"
      :alert-text="alertText"
      :order="order"
      :coordinator="coordinator"
      :customer="customer"
      :products="products"
      @submit:capture="handleSubmitCapture"
      @submit:draft="handleSubmitDraft"
      @submit:complete="handleSubmitComplete"
      @submit:update-fulfillment="handleSubmitUpdateFulfillment"
      @submit:cancel="handleSubmitCancel"
      @submit:refund="handleSubmitRefund"
    />
    
    <div
      class="position-fixed bottom-0 left-0 w-100 bg-white pa-4 text-right elevation-3"
    >
      <div class="d-inline-flex ga-4">
        <v-btn
          color="secondary"
          variant="outlined"
          @click="$router.back()"
        >
          戻る
        </v-btn>
        <v-btn
          v-show="isPreservable()"
          :loading="isLoading"
          variant="outlined"
          color="info"
          @click="handleSubmitDraft()"
        >
          下書き保存
        </v-btn>
        <v-btn
          v-show="isCompletable()"
          :loading="isLoading"
          variant="outlined"
          color="primary"
          @click="handleSubmitComplete()"
        >
          {{ order?.type === OrderType.PRODUCT ? '発送完了を通知' : 'レビュー依頼を送信' }}
        </v-btn>
        <v-btn
          v-show="isRefundable()"
          :loading="isLoading"
          variant="outlined"
          color="error"
          @click="handleClickOpenRefundDialog"
        >
          返金
        </v-btn>
      </div>
    </div>
  </div>
</template>
