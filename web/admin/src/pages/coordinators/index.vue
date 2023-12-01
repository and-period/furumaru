<script lang="ts" setup>
import { storeToRefs } from 'pinia'

import { useAlert, usePagination } from '~/lib/hooks'
import { useCommonStore, useCoordinatorStore, useProductTypeStore } from '~/store'

const router = useRouter()
const commonStore = useCommonStore()
const coordinatorStore = useCoordinatorStore()
const productTypeStore = useProductTypeStore()
const pagination = usePagination()
const { isShow, alertText, alertType, show } = useAlert('error')

const { coordinators, totalItems } = storeToRefs(coordinatorStore)
const { productTypes } = storeToRefs(productTypeStore)

const loading = ref<boolean>(false)
const deleteDialog = ref<boolean>(false)

const fetchState = useAsyncData(async (): Promise<void> => {
  await fetchCoordinators()
})

watch(pagination.itemsPerPage, (): void => {
  fetchCoordinators()
})

const fetchCoordinators = async (): Promise<void> => {
  try {
    await coordinatorStore.fetchCoordinators(pagination.itemsPerPage.value, pagination.offset.value)
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
  await fetchState.refresh()
}

const handleClickAdd = () => {
  router.push('/coordinators/new')
}

const handleClickRow = (coordinatorId: string) => {
  router.push(`/coordinators/${coordinatorId}`)
}

const handleClickDelete = async (coordinatorId: string): Promise<void> => {
  try {
    loading.value = true
    await coordinatorStore.deleteCoordinator(coordinatorId)
    commonStore.addSnackbar({
      color: 'info',
      message: 'コーディネーターの削除が完了しました。'
    })
    fetchState.refresh()
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    window.scrollTo({
      top: 0,
      behavior: 'smooth'
    })
    console.log(err)
  } finally {
    deleteDialog.value = false
    loading.value = true
  }
}

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-coordinator-list
    v-model:delete-dialog="deleteDialog"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :coordinators="coordinators"
    :table-items-per-page="pagination.itemsPerPage.value"
    :table-items-total="totalItems"
    :product-types="productTypes"
    @click:row="handleClickRow"
    @click:add="handleClickAdd"
    @click:delete="handleClickDelete"
    @click:update-page="handleUpdatePage"
    @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
  />
</template>
