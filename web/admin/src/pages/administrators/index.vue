<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useAlert, usePagination } from '~/lib/hooks'
import { useAdministratorStore, useCommonStore } from '~/store'

const router = useRouter()
const commonStore = useCommonStore()
const administratorStore = useAdministratorStore()
const pagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const { administrators, total } = storeToRefs(administratorStore)

const loading = ref<boolean>(false)
const deleteDialog = ref<boolean>(false)

const fetchState = useAsyncData(async (): Promise<void> => {
  await fetchAdministrators()
})

watch(pagination.itemsPerPage, (): void => {
  fetchState.refresh()
})

const fetchAdministrators = async (): Promise<void> => {
  try {
    await administratorStore.fetchAdministrators(pagination.itemsPerPage.value, pagination.offset.value)
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

const handleClickAdd = (): void => {
  router.push('/administrators/new')
}

const handleClickDelete = async (administratorId: string): Promise<void> => {
  try {
    loading.value = true
    await administratorStore.deleteAdministrator(administratorId)
    commonStore.addSnackbar({
      message: '管理者情報の削除が完了しました。',
      color: 'info',
    })
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    deleteDialog.value = false
    loading.value = false
  }
}

const handleClickRow = (administratorId: string): void => {
  router.push(`/administrators/${administratorId}`)
}

try {
  await fetchState.execute()
}
catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-administrator-list
    v-model:delete-dialog="deleteDialog"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :administrators="administrators"
    :table-items-per-page="pagination.itemsPerPage.value"
    :table-items-total="total"
    @click:row="handleClickRow"
    @click:add="handleClickAdd"
    @click:delete="handleClickDelete"
    @click:update-page="handleUpdatePage"
    @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
  />
</template>
