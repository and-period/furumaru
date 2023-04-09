import { defineStore } from 'pinia'

import { OrderResponse, OrdersResponse } from '~/types/api'
import { getAccessToken } from './auth'

export const useOrderStore = defineStore('order', {
  state: () => ({
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
        const accessToken = getAccessToken()
        const res = await this.orderApiClient(accessToken).v1ListOrders(
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
        const accessToken = getAccessToken()
        const res = await this.orderApiClient(accessToken).v1GetOrder(id)
        return res.data
      } catch (error) {
        console.log(error)
        return this.errorHandler(error)
      }
    }
  }
})
