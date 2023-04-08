import { defineStore } from 'pinia'

import { useAuthStore } from '~/store'
import { OrderResponse, OrdersResponse } from '~/types/api'

export const useOrderStore = defineStore('order', {
  state: () => {
    return {
      orders: [] as OrdersResponse['orders'],
      totalItems: 0,
    }
  },

  actions: {
    /**
     * 注文一覧を取得する非同期関数
     * @param limit
     * @param offset
     * @returns
     */
    async fetchOrders(limit: number = 20, offset: number = 0): Promise<void> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(new Error('認証エラー'))
        }

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
    async getOrder(id: string): Promise<OrderResponse> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(new Error('認証エラー'))
        }
        const res = await this.orderApiClient(accessToken).v1GetOrder(id)
        return res.data
      } catch (error) {
        console.log(error)
        return this.errorHandler(error)
      }
    },
  },
})
