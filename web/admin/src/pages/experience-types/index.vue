<script lang="ts" setup>
import { useAlert, usePagination } from '~/lib/hooks'
import { useAuthStore, useCommonStore, useExperienceTypeStore } from '~/store'

const commonStore = useCommonStore()
const authStore = useAuthStore()
const experienceTypeStore = useExperienceTypeStore()
const pagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const { role } = storeToRefs(authStore)
const { experienceTypes, total } = storeToRefs(experienceTypeStore)

const loading = ref<boolean>(false)
const sortBy = ref<VDataTable['sortBy']>([])
const newDialog = ref<boolean>(false)
const editDialog = ref<boolean>(false)
const deleteDialog = ref<boolean>(false)
const newFormData = ref<CreateExperienceTypeRequest>({
  name: '',
})
const editFormData = ref<UpdateExperienceTypeRequest>({
  name: '',
})

const fetchState = useAsyncData(async (): Promise<void> => {
  await fetchExperienceTypes()
})

watch(pagination.itemsPerPage, (): void => {
  fetchExperienceTypes()
})
watch(sortBy, (): void => {
  fetchState.refresh()
})

const fetchExperienceTypes = async (): Promise<void> => {
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

    await ExperienceTypeStore.fetchExperienceTypes(pagination.itemsPerPage.value, pagination.offset.value, orders)
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
  await fetchExperienceTypes()
}

const handleCreate = async (): Promise<void> => {
  try {
    loading.value = true
    await ExperienceTypeStore.createExperienceType(newFormData.value)
    commonStore.addSnackbar({
      message: '体験種別を追加しました。',
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

const handleUpdate = async (ExperienceTypeId: string): Promise<void> => {
  try {
    loading.value = true
    await ExperienceTypeStore.updateExperienceType(ExperienceTypeId, editFormData.value)
    commonStore.addSnackbar({
      message: '体験種別を更新しました。',
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

const handleDelete = async (ExperienceTypeId: string): Promise<void> => {
  try {
    loading.value = true
    await ExperienceTypeStore.deleteExperienceType(ExperienceTypeId)
    commonStore.addSnackbar({
      message: '体験種別を削除しました。',
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
    deleteDialog.value = false
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
  <templates-experience-type-list
    v-model:new-form-data="newFormData"
    v-model:edit-form-data="editFormData"
    v-model:new-dialog="newDialog"
    v-model:edit-dialog="editDialog"
    v-model:delete-dialog="deleteDialog"
    v-model:sort-by="sortBy"
    :loading="isLoading()"
    :role="role"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :experience-types="experienceTypes"
    :table-items-per-page="pagination.itemsPerPage.value"
    :table-items-total="total"
    @click:update-page="handleUpdatePage"
    @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
    @update:sort-by="fetchState.refresh"
    @submit:create="handleCreate"
    @submit:update="handleUpdate"
    @submit:delete="handleDelete"
  />
</template>
