import { ref } from '@nuxtjs/composition-api'
import { defineStore } from 'pinia'

import { useAuthStore } from './auth'
import { useCommonStore } from './common'

import ApiClientFactory from '../plugins/factory'

import { useAuthStore } from './auth'

import {
  CategoriesResponse,
  CategoryApi,
  CreateCategoryRequest,
} from '~/types/api'

export const useCategoryStore = defineStore('Category', {
  state: () => ({
    categories: [] as CategoriesResponse['categories'],
    productTypeCategories: [] as CategoriesResponse['categories'],
  }),

  actions: {
    /**
     * カテゴリを全件取得する非同期関数
     * @param limit 取得上限数
     */
    async fetchCategories(limit?: number): Promise<void> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(new Error('認証エラー'))
        }

        const factory = new ApiClientFactory()
        const categoriesApiClient = factory.create(CategoryApi, accessToken)
        if (limit === undefined) {
          const res = await categoriesApiClient.v1ListCategories()
          this.categories = res.data.categories
        } else {
          const res = await categoriesApiClient.v1ListCategories(limit)
          this.productTypeCategories = res.data.categories
        }
      } catch (e) {
        console.log(e)
        throw new Error('Internal Server Error')
      }
    },

    /**
     * カテゴリを新規登録する非同期関数
     * @param payload
     */
    async createCategory(payload: CreateCategoryRequest): Promise<void> {
      const commonStore = useCommonStore()
      const errorMessage = ref<string>('')
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(new Error('認証エラー'))
        }

        const factory = new ApiClientFactory()
        const categoriesApiClient = factory.create(CategoryApi, accessToken)
        const res = await categoriesApiClient.v1CreateCategory(payload)
        this.categories.unshift(res.data)
        commonStore.addSnackbar({
          message: `カテゴリーを追加しました。`,
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
