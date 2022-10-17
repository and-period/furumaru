<template>
  <div>
    <p class="text-h6">セール情報登録</p>
    <the-promotion-form form-type="edit" :form-data="formDataValue" :form-data-loading="formDataLoading" @submit="handleSubmit" />
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, PropType } from '@vue/composition-api'
import dayjs from 'dayjs'

import { CreatePromotionRequest } from '~/types/api'

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
          discountType: 1,
          discountRate: 0,
          code: '',
          startAt: dayjs().unix(),
          endAt: dayjs().unix(),
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

