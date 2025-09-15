<script lang="ts" setup>
import dayjs, { unix } from 'dayjs'
import { storeToRefs } from 'pinia'
import { useAlert } from '~/lib/hooks'
import {
  useBroadcastStore,
  useCommonStore,
  useCoordinatorStore,
  useLiveStore,
  useProducerStore,
  useProductStore,
  useScheduleStore,
} from '~/store'
import type {
  AuthYoutubeBroadcastRequest,
  CreateLiveRequest,
  Live,
  UpdateLiveRequest,
  UpdateScheduleRequest,
} from '~/types/api/v1'
import type { ImageUploadStatus } from '~/types/props'

const route = useRoute()
const commonStore = useCommonStore()
const scheduleStore = useScheduleStore()
const liveStore = useLiveStore()
const broadcastStore = useBroadcastStore()
const coordinatorStore = useCoordinatorStore()
const producerStore = useProducerStore()
const productStore = useProductStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const scheduleId = route.params.id as string
const tab = route.query.tab as string

const { schedule, viewerLogs, totalViewers } = storeToRefs(scheduleStore)
const { lives } = storeToRefs(liveStore)
const { broadcast } = storeToRefs(broadcastStore)
const { coordinators } = storeToRefs(coordinatorStore)
const { producers } = storeToRefs(producerStore)
const { products } = storeToRefs(productStore)

const initialLive: Live = {
  id: '',
  scheduleId: '',
  producerId: '',
  productIds: [],
  comment: '',
  startAt: dayjs().unix(),
  endAt: dayjs().unix(),
  createdAt: 0,
  updatedAt: 0,
}

const loading = ref<boolean>(false)
const selector = ref<string>(tab ?? 'schedule')
const selectedLive = ref<Live>({ ...initialLive })
const authYoutubeUrl = ref<string>('')
const createLiveDialog = ref<boolean>(false)
const pauseDialog = ref<boolean>(false)
const liveMp4Dialog = ref<boolean>(false)
const archiveMp4Dialog = ref<boolean>(false)
const scheduleFormData = ref<UpdateScheduleRequest>({
  title: '',
  description: '',
  thumbnailUrl: '',
  imageUrl: '',
  openingVideoUrl: '',
  startAt: dayjs().unix(),
  endAt: dayjs().unix(),
})
const createLiveFormData = ref<CreateLiveRequest>({
  producerId: '',
  productIds: [],
  comment: '',
  startAt: dayjs().unix(),
  endAt: dayjs().unix(),
})
const mp4FormData = ref<File | undefined>()
const authYoutubeFormData = ref<AuthYoutubeBroadcastRequest>({
  youtubeHandle: '',
})
const thumbnailUploadStatus = ref<ImageUploadStatus>({
  error: false,
  message: '',
})
const imageUploadStatus = ref<ImageUploadStatus>({
  error: false,
  message: '',
})
const openingVideoUploadStatus = ref<ImageUploadStatus>({
  error: false,
  message: '',
})

