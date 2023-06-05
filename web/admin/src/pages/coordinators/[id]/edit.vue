<script lang="ts" setup>
import { useVuelidate } from '@vuelidate/core'
import { storeToRefs } from 'pinia'

import { convertJapaneseToI18nPhoneNumber } from '~/lib/formatter'
import { useAlert, usePagination, useSearchAddress } from '~/lib/hooks'
import { kana, required, tel, maxLength } from '~/lib/validations'
import { useCoordinatorStore, useProducerStore } from '~/store'
import { UpdateCoordinatorRequest, UploadImageResponse } from '~/types/api'
import { ImageUploadStatus } from '~/types/props'

const route = useRoute()
const router = useRouter()
const coordinatorStore = useCoordinatorStore()
const producerStore = useProducerStore()
const relatedProducersPagination = usePagination()
const unrelatedProducersPagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')
const { loading: searchLoading, errorMessage: searchErrorMessage, searchAddressByPostalCode } = useSearchAddress()

const coordinatorId = route.params.id as string

const relatedProducersDialog = ref<boolean>(false)
const selectedProducerIds = ref<string[]>([])

const { producers: relatedProducers, totalItems } = storeToRefs(coordinatorStore)
const { producers: unrelatedProducers } = storeToRefs(producerStore)

watch(relatedProducersPagination.itemsPerPage, () => {
  fetchRelatedProducers()
})

const thumbnailUploadStatus = reactive<ImageUploadStatus>({
  error: false,
  message: ''
})
const headerUploadStatus = reactive<ImageUploadStatus>({
  error: false,
  message: ''
})
const formData = reactive<UpdateCoordinatorRequest>({
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

const fetchState = useAsyncData(async () => {
  try {
    const coordinator = await coordinatorStore.getCoordinator(coordinatorId)
    formData.storeName = coordinator.storeName
    formData.firstname = coordinator.firstname
    formData.lastname = coordinator.lastname
    formData.firstnameKana = coordinator.firstnameKana
    formData.lastnameKana = coordinator.lastnameKana
    formData.companyName = coordinator.companyName
    formData.thumbnailUrl = coordinator.thumbnailUrl
    formData.headerUrl = coordinator.headerUrl
    formData.twitterAccount = coordinator.twitterAccount
    formData.instagramAccount = coordinator.instagramAccount
    formData.facebookAccount = coordinator.facebookAccount
    formData.phoneNumber = coordinator.phoneNumber.replace('+81', '0')
    formData.postalCode = coordinator.postalCode
    formData.prefecture = coordinator.prefecture
    formData.city = coordinator.city
    formData.addressLine1 = coordinator.addressLine1
    formData.addressLine2 = coordinator.addressLine2
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
})

const rules = computed(() => ({
  storeName: { required, maxLength: maxLength(64) },
  companyName: { required, maxLength: maxLength(64) },
  firstname: { required, maxLength: maxLength(16) },
  lastname: { required, maxLength: maxLength(16) },
  firstnameKana: { required, kana },
  lastnameKana: { required, kana },
  phoneNumber: { required, tel }
}))

const v$ = useVuelidate(rules, formData)

const fetchRelatedProducers = async () => {
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

const fetchUnrelatedProducers = async () => {
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

const searchAddress = async () => {
  searchLoading.value = true
  searchErrorMessage.value = ''
  const res = await searchAddressByPostalCode(Number(formData.postalCode))
  if (res) {
    formData.prefecture = res.prefecture
    formData.city = res.city
    formData.addressLine1 = res.addressLine1
  }
}

const handleUpdateThumbnail = (files: FileList) => {
  if (files.length === 0) {
    return
  }

  coordinatorStore.uploadCoordinatorThumbnail(files[0])
    .then((res: UploadImageResponse) => {
      formData.thumbnailUrl = res.url
    })
    .catch(() => {
      thumbnailUploadStatus.error = true
      thumbnailUploadStatus.message = 'アップロードに失敗しました。'
    })
}

const handleUpdateHeader = async (files: FileList) => {
  if (files.length === 0) {
    return
  }

  await coordinatorStore.uploadCoordinatorHeader(files[0])
    .then((res: UploadImageResponse) => {
      formData.headerUrl = res.url
    })
    .catch(() => {
      headerUploadStatus.error = true
      headerUploadStatus.message = 'アップロードに失敗しました。'
    })
}

const handleSubmitCoordinator = async (): Promise<void> => {
  try {
    const result = await v$.value.$validate()
    if (!result) {
      return
    }

    const req: UpdateCoordinatorRequest = {
      ...formData,
      phoneNumber: convertJapaneseToI18nPhoneNumber(formData.phoneNumber)
    }
    await coordinatorStore.updateCoordinator(coordinatorId, req)
    router.push('/coordinators')
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleSubmitRelateProducers = async (): Promise<void> => {
  try {
    await coordinatorStore.relateProducers(coordinatorId, { producerIds: selectedProducerIds.value })
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  } finally {
    relatedProducersDialog.value = false
  }
}

const isLoading = (): boolean => {
  return fetchState?.pending?.value || false
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
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :related-producers="relatedProducers"
    :unrelated-producers="unrelatedProducers"
    :thumbnail-upload-status="thumbnailUploadStatus"
    :header-upload-status="headerUploadStatus"
    :search-loading="searchLoading"
    :search-error-message="searchErrorMessage"
    :related-producers-table-items-per-page="relatedProducersPagination.itemsPerPage.value"
    :related-producers-table-items-total="totalItems"
    @update:thumbnail-file="handleUpdateThumbnail"
    @update:header-file="handleUpdateHeader"
    @update:related-producers-table-page="handleUpdateRelatedProducersPage"
    @update:related-producers-table-items-per-page="relatedProducersPagination.handleUpdateItemsPerPage"
    @click:search-address="searchAddress"
    @submit:coordinator="handleSubmitCoordinator"
    @submit:related-producers="handleSubmitRelateProducers"
  />
</template>
