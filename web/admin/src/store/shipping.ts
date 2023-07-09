import { defineStore } from 'pinia'

import { apiClient } from '~/plugins/api-client'
import {
  CreateShippingRequest,
  ShippingResponse,
  Shipping,
  UpdateShippingRequest
} from '~/types/api'

export const useShippingStore = defineStore('shipping', {
  state: () => ({
    shipping: {} as Shipping,
    shippings: [] as Shipping[],
    totalItems: 0
  }),

  actions: {
    /**
     * 配送情報一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     * @returns
     */
    async fetchShippings (limit = 20, offset = 0): Promise<void> {
      try {
        const res = await apiClient.shippingApi().v1ListShippings(
          limit,
          offset
        )
        this.shippings = res.data.shippings
        this.totalItems = res.data.total
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 配送情報を検索する非同期関数
     * @param name 配送設定名(あいまい検索)
     * @param shippingIds stateの更新時に残しておく必要がある配送設定情報
     */
    async searchCoordinators (name = '', shippingIds: string[] = []): Promise<void> {
      try {
        const res = await apiClient.shippingApi().v1ListShippings(undefined, undefined, name)
        const shippings: Shipping[] = []
        this.shippings.forEach((shipping: Shipping): void => {
          if (!shippingIds.includes(shipping.id)) {
            return
          }
          shippings.push(shipping)
        })
        res.data.shippings.forEach((shipping: Shipping): void => {
          if (shippings.find((v): boolean => v.id === shipping.id)) {
            return
          }
          shippings.push(shipping)
        })
        this.shippings = shippings
        this.totalItems = res.data.total
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 指定したIDの配送設定情報を取得する非同期関数
     * @param shippingId 配送設定情報ID
     * @returns 配送設定情報
     */
    async getShipping (shippingId: string): Promise<ShippingResponse> {
      try {
        const res = await apiClient.shippingApi().v1GetShipping(shippingId)
        this.shipping = res.data.shipping
        return res.data
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 配送情報を新規作成する非同期関数
     * @param payload
     * @returns
     */
    async createShipping (payload: CreateShippingRequest): Promise<void> {
      try {
        await apiClient.shippingApi().v1CreateShipping(payload)
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 指定したIDの配送情報を更新する非同期関数
     * @param shippingId 配送情報ID
     * @param payload
     * @returns
     */
    async updateShipping (shippingId: string, payload: UpdateShippingRequest): Promise<void> {
      try {
        await apiClient.shippingApi().v1UpdateShipping(shippingId, payload)
      } catch (err) {
        return this.errorHandler(err)
      }
    }
  }
})
