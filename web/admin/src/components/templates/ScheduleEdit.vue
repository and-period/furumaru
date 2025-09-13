<script lang="ts" setup>
import dayjs from 'dayjs'
import { VTabs } from 'vuetify/lib/components/index.mjs'
import type { AlertType } from '~/lib/hooks'
import {

  BroadcastStatus,

  ScheduleStatus,

} from '~/types/api/v1'
import type { Broadcast, Coordinator, CreateLiveRequest, Live, Product, Schedule, UpdateLiveRequest, UpdateScheduleRequest, Producer, AuthYoutubeBroadcastRequest, BroadcastViewerLog } from '~/types/api/v1'
import type { ImageUploadStatus } from '~/types/props'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  updatable: {
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
  pauseDialog: {
    type: Boolean,
    default: false,
  },
  createLiveDialog: {
    type: Boolean,
    default: false,
  },
  liveMp4Dialog: {
    type: Boolean,
    default: false,
  },
  archiveMp4Dialog: {
    type: Boolean,
    default: false,
  },
  scheduleFormData: {
    type: Object as PropType<UpdateScheduleRequest>,
    default: (): UpdateScheduleRequest => ({
      title: '',
      description: '',
      thumbnailUrl: '',
      imageUrl: '',
      openingVideoUrl: '',
      startAt: dayjs().unix(),
      endAt: dayjs().unix(),
    }),
  },
  publicFormData: {
    type: Boolean,
    default: false,
  },
  createLiveFormData: {
    type: Object as PropType<CreateLiveRequest>,
    default: (): CreateLiveRequest => ({
      producerId: '',
      productIds: [],
      comment: '',
      startAt: dayjs().unix(),
      endAt: dayjs().unix(),
    }),
  },
  mp4FormData: {
    type: Object as PropType<File[] | undefined>,
    default: (): File[] | undefined => undefined,
  },
  authYoutubeFormData: {
    type: Object as PropType<AuthYoutubeBroadcastRequest>,
    default: (): AuthYoutubeBroadcastRequest => ({
      youtubeHandle: '',
    }),
  },
  schedule: {
    type: Object as PropType<Schedule>,
    default: (): Schedule => ({
      id: '',
      shopId: '',
      coordinatorId: '',
      title: '',
      description: '',
      status: ScheduleStatus.ScheduleStatusUnknown,
      thumbnailUrl: '',
      imageUrl: '',
      openingVideoUrl: '',
      _public: false,
      approved: false,
      startAt: dayjs().unix(),
      endAt: dayjs().unix(),
      createdAt: 0,
      updatedAt: 0,
    }),
  },
  live: {
    type: Object as PropType<Live>,
    default: (): Live => ({
      id: '',
      scheduleId: '',
      producerId: '',
      productIds: [],
      comment: '',
      startAt: dayjs().unix(),
      endAt: dayjs().unix(),
      createdAt: 0,
      updatedAt: 0,
    }),
  },
  lives: {
    type: Array<Live>,
    default: () => [],
  },
  broadcast: {
    type: Object as PropType<Broadcast>,
    default: (): Broadcast => ({
      id: '',
      scheduleId: '',
      status: BroadcastStatus.BroadcastStatusUnknown,
      inputUrl: '',
      outputUrl: '',
      archiveUrl: '',
      youtubeAccount: '',
      youtubeViewerUrl: '',
      youtubeAdminUrl: '',
      createdAt: 0,
      updatedAt: 0,
    }),
  },
  coordinators: {
    type: Array<Coordinator>,
    default: () => [],
  },
  producers: {
    type: Array<Producer>,
    default: () => [],
  },
  products: {
    type: Array<Product>,
    default: () => [],
  },
  viewerLogs: {
    type: Array<BroadcastViewerLog>,
    default: () => [],
  },
  authYoutubeUrl: {
    type: String,
    default: '',
  },
  video: {
    type: Object as PropType<HTMLVideoElement | undefined>,
    default: (): HTMLVideoElement | undefined => undefined,
  },
  selectedTabItem: {
    type: String,
    default: 'schedule',
  },
  thumbnailUploadStatus: {
    type: Object,
    default: (): ImageUploadStatus => ({
      error: false,
      message: '',
    }),
  },
  imageUploadStatus: {
    type: Object,
    default: (): ImageUploadStatus => ({
      error: false,
      message: '',
    }),
  },
  openingVideoUploadStatus: {
    type: Object,
    default: (): ImageUploadStatus => ({
      error: false,
      message: '',
    }),
  },
})

