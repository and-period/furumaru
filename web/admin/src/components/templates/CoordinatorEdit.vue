<script lang="ts" setup>
import type { VTabs } from 'vuetify/lib/components/index.mjs'

import type { AlertType } from '~/lib/hooks'
import { type UpdateCoordinatorRequest, type ProductType, type Coordinator, AdminStatus, Prefecture, Weekday, type UpsertShippingRequest, type Shipping } from '~/types/api'
import type { ImageUploadStatus } from '~/types/props'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  isAlert: {
    type: Boolean,
    default: false,
  },
  alertType: {
    type: String as PropType<AlertType>,
    default: undefined,
  },
  alertText: {
    type: String,
    default: '',
  },
  coordinatorFormData: {
    type: Object as PropType<UpdateCoordinatorRequest>,
    default: (): UpdateCoordinatorRequest => ({
      lastname: '',
      lastnameKana: '',
      firstname: '',
      firstnameKana: '',
      marcheName: '',
      username: '',
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
      businessDays: [],
    }),
  },
  shippingFormData: {
    type: Object as PropType<UpsertShippingRequest>,
    default: (): UpsertShippingRequest => ({
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
    }),
  },
  coordinator: {
    type: Object as PropType<Coordinator>,
    default: (): Coordinator => ({
      id: '',
      status: AdminStatus.UNKNOWN,
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
      thumbnails: [],
      headerUrl: '',
      headers: [],
      promotionVideoUrl: '',
      bonusVideoUrl: '',
      instagramId: '',
      facebookId: '',
      createdAt: 0,
      updatedAt: 0,
      businessDays: [],
    }),
  },
  shipping: {
    type: Object as PropType<Shipping>,
    default: (): Shipping => ({
      id: '',
      isDefault: false,
      box60Rates: [],
      box60Frozen: 0,
      box80Rates: [],
      box80Frozen: 0,
      box100Rates: [],
      box100Frozen: 0,
      hasFreeShipping: false,
      freeShippingRates: 0,
      createdAt: 0,
      updatedAt: 0,
    }),
  },
  productTypes: {
    type: Array<ProductType>,
    default: () => [],
  },
  thumbnailUploadStatus: {
    type: Object,
    default: (): ImageUploadStatus => ({
      error: false,
      message: '',
    }),
  },
  headerUploadStatus: {
    type: Object,
    default: (): ImageUploadStatus => ({
      error: false,
      message: '',
    }),
  },
  promotionVideoUploadStatus: {
    type: Object,
    default: (): ImageUploadStatus => ({
      error: false,
      message: '',
    }),
  },
  bonusVideoUploadStatus: {
    type: Object,
    default: (): ImageUploadStatus => ({
      error: false,
      message: '',
    }),
  },
  searchErrorMessage: {
    type: String,
    default: '',
  },
  searchLoading: {
    type: Boolean,
    default: false,
  },
  selectedTabItem: {
    type: String,
    default: 'coordinator',
  },
})

const emit = defineEmits<{
  (e: 'update:selected-tab-item', item: string): void
  (e: 'update:coordinator-form-data', formData: UpdateCoordinatorRequest): void
  (e: 'update:shipping-form-data', formData: UpsertShippingRequest): void
  (e: 'update:thumbnail-file', files: FileList): void
  (e: 'update:header-file', files: FileList): void
  (e: 'update:promotion-video', files: FileList): void
  (e: 'update:bonus-video', files: FileList): void
  (e: 'update:search-product-type', name: string): void
  (e: 'click:search-address'): void
  (e: 'submit:coordinator'): void
  (e: 'submit:shipping'): void
}>()

const tabs: VTabs[] = [
  { title: '基本情報', value: 'coordinator' },
  { title: '配送設定', value: 'shipping' },
]

const selectedTabItemValue = computed({
  get: (): string => props.selectedTabItem,
  set: (item: string): void => emit('update:selected-tab-item', item),
})
const coordinatorFormDataValue = computed({
  get: (): UpdateCoordinatorRequest => props.coordinatorFormData,
  set: (val: UpdateCoordinatorRequest): void => emit('update:coordinator-form-data', val),
})
const shippingFormDataValue = computed({
  get: (): UpsertShippingRequest => props.shippingFormData,
  set: (val: UpsertShippingRequest): void => emit('update:shipping-form-data', val),
})

const onChangeThumbnailFile = (files?: FileList) => {
  if (!files) {
    return
  }
  emit('update:thumbnail-file', files)
}

const onChangeHeaderFile = (files?: FileList) => {
  if (!files) {
    return
  }
  emit('update:header-file', files)
}

const onChangePromotionVideo = (files?: FileList) => {
  if (!files) {
    return
  }
  emit('update:promotion-video', files)
}

const onChangeBonusVideo = (files?: FileList) => {
  if (!files) {
    return
  }
  emit('update:bonus-video', files)
}

const onSubmitCoordinator = (): void => {
  emit('submit:coordinator')
}

const onSubmitShipping = (): void => {
  emit('submit:shipping')
}

const onChangeSearchProductType = (name: string): void => {
  emit('update:search-product-type', name)
}

const onClickSearchAddress = (): void => {
  emit('click:search-address')
}
</script>

<template>
  <v-alert
    v-show="props.isAlert"
    :type="props.alertType"
    v-text="props.alertText"
  />

  <v-card class="mb-4">
    <v-card-title>コーディネーター詳細</v-card-title>

    <v-card-text>
      <v-tabs
        v-model="selectedTabItemValue"
        grow
        color="dark"
      >
        <v-tab
          v-for="item in tabs"
          :key="item.value"
          :value="item.value"
        >
          {{ item.title }}
        </v-tab>
      </v-tabs>
    </v-card-text>
  </v-card>

  <v-window v-model="selectedTabItemValue">
    <v-window-item value="coordinator">
      <organisms-coordinator-show
        v-model:form-data="coordinatorFormDataValue"
        :loading="loading"
        :coordinator="coordinator"
        :product-types="productTypes"
        :thumbnail-upload-status="thumbnailUploadStatus"
        :header-upload-status="headerUploadStatus"
        :promotion-video-upload-status="promotionVideoUploadStatus"
        :bonus-video-upload-status="bonusVideoUploadStatus"
        :search-error-message="searchErrorMessage"
        :search-loading="searchLoading"
        @update:thumbnail-file="onChangeThumbnailFile"
        @update:header-file="onChangeHeaderFile"
        @update:promotion-video="onChangePromotionVideo"
        @update:bonus-video="onChangeBonusVideo"
        @update:search-product-type="onChangeSearchProductType"
        @click:search-address="onClickSearchAddress"
        @submit="onSubmitCoordinator"
      />
    </v-window-item>

    <v-window-item value="shipping">
      <organisms-coordinator-shipping
        v-model:form-data="shippingFormDataValue"
        :loading="loading"
        :shipping="shipping"
        @submit="onSubmitShipping"
      />
    </v-window-item>
  </v-window>
</template>
