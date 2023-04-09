import axios from 'axios'
import { defineStore } from 'pinia'

import { getAccessToken } from './auth'
import { useCommonStore } from './common'
import {
  CreatePromotionRequest,
  PromotionResponse,
  PromotionsResponse,
  UpdatePromotionRequest
} from '~/types/api'
import {
  AuthError,
  ConflictError,
  ConnectionError,
  InternalServerError,
  NotFoundError,
  ValidationError
} from '~/types/exception'

export const usePromotionStore = defineStore('Promotion', {
  state: () => ({
    promotions: [] as PromotionsResponse['promotions']
  }),
  actions: {
    /**
     * 登録済みのセール情報一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     */
    async fetchPromotions (limit = 20, offset = 0): Promise<void> {
      try {
        const accessToken = getAccessToken()
        const res = await this.promotionApiClient(accessToken).v1ListPromotions(limit, offset)
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

    /**
     * セール情報を登録する非同期関数
     * @param payload
     */
    async createPromotion (payload: CreatePromotionRequest): Promise<void> {
      try {
        const accessToken = getAccessToken()
        const res = await this.promotionApiClient(accessToken).v1CreatePromotion(payload)
        const commonStore = useCommonStore()
        commonStore.addSnackbar({
          message: `${payload.title}を作成しました。`,
          color: 'info'
        })
      } catch (error) {
        if (axios.isAxiosError(error)) {
          if (!error.response) {
            return Promise.reject(new ConnectionError(error))
          }
          const statusCode = error.response.status
          switch (statusCode) {
            case 400:
              return Promise.reject(
                new ValidationError('入力内容に誤りがあります。', error)
              )
            case 401:
              return Promise.reject(
                new AuthError('認証エラー。再度ログインをしてください。', error)
              )
            case 409:
              return Promise.reject(
                new ConflictError(
                  'このクーポンコードはすでに登録されています。',
                  error
                )
              )
            case 500:
            default:
              return Promise.reject(new InternalServerError(error))
          }
        }

        throw new InternalServerError(error)
      }
    },

    /**
     * セール情報を削除する非同期関数
     * @param id お知らせID
     */
    async deletePromotion (id: string): Promise<void> {
      const commonStore = useCommonStore()
      try {
        const accessToken = getAccessToken()
        await this.promotionApiClient(accessToken).v1DeletePromotion(id)
        commonStore.addSnackbar({
          message: 'セール情報の削除が完了しました',
          color: 'info'
        })
      } catch (error) {
        if (axios.isAxiosError(error)) {
          if (!error.response) {
            return Promise.reject(new ConnectionError(error))
          }
          const statusCode = error.response.status
          switch (statusCode) {
            case 400:
              return Promise.reject(
                new ValidationError(
                  '削除できませんでした。管理者にお問い合わせしてください。',
                  error
                )
              )
            case 401:
              return Promise.reject(
                new AuthError('認証エラー。再度ログインをしてください。', error)
              )
            case 404:
              return Promise.reject(
                new NotFoundError(
                  '削除するセール情報が見つかりませんでした。',
                  error
                )
              )
            case 500:
            default:
              return Promise.reject(new InternalServerError(error))
          }
        }
        throw new InternalServerError(error)
      }
      this.fetchPromotions()
    },

    /**
     * セールIDからセール情報情報を取得する非同期関数
     * @param id セールID
     * @returns セールの情報
     */
    async getPromotion (id: string): Promise<PromotionResponse> {
      try {
        const accessToken = getAccessToken()
        const res = await this.promotionApiClient(accessToken).v1GetPromotion(id)
        return res.data
      } catch (error) {
        if (axios.isAxiosError(error)) {
          if (!error.response) {
            return Promise.reject(new ConnectionError(error))
          }
          const statusCode = error.response.status
          switch (statusCode) {
            case 401:
              return Promise.reject(
                new AuthError('認証エラー。再度ログインをしてください', error)
              )
            case 404:
              return Promise.reject(
                new NotFoundError(
                  '一致するセール情報が見つかりませんでした。',
                  error
                )
              )
            case 500:
            default:
              return Promise.reject(new InternalServerError(error))
          }
        }
        throw new InternalServerError(error)
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
        const accessToken = getAccessToken()
        await this.promotionApiClient(accessToken).v1UpdatePromotion(id, payload)
        commonStore.addSnackbar({
          message: 'セール情報の編集が完了しました',
          color: 'info'
        })
      } catch (error) {
        if (axios.isAxiosError(error)) {
          if (!error.response) {
            return Promise.reject(new ConnectionError(error))
          }
          const statusCode = error.response.status
          switch (statusCode) {
            case 400:
              return Promise.reject(
                new ValidationError('入力内容に誤りがあります。', error)
              )
            case 401:
              return Promise.reject(
                new AuthError('認証エラー。再度ログインをしてください', error)
              )
            case 404:
              return Promise.reject(
                new NotFoundError(
                  '一致するセール情報が見つかりませんでした。',
                  error
                )
              )
            case 409:
              return Promise.reject(
                new ConflictError(
                  'このクーポンコードはすでに登録されています。',
                  error
                )
              )
            case 500:
            default:
              return Promise.reject(new InternalServerError(error))
          }
        }
        throw new InternalServerError(error)
      }
    }
  }
})
