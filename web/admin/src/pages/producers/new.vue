<script lang="ts" setup>
import { convertJapaneseToI18nPhoneNumber } from '~/lib/formatter'
import { useAlert, useSearchAddress } from '~/lib/hooks'
import { useProducerStore } from '~/store'
import { CreateProducerRequest } from '~/types/api'
import { ImageUploadStatus } from '~/types/props'

const router = useRouter()
const producerStore = useProducerStore()
const searchAddress = useSearchAddress()
const { alertType, isShow, alertText, show } = useAlert('error')

const loading = ref<boolean>(false)
const formData = ref<CreateProducerRequest>({
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
const thumbnailUploadStatus = ref<ImageUploadStatus>({
  error: false,
  message: ''
})
const headerUploadStatus = ref<ImageUploadStatus>({
  error: false,
  message: ''
})

const handleSubmit = async (): Promise<void> => {
  try {
    loading.value = true
    const req: CreateProducerRequest = {
      ...formData.value,
      phoneNumber: convertJapaneseToI18nPhoneNumber(formData.value.phoneNumber)
    }
    await producerStore.createProducer(req)
    router.push('/producers')
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  } finally {
    loading.value = false
  }
}

const handleUpdateThumbnail = (files: FileList): void => {
  if (files.length === 0) {
    return
  }

  loading.value = true
  producerStore.uploadProducerThumbnail(files[0])
    .then((res) => {
      formData.value.thumbnailUrl = res.url
    })
    .catch(() => {
      thumbnailUploadStatus.value.error = true
      thumbnailUploadStatus.value.message = 'アップロードに失敗しました。'
    })
    .finally(() => {
      loading.value = false
    })
}

const handleUpdateHeader = (files: FileList): void => {
  if (files.length === 0) {
    return
  }

  loading.value = true
  producerStore.uploadProducerHeader(files[0])
    .then((res) => {
      formData.value.headerUrl = res.url
    })
    .catch(() => {
      headerUploadStatus.value.error = true
      headerUploadStatus.value.message = 'アップロードに失敗しました。'
    })
    .finally(() => {
      loading.value = false
    })
}

const handleSearchAddress = async (): Promise<void> => {
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
</script>

<template>
  <templates-producer-new
    v-model:form-data="formData"
    :loading="loading"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :thumbnail-upload-status="thumbnailUploadStatus"
    :header-upload-status="headerUploadStatus"
    @click:search-address="handleSearchAddress"
    @update:thumbnail-file="handleUpdateThumbnail"
    @update:header-file="handleUpdateHeader"
    @submit="handleSubmit"
  />
</template>
