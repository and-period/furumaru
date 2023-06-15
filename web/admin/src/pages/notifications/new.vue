<script lang="ts" setup>
import dayjs from 'dayjs'
import { useAlert } from '~/lib/hooks'

import { useNotificationStore } from '~/store'
import { CreateNotificationRequest } from '~/types/api'

const router = useRouter()
const notificationStore = useNotificationStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const loading = ref<boolean>(false)
const formData = ref<CreateNotificationRequest>({
  title: '',
  body: '',
  targets: [],
  public: false,
  publishedAt: dayjs().unix()
})

const handleSubmit = async () => {
  try {
    loading.value = true
    await notificationStore.createNotification(formData.value)
    router.push('/notifications')
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <templates-notification-new
    v-model:form-data="formData"
    :loading="loading"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    @submit="handleSubmit"
  />
</template>
