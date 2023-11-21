<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useAlert, usePagination } from '~/lib/hooks'
import { useCustomerStore } from '~/store'

const route = useRoute()
const router = useRouter()
const customerStore = useCustomerStore()
const pagination = usePagination()
const { isShow, alertText, alertType, show } = useAlert('error')

const customerId = route.params.id as string

const { customer, orders, totalOrders, totalAmount } = storeToRefs(customerStore)

const loading = ref<boolean>(false)

const fetchState = useAsyncData(async () => {
  await fetchCustomerOrders()
})

watch(pagination.itemsPerPage, (): void => {
  fetchCustomerOrders()
})

const fetchCustomerOrders = async (): Promise<void> => {
  try {
    await customerStore.fetchCustomerOrders(customerId, pagination.itemsPerPage.value, pagination.offset.value)
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

const handleClickRow = (orderId: string) => {
  router.push(`/orders/${orderId}`)
}

try {
  await customerStore.getCustomer(customerId)
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-customer-edit
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :customer="customer"
    :orders="orders"
    :order-total="totalOrders"
    :order-amount="totalAmount"
    @click:row="handleClickRow"
    @click:update-page="handleUpdatePage"
    @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
  />
</template>
