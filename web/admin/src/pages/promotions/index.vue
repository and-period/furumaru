<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import type { VDataTable } from 'vuetify/lib/components/index.mjs'

import { useAlert, usePagination } from '~/lib/hooks'
import { useAuthStore, useCommonStore, usePromotionStore } from '~/store'

const router = useRouter()
const commonStore = useCommonStore()
const authStore = useAuthStore()
const promotionStore = usePromotionStore()
const pagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const { role } = storeToRefs(authStore)
const { promotions, total } = storeToRefs(promotionStore)

const loading = ref<boolean>(false)
const deleteDialog = ref<boolean>(false)
const sortBy = ref<VDataTable['sortBy']>([])

const fetchState = useAsyncData(async (): Promise<void> => {
  await fetchPromotions()
})

watch(pagination.itemsPerPage, (): void => {
  fetchPromotions()
})

const fetchPromotions = async (): Promise<void> => {
  try {
    const orders: string[] = sortBy.value.map((item: any) => {
      switch (item.order) {
        case 'asc':
          return item.key
        case 'desc':
          return `-${item.key}`
        default:
          return item.order ? item.key : `-${item.key}`
      }
    }) || []

    await promotionStore.fetchPromotions(pagination.itemsPerPage.value, pagination.offset.value, orders)
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
  await fetchPromotions()
}

const handleClickAdd = (): void => {
  router.push('/promotions/new')
}

const handleClickRow = (promotionId: string): void => {
  router.push(`/promotions/${promotionId}`)
}

const handleClickDelete = async (promotionId: string): Promise<void> => {
  try {
    loading.value = true
    await promotionStore.deletePromotion(promotionId)
    commonStore.addSnackbar({
      message: 'セール情報の削除が完了しました',
      color: 'info'
    })
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  } finally {
    deleteDialog.value = false
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
  <templates-promotion-list
    v-model:delete-dialog="deleteDialog"
    v-model:sort-by="sortBy"
    :loading="isLoading()"
    :role="role"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :promotions="promotions"
    :table-items-per-page="pagination.itemsPerPage.value"
    :table-items-total="total"
    @click:row="handleClickRow"
    @click:add="handleClickAdd"
    @click:delete="handleClickDelete"
    @click:update-page="handleUpdatePage"
    @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
    @update:sort-by="fetchState.refresh"
  />
</template>
