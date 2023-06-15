<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { convertI18nToJapanesePhoneNumber, convertJapaneseToI18nPhoneNumber } from '~/lib/formatter'
import { useAlert, useSearchAddress } from '~/lib/hooks'
import { useCommonStore, useProducerStore } from '~/store'
import { ProducerResponse, UpdateProducerRequest } from '~/types/api'
import { ImageUploadStatus } from '~/types/props'

const route = useRoute()
const router = useRouter()
const commonStore = useCommonStore()
const producerStore = useProducerStore()
const searchAddress = useSearchAddress()
const { alertType, isShow, alertText, show } = useAlert('error')

const producerId = route.params.id as string

const { producer } = storeToRefs(producerStore)

const loading = ref<boolean>(false)
const formData = ref<ProducerResponse>({
  id: producerId,
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
  createdAt: 0,
  updatedAt: 0,
  thumbnailUrl: '',
  thumbnails: [],
  email: ''
})
const thumbnailUploadStatus = ref<ImageUploadStatus>({
  error: false,
  message: ''
})
const headerUploadStatus = ref<ImageUploadStatus>({
  error: false,
  message: ''
})

const fetchState = useAsyncData(async (): Promise<void> => {
  try {
    await producerStore.getProducer(producerId)
    formData.value = {
      ...producer.value,
      phoneNumber: convertI18nToJapanesePhoneNumber(producer.value.phoneNumber)
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

const handleSubmit = async () => {
  try {
    loading.value = true
    const req: UpdateProducerRequest = {
      ...formData.value,
      phoneNumber: convertJapaneseToI18nPhoneNumber(formData.value.phoneNumber)
    }
    await producerStore.updateProducer(producerId, req)
    commonStore.addSnackbar({
      color: 'info',
      message: `${formData.value.storeName}を更新しました。`
    })
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

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-producer-edit
    v-model:form-data="formData"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :thumbnail-upload-status="thumbnailUploadStatus"
    :header-upload-status="headerUploadStatus"
    :producer="producer"
    @click:search-address="handleSearchAddress"
    @update:thumbnail-file="handleUpdateThumbnail"
    @update:header-file="handleUpdateHeader"
    @submit="handleSubmit"
  />
</template>