const emit = defineEmits<{
  (e: 'click:link-youtube'): void
  (e: 'click:new-live'): void
  (e: 'click:edit-live', liveId: string): void
  (e: 'update:pause-dialog', v: boolean): void
  (e: 'update:live-mp4-dialog', v: boolean): void
  (e: 'update:archive-mp4-dialog', v: boolean): void
  (e: 'update:selected-tab-item', item: string): void
  (e: 'update:schedule-form-data', formData: UpdateScheduleRequest): void
  (e: 'update:create-live-dialog', v: boolean): void
  (e: 'update:update-live-dialog', v: boolean): void
  (e: 'update:create-live-form-data', formData: CreateLiveRequest): void
  (
    e: 'update:auth-youtube-form-data',
    formData: AuthYoutubeBroadcastRequest,
  ): void
  (e: 'update:mp4-form-data', formData?: File[]): void
  (e: 'update:thumbnail', files: FileList): void
  (e: 'update:image', files: FileList): void
  (e: 'update:opening-video', files: FileList): void
  (e: 'update:public', publish: boolean): void
  (e: 'search:producer', name: string): void
  (e: 'search:product', producerId: string, name: string): void
  (e: 'submit:schedule'): void
  (e: 'submit:create-live'): void
  (e: 'submit:update-live', liveId: string, formData: UpdateLiveRequest): void
  (e: 'submit:delete-live', liveId: string): void
  (e: 'submit:pause'): void
  (e: 'submit:unpause'): void
  (e: 'submit:activate-static-image'): void
  (e: 'submit:deactivate-static-image'): void
  (e: 'submit:change-input-mp4'): void
  (e: 'submit:change-input-rtmp'): void
  (e: 'submit:upload-archive-mp4'): void
}>()

const tabs: VTabs[] = [
  { name: '基本情報', value: 'schedule' },
  { name: 'ライブスケジュール', value: 'lives' },
  { name: 'ライブ配信', value: 'streaming' },
  { name: '分析情報', value: 'analytics' },
]

const selectedTabItemValue = computed({
  get: (): string => props.selectedTabItem,
  set: (item: string): void => emit('update:selected-tab-item', item),
})
const scheduleFormDataValue = computed({
  get: (): UpdateScheduleRequest => props.scheduleFormData,
  set: (formData: UpdateScheduleRequest): void =>
    emit('update:schedule-form-data', formData),
})
const createLiveDialogValue = computed({
  get: (): boolean => props.createLiveDialog,
  set: (val: boolean): void => emit('update:create-live-dialog', val),
})
const pauseDialogValue = computed({
  get: (): boolean => props.pauseDialog,
  set: (val: boolean): void => emit('update:pause-dialog', val),
})
const liveMp4DialogValue = computed({
  get: (): boolean => props.liveMp4Dialog,
  set: (val: boolean): void => emit('update:live-mp4-dialog', val),
})
const archiveMp4DialogValue = computed({
  get: (): boolean => props.archiveMp4Dialog,
  set: (val: boolean): void => emit('update:archive-mp4-dialog', val),
})
const createLiveFormDataValue = computed({
  get: (): CreateLiveRequest => props.createLiveFormData,
  set: (formData: CreateLiveRequest): void =>
    emit('update:create-live-form-data', formData),
})
const mp4FormDataValue = computed({
  get: (): File[] | undefined => props.mp4FormData,
  set: (formData?: File[]): void => emit('update:mp4-form-data', formData),
})
const authYoutubeFormDataValue = computed({
  get: (): AuthYoutubeBroadcastRequest => props.authYoutubeFormData,
  set: (formData: AuthYoutubeBroadcastRequest): void =>
    emit('update:auth-youtube-form-data', formData),
})

const onClickLinkYouTube = (): void => {
  emit('click:link-youtube')
}

