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
      await this.checkoutApiClient().v1Checkout({ body: payload })
      this.checkoutState.isLoading = false
    },
  },
})