const fetchState = useAsyncData(async (): Promise<void> => {
  try {
    await Promise.all([
      scheduleStore.getSchedule(scheduleId),
      scheduleStore.analyzeSchedule(scheduleId),
      liveStore.fetchLives(scheduleId),
      broadcastStore.getBroadcastByScheduleId(scheduleId),
    ])
    scheduleFormData.value = { ...schedule.value }
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

const updatable = (): boolean => {
  const startAt = unix(schedule.value.startAt)
  return dayjs().isBefore(startAt)
}

const handleSearchProducer = async (name: string): Promise<void> => {
  try {
    const producerIds = lives.value.map(
      (live: Live): string => live.producerId,
    )
    await producerStore.searchProducers(name, producerIds)
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleSearchProduct = async (
  producerId: string,
  name: string,
): Promise<void> => {
  try {
    const productIds: string[] = []
    lives.value.forEach((live: Live): void => {
      live.productIds.forEach((productId: string): void => {
        productIds.push(productId)
      })
    })
    await productStore.searchProducts(name, producerId, productIds)
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleUploadThumbnail = (files: FileList): void => {
  if (files.length === 0 || !files[0]) {
    return
  }

  loading.value = true
  scheduleStore
    .uploadScheduleThumbnail(files[0])
    .then((url: string) => {
      scheduleFormData.value.thumbnailUrl = url
    })
    .catch(() => {
      thumbnailUploadStatus.value.error = true
      thumbnailUploadStatus.value.message = 'アップロードに失敗しました。'
    })
    .finally(() => {
      loading.value = false
    })
}

const handleUploadImage = (files: FileList): void => {
  if (files.length === 0 || !files[0]) {
    return
  }

  loading.value = true
  scheduleStore
    .uploadScheduleImage(files[0])
    .then((url: string) => {
      scheduleFormData.value.imageUrl = url
    })
    .catch(() => {
      imageUploadStatus.value.error = true
      imageUploadStatus.value.message = 'アップロードに失敗しました。'
    })
    .finally(() => {
      loading.value = false
    })
}

const handleUploadOpeningVideo = (files: FileList): void => {
  if (files.length === 0 || !files[0]) {
    return
  }

  loading.value = true
  scheduleStore
    .uploadScheduleOpeningVideo(files[0])
    .then((url: string) => {
      scheduleFormData.value.openingVideoUrl = url
    })
    .catch(() => {
      openingVideoUploadStatus.value.error = true
      openingVideoUploadStatus.value.message = 'アップロードに失敗しました。'
    })
    .finally(() => {
      loading.value = false
    })
}

const handleClickLinkYouTube = async (): Promise<void> => {
  try {
    loading.value = true
    const authUrl: string = await broadcastStore.authYouTube(
      scheduleId,
      authYoutubeFormData.value,
    )
    authYoutubeUrl.value = authUrl

    commonStore.addSnackbar({
      message: 'YouTubeと連携用のURLを発行しました。',
      color: 'info',
    })
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

const handleClickNewLive = (): void => {
  handleSearchProducer('')

  createLiveFormData.value.producerId = ''
  createLiveFormData.value.productIds = []
  createLiveFormData.value.startAt = schedule.value.startAt
  createLiveFormData.value.endAt = schedule.value.endAt
  createLiveDialog.value = true
}

const handleSubmitUpdateSchedule = async (): Promise<void> => {
  try {
    loading.value = true
    await scheduleStore.updateSchedule(scheduleId, scheduleFormData.value)
    commonStore.addSnackbar({
      message: `${scheduleFormData.value.title}を更新しました。`,
      color: 'info',
    })
    schedule.value = { ...schedule.value, ...scheduleFormData.value }
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

const handleSubmitPublishSchedule = async (publish: boolean): Promise<void> => {
  try {
    loading.value = true
    await scheduleStore.publishSchedule(scheduleId, publish)
    let message: string
    if (publish) {
      message = `${schedule.value.title}を公開しました。`
    }
    else {
      message = `${schedule.value.title}を非公開にしましました。`
    }
    commonStore.addSnackbar({
      message,
      color: 'info',
    })
    schedule.value = { ...schedule.value, _public: publish }
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

const handleSubmitCreateLive = async (): Promise<void> => {
  try {
    loading.value = true
    await liveStore.createLive(scheduleId, createLiveFormData.value)
    createLiveDialog.value = false
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

const handleSubmitUpdateLive = async (
  liveId: string,
  formData: UpdateLiveRequest,
): Promise<void> => {
  try {
    loading.value = true
    await liveStore.updateLive(scheduleId, liveId, formData)
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

const handleSubmitDeleteLive = async (liveId: string): Promise<void> => {
  try {
    loading.value = true
    await liveStore.deleteLive(scheduleId, liveId)
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

const handleSubmitPause = async (): Promise<void> => {
  try {
    loading.value = true
    await broadcastStore.pause(scheduleId)
    pauseDialog.value = false
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

const handleSubmitUnpause = async (): Promise<void> => {
  try {
    loading.value = true
    await broadcastStore.unpause(scheduleId)
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

const handleSubmitActivateStaticImage = async (): Promise<void> => {
  try {
    loading.value = true
    await broadcastStore.activateStaticImage(scheduleId)
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

const handleSubmitDeactivateStaticImage = async (): Promise<void> => {
  try {
    loading.value = true
    await broadcastStore.deactivateStaticImage(scheduleId)
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

const handleSubmitChangeMp4Input = async (): Promise<void> => {
  if (!mp4FormData.value) {
    return
  }
  try {
    loading.value = true
    await broadcastStore.activateMp4Input(scheduleId, mp4FormData.value)
    liveMp4Dialog.value = false
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

const handleSubmitChangeRtmpInput = async (): Promise<void> => {
  try {
    loading.value = true
    await broadcastStore.activateRtmpInput(scheduleId)
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

const handleSubmitUploadArchiveMp4 = async (): Promise<void> => {
  if (!mp4FormData.value) {
    return
  }
  try {
    loading.value = true
    await broadcastStore.uploadArchiveMp4(scheduleId, mp4FormData.value)
    archiveMp4Dialog.value = false
    fetchState.refresh()
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
  <templates-schedule-edit
    v-model:selected-tab-item="selector"
    v-model:create-live-dialog="createLiveDialog"
    v-model:update-live-dialog="updateLiveDialog"
    v-model:pause-dialog="pauseDialog"
    v-model:live-mp4-dialog="liveMp4Dialog"
    v-model:archive-mp4-dialog="archiveMp4Dialog"
    v-model:schedule-form-data="scheduleFormData"
    v-model:create-live-form-data="createLiveFormData"
    v-model:update-live-form-data="updateLiveFormData"
    v-model:mp4-form-data="mp4FormData"
    v-model:auth-youtube-form-data="authYoutubeFormData"
    :loading="isLoading()"
    :updatable="updatable()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :schedule="schedule"
    :live="selectedLive"
    :lives="lives"
    :broadcast="broadcast"
    :coordinators="coordinators"
    :producers="producers"
    :products="products"
    :viewer-logs="viewerLogs"
    :total-viewers="totalViewers"
    :auth-youtube-url="authYoutubeUrl"
    :thumbnail-upload-status="thumbnailUploadStatus"
    :image-upload-status="imageUploadStatus"
    :opening-video-upload-status="openingVideoUploadStatus"
    @click:link-youtube="handleClickLinkYouTube"
    @click:new-live="handleClickNewLive"
    @update:thumbnail="handleUploadThumbnail"
    @update:image="handleUploadImage"
    @update:opening-video="handleUploadOpeningVideo"
    @update:public="handleSubmitPublishSchedule"
    @search:producer="handleSearchProducer"
    @search:product="handleSearchProduct"
    @submit:schedule="handleSubmitUpdateSchedule"
    @submit:create-live="handleSubmitCreateLive"
    @submit:update-live="handleSubmitUpdateLive"
    @submit:delete-live="handleSubmitDeleteLive"
    @submit:pause="handleSubmitPause"
    @submit:unpause="handleSubmitUnpause"
    @submit:activate-static-image="handleSubmitActivateStaticImage"
    @submit:deactivate-static-image="handleSubmitDeactivateStaticImage"
    @submit:change-input-mp4="handleSubmitChangeMp4Input"
    @submit:change-input-rtmp="handleSubmitChangeRtmpInput"
    @submit:upload-archive-mp4="handleSubmitUploadArchiveMp4"
  />
</template>
