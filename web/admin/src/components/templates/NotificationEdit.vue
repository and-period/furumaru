<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import dayjs, { unix } from 'dayjs'
import { PropType } from 'nuxt/dist/app/compat/capi'
import { AlertType } from '~/lib/hooks'

import { getErrorMessage, maxLength, required } from '~/lib/validations'
import { NotificationResponse, NotificationTargetType, UpdateNotificationRequest } from '~/types/api'
import { NotificationTime } from '~/types/props'

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
  notification: {
    type: Object as PropType<NotificationResponse>,
    default: (): NotificationResponse => ({
      id: '',
      title: '',
      body: '',
      targets: [],
      public: false,
      publishedAt: dayjs().unix(),
      createdBy: '',
      creatorName: '',
      createdAt: 0,
      updatedBy: '',
      updatedAt: 0
    })
  },
  formData: {
    type: Object as PropType<UpdateNotificationRequest>,
    default: (): UpdateNotificationRequest => ({
      title: '',
      body: '',
      targets: [],
      public: false,
      publishedAt: dayjs().unix()
    })
  }
})

const emit = defineEmits<{
  (e: 'update:form-data', formData: UpdateNotificationRequest): void
  (e: 'submit'): void
}>()

const statusList = [
  { public: '公開', value: true },
  { public: '非公開', value: false }
]
const targetList = [
  { title: 'ユーザー', value: NotificationTargetType.USERS },
  { title: '生産者', value: NotificationTargetType.PRODUCERS },
  { title: 'コーディネータ', value: NotificationTargetType.COORDINATORS }
]

const formDataRules = computed(() => ({
  title: { required, maxLength: maxLength(128) },
  body: { required, maxLength: maxLength(2000) },
  targets: {},
  public: {}
}))
const timeDataRules = computed(() => ({
  publishedDate: {},
  publishedTime: {}
}))
const formDataValue = computed({
  get: (): UpdateNotificationRequest => props.formData as UpdateNotificationRequest,
  set: (formData: UpdateNotificationRequest) => emit('update:form-data', formData)
})
const timeDataValue = computed({
  get: (): NotificationTime => ({
    publishedDate: unix(props.formData.publishedAt).format('YYYY-MM-DD'),
    publishedTime: unix(props.formData.publishedAt).format('HH:mm')
  }),
  set: (timeData: NotificationTime): void => {
    const publishedAt = dayjs(`${timeData.publishedDate} ${timeData.publishedTime}`)
    formDataValue.value.publishedAt = publishedAt.unix()
  }
})

const formDataValidate = useVuelidate(formDataRules, formDataValue)
const timeDataValidate = useVuelidate(timeDataRules, timeDataValue)

const onSubmit = async (): Promise<void> => {
  const formDataValid = await formDataValidate.value.$validate()
  const timeDataValid = await timeDataValidate.value.$validate()
  if (!formDataValid || !timeDataValid) {
    return
  }

  emit('submit')
}
</script>

<template>
  <v-card>
    <v-card-title>お知らせ編集</v-card-title>

    <v-form @submit.prevent="onSubmit">
      <v-card-text>
        <v-select
          v-model="formDataValidate.public.$model"
          :error-messages="getErrorMessage(formDataValidate.public.$errors)"
          :items="statusList"
          label="ステータス"
          item-title="public"
          item-value="value"
        />
        <v-text-field
          v-model="formDataValidate.title.$model"
          :error-messages="getErrorMessage(formDataValidate.title.$errors)"
          label="タイトル"
          required
          maxlength="128"
        />
        <v-textarea
          v-model="formDataValidate.body.$model"
          :error-messages="getErrorMessage(formDataValidate.body.$errors)"
          label="本文"
          maxlength="2000"
        />
        <v-autocomplete
          v-model="formDataValidate.targets.$model"
          :error-messages="getErrorMessage(formDataValidate.targets.$errors)"
          :items="targetList"
          label="公開範囲"
          multiple
          item-title="title"
          item-value="value"
        />
        <p class="text-h6">
          投稿予約時間
        </p>
        <div class="d-flex align-center">
          <v-text-field
            v-model="timeDataValidate.publishedDate.$model"
            :error-messages="getErrorMessage(timeDataValidate.publishedDate.$errors)"
            type="date"
            class="mr-2"
            required
            variant="outlined"
          />
          <v-text-field
            v-model="timeDataValidate.publishedTime.$model"
            :error-messages="getErrorMessage(timeDataValidate.publishedTime.$errors)"
            type="time"
            required
            variant="outlined"
          />
          <p class="text-h6 mb-6 ml-4">
            〜
          </p>
          <v-spacer />
        </div>
      </v-card-text>

      <v-card-actions>
        <v-btn
          block
          :loading="loading"
          variant="outlined"
          color="primary"
          type="submit"
          class="mt-4"
        >
          更新
        </v-btn>
      </v-card-actions>
    </v-form>
  </v-card>
</template>
