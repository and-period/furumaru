import { ref } from '@nuxtjs/composition-api'
import { defineStore } from 'pinia'

import ApiClientFactory from '../plugins/factory'

import { useAuthStore } from './auth'
import { useCommonStore } from './common'

import {
  CreateProductTypeRequest,
  ProductTypeApi,
  ProductTypesResponse,
} from '~/types/api'

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
      } catch (e) {
        console.log(e)
        throw new Error('Internal Server Error')
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
      } catch (e) {
        return Promise.reject(new Error('不明なエラーが発生しました。'))
      }
    },

    async deleteProductType(
      categoryId: string,
      productTypeId: string
    ): Promise<void> {
      const commonStore = useCommonStore()
      const errorMessage = ref<string>('')
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
          message: 'カテゴリー削除が完了しました',
          color: 'info',
        })
      } catch (e) {
        // TODO: エラーハンドリングは今後見直していく
        if (e instanceof Error) {
          errorMessage.value = e.message
        } else {
          errorMessage.value =
            '不明なエラーが発生しました。お手数ですがご自身で入力してください。'
        }
      }
      commonStore.addSnackbar({
        message: errorMessage.value,
        color: 'error',
      })
      this.fetchProductTypes()
    },
  },
})
