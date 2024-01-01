import { useAuthStore } from './auth'
import type { CheckoutRequest, CheckoutStateResponse } from '~/types/api'

export const useCheckoutStore = defineStore('checkout', {
  state: () => {
    return {
      checkoutState: {
        isLoading: false,
      },
      checkTransactionStatusState: {
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

    async checkTransactionStatus(
      sessionId: string,
    ): Promise<CheckoutStateResponse> {
      this.checkTransactionStatusState.isLoading = true
      try {
        const authStore = useAuthStore()
        const res = await this.checkoutApiClient(
          authStore.accessToken,
        ).v1GetCheckoutState({ transactionId: sessionId })
        return res
      } catch (error) {
        return this.errorHandler(error)
      } finally {
        this.checkTransactionStatusState.isLoading = false
      }
    },
  },
})
