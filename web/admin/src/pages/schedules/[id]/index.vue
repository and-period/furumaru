<script lang="ts" setup>
import dayjs from 'dayjs'
import { storeToRefs } from 'pinia'
import { useAlert } from '~/lib/hooks'
import { useCoordinatorStore, useLiveStore, useProducerStore, useProductStore, useScheduleStore, useShippingStore } from '~/store'
import { CreateLiveRequest, Live, Producer, UpdateLiveRequest, UpdateScheduleRequest } from '~/types/api'
import { ImageUploadStatus } from '~/types/props'

const route = useRoute()
const router = useRouter()
const scheduleStore = useScheduleStore()
const liveStore = useLiveStore()
const coordinatorStore = useCoordinatorStore()
const producerStore = useProducerStore()
const productStore = useProductStore()
const shippingStore = useShippingStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const scheduleId = route.params.id as string
const tab = route.query.tab as string

const { schedule } = storeToRefs(scheduleStore)
const { lives } = storeToRefs(liveStore)
const { coordinators } = storeToRefs(coordinatorStore)
const { producers } = storeToRefs(producerStore)
const { products } = storeToRefs(productStore)
const { shippings } = storeToRefs(shippingStore)

const loading = ref<boolean>(false)
const selector = ref<string>(tab === 'lives' ? 'lives' : 'schedules')
const selectedLiveId = ref<string>('')
const createLiveDialog = ref<boolean>(false)
const updateLiveDialog = ref<boolean>(false)
const scheduleFormData = ref<UpdateScheduleRequest>({
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
const createLiveFormData = ref<CreateLiveRequest>({
  producerId: '',
  productIds: [],
  comment: '',
  startAt: dayjs().unix(),
  endAt: dayjs().unix()
})
const updateLiveFormData = ref<UpdateLiveRequest>({
  productIds: [],
  comment: '',
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

watch(updateLiveDialog, (): void => {
  if (updateLiveDialog) {
    return
  }
  selectedLiveId.value = ''
})

const fetchState = useAsyncData(async (): Promise<void> => {
  try {
    await Promise.all([
      scheduleStore.getSchedule(scheduleId),
      liveStore.fetchLives(scheduleId)
    ])
    scheduleFormData.value = { ...schedule.value }
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

const handleSearchShipping = async (name: string): Promise<void> => {
  try {
    await shippingStore.searchShippings(name, [scheduleFormData.value.shippingId])
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleSearchProducer = async (name: string): Promise<void> => {
  try {
    const producerIds = lives.value.map((live: Live): string => live.producerId)
    await producerStore.searchProducers(name, producerIds)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleSearchProduct = async (name: string): Promise<void> => {
  try {
    const productIds: string[] = []
    lives.value.forEach((live: Live): void => {
      live.productIds.forEach((productId: string): void => {
        productIds.push(productId)
      })
    })
    await productStore.searchProducts(name, '', productIds)
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
      scheduleFormData.value.thumbnailUrl = res.url
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
      scheduleFormData.value.imageUrl = res.url
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
      scheduleFormData.value.openingVideoUrl = res.url
    })
    .catch(() => {
      openingVideoUploadStatus.value.error = true
      openingVideoUploadStatus.value.message = 'アップロードに失敗しました。'
    })
    .finally(() => {
      loading.value = false
    })
}

const handleClickNewLive = (): void => {
  createLiveDialog.value = true
}

const handleClickEditLive = (liveId: string): void => {
  const live = lives.value.find((live: Live): boolean => {
    return live.id === liveId
  })
  if (!live) {
    return
  }
  selectedLiveId.value = liveId
  updateLiveFormData.value = { ...live }
  updateLiveDialog.value = true
}

const handleSubmitUpdateSchedule = (): void => {
  console.log('submit:update-schedule', { scheduleFormData })
}

const handleSubmitCreateLive = (): void => {
  console.log('submit:create-live', { createLiveFormData })
  createLiveDialog.value = false
}

const handleSubmitUpdateLive = (): void => {
  console.log('submit:update-live', { updateLiveFormData })
  updateLiveDialog.value = false
}

const handleSubmitDeleteLive = (): void => {
  console.log('submit:delete-live', { selectedLiveId })
  updateLiveDialog.value = false
}

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-schedule-show
    v-model:selected-tab-item="selector"
    v-model:create-live-dialog="createLiveDialog"
    v-model:update-live-dialog="updateLiveDialog"
    v-model:schedule-form-data="scheduleFormData"
    v-model:create-live-form-data="createLiveFormData"
    v-model:update-live-form-data="updateLiveFormData"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :schedule="schedule"
    :lives="lives"
    :coordinators="coordinators"
    :producers="producers"
    :products="products"
    :shippings="shippings"
    :thumbnail-upload-status="thumbnailUploadStatus"
    :image-upload-status="imageUploadStatus"
    :opening-video-upload-status="openingVideoUploadStatus"
    @click:new-live="handleClickNewLive"
    @click:edit-live="handleClickEditLive"
    @update:thumbnail="handleUploadThumbnail"
    @update:image="handleUploadImage"
    @update:opening-video="handleUploadOpeningVideo"
    @search:shipping="handleSearchShipping"
    @search:producer="handleSearchProducer"
    @search:product="handleSearchProduct"
    @submit:schedule="handleSubmitUpdateSchedule"
    @submit:create-live="handleSubmitCreateLive"
    @submit:update-live="handleSubmitUpdateLive"
    @submit:delete-live="handleSubmitDeleteLive"
  />
</template>
