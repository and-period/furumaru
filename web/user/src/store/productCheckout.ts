import { useAuthStore } from './auth'
import type { CheckoutProductRequest, CheckoutStateResponse, GuestCheckoutProductRequest, GuestCheckoutStateResponse } from '~/types/api'

export const useProductCheckoutStore = defineStore('product-checkout', {
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
    /**
     * チェックアウト処理を行うメソッド
     * @param payload
     * @returns
     */
    async checkout(payload: CheckoutProductRequest): Promise<string> {
      this.checkoutState.isLoading = true
      try {
        const authStore = useAuthStore()
        const res = await this.checkoutApiClient(
          authStore.accessToken,
        ).v1CheckoutProduct({
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
     * チェックアウト処理を行うメソッド (ゲスト用)
     * @param payload
     * @returns
     */
    async guestCheckout(payload: GuestCheckoutProductRequest): Promise<string> {
      this.checkoutState.isLoading = true
      try {
        const res = await this.checkoutApiClient().v1GuestCheckoutProduct({
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
      }
      catch (error) {
        return this.errorHandler(error)
      }
      finally {
        this.checkTransactionStatusState.isLoading = false
      }
    },

    /**
     * 注文情報の取得を行うメソッド (ゲスト用)
     * @param payload
     * @returns
     */
    async guestCheckTransactionStatus(
      sessionId: string,
    ): Promise<GuestCheckoutStateResponse> {
      this.checkTransactionStatusState.isLoading = true
      try {
        const res = await this.checkoutApiClient().v1GetGuestCheckoutState({ transactionId: sessionId })
        return res
      }
      catch (error) {
        return this.errorHandler(error)
      }
      finally {
        this.checkTransactionStatusState.isLoading = false
      }
    },

  },
})
