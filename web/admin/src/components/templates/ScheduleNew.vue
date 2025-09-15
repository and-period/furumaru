<script lang="ts" setup>
import {
  mdiPlus,
  mdiArrowLeft,
  mdiCalendarToday,
  mdiImageMultiple,
  mdiVideo,
  mdiClock,
  mdiFileDocument,
} from '@mdi/js'
import useVuelidate from '@vuelidate/core'
import dayjs, { unix } from 'dayjs'
import type { AlertType } from '~/lib/hooks'
import { getErrorMessage } from '~/lib/validations'
import type { CreateScheduleRequest } from '~/types/api/v1'
import type { DateTimeInput, ImageUploadStatus } from '~/types/props'
import {
  CreateScheduleValidationRules,
  TimeDataValidationRules,
} from '~/types/validations'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  isAlert: {
    type: Boolean,
    default: false,
  },
  alertType: {
    type: String as PropType<AlertType>,
    default: undefined,
  },
  alertText: {
    type: String,
    default: '',
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
      _public: false,
      startAt: dayjs().unix(),
      endAt: dayjs().unix(),
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
  (e: 'update:form-data', formData: CreateScheduleRequest): void
  (e: 'update:thumbnail', files: FileList): void
  (e: 'update:image', files: FileList): void
  (e: 'update:opening-video', files: FileList): void
  (e: 'submit'): void
}>()

