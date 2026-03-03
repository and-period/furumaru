import { useApiClient } from '~/composables/useApiClient'
import { ExperienceApi, OrderApi, OrderStatus, ProductApi, TopApi } from '~/types/api/v1'
import type { ExperienceStatus, Order, ProductStatus, TopOrdersResponse, V1ExperiencesGetRequest, V1OrdersGetRequest, V1ProductsGetRequest, V1TopOrdersGetRequest } from '~/types/api/v1'
import type { TopOrderPeriodType } from '~/types'

export const useTopStore = defineStore('top', () => {
  const { create, errorHandler } = useApiClient()
  const topApi = () => create(TopApi)
  const orderApi = () => create(OrderApi)
  const productApi = () => create(ProductApi)
  const experienceApi = () => create(ExperienceApi)

  const orders = ref<TopOrdersResponse>({} as TopOrdersResponse)
  const pendingOrders = ref<Order[]>([])
  const pendingOrdersTotal = ref<number>(0)
  const hasPublishedProduct = ref<boolean>(true)
  const hasPublishedExperience = ref<boolean>(true)

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

  async function fetchPendingOrders(): Promise<void> {
    try {
      const params: V1OrdersGetRequest = {
        limit: 5,
        offset: 0,
        statuses: [
          OrderStatus.OrderStatusUnpaid,
          OrderStatus.OrderStatusWaiting,
          OrderStatus.OrderStatusPreparing,
          OrderStatus.OrderStatusShipped,
        ],
      }
      const res = await orderApi().v1OrdersGet(params)
      pendingOrders.value = res.orders
      pendingOrdersTotal.value = res.total
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function fetchPublicationStatus(): Promise<void> {
    try {
      const [productsRes, experiencesRes] = await Promise.all([
        productApi().v1ProductsGet({ limit: 200 } as V1ProductsGetRequest),
        experienceApi().v1ExperiencesGet({ limit: 200 } as V1ExperiencesGetRequest),
      ])

      const publishedProductStatuses: ProductStatus[] = [2, 3] // Presale, ForSale
      hasPublishedProduct.value = productsRes.products.some(
        p => (publishedProductStatuses as number[]).includes(p.status),
      )

      const publishedExperienceStatuses: ExperienceStatus[] = [3, 4] // Accepting, SoldOut
      hasPublishedExperience.value = experiencesRes.experiences.some(
        e => (publishedExperienceStatuses as number[]).includes(e.status),
      )
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  return {
    orders,
    pendingOrders,
    pendingOrdersTotal,
    hasPublishedProduct,
    hasPublishedExperience,
    fetchOrders,
    fetchPendingOrders,
    fetchPublicationStatus,
  }
})
