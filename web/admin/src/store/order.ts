import { useApiClient } from '~/composables/useApiClient'
import { useCoordinatorStore } from './coordinator'
import { useCustomerStore } from './customer'
import { usePromotionStore } from './promotion'
import { OrderApi } from '~/types/api/v1'
import type {
  CompleteOrderRequest,
  DraftOrderRequest,
  ExportOrdersRequest,
  Order,
  OrderResponse,
  RefundOrderRequest,
  UpdateOrderFulfillmentRequest,
  V1OrdersExportPostRequest,
  V1OrdersGetRequest,
  V1OrdersOrderIdCancelPostRequest,
  V1OrdersOrderIdCapturePostRequest,
  V1OrdersOrderIdCompletePostRequest,
  V1OrdersOrderIdDraftPostRequest,
  V1OrdersOrderIdFulfillmentsFulfillmentIdPatchRequest,
  V1OrdersOrderIdGetRequest,
  V1OrdersOrderIdRefundPostRequest,
} from '~/types/api/v1'

export const useOrderStore = defineStore('order', () => {
  const { create, errorHandler } = useApiClient()
  const orderApi = () => create(OrderApi)

  const order = ref<Order>({} as Order)
  const orders = ref<Order[]>([])
  const totalItems = ref<number>(0)

  async function fetchOrders(limit = 20, offset = 0): Promise<void> {
    try {
      const params: V1OrdersGetRequest = { limit, offset }
      const res = await orderApi().v1OrdersGet(params)

      const coordinatorStore = useCoordinatorStore()
      const customerStore = useCustomerStore()
      const promotionStore = usePromotionStore()
      orders.value = res.orders
      totalItems.value = res.total
      coordinatorStore.coordinators = res.coordinators
      customerStore.customers = res.users
      promotionStore.promotions = res.promotions
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function getOrder(orderId: string): Promise<OrderResponse> {
    try {
      const params: V1OrdersOrderIdGetRequest = { orderId }
      return await orderApi().v1OrdersOrderIdGet(params)
    }
    catch (err) {
      return errorHandler(err, {
        403: '対象の注文を閲覧する権限がありません',
        404: '対象の注文が存在しません',
      })
    }
  }

  async function captureOrder(orderId: string): Promise<void> {
    try {
      const params: V1OrdersOrderIdCapturePostRequest = { orderId }
      await orderApi().v1OrdersOrderIdCapturePost(params)
    }
    catch (err) {
      return errorHandler(err, {
        403: '対象の注文を実売り上げ状態にする権限がありません',
        404: '対象の注文が存在しません',
      })
    }
  }

  async function draftOrder(orderId: string, payload: DraftOrderRequest): Promise<void> {
    try {
      const params: V1OrdersOrderIdDraftPostRequest = {
        orderId,
        draftOrderRequest: payload,
      }
      await orderApi().v1OrdersOrderIdDraftPost(params)
    }
    catch (err) {
      return errorHandler(err, {
        403: '対象の注文を下書き保存する権限がありません',
        404: '対象の注文が存在しません',
      })
    }
  }

  async function completeOrder(orderId: string, payload: CompleteOrderRequest): Promise<void> {
    try {
      const params: V1OrdersOrderIdCompletePostRequest = {
        orderId,
        completeOrderRequest: payload,
      }
      await orderApi().v1OrdersOrderIdCompletePost(params)
    }
    catch (err) {
      return errorHandler(err, {
        403: '対象の注文を完了状態にする権限がありません',
        404: '対象の注文が存在しません',
      })
    }
  }

  async function cancelOrder(orderId: string): Promise<void> {
    try {
      const params: V1OrdersOrderIdCancelPostRequest = { orderId }
      await orderApi().v1OrdersOrderIdCancelPost(params)
    }
    catch (err) {
      return errorHandler(err, {
        403: '対象の注文をキャンセルする権限がありません',
        404: '対象の注文が存在しません',
      })
    }
  }

  async function refundOrder(orderId: string, payload: RefundOrderRequest): Promise<void> {
    try {
      const params: V1OrdersOrderIdRefundPostRequest = {
        orderId,
        refundOrderRequest: payload,
      }
      await orderApi().v1OrdersOrderIdRefundPost(params)
    }
    catch (err) {
      return errorHandler(err, {
        403: '対象の注文を返金する権限がありません',
        404: '対象の注文が存在しません',
      })
    }
  }

  async function updateFulfillment(
    orderId: string,
    fulfillmentId: string,
    payload: UpdateOrderFulfillmentRequest,
  ): Promise<void> {
    try {
      const params: V1OrdersOrderIdFulfillmentsFulfillmentIdPatchRequest = {
        orderId,
        fulfillmentId,
        updateOrderFulfillmentRequest: payload,
      }
      await orderApi().v1OrdersOrderIdFulfillmentsFulfillmentIdPatch(params)
    }
    catch (err) {
      return errorHandler(err, {
        403: '対象の注文の配送情報を更新する権限がありません',
        404: '対象の注文または配送情報が存在しません',
      })
    }
  }

  async function exportOrders(payload: ExportOrdersRequest): Promise<string> {
    try {
      const params: V1OrdersExportPostRequest = {
        exportOrdersRequest: payload,
      }
      return await orderApi().v1OrdersExportPost(params)
    }
    catch (err) {
      return errorHandler(err, {
        403: '注文履歴を取得する権限がありません',
      })
    }
  }

  return {
    order,
    orders,
    totalItems,
    fetchOrders,
    getOrder,
    captureOrder,
    draftOrder,
    completeOrder,
    cancelOrder,
    refundOrder,
    updateFulfillment,
    exportOrders,
  }
})
