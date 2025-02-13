<script setup lang="ts">
import { useAlert, useSearchAddress } from '~/lib/hooks'
import {
  useExperienceStore,
  useExperienceTypeStore,
  useProducerStore,
} from '~/store'

import type { UpdateExperienceRequest } from '~/types/api'

const route = useRoute()
const router = useRouter()

const experienceId = computed<string>(() => route.params.id as string)

const experienceTypeStore = useExperienceTypeStore()
const producerStore = useProducerStore()
const experienceStore = useExperienceStore()
const searchAddress = useSearchAddress()
const { alertType, isShow, alertText, show } = useAlert('error')
const { producers } = storeToRefs(producerStore)
const { experienceTypes } = storeToRefs(experienceTypeStore)

const isLoading = ref<boolean>(false)

const formData = ref<UpdateExperienceRequest>({
  title: '',
  description: '',
  public: false,
  soldOut: false,
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
  startAt: 0,
  endAt: 0,
  promotionVideoUrl: '',
  duration: 0,
  direction: '',
  businessOpenTime: '',
  businessCloseTime: '',
})

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
    searchAddress.loading.value = false
  }
}

const handleImageUpload = async (files: FileList): Promise<void> => {
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

  const thumbnailItem = formData.value.media.find(item => item.isThumbnail)
  if (thumbnailItem) {
    return
  }
  formData.value.media = formData.value.media.map((item, i): any => ({
    ...item,
    isThumbnail: i === 0,
  }))
}

const handleSubmit = async () => {
  const req = {
    ...formData.value,
    businessOpenTime: formData.value.businessOpenTime.replace(':', ''),
    businessCloseTime: formData.value.businessCloseTime.replace(':', ''),
  }
  try {
    isLoading.value = true
    await experienceStore.updateExperience(experienceId.value, req)
    router.push('/experiences')
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    isLoading.value = false
  }
}

const convertToTimeFormat = (time: string): string => {
  return time.slice(0, 2) + ':' + time.slice(2)
}

onMounted(async () => {
  isLoading.value = true
  producerStore.fetchProducers()
  experienceTypeStore.fetchExperienceTypes()
  const result = await experienceStore.fetchExperience(experienceId.value)
  formData.value = {
    ...formData.value,
    ...result.experience,
  }
  formData.value.businessOpenTime = convertToTimeFormat(formData.value.businessOpenTime)
  formData.value.businessCloseTime = convertToTimeFormat(formData.value.businessCloseTime)
  isLoading.value = false
})
</script>

<template>
  <templates-experience-new
    v-model:form-data="formData"
    :loading="isLoading"
    :search-loading="searchAddress.loading.value"
    :search-error-message="searchAddress.errorMessage.value"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :producers="producers"
    :experience-types="experienceTypes"
    @click:search-address="handleSearchAddress"
    @update:files="handleImageUpload"
    @submit="handleSubmit"
  />
</template>
