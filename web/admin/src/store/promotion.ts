import { defineStore } from 'pinia'

import { useCommonStore } from './common'
import {
  CreatePromotionRequest,
  PromotionResponse,
  PromotionsResponse,
  UpdatePromotionRequest
} from '~/types/api'
import { apiClient } from '~/plugins/api-client'

export const usePromotionStore = defineStore('promotion', {
  state: () => ({
    promotions: [] as PromotionsResponse['promotions'],
    total: 0
  }),
  actions: {
    /**
     * 登録済みのセール情報一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     */
    async fetchPromotions (limit = 20, offset = 0): Promise<void> {
      try {
        const res = await apiClient.promotionApi().v1ListPromotions(limit, offset)
        this.promotions = res.data.promotions
        this.total = res.data.total
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
        const res = await apiClient.promotionApi().v1CreatePromotion(payload)
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
     * @param id お知らせID
     */
    async deletePromotion (id: string): Promise<void> {
      const commonStore = useCommonStore()
      try {
        await apiClient.promotionApi().v1DeletePromotion(id)
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
     * セールIDからセール情報情報を取得する非同期関数
     * @param id セールID
     * @returns セールの情報
     */
    async getPromotion (id: string): Promise<PromotionResponse> {
      try {
        const res = await apiClient.promotionApi().v1GetPromotion(id)
        return res.data
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * セール情報を編集する非同期関数
     * @param id セールID
     * @param payload
     */
    async editPromotion (
      id: string,
      payload: UpdatePromotionRequest
    ): Promise<void> {
      const commonStore = useCommonStore()
      try {
        await apiClient.promotionApi().v1UpdatePromotion(id, payload)
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
