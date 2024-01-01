import { defineStore } from 'pinia'

import { useCategoryStore } from './category'
import { useCoordinatorStore } from './coordinator'
import { useProductTypeStore } from './product-type'
import { useProductTagStore } from './product-tag'
import { useProducerStore } from './producer'
import { apiClient } from '~/plugins/api-client'
import type {
  CreateProductRequest,
  ProductResponse,
  Product,
  UpdateProductRequest,
  UploadImageResponse
} from '~/types/api'
import { uploadTimeout } from '~/plugins/axios'

export const useProductStore = defineStore('product', {
  state: () => ({
    product: {} as Product,
    products: [] as Product[],
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
        const res = await apiClient.productApi().v1ListProducts(limit, offset)

        const coordinatorStore = useCoordinatorStore()
        const producerStore = useProducerStore()
        const categoryStore = useCategoryStore()
        const productTypeStore = useProductTypeStore()
        const productTagStore = useProductTagStore()
        this.products = res.data.products
        this.totalItems = res.data.total
        coordinatorStore.coordinators = res.data.coordinators
        producerStore.producers = res.data.producers
        categoryStore.categories = res.data.categories
        productTypeStore.productTypes = res.data.productTypes
        productTagStore.productTags = res.data.productTags
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 商品情報を検索する非同期関数
     * @param name 商品名(あいまい検索)
     * @param producerId 生産者ID
     * @param productIds stateの更新時に残しておく必要がある商品情報
     */
    async searchProducts (name = '', producerId = '', productIds: string[] = []): Promise<void> {
      try {
        const res = await apiClient.productApi().v1ListProducts(undefined, undefined, producerId, name)
        const products: Product[] = []
        this.products.forEach((product: Product): void => {
          if (!productIds.includes(product.id)) {
            return
          }
          products.push(product)
        })
        res.data.products.forEach((product: Product): void => {
          if (products.find((v): boolean => v.id === product.id)) {
            return
          }
          products.push(product)
        })
        this.products = products
        this.totalItems = res.data.total
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 商品詳細を取得する非同期関数
     * @param productId
     * @returns
     */
    async getProduct (productId: string): Promise<ProductResponse> {
      try {
        const res = await apiClient.productApi().v1GetProduct(productId)

        const coordinatorStore = useCoordinatorStore()
        const producerStore = useProducerStore()
        const categoryStore = useCategoryStore()
        const productTypeStore = useProductTypeStore()
        const productTagStore = useProductTagStore()
        this.product = res.data.product
        coordinatorStore.coordinators = [res.data.coordinator]
        producerStore.producers = [res.data.producer]
        categoryStore.categories = [res.data.category]
        productTypeStore.productTypes = [res.data.productType]
        productTagStore.productTags = res.data.productTags
        return res.data
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
            },
            timeout: uploadTimeout
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
     * 商品を更新する関数
     * @param productId
     * @param payload
     */
    async updateProduct (productId: string, payload: UpdateProductRequest) {
      try {
        await apiClient.productApi().v1UpdateProduct(productId, payload)
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
