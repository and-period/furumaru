<script lang="ts" setup>
import { useAlert, useSearchAddress } from '~/lib/hooks'
import { useCoordinatorStore } from '~/store'
import { CreateCoordinatorRequest } from '~/types/api'
import { ImageUploadStatus } from '~/types/props'

const router = useRouter()
const {
  createCoordinator,
  uploadCoordinatorThumbnail,
  uploadCoordinatorHeader
} = useCoordinatorStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const formData = reactive<CreateCoordinatorRequest>({
  lastname: '',
  lastnameKana: '',
  firstname: '',
  firstnameKana: '',
  companyName: '',
  storeName: '',
  thumbnailUrl: '',
  headerUrl: '',
  email: '',
  phoneNumber: '',
  postalCode: '',
  prefecture: '',
  city: '',
  addressLine1: '',
  addressLine2: ''
})

const thumbnailUploadStatus = reactive<ImageUploadStatus>({
  error: false,
  message: ''
})

const headerUploadStatus = reactive<ImageUploadStatus>({
  error: false,
  message: ''
})

const handleSubmit = async () => {
  try {
    await createCoordinator({
      ...formData,
      phoneNumber: formData.phoneNumber.replace('0', '+81')
    })
    router.push('/coordinators')
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleUpdateThumbnail = (files?: FileList) => {
  if (!files || files.length === 0) {
    return
  }
  uploadCoordinatorThumbnail(files[0])
    .then((res) => {
      formData.thumbnailUrl = res.url
    })
    .catch(() => {
      thumbnailUploadStatus.error = true
      thumbnailUploadStatus.message = 'アップロードに失敗しました。'
    })
}

const handleUpdateHeader = (files?: FileList) => {
  if (!files || files.length === 0) {
    return
  }
  uploadCoordinatorHeader(files[0])
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
</script>

<template>
  <templates-coordinator-new
    v-model:form-data="formData"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :thumbnail-upload-status="thumbnailUploadStatus"
    :header-upload-status="headerUploadStatus"
    :search-loading="searchLoading"
    :search-error-message="searchErrorMessage"
    @update:thumbnail-file="handleUpdateThumbnail"
    @update:header-file="handleUpdateHeader"
    @submit="handleSubmit"
    @click:search-address="searchAddress"
  />
</template>
