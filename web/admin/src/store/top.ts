import { defineStore } from 'pinia'

import { apiClient } from '~/plugins/api-client'
import type { TopOrderPeriodType, TopOrdersResponse } from '~/types/api'

export const useTopStore = defineStore('top', {
  state: () => ({
    orders: {} as TopOrdersResponse,
  }),

  actions: {
    /**
     * 注文集計結果を取得する非同期関数
     */
    async fetchOrders(startAt?: number, endAt?: number, periodType?: TopOrderPeriodType): Promise<void> {
      try {
        const res = await apiClient.topApi().v1TopOrders(startAt, endAt, periodType)

        this.orders = res.data
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },
  },
})
