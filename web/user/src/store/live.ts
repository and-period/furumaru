import { defineStore } from 'pinia'
import type { ArchiveSchedulesResponse } from '~/types/api'

export const useLiveStore = defineStore('live', {
  state: () => {
    return {
      archivesFetchState: {
        isLoading: false,
      },
      archiveResponse: {} as ArchiveSchedulesResponse
    }
  },

  actions: {
    async fetchArchves(limit = 20, offset = 0): Promise<void> {
      this.archivesFetchState.isLoading = true
      const response: ArchiveSchedulesResponse =
        await this.scheduleApiClient().v1ArchiveSchedules({
          limit,
          offset
        })
      this.archiveResponse = response
      this.archivesFetchState.isLoading = false
    }
  },

  getters: {
    totalArchivesCount(state) {
      return state.archiveResponse.total
    }
  }
})
