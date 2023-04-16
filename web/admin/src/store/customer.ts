import axios from 'axios'
import { defineStore } from 'pinia'

import { UsersResponse } from '~/types/api'
import {
  AuthError,
  ConnectionError,
  InternalServerError
} from '~/types/exception'
import { apiClient } from '~/plugins/api-client'

export const useUserStore = defineStore('user', {
  state: () => ({
    users: [] as UsersResponse['users'],
    totalItems: 0
  }),
  actions: {
    /**
     * 顧客の一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     */
    async fetchUsers (limit = 20, offset = 0): Promise<void> {
      try {
        const res = await apiClient.userApi().v1ListUsers(limit, offset)
        this.users = res.data.users
        this.totalItems = res.data.total
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
    }
  }
})
