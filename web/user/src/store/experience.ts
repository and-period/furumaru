import type { ExperiencesResponse } from '~/types/api'

export const useExperienceStore = defineStore('experience', {
  state: () => {
    return {
      experiencesFetchState: {
        isLoading: false,
      },
      experiencesResponse: {
        experiences: [],
        coordinators: [],
        producers: [],
        experienceTypes: [],
        total: 0,
      } as ExperiencesResponse,
      experienceFetchState: {
        isLoading: false,
      },
    }
  },

  actions: {
    /**
     * 体験一覧取得
     * @returns
     */
    async fetchExperiences() {
      this.experiencesFetchState.isLoading = true
      try {
        const response = await this.experienceApiClient().v1ListExperiences()
        this.experiencesResponse = response
      }
      catch (error) {
        console.error(e)
        return this.errorHandler(error)
      }
      finally {
        this.experiencesFetchState.isLoading = false
      }
    },
  },
})
