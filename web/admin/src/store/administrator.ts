import { defineStore } from 'pinia'

import ApiClientFactory from '../plugins/factory'

import { useAuthStore } from './auth'

import { AdministratorApi, AdministratorsResponse } from '~/types/api'

export const useAdministratorStore = defineStore('administrator', {
  state: () => ({
    administrators: [] as AdministratorsResponse['administrators'],
  }),
  actions: {
    async fetchAdministrators(): Promise<void> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) throw new Error('認証エラー')

        const factory = new ApiClientFactory()
        const administratorApiClient = factory.create(
          AdministratorApi,
          accessToken
        )
        const res = await administratorApiClient.v1ListAdministrators()
        this.administrators = res.data.administrators
      } catch (error) {
        // TODO: エラーハンドリング
        throw new Error('Internal Server Error')
      }
    },
  },
})
