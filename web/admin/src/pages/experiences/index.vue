<script lang="ts" setup>
import { storeToRefs } from 'pinia'

import { useAlert, usePagination } from '~/lib/hooks'
import {
  useAuthStore,
  useExperienceStore,
  useExperienceTypeStore,
  useProducerStore,
} from '~/store'

const router = useRouter()
const authStore = useAuthStore()
const experienceStore = useExperienceStore()
const experienceTypeStore = useExperienceTypeStore()
const producerStore = useProducerStore()
const pagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const { role } = storeToRefs(authStore)
const { experiencesResponse, totalItems } = storeToRefs(experienceStore)
const { experienceTypes } = storeToRefs(experienceTypeStore)
const { producers } = storeToRefs(producerStore)

const loading = ref<boolean>(false)
const deleteDialog = ref<boolean>(false)
const selectedItemId = ref<string>('')

const fetchState = useAsyncData(async (): Promise<void> => {
  await Promise.all([
    fetchExperiences(),
    fetchExperienceTypes(),
    fetchProducers(),
  ])
})

watch(pagination.itemsPerPage, (): void => {
  fetchState.refresh()
})

const fetchProducers = async (): Promise<void> => {
  try {
    await producerStore.fetchProducers(pagination.itemsPerPage.value, pagination.offset.value)
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

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

const fetchExperienceTypes = async (): Promise<void> => {
  try {
    await experienceTypeStore.fetchExperienceTypes()
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
  <templates-experience-list
    v-model:delete-dialog="deleteDialog"
    v-model:selected-item-id="selectedItemId"
    :loading="isLoading()"
    :role="role"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :producers="producers"
    :experiences-response="experiencesResponse"
    :experience-types="experienceTypes"
    :table-items-per-page="pagination.itemsPerPage.value"
    :table-items-total="totalItems"
    @click:show="handleClickShow"
    @click:new="handleClickNew"
    @click:delete="handleClickDelete"
    @click:update-page="handleUpdatePage"
    @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
  />
</template>
