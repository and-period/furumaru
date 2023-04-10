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

export const useProductTypeStore = defineStore('ProductType', {
  state: () => ({
    productTypes: [] as ProductTypesResponse['productTypes'],
    totalItems: 0,
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
                  'この品目はすでに登録されているため、登録できません。',
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
                  'この品目はすでに登録されているため、登録できません。',
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
                new NotFoundError('削除する品目が見つかりませんでした。', error)
              )
            case 500:
            default:
              return Promise.reject(new InternalServerError(error))
          }
        }
        throw new InternalServerError(error)
      }
      this.fetchProductTypes()
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
    }
  }
})
