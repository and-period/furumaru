<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { VDataTable } from 'vuetify/lib/labs/components'
import { useAlert, usePagination } from '~/lib/hooks'
import { useAdministratorStore } from '~/store'

const administratorStore = useAdministratorStore()
const pagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const fetchState = useAsyncData(async () => {
  await fetchAdministrators()
})

const { administrators, total } = storeToRefs(administratorStore)

const deleteDialog = ref<boolean>(false)

const sortBy = reactive<VDataTable['sortBy']>([])

watch(pagination.itemsPerPage, () => {
  fetchState.refresh()
})
watch(sortBy, () => {
  fetchState.refresh()
})

const fetchAdministrators = async () => {
  try {
    const orders: string[] = sortBy?.map((item) => {
      switch (item.order) {
        case 'asc':
          return item.key
        case 'desc':
          return `-${item.key}`
        default:
          return item.order ? item.key : `-${item.key}`
      }
    }) || []

    await administratorStore.fetchAdministrators(pagination.itemsPerPage.value, pagination.offset.value, orders)
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

const handleUpdatePage = async (page: number) => {
  pagination.updateCurrentPage(page)
  await fetchState.refresh()
}

const handleClickAdd = () => {
  console.log('click:add')
}

const handleClickDelete = (administratorId: string) => {
  console.log('click:delete', administratorId)
}

const handleClickRow = (administratorId: string) => {
  console.log('click:row', administratorId)
}

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-administrator-list
    v-model:delete-dialog="deleteDialog"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :administrators="administrators"
    :table-items-per-page="pagination.itemsPerPage.value"
    :table-items-total="total"
    :table-sort-by="sortBy"
    @click:row="handleClickRow"
    @click:add="handleClickAdd"
    @click:delete="handleClickDelete"
    @click:update-page="handleUpdatePage"
    @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
    @update:sort-by="fetchState.refresh"
  />
</template>
