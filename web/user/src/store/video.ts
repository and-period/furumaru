import { defineStore } from 'pinia'
import { useAuthStore } from './auth'
import type { Category, Coordinator, Producer, Product, ProductTag, ProductType, VideoCommentsResponse, VideoResponse } from '~/types/api'

export const useVideoStore = defineStore('video', {
  state: () => {
    return {
      videoResponse: {} as VideoResponse,
      productsResponse: {
        total: 0,
        products: [] as Product[],
        producers: [] as Producer[],
        coordinators: [] as Coordinator[],
        categories: [] as Category[],
        productTypes: [] as ProductType[],
        productTags: [] as ProductTag[],
      },
      productsFetchState: {
        isLoading: false,
      },
    }
  },

  actions: {
    async getVideo(id: string) {
      const res = await this.videoApiClient().v1GetVideo({
        videoId: id,
      })
      this.videoResponse = res
      return res
    },

    async getComments(id: string): Promise<VideoCommentsResponse> {
      const res = await this.videoApiClient().v1ListVideoComments({
        videoId: id,
      })
      return res
    },

    async postComment(id: string, comment: string) {
      const authStore = useAuthStore()
      const { isAuthenticated, accessToken } = authStore
      try {
        if (isAuthenticated) {
          await this.videoApiClient(accessToken).v1CreateVideoComment({
            videoId: id,
            body: { comment },
          })
        }
        else {
          await this.videoApiClient().v1CreateGuestVideoComment({
            videoId: id,
            body: { comment },
          })
        }
      }
      catch (e) {
        return this.errorHandler(e)
      }
    },

  },

  getters: {
    products(state) {
      const products = Array.isArray(state.videoResponse.products) ? state.videoResponse.products : []
      return products.map((product) => {
        const thumbnail = product.media.find(m => m.isThumbnail)
        return {
          ...product,
          // 在庫があるかのフラグ
          hasStock: product.inventory > 0,
          // サムネイル画像のマッピング
          thumbnail,
          // サムネイルが動画かどうかのフラグ
          thumbnailIsVideo: thumbnail ? thumbnail.url.endsWith('.mp4') : false,
        }
      })
    },
  },
})
