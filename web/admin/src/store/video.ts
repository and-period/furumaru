import { apiClient } from "~/plugins/api-client";
import type {
  CreateVideoRequest,
  VideoResponse,
  VideosResponse,
} from "~/types/api";

export const useVideoStore = defineStore("video", {
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
        const res = await apiClient.videoApi().v1ListVideos(limit, offset);
        this.videoResponse = res.data;
      } catch (error) {
        console.log(error);
        return this.errorHandler(error);
      }
    },

    /**
     * 動画詳細取得関数
     * @param id 動画ID
     * @returns
     */
    async fetchVideo(id: string): Promise<VideoResponse> {
      try {
        const res = await apiClient.videoApi().v1GetVideo(id);
        return res.data;
      } catch (error) {
        console.log(error);
        return this.errorHandler(error);
      }
    },

    /**
     * 動画の新規登録関数
     * @param payload
     * @returns
     */
    async createVideo(payload: CreateVideoRequest): Promise<void> {
      try {
        const res = await apiClient.videoApi().v1CreateVideo(payload);
      } catch (error) {
        console.log(error);
        return this.errorHandler(error);
      }
    },

    /**
     * 動画の更新関数
     * @param id  動画ID
     * @param payload
     * @returns
     */
    async updateVideo(id: string, payload: CreateVideoRequest): Promise<void> {
      try {
        await apiClient.videoApi().v1UpdateVideo(id, payload);
      } catch (error) {
        console.log(error);
        return this.errorHandler(error);
      }
    },
  },
});
