import { defineStore } from 'pinia'

import { useCoordinatorStore } from './coordinator'
import { useCustomerStore } from './customer'
import { usePromotionStore } from './promotion'
import { useProductStore } from './product'
import { apiClient } from '~/plugins/api-client'
import type { Order } from '~/types/api'

export const useOrderStore = defineStore('order', {
  state: () => ({
    order: {} as Order,
    orders: [] as Order[],
    totalItems: 0
  }),

  actions: {
    /**
     * 注文一覧を取得する非同期関数
     * @param limit
     * @param offset
     * @returns
     */
    async fetchOrders (limit = 20, offset = 0): Promise<void> {
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
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 注文IDから注文情報を取得する非同期関数
     * @param orderId 注文ID
     * @returns 注文情報
     */
    async getOrder (orderId: string): Promise<void> {
      try {
        const res = await apiClient.orderApi().v1GetOrder(orderId)

        const coordinatorStore = useCoordinatorStore()
        const customerStore = useCustomerStore()
        const promotionStore = usePromotionStore()
        const productStore = useProductStore()
        this.order = res.data.order
        coordinatorStore.coordinator = res.data.coordinator
        customerStore.customer = res.data.user
        promotionStore.promotions.push(res.data.promotion)
        productStore.products = res.data.products
      } catch (err) {
        return this.errorHandler(err)
      }
    }
  }
})
