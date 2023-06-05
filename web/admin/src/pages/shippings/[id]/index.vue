<script lang="ts" setup>
import { useAlert } from '~/lib/hooks'
import { useCommonStore, useShippingStore } from '~/store'
import { UpdateShippingRequest } from '~/types/api'

const route = useRoute()
const router = useRouter()
const commonStore = useCommonStore()
const shippingStore = useShippingStore()
const id = route.params.id as string
const { alertType, isShow, alertText, show } = useAlert('error')

const formData = ref<UpdateShippingRequest>({
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

const fetchState = useAsyncData(async () => {
  try {
    const shipping = await shippingStore.getShipping(id)
    formData.value = { ...shipping }
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
})

const isLoading = (): boolean => {
  return fetchState?.pending?.value || false
}

const handleSubmit = async () => {
  try {
    await shippingStore.updateShipping(id, formData.value)
    commonStore.addSnackbar({
      color: 'info',
      message: `${formData.value.name}を更新しました。`
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
  }
}

try {
  await fetchState.execute()
} catch (err) {
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
    @submit="handleSubmit"
  />
</template>
