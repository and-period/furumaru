<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useAlert, usePagination } from '~/lib/hooks'
import { useAuthStore, useShippingStore } from '~/store'

const router = useRouter()
const authStore = useAuthStore()
const shippingStore = useShippingStore()
const pagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const { role } = storeToRefs(authStore)
const { shippings, totalItems } = storeToRefs(shippingStore)

const loading = ref<boolean>(false)

const fetchState = useAsyncData(async (): Promise<void> => {
  await fetchShippings()
})

watch(pagination.itemsPerPage, (): void => {
  fetchShippings()
})

const fetchShippings = async (): Promise<void> => {
  try {
    await shippingStore.fetchShippings(pagination.itemsPerPage.value, pagination.offset.value)
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
  await fetchShippings()
}

const handleClickAdd = (): void => {
  router.push('/shippings/new')
}

const handleClickRow = (shippingId: string): void => {
  router.push(`/shippings/${shippingId}`)
}

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-shipping-list
    :loading="isLoading()"
    :role="role"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :shippings="shippings"
    :table-items-per-page="pagination.itemsPerPage.value"
    :table-items-total="totalItems"
    @click:row="handleClickRow"
    @click:add="handleClickAdd"
    @click:update-page="handleUpdatePage"
    @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
  />
</template>
