import { defineStore } from 'pinia'

import { getAccessToken } from './auth'
import { AdministratorsResponse } from '~/types/api'

export const useAdministratorStore = defineStore('administrator', {
  state: () => ({
    administrators: [] as AdministratorsResponse['administrators'],
  }),
  actions: {
    async fetchAdministrators(): Promise<void> {
      try {
        const accessToken = getAccessToken()
        const res = await this.administratorApiClient(accessToken).v1ListAdministrators()
        this.administrators = res.data.administrators
      } catch (error) {
        // TODO: エラーハンドリング
        throw new Error('Internal Server Error')
      }
    },
  },
})
