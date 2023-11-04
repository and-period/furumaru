import { TopCommonResponse } from '~/types/api'

export const useTopPageStore = defineStore('top-page', {
  state: () => {
    return {
      _lives: [] as object[],
      _archives: [] as object[],
    }
  },

  actions: {
    async getHomeContent() {
      const response: TopCommonResponse =
        await this.topPageApiClient().v1TopCommon()

      this._lives = response.lives
      this._archives = response.archives
    },
  },
})
