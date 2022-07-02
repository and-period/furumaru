import { defineStore } from 'pinia'

import { useAuthStore } from './auth'

import { ApiClientFactory } from '.'

import { ProducerApi, ProducersResponse } from '~/types/api'

export const useProducerStore = defineStore('Producer', {
  state: () => ({
    producers: [] as ProducersResponse['producers'],
  }),
  actions: {
    async fetchProducers(): Promise<void> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) throw new Error('認証エラー')

        const factory = new ApiClientFactory()
        const producersApiClient = factory.create(ProducerApi, accessToken)
        const res = await producersApiClient.v1ListProducers()
        this.producers = res.data.producers
      } catch (error) {
        // TODO: エラーハンドリング
        throw new Error('Internal Server Error')
      }
    },
  },
})
