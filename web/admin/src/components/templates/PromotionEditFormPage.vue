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
      publishedAt: dayjs().unix(),
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
      publishedDate: '',
      publishedTime: '',
      startDate: '',
      startTime: '',
      endDate: '',
      endTime: ''
    })
  },
  formDataLoading: {
    type: Boolean,
    default: false
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
      セール情報編集
    </p>
    <organisms-promotion-form
      form-type="edit"
      :form-data="formDataValue"
      :time-data="timeDataValue"
      :form-data-loading="props.formDataLoading"
      @submit="handleSubmit"
    />
  </div>
</template>
