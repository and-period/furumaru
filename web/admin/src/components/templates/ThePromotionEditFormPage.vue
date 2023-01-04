<template>
  <div>
    <p class="text-h6">セール情報編集</p>
    <the-promotion-form
      form-type="edit"
      :form-data="formDataValue"
      :time-data="timeDataValue"
      :form-data-loading="formDataLoading"
      @submit="handleSubmit"
    />
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, PropType } from '@vue/composition-api'
import dayjs from 'dayjs'

import { CreatePromotionRequest, DiscountType } from '~/types/api'
import { PromotionTime } from '~/types/props'

export default defineComponent({
  props: {
    formData: {
      type: Object as PropType<CreatePromotionRequest>,
      default: () => {
        return {
          title: '',
          description: '',
          public: false,
          publishedAt: dayjs().unix(),
          discountType: DiscountType.AMOUNT,
          discountRate: 0,
          code: '',
          startAt: dayjs().unix(),
          endAt: dayjs().unix(),
        }
      },
    },
    timeData: {
      type: Object as PropType<PromotionTime>,
      default: () => {
        return {
          publishedDate: '',
          publishedTime: '',
          startDate: '',
          startTime: '',
          endDate: '',
          endTime: '',
        }
      },
    },
    formDataLoading: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const formDataValue = computed({
      get: (): CreatePromotionRequest => props.formData,
      set: (val: CreatePromotionRequest) => emit('update:formData', val),
    })

    const timeDataValue = computed({
      get: (): PromotionTime => props.timeData,
      set: (val: PromotionTime) => emit('update:timeData', val),
    })

    const handleSubmit = () => {
      emit('submit')
    }

    return {
      formDataValue,
      timeDataValue,
      handleSubmit,
    }
  },
})
</script>
