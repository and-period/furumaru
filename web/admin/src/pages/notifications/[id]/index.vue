<script lang="ts" setup>
import { storeToRefs } from 'pinia'

import { useAlert } from '~/lib/hooks'
import { useAdminStore, useAuthStore, useCommonStore, useNotificationStore, usePromotionStore } from '~/store'
import { AdminType, NotificationType } from '~/types/api'
import type { UpdateNotificationRequest } from '~/types/api'

const route = useRoute()
const router = useRouter()
const commonStore = useCommonStore()
const authStore = useAuthStore()
const adminStore = useAdminStore()
const notificationStore = useNotificationStore()
const promotionStore = usePromotionStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const notificationId = route.params.id as string

const { adminType } = storeToRefs(authStore)
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

const isEditable = (): boolean => {
  return adminType.value === AdminType.ADMINISTRATOR
}

try {
  await fetchState.execute()
}
catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <div>
    <templates-notification-edit
      v-model:form-data="formData"
      :loading="isLoading()"
      :admin-type="adminType"
      :is-alert="isShow"
      :alert-type="alertType"
      :alert-text="alertText"
      :notification="notification"
      :promotion="promotion"
      :admin="admin"
      @submit="handleSubmit"
    />
    <div
      class="position-fixed bottom-0 left-0 w-100 bg-white pa-4 text-right elevation-3"
    >
      <div class="d-inline-flex ga-4">
        <v-btn
          color="secondary"
          variant="outlined"
          @click="$router.back()"
        >
          戻る
        </v-btn>
        <v-btn
          v-show="isEditable()"
          color="primary"
          variant="outlined"
          type="submit"
          form="update-notification-form"
        >
          更新
        </v-btn>
      </div>
    </div>
  </div>
</template>
