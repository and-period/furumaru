<script lang="ts" setup>
import dayjs from 'dayjs';
import { VTabs } from 'vuetify/lib/components/index.mjs';
import { AlertType } from '~/lib/hooks';
import { Coordinator, CreateLiveRequest, Live, Product, Schedule, ScheduleStatus, Shipping, UpdateLiveRequest, UpdateScheduleRequest } from '~/types/api';
import { ImageUploadStatus } from '~/types/props';

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  isAlert: {
    type: Boolean,
    default: false
  },
  alertType: {
    type: String as PropType<AlertType>,
    default: undefined
  },
  alertText: {
    type: String,
    default: ''
  },
  scheduleFormData: {
    type: Object as PropType<UpdateScheduleRequest>,
    default: (): UpdateScheduleRequest => ({
      shippingId: '',
      title: '',
      description: '',
      thumbnailUrl: '',
      imageUrl: '',
      openingVideoUrl: '',
      public: false,
      startAt: dayjs().unix(),
      endAt: dayjs().unix()
    })
  },
  createLiveFormData: {
    type: Object as PropType<CreateLiveRequest>,
    default: (): CreateLiveRequest => ({
      producerId: '',
      productIds: [],
      comment: '',
      startAt: dayjs().unix(),
      endAt: dayjs().unix()
    })
  },
  updateLiveFormData: {
    type: Object as PropType<UpdateLiveRequest>,
    default: (): UpdateLiveRequest => ({
      productIds: [],
      comment: '',
      startAt: dayjs().unix(),
      endAt: dayjs().unix()
    })
  },
  schedule: {
    type: Object as PropType<Schedule>,
    default: (): Schedule => ({
      id: '',
      coordinatorId: '',
      shippingId: '',
      title: '',
      description: '',
      status: ScheduleStatus.UNKNOWN,
      thumbnailUrl: '',
      thumbnails: [],
      imageUrl: '',
      openingVideoUrl: '',
      public: false,
      approved: false,
      startAt: dayjs().unix(),
      endAt: dayjs().unix(),
      createdAt: 0,
      updatedAt: 0,
    })
  },
  lives: {
    type: Array<Live>,
    default: () => []
  },
  coordinators: {
    type: Array<Coordinator>,
    default: () => []
  },
  producers: {
    type: Array<Product>,
    default: () => []
  },
  shippings: {
    type: Array<Shipping>,
    default: () => []
  },
  products: {
    type: Array<Product>,
    default: () => []
  },
  selectedTabItem: {
    type: String,
    default: 'schedule'
  },
  thumbnailUploadStatus: {
    type: Object,
    default: (): ImageUploadStatus => ({
      error: false,
      message: ''
    })
  },
  imageUploadStatus: {
    type: Object,
    default: (): ImageUploadStatus => ({
      error: false,
      message: ''
    })
  },
  openingVideoUploadStatus: {
    type: Object,
    default: (): ImageUploadStatus => ({
      error: false,
      message: ''
    })
  }
})

const emit = defineEmits<{
  (e: 'click:edit-live', liveId: string): void
  (e: 'update:selected-tab-item', item: string): void
  (e: 'update:schedule-form-data', formData: UpdateScheduleRequest): void
  (e: 'update:create-live-form-data', formData: CreateLiveRequest): void
  (e: 'update:update-live-form-data', formData: UpdateLiveRequest): void
  (e: 'update:thumbnail', files: FileList): void
  (e: 'update:image', files: FileList): void
  (e: 'update:opening-video', files: FileList): void
  (e: 'search:shipping', name: string): void
  (e: 'search:producer', name: string): void
  (e: 'search:product', name: string): void
  (e: 'submit:schedule'): void
  (e: 'submit:create-live'): void
  (e: 'submit:update-live'): void
}>()

const tabs: VTabs[] = [
  { name: '基本情報', value: 'schedule' },
  { name: 'ライブスケジュール', value: 'lives' }
]

const selectedTabItemValue = computed({
  get: (): string => props.selectedTabItem,
  set: (item: string): void => emit('update:selected-tab-item', item)
})
const scheduleFormDataValue = computed({
  get: (): UpdateScheduleRequest => props.scheduleFormData,
  set: (formData: UpdateScheduleRequest): void => emit('update:schedule-form-data', formData)
})
const createLiveFormDataValue = computed({
  get: (): CreateLiveRequest => props.createLiveFormData,
  set: (formData: CreateLiveRequest): void => emit('update:create-live-form-data', formData)
})
const updateLiveFormDataValue = computed({
  get: (): UpdateLiveRequest => props.updateLiveFormData,
  set: (formData: UpdateLiveRequest): void => emit('update:update-live-form-data', formData)
})

const onClickEditLive = (liveId: string): void => {
  emit('click:edit-live', liveId)
}

const onChangeThumbnailFile = (files: FileList): void => {
  emit('update:thumbnail', files)
}

const onChangeImageFile = (files: FileList): void => {
  emit('update:image', files)
}

const onChangeOpeningVideo = (files: FileList): void => {
  emit('update:opening-video', files)
}

const onSearchShipping = (name: string): void => {
  emit('search:shipping', name)
}

const onSearchProducer = (name: string): void => {
  emit('search:producer', name)
}

const onSearchProduct = (name: string): void => {
  emit('search:product', name)
}

const onSubmitSchedule = (): void => {
  emit('submit:schedule')
}

const onSubmitCreateLive = (): void => {
  emit('submit:create-live')
}

const onSubmitUpdateLive = (): void => {
  emit('submit:update-live')
}
</script>

<template>
  <v-alert v-show="props.isAlert" :type="props.alertType" v-text="props.alertText" />

  <v-card>
    <v-card-title>ライブ配信詳細</v-card-title>

    <v-tabs v-model="selectedTabItemValue" grow color="dark">
      <v-tab v-for="item in tabs" :key="item.value" :value="item.value">
        {{ item.name }}
      </v-tab>
    </v-tabs>

    <v-card-text>
      <v-window v-model="props.selectedTabItem">
        <v-window-item value="schedule">
          <organisms-schedule-show
            v-model:form-data="scheduleFormDataValue"
            :loading="loading"
            :schedule="schedule"
            :coordinators="coordinators"
            :shippings="shippings"
            :thumbnail-upload-status="thumbnailUploadStatus"
            :image-upload-status="imageUploadStatus"
            :opening-video-upload-status="openingVideoUploadStatus"
            @update:thumbnail="onChangeThumbnailFile"
            @update:image="onChangeImageFile"
            @update:opening-video="onChangeOpeningVideo"
            @search:shipping="onSearchShipping"
            @submit="onSubmitSchedule"
          />
        </v-window-item>

        <v-window-item value="lives">
          <organisms-live-list
            v-model:create-form-data="createLiveFormDataValue"
            v-model:update-form-data="updateLiveFormDataValue"
            :loading="loading"
            :lives="lives"
            :producers="producers"
            :products="products"
            @click:edit="onClickEditLive"
            @search:producer="onSearchProducer"
            @search:product="onSearchProduct"
            @submit:create="onSubmitCreateLive"
            @submit:update="onSubmitUpdateLive"
          />
        </v-window-item>
      </v-window>
    </v-card-text>
  </v-card>
</template>
