<script lang="ts" setup>
import { storeToRefs } from 'pinia'

import { useAlert } from '~/lib/hooks'
import { useAdminStore, useAuthStore, useCommonStore, useNotificationStore, usePromotionStore } from '~/store'
import { NotificationType, type UpdateNotificationRequest } from '~/types/api'

const route = useRoute()
const router = useRouter()
const commonStore = useCommonStore()
const authStore = useAuthStore()
const adminStore = useAdminStore()
const notificationStore = useNotificationStore()
const promotionStore = usePromotionStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const notificationId = route.params.id as string

const { role } = storeToRefs(authStore)
const { admin } = storeToRefs(adminStore)
const { notification } = storeToRefs(notificationStore)
const { promotion } = storeToRefs(promotionStore)

const loading = ref<boolean>(false)
const formData = ref<UpdateNotificationRequest>({
  targets: [],
  title: '',
  body: '',
  note: '',
  publishedAt: 0,
})

const fetchState = useAsyncData(async (): Promise<void> => {
  try {
    await notificationStore.getNotification(notificationId)
    if (notification.value.type === NotificationType.PROMOTION) {
      await promotionStore.getPromotion(notification.value.promotionId)
    }
    formData.value = { ...notification.value }
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
})

const isLoading = (): boolean => {
  return fetchState?.pending?.value || loading.value
}

const handleSubmit = async (): Promise<void> => {
  try {
    loading.value = true
    await notificationStore.updateNotification(notificationId, formData.value)
    commonStore.addSnackbar({
      message: 'お知らせ情報の編集が完了しました',
      color: 'info',
    })
    router.push('/notifications')
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    loading.value = false
  }
}

try {
  await fetchState.execute()
}
catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-notification-edit
    v-model:form-data="formData"
    :loading="isLoading()"
    :role="role"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :notification="notification"
    :promotion="promotion"
    :admin="admin"
    @submit="handleSubmit"
  />
</template>
