import axios from 'axios'
import { defineStore } from 'pinia'

import { useAuthStore } from './auth'

import ApiClientFactory from '~/plugins/factory'
import { ProductApi, ProductsResponseProductsInner } from '~/types/api'
import {
  AuthError,
  ConnectionError,
  InternalServerError,
} from '~/types/exception'

export const useProductStore = defineStore('product', {
  state: () => {
    const apiClient = (token: string) => {
      const factory = new ApiClientFactory()
      return factory.create(ProductApi, token)
    }
    return {
      products: [] as ProductsResponseProductsInner[],
      apiClient,
    }
  },

  actions: {
    /**
     * 商品一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     * @returns
     */
    async fetchProducts(limit: number = 20, offset: number = 0): Promise<void> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(
            new AuthError('認証エラー。再度ログインをしてください。')
          )
        }
        const res = await this.apiClient(accessToken).v1ListProducts(
          limit,
          offset
        )
        this.products = res.data.products
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
