<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useAlert } from '~/lib/hooks'
import { useAuthStore, useCommonStore, useShippingStore } from '~/store'
import type { UpsertShippingRequest } from '~/types/api'

const route = useRoute()

const coordinatorId = route.params.id as string

const authStore = useAuthStore()
const commonStore = useCommonStore()
const shippingStore = useShippingStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const { adminId } = storeToRefs(authStore)

const loading = ref<boolean>(false)
const formData = ref<UpsertShippingRequest>({
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

const { data, status, error } = useAsyncData(async () => {
  return await shippingStore.fetchShipping(adminId.value, coordinatorId)
})

watch(error, (newError) => {
  if (newError) {
    if (newError instanceof Error) {
      show(newError.message)
    }
    console.log(newError)
  }
})

watch(data, (newData) => {
  formData.value = { ...newData }
})

const isLoading = (): boolean => {
  return status.value === 'pending'
}

const handleSubmit = async (): Promise<void> => {
  try {
    loading.value = true
    await shippingStore.upsertShipping(adminId.value, formData.value)
    commonStore.addSnackbar({
      color: 'info',
      message: '配送設定を更新しました。',
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
</script>

<template>
  <templates-shipping-edit
    v-model:form-data="formData"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :shipping="data"
    @submit="handleSubmit"
  />
</template>
