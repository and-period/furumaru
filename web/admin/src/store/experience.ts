import { apiClient } from "~/plugins/api-client";
import type { Experience, ExperienceResponse, ExperiencesResponse, Producer } from "~/types/api";

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
        const res = await apiClient.experienceApi().v1ListExperiences(limit, offset)

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
     * 体験を削除する関数
     * @param experienceId
     * @returns
     */
    async deleteExperience(experienceId: string) {
      try {
        await apiClient.experienceApi().v1DeleteExperience(experienceId)
        const index = this.experiences.findIndex(experience => experience.id === experienceId)
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
