<script lang="ts" setup>
import dayjs from 'dayjs'

import { useNotificationStore } from '~/store'
import { NotificationResponse } from '~/types/api'
import { NotificationTime } from '~/types/props'

const router = useRouter()
const route = useRoute()
const id = route.params.id

const { getNotification, editNotification } = useNotificationStore()

const formData = reactive<NotificationResponse>({
  id,
  createdBy: '',
  creatorName: '',
  updatedBy: '',
  title: '',
  body: '',
  targets: [],
  public: false,
  publishedAt: -1,
  createdAt: -1,
  updatedAt: -1
})

const timeData = reactive<NotificationTime>({
  publishedDate: '',
  publishedTime: ''
})

const fetchState = useAsyncData(async () => {
  const notification = await getNotification(id)
  formData.title = notification.title
  formData.body = notification.body
  formData.targets = notification.targets
  formData.public = notification.public
  formData.publishedAt = notification.publishedAt
  timeData.publishedDate = dayjs
    .unix(notification.publishedAt)
    .format('YYYY-MM-DD')
  timeData.publishedTime = dayjs.unix(notification.publishedAt).format('HH:mm')
})

const handleSubmit = async () => {
  try {
    await editNotification(id, formData)
    router.push('/notifications')
  } catch (error) {
    console.log(error)
  }
}
</script>

<template>
  <the-notification-edit-form-page
    :form-data="formData"
    :time-data="timeData"
    :form-data-loading="fetchState.pending"
    @submit="handleSubmit"
  />
</template>
