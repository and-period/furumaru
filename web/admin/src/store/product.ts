import axios from 'axios'
import { defineStore } from 'pinia'

import { useAuthStore } from './auth'

import ApiClientFactory from '~/plugins/factory'
import {
  CreateProductRequest,
  ProductApi,
  ProductsResponseProductsInner,
  UploadImageResponse,
} from '~/types/api'
import {
  AuthError,
  ConnectionError,
  InternalServerError,
  ValidationError,
} from '~/types/exception'

export const useProductStore = defineStore('product', {
  state: () => {
    const apiClient = (token: string) => {
      const factory = new ApiClientFactory()
      return factory.create(ProductApi, token)
    }
    return {
      products: [] as ProductsResponseProductsInner[],
      totalItems: 0,
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

    /**
     * 商品画像をアップロードする非同期関数
     * @param payload
     * @returns
     */
    async uploadProductImage(payload: File): Promise<UploadImageResponse> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(
            new AuthError('認証エラー。再度ログインをしてください。')
          )
        }
        const res = await this.apiClient(accessToken).v1UploadProductImage(
          payload,
          {
            headers: {
              'Content-Type': 'multipart/form-data',
            },
          }
        )
        return res.data
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
            case 400:
              return Promise.reject(
                new ValidationError(
                  'このファイルはアップロードできません。',
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
     * 商品を作成する非同期関数
     */
    async createProduct(payload: CreateProductRequest): Promise<void> {
      try {
        const authStore = useAuthStore()
        const user = authStore.user
        const accessToken = authStore.accessToken
        if (!user || !accessToken) {
          return Promise.reject(
            new AuthError('認証エラー。再度ログインをしてください。')
          )
        }
        await this.apiClient(accessToken).v1CreateProduct({
          ...payload,
          inventory: Number(payload.inventory),
        })
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
            case 400:
              return Promise.reject(
                new ValidationError('入力項目に誤りがあります。', error)
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
