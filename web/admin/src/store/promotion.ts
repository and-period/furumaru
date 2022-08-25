import axios from 'axios'
import { defineStore } from 'pinia'

import { useAuthStore } from './auth'

import ApiClientFactory from '~/plugins/factory'
import { PromotionApi, PromotionsResponse } from '~/types/api'
import {
  AuthError,
  ConnectionError,
  InternalServerError,
} from '~/types/exception'

export const usePromotionStore = defineStore('Promotion', {
  state: () => ({
    promotions: [] as PromotionsResponse['promotions'],
  }),
  actions: {
    /**
     * 登録済みのセール情報一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     */
    async fetchPromotions(
      limit: number = 20,
      offset: number = 0
    ): Promise<void> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(
            new AuthError('認証エラー。再度ログインをしてください。')
          )
        }

        const factory = new ApiClientFactory()
        const promotionsApiClient = factory.create(PromotionApi, accessToken)
        const res = await promotionsApiClient.v1ListPromotions(limit, offset)
        this.promotions = res.data.promotions
      } catch (error) {
        if (axios.isAxiosError(error)) {
          if (!error.response) {
            return Promise.reject(new ConnectionError(error))
          }

          switch (error.response.status) {
            case 401:
              return Promise.reject(
                new AuthError('認証エラー。再度ログインをしてください。', error)
              )
            case 500:
            default:
              return Promise.reject(new InternalServerError(error))
          }
        }
        throw new InternalServerError(error)
      }
    },
  },
})
