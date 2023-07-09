<script lang="ts" setup>
import { storeToRefs } from 'pinia'

import { convertI18nToJapanesePhoneNumber, convertJapaneseToI18nPhoneNumber } from '~/lib/formatter'
import { useAlert, useSearchAddress } from '~/lib/hooks'
import { useCoordinatorStore, useProductTypeStore } from '~/store'
import { Prefecture, UpdateCoordinatorRequest } from '~/types/api'
import { ImageUploadStatus } from '~/types/props'

const route = useRoute()
const router = useRouter()
const coordinatorStore = useCoordinatorStore()
const productTypeStore = useProductTypeStore()
const searchAddress = useSearchAddress()
const { alertType, isShow, alertText, show } = useAlert('error')

const coordinatorId = route.params.id as string

const { coordinator } = storeToRefs(coordinatorStore)
const { productTypes } = storeToRefs(productTypeStore)

const loading = ref<boolean>(false)
const formData = ref<UpdateCoordinatorRequest>({
  lastname: '',
  lastnameKana: '',
  firstname: '',
  firstnameKana: '',
  marcheName: '',
  username: '',
  phoneNumber: '',
  postalCode: '',
  prefecture: Prefecture.UNKNOWN,
  city: '',
  addressLine1: '',
  addressLine2: '',
  profile: '',
  productTypeIds: [],
  thumbnailUrl: '',
  headerUrl: '',
  promotionVideoUrl: '',
  bonusVideoUrl: '',
  instagramId: '',
  facebookId: ''
})
const thumbnailUploadStatus = ref<ImageUploadStatus>({
  error: false,
  message: ''
})
const headerUploadStatus = ref<ImageUploadStatus>({
  error: false,
  message: ''
})
const promotionVideoUploadStatus = ref<ImageUploadStatus>({
  error: false,
  message: ''
})
const bonusVideoUploadStatus = ref<ImageUploadStatus>({
  error: false,
  message: ''
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

const handleSubmit = async (): Promise<void> => {
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
    window.scrollTo({
      top: 0,
      behavior: 'smooth'
    })
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
  coordinatorStore.uploadCoordinatorThumbnail(files[0])
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
  coordinatorStore.uploadCoordinatorHeader(files[0])
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

const handleUpdatePromotionVideo = (files: FileList): void => {
  if (files.length === 0) {
    return
  }

  loading.value = true
  coordinatorStore.uploadCoordinatorPromotionVideo(files[0])
    .then((res) => {
      formData.value.promotionVideoUrl = res.url
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
  coordinatorStore.uploadCoordinatorBonusVideo(files[0])
    .then((res) => {
      formData.value.bonusVideoUrl = res.url
    })
    .catch(() => {
      bonusVideoUploadStatus.value.error = true
      bonusVideoUploadStatus.value.message = 'アップロードに失敗しました。'
    })
    .finally(() => {
      loading.value = false
    })
}

const handleSearchAddress = async () => {
  try {
    const res = await searchAddress.searchAddressByPostalCode(formData.value.postalCode)
    formData.value = {
      ...formData.value,
      prefecture: res.prefecture,
      city: res.city,
      addressLine1: res.town
    }
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
  <templates-coordinator-edit
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
    :coordinator="coordinator"
    :product-types="productTypes"
    @click:search-address="handleSearchAddress"
    @update:thumbnail-file="handleUpdateThumbnail"
    @update:header-file="handleUpdateHeader"
    @update:promotion-video="handleUpdatePromotionVideo"
    @update:bonus-video="handleUpdateBonusVideo"
    @submit="handleSubmit"
  />
</template>
