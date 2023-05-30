<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { VDataTable } from 'vuetify/labs/components'

import { useAlert, usePagination } from '~/lib/hooks'
import { useNotificationStore } from '~/store'

const router = useRouter()
const notificationStore = useNotificationStore()
const pagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const fetchState = useAsyncData(async () => {
  await fetchNotifications()
})

const { notifications, totalItems } = storeToRefs(notificationStore)

const deleteDialog = ref<boolean>(false)
const sortBy = ref<VDataTable['sortBy']>([])

watch(pagination.itemsPerPage, () => {
  fetchNotifications()
})

const handleUpdatePage = async (page: number) => {
  pagination.updateCurrentPage(page)
  await fetchNotifications()
}

const fetchNotifications = async () => {
  try {
    await notificationStore.fetchNotifications(pagination.itemsPerPage.value, pagination.offset.value)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleClickAdd = () => {
  router.push('/notifications/add')
}

const handleClickRow = (notificationId: string) => {
  router.push(`/notifications/edit/${notificationId}`)
}

const handleClickDelete = async (notificationId: string): Promise<void> => {
  try {
    await notificationStore.deleteNotification(notificationId)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  deleteDialog.value = false
}

const isLoading = (): boolean => {
  return fetchState?.pending?.value || false
}

try {
  fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-notification-list
    :loading="isLoading()"
    :is-alrt="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :notifications="notifications"
    :table-items-per-page="pagination.itemsPerPage.value"
    :table-items-total="totalItems"
    :table-sort-by="sortBy"
    @click:row="handleClickRow"
    @click:add="handleClickAdd"
    @click:delete="handleClickDelete"
    @click:update-page="handleUpdatePage"
    @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
    @update:sort-by="fetchState.refresh"
  />
</template>
