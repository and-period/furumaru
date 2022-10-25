<template>
  <div>
    <p class="text-h6">お知らせ編集</p>
    <the-notification-form
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

import { CreateNotificationRequest } from '~/types/api'
import { NotificationTime } from '~/types/props'

export default defineComponent({
  props: {
    formData: {
      type: Object as PropType<CreateNotificationRequest>,
      default: () => {
        return {
          title: '',
          body: '',
          targets: [],
          public: false,
          publishedAt: dayjs().unix(),
        }
      },
    },
    timeData: {
      type: Object as PropType<NotificationTime>,
      default: () => {
        return {
          publishedDate: '',
          publishedTime: '',
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
      get: (): CreateNotificationRequest => props.formData,
      set: (val: CreateNotificationRequest) => emit('update:formData', val),
    })

    const timeDataValue = computed({
      get: (): NotificationTime => props.timeData,
      set: (val: NotificationTime) => emit('update:timeData', val),
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
