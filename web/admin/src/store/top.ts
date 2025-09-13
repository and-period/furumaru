import type { TopOrderPeriodType } from '~/types'
import type { TopOrdersResponse, V1TopOrdersGetRequest } from '~/types/api/v1'

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
        const params: V1TopOrdersGetRequest = {
          startAt,
          endAt,
          periodType,
        }
        const res = await this.topApi().v1TopOrdersGet(params)

        this.orders = res
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },
  },
})
