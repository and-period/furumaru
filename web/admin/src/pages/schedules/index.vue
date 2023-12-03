<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useAlert, usePagination } from '~/lib/hooks'
import { useAuthStore, useCommonStore, useCoordinatorStore, useScheduleStore } from '~/store'
import type { Schedule } from '~/types/api'

const router = useRouter()
const commonStore = useCommonStore()
const authStore = useAuthStore()
const coordinatorStore = useCoordinatorStore()
const scheduleStore = useScheduleStore()
const pagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const { role } = storeToRefs(authStore)
const { coordinators } = storeToRefs(coordinatorStore)
const { schedules, total } = storeToRefs(scheduleStore)

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

const handleClickApproval = async (scheduleId: string): Promise<void> => {
  try {
    const schedule = schedules.value.find((schedule: Schedule): boolean => {
      return schedule.id === scheduleId
    })
    if (!schedule) {
      throw new Error(`failed to find schedule. scheduleId=${scheduleId}`)
    }
    await scheduleStore.approveSchedule(schedule)
    commonStore.addSnackbar({
      message: `${schedule.title}を更新しました。`,
      color: 'info'
    })
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
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
    :role="role"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :schedules="schedules"
    :coordinators="coordinators"
    :table-items-per-page="pagination.itemsPerPage.value"
    :table-items-total="total"
    @click:row="handleClickRow"
    @click:add="handleClickAdd"
    @click:delete="handleClickDelete"
    @click:approval="handleClickApproval"
    @click:update-page="handleUpdatePage"
    @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
  />
</template>
