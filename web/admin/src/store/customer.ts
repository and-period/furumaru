import type { User, UserOrder, UserToList, V1UsersGetRequest, V1UsersUserIdDeleteRequest, V1UsersUserIdGetRequest, V1UsersUserIdOrdersGetRequest } from '~/types/api/v1'
import { useAddressStore } from '~/store'

export const useCustomerStore = defineStore('customer', {
  state: () => ({
    customer: {} as User,
    customers: [] as User[],
    customersToList: [] as UserToList[],
    orders: [] as UserOrder[],
    totalItems: 0,
    totalOrderCount: 0,
    totalPaymentCount: 0,
    totalProductAmount: 0,
    totalPaymentAmount: 0,
  }),

  actions: {
    /**
     * 顧客の一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     */
    async fetchCustomers(limit = 20, offset = 0): Promise<void> {
      try {
        const params: V1UsersGetRequest = {
          limit,
          offset,
        }
        const res = await this.userApi().v1UsersGet(params)
        this.customersToList = res.users
        this.totalItems = res.total
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 顧客を取得する非同期関数
     * @param customerId 顧客ID
     */
    async getCustomer(customerId: string): Promise<void> {
      try {
        const params: V1UsersUserIdGetRequest = {
          userId: customerId,
        }
        const res = await this.userApi().v1UsersUserIdGet(params)
        const addressStore = useAddressStore()
        this.customer = res.user
        addressStore.address = res.address
      }
      catch (err) {
        return this.errorHandler(err, {
          403: '購入者を閲覧する権限がありません',
          404: '対象の購入者が存在しません',
        })
      }
    },

    /**
     * 顧客の注文履歴一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     */
    async fetchCustomerOrders(customerId: string, limit = 20, offset = 0): Promise<void> {
      try {
        const params: V1UsersUserIdOrdersGetRequest = {
          userId: customerId,
          limit,
          offset,
        }
        const res = await this.userApi().v1UsersUserIdOrdersGet(params)
        this.orders = res.orders
        this.totalOrderCount = res.orderTotalCount
        this.totalPaymentCount = res.paymentTotalCount
        this.totalProductAmount = res.productTotalAmount
        this.totalPaymentAmount = res.paymentTotalAmount
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 顧客を削除する非同期関数
     */
    async deleteCustomer(customerId: string): Promise<void> {
      try {
        const params: V1UsersUserIdDeleteRequest = {
          userId: customerId,
        }
        await this.userApi().v1UsersUserIdDelete(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          403: '購入者を削除する権限がありません',
          404: '対象の購入者が存在しません',
        })
      }
    },
  },
})
