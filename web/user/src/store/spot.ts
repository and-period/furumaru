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

  getters: {
    spots: (state) => {
      return state.spotsResponse.spots.map((spot) => {
        return {
          ...spot,
          spotType: state.spotsResponse.spotTypes.find((spotType) => {
            return spotType.id === spot.spotTypeId
          }),
        }
      })
    },
  },

  actions: {
    /**
     * スポット一覧取得
     * @param longitude 経度
     * @param latitude 軽度
     * @param radius 取得半径（km）
     * @returns
     */
    async fetchSpots(
      longitude: number,
      latitude: number,
      radius?: number,
    ): Promise<void> {
      try {
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
