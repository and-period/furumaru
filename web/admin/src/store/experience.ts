import type { AxiosResponse } from 'axios'
import { fileUpload } from './helper'
import { apiClient } from '~/plugins/api-client'
import type {
  CreateExperienceRequest,
  Experience,
  ExperiencesResponse,
  GetUploadUrlRequest,
  UpdateExperienceRequest,
  UploadUrlResponse,
} from '~/types/api'

export const useExperienceStore = defineStore('experience', {
  state: () => ({
    experience: {} as Experience,
    experiences: [] as Experience[],
    experiencesResponse: null as ExperiencesResponse | null,
    totalItems: 0,
  }),

  actions: {
    /**
     * 体験一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     * @returns
     */
    async fetchExperiences(limit = 20, offset = 0): Promise<void> {
      try {
        const res = await apiClient
          .experienceApi()
          .v1ListExperiences(limit, offset)

        const experienceStore = useExperienceStore()
        this.experiencesResponse = res.data
        this.totalItems = res.data.total
        experienceStore.experiences = res.data.experiences
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 体験詳細を取得する非同期関数
     * @param experienceId 取得対象の体験ID
     * @returns
     */
    async fetchExperience(experienceId: string) {
      try {
        const res = await apiClient
          .experienceApi()
          .v1GetExperience(experienceId)
        return res.data
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    async createExperience(payload: CreateExperienceRequest) {
      try {
        const res = await apiClient.experienceApi().v1CreateExperience(payload)
        return res.data
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、入力内容に誤りがあります。',
          403: '登録する権限がありません。管理者にお問い合わせください',
        })
      }
    },

    /**
     * 体験を更新する非同期関数
     * @param experienceId 更新対象の体験ID
     * @param payload 更新内容
     * @returns
     */
    async updateExperience(
      experienceId: string,
      payload: UpdateExperienceRequest,
    ) {
      try {
        const res = await apiClient
          .experienceApi()
          .v1UpdateExperience(experienceId, payload)
        return res.data
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 体験メディアファイルをアップロードする非同期関数
     * @param payload
     * @returns
     */
    async uploadExperienceMedia(payload: File): Promise<string> {
      const contentType = payload.type
      try {
        const res = await this.getExperienceMediaUploadUrl(contentType)

        return await fileUpload(payload, res.data.key, res.data.url)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: 'このファイルはアップロードできません。',
        })
      }
    },

    async getExperienceMediaUploadUrl(
      contentType: string,
    ): Promise<AxiosResponse<UploadUrlResponse, any>> {
      const body: GetUploadUrlRequest = {
        fileType: contentType,
      }
      if (contentType.includes('image/')) {
        try {
          const res = await apiClient
            .productApi()
            .v1GetProductImageUploadUrl(body)
          return res
        }
        catch (err) {
          return this.errorHandler(err, {
            400: 'このファイルはアップロードできません。',
          })
        }
      }
      if (contentType.includes('video/')) {
        try {
          const res = await apiClient
            .productApi()
            .v1GetProductVideoUploadUrl(body)
          return res
        }
        catch (err) {
          return this.errorHandler(err, {
            400: 'このファイルはアップロードできません。',
          })
        }
      }
      throw new Error('不明なMINEタイプです。')
    },
    /**
     * 体験を削除する関数
     * @param experienceId
     * @returns
     */
    async deleteExperience(experienceId: string) {
      try {
        await apiClient.experienceApi().v1DeleteExperience(experienceId)
        const index = this.experiences.findIndex(
          experience => experience.id === experienceId,
        )
        this.experiences.splice(index, 1)
        this.totalItems--
      }
      catch (err) {
        return this.errorHandler(err, {
          403: '体験を削除する権限がありません',
          404: '対象の商品が存在しません',
        })
      }
    },
  },
})
