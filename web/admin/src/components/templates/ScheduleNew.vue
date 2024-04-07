<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import dayjs, { unix } from 'dayjs'
import type { AlertType } from '~/lib/hooks'
import { getErrorMessage } from '~/lib/validations'
import type { CreateScheduleRequest } from '~/types/api'
import type { DateTimeInput, ImageUploadStatus } from '~/types/props'
import { CreateScheduleValidationRules, TimeDataValidationRules } from '~/types/validations'

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
  formData: {
    type: Object as PropType<CreateScheduleRequest>,
    default: (): CreateScheduleRequest => ({
      coordinatorId: '',
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
  (e: 'update:form-data', formData: CreateScheduleRequest): void
  (e: 'update:thumbnail', files: FileList): void
  (e: 'update:image', files: FileList): void
  (e: 'update:opening-video', files: FileList): void
  (e: 'submit'): void
}>()

const formDataValue = computed({
  get: (): CreateScheduleRequest => props.formData,
  set: (formData: CreateScheduleRequest): void => emit('update:form-data', formData)
})
const startTimeDataValue = computed({
  get: (): DateTimeInput => ({
    date: unix(props.formData.startAt).format('YYYY-MM-DD'),
    time: unix(props.formData.startAt).format('HH:mm')
  }),
  set: (timeData: DateTimeInput): void => {
    const startAt = dayjs(`${timeData.date} ${timeData.time}`)
    formDataValue.value.startAt = startAt.unix()
  }
})
const endTimeDataValue = computed({
  get: (): DateTimeInput => ({
    date: unix(props.formData.endAt).format('YYYY-MM-DD'),
    time: unix(props.formData.endAt).format('HH:mm')
  }),
  set: (timeData: DateTimeInput): void => {
    const endAt = dayjs(`${timeData.date} ${timeData.time}`)
    formDataValue.value.endAt = endAt.unix()
  }
})

const formDataValidate = useVuelidate(CreateScheduleValidationRules, formDataValue)
const startTimeDataValidate = useVuelidate(TimeDataValidationRules, startTimeDataValue)
const endTimeDataValidate = useVuelidate(TimeDataValidationRules, endTimeDataValue)

const onChangeStartAt = (): void => {
  const startAt = dayjs(`${startTimeDataValue.value.date} ${startTimeDataValue.value.time}`)
  formDataValue.value.startAt = startAt.unix()
}

const onChangeEndAt = (): void => {
  const endAt = dayjs(`${endTimeDataValue.value.date} ${endTimeDataValue.value.time}`)
  formDataValue.value.endAt = endAt.unix()
}

const onChangeThumbnailFile = (files?: FileList) => {
  if (!files) {
    return
  }
  emit('update:thumbnail', files)
}

const onChangeImageFile = (files?: FileList) => {
  if (!files) {
    return
  }
  emit('update:image', files)
}

const onChangeOpeningVideo = (files?: FileList) => {
  if (!files) {
    return
  }
  emit('update:opening-video', files)
}

const onSubmit = async (): Promise<void> => {
  const formDataValid = await formDataValidate.value.$validate()
  const startTimeDataValid = await startTimeDataValidate.value.$validate()
  const endTimeDataValid = await endTimeDataValidate.value.$validate()
  if (!formDataValid || !startTimeDataValid || !endTimeDataValid) {
    return
  }

  emit('submit')
}
</script>

<template>
  <v-alert v-show="props.isAlert" :type="props.alertType" v-text="props.alertText" />

  <v-card>
    <v-card-title>ライブ配信登録</v-card-title>

    <v-form @submit.prevent="onSubmit">
      <v-card-text>
        <v-text-field
          v-model="formDataValidate.title.$model"
          :error-messages="getErrorMessage(formDataValidate.title.$errors)"
          label="タイトル"
        />
        <p class="text-subtitle-2 text-grey py-2">
          開催期間
        </p>
        <div class="d-flex flex-column flex-md-row justify-center">
          <v-text-field
            v-model="startTimeDataValidate.date.$model"
            :error-messages="getErrorMessage(startTimeDataValidate.date.$errors)"
            type="date"
            variant="outlined"
            density="compact"
            class="mr-md-2"
            @update:model-value="onChangeStartAt"
          />
          <v-text-field
            v-model="startTimeDataValidate.time.$model"
            :error-messages="getErrorMessage(startTimeDataValidate.time.$errors)"
            type="time"
            variant="outlined"
            density="compact"
            @update:model-value="onChangeStartAt"
          />
          <p class="text-subtitle-2 mx-4 pt-md-3 pb-4 pb-md-6">
            〜
          </p>
          <v-text-field
            v-model="endTimeDataValidate.date.$model"
            :error-messages="getErrorMessage(endTimeDataValidate.date.$errors)"
            type="date"
            variant="outlined"
            density="compact"
            class="mr-md-2"
            @update:model-value="onChangeEndAt"
          />
          <v-text-field
            v-model="endTimeDataValidate.time.$model"
            :error-messages="getErrorMessage(endTimeDataValidate.time.$errors)"
            type="time"
            variant="outlined"
            density="compact"
            @update:model-value="onChangeEndAt"
          />
        </div>
        <v-textarea
          v-model="formDataValidate.description.$model"
          :error-message="getErrorMessage(formDataValidate.description.$errors)"
          label="詳細"
          maxlength="2000"
        />
        <v-row>
          <v-col cols="12" sm="12" md="4">
            <molecules-image-select-form
              label="サムネイル画像"
              :loading="loading"
              :img-url="formDataValue.thumbnailUrl"
              :validation-error-message="getErrorMessage(formDataValidate.thumbnailUrl.$errors)"
              :error="props.thumbnailUploadStatus.error"
              :message="props.thumbnailUploadStatus.message"
              @update:file="onChangeThumbnailFile"
            />
          </v-col>
          <v-col cols="12" sm="12" md="4">
            <molecules-video-select-form
              label="オープニング動画"
              :loading="loading"
              :video-url="formDataValue.openingVideoUrl"
              :validation-error-message="getErrorMessage(formDataValidate.openingVideoUrl.$errors)"
              :error="props.openingVideoUploadStatus.error"
              :message="props.openingVideoUploadStatus.message"
              @update:file="onChangeOpeningVideo"
            />
          </v-col>
          <v-col cols="12" sm="12" md="4">
            <molecules-image-select-form
              label="待機中の画像"
              :loading="loading"
              :img-url="formDataValue.imageUrl"
              :validation-error-message="getErrorMessage(formDataValidate.imageUrl.$errors)"
              :error="props.imageUploadStatus.error"
              :message="props.imageUploadStatus.message"
              @update:file="onChangeImageFile"
            />
          </v-col>
        </v-row>
      </v-card-text>

      <v-card-actions>
        <v-btn block :loading="props.loading" variant="outlined" color="primary" type="submit">
          登録
        </v-btn>
      </v-card-actions>
    </v-form>
  </v-card>
</template>
