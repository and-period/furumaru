<script lang="ts" setup>
import { unix } from 'dayjs'
import { useAlert } from '~/lib/hooks'

import { useNotificationStore } from '~/store'
import { NotificationResponse } from '~/types/api'
import { NotificationTime } from '~/types/props'

const router = useRouter()
const route = useRoute()
const id = route.params.id as string

const { getNotification, updateNotification } = useNotificationStore()
const { alertType, isShow, alertText, show } = useAlert('error')

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
  timeData.publishedDate = unix(notification.publishedAt).format('YYYY-MM-DD')
  timeData.publishedTime = unix(notification.publishedAt).format('HH:mm')
})

const isLoading = (): boolean => {
  return fetchState?.pending?.value || false
}

const handleSubmit = async () => {
  try {
    await updateNotification(id, formData)
    router.push('/notifications')
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-notification-edit
    :form-data="formData"
    :time-data="timeData"
    :form-data-loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    @submit="handleSubmit"
  />
</template>
