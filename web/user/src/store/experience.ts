import type { ExperiencesResponse } from '~/types/api'

export const useExperienceStore = defineStore('experience', {
  state: () => {
    return {
      // 体験一覧取得ステータス
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
      // 体験詳細取得ステータス
      experienceFetchState: {
        isLoading: false,
      },
    }
  },

  getters: {
    // 体験一覧
    experiences: (state) => {
      return state.experiencesResponse.experiences.map((experience) => {
        return {
          ...experience,
          coordinator: state.experiencesResponse.coordinators.find((coordinator) => {
            return coordinator.id === experience.coordinatorId
          }),
          producer: state.experiencesResponse.producers.find((producer) => {
            return producer.id === experience.producerId
          }),
          experienceType: state.experiencesResponse.experienceTypes.find((experienceType) => {
            return experienceType.id === experience.experienceTypeId
          }),
        }
      })
    },
  },

  actions: {
    /**
     * 体験一覧取得
     * @returns
     */
    async fetchExperiences(
      longitude: number,
      latitude: number,
    ) {
      this.experiencesFetchState.isLoading = true
      try {
        const response = await this.experienceApiClient().v1ListExperiencesByGeolocation(
          {
            longitude,
            latitude,
          },
        )
        this.experiencesResponse = response
      }
      catch (error) {
        return this.errorHandler(error)
      }
      finally {
        this.experiencesFetchState.isLoading = false
      }
    },
  },
})
