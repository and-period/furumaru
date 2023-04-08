<script lang="ts" setup>
import { useSearchAddress } from '~/lib/hooks'
import { useProducerStore } from '~/store/producer'
import { CreateProducerRequest } from '~/types/api'
import { ImageUploadStatus } from '~/types/props'

const router = useRouter()
const { createProducer, uploadProducerThumbnail, uploadProducerHeader } =
  useProducerStore()

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
  addressLine2: '',
})

const thumbnailUploadStatus = reactive<ImageUploadStatus>({
  error: false,
  message: '',
})

const headerUploadStatus = reactive<ImageUploadStatus>({
  error: false,
  message: '',
})

const handleSubmit = async () => {
  try {
    await createProducer({
      ...formData,
      phoneNumber: formData.phoneNumber.replace('0', '+81'),
    })
    router.push('/producers')
  } catch (error) {
    console.log(error)
  }
}

const handleUpdateThumbnail = (files: FileList) => {
  if (files.length > 0) {
    uploadProducerThumbnail(files[0])
      .then((res) => {
        formData.thumbnailUrl = res.url
      })
      .catch(() => {
        thumbnailUploadStatus.error = true
        thumbnailUploadStatus.message = 'アップロードに失敗しました。'
      })
  }
}

const handleUpdateHeader = async (files: FileList) => {
  if (files.length > 0) {
    await uploadProducerHeader(files[0])
      .then((res) => {
        formData.headerUrl = res.url
      })
      .catch(() => {
        headerUploadStatus.error = true
        headerUploadStatus.message = 'アップロードに失敗しました。'
      })
  }
}

const {
  loading: searchLoading,
  errorMessage: searchErrorMessage,
  searchAddressByPostalCode,
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
  <the-producer-create-form-page
    :form-data="formData"
    :thumbnail-upload-status="thumbnailUploadStatus"
    :header-upload-status="headerUploadStatus"
    :search-loading="searchLoading"
    :search-error-message="searchErrorMessage"
    @update:thumbnailFile="handleUpdateThumbnail"
    @update:headerFile="handleUpdateHeader"
    @submit="handleSubmit"
    @click:search="searchAddress"
  />
</template>
