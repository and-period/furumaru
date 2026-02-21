<script lang="ts" setup>
import { storeToRefs } from 'pinia'

import { convertI18nToJapanesePhoneNumber, convertJapaneseToI18nPhoneNumber } from '~/lib/formatter'
import { useAlert, useSearchAddress } from '~/lib/hooks'
import { useCommonStore, useProducerStore } from '~/store'
import { Prefecture } from '~/types'
import type { UpdateProducerRequest } from '~/types/api/v1'
import type { ImageUploadStatus } from '~/types/props'

const route = useRoute()
const router = useRouter()
const commonStore = useCommonStore()
const producerStore = useProducerStore()
const searchAddress = useSearchAddress()
const { alertType, isShow, alertText, show } = useAlert('error')

const producerId = route.params.id as string

const { producer } = storeToRefs(producerStore)

const loading = ref<boolean>(false)
const formData = ref<UpdateProducerRequest>({
  lastname: '',
  lastnameKana: '',
  firstname: '',
  firstnameKana: '',
  username: '',
  email: '',
  phoneNumber: '',
  postalCode: '',
  prefectureCode: Prefecture.UNKNOWN,
  city: '',
  addressLine1: '',
  addressLine2: '',
  profile: '',
  thumbnailUrl: '',
  headerUrl: '',
  promotionVideoUrl: '',
  bonusVideoUrl: '',
  instagramId: '',
  facebookId: '',
})
const thumbnailUploadStatus = ref<ImageUploadStatus>({
  error: false,
  message: '',
})
const headerUploadStatus = ref<ImageUploadStatus>({
  error: false,
  message: '',
})
const promotionVideoUploadStatus = ref<ImageUploadStatus>({
  error: false,
  message: '',
})
const bonusVideoUploadStatus = ref<ImageUploadStatus>({
  error: false,
  message: '',
})

const fetchState = useAsyncData('producer-detail', async (): Promise<void> => {
  try {
    await producerStore.getProducer(producerId)
    formData.value = {
      ...producer.value,
      phoneNumber: convertI18nToJapanesePhoneNumber(producer.value.phoneNumber),
    }
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
})

const isLoading = (): boolean => {
  return fetchState?.pending?.value || loading.value
}

const handleSubmit = async () => {
  try {
    loading.value = true
    const req: UpdateProducerRequest = {
      ...formData.value,
      phoneNumber: convertJapaneseToI18nPhoneNumber(formData.value.phoneNumber),
    }
    await producerStore.updateProducer(producerId, req)
    commonStore.addSnackbar({
      color: 'info',
      message: `${formData.value.username}を更新しました。`,
    })
    router.push('/producers')
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    loading.value = false
  }
}

const handleUpdateThumbnail = (files: FileList): void => {
  if (files.length === 0) {
    return
  }

  loading.value = true
  producerStore.uploadProducerThumbnail(files[0])
    .then((url: string) => {
      formData.value.thumbnailUrl = url
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
    .then((url: string) => {
      formData.value.headerUrl = url
    })
    .catch(() => {
      headerUploadStatus.value.error = true
      headerUploadStatus.value.message = 'アップロードに失敗しました。'
    })
    .finally(() => {
      loading.value = false
    })
}

const handleUpdatePromotionVideo = (files: FileList): void => {
  if (files.length === 0) {
    return
  }

  loading.value = true
  producerStore.uploadProducerPromotionVideo(files[0])
    .then((url: string) => {
      formData.value.promotionVideoUrl = url
    })
    .catch(() => {
      promotionVideoUploadStatus.value.error = true
      promotionVideoUploadStatus.value.message = 'アップロードに失敗しました。'
    })
    .finally(() => {
      loading.value = false
    })
}

const handleUpdateBonusVideo = (files: FileList): void => {
  if (files.length === 0) {
    return
  }

  loading.value = true
  producerStore.uploadProducerBonusVideo(files[0])
    .then((url: string) => {
      formData.value.bonusVideoUrl = url
    })
    .catch(() => {
      bonusVideoUploadStatus.value.error = true
      bonusVideoUploadStatus.value.message = 'アップロードに失敗しました。'
    })
    .finally(() => {
      loading.value = false
    })
}

const handleSearchAddress = async (): Promise<void> => {
  try {
    searchAddress.loading.value = true
    searchAddress.errorMessage.value = ''
    const res = await searchAddress.searchAddressByPostalCode(formData.value.postalCode)
    formData.value = {
      ...formData.value,
      prefectureCode: res.prefecture,
      city: res.city,
      addressLine1: res.town,
    }
  }
  catch (err) {
    console.log(err)
  }
  finally {
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
  <templates-producer-edit
    v-model:form-data="formData"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :thumbnail-upload-status="thumbnailUploadStatus"
    :header-upload-status="headerUploadStatus"
    :promotion-video-upload-status="promotionVideoUploadStatus"
    :bonus-video-upload-status="bonusVideoUploadStatus"
    :search-loading="searchAddress.loading.value"
    :search-error-message="searchAddress.errorMessage.value"
    :producer="producer"
    @click:search-address="handleSearchAddress"
    @update:thumbnail-file="handleUpdateThumbnail"
    @update:header-file="handleUpdateHeader"
    @update:promotion-video="handleUpdatePromotionVideo"
    @update:bonus-video="handleUpdateBonusVideo"
    @submit="handleSubmit"
  />
</template>
