import type { V1VideosVideoIdCommentsCommentIdPatchRequest, V1VideosVideoIdCommentsGetRequest, VideoComment } from "~/types/api/v1";

export const useVideoCommentStore = defineStore('video-comment', {
  state: () => ({
    comments: [] as VideoComment[],
    nextToken: null as string | null,
  }),

  actions: {
    async fetchComments(videoId: string, limit = 20, nextToken: string | null = null): Promise<void> {
      try {
        const params: V1VideosVideoIdCommentsGetRequest = {
          videoId,
          limit,
          next: nextToken || undefined,
        }
        const res = await this.videoCommentApi().v1VideosVideoIdCommentsGet(params)
        this.comments = res.comments
        this.nextToken = res.nextToken
      }
      catch (err) {
        console.log(err)
        return this.errorHandler(err)
      }
    },

    async fetchAllComments(videoId: string, limit = 100): Promise<void> {
      try {
        let allComments: VideoComment[] = []
        let nextToken: string | null = null
        
        do {
          const params: V1VideosVideoIdCommentsGetRequest = {
            videoId,
            limit,
            next: nextToken || undefined,
          }
          const res = await this.videoCommentApi().v1VideosVideoIdCommentsGet(params)
          allComments = [...allComments, ...res.comments]
          nextToken = res.nextToken
        } while (nextToken)
        
        this.comments = allComments
        this.nextToken = null
      }
      catch (err) {
        console.log(err)
        return this.errorHandler(err)
      }
    },

    async disableComment(videoId: string, commentId: string, disabled: boolean): Promise<void> {
      try {
        const params: V1VideosVideoIdCommentsCommentIdPatchRequest = {
          videoId,
          commentId,
          updateVideoCommentRequest: {
            disabled,
          },
        }
        await this.videoCommentApi().v1VideosVideoIdCommentsCommentIdPatch(params)
      }
      catch (err) {
        console.log(err)
        return this.errorHandler(err)
      }
    },
  },
})
