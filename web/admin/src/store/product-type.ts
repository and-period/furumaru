import { defineStore } from 'pinia'

import { useCommonStore } from './common'
import {
  CreateProductTypeRequest,
  ProductTagsResponseProductTagsInner,
  ProductTypesResponse,
  ProductTypesResponseProductTypesInner,
  UpdateProductTypeRequest,
  UploadImageResponse
} from '~/types/api'
import { apiClient } from '~/plugins/api-client'

export const useProductTypeStore = defineStore('productType', {
  state: () => ({
    productTypes: [] as ProductTypesResponse['productTypes'],
    totalItems: 0
  }),

  actions: {
    /**
     * 品目一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     * @param orders ソートキー
     */
    async fetchProductTypes (limit = 20, offset = 0, orders = []): Promise<void> {
      try {
        const res = await apiClient.productTypeApi().v1ListAllProductTypes(limit, offset, orders.join(','))
        this.productTypes = res.data.productTypes
        this.totalItems = res.data.total
      } catch (err) {
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
    async fetchProductTypesByCategoryId (categoryId: string, limit = 20, offset = 0): Promise<void> {
      try {
        const res = await apiClient.productTypeApi().v1ListProductTypes(categoryId)
        this.productTypes = res.data.productTypes
        this.totalItems = res.data.total
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 品目を検索する非同期関数
     * @param name 品目名(あいまい検索)
     * @param categoryId カテゴリ名
     * @param productTypeIds stateの更新時に残しておく必要がある品目情報
     */
    async searchProductTypes (name = '', categoryId = '', productTypeIds: string[] = []): Promise<void> {
      try {
        const res = await apiClient.productTypeApi().v1ListProductTypes(categoryId, undefined, undefined, name)
        const productTypes: ProductTypesResponseProductTypesInner[] = []
        this.productTypes.forEach((productType: ProductTypesResponseProductTypesInner): void => {
          if (!productTypeIds.includes(productType.id)) {
            return
          }
          productTypes.push(productType)
        })
        res.data.productTypes.forEach((productType: ProductTypesResponseProductTypesInner): void => {
          if (productTypes.find((v): boolean => v.id === productType.id)) {
            return
          }
          productTypes.push(productType)
        })
        this.productTypes = productTypes
        this.totalItems = res.data.total
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 品目を新規登録する非同期関数
     * @param categoryId 品目の親となるカテゴリのID
     * @param payload
     * @returns
     */
    async createProductType (
      categoryId: string,
      payload: CreateProductTypeRequest
    ): Promise<void> {
      const commonStore = useCommonStore()
      try {
        const res = await apiClient.productTypeApi().v1CreateProductType(
          categoryId,
          payload
        )
        this.productTypes.unshift(res.data)
        commonStore.addSnackbar({
          message: '品目を追加しました。',
          color: 'info'
        })
      } catch (err) {
        return this.errorHandler(err, { 409: 'この品目はすでに登録されているため、登録できません。' })
      }
    },

    /**
     * 登録済みの品目を更新する非同期関数
     * @param categoryId カテゴリID
     * @param productTypeId 品目ID
     * @param payload 品目情報
     * @returns
     */
    async updateProductType (
      categoryId: string,
      productTypeId: string,
      payload: UpdateProductTypeRequest
    ) {
      try {
        await apiClient.productTypeApi().v1UpdateProductType(
          categoryId,
          productTypeId,
          payload
        )
      } catch (err) {
        return this.errorHandler(err, { 409: 'この品目はすでに登録されているため、登録できません。' })
      }
    },

    /**
     * 品目を削除する非同期関数
     * @param categoryId カテゴリID
     * @param productTypeId 品目ID
     * @returns
     */
    async deleteProductType (
      categoryId: string,
      productTypeId: string
    ): Promise<void> {
      const commonStore = useCommonStore()
      try {
        await apiClient.productTypeApi().v1DeleteProductType(
          categoryId,
          productTypeId
        )
        commonStore.addSnackbar({
          message: '品目削除が完了しました',
          color: 'info'
        })
        this.fetchProductTypes()
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 品目画像をアップロードする関数
     * @param payload 品目画像のファイルオブジェクト
     * @returns アップロード後の品目画像のパスを含んだオブジェクト
     */
    async uploadProductTypeIcon (payload: File): Promise<UploadImageResponse> {
      try {
        const res = await apiClient.productTypeApi().v1UploadProductTypeIcon(
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
    }
  }
})
