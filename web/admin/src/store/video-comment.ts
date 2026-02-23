import { useApiClient } from '~/composables/useApiClient'
import { VideoCommentApi } from '~/types/api/v1'
import type { V1VideosVideoIdCommentsCommentIdPatchRequest, V1VideosVideoIdCommentsGetRequest, VideoComment } from '~/types/api/v1'

export const useVideoCommentStore = defineStore('video-comment', () => {
  const { create, errorHandler } = useApiClient()
  const videoCommentApi = () => create(VideoCommentApi)

  const comments = ref<VideoComment[]>([])
  const nextToken = ref<string | null>(null)

  async function fetchComments(videoId: string, limit = 20, token: string | null = null): Promise<void> {
    try {
      const params: V1VideosVideoIdCommentsGetRequest = {
        videoId,
        limit,
        next: token || undefined,
      }
      const res = await videoCommentApi().v1VideosVideoIdCommentsGet(params)
      comments.value = res.comments
      nextToken.value = res.nextToken
    }
    catch (err) {
      console.log(err)
      return errorHandler(err)
    }
  }

  async function fetchAllComments(videoId: string, limit = 100): Promise<void> {
    try {
      let allComments: VideoComment[] = []
      let token: string | null = null

      do {
        const params: V1VideosVideoIdCommentsGetRequest = {
          videoId,
          limit,
          next: token || undefined,
        }
        const res = await videoCommentApi().v1VideosVideoIdCommentsGet(params)
        allComments = [...allComments, ...res.comments]
        token = res.nextToken
      } while (token)

      comments.value = allComments
      nextToken.value = null
    }
    catch (err) {
      console.log(err)
      return errorHandler(err)
    }
  }

  async function disableComment(videoId: string, commentId: string, disabled: boolean): Promise<void> {
    try {
      const params: V1VideosVideoIdCommentsCommentIdPatchRequest = {
        videoId,
        commentId,
        updateVideoCommentRequest: { disabled },
      }
      await videoCommentApi().v1VideosVideoIdCommentsCommentIdPatch(params)
    }
    catch (err) {
      console.log(err)
      return errorHandler(err)
    }
  }

  return {
    comments,
    nextToken,
    fetchComments,
    fetchAllComments,
    disableComment,
  }
})
