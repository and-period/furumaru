<script lang="ts" setup>
import dayjs from 'dayjs'

import { CreatePromotionRequest, DiscountType } from '~/types/api'

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
  }
})

const emit = defineEmits<{
  (e: 'update:formData', formData: CreatePromotionRequest): void
  (e: 'submit'): void
}>()

const formDataValue = computed({
  get: (): CreatePromotionRequest => props.formData as CreatePromotionRequest,
  set: (val: CreatePromotionRequest) => emit('update:formData', val)
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
    <the-promotion-form :form-data="formDataValue" @submit="handleSubmit" />
  </div>
</template>
