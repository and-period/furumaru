<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useAuthStore, useFeatureRequestStore } from '~/store'
import { AdminType } from '~/types/api/v1'

const router = useRouter()
const authStore = useAuthStore()
const featureRequestStore = useFeatureRequestStore()

const { featureRequests, total } = storeToRefs(featureRequestStore)
const { adminType, adminId } = storeToRefs(authStore)

const {
  pagination,
  isShow,
  alertType,
  alertText,
  isLoading,
  handleUpdatePage,
  execute,
} = useListPage({
  key: 'feature-requests-list',
  fetchFn: (limit, offset) => {
    // コーディネーターは自分の提出のみ表示
    const submittedBy
      = adminType.value === AdminType.AdminTypeAdministrator
        ? undefined
        : adminId.value
    return featureRequestStore.fetchFeatureRequests(limit, offset, submittedBy)
  },
})

const handleClickRow = (id: string): void => {
  router.push(`/feature-requests/${id}`)
}

const handleClickNew = (): void => {
  router.push('/feature-requests/new')
}

await execute()
</script>

<template>
  <templates-feature-request-list
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :feature-requests="featureRequests"
    :table-items-per-page="pagination.itemsPerPage.value"
    :table-items-total="total"
    @click:row="handleClickRow"
    @click:new="handleClickNew"
    @click:update-page="handleUpdatePage"
    @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
  />
</template>
