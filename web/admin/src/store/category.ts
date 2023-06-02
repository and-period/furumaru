import { defineStore } from 'pinia'

import { useCommonStore } from './common'
import {
  CategoriesResponse,
  CreateCategoryRequest,
  UpdateCategoryRequest
} from '~/types/api'
import { apiClient } from '~/plugins/api-client'

export const useCategoryStore = defineStore('category', {
  state: () => ({
    categories: [] as CategoriesResponse['categories'],
    totalCategoryItems: 0
  }),

  actions: {
    /**
     * カテゴリを全件取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     */
    async fetchCategories (limit = 20, offset = 0): Promise<void> {
      try {
        const res = await listCategories(limit, offset)
        this.categories = res.categories
        this.totalCategoryItems = res.total
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * カテゴリを追加取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     */
    async moreCategories (limit = 20, offset = 0): Promise<void> {
      try {
        const res = await listCategories(limit, offset)
        this.categories.push(...res.categories)
        this.totalCategoryItems = res.total
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * カテゴリを新規登録する非同期関数
     * @param payload
     */
    async createCategory (payload: CreateCategoryRequest): Promise<void> {
      const commonStore = useCommonStore()
      try {
        const res = await apiClient.categoryApi().v1CreateCategory(payload)
        this.categories.unshift(res.data)
        commonStore.addSnackbar({
          message: 'カテゴリーを追加しました。',
          color: 'info'
        })
      } catch (err) {
        return this.errorHandler(err, { 409: 'このカテゴリー名はすでに登録されています。' })
      }
    },

    /**
     * カテゴリを編集する非同期関数
     * @param payload
     * @param categoryId
     */
    async editCategory (categoryId: string, payload: UpdateCategoryRequest) {
      const commonStore = useCommonStore()
      try {
        await apiClient.categoryApi().v1UpdateCategory(categoryId, payload)
        commonStore.addSnackbar({
          message: '変更しました。',
          color: 'info'
        })
      } catch (err) {
        return this.errorHandler(err, { 409: 'このカテゴリー名はすでに登録されています。' })
      }
      this.fetchCategories()
    },

    /**
     * カテゴリを削除する非同期関数
     * @param categoryId
     */
    async deleteCategory (categoryId: string): Promise<void> {
      const commonStore = useCommonStore()
      try {
        await apiClient.categoryApi().v1DeleteCategory(categoryId)
        commonStore.addSnackbar({
          message: 'カテゴリー削除が完了しました',
          color: 'info'
        })
      } catch (err) {
        return this.errorHandler(err)
      }
      this.fetchCategories()
    }
  }
})

async function listCategories (limit = 20, offset = 0): Promise<CategoriesResponse> {
  const res = await apiClient.categoryApi().v1ListCategories(limit, offset)
  return { ...res.data }
}
