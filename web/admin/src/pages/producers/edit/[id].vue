<script lang="ts" setup>
import { convertIntlToJapanesePhoneNumber, convertJapaneseToIntlPhoneNumber } from '~/lib/formatter'
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
  console.log('ここ呼ばれる？')
  const producer = await getProducer(id)
  Object.assign(formData, producer)
  formData.phoneNumber = convertIntlToJapanesePhoneNumber(producer.phoneNumber)
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
    await updateProducer(id, { ...formData, phoneNumber: convertJapaneseToIntlPhoneNumber(formData.phoneNumber) })
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
      :form-data-loading="isLoading()"
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
