<script lang="ts" setup>
import dayjs from 'dayjs'
import { useAlert } from '~/lib/hooks'

import { useNotificationStore } from '~/store'
import { CreateNotificationRequest } from '~/types/api'
import { NotificationTime } from '~/types/props'

const router = useRouter()
const notificationStore = useNotificationStore()
const { alertType, isShow, alertText, show } = useAlert('error')

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
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}
</script>

<template>
  <templates-notification-new
    :form-data="formData"
    :time-data="timeData"
    :is-alrt="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    @submit="handleSubmit"
  />
</template>
