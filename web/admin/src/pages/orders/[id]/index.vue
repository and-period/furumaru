<script lang="ts" setup>
import { storeToRefs } from 'pinia'

import { useAlert } from '~/lib/hooks'
import { useCoordinatorStore, useCustomerStore, useOrderStore, useProductStore, usePromotionStore } from '~/store'
import type { DraftOrderRequest, CompleteOrderRequest, RefundOrderRequest, UpdateOrderFulfillmentRequest, OrderFulfillment } from '~/types/api'
import type { FulfillmentInput } from '~/types/props'

const route = useRoute()
const orderStore = useOrderStore()
const customerStore = useCustomerStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const orderId = route.params.id as string

// デバッグ: ページ初期化時のログ
console.log('[DEBUG] 注文詳細ページ初期化:', { orderId })

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
  // デバッグ: API呼び出し開始
  console.log('[DEBUG] 注文データ取得開始:', { orderId })

  // fetchOrder関数を呼び出して注文情報を取得
  return orderStore.getOrder(orderId)
})

// デバッグ: statusの変化を監視
watch(status, (newStatus, oldStatus) => {
  console.log('[DEBUG] useAsyncData status変化:', {
    orderId,
    oldStatus,
    newStatus,
    hasData: !!data.value,
    hasError: !!error.value,
  })
})

watch(data, (newData, oldData) => {
  // デバッグ: データの変化をログ出力
  console.log('[DEBUG] 注文データ更新:', {
    orderId,
    hasNewData: !!newData,
    hasOldData: !!oldData,
    orderStatus: newData?.order?.status,
    fulfillmentsCount: newData?.order?.fulfillments?.length,
    productsCount: newData?.products?.length,
  })

  // 注文情報が更新されたら、フォームデータを初期化
  if (newData) {
    console.log('[DEBUG] フォームデータ初期化開始')

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

    console.log('[DEBUG] フォームデータ初期化完了:', {
      shippingMessage: completeFormData.value.shippingMessage,
      refundDescription: refundFormData.value.description,
      fulfillmentsCount: fulfillmentsFormData.value.length,
    })
  }
})

// 注文情報を取得
const order = computed(() => {
  const orderData = data.value?.order || null
  console.log('[DEBUG] order computed実行:', {
    hasOrder: !!orderData,
    status: orderData?.status,
  })
  return orderData
})

// コーディネーター情報を取得
const coordinator = computed(() => {
  const coordinatorData = data.value?.coordinator
  console.log('[DEBUG] coordinator computed実行:', {
    hasCoordinator: !!coordinatorData,
  })
  return coordinatorData
})

// 商品情報を取得
const products = computed(() => {
  const productsData = data.value?.products || []
  console.log('[DEBUG] products computed実行:', {
    productsCount: productsData.length,
    productIds: productsData.map(p => p.id),
  })
  return productsData
})

watch(error, (err) => {
  // デバッグ: エラーの詳細をログ出力
  console.log('[DEBUG] エラー発生:', {
    orderId,
    error: err,
    errorMessage: err?.message,
    errorStack: err?.stack,
  })

  // エラーが発生した場合、アラートを表示
  if (err) {
    show(err.message)
    console.error('[ERROR] 注文データ取得エラー:', err)
  }
})

// ローディング状態を返す関数
const isLoading = computed<boolean>(() => {
  const isPending = status.value === 'pending'
  const isLocalLoading = loading.value
  const result = isPending || isLocalLoading

  console.log('[DEBUG] isLoading computed実行:', {
    orderId,
    status: status.value,
    isPending,
    isLocalLoading,
    result,
  })

  // statusが'pending'の場合はローディング中と判断
  return result
})

// 売上確定処理のハンドラー
const handleSubmitCapture = async (): Promise<void> => {
  console.log('[DEBUG] 売上確定処理開始:', { orderId })

  try {
    loading.value = true
    console.log('[DEBUG] 売上確定API呼び出し開始')
    await orderStore.captureOrder(orderId)
    console.log('[DEBUG] 売上確定API呼び出し成功')
    refresh() // 処理成功後にデータを再取得
  }
  catch (err) {
    console.log('[DEBUG] 売上確定処理エラー:', { orderId, error: err })
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    loading.value = false
    console.log('[DEBUG] 売上確定処理完了:', { orderId })
  }
}

