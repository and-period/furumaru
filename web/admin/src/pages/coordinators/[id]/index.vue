<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { convertI18nToJapanesePhoneNumber, convertJapaneseToI18nPhoneNumber } from '~/lib/formatter'
import { useAlert, usePagination, useSearchAddress } from '~/lib/hooks'
import { useUnsavedChangesGuard } from '~/composables/useUnsavedChangesGuard'
import { useCommonStore, useCoordinatorStore, useProducerStore, useProductTypeStore, useShippingStore, useShopStore } from '~/store'
import { Prefecture } from '~/types'
import type { UpdateCoordinatorRequest, UpdateShopRequest, Shipping, CreateShippingRequest, UpdateShippingRequest } from '~/types/api/v1'
import type { ImageUploadStatus } from '~/types/props'

const route = useRoute()
const commonStore = useCommonStore()
const coordinatorStore = useCoordinatorStore()
const producerStore = useProducerStore()
const productTypeStore = useProductTypeStore()
const searchAddress = useSearchAddress()
const shippingStore = useShippingStore()
const shopStore = useShopStore()
const pagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const { coordinator } = storeToRefs(coordinatorStore)
const { producers } = storeToRefs(producerStore)
const { productTypes } = storeToRefs(productTypeStore)
const { shippings, total } = storeToRefs(shippingStore)
const { shop } = storeToRefs(shopStore)

const coordinatorId = route.params.id as string

const loading = ref<boolean>(false)
const selector = ref<string>('coordinator')
const selectedShipping = ref<Shipping>()

const initialShippingFormData: CreateShippingRequest | UpdateShippingRequest = {
  name: '',
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
}

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
const createShippingFormData = ref<CreateShippingRequest>({ ...initialShippingFormData })
const updateShippingFormData = ref<UpdateShippingRequest>({ ...initialShippingFormData })
const createShippingDialog = ref<boolean>(false)
const updateShippingDialog = ref<boolean>(false)
const deleteShippingDialog = ref<boolean>(false)
const activeShippingDialog = ref<boolean>(false)
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

const guardedFormData = computed(() => ({
  coordinator: coordinatorFormData.value,
  shop: shopFormData.value,
}))
const { captureSnapshot, markAsSaved, showLeaveDialog, confirmLeave, cancelLeave }
  = useUnsavedChangesGuard(guardedFormData)

