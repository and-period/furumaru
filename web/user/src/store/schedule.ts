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
  },
})
