import { useAuthStore } from './auth'
import type { OrderHistory } from '~/types/store/order'
import type {
  Coordinator,
  Order,
  OrderResponse,
  Product,
  Promotion,
} from '~/types/api'

export const useOrderStore = defineStore('order', {
  state: () => {
    return {
      _orders: [] as Order[],
      _coordinators: [] as Coordinator[],
      _promotions: [] as Promotion[],
      _products: [] as Product[],
      total: 0,
      orederResponse: undefined as OrderResponse | undefined,
      fetchState: {
        isLoading: false,
      },
    }
  },

  actions: {
    /**
     * 注文履歴を取得する非同期関数
     */
    async fetchOrderHistoryList(limit: number = 20, offset: number = 0) {
      try {
        this.fetchState.isLoading = true
        const authStore = useAuthStore()
        const res = await this.orderApiClient(
          authStore.accessToken,
        ).v1ListOrders({
          limit,
          offset,
        })
        this._orders = res.orders
        this._coordinators = res.coordinators
        this._promotions = res.promotions
        this._products = res.products
        this.total = res.total
      }
      catch (error) {
        return this.errorHandler(error)
      }
      finally {
        this.fetchState.isLoading = false
      }
    },

    /**
     * 指定した注文IDの注文履歴を取得する非同期関数
     */
    async fetchOrderHistory(orderId: string) {
      try {
        const authStore = useAuthStore()
        const res = await this.orderApiClient(authStore.accessToken).v1GetOrder(
          {
            orderId,
          },
        )
        this.orederResponse = res
      }
      catch (error) {
        return this.errorHandler(error)
      }
    },
  },

  getters: {
    orderHistories(): OrderHistory[] {
      return this._orders.map((order) => {
        return {
          ...order,
          // コーディネーターをマッピング
          coordinator: this._coordinators.find(
            coordinator => coordinator.id === order.coordinatorId,
          ),
          items: order.items.map((item) => {
            return {
              ...item,
              // 商品をマッピング
              product: this._products.find(
                product => product.id === item.productId,
              ),
            }
          }),
        }
      })
    },

    orderHistory(): OrderHistory | undefined {
      if (!this.orederResponse) {
        return undefined
      }

      return {
        ...this.orederResponse.order,
        coordinator: this.orederResponse.coordinator,
        items: this.orederResponse.order.items.map((item) => {
          return {
            ...item,
            product: this.orederResponse?.products.find(
              product => product.id === item.productId,
            ),
          }
        }),
      }
    },
  },
})
