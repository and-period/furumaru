<script lang="ts" setup>
import { useAlert, usePagination } from '~/lib/hooks'
import { useOrderStore } from '~/store'

const router = useRouter()
const orderStore = useOrderStore()
const pagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const dialog = reactive({
  import: false,
  export: false
})
const importFormData = reactive({ // TODO: API設計が決まり次第型定義の厳格化
  company: ''
})
const exportFormData = reactive({ // TODO: API設計が決まり次第型定義の厳格化
  company: ''
})

const orders = computed(() => {
  return orderStore.orders
})
const ordersTotal = computed(() => {
  return orderStore.totalItems
})

watch(pagination.itemsPerPage, () => {
  fetchState.refresh()
})

const fetchState = useAsyncData(async () => {
  await fetchOrders()
})

const isLoading = (): boolean => {
  return fetchState?.pending?.value || false
}

const fetchOrders = async (): Promise<void> => {
  try {
    await orderStore.fetchOrders(pagination.itemsPerPage.value, pagination.offset.value)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log('failed to fetch orders', err)
  }
}

const handleUpdateTablePage = async (page: number): Promise<void> => {
  pagination.updateCurrentPage(page)
  await fetchState.refresh()
}

const handleEdit = (orderId: string): void => {
  router.push(`/orders/${orderId}`)
}

const handleImport = () => {
  // TODO: APIの実装が完了後に対応
  console.log('debug', 'submit:import')
  dialog.import = false
}

const handleExport = () => {
  // TODO: APIの実装が完了後に対応
  console.log('debug', 'submit:export')
  dialog.export = false
}

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-order-list
    v-model:dialog="dialog"
    v-model:import-form-data="importFormData"
    v-model:export-form-data="exportFormData"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :orders="orders"
    :table-items-length="ordersTotal"
    :table-items-per-page="pagination.itemsPerPage"
    @click:edit="handleEdit"
    @click:update-page="handleUpdateTablePage"
    @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
    @submit:import="handleImport"
    @submit:export="handleExport"
  />
</template>