const formDataValue = computed({
  get: (): CreateScheduleRequest => props.formData,
  set: (formData: CreateScheduleRequest): void =>
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
  CreateScheduleValidationRules,
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
  <v-container class="pa-6">
    <v-alert
      v-show="props.isAlert"
      :type="props.alertType"
      class="mb-6"
      v-text="props.alertText"
    />

    <div class="mb-6">
      <v-btn
        variant="text"
        :prepend-icon="mdiArrowLeft"
        @click="$router.back()"
      >
        戻る
      </v-btn>
      <h1 class="text-h4 font-weight-bold mt-2 mb-2">
        ライブ配信登録
      </h1>
      <p class="text-body-1 text-grey-darken-1">
        新しいライブ配信の情報を登録します。各セクションを順番に入力してください。
      </p>
    </div>

    <v-row>
      <v-col
        cols="12"
        lg="8"
      >
        <!-- 基本情報セクション -->
        <v-card
          class="form-section-card mb-6"
          elevation="2"
        >
          <v-card-title class="d-flex align-center section-header">
            <v-icon
              :icon="mdiCalendarToday"
              size="24"
              class="mr-3 text-primary"
            />
            <span class="text-h6 font-weight-medium">基本情報</span>
          </v-card-title>
          <v-card-text class="pa-6">
            <v-text-field
              v-model="formDataValidate.title.$model"
              :error-messages="getErrorMessage(formDataValidate.title.$errors)"
              label="タイトル *"
              variant="outlined"
              density="comfortable"
              class="mb-4"
            />
            <v-textarea
              v-model="formDataValidate.description.$model"
              :error-messages="
                getErrorMessage(formDataValidate.description.$errors)
              "
              label="詳細説明 *"
              maxlength="2000"
              variant="outlined"
              density="comfortable"
              rows="4"
              counter
            />
          </v-card-text>
        </v-card>

        <!-- メディアファイル管理セクション -->
        <v-card
          class="form-section-card mb-6"
          elevation="2"
        >
          <v-card-title class="d-flex align-center section-header">
            <v-icon
              :icon="mdiImageMultiple"
              size="24"
              class="mr-3 text-primary"
            />
            <span class="text-h6 font-weight-medium">メディアファイル管理</span>
          </v-card-title>
          <v-card-text class="pa-6">
            <v-row>
              <v-col
                cols="12"
                md="4"
              >
                <molecules-image-select-form
                  label="サムネイル画像"
                  :loading="loading"
                  :img-url="formDataValue.thumbnailUrl"
                  :validation-error-message="
                    getErrorMessage(formDataValidate.thumbnailUrl.$errors)
                  "
                  :error="props.thumbnailUploadStatus.error"
                  :message="props.thumbnailUploadStatus.message"
                  @update:file="onChangeThumbnailFile"
                />
              </v-col>
              <v-col
                cols="12"
                md="4"
              >
                <molecules-video-select-form
                  label="オープニング動画"
                  :loading="loading"
                  :video-url="formDataValue.openingVideoUrl"
                  :validation-error-message="
                    getErrorMessage(formDataValidate.openingVideoUrl.$errors)
                  "
                  :error="props.openingVideoUploadStatus.error"
                  :message="props.openingVideoUploadStatus.message"
                  @update:file="onChangeOpeningVideo"
                />
              </v-col>
              <v-col
                cols="12"
                md="4"
              >
                <molecules-image-select-form
                  label="待機中の画像"
                  :loading="loading"
                  :img-url="formDataValue.imageUrl"
                  :validation-error-message="
                    getErrorMessage(formDataValidate.imageUrl.$errors)
                  "
                  :error="props.imageUploadStatus.error"
                  :message="props.imageUploadStatus.message"
                  @update:file="onChangeImageFile"
                />
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col
        cols="12"
        lg="4"
      >
        <!-- 開催期間設定セクション -->
        <v-card
          class="form-section-card mb-6"
          elevation="2"
        >
          <v-card-title class="d-flex align-center section-header">
            <v-icon
              :icon="mdiClock"
              size="24"
              class="mr-3 text-primary"
            />
            <span class="text-h6 font-weight-medium">開催期間設定</span>
          </v-card-title>
          <v-card-text class="pa-6">
            <div class="date-time-section">
              <p class="text-subtitle-2 mb-3 text-grey-darken-1">
                開始日時 *
              </p>
              <v-row>
                <v-col
                  cols="12"
                  sm="6"
                >
                  <v-text-field
                    v-model="startTimeDataValidate.date.$model"
                    :error-messages="
                      getErrorMessage(startTimeDataValidate.date.$errors)
                    "
                    label="日付"
                    type="date"
                    variant="outlined"
                    density="comfortable"
                    @update:model-value="onChangeStartAt"
                  />
                </v-col>
                <v-col
                  cols="12"
                  sm="6"
                >
                  <v-text-field
                    v-model="startTimeDataValidate.time.$model"
                    :error-messages="
                      getErrorMessage(startTimeDataValidate.time.$errors)
                    "
                    label="時刻"
                    type="time"
                    variant="outlined"
                    density="comfortable"
                    @update:model-value="onChangeStartAt"
                  />
                </v-col>
              </v-row>

              <p class="text-subtitle-2 mb-3 mt-4 text-grey-darken-1">
                終了日時 *
              </p>
              <v-row>
                <v-col
                  cols="12"
                  sm="6"
                >
                  <v-text-field
                    v-model="endTimeDataValidate.date.$model"
                    :error-messages="getErrorMessage(endTimeDataValidate.date.$errors)"
                    label="日付"
                    type="date"
                    variant="outlined"
                    density="comfortable"
                    @update:model-value="onChangeEndAt"
                  />
                </v-col>
                <v-col
                  cols="12"
                  sm="6"
                >
                  <v-text-field
                    v-model="endTimeDataValidate.time.$model"
                    :error-messages="getErrorMessage(endTimeDataValidate.time.$errors)"
                    label="時刻"
                    type="time"
                    variant="outlined"
                    density="comfortable"
                    @update:model-value="onChangeEndAt"
                  />
                </v-col>
              </v-row>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- 送信ボタン -->
    <div class="d-flex justify-end gap-3 mt-6">
      <v-btn
        variant="text"
        size="large"
        @click="$router.back()"
      >
        キャンセル
      </v-btn>
      <v-btn
        :loading="props.loading"
        color="primary"
        variant="elevated"
        size="large"
        @click="onSubmit"
      >
        <v-icon
          :icon="mdiPlus"
          start
        />
        ライブ配信を登録
      </v-btn>
    </div>
  </v-container>
</template>

<style scoped>
.form-section-card {
  border-radius: 12px;
  max-width: none;
}

.section-header {
  background: linear-gradient(90deg, rgb(33 150 243 / 5%) 0%, rgb(33 150 243 / 0%) 100%);
  border-bottom: 1px solid rgb(0 0 0 / 5%);
  padding: 20px 24px;
}

.date-time-section {
  border-top: 1px solid rgb(0 0 0 / 10%);
  padding-top: 16px;
}

@media (width <= 600px) {
  .form-section-card {
    border-radius: 8px;
  }

  .section-header {
    padding: 16px 20px;
  }
}
</style>
