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
    async checkout(payload: CheckoutRequest) {
      this.checkoutState.isLoading = true
      const authStore = useAuthStore()
      const _ = await this.checkoutApiClient(authStore.accessToken).v1Checkout({
        body: payload,
      })
      this.checkoutState.isLoading = false
    },
  },
})
