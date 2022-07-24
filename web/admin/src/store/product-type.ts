import { ref } from '@nuxtjs/composition-api'
import { defineStore } from 'pinia'

import { useAuthStore } from './auth'
import { useCommonStore } from './common'

import { ApiClientFactory } from '.'

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
    async fetchProductTypes(): Promise<void> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) throw new Error('認証エラー')

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
    async createProductType(
      categoryId: string,
      payload: CreateProductTypeRequest
    ): Promise<void> {
      const commonStore = useCommonStore()
      const errorMessage = ref<string>('')
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) throw new Error('認証エラー')

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
    },
  },
})
