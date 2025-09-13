import type { PaymentMethodType, PaymentSystemStatus, PaymentSystem, UpdatePaymentSystemRequest, V1PaymentSystemsMethodTypePatchRequest } from '~/types/api/v1'

export const usePaymentSystemStore = defineStore('paymentSystem', {
  state: () => ({
    systems: [] as PaymentSystem[],
  }),

  actions: {
    /**
     * 決済システム状態一覧を取得する非同期関数
     * @returns
     */
    async fetchPaymentSystems(): Promise<void> {
      try {
        const res = await this.paymentSystemApi().v1PaymentSystemsGet()
        this.systems = res.systems
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    async updatePaymentStatus(methodType: PaymentMethodType, status: PaymentSystemStatus) {
      try {
        const params: V1PaymentSystemsMethodTypePatchRequest = {
          methodType,
          updatePaymentSystemRequest: {
            status,
          },
        }
        await this.paymentSystemApi().v1PaymentSystemsMethodTypePatch(params)
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },
  },
})