// 下書き保存処理のハンドラー
const handleSubmitDraft = async (): Promise<void> => {
  console.log('[DEBUG] 下書き保存処理開始:', {
    orderId,
    formData: completeFormData.value,
  })

  try {
    loading.value = true
    const req: DraftOrderRequest = { ...completeFormData.value }
    console.log('[DEBUG] 下書き保存API呼び出し開始:', { request: req })
    await orderStore.draftOrder(orderId, req)
    console.log('[DEBUG] 下書き保存API呼び出し成功')
  }
  catch (err) {
    console.log('[DEBUG] 下書き保存処理エラー:', { orderId, error: err })
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    loading.value = false
    console.log('[DEBUG] 下書き保存処理完了:', { orderId })
  }
}

// 完了処理のハンドラー
const handleSubmitComplete = async (): Promise<void> => {
  console.log('[DEBUG] 完了処理開始:', {
    orderId,
    formData: completeFormData.value,
  })

  try {
    loading.value = true
    console.log('[DEBUG] 完了API呼び出し開始')
    await orderStore.completeOrder(orderId, completeFormData.value)
    console.log('[DEBUG] 完了API呼び出し成功')
    refresh() // 処理成功後にデータを再取得
  }
  catch (err) {
    console.log('[DEBUG] 完了処理エラー:', { orderId, error: err })
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    loading.value = false
    console.log('[DEBUG] 完了処理完了:', { orderId })
  }
}

// キャンセル処理のハンドラー
const handleSubmitCancel = async (): Promise<void> => {
  console.log('[DEBUG] キャンセル処理開始:', { orderId })

  try {
    loading.value = true
    console.log('[DEBUG] キャンセルAPI呼び出し開始')
    await orderStore.cancelOrder(orderId)
    console.log('[DEBUG] キャンセルAPI呼び出し成功')
    cancelDialog.value = false
    refresh() // 処理成功後にデータを再取得
  }
  catch (err) {
    console.log('[DEBUG] キャンセル処理エラー:', { orderId, error: err })
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    loading.value = false
    console.log('[DEBUG] キャンセル処理完了:', { orderId })
  }
}

// 返金処理のハンドラー
const handleSubmitRefund = async (): Promise<void> => {
  console.log('[DEBUG] 返金処理開始:', {
    orderId,
    formData: refundFormData.value,
  })

  try {
    loading.value = true
    console.log('[DEBUG] 返金API呼び出し開始')
    await orderStore.refundOrder(orderId, refundFormData.value)
    console.log('[DEBUG] 返金API呼び出し成功')
    refundDialog.value = false
    refresh() // 処理成功後にデータを再取得
  }
  catch (err) {
    console.log('[DEBUG] 返金処理エラー:', { orderId, error: err })
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    loading.value = false
    console.log('[DEBUG] 返金処理完了:', { orderId })
  }
}

// 配送情報更新処理のハンドラー
const handleSubmitUpdateFulfillment = async (fulfillmentId: string): Promise<void> => {
  console.log('[DEBUG] 配送情報更新処理開始:', {
    orderId,
    fulfillmentId,
    allFulfillments: fulfillmentsFormData.value,
  })

  const payload = fulfillmentsFormData.value.find((formData: FulfillmentInput): boolean => {
    return formData.fulfillmentId === fulfillmentId
  })

  console.log('[DEBUG] 更新対象の配送情報:', {
    orderId,
    fulfillmentId,
    foundPayload: !!payload,
    payload,
  })

  if (!payload) {
    console.log('[DEBUG] 更新対象の配送情報が見つかりません:', {
      orderId,
      fulfillmentId,
    })
    return
  }

  try {
    loading.value = true
    const req: UpdateOrderFulfillmentRequest = { ...payload }
    console.log('[DEBUG] 配送情報更新API呼び出し開始:', { request: req })
    await orderStore.updateFulfillment(orderId, fulfillmentId, req)
    console.log('[DEBUG] 配送情報更新API呼び出し成功')
    refresh() // 処理成功後にデータを再取得
  }
  catch (err) {
    console.log('[DEBUG] 配送情報更新処理エラー:', {
      orderId,
      fulfillmentId,
      error: err,
    })
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    loading.value = false
    console.log('[DEBUG] 配送情報更新処理完了:', {
      orderId,
      fulfillmentId,
    })
  }
}
</script>

<template>
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
</template>
