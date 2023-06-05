<script lang="ts" setup>
import dayjs, { unix } from 'dayjs'
import { useAlert } from '~/lib/hooks'

import { usePromotionStore } from '~/store'
import { UpdatePromotionRequest } from '~/types/api'
import { PromotionTime } from '~/types/props'

const router = useRouter()
const route = useRoute()
const id = route.params.id as string

const { getPromotion, updatePromotion } = usePromotionStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const formData = reactive<UpdatePromotionRequest>({
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

const fetchState = useAsyncData(async () => {
  const promotion = await getPromotion(id)
  formData.title = promotion.title
  formData.description = promotion.description
  formData.public = promotion.public
  formData.discountType = promotion.discountType
  formData.discountRate = promotion.discountRate
  formData.code = promotion.code
  timeData.startDate = unix(promotion.startAt).format('YYYY-MM-DD')
  timeData.startTime = unix(promotion.startAt).format('HH:mm')
  timeData.endDate = unix(promotion.endAt).format('YYYY-MM-DD')
  timeData.endTime = unix(promotion.endAt).format('HH:mm')
})

const isLoading = (): boolean => {
  return fetchState?.pending?.value || false
}

const handleSubmit = async () => {
  try {
    await updatePromotion(id, {
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

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-promotion-edit
    v-model:form-data="formData"
    v-model:time-data="timeData"
    :is-alrt="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :form-data-loading="isLoading()"
    @submit="handleSubmit"
  />
</template>
