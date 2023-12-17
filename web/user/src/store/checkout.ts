import { useAuthStore } from './auth'
import type { CheckoutRequest } from '~/types/api'

export const useCheckoutStore = defineStore('checkout', {
  state: () => {
    return {
      checkoutState: {
        isLoading: false,
      },
    }
  },

  actions: {
    async checkout(payload: CheckoutRequest): Promise<string> {
      this.checkoutState.isLoading = true
      try {
        const authStore = useAuthStore()
        const res = await this.checkoutApiClient(
          authStore.accessToken,
        ).v1Checkout({
          body: payload,
        })
        return res.url
      } catch (error) {
        return this.errorHandler(error)
      } finally {
        this.checkoutState.isLoading = false
      }
    },
  },
})
