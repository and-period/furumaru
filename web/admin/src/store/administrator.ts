import { defineStore } from 'pinia'

import { AdministratorsResponse } from '~/types/api'
import { apiClient } from '~/plugins/api-client'

export const useAdministratorStore = defineStore('administrator', {
  state: () => ({
    administrators: [] as AdministratorsResponse['administrators'],
    total: 0
  }),

  actions: {
    /**
     * 管理者一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     */
    async fetchAdministrators (limit = 20, offset = 0): Promise<void> {
      try {
        const res = await apiClient.administratorApi().v1ListAdministrators(
          limit,
          offset
        )
        this.administrators = res.data.administrators
        this.total = res.data.total
      } catch (err) {
        return this.errorHandler(err)
      }
    }
  }
})
