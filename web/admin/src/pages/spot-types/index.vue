<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import type { VDataTable } from 'vuetify/components'

import { useAlert, usePagination } from '~/lib/hooks'
import { useAuthStore, useCommonStore, useSpotTypeStore } from '~/store'
import type { CreateSpotTypeRequest, UpdateSpotTypeRequest } from '~/types/api/v1'

const commonStore = useCommonStore()
const authStore = useAuthStore()
const spotTypeStore = useSpotTypeStore()
const pagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const { adminType } = storeToRefs(authStore)
const { spotTypes, total } = storeToRefs(spotTypeStore)

const loading = ref<boolean>(false)
const newDialog = ref<boolean>(false)
const editDialog = ref<boolean>(false)
const newFormData = ref<CreateSpotTypeRequest>({
  name: '',
})
const editFormData = ref<UpdateSpotTypeRequest>({
  name: '',
})

const fetchState = useAsyncData('spot-types', async (): Promise<void> => {
  await fetchSpotTypes()
})

watch(pagination.itemsPerPage, (): void => {
  fetchSpotTypes()
})

const fetchSpotTypes = async (): Promise<void> => {
  try {
    await spotTypeStore.fetchSpotTypes(pagination.itemsPerPage.value, pagination.offset.value)
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
  await fetchSpotTypes()
}

const handleCreate = async (): Promise<void> => {
  try {
    loading.value = true
    await spotTypeStore.createSpotType(newFormData.value)
    commonStore.addSnackbar({
      message: 'スポット種別を追加しました。',
      color: 'info',
    })
    fetchState.refresh()
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    newDialog.value = false
    loading.value = false
  }
}

const handleUpdate = async (spotTypeId: string): Promise<void> => {
  try {
    loading.value = true
    await spotTypeStore.updateSpotType(spotTypeId, editFormData.value)
    commonStore.addSnackbar({
      message: 'スポット種別を更新しました。',
      color: 'info',
    })
    fetchState.refresh()
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    editDialog.value = false
    loading.value = false
  }
}

const handleDelete = async (spotTypeId: string): Promise<void> => {
  try {
    loading.value = true
    await spotTypeStore.deleteSpotType(spotTypeId)
    commonStore.addSnackbar({
      message: 'スポット種別を削除しました。',
      color: 'info',
    })
    fetchState.refresh()
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

try {
  await fetchState.execute()
}
catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-product-tag-list
    v-model:new-form-data="newFormData"
    v-model:edit-form-data="editFormData"
    v-model:new-dialog="newDialog"
    v-model:edit-dialog="editDialog"
    :loading="isLoading()"
    :admin-type="adminType"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :product-tags="spotTypes"
    :table-items-per-page="pagination.itemsPerPage.value"
    :table-items-total="total"
    @click:update-page="handleUpdatePage"
    @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
    @submit:create="handleCreate"
    @submit:update="handleUpdate"
    @submit:delete="handleDelete"
  />
</template>
