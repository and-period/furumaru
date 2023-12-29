<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import dayjs, { unix } from 'dayjs'
import type { AlertType } from '~/lib/hooks'
import { getErrorMessage, maxLength, required } from '~/lib/validations'
import type { CreateScheduleRequest } from '~/types/api'
import type { ImageUploadStatus, ScheduleTime } from '~/types/props'

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

const formDataRules = computed(() => ({
  title: { required, maxLength: maxLength(200) },
  description: { required, maxLength: maxLength(2000) }
}))
const timeDataRules = computed(() => ({
  startDate: { required },
  startTime: { required },
  endDate: { required },
  endTime: { required }
}))
const formDataValue = computed({
  get: (): CreateScheduleRequest => props.formData,
  set: (formData: CreateScheduleRequest): void => emit('update:form-data', formData)
})
const timeDataValue = computed({
  get: (): ScheduleTime => ({
    startDate: unix(props.formData.startAt).format('YYYY-MM-DD'),
    startTime: unix(props.formData.startAt).format('HH:mm'),
    endDate: unix(props.formData.endAt).format('YYYY-MM-DD'),
    endTime: unix(props.formData.endAt).format('HH:mm')
  }),
  set: (timeData: ScheduleTime): void => {
    const startAt = dayjs(`${timeData.startDate} ${timeData.startTime}`)
    const endAt = dayjs(`${timeData.endDate} ${timeData.endTime}`)
    formDataValue.value.startAt = startAt.unix()
    formDataValue.value.endAt = endAt.unix()
  }
})

const formDataValidate = useVuelidate(formDataRules, formDataValue)
const timeDataValidate = useVuelidate(timeDataRules, timeDataValue)

const onChangeStartAt = (): void => {
  const startAt = dayjs(`${timeDataValue.value.startDate} ${timeDataValue.value.startTime}`)
  formDataValue.value.startAt = startAt.unix()
}

const onChangeEndAt = (): void => {
  const endAt = dayjs(`${timeDataValue.value.endDate} ${timeDataValue.value.endTime}`)
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
  const timeDataValid = await timeDataValidate.value.$validate()
  if (!formDataValid || !timeDataValid) {
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
            v-model="timeDataValidate.startDate.$model"
            :error-messages="getErrorMessage(timeDataValidate.startDate.$errors)"
            type="date"
            variant="outlined"
            density="compact"
            class="mr-md-2"
            @update:model-value="onChangeStartAt"
          />
          <v-text-field
            v-model="timeDataValidate.startTime.$model"
            :error-messages="getErrorMessage(timeDataValidate.startTime.$errors)"
            type="time"
            variant="outlined"
            density="compact"
            @update:model-value="onChangeStartAt"
          />
          <p class="text-subtitle-2 mx-4 pt-md-3 pb-4 pb-md-6">
            〜
          </p>
          <v-text-field
            v-model="timeDataValidate.endDate.$model"
            :error-messages="getErrorMessage(timeDataValidate.endDate.$errors)"
            type="date"
            variant="outlined"
            density="compact"
            class="mr-md-2"
            @update:model-value="onChangeEndAt"
          />
          <v-text-field
            v-model="timeDataValidate.endTime.$model"
            :error-messages="getErrorMessage(timeDataValidate.endTime.$errors)"
            type="time"
            variant="outlined"
            density="compact"
            @update:model-value="onChangeEndAt"
          />
        </div>
        <p class="text-subtitle-2 text-grey py-2">
          詳細
        </p>
        <client-only>
          <tiptap-editor
            v-model="formDataValidate.description.$model"
            :error-message="getErrorMessage(formDataValidate.description.$errors)"
            class="mb-4"
          />
        </client-only>
        <v-row>
          <v-col cols="12" sm="12" md="4">
            <molecules-image-select-form
              label="サムネイル画像"
              :loading="loading"
              :img-url="formDataValue.thumbnailUrl"
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
              :error="props.openingVideoUploadStatus.error"
              :message="props.openingVideoUploadStatus.message"
              @update:file="onChangeOpeningVideo"
            />
          </v-col>
          <v-col cols="12" sm="12" md="4">
            <molecules-image-select-form
              label="待機中の画像"
              :loading="loading"
              :accept="['image/png']"
              :img-url="formDataValue.imageUrl"
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
