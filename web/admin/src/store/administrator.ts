import { defineStore } from 'pinia'

import { AdministratorApi, AdministratorsResponse } from '~/types/api'

export const useAdministratorStore = defineStore('administrator', {
  state: () => ({
    administrators: [] as AdministratorsResponse['administrators'],
  }),
  actions: {
    async fetchAdministrators(): Promise<void> {
      try {
        const administratorApiClient = new AdministratorApi()
        const res = await administratorApiClient.v1ListAdministrators()
        this.administrators = res.data.administrators
      } catch (error) {
        // TODO: エラーハンドリング
        throw new Error('Internal Server Error')
      }
    },
  },
})
