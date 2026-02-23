import { useApiClient } from '~/composables/useApiClient'
import { fileUpload } from './helper'
import { UploadApi, VideoApi } from '~/types/api/v1'
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

export const useVideoStore = defineStore('video', () => {
  const { create, errorHandler } = useApiClient()
  const videoApi = () => create(VideoApi)
  const uploadApi = () => create(UploadApi)

  const video = ref<Video | null>(null)
  const products = ref<Product[]>([])
  const experiences = ref<Experience[]>([])
  const viewerLogs = ref<VideoViewerLog[]>([])
  const videoResponse = ref<VideosResponse | null>(null)

  async function fetchVideos(limit = 20, offset = 0): Promise<void> {
    try {
      const params: V1VideosGetRequest = { limit, offset }
      const res = await videoApi().v1VideosGet(params)
      videoResponse.value = res
    }
    catch (error) {
      console.log(error)
      return errorHandler(error)
    }
  }

  async function fetchVideo(id: string): Promise<void> {
    try {
      const params: V1VideosVideoIdGetRequest = { videoId: id }
      const res = await videoApi().v1VideosVideoIdGet(params)
      video.value = res.video
      products.value = res.products
      experiences.value = res.experiences
    }
    catch (error) {
      console.log(error)
      return errorHandler(error)
    }
  }

  async function analyzeVideo(videoId: string, start?: number, end?: number): Promise<void> {
    try {
      const params: V1VideosVideoIdAnalyticsGetRequest = { videoId, start, end }
      const res = await videoApi().v1VideosVideoIdAnalyticsGet(params)
      viewerLogs.value = res.viewerLogs
    }
    catch (err) {
      return errorHandler(err, { 404: '対象の動画が見つかりません。' })
    }
  }

  async function createVideo(payload: CreateVideoRequest): Promise<void> {
    try {
      const params: V1VideosPostRequest = { createVideoRequest: payload }
      await videoApi().v1VideosPost(params)
    }
    catch (error) {
      console.log(error)
      return errorHandler(error)
    }
  }

  async function updateVideo(id: string, payload: UpdateVideoRequest): Promise<void> {
    try {
      const params: V1VideosVideoIdPatchRequest = { videoId: id, updateVideoRequest: payload }
      await videoApi().v1VideosVideoIdPatch(params)
    }
    catch (error) {
      console.log(error)
      return errorHandler(error)
    }
  }

  async function uploadVideoFile(file: File): Promise<string> {
    try {
      const params: V1UploadVideosFilePostRequest = {
        getUploadURLRequest: { fileType: file.type },
      }
      const res = await uploadApi().v1UploadVideosFilePost(params)
      return await fileUpload(uploadApi(), file, res.key, res.url)
    }
    catch (error) {
      console.log(error)
      return errorHandler(error)
    }
  }

  async function uploadThumbnailFile(file: File): Promise<string> {
    try {
      const params: V1UploadVideosThumbnailPostRequest = {
        getUploadURLRequest: { fileType: file.type },
      }
      const res = await uploadApi().v1UploadVideosThumbnailPost(params)
      return await fileUpload(uploadApi(), file, res.key, res.url)
    }
    catch (error) {
      console.log(error)
      return errorHandler(error)
    }
  }

  async function deleteVideo(id: string): Promise<void> {
    try {
      const params: V1VideosVideoIdDeleteRequest = { videoId: id }
      await videoApi().v1VideosVideoIdDelete(params)
    }
    catch (error) {
      console.log(error)
      return errorHandler(error)
    }
  }

  return {
    video,
    products,
    experiences,
    viewerLogs,
    videoResponse,
    fetchVideos,
    fetchVideo,
    analyzeVideo,
    createVideo,
    updateVideo,
    uploadVideoFile,
    uploadThumbnailFile,
    deleteVideo,
  }
})
