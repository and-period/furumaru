<script lang="ts" setup>
import dayjs from 'dayjs'
import { storeToRefs } from 'pinia'
import { useAlert } from '~/lib/hooks'
import { useAuthStore, useScheduleStore, useShippingStore } from '~/store'
import { CreateScheduleRequest } from '~/types/api'
import { ImageUploadStatus } from '~/types/props'

const router = useRouter()
const authStore = useAuthStore()
const scheduleStore = useScheduleStore()
const shippingStore = useShippingStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const { auth } = storeToRefs(authStore)
const { shippings } = storeToRefs(shippingStore)

const loading = ref<boolean>(false)
const formData = ref<CreateScheduleRequest>({
  coordinatorId: '',
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
const thumbnailUploadStatus = ref<ImageUploadStatus>({
  error: false,
  message: ''
})
const imageUploadStatus = ref<ImageUploadStatus>({
  error: false,
  message: ''
})
const openingVideoUploadStatus = ref<ImageUploadStatus>({
  error: false,
  message: ''
})

const fetchState = useAsyncData(async (): Promise<void> => {
  await shippingStore.fetchShippings()
})

const isLoading = (): boolean => {
  return fetchState?.pending?.value || loading.value
}

const handleSearchShipping = async (name: string): Promise<void> => {
  try {
    await shippingStore.searchCoordinators(name)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleUploadThumbnail = (files: FileList): void => {
  if (files.length === 0) {
    return
  }

  loading.value = true
  scheduleStore.uploadScheduleThumbnail(files[0])
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

const handleUploadImage = (files: FileList): void => {
  if (files.length === 0) {
    return
  }

  loading.value = true
  scheduleStore.uploadScheduleImage(files[0])
    .then((res) => {
      formData.value.imageUrl = res.url
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
  if (files.length === 0) {
    return
  }

  loading.value = true
  scheduleStore.uploadScheduleOpeningVideo(files[0])
    .then((res) => {
      formData.value.openingVideoUrl = res.url
    })
    .catch(() => {
      openingVideoUploadStatus.value.error = true
      openingVideoUploadStatus.value.message = 'アップロードに失敗しました。'
    })
    .finally(() => {
      loading.value = false
    })
}

const handleSubmit = async (): Promise<void> => {
  try {
    loading.value = true
    const req: CreateScheduleRequest = {
      ...formData.value,
      coordinatorId: auth.value?.adminId || ''
    }
    const schedule = await scheduleStore.createSchedule(req)
    router.push(`/schedules/${schedule.id}`)
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
</script>

<template>
  <templates-schedule-new
    v-model:form-data="formData"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :shippings="shippings"
    :thumbnail-upload-status="thumbnailUploadStatus"
    :image-upload-status="imageUploadStatus"
    :opening-video-upload-status="openingVideoUploadStatus"
    @update:thumbnail="handleUploadThumbnail"
    @update:image="handleUploadImage"
    @update:opening-video="handleUploadOpeningVideo"
    @udpate:search-shipping="handleSearchShipping"
    @submit="handleSubmit"
  />
</template>
