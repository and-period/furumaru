import { defineStore } from 'pinia'

import { useAuthStore } from './auth'

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
      } catch (error) {
        // TODO: エラーハンドリング
        throw new Error('Internal Server Error')
      }
    },

    async createProductType(
      categoryId: string,
      payload: CreateProductTypeRequest
    ): Promise<void> {
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
        console.log(res)
      } catch (error) {
        // TODO: エラーハンドリング
        throw new Error('Internal Server Error')
      }
    },
  },
})
