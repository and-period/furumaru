<script lang="ts" setup>
import dayjs from 'dayjs'
import { useAlert } from '~/lib/hooks'

import { usePromotionStore } from '~/store'
import { CreatePromotionRequest } from '~/types/api'
import { PromotionTime } from '~/types/props'

const router = useRouter()
const promotionStore = usePromotionStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const formData = reactive<CreatePromotionRequest>({
  title: '',
  description: '',
  public: false,
  discountType: 1,
  discountRate: 0,
  code: '',
  startAt: dayjs().unix(),
  endAt: dayjs().unix()
})

const timeData = reactive<PromotionTime>({
  startDate: '',
  startTime: '',
  endDate: '',
  endTime: ''
})

const handleSubmit = async () => {
  try {
    await promotionStore.createPromotion({
      ...formData,
      discountRate: Number(formData.discountRate)
    })
    router.push('/promotions')
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}
</script>

<template>
  <templates-promotion-new
    v-model:form-data="formData"
    v-model:time-data="timeData"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    @submit="handleSubmit"
  />
</template>
