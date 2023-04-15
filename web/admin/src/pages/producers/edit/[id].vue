<script lang="ts" setup>
import { useAlert, useSearchAddress } from '~/lib/hooks'
import { useCommonStore, useProducerStore } from '~/store'
import { ProducerResponse } from '~/types/api'
import { ApiBaseError } from '~/types/exception'
import { ImageUploadStatus } from '~/types/props'

const route = useRoute()
const router = useRouter()
const id = route.params.id as string

const { getProducer } = useProducerStore()
const { addSnackbar } = useCommonStore()

const { uploadProducerThumbnail, uploadProducerHeader, updateProducer } =
  useProducerStore()

const thumbnailUploadStatus = reactive<ImageUploadStatus>({
  error: false,
  message: ''
})

const headerUploadStatus = reactive<ImageUploadStatus>({
  error: false,
  message: ''
})

const { alertType, isShow, alertText, show } = useAlert('error')

const formData = reactive<ProducerResponse>({
  id,
  coordinatorId: '',
  lastname: '',
  lastnameKana: '',
  firstname: '',
  firstnameKana: '',
  addressLine1: '',
  addressLine2: '',
  city: '',
  prefecture: '',
  phoneNumber: '',
  postalCode: '',
  storeName: '',
  headerUrl: '',
  headers: [],
  createdAt: -1,
  updatedAt: -1,
  thumbnailUrl: '',
  thumbnails: [],
  email: ''
})

const fetchState = useAsyncData(async () => {
  const producer = await getProducer(id)
  formData.coordinatorId = producer.coordinatorId
  formData.lastname = producer.lastname
  formData.lastnameKana = producer.lastnameKana
  formData.firstname = producer.firstname
  formData.firstnameKana = producer.firstnameKana
  formData.addressLine1 = producer.addressLine1
  formData.addressLine2 = producer.addressLine2
  formData.city = producer.city
  formData.prefecture = producer.prefecture
  formData.phoneNumber = producer.phoneNumber
  formData.postalCode = producer.postalCode
  formData.storeName = producer.storeName
  formData.headerUrl = producer.headerUrl
  formData.thumbnailUrl = producer.thumbnailUrl
  formData.email = producer.email
  formData.createdAt = producer.createdAt
  formData.updatedAt = producer.updatedAt
})

const handleUpdateThumbnail = (files?: FileList) => {
  if (!files || files.length === 0) {
    return
  }

  uploadProducerThumbnail(files[0])
    .then((res) => {
      formData.thumbnailUrl = res.url
    })
    .catch(() => {
      thumbnailUploadStatus.error = true
      thumbnailUploadStatus.message = 'アップロードに失敗しました。'
    })
}

const handleUpdateHeader = async (files?: FileList) => {
  if (!files || files.length === 0) {
    return
  }

  await uploadProducerHeader(files[0])
    .then((res) => {
      formData.headerUrl = res.url
    })
    .catch(() => {
      headerUploadStatus.error = true
      headerUploadStatus.message = 'アップロードに失敗しました。'
    })
}

const {
  loading: searchLoading,
  errorMessage: searchErrorMessage,
  searchAddressByPostalCode
} = useSearchAddress()

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

const isLoading = (): boolean => {
  return fetchState?.pending?.value || false
}

const handleSubmit = async () => {
  try {
    await updateProducer(id, formData)
    addSnackbar({
      color: 'info',
      message: `${formData.storeName}を更新しました。`
    })
    router.push('/producers')
  } catch (error) {
    if (error instanceof ApiBaseError) {
      show(error.message)
      window.scrollTo({
        top: 0,
        behavior: 'smooth'
      })
    }
  }
}

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <div>
    <v-alert v-model="isShow" :type="alertType" v-text="alertText" />

    <templates-producer-edit-form-page
      :form-data="formData"
      :form-data-loading="isLoading"
      :thumbnail-upload-status="thumbnailUploadStatus"
      :header-upload-status="headerUploadStatus"
      :search-loading="searchLoading"
      :search-error-message="searchErrorMessage"
      @update:thumbnail-file="handleUpdateThumbnail"
      @update:header-file="handleUpdateHeader"
      @submit="handleSubmit"
      @click:search="searchAddress"
    />
  </div>
</template>
