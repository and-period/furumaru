<script lang="ts" setup>
import { storeToRefs } from 'pinia'

import { useAlert, usePagination } from '~/lib/hooks'
import {
  useAuthStore,
  useExperienceStore,
  useExperienceTypeStore,
  useProducerStore,
} from '~/store'
import type { Experience } from '~/types/api/v1'

const router = useRouter()
const authStore = useAuthStore()
const experienceStore = useExperienceStore()
const experienceTypeStore = useExperienceTypeStore()
const producerStore = useProducerStore()
const pagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const { adminType } = storeToRefs(authStore)
const { experiences, totalItems } = storeToRefs(experienceStore)
const { experienceTypes } = storeToRefs(experienceTypeStore)
const { producers } = storeToRefs(producerStore)

const loading = ref<boolean>(false)
const selectedItemId = ref<string>('')
const isSortMode = ref<boolean>(false)
const sortableExperiences = ref<Experience[]>([])

const fetchState = useAsyncData('experiences', async (): Promise<void> => {
  await fetchExperiences()
})

watch(pagination.itemsPerPage, (): void => {
  fetchState.refresh()
})

const fetchExperiences = async (): Promise<void> => {
  try {
    await experienceStore.fetchExperiences(
      pagination.itemsPerPage.value,
      pagination.offset.value,
    )
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

const handleClickShow = (experienceId: string): void => {
  router.push(`/experiences/${experienceId}`)
}

const handleClickNew = (): void => {
  router.push('/experiences/new')
}

const handleClickCopyItem = (): void => {
  if (selectedItemId.value !== '') {
    router.push(`/experiences/new?from=${selectedItemId.value}`)
  }
}

const handleToggleSortMode = async (value: boolean): Promise<void> => {
  if (value) {
    try {
      loading.value = true
      await experienceStore.fetchAllExperiences()
      sortableExperiences.value = [...experienceStore.allExperiences]
    }
    catch (err) {
      if (err instanceof Error) {
        show(err.message)
      }
      console.log(err)
      return
    }
    finally {
      loading.value = false
    }
  }
  isSortMode.value = value
}

const handleUpdateSortableExperiences = (value: Experience[]): void => {
  sortableExperiences.value = value
}

const handleSaveSortOrder = async (): Promise<void> => {
  try {
    loading.value = true
    const experienceIds = sortableExperiences.value.map(e => e.id)
    await experienceStore.sortExperiences(experienceIds)
    isSortMode.value = false
    await fetchState.refresh()
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

const handleClickDelete = async (experienceId: string): Promise<void> => {
  try {
    loading.value = true
    await experienceStore.deleteExperience(experienceId)
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
  <templates-experience-list
    v-model:selected-item-id="selectedItemId"
    :loading="isLoading()"
    :admin-type="adminType"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :experiences="experiences"
    :experience-types="experienceTypes"
    :producers="producers"
    :table-items-per-page="pagination.itemsPerPage.value"
    :table-items-total="totalItems"
    :is-sort-mode="isSortMode"
    :sortable-experiences="sortableExperiences"
    @click:show="handleClickShow"
    @click:new="handleClickNew"
    @click:delete="handleClickDelete"
    @click:update-page="handleUpdatePage"
    @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
    @click:copy-item="handleClickCopyItem"
    @update:sort-mode="handleToggleSortMode"
    @update:sortable-experiences="handleUpdateSortableExperiences"
    @click:save-sort="handleSaveSortOrder"
  />
</template>
