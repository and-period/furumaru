import { defineStore } from 'pinia'
import { useShopStore } from './shop'

import { apiClient } from '~/plugins/api-client'
import type {
  CreatePromotionRequest,
  Promotion,
  UpdatePromotionRequest,
} from '~/types/api'

export const usePromotionStore = defineStore('promotion', {
  state: () => ({
    promotion: {} as Promotion,
    promotions: [] as Promotion[],
    total: 0,
  }),
  actions: {
    /**
     * 登録済みのセール情報一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     * @param orders ソートキー
     */
    async fetchPromotions(limit = 20, offset = 0, orders: string[] = []): Promise<void> {
      try {
        const res = await apiClient.promotionApi().v1ListPromotions(limit, offset, orders.join(','))

        const shopStore = useShopStore()
        this.promotions = res.data.promotions
        this.total = res.data.total
        shopStore.shops = res.data.shops
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * セール情報を検索する非同期関数
     * @param name タイトル(あいまい検索)
     * @param promotionIds stateの更新時に残しておく必要があるセール情報
     */
    async searchPromotions(name = '', promotionIds: string[] = []): Promise<void> {
      try {
        const res = await apiClient.promotionApi().v1ListPromotions(undefined, undefined, name)
        const promotions: Promotion[] = []
        this.promotions.forEach((promotion: Promotion): void => {
          if (!promotionIds.includes(promotion.id)) {
            return
          }
          promotions.push(promotion)
        })
        res.data.promotions.forEach((promotion: Promotion): void => {
          if (promotions.find((v): boolean => v.id === promotion.id)) {
            return
          }
          promotions.push(promotion)
        })

        const shopStore = useShopStore()
        this.promotions = promotions
        this.total = res.data.total
        shopStore.shops = res.data.shops
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * セールIDからセール情報情報を取得する非同期関数
     * @param promotionId セールID
     * @returns セールの情報
     */
    async getPromotion(promotionId: string): Promise<void> {
      try {
        const res = await apiClient.promotionApi().v1GetPromotion(promotionId)
        this.promotion = res.data.promotion
        if (!res.data.shop) {
          return
        }

        const shopStore = useShopStore()
        shopStore.shop = res.data.shop
      }
      catch (err) {
        return this.errorHandler(err, { 404: '対象のセール情報が存在しません。' })
      }
    },

    /**
     * セール情報を登録する非同期関数
     * @param payload
     */
    async createPromotion(payload: CreatePromotionRequest): Promise<void> {
      try {
        await apiClient.promotionApi().v1CreatePromotion(payload)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります。',
          409: 'このクーポンコードはすでに登録されています。',
        })
      }
    },

    /**
     * セール情報を削除する非同期関数
     * @param promotionId お知らせID
     */
    async deletePromotion(promotionId: string): Promise<void> {
      try {
        await apiClient.promotionApi().v1DeletePromotion(promotionId)
        this.fetchPromotions()
      }
      catch (err) {
        return this.errorHandler(err, { 404: '対象のセール情報が存在しません。' })
      }
    },

    /**
     * セール情報を編集する非同期関数
     * @param promotionId セールID
     * @param payload
     */
    async updatePromotion(promotionId: string, payload: UpdatePromotionRequest): Promise<void> {
      try {
        await apiClient.promotionApi().v1UpdatePromotion(promotionId, payload)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります。',
          404: '対象のセール情報が存在しません。',
          409: 'このクーポンコードはすでに登録されています。',
        })
      }
    },
  },
})
