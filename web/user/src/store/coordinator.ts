import { defineStore } from 'pinia'

import type { Coordinator } from '~/types/api'

export const useCoordinatorStore = defineStore('coordinator', {
  state: () => ({
    coordinator: {} as Coordinator,
  }),

  actions: {
    /**
     * コーディネーターの詳細情報を取得する非同期関数
     * @param coordinatorId 対象のコーディネーターのID
     */
    async getCoordinator (coordinatorId: string): Promise<void> {
    },
  }
})
