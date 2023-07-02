import { defineStore } from 'pinia'
import { apiClient } from '~/plugins/api-client'
import { SchedulesResponse } from '~/types/api'

export const useScheduleStore = defineStore('schedule', {
  state: () => ({
    schedules: [] as SchedulesResponse['schedules'],
    total: 0
  }),

  actions: {
    /**
     * マルシェ開催スケジュール一覧を取得する非同期関数
     * @param limit
     * @param offset
     */
    async fetchSchedules (limit = 20, offset = 0): Promise<void> {
      try {
        const res = await apiClient.scheduleApi().v1ListSchedules()
        this.schedules = res.data.schedules
        this.total = res.data.total
      } catch (err) {
        return this.errorHandler(err)
      }
    }
  }
})
