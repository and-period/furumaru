<script lang="ts" setup>
import dayjs from 'dayjs'

import { useNotificationStore } from '~/store'
import { CreateNotificationRequest } from '~/types/api'
import { NotificationTime } from '~/types/props'

const router = useRouter()
const notificationStore = useNotificationStore()

const timeData = reactive<NotificationTime>({
  publishedDate: '',
  publishedTime: ''
})

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
  <templates-notification-new
    :form-data="formData"
    :time-data="timeData"
    @submit="handleSubmit"
  />
</template>
