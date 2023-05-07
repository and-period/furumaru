<script lang="ts" setup>
import dayjs from 'dayjs'

import { usePromotionStore } from '~/store'
import { UpdatePromotionRequest } from '~/types/api'
import { PromotionTime } from '~/types/props'

const router = useRouter()
const route = useRoute()
const id = route.params.id as string

const { getPromotion, editPromotion } = usePromotionStore()

const formData = reactive<UpdatePromotionRequest>({
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

const timeData = reactive<PromotionTime>({
  publishedDate: '',
  publishedTime: '',
  startDate: '',
  startTime: '',
  endDate: '',
  endTime: ''
})

const fetchState = useAsyncData(async () => {
  const promotion = await getPromotion(id)
  formData.title = promotion.title
  formData.description = promotion.description
  formData.public = promotion.public
  formData.discountType = promotion.discountType
  formData.discountRate = promotion.discountRate
  formData.code = promotion.code
  timeData.publishedDate = dayjs
    .unix(promotion.publishedAt)
    .format('YYYY-MM-DD')
  timeData.publishedTime = dayjs.unix(promotion.publishedAt).format('HH:mm')
  timeData.startDate = dayjs.unix(promotion.startAt).format('YYYY-MM-DD')
  timeData.startTime = dayjs.unix(promotion.startAt).format('HH:mm')
  timeData.endDate = dayjs.unix(promotion.endAt).format('YYYY-MM-DD')
  timeData.endTime = dayjs.unix(promotion.endAt).format('HH:mm')
})

const isLoading = (): boolean => {
  return fetchState?.pending?.value || false
}

const handleSubmit = async () => {
  try {
    await editPromotion(id, {
      ...formData,
      discountRate: Number(formData.discountRate)
    })
    router.push('/promotions')
  } catch (error) {
    console.log(error)
  }
}

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-promotion-edit-form-page
    v-model:form-data="formData"
    v-model:time-data="timeData"
    :form-data-loading="isLoading()"
    @submit="handleSubmit"
  />
</template>
