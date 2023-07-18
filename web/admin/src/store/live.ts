import { defineStore } from 'pinia'
import { useProducerStore } from './producer'
import { useProductStore } from './product'
import { apiClient } from '~/plugins/api-client'
import { Live } from '~/types/api'

export const useLiveStore = defineStore('live', {
  state: () => ({
    lives: [] as Live[],
    total: 0
  }),

  actions: {
    async fetchLives (scheduleId: string): Promise<void> {
      try {
        const res = await apiClient.liveApi().v1ListLives(scheduleId)

        const producerStore = useProducerStore()
        const productStore = useProductStore()
        this.lives = res.data.lives
        this.total = res.data.total
        producerStore.producers = res.data.producers
        productStore.products = res.data.products
      } catch (err) {
        return this.errorHandler(err)
      }
    }
  }
})
