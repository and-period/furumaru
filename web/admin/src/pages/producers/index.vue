<script lang="ts" setup>
import { storeToRefs } from 'pinia'

import { useAlert, usePagination } from '~/lib/hooks'
import { useAuthStore, useCommonStore, useCoordinatorStore, useProducerStore, useShopStore } from '~/store'

const router = useRouter()
const authStore = useAuthStore()
const commonStore = useCommonStore()
const coordinatorStore = useCoordinatorStore()
const producerStore = useProducerStore()
const shopStore = useShopStore()
const pagination = usePagination()
const { isShow, alertText, alertType, show } = useAlert('error')

const { adminType } = storeToRefs(authStore)
const { coordinators } = storeToRefs(coordinatorStore)
const { producers, totalItems } = storeToRefs(producerStore)
const { shops } = storeToRefs(shopStore)

const loading = ref<boolean>(false)
const deleteDialog = ref<boolean>(false)

const fetchState = useAsyncData('producers', async (): Promise<void> => {
  await fetchProducers()
})

watch(pagination.itemsPerPage, (): void => {
  fetchState.refresh()
})

const fetchProducers = async (): Promise<void> => {
  try {
    await producerStore.fetchProducers(pagination.itemsPerPage.value, pagination.offset.value)
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
  await producerStore.fetchProducers(pagination.itemsPerPage.value, pagination.offset.value)
}

const handleClickAdd = (): void => {
  router.push('/producers/new')
}

const handleClickRow = (producerId: string): void => {
  router.push(`/producers/${producerId}`)
}

const handleClickDelete = async (producerId: string): Promise<void> => {
  try {
    loading.value = true
    await producerStore.deleteProducer(producerId)
    commonStore.addSnackbar({
      color: 'info',
      message: '生産者を削除しました。',
    })
    fetchState.refresh()
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
  await fetchState.execute()
}
catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-producer-list
    v-model:delete-dialog="deleteDialog"
    :loading="isLoading()"
    :admin-type="adminType"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :producers="producers"
    :shops="shops"
    :coordinators="coordinators"
    :table-items-per-page="pagination.itemsPerPage.value"
    :table-items-total="totalItems"
    @click:row="handleClickRow"
    @click:add="handleClickAdd"
    @click:delete="handleClickDelete"
    @click:update-page="handleUpdatePage"
    @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
  />
</template>
