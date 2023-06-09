<script lang="ts" setup>
import { storeToRefs } from 'pinia'

import { convertI18nToJapanesePhoneNumber, convertJapaneseToI18nPhoneNumber } from '~/lib/formatter'
import { useAlert, usePagination, useSearchAddress } from '~/lib/hooks'
import { useCoordinatorStore, useProducerStore } from '~/store'
import { RelateProducersRequest, UpdateCoordinatorRequest, UploadImageResponse } from '~/types/api'
import { ImageUploadStatus } from '~/types/props'

const route = useRoute()
const router = useRouter()
const coordinatorStore = useCoordinatorStore()
const producerStore = useProducerStore()
const relatedProducersPagination = usePagination()
const unrelatedProducersPagination = usePagination()
const searchAddress = useSearchAddress()
const { alertType, isShow, alertText, show } = useAlert('error')

const coordinatorId = route.params.id as string

const { coordinator } = storeToRefs(coordinatorStore)
const { producers: relatedProducers, totalItems } = storeToRefs(coordinatorStore)
const { producers: unrelatedProducers } = storeToRefs(producerStore)

const loading = ref<boolean>(false)
const relatedProducersDialog = ref<boolean>(false)
const selectedProducerIds = ref<string[]>([])
const formData = ref<UpdateCoordinatorRequest>({
  storeName: '',
  firstname: '',
  lastname: '',
  firstnameKana: '',
  lastnameKana: '',
  companyName: '',
  thumbnailUrl: '',
  headerUrl: '',
  twitterAccount: '',
  instagramAccount: '',
  facebookAccount: '',
  phoneNumber: '',
  postalCode: '',
  prefecture: '',
  city: '',
  addressLine1: '',
  addressLine2: ''
})
const thumbnailUploadStatus = ref<ImageUploadStatus>({
  error: false,
  message: ''
})
const headerUploadStatus = ref<ImageUploadStatus>({
  error: false,
  message: ''
})

watch(relatedProducersPagination.itemsPerPage, () => {
  fetchRelatedProducers()
})

const fetchState = useAsyncData(async (): Promise<void> => {
  try {
    await coordinatorStore.getCoordinator(coordinatorId)
    formData.value = {
      ...coordinator.value,
      phoneNumber: convertI18nToJapanesePhoneNumber(coordinator.value.phoneNumber)
    }
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
})

const isLoading = (): boolean => {
  return fetchState?.pending?.value || loading.value
}

const fetchRelatedProducers = async (): Promise<void> => {
  try {
    await coordinatorStore.fetchRelatedProducers(
      coordinatorId,
      relatedProducersPagination.itemsPerPage.value,
      relatedProducersPagination.offset.value
    )
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const fetchUnrelatedProducers = async (): Promise<void> => {
  try {
    await producerStore.fetchProducers(
      unrelatedProducersPagination.itemsPerPage.value,
      unrelatedProducersPagination.offset.value,
      'unrelated'
    )
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleUpdateRelatedProducersPage = async (page: number) => {
  relatedProducersPagination.updateCurrentPage(page)
  await fetchRelatedProducers()
}

const handleUpdateThumbnail = (files: FileList): void => {
  if (files.length === 0) {
    return
  }

  coordinatorStore.uploadCoordinatorThumbnail(files[0])
    .then((res) => {
      formData.value.thumbnailUrl = res.url
    })
    .catch(() => {
      thumbnailUploadStatus.value.error = true
      thumbnailUploadStatus.value.message = 'アップロードに失敗しました。'
    })
}

const handleUpdateHeader = (files: FileList): void => {
  if (files.length === 0) {
    return
  }

  coordinatorStore.uploadCoordinatorHeader(files[0])
    .then((res) => {
      formData.value.headerUrl = res.url
    })
    .catch(() => {
      headerUploadStatus.value.error = true
      headerUploadStatus.value.message = 'アップロードに失敗しました。'
    })
}

const handleSubmitCoordinator = async (): Promise<void> => {
  try {
    loading.value = true
    const req: UpdateCoordinatorRequest = {
      ...formData.value,
      phoneNumber: convertJapaneseToI18nPhoneNumber(formData.value.phoneNumber)
    }
    await coordinatorStore.updateCoordinator(coordinatorId, req)
    router.push('/coordinators')
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  } finally {
    loading.value = false
  }
}

const handleSubmitRelateProducers = async (): Promise<void> => {
  try {
    loading.value = true
    const req: RelateProducersRequest = {
      producerIds: selectedProducerIds.value
    }
    await coordinatorStore.relateProducers(coordinatorId, req)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  } finally {
    relatedProducersDialog.value = false
    loading.value = false
  }
}

const handleSearchAddress = async () => {
  try {
    const res = await searchAddress.searchAddressByPostalCode(Number(formData.value.postalCode))
    formData.value = { ...formData.value, ...res }
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
  Promise.all([fetchState.execute(), fetchRelatedProducers(), fetchUnrelatedProducers()])
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-coordinator-edit
    v-model:form-data="formData"
    v-model:selected-producer-ids="selectedProducerIds"
    v-model:related-producers-dialog="relatedProducersDialog"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :related-producers="relatedProducers"
    :unrelated-producers="unrelatedProducers"
    :thumbnail-upload-status="thumbnailUploadStatus"
    :header-upload-status="headerUploadStatus"
    :related-producers-table-items-per-page="relatedProducersPagination.itemsPerPage.value"
    :related-producers-table-items-total="totalItems"
    @update:thumbnail-file="handleUpdateThumbnail"
    @update:header-file="handleUpdateHeader"
    @update:related-producers-table-page="handleUpdateRelatedProducersPage"
    @update:related-producers-table-items-per-page="relatedProducersPagination.handleUpdateItemsPerPage"
    @click:search-address="handleSearchAddress"
    @submit:coordinator="handleSubmitCoordinator"
    @submit:related-producers="handleSubmitRelateProducers"
  />
</template>
