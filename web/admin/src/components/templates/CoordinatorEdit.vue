<script lang="ts" setup>
import type { VTabs } from 'vuetify/lib/components/index.mjs'
import {
  mdiAccount,
  mdiStore,
  mdiPackageVariant,
} from '@mdi/js'

import type { AlertType } from '~/lib/hooks'
import { Prefecture } from '~/types'
import { AdminStatus } from '~/types/api/v1'
import type { UpdateCoordinatorRequest, ProductType, Coordinator, UpsertShippingRequest, Shipping, Shop, UpdateShopRequest, TimeWeekday, CreateShippingRequest, UpdateShippingRequest } from '~/types/api/v1'
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
    }),
  },
  shopFormData: {
    type: Object as PropType<UpdateShopRequest>,
    default: (): UpdateShopRequest => ({
      name: '',
      productTypeIds: [],
      businessDays: new Set<TimeWeekday>(),
    }),
  },
  createShippingFormData: {
    type: Object as PropType<CreateShippingRequest>,
    default: (): CreateShippingRequest => ({
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
    }),
  },
  updateShippingFormData: {
    type: Object as PropType<UpdateShippingRequest>,
    default: (): UpdateShippingRequest => ({
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
    }),
  },
  coordinator: {
    type: Object as PropType<Coordinator>,
    default: (): Coordinator => ({
      id: '',
      shopId: '',
      status: AdminStatus.AdminStatusUnknown,
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
      producerTotal: 0,
      createdAt: 0,
      updatedAt: 0,
    }),
  },
  shop: {
    type: Object as PropType<Shop>,
    default: (): Shop => ({
      id: '',
      name: '',
      coordinatorId: '',
      producerIds: [],
      productTypeIds: [],
      businessDays: [],
      createdAt: 0,
      updatedAt: 0,
    }),
  },
  shipping: {
    type: Object as PropType<Shipping>,
    default: (): Shipping => ({
      id: '',
      name: '',
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
  shippings: {
    type: Array<Shipping>,
    default: () => [],
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
  shippingTableItemsPerPage: {
    type: Number,
    default: 20,
  },
  shippingTableItemsTotal: {
    type: Number,
    default: 0,
  },
  createShippingDialog: {
    type: Boolean,
    default: false,
  },
  updateShippingDialog: {
    type: Boolean,
    default: false,
  },
  deleteShippingDialog: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits<{
  (e: 'update:selected-tab-item', item: string): void
  (e: 'update:coordinator-form-data', formData: UpdateCoordinatorRequest): void
  (e: 'update:shop-form-data', formData: UpdateShopRequest): void
  (e: 'update:create-shipping-form-data', formData: CreateShippingRequest): void
  (e: 'update:update-shipping-form-data', formData: UpdateShippingRequest): void
  (e: 'update:thumbnail-file', files: FileList): void
  (e: 'update:header-file', files: FileList): void
  (e: 'update:promotion-video', files: FileList): void
  (e: 'update:bonus-video', files: FileList): void
  (e: 'update:search-product-type', name: string): void
  (e: 'update:create-shipping-dialog', val: boolean): void
  (e: 'update:update-shipping-dialog', val: boolean): void
  (e: 'update:delete-shipping-dialog', val: boolean): void
  (e: 'click:search-address'): void
  (e: 'click:update-shipping-page', page: number): void
  (e: 'click:update-shipping-items-per-page', itemsPerPage: number): void
  (e: 'click:create-shipping'): void
  (e: 'click:update-shipping', shippingId: string): void
  (e: 'click:copy-shipping', shippingId: string): void
  (e: 'click:delete-shipping', shippingId: string): void
  (e: 'submit:coordinator'): void
  (e: 'submit:shop'): void
  (e: 'submit:create-shipping'): void
  (e: 'submit:update-shipping'): void
  (e: 'submit:delete-shipping'): void
}>()

const tabs = [
  { title: '基本情報', value: 'coordinator', icon: mdiAccount },
  { title: '店舗情報', value: 'shop', icon: mdiStore },
  { title: '配送設定', value: 'shipping', icon: mdiPackageVariant },
]

const selectedTabItemValue = computed({
  get: (): string => props.selectedTabItem,
  set: (item: string): void => emit('update:selected-tab-item', item),
})
const coordinatorFormDataValue = computed({
  get: (): UpdateCoordinatorRequest => props.coordinatorFormData,
  set: (val: UpdateCoordinatorRequest): void => emit('update:coordinator-form-data', val),
})
const shopFormDataValue = computed({
  get: (): UpdateShopRequest => props.shopFormData,
  set: (val: UpdateShopRequest): void => emit('update:shop-form-data', val),
})
const createShippingFormDataValue = computed({
  get: (): CreateShippingRequest => props.createShippingFormData,
  set: (val: CreateShippingRequest): void => emit('update:create-shipping-form-data', val),
})
const updateShippingFormDataValue = computed({
  get: (): UpdateShippingRequest => props.updateShippingFormData,
  set: (val: UpdateShippingRequest): void => emit('update:update-shipping-form-data', val),
})
const createShippingDialogValue = computed({
  get: (): boolean => props.createShippingDialog,
  set: (val: boolean): void => emit('update:create-shipping-dialog', val),
})
const updateShippingDialogValue = computed({
  get: (): boolean => props.updateShippingDialog,
  set: (val: boolean): void => emit('update:update-shipping-dialog', val),
})
const deleteShippingDialogValue = computed({
  get: (): boolean => props.deleteShippingDialog,
  set: (val: boolean): void => emit('update:delete-shipping-dialog', val),
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

const onChangeSearchProductType = (name: string): void => {
  emit('update:search-product-type', name)
}

const onClickSearchAddress = (): void => {
  emit('click:search-address')
}

const onClickUpdateShippingPage = (page: number): void => {
  emit('click:update-shipping-page', page)
}

const onClickUpdateShippingItemsPerPage = (itemsPerPage: number): void => {
  emit('click:update-shipping-items-per-page', itemsPerPage)
}

const onClickCreateShipping = (): void => {
  emit('click:create-shipping')
}

const onClickUpdateShipping = (shippingId: string): void => {
  emit('click:update-shipping', shippingId)
}

const onClickCopyShipping = (shippingId: string): void => {
  emit('click:copy-shipping', shippingId)
}

const onClickDeleteShipping = (shippingId: string): void => {
  emit('click:delete-shipping', shippingId)
}

const onSubmitCoordinator = (): void => {
  emit('submit:coordinator')
}

const onSubmitShop = (): void => {
  emit('submit:shop')
}

const onSubmitCreateShipping = (): void => {
  emit('submit:create-shipping')
}

const onSubmitUpdateShipping = (): void => {
  emit('submit:update-shipping')
}

const onSubmitDeleteShipping = (): void => {
  emit('submit:delete-shipping')
}
</script>

<template>
  <transition
    name="fade-slide"
    mode="out-in"
  >
    <v-alert
      v-show="props.isAlert"
      :type="props.alertType"
      class="mb-6"
      closable
      v-text="props.alertText"
    />
  </transition>

  <v-card
    class="form-section-card"
    elevation="0"
    :loading="loading"
  >
    <v-card-title class="section-header pa-0">
      <v-tabs
        v-model="selectedTabItemValue"
        class="w-100"
        density="comfortable"
      >
        <v-tab
          v-for="item in tabs"
          :key="item.value"
          :value="item.value"
          class="tab-item"
        >
          <v-icon
            :icon="item.icon"
            size="20"
            class="mr-2"
          />
          {{ item.title }}
        </v-tab>
      </v-tabs>
    </v-card-title>
    <v-card-text class="pa-0">
      <v-window
        v-model="selectedTabItemValue"
        class="tab-content"
      >
        <v-window-item value="coordinator">
          <div class="pa-6">
            <organisms-coordinator-show
              v-model:form-data="coordinatorFormDataValue"
              :loading="loading"
              :coordinator="coordinator"
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
              @click:search-address="onClickSearchAddress"
              @submit="onSubmitCoordinator"
            />
          </div>
        </v-window-item>

        <v-window-item value="shop">
          <div class="pa-6">
            <organisms-coordinator-shop
              v-model:form-data="shopFormDataValue"
              :loading="loading"
              :shop="shop"
              :product-types="productTypes"
              @update:search-product-type="onChangeSearchProductType"
              @submit="onSubmitShop"
            />
          </div>
        </v-window-item>

        <v-window-item value="shipping">
          <div class="pa-6">
            <organisms-coordinator-shipping
              v-model:create-form-data="createShippingFormDataValue"
              v-model:update-form-data="updateShippingFormDataValue"
              v-model:create-dialog="createShippingDialogValue"
              v-model:update-dialog="updateShippingDialogValue"
              v-model:delete-dialog="deleteShippingDialogValue"
              :loading="props.loading"
              :shipping="props.shipping"
              :shippings="props.shippings"
              :table-items-per-page="props.shippingTableItemsPerPage"
              :table-items-total="props.shippingTableItemsTotal"
              @click:update-page="onClickUpdateShippingPage"
              @click:update-items-per-page="onClickUpdateShippingItemsPerPage"
              @click:create="onClickCreateShipping"
              @click:update="onClickUpdateShipping"
              @click:copy="onClickCopyShipping"
              @click:delete="onClickDeleteShipping"
              @submit:create="onSubmitCreateShipping"
              @submit:update="onSubmitUpdateShipping"
              @submit:delete="onSubmitDeleteShipping"
            />
          </div>
        </v-window-item>
      </v-window>
    </v-card-text>
  </v-card>
</template>

<style scoped>
.form-section-card {
  border-radius: 16px;
  max-width: none;
  border: 1px solid rgb(0 0 0 / 8%);
  box-shadow: 0 2px 8px rgb(0 0 0 / 4%);
  transition: all 0.3s ease;
}

.form-section-card:hover {
  box-shadow: 0 4px 12px rgb(0 0 0 / 8%);
}

.section-header {
  background: linear-gradient(135deg, rgb(33 150 243 / 8%) 0%, rgb(33 150 243 / 0%) 100%);
  border-bottom: 1px solid rgb(33 150 243 / 12%);
  padding: 0;
}

.tab-item {
  text-transform: none;
  font-weight: 500;
  letter-spacing: 0.025rem;
  transition: all 0.2s ease;
}

.tab-content {
  min-height: 500px;
  animation: fadeIn 0.3s ease;
}

.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.3s ease;
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(-10px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(10px);
}
</style>
