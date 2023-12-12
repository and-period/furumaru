import type { TopArchive, TopCommonResponse, TopLive } from '~/types/api'

export const useTopPageStore = defineStore('top-page', {
  state: () => {
    return {
      _lives: [] as TopLive[],
      archives: [] as TopArchive[],
    }
  },

  actions: {
    async getHomeContent() {
      const response: TopCommonResponse =
        await this.topPageApiClient().v1TopCommon()

      this._lives = response.lives
      this.archives = response.archives
    },
  },

  getters: {
    lives: (state) => {
      return [
        ...state._lives.map((live) => {
          return {
            ...live,
          }
        }),
      ]
    },
  },
})
