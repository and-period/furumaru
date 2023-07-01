<script lang="ts" setup>
import dayjs from 'dayjs'
import { useAlert } from '~/lib/hooks'

import { usePromotionStore } from '~/store'
import { CreatePromotionRequest } from '~/types/api'

const router = useRouter()
const promotionStore = usePromotionStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const loading = ref<boolean>(false)
const formData = ref<CreatePromotionRequest>({
  title: '',
  description: '',
  public: false,
  discountType: 1,
  discountRate: 0,
  code: '',
  startAt: dayjs().unix(),
  endAt: dayjs().unix()
})

const handleSubmit = async (): Promise<void> => {
  try {
    loading.value = true
    const req: CreatePromotionRequest = {
      ...formData.value,
      discountRate: Number(formData.value.discountRate)
    }
    await promotionStore.createPromotion(req)
    router.push('/promotions')
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <templates-promotion-new
    v-model:form-data="formData"
    :loading="loading"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    @submit="handleSubmit"
  />
</template>
