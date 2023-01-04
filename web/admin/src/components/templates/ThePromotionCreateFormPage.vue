<template>
  <div>
    <p class="text-h6">セール情報登録</p>
    <the-promotion-form :form-data="formDataValue" @submit="handleSubmit" />
  </div>
</template>

<script lang="ts">
import { computed, PropType } from '@nuxtjs/composition-api'
import { defineComponent } from '@vue/composition-api'
import dayjs from 'dayjs'

import { CreatePromotionRequest, DiscountType } from '~/types/api'

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
  },
  setup(props, { emit }) {
    const formDataValue = computed({
      get: (): CreatePromotionRequest => props.formData,
      set: (val: CreatePromotionRequest) => emit('update:formData', val),
    })

    const handleSubmit = () => {
      emit('submit')
    }

    return {
      formDataValue,
      handleSubmit,
    }
  },
})
</script>
