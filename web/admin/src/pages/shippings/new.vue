<script lang="ts" setup>
import { useAlert } from '~/lib/hooks'
import { useCommonStore, useShippingStore } from '~/store'
import { CreateShippingRequest } from '~/types/api'

const router = useRouter()
const commonStore = useCommonStore()
const shippingStore = useShippingStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const loading = ref<boolean>(false)
const formData = ref<CreateShippingRequest>({
  name: '',
  box60Rates: [
    {
      name: '',
      price: 0,
      prefectures: []
    }
  ],
  box60Refrigerated: 0,
  box60Frozen: 0,
  box80Rates: [
    {
      name: '',
      price: 0,
      prefectures: []
    }
  ],
  box80Refrigerated: 0,
  box80Frozen: 0,
  box100Rates: [
    {
      name: '',
      price: 0,
      prefectures: []
    }
  ],
  box100Refrigerated: 0,
  box100Frozen: 0,
  hasFreeShipping: false,
  freeShippingRates: 0
})

const handleSubmit = async (): Promise<void> => {
  try {
    loading.value = true
    await shippingStore.createShipping(formData.value)
    commonStore.addSnackbar({
      color: 'info',
      message: `${formData.value.name}を登録しました。`
    })
    router.push('/shippings')
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)

    window.scrollTo({
      top: 0,
      behavior: 'smooth'
    })
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <templates-shipping-new
    v-model:form-data="formData"
    :loading="loading"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    @submit="handleSubmit"
  />
</template>
