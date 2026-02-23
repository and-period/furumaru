import { useApiClient } from '~/composables/useApiClient'
import { PaymentSystemApi } from '~/types/api/v1'
import type { PaymentMethodType, PaymentSystem, PaymentSystemStatus, V1PaymentSystemsMethodTypePatchRequest } from '~/types/api/v1'

export const usePaymentSystemStore = defineStore('paymentSystem', () => {
  const { create, errorHandler } = useApiClient()
  const paymentSystemApi = () => create(PaymentSystemApi)

  const systems = ref<PaymentSystem[]>([])

  async function fetchPaymentSystems(): Promise<void> {
    try {
      const res = await paymentSystemApi().v1PaymentSystemsGet()
      systems.value = res.systems
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function updatePaymentStatus(methodType: PaymentMethodType, status: PaymentSystemStatus) {
    try {
      const params: V1PaymentSystemsMethodTypePatchRequest = {
        methodType,
        updatePaymentSystemRequest: { status },
      }
      await paymentSystemApi().v1PaymentSystemsMethodTypePatch(params)
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  return {
    systems,
    fetchPaymentSystems,
    updatePaymentStatus,
  }
})
