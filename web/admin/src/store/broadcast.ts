import { apiClient } from '~/plugins/api-client'
import { Broadcast } from '~/types/api'

export const useBroadcastStore = defineStore('broadcast', {
  state: () => ({
    broadcast: {} as Broadcast
  }),

  actions: {
    /**
     * ライブ配信情報を取得する非同期関数
     * @param scheduleId マルシェ開催スケジュールID
     * @returns
     */
    async getBroadcastByScheduleId (scheduleId: string): Promise<void> {
      try {
        const res = await apiClient.broadcastApi().v1GetBroadcast(scheduleId)
        this.broadcast = res.data.broadcast
      } catch (err) {
        return this.errorHandler(err)
      }
    }
  }
})
