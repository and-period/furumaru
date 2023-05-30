<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { VDataTable } from 'vuetify/labs/components'

import { useAlert, usePagination } from '~/lib/hooks'
import { useCommonStore, useCoordinatorStore } from '~/store'

const router = useRouter()
const commonStore = useCommonStore()
const coordinatorStore = useCoordinatorStore()
const pagination = usePagination()
const { isShow, alertText, alertType, show } = useAlert('error')

const fetchState = useAsyncData(async () => {
  await fetchCoordinators()
})

const { coordinators, totalItems } = storeToRefs(coordinatorStore)

const deleteDialog = ref<boolean>(false)
const sortBy = ref<VDataTable['sortBy']>([])

watch(pagination.itemsPerPage, () => {
  fetchCoordinators()
})

const fetchCoordinators = async () => {
  try {
    await coordinatorStore.fetchCoordinators(pagination.itemsPerPage.value, pagination.offset.value)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleUpdatePage = async (page: number) => {
  pagination.updateCurrentPage(page)
  await fetchCoordinators()
}

const handleClickAdd = () => {
  router.push('/coordinators/add')
}

const handleClickRow = (coordinatorId: string) => {
  router.push(`/coordinators/edit/${coordinatorId}`)
}

const handleClickDelete = async (coordinatorId: string): Promise<void> => {
  try {
    await coordinatorStore.deleteCoordinator(coordinatorId)
    commonStore.addSnackbar({
      color: 'info',
      message: 'コーディネータを削除しました。'
    })
    fetchState.refresh()
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
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-coordinator-list
    :loading="isLoading()"
    :is-alrt="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :coordinators="coordinators"
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
