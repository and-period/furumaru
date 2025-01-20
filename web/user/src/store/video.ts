import { defineStore } from 'pinia'
import { useAuthStore } from './auth'
import type { VideoCommentsResponse } from '~/types/api'

export const useVideoStore = defineStore('video', {
  state: () => {
    return {}
  },

  actions: {
    async getVideo(id: string) {
      const res = await this.videoApiClient().v1GetVideo({
        videoId: id,
      })
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
})
