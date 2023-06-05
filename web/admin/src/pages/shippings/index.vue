<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { VDataTable } from 'vuetify/lib/labs/components'
import { useAlert, usePagination } from '~/lib/hooks'
import { useShippingStore } from '~/store'

const router = useRouter()
const shippingStore = useShippingStore()
const pagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const fetchState = useAsyncData(async () => {
  await fetchShippings()
})

const { shippings, totalItems } = storeToRefs(shippingStore)

const deleteDialog = ref<boolean>(false)
const sortBy = ref<VDataTable['sortBy']>([])

watch(pagination.itemsPerPage, () => {
  fetchShippings()
})

const fetchShippings = async () => {
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
  return fetchState?.pending?.value || false
}

const handleUpdatePage = async (page: number) => {
  pagination.updateCurrentPage(page)
  await fetchShippings()
}

const handleClickAdd = () => {
  router.push('/shippings/new')
}

const handleClickRow = (shippingId: string) => {
  router.push(`/shippings/${shippingId}`)
}

const handleClickDelete = (shippingId: string) => {
  console.log('delete', shippingId)
}

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-shipping-list
    v-model:delete-dialog="deleteDialog"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :shippings="shippings"
    :table-items-per-page="pagination.itemsPerPage.value"
    :table-items-total="totalItems"
    :table-sort-by="sortBy"
    @click:row="handleClickRow"
    @click:add="handleClickAdd"
    @click:delete="handleClickDelete"
    @click:update-page="handleUpdatePage"
    @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
    @update:sort-by="fetchState.refresh"
  />
</template>
