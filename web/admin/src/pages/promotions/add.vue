<script lang="ts" setup>
import dayjs from 'dayjs'

import { usePromotionStore } from '~/store'
import { CreatePromotionRequest } from '~/types/api'

const router = useRouter()
const promotionStore = usePromotionStore()

const formData = reactive<CreatePromotionRequest>({
  title: '',
  description: '',
  public: false,
  publishedAt: dayjs().unix(),
  discountType: 1,
  discountRate: 0,
  code: '',
  startAt: dayjs().unix(),
  endAt: dayjs().unix()
})

const handleSubmit = async () => {
  try {
    await promotionStore.createPromotion(formData)
    router.push('/promotions')
  } catch (error) {
    console.log(error)
  }
}
</script>

<template>
  <the-promotion-create-form-page
    :form-data="formData"
    @submit="handleSubmit"
  />
</template>
