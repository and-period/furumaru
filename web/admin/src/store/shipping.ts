import { defineStore } from 'pinia'

import { apiClient } from '~/plugins/api-client'
import type { Shipping, UpdateDefaultShippingRequest, UpsertShippingRequest } from '~/types/api'

export const useShippingStore = defineStore('shipping', {
  state: () => ({
    shipping: {} as Shipping
  }),

  actions: {
    /**
     * デフォルト配送設定を取得する非同期関数
     * @returns
     */
    async fetchDefaultShipping (): Promise<void> {
      try {
        const res = await apiClient.shippingApi().v1GetDefaultShipping()
        this.shipping = res.data.shipping
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * デフォルトの配送設定を変更する非同期関数
     * @param payload
     * @returns
     */
    async updateDefaultShipping (payload: UpdateDefaultShippingRequest): Promise<void> {
      try {
        await apiClient.shippingApi().v1UpdateDefaultShipping(payload)
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 指定したコーディネータの配送設定を取得する非同期関数
     * @param coordinatorId
     * @returns
     */
    async fetchShipping (coordinatorId: string): Promise<void> {
      try {
        const res = await apiClient.shippingApi().v1GetShipping(coordinatorId)
        this.shipping = res.data.shipping
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 指定したコーディネータの配送設定を変更する非同期関数
     * @param coordinatorId コーディネータID
     * @param payload
     * @returns
     */
    async upsertShipping (coordinatorId: string, payload: UpsertShippingRequest): Promise<void> {
      try {
        await apiClient.shippingApi().v1UpsertShipping(coordinatorId, payload)
      } catch (err) {
        return this.errorHandler(err)
      }
    }
  }
})
