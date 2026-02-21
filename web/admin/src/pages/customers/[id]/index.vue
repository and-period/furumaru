<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useAlert, usePagination } from '~/lib/hooks'
import { useAddressStore, useAuthStore, useCommonStore, useCustomerStore } from '~/store'

const route = useRoute()
const router = useRouter()
const commonStore = useCommonStore()
const authStore = useAuthStore()
const addressStore = useAddressStore()
const customerStore = useCustomerStore()
const pagination = usePagination()
const { isShow, alertText, alertType, show } = useAlert('error')

const customerId = route.params.id as string

const { adminType } = storeToRefs(authStore)
const { address } = storeToRefs(addressStore)
const {
  customer,
  orders,
  totalOrderCount,
  totalPaymentCount,
  totalProductAmount,
  totalPaymentAmount,
} = storeToRefs(customerStore)

const loading = ref<boolean>(false)
const deleteDialog = ref<boolean>(false)

const fetchState = useAsyncData('customer-detail', async () => {
  await fetchCustomerOrders()
})

watch(pagination.itemsPerPage, (): void => {
  fetchCustomerOrders()
})

const fetchCustomerOrders = async (): Promise<void> => {
  try {
    await customerStore.fetchCustomerOrders(customerId, pagination.itemsPerPage.value, pagination.offset.value)
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
  await fetchState.refresh()
}

const handleDelete = async (): Promise<void> => {
  try {
    loading.value = true
    await customerStore.deleteCustomer(customerId)
    commonStore.addSnackbar({
      color: 'info',
      message: '顧客の削除が完了しました。',
    })
    router.push('/customers')
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    loading.value = false
  }
}

const handleClickRow = (orderId: string) => {
  router.push(`/orders/${orderId}`)
}

try {
  await customerStore.getCustomer(customerId)
  await fetchState.execute()
}
catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-customer-show
    v-model:delete-dialog="deleteDialog"
    :loading="isLoading()"
    :admin-type="adminType"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :customer="customer"
    :address="address"
    :orders="orders"
    :total-order-count="totalOrderCount"
    :total-payment-count="totalPaymentCount"
    :total-product-amount="totalProductAmount"
    :total-payment-amount="totalPaymentAmount"
    @click:row="handleClickRow"
    @click:update-page="handleUpdatePage"
    @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
    @submit:delete="handleDelete"
  />
</template>
