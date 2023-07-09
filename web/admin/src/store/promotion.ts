import { defineStore } from 'pinia'

import { useCommonStore } from './common'
import {
  CreatePromotionRequest,
  PromotionResponse,
  PromotionsResponse,
  PromotionsResponsePromotionsInner,
  UpdatePromotionRequest
} from '~/types/api'
import { apiClient } from '~/plugins/api-client'

export const usePromotionStore = defineStore('promotion', {
  state: () => ({
    promotion: {} as PromotionResponse,
    promotions: [] as PromotionsResponse['promotions'],
    total: 0
  }),
  actions: {
    /**
     * 登録済みのセール情報一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     * @param orders ソートキー
     */
    async fetchPromotions (limit = 20, offset = 0, orders: string[] = []): Promise<void> {
      try {
        const res = await apiClient.promotionApi().v1ListPromotions(limit, offset, orders.join(','))
        this.promotions = res.data.promotions
        this.total = res.data.total
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * セール情報を検索する非同期関数
     * @param name タイトル(あいまい検索)
     * @param promotionIds stateの更新時に残しておく必要があるセール情報
     */
    async searchPromotions (name = '', promotionIds: string[] = []): Promise<void> {
      try {
        const res = await apiClient.promotionApi().v1ListPromotions(undefined, undefined, name)
        const promotions: PromotionsResponsePromotionsInner[] = []
        this.promotions.forEach((promotion: PromotionsResponsePromotionsInner): void => {
          if (!promotionIds.includes(promotion.id)) {
            return
          }
          promotions.push(promotion)
        })
        res.data.promotions.forEach((promotion: PromotionsResponsePromotionsInner): void => {
          if (promotions.find((v): boolean => v.id === promotion.id)) {
            return
          }
          promotions.push(promotion)
        })
        this.promotions = promotions
        this.total = res.data.total
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * セールIDからセール情報情報を取得する非同期関数
     * @param promotionId セールID
     * @returns セールの情報
     */
    async getPromotion (promotionId: string): Promise<PromotionResponse> {
      try {
        const res = await apiClient.promotionApi().v1GetPromotion(promotionId)
        this.promotion = res.data
        return res.data
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * セール情報を登録する非同期関数
     * @param payload
     */
    async createPromotion (payload: CreatePromotionRequest): Promise<void> {
      try {
        await apiClient.promotionApi().v1CreatePromotion(payload)
        const commonStore = useCommonStore()
        commonStore.addSnackbar({
          message: `${payload.title}を作成しました。`,
          color: 'info'
        })
      } catch (err) {
        return this.errorHandler(err, { 409: 'このクーポンコードはすでに登録されています。' })
      }
    },

    /**
     * セール情報を削除する非同期関数
     * @param promotionId お知らせID
     */
    async deletePromotion (promotionId: string): Promise<void> {
      const commonStore = useCommonStore()
      try {
        await apiClient.promotionApi().v1DeletePromotion(promotionId)
        commonStore.addSnackbar({
          message: 'セール情報の削除が完了しました',
          color: 'info'
        })
        this.fetchPromotions()
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * セール情報を編集する非同期関数
     * @param promotionId セールID
     * @param payload
     */
    async updatePromotion (promotionId: string, payload: UpdatePromotionRequest): Promise<void> {
      const commonStore = useCommonStore()
      try {
        await apiClient.promotionApi().v1UpdatePromotion(promotionId, payload)
        commonStore.addSnackbar({
          message: 'セール情報の編集が完了しました',
          color: 'info'
        })
      } catch (err) {
        return this.errorHandler(err, { 409: 'このクーポンコードはすでに登録されています。' })
      }
    }
  }
})
