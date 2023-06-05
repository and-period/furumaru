<script lang="ts" setup>
import { useAlert, useSearchAddress } from '~/lib/hooks'
import { useProducerStore } from '~/store'
import { CreateProducerRequest } from '~/types/api'
import { ImageUploadStatus } from '~/types/props'

const router = useRouter()
const { createProducer, uploadProducerThumbnail, uploadProducerHeader } =
  useProducerStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const formData = reactive<CreateProducerRequest>({
  lastname: '',
  lastnameKana: '',
  firstname: '',
  firstnameKana: '',
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
    await createProducer({
      ...formData,
      phoneNumber: formData.phoneNumber.replace('0', '+81')
    })
    router.push('/producers')
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
</script>

<template>
  <templates-producer-new
    :form-data="formData"
    :is-alrt="isShow"
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