const fetchState = useAsyncData('coordinator-detail', async (): Promise<void> => {
  try {
    await coordinatorStore.getCoordinator(coordinatorId)
    await Promise.all([
      fetchShippings(),
      fetchShop(),
    ])

    coordinatorFormData.value = {
      ...coordinator.value,
      phoneNumber: convertI18nToJapanesePhoneNumber(coordinator.value.phoneNumber),
    }
    shopFormData.value = { ...shop.value }
    captureSnapshot()
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

watch(pagination.itemsPerPage, async (): Promise<void> => {
  loading.value = true
  await fetchShippings()
  loading.value = false
})

const fetchShippings = async (): Promise<void> => {
  try {
    await shippingStore.fetchShippings(coordinator.value.id, pagination.itemsPerPage.value, pagination.offset.value)
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const fetchShop = async (): Promise<void> => {
  try {
    await shopStore.fetchShop(coordinator.value.shopId)
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const isLoading = (): boolean => {
  return fetchState?.pending?.value || loading.value
}

const handleUpdateShippingPage = async (page: number): Promise<void> => {
  loading.value = true
  pagination.updateCurrentPage(page)
  await fetchShippings()
  loading.value = false
}

const handleClickCreateShipping = (): void => {
  createShippingDialog.value = true
}

const handleClickUpdateShipping = (shippingId: string): void => {
  const selected = shippings.value.find(s => s.id === shippingId)
  if (!selected) {
    return
  }
  selectedShipping.value = selected
  updateShippingFormData.value = { ...selected }
  updateShippingDialog.value = true
}

const handleClickDeleteShipping = (shippingId: string): void => {
  const selected = shippings.value.find(s => s.id === shippingId)
  if (!selected) {
    return
  }
  selectedShipping.value = selected
  deleteShippingDialog.value = true
}

const handleClickCopyShipping = (shippingId: string) => {
  const selected = shippings.value.find(s => s.id === shippingId)
  if (!selected) {
    return
  }
  createShippingFormData.value = { ...selected }
  createShippingDialog.value = true
}

const handleClickActiveShipping = (shippingId: string) => {
  const selected = shippings.value.find(s => s.id === shippingId)
  if (!selected) {
    return
  }
  selectedShipping.value = selected
  activeShippingDialog.value = true
}

const handleSubmitCoordinator = async (): Promise<void> => {
  try {
    loading.value = true
    const req: UpdateCoordinatorRequest = {
      ...coordinatorFormData.value,
      phoneNumber: convertJapaneseToI18nPhoneNumber(coordinatorFormData.value.phoneNumber),
    }
    await coordinatorStore.updateCoordinator(coordinatorId, req)
    markAsSaved()
    captureSnapshot()
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
    markAsSaved()
    captureSnapshot()
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

const handleSubmitCreateShipping = async (): Promise<void> => {
  try {
    loading.value = true
    await shippingStore.createShipping(coordinator.value.id, createShippingFormData.value)
    commonStore.addSnackbar({
      color: 'info',
      message: '配送設定を作成しました。',
    })
    fetchShippings()
    createShippingDialog.value = false
    createShippingFormData.value = { ...initialShippingFormData }
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

const handleSubmitUpdateShipping = async (): Promise<void> => {
  try {
    if (!selectedShipping.value) {
      throw new Error('配送設定が選択されていません。')
    }
    loading.value = true
    await shippingStore.updateShipping(coordinator.value.id, selectedShipping.value.id, updateShippingFormData.value)
    commonStore.addSnackbar({
      color: 'info',
      message: '配送設定を更新しました。',
    })
    fetchShippings()
    updateShippingDialog.value = false
    updateShippingFormData.value = { ...initialShippingFormData }
    selectedShipping.value = undefined
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

const handleSubmitDeleteShipping = async (): Promise<void> => {
  try {
    loading.value = true
    if (!selectedShipping.value) {
      throw new Error('配送設定が選択されていません。')
    }
    await shippingStore.deleteShipping(coordinator.value.id, selectedShipping.value.id)
    commonStore.addSnackbar({
      color: 'info',
      message: '配送設定を削除しました。',
    })
    fetchShippings()
    deleteShippingDialog.value = false
    selectedShipping.value = undefined
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

const handleSubmitActiveShipping = async (): Promise<void> => {
  try {
    loading.value = true
    if (!selectedShipping.value) {
      throw new Error('配送設定が選択されていません。')
    }
    await shippingStore.activeShipping(coordinator.value.id, selectedShipping.value.id)
    commonStore.addSnackbar({
      color: 'info',
      message: '配送設定を有効化しました。',
    })
    fetchShippings()
    activeShippingDialog.value = false
    selectedShipping.value = undefined
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
  if (files.length === 0 || files[0] == null) {
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
  if (files.length === 0 || files[0] == null) {
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
  if (files.length === 0 || files[0] == null) {
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
  if (files.length === 0 || files[0] == null) {
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
  <div>
    <templates-coordinator-edit
      v-model:selected-tab-item="selector"
      v-model:coordinator-form-data="coordinatorFormData"
      v-model:shop-form-data="shopFormData"
      v-model:create-shipping-form-data="createShippingFormData"
      v-model:update-shipping-form-data="updateShippingFormData"
      v-model:create-shipping-dialog="createShippingDialog"
      v-model:update-shipping-dialog="updateShippingDialog"
      v-model:delete-shipping-dialog="deleteShippingDialog"
      v-model:active-shipping-dialog="activeShippingDialog"
      :loading="isLoading()"
      :is-alert="isShow"
      :alert-type="alertType"
      :alert-text="alertText"
      :bonus-video-upload-status="bonusVideoUploadStatus"
      :header-upload-status="headerUploadStatus"
      :promotion-video-upload-status="promotionVideoUploadStatus"
      :thumbnail-upload-status="thumbnailUploadStatus"
      :search-loading="searchAddress.loading.value"
      :search-error-message="searchAddress.errorMessage.value"
      :coordinator="coordinator"
      :producers="producers"
      :product-types="productTypes"
      :shipping="selectedShipping"
      :shippings="shippings"
      :shop="shop"
      :table-items-per-page="pagination.itemsPerPage.value"
      :table-items-total="total"
      @click:search-address="handleSearchAddress"
      @click:update-shipping-page="handleUpdateShippingPage"
      @click:update-shipping-items-per-page="pagination.handleUpdateItemsPerPage"
      @click:create-shipping="handleClickCreateShipping"
      @click:update-shipping="handleClickUpdateShipping"
      @click:delete-shipping="handleClickDeleteShipping"
      @click:copy-shipping="handleClickCopyShipping"
      @click:active-shipping="handleClickActiveShipping"
      @update:search-product-type="handleSearchProductType"
      @update:thumbnail-file="handleUpdateThumbnail"
      @update:header-file="handleUpdateHeader"
      @update:promotion-video="handleUpdatePromotionVideo"
      @update:bonus-video="handleUpdateBonusVideo"
      @submit:coordinator="handleSubmitCoordinator"
      @submit:shop="handleSubmitShop"
      @submit:create-shipping="handleSubmitCreateShipping"
      @submit:update-shipping="handleSubmitUpdateShipping"
      @submit:delete-shipping="handleSubmitDeleteShipping"
      @submit:active-shipping="handleSubmitActiveShipping"
    />
    <atoms-app-confirm-dialog
      v-model="showLeaveDialog"
      title="未保存の変更があります"
      message="ページを離れると入力内容が失われます。よろしいですか？"
      confirm-text="破棄して離れる"
      confirm-color="warning"
      @confirm="confirmLeave"
      @cancel="cancelLeave"
    />
  </div>
</template>
