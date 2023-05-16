<script lang="ts" setup>
import dayjs from 'dayjs'

import { CreatePromotionRequest, DiscountType } from '~/types/api'
import { PromotionTime } from '~/types/props'

const props = defineProps({
  formData: {
    type: Object,
    default: (): CreatePromotionRequest => ({
      title: '',
      description: '',
      public: false,
      discountType: DiscountType.AMOUNT,
      discountRate: 0,
      code: '',
      startAt: dayjs().unix(),
      endAt: dayjs().unix()
    })
  },
  timeData: {
    type: Object,
    default: (): PromotionTime => ({
      startDate: '',
      startTime: '',
      endDate: '',
      endTime: ''
    })
  }
})

const emit = defineEmits<{
  (e: 'update:formData', formData: CreatePromotionRequest): void
  (e: 'update:timeData', timeData: PromotionTime): void
  (e: 'submit'): void
}>()

const formDataValue = computed({
  get: (): CreatePromotionRequest => props.formData as CreatePromotionRequest,
  set: (val: CreatePromotionRequest) => emit('update:formData', val)
})

const timeDataValue = computed({
  get: (): PromotionTime => props.timeData as PromotionTime,
  set: (val: PromotionTime) => emit('update:timeData', val)
})

const handleSubmit = () => {
  emit('submit')
}
</script>

<template>
  <div>
    <p class="text-h6">
      セール情報登録
    </p>
    <organisms-promotion-form v-model:form-data="formDataValue" v-model:time-data="timeDataValue" @submit="handleSubmit" />
  </div>
</template>
