import type { Spot, SpotResponse, SpotType } from '~/types/api'
import type { GoogleMapSearchResult } from '~/types/store'

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
      spotFetchState: {
        isLoading: false,
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

    /**
     * スポット詳細取得
     * @param id スポットID
     * @returns
     */
    async fetchSpot(id: string): Promise<SpotResponse> {
      try {
        this.spotFetchState.isLoading = true
        return await this.spotApiClient().v1GetSpot({ spotId: id })
      }
      catch (error) {
        console.error(error)
        return this.errorHandler(error)
      }
      finally {
        this.spotFetchState.isLoading = false
      }
    },

    async search(address: string): Promise<GoogleMapSearchResult[]> {
      const geocoder = new google.maps.Geocoder()
      const response = await geocoder.geocode({ address })

      return response.results.map((result) => {
        return {
          formattedAddress: result.formatted_address,
          longitude: result.geometry.location.lng(),
          latitude: result.geometry.location.lat(),
        }
      })
    },
  },
})
