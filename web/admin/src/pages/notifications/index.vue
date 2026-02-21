<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import type { VDataTable } from 'vuetify/components'
import { useAlert, usePagination } from '~/lib/hooks'
import { useAdminStore, useAuthStore, useCommonStore, useNotificationStore } from '~/store'

const router = useRouter()
const commonStore = useCommonStore()
const authStore = useAuthStore()
const adminStore = useAdminStore()
const notificationStore = useNotificationStore()
const pagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const { adminType } = storeToRefs(authStore)
const { admins } = storeToRefs(adminStore)
const { notifications, totalItems } = storeToRefs(notificationStore)

const loading = ref<boolean>(false)
const deleteDialog = ref<boolean>(false)
const sortBy = ref<VDataTable['sortBy']>([])

const fetchState = useAsyncData('notifications', async (): Promise<void> => {
  await fetchNotifications()
})

watch(pagination.itemsPerPage, (): void => {
  fetchNotifications()
})
watch(sortBy, (): void => {
  fetchNotifications()
})

const fetchNotifications = async (): Promise<void> => {
  try {
    await notificationStore.fetchNotifications(pagination.itemsPerPage.value, pagination.offset.value)
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const isLoading = (): boolean => {
  return fetchState?.pending?.value || loading.value
}

const handleUpdatePage = async (page: number): Promise<void> => {
  pagination.updateCurrentPage(page)
  await fetchNotifications()
}

const handleClickAdd = (): void => {
  router.push('/notifications/new')
}

const handleClickRow = (notificationId: string) => {
  router.push(`/notifications/${notificationId}`)
}

const handleClickDelete = async (notificationId: string): Promise<void> => {
  try {
    loading.value = true
    await notificationStore.deleteNotification(notificationId)
    commonStore.addSnackbar({
      message: 'お知らせの削除が完了しました',
      color: 'info',
    })
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    deleteDialog.value = false
    loading.value = false
  }
}

try {
  fetchState.execute()
}
catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-notification-list
    v-model:delete-dialog="deleteDialog"
    :loading="isLoading()"
    :admin-type="adminType"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :notifications="notifications"
    :admins="admins"
    :table-items-per-page="pagination.itemsPerPage.value"
    :table-items-total="totalItems"
    :table-sort-by="sortBy"
    @click:row="handleClickRow"
    @click:add="handleClickAdd"
    @click:delete="handleClickDelete"
    @click:update-page="handleUpdatePage"
    @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
    @update:sort-by="fetchNotifications"
  />
</template>
