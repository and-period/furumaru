import axios from 'axios'
import { defineStore } from 'pinia'
import { apiClient } from '~/plugins/api-client'

import {
  CreateProductRequest,
  ProductResponse,
  ProductsResponseProductsInner,
  UpdateProductRequest,
  UploadImageResponse
} from '~/types/api'
import {
  AuthError,
  ConnectionError,
  InternalServerError,
  ValidationError
} from '~/types/exception'

export const useProductStore = defineStore('product', {
  state: () => ({
    products: [] as ProductsResponseProductsInner[],
    totalItems: 0
  }),

  actions: {
    /**
     * 商品一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     * @returns
     */
    async fetchProducts (limit = 20, offset = 0): Promise<void> {
      try {
        const res = await apiClient.productApi().v1ListProducts(
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
    async uploadProductImage (payload: File): Promise<UploadImageResponse> {
      try {
        const res = await apiClient.productApi().v1UploadProductImage(
          payload,
          {
            headers: {
              'Content-Type': 'multipart/form-data'
            }
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
    async createProduct (payload: CreateProductRequest): Promise<void> {
      try {
        await apiClient.productApi().v1CreateProduct({
          ...payload,
          inventory: Number(payload.inventory)
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

    /**
     * 商品詳細を取得する非同期関数
     * @param id
     * @returns
     */
    async getProduct (id: string): Promise<ProductResponse> {
      try {
        const res = await apiClient.productApi().v1GetProduct(id)
        return res.data
      } catch (error) {
        return this.errorHandler(error)
      }
    },

    /**
     * 商品を更新する関数
     * @param id
     * @param payload
     */
    async updateProduct (id: string, payload: UpdateProductRequest) {
      try {
        await apiClient.productApi().v1UpdateProduct(id, payload)
      } catch (error) {
        return this.errorHandler(error)
      }
    },

    /**
     * 商品を削除する関数
     * @param productId
     * @returns
     */
    async deleteProduct (productId: string) {
      try {
        await apiClient.productApi().v1DeleteProduct(productId)
        const index = this.products.findIndex(product => product.id === productId)
        this.products.splice(index, 1)
        this.totalItems--
      } catch (err) {
        return this.errorHandler(err)
      }
    }
  }
})
