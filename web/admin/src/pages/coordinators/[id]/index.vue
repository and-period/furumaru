<script lang="ts" setup>
import { storeToRefs } from 'pinia'

import { convertI18nToJapanesePhoneNumber, convertJapaneseToI18nPhoneNumber } from '~/lib/formatter'
import { useAlert, useSearchAddress } from '~/lib/hooks'
import { useCommonStore, useCoordinatorStore, useProductTypeStore, useShippingStore, useShopStore } from '~/store'
import { Prefecture } from '~/types/api'
import type { UpsertShippingRequest, UpdateCoordinatorRequest, UpdateShopRequest } from '~/types/api'
import type { ImageUploadStatus } from '~/types/props'

const route = useRoute()
const commonStore = useCommonStore()
const coordinatorStore = useCoordinatorStore()
const productTypeStore = useProductTypeStore()
const searchAddress = useSearchAddress()
const shippingStore = useShippingStore()
const shopStore = useShopStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const coordinatorId = route.params.id as string

const { coordinator } = storeToRefs(coordinatorStore)
const { productTypes } = storeToRefs(productTypeStore)
const { shipping } = storeToRefs(shippingStore)
const { shop } = storeToRefs(shopStore)

const loading = ref<boolean>(false)
const selector = ref<string>('coordinator')

const coordinatorFormData = ref<UpdateCoordinatorRequest>({
  lastname: '',
  lastnameKana: '',
  firstname: '',
  firstnameKana: '',
  username: '',
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
const shopFormData = ref<UpdateShopRequest>({
  name: '',
  productTypeIds: [],
  businessDays: [],
})
const shippingFormData = ref<UpsertShippingRequest>({
  box60Rates: [
    {
      name: '',
      price: 0,
      prefectureCodes: [],
    },
  ],
  box60Frozen: 0,
  box80Rates: [
    {
      name: '',
      price: 0,
      prefectureCodes: [],
    },
  ],
  box80Frozen: 0,
  box100Rates: [
    {
      name: '',
      price: 0,
      prefectureCodes: [],
    },
  ],
  box100Frozen: 0,
  hasFreeShipping: false,
  freeShippingRates: 0,
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
  try {
    await Promise.all([
      coordinatorStore.getCoordinator(coordinatorId),
      shippingStore.fetchShipping(coordinatorId),
    ])
    coordinatorFormData.value = {
      ...coordinator.value,
      phoneNumber: convertI18nToJapanesePhoneNumber(coordinator.value.phoneNumber),
    }
    shopFormData.value = { ...shop.value }
    shippingFormData.value = { ...shipping.value }
    if (productTypes.value.length === 0) {
      productTypeStore.fetchProductTypes(20)
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

const handleSubmitCoordinator = async (): Promise<void> => {
  try {
    loading.value = true
    const req: UpdateCoordinatorRequest = {
      ...coordinatorFormData.value,
      phoneNumber: convertJapaneseToI18nPhoneNumber(coordinatorFormData.value.phoneNumber),
    }
    await coordinatorStore.updateCoordinator(coordinatorId, req)
    commonStore.addSnackbar({
      color: 'info',
      message: 'コーディネーター情報を更新しました。',
    })
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

const handleSubmitShop = async (): Promise<void> => {
  try {
    loading.value = true
    await shopStore.updateShop(shop.value.id, shopFormData.value)
    commonStore.addSnackbar({
      color: 'info',
      message: '店舗情報を更新しました。',
    })
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

const handleSubmitShipping = async (): Promise<void> => {
  try {
    loading.value = true
    await shippingStore.upsertShipping(coordinatorId, shippingFormData.value)
    commonStore.addSnackbar({
      color: 'info',
      message: '配送設定を更新しました。',
    })
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
      coordinatorFormData.value.thumbnailUrl = url
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
      coordinatorFormData.value.headerUrl = url
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
      coordinatorFormData.value.promotionVideoUrl = url
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
      coordinatorFormData.value.bonusVideoUrl = url
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
    await productTypeStore.searchProductTypes(name, '', shopFormData.value.productTypeIds)
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleSearchAddress = async () => {
  try {
    const res = await searchAddress.searchAddressByPostalCode(coordinatorFormData.value.postalCode)
    coordinatorFormData.value = {
      ...coordinatorFormData.value,
      prefectureCode: res.prefecture,
      city: res.city,
      addressLine1: res.town,
    }
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

try {
  await fetchState.execute()
}
catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-coordinator-edit
    v-model:selected-tab-item="selector"
    v-model:coordinator-form-data="coordinatorFormData"
    v-model:shop-form-data="shopFormData"
    v-model:shipping-form-data="shippingFormData"
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
    :shipping="shipping"
    @click:search-address="handleSearchAddress"
    @update:search-product-type="handleSearchProductType"
    @update:thumbnail-file="handleUpdateThumbnail"
    @update:header-file="handleUpdateHeader"
    @update:promotion-video="handleUpdatePromotionVideo"
    @update:bonus-video="handleUpdateBonusVideo"
    @submit:coordinator="handleSubmitCoordinator"
    @submit:shop="handleSubmitShop"
    @submit:shipping="handleSubmitShipping"
  />
</template>
