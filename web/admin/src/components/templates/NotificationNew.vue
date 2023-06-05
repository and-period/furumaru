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
  <v-card>
    <v-card-title>お知らせ登録</v-card-title>

    <organisms-notification-form :form-data="formDataValue" :time-data="timeDataValue" @submit="handleSubmit" />
  </v-card>
</template>
