<script setup lang="ts">
import type { CreateShippingRequest } from '~/types/api'
import { useAuthStore } from '~/store/auth'
import { useShippingStore } from '~/store'
import { useAlert } from '~/lib/hooks'

const router = useRouter()

const authStore = useAuthStore()
const { adminId } = storeToRefs(authStore)

const shippingStore = useShippingStore()
const { createShipping } = shippingStore

const { alertType, isShow, alertText, show, hide } = useAlert('error')

const formData = ref<CreateShippingRequest>({
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

const submitting = ref<boolean>(false)

const handleSubmit = async () => {
  hide()
  submitting.value = true

  try {
    await createShipping(adminId.value, formData.value)
    router.push('/shippings')
  }
  catch (error) {
    // エラーハンドリング
    console.error(error)
    if (error instanceof Error) {
      show(error.message)
    }
    else {
      show('予期しないエラーが発生しました。')
    }

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
      配送情報新規作成
    </v-card-title>
    <organisms-shipping-form
      v-model="formData"
      :loading="false"
      form-type="create"
      :submitting="submitting"
      @submit="handleSubmit"
    />
  </div>
</template>
