<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useAlert, usePagination } from '~/lib/hooks'
import { useCoordinatorStore, useScheduleStore, useShippingStore } from '~/store'

const router = useRouter()
const coordinatorStore = useCoordinatorStore()
const scheduleStore = useScheduleStore()
const shippingStore = useShippingStore()
const pagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const { coordinators } = storeToRefs(coordinatorStore)
const { schedules, total } = storeToRefs(scheduleStore)
const { shippings } = storeToRefs(shippingStore)

const loading = ref<boolean>(false)
const deleteDialog = ref<boolean>(false)

const fetchState = useAsyncData(async (): Promise<void> => {
  await fetchSchedules()
})

watch(pagination.itemsPerPage, (): void => {
  fetchSchedules()
})

const fetchSchedules = async (): Promise<void> => {
  try {
    await scheduleStore.fetchSchedules(pagination.itemsPerPage.value, pagination.offset.value)
  } catch (err) {
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
  await fetchSchedules()
}

const handleClickAdd = (): void => {
  router.push('/schedules/new')
}

const handleClickRow = (scheduleId: string): void => {
  router.push(`/schedules/${scheduleId}`)
}

const handleClickDelete = (): void => {
  console.log('debug', 'click:delete-schedule')
}

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-schedule-list
    v-model:delete-dialog="deleteDialog"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :schedules="schedules"
    :coordinators="coordinators"
    :shippings="shippings"
    :table-items-per-page="pagination.itemsPerPage.value"
    :table-items-total="total"
    @click:row="handleClickRow"
    @click:add="handleClickAdd"
    @click:delete="handleClickDelete"
    @click:update-page="handleUpdatePage"
    @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
  />
</template>
