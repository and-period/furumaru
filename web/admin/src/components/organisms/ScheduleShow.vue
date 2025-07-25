<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import dayjs, { unix } from 'dayjs'
import { getErrorMessage } from '~/lib/validations'
import {

  ScheduleStatus,

} from '~/types/api'
import type { Schedule, UpdateScheduleRequest } from '~/types/api'
import type { DateTimeInput, ImageUploadStatus } from '~/types/props'
import {
  TimeDataValidationRules,
  UpdateScheduleValidationRules,
} from '~/types/validations'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  updatable: {
    type: Boolean,
    default: false,
  },
  formData: {
    type: Object as PropType<UpdateScheduleRequest>,
    default: (): UpdateScheduleRequest => ({
      title: '',
      description: '',
      thumbnailUrl: '',
      imageUrl: '',
      openingVideoUrl: '',
      startAt: dayjs().unix(),
      endAt: dayjs().unix(),
    }),
  },
  schedule: {
    type: Object as PropType<Schedule>,
    default: (): Schedule => ({
      id: '',
      coordinatorId: '',
      title: '',
      description: '',
      status: ScheduleStatus.UNKNOWN,
      thumbnailUrl: '',
      imageUrl: '',
      openingVideoUrl: '',
      public: false,
      approved: false,
      startAt: dayjs().unix(),
      endAt: dayjs().unix(),
      createdAt: 0,
      updatedAt: 0,
    }),
  },
  thumbnailUploadStatus: {
    type: Object,
    default: (): ImageUploadStatus => ({
      error: false,
      message: '',
    }),
  },
  imageUploadStatus: {
    type: Object,
    default: (): ImageUploadStatus => ({
      error: false,
      message: '',
    }),
  },
  openingVideoUploadStatus: {
    type: Object,
    default: (): ImageUploadStatus => ({
      error: false,
      message: '',
    }),
  },
})

const emit = defineEmits<{
  (e: 'update:form-data', formData: UpdateScheduleRequest): void
  (e: 'update:schedule', formData: UpdateScheduleRequest): void
  (e: 'update:thumbnail', files: FileList): void
  (e: 'update:image', files: FileList): void
  (e: 'update:opening-video', files: FileList): void
  (e: 'update:public', publish: boolean): void
  (e: 'submit'): void
}>()

const statuses = [
  { title: '非公開', value: ScheduleStatus.PRIVATE },
  { title: '申請中', value: ScheduleStatus.IN_PROGRESS },
  { title: '開催前', value: ScheduleStatus.WAITING },
  { title: '開催中', value: ScheduleStatus.LIVE },
  { title: '終了', value: ScheduleStatus.CLOSED },
  { title: '不明', value: ScheduleStatus.UNKNOWN },
]

const scheduleValue = computed({
  get: (): Schedule => props.schedule,
  set: (schedule: Schedule): void => emit('update:schedule', schedule),
})
const formDataValue = computed({
  get: (): UpdateScheduleRequest => props.formData,
  set: (formData: UpdateScheduleRequest): void =>
    emit('update:form-data', formData),
})
const startTimeDataValue = computed({
  get: (): DateTimeInput => ({
    date: unix(props.formData.startAt).format('YYYY-MM-DD'),
    time: unix(props.formData.startAt).format('HH:mm'),
  }),
  set: (timeData: DateTimeInput): void => {
    const startAt = dayjs(`${timeData.date} ${timeData.time}`)
    formDataValue.value.startAt = startAt.unix()
  },
})
const endTimeDataValue = computed({
  get: (): DateTimeInput => ({
    date: unix(props.formData.endAt).format('YYYY-MM-DD'),
    time: unix(props.formData.endAt).format('HH:mm'),
  }),
  set: (timeData: DateTimeInput): void => {
    const endAt = dayjs(`${timeData.date} ${timeData.time}`)
    formDataValue.value.endAt = endAt.unix()
  },
})

const formDataValidate = useVuelidate(
  UpdateScheduleValidationRules,
  formDataValue,
)
const startTimeDataValidate = useVuelidate(
  TimeDataValidationRules,
  startTimeDataValue,
)
const endTimeDataValidate = useVuelidate(
  TimeDataValidationRules,
  endTimeDataValue,
)

