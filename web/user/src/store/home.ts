import type {
  Coordinator,
  TopArchive,
  TopCommonResponse,
  TopLive,
  VideoSummary,
} from '~/types/api'

export const useTopPageStore = defineStore('top-page', {
  state: () => {
    return {
      _coordinators: [] as Coordinator[],
      _lives: [] as TopLive[],
      archives: [] as TopArchive[],
      experienceVideos: [] as VideoSummary[],
      productVideos: [] as VideoSummary[],
    }
  },

  actions: {
    async getHomeContent() {
      const response: TopCommonResponse
        = await this.topPageApiClient().v1TopCommon()

      this._coordinators = response.coordinators
      this._lives = response.lives
      this.archives = response.archives
      this.experienceVideos = response.experienceVideos
      this.productVideos = response.productVideos
    },
  },

  getters: {
    lives: (state) => {
      return [
        ...state._lives.map((live) => {
          return {
            ...live,
            coordinator: state._coordinators.find(
              c => c.id === live.coordinatorId,
            ),
          }
        }),
      ]
    },
  },
})
