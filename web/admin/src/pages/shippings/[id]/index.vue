<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useAlert } from '~/lib/hooks'
import { useAuthStore, useCommonStore, useShippingStore } from '~/store'
import type { UpdateShippingRequest } from '~/types/api'

const router = useRouter()
const route = useRoute()

const shippingId = route.params.id as string

const authStore = useAuthStore()
const commonStore = useCommonStore()
const shippingStore = useShippingStore()
const { alertType, isShow, alertText, show, hide } = useAlert('error')

const { adminId } = storeToRefs(authStore)

const formData = ref<UpdateShippingRequest>({
  name: '',
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
  return await shippingStore.fetchShipping(adminId.value, shippingId)
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

const isLoading = computed(() => {
  return status.value === 'pending'
})

const submitting = ref<boolean>(false)

const handleSubmit = async (): Promise<void> => {
  try {
    hide()
    submitting.value = true
    await shippingStore.updateShipping(adminId.value, shippingId, formData.value)
    commonStore.addSnackbar({
      color: 'info',
      message: '配送設定を更新しました。',
    })
    router.push('/shippings')
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
    submitting.value = false
  }
}
</script>

<template>
  <div>
    <v-alert
      v-show="isShow"
      :type="alertType"
      v-text="alertText"
    />

    <v-card-title>
      配送情報詳細
    </v-card-title>
    <organisms-shipping-form
      v-model="formData"
      :loading=" status === 'pending' "
      form-type="update"
      :submitting="submitting"
      @submit="handleSubmit"
    />
  </div>
</template>