const onChangeStartAt = (): void => {
  const startAt = dayjs(
    `${startTimeDataValue.value.date} ${startTimeDataValue.value.time}`,
  )
  formDataValue.value.startAt = startAt.unix()
}

const onChangeEndAt = (): void => {
  const endAt = dayjs(
    `${endTimeDataValue.value.date} ${endTimeDataValue.value.time}`,
  )
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
  <v-row>
    <v-col
      sm="12"
      md="12"
      lg="8"
    >
      <div class="mb-4">
        <v-text-field
          v-model="formDataValidate.title.$model"
          variant="outlined"
          :readonly="!updatable"
          :error-messages="getErrorMessage(formDataValidate.title.$errors)"
          label="タイトル"
        />
        <v-textarea
          v-model="formDataValidate.description.$model"
          :readonly="!updatable"
          variant="outlined"
          :error-messages="
            getErrorMessage(formDataValidate.description.$errors)
          "
          label="詳細"
          maxlength="2000"
        />
        <v-row>
          <v-col
            cols="12"
            sm="12"
            md="4"
          >
            <molecules-image-select-form
              label="サムネイル画像"
              :loading="loading"
              :img-url="formDataValue.thumbnailUrl"
              :error="props.thumbnailUploadStatus.error"
              :message="props.thumbnailUploadStatus.message"
              @update:file="onChangeThumbnailFile"
            />
          </v-col>
          <v-col
            cols="12"
            sm="12"
            md="4"
          >
            <molecules-video-select-form
              label="オープニング動画"
              :loading="loading"
              :video-url="formDataValue.openingVideoUrl"
              :error="props.openingVideoUploadStatus.error"
              :message="props.openingVideoUploadStatus.message"
              @update:file="onChangeOpeningVideo"
            />
          </v-col>
          <v-col
            cols="12"
            sm="12"
            md="4"
          >
            <molecules-image-select-form
              label="待機中の画像"
              :loading="loading"
              :img-url="formDataValue.imageUrl"
              :error="props.imageUploadStatus.error"
              :message="props.imageUploadStatus.message"
              @update:file="onChangeImageFile"
            />
          </v-col>
        </v-row>
      </div>
    </v-col>

    <v-col
      sm="12"
      md="12"
      lg="4"
    >
      <v-select
        v-model="scheduleValue.status"
        label="開催ステータス"
        :items="statuses"
        item-title="title"
        item-value="value"
        variant="outlined"
        readonly
      />
      <p class="text-subtitle-2 text-grey py-2">
        開催開始日時
      </p>
      <div class="d-flex flex-column flex-md-row justify-center">
        <v-text-field
          v-model="startTimeDataValidate.date.$model"
          :readonly="!updatable"
          :error-messages="getErrorMessage(startTimeDataValidate.date.$errors)"
          type="date"
          variant="outlined"
          density="compact"
          class="mr-md-2"
          @update:model-value="onChangeStartAt"
        />
        <v-text-field
          v-model="startTimeDataValidate.time.$model"
          :readonly="!updatable"
          :error-messages="getErrorMessage(startTimeDataValidate.time.$errors)"
          type="time"
          variant="outlined"
          density="compact"
          @update:model-value="onChangeStartAt"
        />
      </div>
      <p class="text-subtitle-2 text-grey py-2">
        開催終了日時
      </p>
      <div class="d-flex flex-column flex-md-row justify-center">
        <v-text-field
          v-model="endTimeDataValidate.date.$model"
          :readonly="!updatable"
          :error-messages="getErrorMessage(endTimeDataValidate.date.$errors)"
          type="date"
          variant="outlined"
          density="compact"
          class="mr-md-2"
          @update:model-value="onChangeEndAt"
        />
        <v-text-field
          v-model="endTimeDataValidate.time.$model"
          :readonly="!updatable"
          :error-messages="getErrorMessage(endTimeDataValidate.time.$errors)"
          type="time"
          variant="outlined"
          density="compact"
          @update:model-value="onChangeEndAt"
        />
      </div>
    </v-col>
  </v-row>

  <v-btn
    v-if="updatable"
    block
    :loading="loading"
    variant="outlined"
    color="primary"
    @click="onSubmit"
  >
    更新
  </v-btn>
</template>
