<script lang="ts" setup>
import { VDataTable } from 'vuetify/labs/components'

import { storeToRefs } from 'pinia'
import { useAlert, usePagination } from '~/lib/hooks'
import { useCustomerStore } from '~/store/customer'

const router = useRouter()
const customerStore = useCustomerStore()
const pagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const fetchState = useAsyncData(async () => {
  await fetchUsers()
})

const { customers, totalItems } = storeToRefs(customerStore)

const sortBy = ref<VDataTable['sortBy']>([])

watch(pagination.itemsPerPage, () => {
  fetchUsers()
})

const handleUpdatePage = async (page: number) => {
  pagination.updateCurrentPage(page)
  await fetchUsers()
}

const isLoading = (): boolean => {
  return fetchState?.pending?.value || false
}

const fetchUsers = async () => {
  try {
    await customerStore.fetchCustomers(pagination.itemsPerPage.value, pagination.offset.value)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleClickEdit = (customerId: string): void => {
  router.push(`/customers/${customerId}`)
}

const handleClickDelete = (customerId: string) => {
  console.log('削除ボタンクリック', customerId)
}

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-customer-list
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :customers="customers"
    :table-items-total="totalItems"
    :table-items-per-page="pagination.itemsPerPage.value"
    :table-sort-by="sortBy"
    @click:row="handleClickEdit"
    @click:delete="handleClickDelete"
    @click:update-page="handleUpdatePage"
    @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
    @update:sort-by="fetchState.refresh"
  />
</template>
