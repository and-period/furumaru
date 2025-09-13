<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useAlert } from '~/lib/hooks'
import { useCommonStore, useShippingStore } from '~/store'
import type { UpdateDefaultShippingRequest } from '~/types/api/v1'

const commonStore = useCommonStore()
const shippingStore = useShippingStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const { shipping } = storeToRefs(shippingStore)

const loading = ref<boolean>(false)
const formData = ref<UpdateDefaultShippingRequest>({
  box60Rates: [
    {
      name: '',
      price: 0,
      prefectureCodes: [],
    },
  ],
  box60Frozen: 0,
  box80Rates: [
    {
      name: '',
      price: 0,
      prefectureCodes: [],
    },
  ],
  box80Frozen: 0,
  box100Rates: [
    {
      name: '',
      price: 0,
      prefectureCodes: [],
    },
  ],
  box100Frozen: 0,
  hasFreeShipping: false,
  freeShippingRates: 0,
})

const fetchState = useAsyncData(async (): Promise<void> => {
  try {
    await shippingStore.fetchDefaultShipping()
    formData.value = { ...shipping.value }
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
})

const isLoading = (): boolean => {
  return fetchState?.pending?.value || loading.value
}

const handleSubmit = async (): Promise<void> => {
  try {
    loading.value = true
    await shippingStore.updateDefaultShipping(formData.value)
    commonStore.addSnackbar({
      color: 'info',
      message: 'デフォルト配送設定を更新しました。',
    })
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
  <templates-shipping-edit
    v-model:form-data="formData"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :shipping="shipping"
    @submit="handleSubmit"
  />
</template>
