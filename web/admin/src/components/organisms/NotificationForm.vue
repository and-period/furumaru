<script lang="ts" setup>
import * as dayjs from 'dayjs'

import { CreateNotificationRequest, NotificationTargetType } from '~/types/api'
import { NotificationTime } from '~/types/props'

const props = defineProps({
  formType: {
    type: String,
    default: 'create',
    validator: (value: string) => {
      return ['create', 'edit'].includes(value)
    }
  },
  formData: {
    type: Object,
    default: (): CreateNotificationRequest => ({
      title: '',
      body: '',
      targets: [],
      public: false,
      publishedAt: dayjs().unix()
    })
  },
  timeData: {
    type: Object,
    default: (): NotificationTime => ({
      publishedDate: '',
      publishedTime: ''
    })
  }
})

const emit = defineEmits<{
  (e: 'update:formData', formData: CreateNotificationRequest): void
  (e: 'update:timeData', timeData: NotificationTime): void
  (e: 'submit'): void
}>()

const formDataValue = computed({
  get: (): CreateNotificationRequest =>
    props.formData as CreateNotificationRequest,
  set: (val: CreateNotificationRequest) => emit('update:formData', val)
})

const timeDataValue = computed({
  get: (): NotificationTime => props.timeData as NotificationTime,
  set: (val: NotificationTime) => emit('update:timeData', val)
})

const btnText = computed(() => {
  return props.formType === 'create' ? '登録' : '更新'
})
const postMenu = ref<boolean>(false)

const statusList = [
  { public: '公開', value: true },
  { public: '非公開', value: false }
]

const handleSubmit = () => {
  formDataValue.value.publishedAt = dayjs(
    timeDataValue.value.publishedDate + ' ' + timeDataValue.value.publishedTime
  ).unix()
  emit('submit')
}
</script>

<template>
  <form @submit.prevent="handleSubmit">
    <v-card elevation="0">
      <v-card-text>
        <v-select
          v-model="props.formData.public"
          :items="statusList"
          label="ステータス"
          item-title="public"
          item-value="value"
        />
        <v-text-field
          v-model="props.formData.title"
          label="タイトル"
          required
          maxlength="128"
        />
        <v-textarea
          v-model="props.formData.body"
          label="本文"
          maxlength="2000"
        />
      </v-card-text>
      <v-container class="ml-2">
        <p class="text-h6">
          公開範囲
        </p>
        <v-checkbox
          v-model="props.formData.targets"
          label="ユーザー"
          :value="NotificationTargetType.USERS"
        />
        <v-checkbox
          v-model="props.formData.targets"
          label="生産者"
          :value="NotificationTargetType.PRODUCERS"
        />
        <v-checkbox
          v-model="props.formData.targets"
          label="コーディネータ"
          :value="NotificationTargetType.COORDINATORS"
        />
        <p class="text-h6">
          投稿予約時間
        </p>
        <div class="d-flex align-center justify-center">
          <v-text-field
            v-model="timeDataValue.publishedDate"
            type="date"
            class="mr-2"
            required
            variant="outlined"
          />
          <v-text-field
            v-model="timeDataValue.publishedTime"
            type="time"
            required
            variant="outlined"
          />
          <p class="text-h6 mb-6 ml-4">
            〜
          </p>
          <v-spacer />
        </div>
      </v-container>
      <v-card-actions>
        <v-btn block variant="outlined" color="primary" type="submit" class="mt-4">
          {{ btnText }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </form>
</template>
