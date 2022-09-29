<template>
  <form @submit.prevent="handleSubmit">
  </form>
</template>

<script lang="ts">
import { computed, PropType } from '@nuxtjs/composition-api';
import { defineComponent } from '@vue/composition-api'
import dayjs from 'dayjs';

import { CreateNotificationRequest } from '~/types/api';

export default defineComponent({
  props: {
    formType: {
      type: String,
      default: 'create',
      validator: (value: string) => {
        return ['create', 'edit'].includes(value)
      },
    },
    formData: {
      type: Object as PropType<CreateNotificationRequest>,
        default: () => {
          return {
            title: '',
            body: '',
            targets: [0],
            public: false,
            publishedAt: dayjs().unix(),
          }
        },
    },
  },

  setup(props, { emit }) {
    const formDataValue = computed({
        get: (): CreateNotificationRequest => props.formData,
        set: (val: CreateNotificationRequest) => emit('update:formData', val)
      })

    const btnText = computed(() => {
      return props.formType === 'create' ? '登録' : '更新'
    })

    const handleSubmit = () => {
      emit('submit')
    }

    return {
      formDataValue,
      btnText,
      handleSubmit,
    }
  },
})
</script>
