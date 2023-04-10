import axios from 'axios'
import { defineStore } from 'pinia'

import { useCommonStore } from './common'
import {
  CategoriesResponse,
  CreateCategoryRequest,
  UpdateCategoryRequest
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

export const useCategoryStore = defineStore('Category', {
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
        const res = await apiClient.categoryApi().v1ListCategories(
          limit,
          offset
        )
        this.categories = res.data.categories
        this.totalCategoryItems = res.data.total
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
                  'このカテゴリー名はすでに登録されています。',
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
            case 404:
              return Promise.reject(
                new NotFoundError(
                  '編集するカテゴリーが見つかりませんでした。',
                  error
                )
              )
            case 409:
              return Promise.reject(
                new ConflictError(
                  'このカテゴリー名はすでに登録されています。',
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
                new NotFoundError(
                  '削除するカテゴリーが見つかりませんでした。',
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
      this.fetchCategories()
    }
  }
})
