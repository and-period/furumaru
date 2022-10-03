<template>
  <div>
    <p class="text-h6">生産者登録</p>
    <the-notification-form :form-data="formDataValue" @submit="handleSubmit" />
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, PropType } from '@nuxtjs/composition-api'
import dayjs from 'dayjs'

import TheNotificationForm from '../organisms/TheNotificationForm.vue'

import { CreateNotificationRequest } from '~/types/api'

export default defineComponent({
  components: { TheNotificationForm },
  props: {
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
      set: (val: CreateNotificationRequest) => emit('update:formData', val),
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
