import { defineStore } from 'pinia'
import { apiClient } from '~/plugins/api-client'
import type { PaymentMethodType, PaymentSystemStatus, PaymentSystem, UpdatePaymentSystemRequest } from '~/types/api'

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
        const res = await apiClient.paymentSystemApi().v1ListPaymentSystems()
        this.systems = res.data.systems
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    async updatePaymentStatus(methodType: PaymentMethodType, status: PaymentSystemStatus) {
      try {
        const req: UpdatePaymentSystemRequest = { status }
        await apiClient.paymentSystemApi().v1UpdatePaymentSystem(methodType, req)
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },
  },
})
