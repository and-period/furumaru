<script lang="ts" setup>
import dayjs from 'dayjs'
import { storeToRefs } from 'pinia'
import type { VDataTable } from 'vuetify/components'

import { useAlert, usePagination } from '~/lib/hooks'
import { useCommonStore, useCoordinatorStore, useCustomerStore, useOrderStore, usePromotionStore } from '~/store'
import { CharacterEncodingType } from '~/types'
import { ShippingCarrier } from '~/types/api/v1'
import type { ExportOrdersRequest } from '~/types/api/v1'

const router = useRouter()
const commonStore = useCommonStore()
const orderStore = useOrderStore()
const coordinatorStore = useCoordinatorStore()
const customerStore = useCustomerStore()
const promotionStore = usePromotionStore()
const pagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const { orders, totalItems } = storeToRefs(orderStore)
const { coordinators } = storeToRefs(coordinatorStore)
const { customers } = storeToRefs(customerStore)
const { promotions } = storeToRefs(promotionStore)

const loading = ref<boolean>(false)
const importDialog = ref<boolean>(false)
const exportDialog = ref<boolean>(false)
const selectedItemId = ref<string>('')
const sortBy = ref<VDataTable['sortBy']>([])

// TODO: API設計が決まり次第型定義の厳格化
const importFormData = ref({
  company: '',
})
const exportFormData = ref<ExportOrdersRequest>({
  shippingCarrier: ShippingCarrier.ShippingCarrierUnknown,
  characterEncodingType: CharacterEncodingType.UTF8,
})

const fetchState = useAsyncData('orders', async (): Promise<void> => {
  await fetchOrders()
})

watch(pagination.itemsPerPage, (): void => {
  fetchState.refresh()
})

watch(sortBy, (): void => {
  fetchState.refresh()
})

const fetchOrders = async (): Promise<void> => {
  try {
    await orderStore.fetchOrders(pagination.itemsPerPage.value, pagination.offset.value)
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

const handleUpdateTablePage = async (page: number): Promise<void> => {
  pagination.updateCurrentPage(page)
  await fetchState.refresh()
}

const handleClickRow = (orderId: string): void => {
  router.push(`/orders/${orderId}`)
}

const handleImport = (): void => {
  // TODO: APIの実装が完了後に対応
  console.log('debug', 'submit:import')
  importDialog.value = false
}

const handleExport = async (): Promise<void> => {
  try {
    const data = await orderStore.exportOrders(exportFormData.value)
    const a = document.createElement('a')
    a.href = window.URL.createObjectURL(new Blob([data]))
    a.setAttribute('download', `orders_${dayjs().format('YYYYMMDDhhmmss')}.csv`)
    a.click()
    a.remove()
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  exportDialog.value = false
}

try {
  await fetchState.execute()
}
catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-order-list
    v-model:import-dialog="importDialog"
    v-model:export-dialog="exportDialog"
    v-model:import-form-data="importFormData"
    v-model:export-form-data="exportFormData"
    v-model:selected-item-id="selectedItemId"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :orders="orders"
    :coordinators="coordinators"
    :customers="customers"
    :promotions="promotions"
    :table-items-per-page="pagination.itemsPerPage.value"
    :table-items-total="totalItems"
    :table-sort-by="sortBy"
    @click:row="handleClickRow"
    @click:update-page="handleUpdateTablePage"
    @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
    @submit:import="handleImport"
    @submit:export="handleExport"
    @update:sort-by="fetchState.refresh"
  />
</template>
