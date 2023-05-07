import { defineStore } from 'pinia'
import { apiClient } from '~/plugins/api-client'

import { OrderResponse, OrdersResponse } from '~/types/api'

export const useOrderStore = defineStore('order', {
  state: () => ({
    order: {} as OrderResponse,
    orders: [] as OrdersResponse['orders'],
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
        const res = await apiClient.orderApi().v1ListOrders(
          limit,
          offset
        )
        this.orders = res.data.orders
        this.totalItems = res.data.total
      } catch (error) {
        console.log(error)
        this.errorHandler(error)
      }
    },
    /**
     * 注文IDから注文情報を取得する非同期関数
     * @param id 注文ID
     * @returns 注文情報
     */
    async getOrder (id: string): Promise<OrderResponse> {
      try {
        const res = await apiClient.orderApi().v1GetOrder(id)
        this.order = res.data
        return res.data
      } catch (error) {
        console.log(error)
        return this.errorHandler(error)
      }
    }
  }
})
