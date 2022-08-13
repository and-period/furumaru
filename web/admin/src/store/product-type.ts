import axios from 'axios'
import { defineStore } from 'pinia'

import ApiClientFactory from '../plugins/factory'

import { useAuthStore } from './auth'
import { useCommonStore } from './common'

import {
  CreateProductTypeRequest,
  ProductTypeApi,
  ProductTypesResponse,
} from '~/types/api'
import { AuthError, ConflictError, ConnectionError, InternalServerError, NotFoundError, ValidationError } from '~/types/exception'

export const useProductTypeStore = defineStore('ProductType', {
  state: () => ({
    productTypes: [] as ProductTypesResponse['productTypes'],
  }),
  actions: {
    /**
     * 品目を全件取得する非同期関数
     */
    async fetchProductTypes(): Promise<void> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(new Error('認証エラー'))
        }

        const factory = new ApiClientFactory()
        const productTypeApiClient = factory.create(ProductTypeApi, accessToken)
        const res = await productTypeApiClient.v1ListAllProductTypes()
        console.log(res)
        this.productTypes = res.data.productTypes
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
    async createProductType(
      categoryId: string,
      payload: CreateProductTypeRequest
    ): Promise<void> {
      const commonStore = useCommonStore()
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(new Error('認証エラー'))
        }

        const factory = new ApiClientFactory()
        const productTypeApiClient = factory.create(ProductTypeApi, accessToken)
        const res = await productTypeApiClient.v1CreateProductType(
          categoryId,
          payload
        )
        this.productTypes.unshift(res.data)
        commonStore.addSnackbar({
          message: `品目を追加しました。`,
          color: 'info',
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

    async deleteProductType(
      categoryId: string,
      productTypeId: string
    ): Promise<void> {
      const commonStore = useCommonStore()
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(new Error('認証エラー'))
        }

        const factory = new ApiClientFactory()
        const categoriesApiClient = factory.create(ProductTypeApi, accessToken)
        await categoriesApiClient.v1DeleteProductType(categoryId, productTypeId)
        commonStore.addSnackbar({
          message: '品目削除が完了しました',
          color: 'info',
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
                new ValidationError('削除できませんでした。管理者にお問い合わせしてください。', error)
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
  },
})
