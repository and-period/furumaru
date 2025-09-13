<script lang="ts" setup>
import axios from 'axios'
import type { RawAxiosRequestHeaders } from 'axios'
import { storeToRefs } from 'pinia'

import { convertJapaneseToI18nPhoneNumber } from '~/lib/formatter'
import { useAlert, useSearchAddress } from '~/lib/hooks'
import { useCommonStore, useCoordinatorStore, useProductTypeStore } from '~/store'
import { Prefecture } from '~/types'
import type { CreateCoordinatorRequest } from '~/types/api/v1'
import type { ImageUploadStatus } from '~/types/props'

const router = useRouter()
const commonStore = useCommonStore()
const coordinatorStore = useCoordinatorStore()
const productTypeStore = useProductTypeStore()
const searchAddress = useSearchAddress()
const { alertType, isShow, alertText, show } = useAlert('error')

const { productTypes } = storeToRefs(productTypeStore)

const loading = ref<boolean>(false)
const formData = ref<CreateCoordinatorRequest>({
  lastname: '',
  lastnameKana: '',
  firstname: '',
  firstnameKana: '',
  marcheName: '',
  username: '',
  email: '',
  phoneNumber: '',
  postalCode: '',
  prefectureCode: Prefecture.UNKNOWN,
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
  facebookId: '',
  businessDays: new Set<number>(),
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

const fetchState = useAsyncData(async (): Promise<void> => {
  await productTypeStore.fetchProductTypes(200)
})

const isLoading = (): boolean => {
  return fetchState?.pending?.value || loading.value
}

const handleSubmit = async (): Promise<void> => {
  try {
    loading.value = true
    const req: CreateCoordinatorRequest = {
      ...formData.value,
      phoneNumber: convertJapaneseToI18nPhoneNumber(formData.value.phoneNumber),
    }
    const res = await coordinatorStore.createCoordinator(req)
    commonStore.addSnackbar({
      message: `コーディネータの登録が完了しました。（初期パスワード：${res.password}）`,
      color: 'info',
    })
    router.push('/coordinators')
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    window.scrollTo({
      top: 0,
      behavior: 'smooth',
    })
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
  coordinatorStore.uploadCoordinatorThumbnail(files[0])
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
  coordinatorStore.uploadCoordinatorHeader(files[0])
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
  coordinatorStore.uploadCoordinatorPromotionVideo(files[0])
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
  coordinatorStore.uploadCoordinatorBonusVideo(files[0])
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

const handleSearchProductType = async (name: string): Promise<void> => {
  try {
    await productTypeStore.searchProductTypes(name, '', formData.value.productTypeIds)
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
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
  <templates-coordinator-new
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
    :product-types="productTypes"
    @click:search-address="handleSearchAddress"
    @update:search-product-type="handleSearchProductType"
    @update:thumbnail-file="handleUpdateThumbnail"
    @update:header-file="handleUpdateHeader"
    @update:promotion-video="handleUpdatePromotionVideo"
    @update:bonus-video="handleUpdateBonusVideo"
    @submit="handleSubmit"
  />
</template>
