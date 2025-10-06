import { defineStore } from 'pinia'

import type { CoordinatorResponse } from '~/types/api'

export const useCoordinatorStore = defineStore('coordinator', {
  state: () => {
    return {
      coordinatorFetchState: {
        isLoading: false,
      },
      coordinatorResponse: {} as CoordinatorResponse,
    }
  },

  actions: {
    /**
     * コーディネーターの詳細情報を取得する非同期関数
     * @param coordinatorId 対象のコーディネーターのID
     */
    async fetchCoordinator(id: string): Promise<void> {
      const response: CoordinatorResponse = await this.coordinatorApiClient().v1GetCoordinator({ coordinatorId: id })
      this.coordinatorResponse = response
      this.coordinatorFetchState.isLoading = false
    },
  },

  getters: {
    coordinatorInfo(state) {
      return {
        ...state.coordinatorResponse.coordinator,
        // 関連product
        product: state.coordinatorResponse.productTypes,
      }
    },
    products(state) {
      const products = state.coordinatorResponse.products ?? []
      return products.map((product) => {
        const thumbnail = product.media.find(m => m.isThumbnail)
        return {
          ...product,
          // 在庫があるかのフラグ
          hasStock: product.inventory > 0,
          // サムネイル画像のマッピング
          thumbnail,
          // 生産者情報をマッピング
          producer: state.coordinatorResponse.producers?.find(
            producer => producer.id === product.producerId,
          ),
          // 商品タイプをマッピング
          productType: state.coordinatorResponse.productTypes?.find(
            productType => productType.id === product.productTypeId,
          ),
        }
      })
    },
    archives(state) {
      const archives = state.coordinatorResponse.archives ?? []
      return archives.map((archive) => {
        return {
          ...archive,
        }
      })
    },
    lives(state) {
      const lives = state.coordinatorResponse.lives ?? []
      return lives.map((live) => {
        return {
          ...live,
        }
      })
    },
    producers(state) {
      return state.coordinatorResponse.producers?.map((producer) => {
        return {
          ...producer,
          products: state.coordinatorResponse.products.filter((product) => {
            return product.producerId === producer.id
          },
          ),
        }
      })
    },
    experiences(state) {
      const experiences = state.coordinatorResponse.experiences ?? []
      return experiences.map((experience) => {
        return {
          ...experience,
        }
      })
    },
    videos(state) {
      const videos = state.coordinatorResponse.videos ?? []
      return videos.map((video) => {
        return {
          ...video,
        }
      })
    },
  },
})
