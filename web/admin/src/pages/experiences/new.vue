<script lang="ts" setup>
import dayjs from 'dayjs'
import { storeToRefs } from 'pinia'

import { useAlert } from '~/lib/hooks'
import {
  useAuthStore,
  useProducerStore,
} from '~/store'
import type {
  CreateExperienceRequest,
} from '~/types/api'

const router = useRouter()
const authStore = useAuthStore()
const producerStore = useProducerStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const { auth } = storeToRefs(authStore)
const { producers } = storeToRefs(producerStore)

const loading = ref<boolean>(false)
const formData = ref<CreateExperienceRequest>({
  title: '',
  description: '',
  public: false,
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
})

const isLoading = (): boolean => {
  return loading.value
}
</script>

<template>
  <templates-experience-new
    v-model:form-data="formData"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :producers="producers"
  />
</template>
