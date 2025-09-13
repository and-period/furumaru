import { fileUpload } from './helper'
import type {
  CreateExperienceRequest,
  Experience,
  ExperiencesResponse,
  GetUploadURLRequest,
  UpdateExperienceRequest,
  UploadURLResponse,
  V1ExperiencesExperienceIdDeleteRequest,
  V1ExperiencesExperienceIdGetRequest,
  V1ExperiencesExperienceIdPatchRequest,
  V1ExperiencesGetRequest,
  V1ExperiencesPostRequest,
  V1UploadProductsImagePostRequest,
  V1UploadProductsVideoPostRequest,
} from '~/types/api/v1'

export const useExperienceStore = defineStore('experience', {
  state: () => ({
    experience: {} as Experience,
    experiences: [] as Experience[],
    experiencesResponse: null as ExperiencesResponse | null,
    totalItems: 0,
  }),

  actions: {
    /**
     * 体験検索関数
     */
    async searchExperiences(
      name: string = '',
      producerId: string = '',
      experienceId: string[] = [],
    ) {
      try {
        const params: V1ExperiencesGetRequest = {
          name,
          producerId,
        }
        const res = await this.experienceApi().v1ExperiencesGet(params)
        const experiences: Experience[] = []
        this.experiences.forEach((experience) => {
          if (experienceId.includes(experience.id)) {
            experiences.push(experience)
          }
        })
        res.experiences.forEach((experience) => {
          if (experiences.find(e => e.id === experience.id)) {
            return
          }
          experiences.push(experience)
        })
        this.experiences = experiences
        this.totalItems = res.total
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 体験一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     * @returns
     */
    async fetchExperiences(limit = 20, offset = 0): Promise<void> {
      try {
        const params: V1ExperiencesGetRequest = {
          limit,
          offset,
        }
        const res = await this.experienceApi().v1ExperiencesGet(params)

        const experienceStore = useExperienceStore()
        this.experiencesResponse = res
        this.totalItems = res.total
        experienceStore.experiences = res.experiences
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
        const params: V1ExperiencesExperienceIdGetRequest = {
          experienceId,
        }
        const res = await this.experienceApi().v1ExperiencesExperienceIdGet(params)
        return res
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    async createExperience(payload: CreateExperienceRequest) {
      try {
        const params: V1ExperiencesPostRequest = {
          createExperienceRequest: payload,
        }
        const res = await this.experienceApi().v1ExperiencesPost(params)
        return res
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
        const params: V1ExperiencesExperienceIdPatchRequest = {
          experienceId,
          updateExperienceRequest: payload,
        }
        const res = await this.experienceApi().v1ExperiencesExperienceIdPatch(params)
        return res
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

        return await fileUpload(this.uploadApi(), payload, res.key, res.url)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: 'このファイルはアップロードできません。',
        })
      }
    },

    async getExperienceMediaUploadUrl(contentType: string): Promise<UploadURLResponse> {
      const body: GetUploadURLRequest = {
        fileType: contentType,
      }
      if (contentType.includes('image/')) {
        try {
          const params: V1UploadProductsImagePostRequest = {
            getUploadURLRequest: body,
          }
          const res = await this.uploadApi().v1UploadProductsImagePost(params)
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
          const params: V1UploadProductsVideoPostRequest = {
            getUploadURLRequest: body,
          }
          const res = await this.uploadApi().v1UploadProductsVideoPost(params)
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
        const params: V1ExperiencesExperienceIdDeleteRequest = {
          experienceId,
        }
        await this.experienceApi().v1ExperiencesExperienceIdDelete(params)
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
