import { fileUpload } from './helper'
import { useCategoryStore } from './category'
import { defineStore } from 'pinia'
import { apiClient } from '~/plugins/api-client'
import type {
  CreateProductTypeRequest,
  GetUploadUrlRequest,
  ProductType,
  UpdateProductTypeRequest,
} from '~/types/api'

export const useProductTypeStore = defineStore('productType', {
  state: () => ({
    productTypes: [] as ProductType[],
    totalItems: 0,
  }),

  actions: {
    /**
     * 品目一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     * @param orders ソートキー
     */
    async fetchProductTypes(
      limit = 20,
      offset = 0,
      orders = [],
    ): Promise<void> {
      try {
        const res = await apiClient
          .productTypeApi()
          .v1ListAllProductTypes(limit, offset, orders.join(','))

        const categoryStore = useCategoryStore()
        this.productTypes = res.data.productTypes
        this.totalItems = res.data.total
        categoryStore.categories = res.data.categories || []
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * カテゴリに紐づく品目一覧を取得する非同期関数
     * @param categoryId カテゴリID
     * @param limit 取得上限数
     * @param offset 取得開始位置
     * @returns
     */
    async fetchProductTypesByCategoryId(
      categoryId: string,
      limit = 20,
      offset = 0,
    ): Promise<void> {
      if (categoryId === '') {
        this.productTypes = []
        this.totalItems = 0
        return
      }

      try {
        const res = await apiClient
          .productTypeApi()
          .v1ListProductTypes(categoryId)

        this.productTypes = res.data.productTypes
        this.totalItems = res.data.total
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 品目を検索する非同期関数
     * @param name 品目名(あいまい検索)
     * @param categoryId カテゴリ名
     * @param productTypeIds stateの更新時に残しておく必要がある品目情報
     */
    async searchProductTypes(
      name = '',
      categoryId = '',
      productTypeIds: string[] = [],
    ): Promise<void> {
      try {
        const res = await apiClient
          .productTypeApi()
          .v1ListProductTypes(categoryId, undefined, undefined, name)
        const productTypes: ProductType[] = []
        this.productTypes.forEach((productType: ProductType): void => {
          if (!productTypeIds.includes(productType.id)) {
            return
          }
          productTypes.push(productType)
        })
        res.data.productTypes.forEach((productType: ProductType): void => {
          if (productTypes.find((v): boolean => v.id === productType.id)) {
            return
          }
          productTypes.push(productType)
        })
        this.productTypes = productTypes
        this.totalItems = res.data.total
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 品目を新規登録する非同期関数
     * @param categoryId 品目の親となるカテゴリのID
     * @param payload
     * @returns
     */
    async createProductType(
      categoryId: string,
      payload: CreateProductTypeRequest,
    ): Promise<void> {
      try {
        const res = await apiClient
          .productTypeApi()
          .v1CreateProductType(categoryId, payload)
        this.productTypes.unshift(res.data.productType)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります。',
          409: '対象の品目名はすでに登録されているため、登録できません。',
        })
      }
    },

    /**
     * 登録済みの品目を更新する非同期関数
     * @param categoryId カテゴリID
     * @param productTypeId 品目ID
     * @param payload 品目情報
     * @returns
     */
    async updateProductType(
      categoryId: string,
      productTypeId: string,
      payload: UpdateProductTypeRequest,
    ) {
      try {
        await apiClient
          .productTypeApi()
          .v1UpdateProductType(categoryId, productTypeId, payload)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります。',
          404: '対象の商品種別または品目が存在しません。',
          409: '対象の品目名はすでに登録されているため、登録できません。',
        })
      }
    },

    /**
     * 品目を削除する非同期関数
     * @param categoryId カテゴリID
     * @param productTypeId 品目ID
     * @returns
     */
    async deleteProductType(
      categoryId: string,
      productTypeId: string,
    ): Promise<void> {
      try {
        await apiClient
          .productTypeApi()
          .v1DeleteProductType(categoryId, productTypeId)
        this.fetchProductTypes()
      }
      catch (err) {
        return this.errorHandler(err, {
          404: '対象の商品種別または品目が存在しません。',
        })
      }
    },

    /**
     * 品目画像をアップロードする関数
     * @param payload 品目画像のファイルオブジェクト
     * @returns アップロード後の品目画像のパスを含んだオブジェクト
     */
    async uploadProductTypeIcon(payload: File): Promise<string> {
      const contentType = payload.type
      try {
        const body: GetUploadUrlRequest = {
          fileType: contentType,
        }
        const res = await apiClient
          .productTypeApi()
          .v1GetProductTypeIconUploadUrl(body)

        return await fileUpload(payload, res.data.key, res.data.url)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: 'このファイルはアップロードできません。',
        })
      }
    },
  },
})
