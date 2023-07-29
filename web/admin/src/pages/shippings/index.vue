<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useAlert, usePagination } from '~/lib/hooks'
import { useAuthStore, useCommonStore, useCoordinatorStore, useShippingStore } from '~/store'
import { Shipping } from '~/types/api'

const router = useRouter()
const authStore = useAuthStore()
const commonStore = useCommonStore()
const coordinatorStore = useCoordinatorStore()
const shippingStore = useShippingStore()
const pagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const { role } = storeToRefs(authStore)
const { coordinators } = storeToRefs(coordinatorStore)
const { shippings, totalItems } = storeToRefs(shippingStore)

const loading = ref<boolean>(false)
const deleteDialog = ref<boolean>(false)
const selectedItem = ref<Shipping>()

const fetchState = useAsyncData(async (): Promise<void> => {
  await fetchShippings()
})

watch(pagination.itemsPerPage, (): void => {
  fetchState.refresh()
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

const handleClickDelete = (shippingId: string): void => {
  const shipping = shippings.value.find((shipping: Shipping): boolean => {
    return shipping.id === shippingId
  })
  if (!shipping) {
    return
  }
  selectedItem.value = shipping
  deleteDialog.value = true
}

const handleSubmitDelete = async (): Promise<void> => {
  try {
    loading.value = true
    await shippingStore.deleteShipping(selectedItem.value?.id || '')
    commonStore.addSnackbar({
      color: 'info',
      message: '配送設定を削除しました。'
    })
    fetchState.refresh()
    deleteDialog.value = false
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  } finally {
    loading.value = false
  }
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
    :role="role"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :shippings="shippings"
    :coordinators="coordinators"
    :table-items-per-page="pagination.itemsPerPage.value"
    :table-items-total="totalItems"
    @click:row="handleClickRow"
    @click:add="handleClickAdd"
    @click:delete="handleClickDelete"
    @update:page="handleUpdatePage"
    @update:items-per-page="pagination.handleUpdateItemsPerPage"
    @submit:delete="handleSubmitDelete"
  />
</template>
