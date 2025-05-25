import { defineStore } from 'pinia'

import { apiClient } from '~/plugins/api-client'
import type { Shipping, ShippingsResponse, UpdateDefaultShippingRequest, UpsertShippingRequest } from '~/types/api'

export const useShippingStore = defineStore('shipping', {
  state: () => ({
    shipping: {} as Shipping,
  }),

  actions: {
    /**
     * コーディネーターが登録している配送先一覧を取得する非同期関数
     */
    async fetchShippings(coordinatorId: string): Promise<ShippingsResponse> {
      try {
        const res = await apiClient.shippingApi().v1ListShippings(coordinatorId)
        return res.data
      }
      catch (err) {
        return this.errorHandler(err, { 404: '対象のコーディネーターが見つかりません。' })
      }
    },

    /**
     * 指定した配送設定を取得する非同期関数
     * @param coordinatorId コーディネーターID
     * @param shippingId 配送設定ID
     * @returns
     */
    async fetchShipping(coordinatorId: string, shippingId: string): Promise<Shipping> {
      try {
        const res = await apiClient.shippingApi().v1GetShipping(coordinatorId, shippingId)
        return res.data.shipping
      }
      catch (err) {
        return this.errorHandler(err, {
          404: '配送設定が見つかりません。',
        })
      }
    },

    /**
     * デフォルト配送設定を取得する非同期関数
     * @returns
     */
    async fetchDefaultShipping(): Promise<void> {
      try {
        const res = await apiClient.shippingApi().v1GetDefaultShipping()
        this.shipping = res.data.shipping
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * デフォルトの配送設定を変更する非同期関数
     * @param payload
     * @returns
     */
    async updateDefaultShipping(payload: UpdateDefaultShippingRequest): Promise<void> {
      try {
        await apiClient.shippingApi().v1UpdateDefaultShipping(payload)
      }
      catch (err) {
        return this.errorHandler(err, { 400: '必須項目が不足しているか、入力内容に誤りがあります。' })
      }
    },

    /**
     * 指定したコーディネーターの配送設定を取得する非同期関数
     * @param coordinatorId
     * @returns
     */
    async fetchActiveShipping(coordinatorId: string): Promise<void> {
      try {
        const res = await apiClient.shippingApi().v1GetActiveShipping(coordinatorId)
        this.shipping = res.data.shipping
      }
      catch (err) {
        return this.errorHandler(err, { 404: '対象のコーディネーターが見つかりません。' })
      }
    },

    /**
     * 指定したコーディネーターの配送設定を変更する非同期関数
     * @param coordinatorId コーディネーターID
     * @param payload
     * @returns
     */
    async upsertShipping(coordinatorId: string, payload: UpsertShippingRequest): Promise<void> {
      try {
        await apiClient.shippingApi().v1UpsertShipping(coordinatorId, payload)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、入力内容に誤りがあります。',
          404: '対象のコーディネーターが見つかりません。',
        })
      }
    },
  },
})
