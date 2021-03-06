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
    async fetchCategories(): Promise<void> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) throw new Error('認証エラー')

        const factory = new ApiClientFactory()
        const categoriesApiClient = factory.create(CategoryApi, accessToken)
        const res = await categoriesApiClient.v1ListCategories()
        console.log(res)
        this.categories = res.data.categories
      } catch (error) {
        // TODO: エラーハンドリング
        throw new Error('Internal Server Error')
      }
    },

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
