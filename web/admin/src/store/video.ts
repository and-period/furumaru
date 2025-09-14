import { fileUpload } from './helper'
import type {
  CreateVideoRequest,
  Experience,
  Product,
  UpdateVideoRequest,
  V1UploadVideosFilePostRequest,
  V1UploadVideosThumbnailPostRequest,
  V1VideosGetRequest,
  V1VideosPostRequest,
  V1VideosVideoIdAnalyticsGetRequest,
  V1VideosVideoIdDeleteRequest,
  V1VideosVideoIdGetRequest,
  V1VideosVideoIdPatchRequest,
  Video,
  VideosResponse,
  VideoViewerLog,
} from '~/types/api/v1'

export const useVideoStore = defineStore('video', {
  state: () => ({
    video: null as Video | null,
    products: [] as Product[],
    experiences: [] as Experience[],
    viewerLogs: [] as VideoViewerLog[],
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
        const params: V1VideosGetRequest = {
          limit,
          offset,
        }
        const res = await this.videoApi().v1VideosGet(params)
        this.videoResponse = res
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
    async fetchVideo(id: string): Promise<void> {
      try {
        const params: V1VideosVideoIdGetRequest = {
          videoId: id,
        }
        const res = await this.videoApi().v1VideosVideoIdGet(params)
        this.video = res.video
        this.products = res.products
        this.experiences = res.experiences
      }
      catch (error) {
        console.log(error)
        return this.errorHandler(error)
      }
    },

    /**
     * 動画の分析情報を取得する関数
     * @param videoId 動画ID
     * @returns
     */
    async analyzeVideo(videoId: string): Promise<void> {
      try {
        const params: V1VideosVideoIdAnalyticsGetRequest = {
          videoId,
        }
        const res = await this.videoApi().v1VideosVideoIdAnalyticsGet(params)
        this.viewerLogs = res.viewerLogs
      }
      catch (err) {
        return this.errorHandler(err, {
          404: '対象の動画が見つかりません。',
        })
      }
    },

    /**
     * 動画の新規登録関数
     * @param payload
     * @returns
     */
    async createVideo(payload: CreateVideoRequest): Promise<void> {
      try {
        const params: V1VideosPostRequest = {
          createVideoRequest: payload,
        }
        await this.videoApi().v1VideosPost(params)
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
        const params: V1VideosVideoIdPatchRequest = {
          videoId: id,
          updateVideoRequest: payload,
        }
        await this.videoApi().v1VideosVideoIdPatch(params)
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
        const params: V1UploadVideosFilePostRequest = {
          getUploadURLRequest: {
            fileType: file.type,
          },
        }
        const res = await this.uploadApi().v1UploadVideosFilePost(params)

        return await fileUpload(this.uploadApi(), file, res.key, res.url)
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
        const params: V1UploadVideosThumbnailPostRequest = {
          getUploadURLRequest: {
            fileType: file.type,
          },
        }
        const res = await this.uploadApi().v1UploadVideosThumbnailPost(params)

        return await fileUpload(this.uploadApi(), file, res.key, res.url)
      }
      catch (error) {
        console.log(error)
        return this.errorHandler(error)
      }
    },

    /**
     * 動画の削除関数
     * @param id  動画ID
     * @returns
     */
    async deleteVideo(id: string): Promise<void> {
      try {
        const params: V1VideosVideoIdDeleteRequest = {
          videoId: id,
        }
        await this.videoApi().v1VideosVideoIdDelete(params)
      }
      catch (error) {
        console.log(error)
        return this.errorHandler(error)
      }
    },
  },
})
