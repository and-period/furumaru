import { defineStore } from 'pinia'

import { apiClient } from '~/plugins/api-client'
import type {
  CategoriesResponse,
  Category,
  CreateCategoryRequest,
  UpdateCategoryRequest,
} from '~/types/api'

export const useCategoryStore = defineStore('category', {
  state: () => ({
    categories: [] as Category[],
    total: 0,
  }),

  actions: {
    /**
     * カテゴリ一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     * @param orders ソートキー
     */
    async fetchCategories(limit = 20, offset = 0, orders = []): Promise<void> {
      try {
        const res = await listCategories(limit, offset, '', orders)
        this.categories = res.categories
        this.total = res.total
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * カテゴリを検索をする非同期関数
     * @param name カテゴリ名(あいまい検索)
     * @param categoryIds stateの更新時に残しておく必要があるカテゴリ情報
     */
    async searchCategories(name = '', categoryIds: string[] = []): Promise<void> {
      try {
        const res = await listCategories(undefined, undefined, name, [])
        const categories: Category[] = []
        this.categories.forEach((category: Category): void => {
          if (!categoryIds.includes(category.id)) {
            return
          }
          categories.push(category)
        })
        res.categories.forEach((category: Category): void => {
          if (categories.find((v): boolean => v.id === category.id)) {
            return
          }
          categories.push(category)
        })
        this.categories = categories
        this.total = res.total
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * カテゴリを追加取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     * @param orders ソートキー
     */
    async moreCategories(limit = 20, offset = 0, orders = []): Promise<void> {
      try {
        const res = await listCategories(limit, offset, '', orders)
        this.categories.push(...res.categories)
        this.total = res.total
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * カテゴリを新規登録する非同期関数
     * @param payload
     */
    async createCategory(payload: CreateCategoryRequest): Promise<void> {
      try {
        const res = await apiClient.categoryApi().v1CreateCategory(payload)
        this.categories.unshift(res.data.category)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります',
          409: 'このカテゴリー名はすでに登録されています。',
        })
      }
    },

    /**
     * カテゴリを編集する非同期関数
     * @param categoryId カテゴリID
     * @param payload
     */
    async updateCategory(categoryId: string, payload: UpdateCategoryRequest) {
      try {
        await apiClient.categoryApi().v1UpdateCategory(categoryId, payload)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります',
          404: '対象のカテゴリーが存在しません',
          409: 'このカテゴリー名はすでに登録されています。',
        })
      }
      this.fetchCategories()
    },

    /**
     * カテゴリを削除する非同期関数
     * @param categoryId カテゴリID
     */
    async deleteCategory(categoryId: string): Promise<void> {
      try {
        await apiClient.categoryApi().v1DeleteCategory(categoryId)
      }
      catch (err) {
        return this.errorHandler(err, {
          404: '対象のカテゴリーが存在しません',
          412: '品目と紐付いているため削除できません',
        })
      }
      this.fetchCategories()
    },
  },
})

async function listCategories(limit = 20, offset = 0, name = '', orders: string[] = []): Promise<CategoriesResponse> {
  const res = await apiClient.categoryApi().v1ListCategories(limit, offset, name, orders.join(','))
  return { ...res.data }
}
