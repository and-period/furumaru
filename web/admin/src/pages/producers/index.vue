<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { VDataTable } from 'vuetify/lib/labs/components'

import { useAlert, usePagination } from '~/lib/hooks'
import { useCommonStore, useProducerStore } from '~/store'

const router = useRouter()

const { addSnackbar } = useCommonStore()
const producerStore = useProducerStore()
const pagination = usePagination()
const { isShow, alertText, alertType, show } = useAlert('error')

const sortBy = ref<VDataTable['sortBy']>([])

const { producers, totalItems } = storeToRefs(producerStore)

watch(pagination.itemsPerPage, () => {
  producerStore.fetchProducers(pagination.itemsPerPage.value, 0, '')
})

const handleUpdatePage = async (page: number) => {
  pagination.updateCurrentPage(page)
  await producerStore.fetchProducers(pagination.itemsPerPage.value, pagination.offset.value, '')
}

const fetchState = useAsyncData(async () => {
  try {
    await producerStore.fetchProducers(pagination.itemsPerPage.value, pagination.offset.value)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
})

const handleClickAdd = () => {
  router.push('/producers/add')
}

const handleClickAddVideo = (producerId: string) => {
  console.log(producerId)
}

const handleClickRow = (producerId: string) => {
  router.push(`/producers/edit/${producerId}`)
}

const handleClickDelete = async (producerId: string) => {
  try {
    await producerStore.deleteProducer(producerId)
    addSnackbar({ color: 'info', message: '生産者を削除しました。' })
    fetchState.refresh()
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
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
  <templates-producer-list
    :loading="isLoading()"
    :is-alrt="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :producers="producers"
    :table-items-per-page="pagination.itemsPerPage.value"
    :table-items-total="totalItems"
    :table-sort-by="sortBy"
    @click:row="handleClickRow"
    @click:add="handleClickAdd"
    @click:add-video="handleClickAddVideo"
    @click:delete="handleClickDelete"
    @click:update-page="handleUpdatePage"
    @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
    @update:sort-by="fetchState.refresh"
  />
</template>
