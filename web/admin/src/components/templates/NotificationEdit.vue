<script lang="ts" setup>
import dayjs from 'dayjs'

import { CreateNotificationRequest } from '~/types/api'
import { NotificationTime } from '~/types/props'

const props = defineProps({
  formData: {
    type: Object,
    default: (): CreateNotificationRequest => ({
      title: '',
      body: '',
      targets: [],
      public: false,
      publishedAt: dayjs().unix()
    })
  },
  timeData: {
    type: Object,
    default: (): NotificationTime => ({
      publishedDate: '',
      publishedTime: ''
    })
  },
  formDataLoading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits<{
  (e: 'update:formData', formData: CreateNotificationRequest): void
  (e: 'update:timeData', timeData: NotificationTime): void
  (e: 'submit'): void
}>()

const formDataValue = computed({
  get: (): CreateNotificationRequest =>
    props.formData as CreateNotificationRequest,
  set: (val: CreateNotificationRequest) => emit('update:formData', val)
})

const timeDataValue = computed({
  get: (): NotificationTime => props.timeData as NotificationTime,
  set: (val: NotificationTime) => emit('update:timeData', val)
})

const handleSubmit = () => {
  emit('submit')
}
</script>

<template>
  <div>
    <p class="text-h6">
      お知らせ編集
    </p>
    <organisms-notification-form
      form-type="edit"
      :form-data="formDataValue"
      :time-data="timeDataValue"
      :form-data-loading="props.formDataLoading"
      @submit="handleSubmit"
    />
  </div>
</template>
