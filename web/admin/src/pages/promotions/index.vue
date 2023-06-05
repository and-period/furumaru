<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { VDataTable } from 'vuetify/lib/labs/components'
import { useAlert, usePagination } from '~/lib/hooks'
import { usePromotionStore } from '~/store'

const router = useRouter()
const promotionStore = usePromotionStore()
const pagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const fetchState = useAsyncData(async () => {
  await fetchPromotions()
})

const { promotions, total } = storeToRefs(promotionStore)

const deleteDialog = ref<boolean>(false)
const sortBy = ref<VDataTable['sortBy']>([])

watch(pagination.itemsPerPage, () => {
  fetchPromotions()
})

const fetchPromotions = async () => {
  try {
    await promotionStore.fetchPromotions(pagination.itemsPerPage.value, pagination.offset.value)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleUpdatePage = async (page: number) => {
  pagination.updateCurrentPage(page)
  await fetchPromotions()
}

const handleClickAdd = () => {
  router.push('/promotions/new')
}

const handleClickRow = (promotionId: string) => {
  router.push(`/promotions/${promotionId}`)
}

const handleClickDelete = async (promotionId: string): Promise<void> => {
  try {
    await promotionStore.deletePromotion(promotionId)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  deleteDialog.value = false
}

const isLoading = (): boolean => {
  return fetchState?.pending?.value || false
}

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-promotion-list
    :loading="isLoading()"
    :is-alrt="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :promotions="promotions"
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
