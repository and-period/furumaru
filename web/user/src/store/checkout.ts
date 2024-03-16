import { useAuthStore } from './auth'
import type { CheckoutRequest, CheckoutStateResponse, GuestCheckoutRequest, GuestCheckoutStateResponse } from '~/types/api'

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
    /**
     * チェックアウト処理を行うメソッド
     * @param payload
     * @returns
     */
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

    /**
     * チェックアウト処理を行うメソッド (ゲスト用)
     * @param payload
     * @returns
     */
    async guestCheckout(payload: GuestCheckoutRequest): Promise<string> {
      this.checkoutState.isLoading = true
      try {
        const res = await this.checkoutApiClient().v1GuestCheckout({
          body: payload
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
      } catch (error) {
        return this.errorHandler(error)
      } finally {
        this.checkTransactionStatusState.isLoading = false
      }
    },

  },
})
