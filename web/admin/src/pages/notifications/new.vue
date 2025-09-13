<script lang="ts" setup>
import dayjs from 'dayjs'
import { storeToRefs } from 'pinia'

import { useAlert } from '~/lib/hooks'
import { useAdminStore, useCommonStore, useNotificationStore, usePromotionStore } from '~/store'
import { NotificationType } from '~/types/api/v1'
import type { CreateNotificationRequest } from '~/types/api/v1'

const router = useRouter()
const commonStore = useCommonStore()
const adminStore = useAdminStore()
const notificationStore = useNotificationStore()
const promotionStore = usePromotionStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const { admins } = storeToRefs(adminStore)
const { promotions } = storeToRefs(promotionStore)

const loading = ref<boolean>(false)
const formData = ref<CreateNotificationRequest>({
  type: NotificationType.NotificationTypeOther,
  targets: [],
  title: '',
  body: '',
  note: '',
  publishedAt: dayjs().unix(),
  promotionId: '',
})

const fetchPromotions = async (): Promise<void> => {
  try {
    await promotionStore.fetchPromotions()
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const updateNotificationType = async (type: NotificationType): Promise<void> => {
  loading.value = true
  switch (type) {
    case NotificationType.NotificationTypePromotion:
      await fetchPromotions()
      break
    default:
    // 何もしない
  }
  loading.value = false
}

const handleSubmit = async () => {
  try {
    loading.value = true
    await notificationStore.createNotification(formData.value)
    commonStore.addSnackbar({
      message: `${formData.value.title}を作成しました。`,
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
</script>

<template>
  <templates-notification-new
    v-model:form-data="formData"
    :loading="loading"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :admins="admins"
    :promotions="promotions"
    @update:notification-type="updateNotificationType"
    @submit="handleSubmit"
  />
</template>
