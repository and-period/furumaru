import { defineStore } from 'pinia'

import { apiClient } from '~/plugins/api-client'
import {
  CreateProductRequest,
  ProductResponse,
  ProductsResponseProductsInner,
  UpdateProductRequest,
  UploadImageResponse
} from '~/types/api'

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
      } catch (err) {
        return this.errorHandler(err)
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
      } catch (err) {
        return this.errorHandler(err, { 400: 'このファイルはアップロードできません。' })
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
      } catch (err) {
        return this.errorHandler(err)
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
      } catch (err) {
        return this.errorHandler(err)
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
      } catch (err) {
        return this.errorHandler(err)
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
