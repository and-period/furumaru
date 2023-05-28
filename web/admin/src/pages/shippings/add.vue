<script lang="ts" setup>
import { useAlert } from '~/lib/hooks'
import { useCommonStore, useShippingStore } from '~/store'
import { CreateShippingRequest } from '~/types/api'
import { ApiBaseError } from '~/types/exception'

const router = useRouter()
const commonStore = useCommonStore()
const shippingStore = useShippingStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const formData = reactive<CreateShippingRequest>({
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
    await shippingStore.createShipping(formData)
    commonStore.addSnackbar({
      color: 'info',
      message: `${formData.name}を登録しました。`
    })
    router.push('/shippings')
  } catch (error) {
    if (error instanceof ApiBaseError) {
      show(error.message)
      window.scrollTo({
        top: 0,
        behavior: 'smooth'
      })
    }
  }
}
</script>

<template>
  <templates-shopping-add
    v-model:form-data="formData"
    :is-alrt="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    @submit="handleSubmit"
  />
</template>
