import { useApiClient } from '~/composables/useApiClient'
import { TopApi } from '~/types/api/v1'
import type { TopOrdersResponse, V1TopOrdersGetRequest } from '~/types/api/v1'
import type { TopOrderPeriodType } from '~/types'

export const useTopStore = defineStore('top', () => {
  const { create, errorHandler } = useApiClient()
  const topApi = () => create(TopApi)

  const orders = ref<TopOrdersResponse>({} as TopOrdersResponse)

  async function fetchOrders(startAt?: number, endAt?: number, periodType?: TopOrderPeriodType): Promise<void> {
    try {
      const params: V1TopOrdersGetRequest = { startAt, endAt, periodType }
      const res = await topApi().v1TopOrdersGet(params)
      orders.value = res
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  return {
    orders,
    fetchOrders,
  }
})
