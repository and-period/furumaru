<script lang="ts" setup>
import { storeToRefs } from 'pinia'

import { useAlert } from '~/lib/hooks'
import { usePaymentSystemStore } from '~/store'
import { PaymentSystemStatus } from '~/types/api/v1'
import type { PaymentMethodType, PaymentSystem } from '~/types/api/v1'

const paymentSystemStore = usePaymentSystemStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const { systems } = storeToRefs(paymentSystemStore)

const loading = ref<boolean>(false)

const fetchState = useAsyncData('payment-systems', async (): Promise<void> => {
  try {
    await paymentSystemStore.fetchPaymentSystems()
  }
  catch (err) {
    if (err instanceof Error) {
      showError(err.message)
    }
    console.log(err)
  }
})

const isLoading = (): boolean => {
  return fetchState?.pending?.value || loading.value
}

const handleUpdateStatus = async (methodType: PaymentMethodType): Promise<void> => {
  let status: PaymentSystemStatus
  const system = systems.value.find((system: PaymentSystem): boolean => {
    return system.methodType === methodType
  })
  if (!system) {
    return
  }

  switch (system.status) {
    case PaymentSystemStatus.PaymentSystemStatusInUse:
      status = PaymentSystemStatus.PaymentSystemStatusOutage
      break
    case PaymentSystemStatus.PaymentSystemStatusOutage:
      status = PaymentSystemStatus.PaymentSystemStatusInUse
      break
    default:
      status = PaymentSystemStatus.PaymentSystemStatusUnknown
  }

  try {
    loading.value = true
    await paymentSystemStore.updatePaymentStatus(methodType, status)
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)

    window.scrollTo({
      top: 0,
      behavior: 'smooth',
    })
  }
  finally {
    loading.value = false
    fetchState.refresh()
  }
}

try {
  await fetchState.execute()
}
catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-payment-system-list
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :systems="systems"
    @submit="handleUpdateStatus"
  />
</template>
