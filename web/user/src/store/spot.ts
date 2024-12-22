import type { Spot, SpotType } from '~/types/api'

export const useSpotStore = defineStore('spot', {
  state: () => {
    return {
      spotsFetchState: {
        isLoading: false,
      },
      spotsResponse: {
        spots: [] as Spot[],
        spotTypes: [] as SpotType[],
      },
    }
  },

  actions: {
    async fetchSpots(
      longitude: number,
      latitude: number,
      radius?: number,
    ): Promise<void> {
      try {
        console.log('fetchSpots', longitude, latitude, radius)
        this.spotsFetchState.isLoading = true
        const response = await this.spotApiClient().v1ListSpots({
          longitude,
          latitude,
          radius,
        })
        this.spotsResponse = response
      }
      catch (error) {
        console.error(error)
        return this.errorHandler(error)
      }
      finally {
        this.spotsFetchState.isLoading = false
      }
    },
  },
})
