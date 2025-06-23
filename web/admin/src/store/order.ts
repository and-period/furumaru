import { defineStore } from 'pinia'
import { useCoordinatorStore } from './coordinator'
import { useCustomerStore } from './customer'
import { usePromotionStore } from './promotion'
import { useProductStore } from './product'
import { apiClient } from '~/plugins/api-client'
import type {
  CompleteOrderRequest,
  DraftOrderRequest,
  ExportOrdersRequest,
  Order,
  OrderResponse,
  RefundOrderRequest,
  UpdateOrderFulfillmentRequest,
} from '~/types/api'

export const useOrderStore = defineStore('order', {
  state: () => ({
    order: {} as Order,
    orders: [] as Order[],
    totalItems: 0,
  }),

  actions: {
    /**
     * 注文一覧を取得する非同期関数
     * @param limit
     * @param offset
     * @returns
     */
    async fetchOrders(limit = 20, offset = 0): Promise<void> {
      try {
        const res = await apiClient.orderApi().v1ListOrders(limit, offset)

        const coordinatorStore = useCoordinatorStore()
        const customerStore = useCustomerStore()
        const promotionStore = usePromotionStore()
        this.orders = res.data.orders
        this.totalItems = res.data.total
        coordinatorStore.coordinators = res.data.coordinators
        customerStore.customers = res.data.users
        promotionStore.promotions = res.data.promotions
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 注文IDから注文情報を取得する非同期関数
     * @param orderId 注文ID
     * @returns 注文情報
     */
    async getOrder(orderId: string): Promise<OrderResponse> {
      try {
        const res = await apiClient.orderApi().v1GetOrder(orderId)
        return res.data
      }
      catch (err) {
        return this.errorHandler(err, {
          403: '対象の注文を閲覧する権限がありません',
          404: '対象の注文が存在しません',
        })
      }
    },

    /**
     * 実売上状態にする非同期関数
     * @param orderId 注文ID
     * @returns
     */
    async captureOrder(orderId: string): Promise<void> {
      try {
        await apiClient.orderApi().v1CaptureOrder(orderId)
      }
      catch (err) {
        return this.errorHandler(err, {
          403: '対象の注文を実売り上げ状態にする権限がありません',
          404: '対象の注文が存在しません',
        })
      }
    },

    /**
     * 下書き保存する非同期関数
     * @param orderId 注文ID
     * @param payload 下書き情報
     * @returns
     */
    async draftOrder(
      orderId: string,
      payload: DraftOrderRequest,
    ): Promise<void> {
      try {
        await apiClient.orderApi().v1DraftOrder(orderId, payload)
      }
      catch (err) {
        return this.errorHandler(err, {
          403: '対象の注文を下書き保存する権限がありません',
          404: '対象の注文が存在しません',
        })
      }
    },

    /**
     * 注文の対応を完了にする非同期関数
     * @param orderId 注文ID
     * @param payload 対応完了時に必要な情報
     * @returns
     */
    async completeOrder(
      orderId: string,
      payload: CompleteOrderRequest,
    ): Promise<void> {
      try {
        await apiClient.orderApi().v1CompleteOrder(orderId, payload)
      }
      catch (err) {
        return this.errorHandler(err, {
          403: '対象の注文を完了状態にする権限がありません',
          404: '対象の注文が存在しません',
        })
      }
    },

    /**
     * 実売上前の注文に対してキャンセル処理をする非同期関数
     * @param orderId 注文ID
     * @returns
     */
    async cancelOrder(orderId: string): Promise<void> {
      try {
        await apiClient.orderApi().v1CancelOrder(orderId)
      }
      catch (err) {
        return this.errorHandler(err, {
          403: '対象の注文をキャンセルする権限がありません',
          404: '対象の注文が存在しません',
        })
      }
    },

    /**
     * 実売上語の注文に対して返金処理をする非同期関数
     * @param orderId 注文ID
     * @param payload 返金時に必要な情報
     * @returns
     */
    async refundOrder(
      orderId: string,
      payload: RefundOrderRequest,
    ): Promise<void> {
      try {
        await apiClient.orderApi().v1RefundOrder(orderId, payload)
      }
      catch (err) {
        return this.errorHandler(err, {
          403: '対象の注文を返金する権限がありません',
          404: '対象の注文が存在しません',
        })
      }
    },

    /**
     * 配送情報を更新する非同期関数
     * @param orderId 注文ID
     * @param fulfillmentId 配送ID
     * @param payload 配送情報
     * @returns
     */
    async updateFulfillment(
      orderId: string,
      fulfillmentId: string,
      payload: UpdateOrderFulfillmentRequest,
    ): Promise<void> {
      try {
        await apiClient
          .orderApi()
          .v1UpdateOrderFulfillment(orderId, fulfillmentId, payload)
      }
      catch (err) {
        return this.errorHandler(err, {
          403: '対象の注文の配送情報を更新する権限がありません',
          404: '対象の注文または配送情報が存在しません',
        })
      }
    },

    /**
     * 注文履歴一覧を一括取得
     * @param payload
     * @returns
     */
    async exportOrders(payload: ExportOrdersRequest): Promise<string> {
      try {
        const res = await apiClient.orderApi().v1ExportOrders(payload)
        return res.data
      }
      catch (err) {
        return this.errorHandler(err, {
          403: '注文履歴を取得する権限がありません',
        })
      }
    },
  },
})
