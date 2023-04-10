import { defineStore } from 'pinia'

import { AdministratorsResponse } from '~/types/api'
import { apiClient } from '~/plugins/api-client'

export const useAdministratorStore = defineStore('administrator', {
  state: () => ({
    administrators: [] as AdministratorsResponse['administrators']
  }),
  actions: {
    async fetchAdministrators (): Promise<void> {
      try {
        const res = await apiClient.administratorApi().v1ListAdministrators()
        this.administrators = res.data.administrators
      } catch (error) {
        // TODO: エラーハンドリング
        throw new Error('Internal Server Error')
      }
    }
  }
})
