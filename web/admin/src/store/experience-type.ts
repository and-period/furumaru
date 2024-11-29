import { ca } from 'vuetify/locale'
import { apiClient } from '~/plugins/api-client'
import type { CreateExperienceTypeRequest, ExperienceType } from '~/types/api'

export const useExperienceTypeStore = defineStore('experienceType', {
  state: () => ({
    experienceType: {} as ExperienceType,
    experienceTypes: [] as ExperienceType[],
    totalItems: 0,
  }),

  actions: {
    /**
     * 体験カテゴリ一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     * @returns
     */
    async fetchExperienceTypes(limit = 20, offset = 0): Promise<void> {
      try {
        const res = await apiClient.experienceTypeApi().v1ListExperienceTypes(limit, offset)

        const experienceTypeStore = useExperienceTypeStore()
        this.experienceTypes = res.data.experienceTypes
        this.totalItems = res.data.total
        experienceTypeStore.experienceTypes = res.data.experienceTypes
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 体験カテゴリを検索する非同期関数
     */
    async searchExperienceTypes(name = '', experienceTypeIds: string[] = []): Promise<void> {
      try {
        const res = await apiClient.experienceTypeApi().v1ListExperienceTypes(undefined, undefined, name)
        const experienceTypes: ExperienceType[] = []
        this.experienceTypes.forEach((experienceType: ExperienceType): void => {
          if (!experienceTypeIds.includes(experienceType.id)) {
            return
          }
          experienceTypes.push(experienceType)
        })
        res.data.experienceTypes.forEach((experienceType: ExperienceType): void => {
          if (experienceTypes.find((v): boolean => v.id === experienceType.id)) {
            return
          }
          experienceTypes.push(experienceType)
        })
        this.experienceTypes = experienceTypes
        this.totalItems = res.data.total
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 体験カテゴリを新規登録する非同期関数
     */
    async createExperienceType(params: CreateExperienceTypeRequest): Promise<void> {
      try {
        await apiClient.experienceTypeApi().v1CreateExperienceType(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります。',
          409: 'この体験カテゴリ名はすでに登録されています。',
        })
      }
    },

    /**
     * 体験カテゴリを更新する非同期関数
     */
    async updateExperienceType(experienceTypeId: string, payload: CreateExperienceTypeRequest): Promise<void> {
      try {
        await apiClient.experienceTypeApi().v1UpdateExperienceType(experienceTypeId, payload)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります。',
          404: 'この体験カテゴリは存在しません。',
          409: 'この体験カテゴリ名はすでに登録されています。',
        })
      }
    },

    /**
     * 体験カテゴリを削除する非同期関数
     */
    async deleteExperienceType(experienceTypeId: string): Promise<void> {
      try {
        await apiClient.experienceTypeApi().v1DeleteExperienceType(experienceTypeId)
      }
      catch (err) {
        return this.errorHandler(err, {
          404: 'この体験カテゴリは存在しません。',
        })
      }
    },
  },
})
