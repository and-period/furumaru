import axios from 'axios'
import { defineStore } from 'pinia'

import { useCommonStore } from './common'
import {
  CreateProductTypeRequest,
  ProductTypesResponse,
  UpdateProductTypeRequest,
  UploadImageResponse
} from '~/types/api'
import {
  AuthError,
  ConflictError,
  ConnectionError,
  InternalServerError,
  NotFoundError,
  ValidationError
} from '~/types/exception'
import { apiClient } from '~/plugins/api-client'

export const useProductTypeStore = defineStore('productType', {
  state: () => ({
    productTypes: [] as ProductTypesResponse['productTypes'],
    totalItems: 0
  }),
  actions: {
    /**
     * 品目を全件取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     */
    async fetchProductTypes (limit = 20, offset = 0): Promise<void> {
      try {
        const res = await apiClient.productTypeApi().v1ListAllProductTypes(
          limit,
          offset
        )
        this.productTypes = res.data.productTypes
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
    async editProductType (
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
