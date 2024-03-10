import { useAuthStore } from './auth'

export const useScheduleStore = defineStore('schedule', {
  state: () => {
    return {}
  },

  actions: {
    async getSchedule(id: string) {
      const res = await this.scheduleApiClient().v1GetSchedule({
        scheduleId: id,
      })
      return res
    },

    async getComments(id: string) {
      const res = await this.scheduleApiClient().v1ListLiveComments({
        scheduleId: id,
      })
      return res
    },

    async postComment(id: string, comment: string) {
      const authStore = useAuthStore()
      const { isAuthenticated, accessToken } = authStore

      if (isAuthenticated) {
        await this.scheduleApiClient(accessToken).v1CreateLiveComment({
          scheduleId: id,
          body: { comment },
        })
      } else {
        await this.scheduleApiClient().v1CreateGuestLiveComment({
          scheduleId: id,
          body: { comment },
        })
      }
    },
  },
})
