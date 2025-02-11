import { defineStore } from 'pinia'

import { apiClient } from '~/plugins/api-client'
import type { Shop, UpdateShopRequest } from '~/types/api'

export const useShopStore = defineStore('shop', {
  state: () => ({
    shop: {} as Shop,
    shops: [] as Shop[],
  }),

  actions: {
    /**
     * 店舗情報を取得する非同期関数
     * @param shopId 店舗ID
     */
    async fetchShop(shopId: string): Promise<void> {
      try {
        const res = await apiClient.shopApi().v1GetShop(shopId)
        this.shop = res.data.shop
      }
      catch (err) {
        return this.errorHandler(err, {
          404: '店舗が見つかりません。',
        })
      }
    },

    /**
     * 店舗情報を取得する非同期関数
     * @param shopId 店舗ID
     */
    async updateShop(shopId: string, payload: UpdateShopRequest): Promise<void> {
      try {
        const res = await apiClient.shopApi().v1UpdateShop(shopId, payload)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、入力内容に誤りがあります。',
        })
      }
    },
  },
})
