<script lang="ts" setup>
import dayjs from 'dayjs'

import { useNotificationStore } from '~/store'
import { CreateNotificationRequest } from '~/types/api'

const router = useRouter()
const notificationStore = useNotificationStore()

const formData = reactive<CreateNotificationRequest>({
  title: '',
  body: '',
  targets: [],
  public: false,
  publishedAt: dayjs().unix()
})

const handleSubmit = async () => {
  try {
    await notificationStore.createNotification(formData)
    router.push('/notifications')
  } catch (error) {
    console.log(error)
  }
}
</script>

<template>
  <templates-notification-create-form-page
    :form-data="formData"
    @submit="handleSubmit"
  />
</template>
