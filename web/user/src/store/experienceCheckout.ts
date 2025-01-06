import { useAuthStore } from './auth'
import type { CheckoutExperienceRequest, GuestCheckoutExperienceRequest } from '~/types/api'

export const useExperienceCheckoutStore = defineStore('experience-checkout', {
  state: () => {
    return {
      checkoutState: {
        isLoading: false,
      },
    }
  },

  actions: {
    /**
     * 体験購入（認証済みユーザー）
     * @param id 購入対象の体験ID
     * @param payload 購入情報
     * @returns リダイレクト先URL
     */
    async checkout(id: string, payload: CheckoutExperienceRequest): Promise<string> {
      this.checkoutState.isLoading = true
      try {
        const authStore = useAuthStore()
        const res = await this.checkoutApiClient(authStore.accessToken).v1CheckoutsExperiencesExperienceIdPost({
          experienceId: id,
          body: payload,
        })
        return res.url
      }
      catch (error) {
        return this.errorHandler(error)
      }
      finally {
        this.checkoutState.isLoading = false
      }
    },

    /**
     * 体験購入（ゲストユーザー）
     * @param id 購入対象の体験ID
     * @param payload 購入情報
     * @returns リダイレクト先URL
     */
    async checkoutByGuest(id: string, payload: GuestCheckoutExperienceRequest): Promise<string> {
      this.checkoutState.isLoading = true
      try {
        const res = await this.checkoutApiClient().v1GuestCheckoutExperience({
          experienceId: id,
          body: payload,
        })
        return res.url
      }
      catch (error) {
        return this.errorHandler(error)
      }
      finally {
        this.checkoutState.isLoading = false
      }
    },
  },
})
