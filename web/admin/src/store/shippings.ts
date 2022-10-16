import axios from 'axios'
import { defineStore } from 'pinia'

import { useAuthStore } from './auth'

import ApiClientFactory from '~/plugins/factory'
import { ShippingApi, ShippingsResponseShippingsInner } from '~/types/api'
import {
  AuthError,
  ConnectionError,
  InternalServerError,
} from '~/types/exception'

export const useShippingStore = defineStore('shippings', {
  state: () => {
    const apiClient = (token: string) => {
      const factory = new ApiClientFactory()
      return factory.create(ShippingApi, token)
    }
    return {
      shippings: [] as ShippingsResponseShippingsInner[],
      totalItems: 0,
      apiClient,
    }
  },

  actions: {
    /**
     * 配送情報一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     * @returns
     */
    async fetchShippings(
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
        const res = await this.apiClient(accessToken).v1ListShippings(
          limit,
          offset
        )
        this.shippings = res.data.shippings
        this.totalItems = res.data.total
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
