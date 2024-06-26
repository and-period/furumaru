import { defineStore } from 'pinia'

import { apiClient } from '~/plugins/api-client'
import type { User, UserOrder, UserToList } from '~/types/api'
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
        const res = await apiClient.userApi().v1ListUsers(limit, offset)
        this.customersToList = res.data.users
        this.totalItems = res.data.total
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
        const res = await apiClient.userApi().v1GetUser(customerId)
        const addressStore = useAddressStore()
        this.customer = res.data.user
        addressStore.address = res.data.address
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
        const res = await apiClient.userApi().v1ListUserOrders(customerId, limit, offset)
        this.orders = res.data.orders
        this.totalOrderCount = res.data.orderTotalCount
        this.totalPaymentCount = res.data.paymentTotalAmount
        this.totalProductAmount = res.data.productTotalAmount
        this.totalPaymentAmount = res.data.paymentTotalAmount
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
        await apiClient.userApi().v1DeleteUser(customerId)
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
