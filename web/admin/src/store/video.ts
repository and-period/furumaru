import { fileUpload } from './helper'
import { apiClient } from '~/plugins/api-client'
import type {
  CreateVideoRequest,
  UpdateVideoRequest,
  VideoResponse,
  VideosResponse,
} from '~/types/api'

export const useVideoStore = defineStore('video', {
  state: () => ({
    videoResponse: null as VideosResponse | null,
  }),

  getters: {},

  actions: {
    /**
     * 動画一覧取得関数
     * @param limit
     * @param offset
     * @returns
     */
    async fetchVideos(limit = 20, offset = 0): Promise<void> {
      try {
        const res = await apiClient.videoApi().v1ListVideos(limit, offset)
        this.videoResponse = res.data
      }
      catch (error) {
        console.log(error)
        return this.errorHandler(error)
      }
    },

    /**
     * 動画詳細取得関数
     * @param id 動画ID
     * @returns
     */
    async fetchVideo(id: string): Promise<VideoResponse> {
      try {
        const res = await apiClient.videoApi().v1GetVideo(id)
        return res.data
      }
      catch (error) {
        console.log(error)
        return this.errorHandler(error)
      }
    },

    /**
     * 動画の新規登録関数
     * @param payload
     * @returns
     */
    async createVideo(payload: CreateVideoRequest): Promise<void> {
      try {
        const res = await apiClient.videoApi().v1CreateVideo(payload)
      }
      catch (error) {
        console.log(error)
        return this.errorHandler(error)
      }
    },

    /**
     * 動画の更新関数
     * @param id  動画ID
     * @param payload
     * @returns
     */
    async updateVideo(id: string, payload: UpdateVideoRequest): Promise<void> {
      try {
        await apiClient.videoApi().v1UpdateVideo(id, payload)
      }
      catch (error) {
        console.log(error)
        return this.errorHandler(error)
      }
    },

    /**
     * 動画ファイルアップロード関数
     * @param file 動画ファイル
     * @returns 参照先のURL
     */
    async uploadVideoFile(file: File): Promise<string> {
      try {
        const contentType = file.type
        const res = await apiClient
          .videoApi()
          .v1GetVideoFileUploadUrl({ fileType: contentType })

        return await fileUpload(file, res.data.key, res.data.url)
      }
      catch (error) {
        console.log(error)
        return this.errorHandler(error)
      }
    },

    /**
     * サムネイルファイルアップロード関数
     * @param file 画像ファイル
     * @returns 参照先のURL
     */
    async uploadThumbnailFile(file: File): Promise<string> {
      try {
        const contentType = file.type
        const res = await apiClient
          .videoApi()
          .v1GetVideoThumbnailUploadUrl({ fileType: contentType })

        return await fileUpload(file, res.data.key, res.data.url)
      }
      catch (error) {
        console.log(error)
        return this.errorHandler(error)
      }
    },

    /**
     * 動画の更新関数
     * @param id  動画ID
     * @returns
     */
    async deleteVideo(id: string): Promise<void> {
      try {
        await apiClient.videoApi().v1DeleteVideo(id)
      }
      catch (error) {
        console.log(error)
        return this.errorHandler(error)
      }
    },
  },
})
