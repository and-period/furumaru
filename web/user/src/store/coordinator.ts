import { defineStore } from 'pinia'

import type { CoordinatorResponse } from '~/types/api'

export const useCoordinatorStore = defineStore('coordinator', {
  state: () => {
    return {
      coordinatorFetchState: {
        isLoading: false,
      },
      coordinatorResponse: {} as CoordinatorResponse,
    }
  },

  actions: {
    /**
     * コーディネーターの詳細情報を取得する非同期関数
     * @param coordinatorId 対象のコーディネーターのID
     */
    async fetchCoordinator (id: string): Promise<void> {
      const response : CoordinatorResponse = await this.coordinatorApiClient().v1GetCoordinator({coordinatorId: id})
      this.coordinatorResponse = response
      this.coordinatorFetchState.isLoading = false
    },
  },

  getters: {
    coordnatorInfo(state) {
      return {
        ...state.coordinatorResponse.coordinator,
        // 関連product
        product: state.coordinatorResponse.productTypes,
      }
    }
  }
})
