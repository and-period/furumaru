<script lang="ts" setup>
import dayjs from 'dayjs'
import { storeToRefs } from 'pinia'

import { useAlert, useSearchAddress } from '~/lib/hooks'
import {
  useAuthStore,
  useCommonStore,
  useExperienceStore,
  useExperienceTypeStore,
  useProducerStore,
} from '~/store'
import type { CreateExperienceRequest } from '~/types/api/v1'

const router = useRouter()
const commonStore = useCommonStore()
const authStore = useAuthStore()
const producerStore = useProducerStore()
const experienceStore = useExperienceStore()
const experienceTypeStore = useExperienceTypeStore()
const searchAddress = useSearchAddress()
const { alertType, isShow, alertText, show } = useAlert('error')

const { auth } = storeToRefs(authStore)
const { producers } = storeToRefs(producerStore)
const { experienceTypes } = storeToRefs(experienceTypeStore)

const loading = ref<boolean>(false)
const formData = ref<CreateExperienceRequest>({
  title: '',
  description: '',
  _public: false,
  soldOut: false,
  coordinatorId: '',
  producerId: '',
  experienceTypeId: '',
  media: [],
  priceAdult: 0,
  priceJuniorHighSchool: 0,
  priceElementarySchool: 0,
  pricePreschool: 0,
  priceSenior: 0,
  recommendedPoint1: '',
  recommendedPoint2: '',
  recommendedPoint3: '',
  hostPostalCode: '',
  hostPrefectureCode: 0,
  hostCity: '',
  hostAddressLine1: '',
  hostAddressLine2: '',
  startAt: dayjs().unix(),
  endAt: dayjs().unix(),
  promotionVideoUrl: '',
  duration: 0,
  direction: '',
  businessOpenTime: '',
  businessCloseTime: '',
})

onMounted(() => {
  producerStore.fetchProducers(20, 0, '')
  fetchExperienceTypes()
})

const fetchExperienceTypes = async (): Promise<void> => {
  try {
    await experienceTypeStore.fetchExperienceTypes()
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const videoUploading = ref<boolean>(false)

/**
 * 紹介動画アップロード処理
 * @param files
 */
const handleVideoUpload = async (files: FileList): Promise<void> => {
  if (!files) {
    return
  }
  try {
    videoUploading.value = true
    const url: string = await experienceStore.uploadExperienceMedia(files[0])
    formData.value.promotionVideoUrl = url
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    videoUploading.value = false
  }
}

const handleSubmit = async (): Promise<void> => {
  const req = {
    ...formData.value,
    coordinatorId: auth.value?.adminId || '',
    businessOpenTime: formData.value.businessOpenTime.replace(':', ''),
    businessCloseTime: formData.value.businessCloseTime.replace(':', ''),
  }
  try {
    loading.value = true
    await experienceStore.createExperience(req)
    commonStore.addSnackbar({
      message: `体験の登録が完了しました。`,
      color: 'info',
    })
    router.push('/experiences')
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

const handleImageUpload = async (files: FileList): Promise<void> => {
  loading.value = true
  for (const [index, file] of Array.from(files).entries()) {
    try {
      const url: string = await experienceStore.uploadExperienceMedia(file)
      formData.value.media.push({ url, isThumbnail: index === 0 })
    }
    catch (err) {
      if (err instanceof Error) {
        show(err.message)
      }
      console.log(err)
    }
  }
  loading.value = false

  const thumbnailItem = formData.value.media.find(item => item.isThumbnail)
  if (thumbnailItem) {
    return
  }
  formData.value.media = formData.value.media.map((item, i): any => ({
    ...item,
    isThumbnail: i === 0,
  }))
}

const handleSearchAddress = async (): Promise<void> => {
  try {
    searchAddress.loading.value = true
    searchAddress.errorMessage.value = ''
    const res = await searchAddress.searchAddressByPostalCode(
      formData.value.hostPostalCode,
    )
    formData.value = {
      ...formData.value,
      hostPrefectureCode: res.prefecture,
      hostCity: res.city,
      hostAddressLine1: res.town,
    }
  }
  catch (err) {
    console.log(err)
  }
  finally {
    loading.value = false
  }
}

const isLoading = (): boolean => {
  return false
}
</script>

<template>
  <templates-experience-new
    v-model:form-data="formData"
    :loading="isLoading()"
    :search-loading="searchAddress.loading.value"
    :search-error-message="searchAddress.errorMessage.value"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :producers="producers"
    :experience-types="experienceTypes"
    :video-uploading="videoUploading"
    @click:search-address="handleSearchAddress"
    @update:files="handleImageUpload"
    @update:video="handleVideoUpload"
    @submit="handleSubmit"
  />
</template>
