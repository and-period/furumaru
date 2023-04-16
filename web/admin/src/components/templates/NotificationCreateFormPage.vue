<script lang="ts" setup>
import * as dayjs from 'dayjs'

import { CreateNotificationRequest } from '~/types/api'

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
  }
})

const emit = defineEmits<{
  (e: 'update:formData', formData: CreateNotificationRequest): void
  (e: 'submit'): void
}>()

const formDataValue = computed({
  get: (): CreateNotificationRequest =>
    props.formData as CreateNotificationRequest,
  set: (val: CreateNotificationRequest) => emit('update:formData', val)
})

const handleSubmit = () => {
  emit('submit')
}
</script>

<template>
  <div>
    <p class="text-h6">
      生産者登録
    </p>
    <organisms-notification-form :form-data="formDataValue" @submit="handleSubmit" />
  </div>
</template>