const onClickNewLive = (): void => {
  emit('click:new-live')
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

const onChangePublic = (publish: boolean): void => {
  emit('update:public', publish)
}

const onSearchProducer = (name: string): void => {
  emit('search:producer', name)
}

const onSearchProduct = (producerId: string, name: string): void => {
  emit('search:product', producerId, name)
}

const onSubmitSchedule = (): void => {
  emit('submit:schedule')
}

const onSubmitCreateLive = (): void => {
  emit('submit:create-live')
}

const onSubmitUpdateLive = (
  liveId: string,
  formData: UpdateLiveRequest,
): void => {
  emit('submit:update-live', liveId, formData)
}

const onSubmitDeleteLive = (liveId: string): void => {
  emit('submit:delete-live', liveId)
}

const onSubmitPause = (): void => {
  emit('submit:pause')
}

const onSubmitUnpause = (): void => {
  emit('submit:unpause')
}

const onSubmitActivateStaticImage = (): void => {
  emit('submit:activate-static-image')
}

const onSubmitDeactivateStaticImage = (): void => {
  emit('submit:deactivate-static-image')
}

const onSubmitChangeMp4Input = (): void => {
  emit('submit:change-input-mp4')
}

const onSubmitChangeRtmpInput = (): void => {
  emit('submit:change-input-rtmp')
}

const onSubmitUploadArchiveMp4 = (): void => {
  emit('submit:upload-archive-mp4')
}
</script>

<template>
  <v-alert
    v-show="props.isAlert"
    :type="props.alertType"
    class="mb-4"
    v-text="props.alertText"
  />

  <v-card class="mb-4">
    <v-card-title>ライブ配信詳細</v-card-title>

    <v-card-text>
      <v-tabs
        v-model="selectedTabItemValue"
        grow
      >
        <v-tab
          v-for="item in tabs"
          :key="item.value"
          :value="item.value"
        >
          {{ item.name }}
        </v-tab>
      </v-tabs>

      <v-window
        v-model="selectedTabItemValue"
        class="py-4"
      >
        <v-window-item value="schedule">
          <organisms-schedule-show
            v-model:form-data="scheduleFormDataValue"
            :loading="loading"
            :updatable="updatable"
            :schedule="schedule"
            :coordinators="coordinators"
            :thumbnail-upload-status="thumbnailUploadStatus"
            :image-upload-status="imageUploadStatus"
            :opening-video-upload-status="openingVideoUploadStatus"
            @update:thumbnail="onChangeThumbnailFile"
            @update:image="onChangeImageFile"
            @update:opening-video="onChangeOpeningVideo"
            @update:public="onChangePublic"
            @submit="onSubmitSchedule"
          />
        </v-window-item>

        <v-window-item value="lives">
          <organisms-live-list
            v-model:create-dialog="createLiveDialogValue"
            v-model:create-form-data="createLiveFormDataValue"
            :loading="loading"
            :live="live"
            :lives="lives"
            :schedule="schedule"
            :producers="producers"
            :products="products"
            @click:new="onClickNewLive"
            @search:producer="onSearchProducer"
            @search:product="onSearchProduct"
            @submit:create="onSubmitCreateLive"
            @submit:update="onSubmitUpdateLive"
            @submit:delete="onSubmitDeleteLive"
          />
        </v-window-item>

        <v-window-item value="streaming">
          <organisms-schedule-streaming
            v-model:pause-dialog="pauseDialogValue"
            v-model:live-mp4-dialog="liveMp4DialogValue"
            v-model:archive-mp4-dialog="archiveMp4DialogValue"
            v-model:mp4-form-data="mp4FormDataValue"
            v-model:auth-youtube-form-data="authYoutubeFormDataValue"
            :loading="loading"
            :selected-tab-item="selectedTabItem"
            :broadcast="broadcast"
            :auth-youtube-url="authYoutubeUrl"
            @click:link-youtube="onClickLinkYouTube"
            @click:activate-static-image="onSubmitActivateStaticImage"
            @click:deactivate-static-image="onSubmitDeactivateStaticImage"
            @submit:pause="onSubmitPause"
            @submit:unpause="onSubmitUnpause"
            @submit:change-input-mp4="onSubmitChangeMp4Input"
            @submit:change-input-rtmp="onSubmitChangeRtmpInput"
            @submit:upload-archive-mp4="onSubmitUploadArchiveMp4"
          />
        </v-window-item>

        <v-window-item value="analytics">
          <organisms-schedule-analytics
            :loading="loading"
            :viewer-logs="viewerLogs"
          />
        </v-window-item>
      </v-window>
    </v-card-text>
  </v-card>
</template>
