import { defineStore } from 'pinia'

import { UserResponse, UsersResponse } from '~/types/api'
import { apiClient } from '~/plugins/api-client'

export const useCustomerStore = defineStore('customer', {
  state: () => ({
    customer: {} as UserResponse,
    customers: [] as UsersResponse['users'],
    totalItems: 0
  }),
  actions: {
    /**
     * 顧客の一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     */
    async fetchCustomers (limit = 20, offset = 0): Promise<void> {
      try {
        const res = await apiClient.userApi().v1ListUsers(limit, offset)
        this.customers = res.data.users
        this.totalItems = res.data.total
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 顧客を取得する非同期関数
     * @param customerId 顧客ID
     */
    async fetchCustomer (customerId: string): Promise<void> {
      try {
        const res = await apiClient.userApi().v1GetUser(customerId)
        this.customer = res.data
      } catch (err) {
        return this.errorHandler(err)
      }
    }
  }
})
