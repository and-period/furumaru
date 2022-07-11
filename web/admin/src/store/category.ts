import { defineStore } from 'pinia'

import { useAuthStore } from './auth'

import { ApiClientFactory } from '.'

import {
  CategoriesResponse,
  CategoryApi,
  CreateCategoryRequest,
} from '~/types/api'

export const useCategoryStore = defineStore('Category', {
  state: () => ({
    categories: [] as CategoriesResponse['categories'],
  }),
  actions: {
    async createCategory(payload: CreateCategoryRequest): Promise<void> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) throw new Error('認証エラー')

        const factory = new ApiClientFactory()
        const categoriesApiClient = factory.create(CategoryApi, accessToken)
        await categoriesApiClient.v1CreateCategory(payload)
      } catch (error) {
        // TODO: エラーハンドリング
        console.log(error)
        throw new Error('Internal Server Error')
      }
    },
  },
})
