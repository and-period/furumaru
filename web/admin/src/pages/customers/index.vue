<script lang="ts" setup>
import type { VDataTable } from 'vuetify/components'
import { storeToRefs } from 'pinia'
import { useAlert, usePagination } from '~/lib/hooks'
import { useCommonStore, useCustomerStore } from '~/store'

const router = useRouter()
const commonStore = useCommonStore()
const customerStore = useCustomerStore()
const pagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const { customersToList, totalItems } = storeToRefs(customerStore)

const loading = ref<boolean>(false)
const sortBy = ref<VDataTable['sortBy']>([])

const fetchState = useAsyncData('customers', async () => {
  await fetchUsers()
})

watch(pagination.itemsPerPage, () => {
  fetchUsers()
})

watch(sortBy, (): void => {
  fetchState.refresh()
})

const fetchUsers = async (): Promise<void> => {
  try {
    await customerStore.fetchCustomers(pagination.itemsPerPage.value, pagination.offset.value)
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

const handleUpdatePage = async (page: number) => {
  pagination.updateCurrentPage(page)
  await fetchUsers()
}

const handleClickEdit = (customerId: string): void => {
  router.push(`/customers/${customerId}`)
}

try {
  await fetchState.execute()
}
catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-customer-list
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :customers="customersToList"
    :table-items-total="totalItems"
    :table-items-per-page="pagination.itemsPerPage.value"
    :table-sort-by="sortBy"
    @click:row="handleClickEdit"
    @click:update-page="handleUpdatePage"
    @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
    @update:sort-by="fetchState.refresh"
  />
</template>
