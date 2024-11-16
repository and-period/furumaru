import { apiClient } from '~/plugins/api-client'
import type { ExperienceType } from '~/types/api'

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
  },
})
