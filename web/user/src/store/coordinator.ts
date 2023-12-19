import { defineStore } from 'pinia'

import type { CoordinatorResponse } from '~/types/api'

export const useCoordinatorStore = defineStore('coordinator', {
  state: () => ({
    coordinatorResponse: {} as CoordinatorResponse,
  }),

  actions: {
    /**
     * コーディネーターの詳細情報を取得する非同期関数
     * @param coordinatorId 対象のコーディネーターのID
     */
    async fetchCoordinator (id: string): Promise<void> {
      const response : CoordinatorResponse = await this.coordinatorApiClient().v1GetCoordinator({
        coordinatorId: id,
      })
      console.log(response)
      this.coordinatorResponse = response
    },
  }
})
