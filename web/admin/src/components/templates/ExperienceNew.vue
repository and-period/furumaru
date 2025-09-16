<script lang="ts" setup>
import {
  mdiClose,
  mdiPlus,
  mdiCalendarCheck,
  mdiImageMultiple,
  mdiStar,
  mdiCurrencyJpy,
  mdiMapMarker,
  mdiClock,
  mdiVideo,
  mdiArrowLeft,
  mdiTagMultiple,
} from '@mdi/js'
import useVuelidate from '@vuelidate/core'

import dayjs, { unix } from 'dayjs'
import type { AlertType } from '~/lib/hooks'
import type {
  CreateExperienceRequest,
  ExperienceType,
  Producer,
  UpdateExperienceRequest,
} from '~/types/api/v1'
import { ExperienceStatus } from '~/types/api/v1'
import type { DateTimeInput } from '~/types/props'
import {
  CreateExperienceValidationRules,
  NotSameTimeDataValidationRules,
  TimeDataValidationRules,
} from '~/types/validations'
import { getErrorMessage } from '~/lib/validations'
import { experienceStatues, experiencePublicationStatuses, experienceSoldStatus } from '~/constants'

interface Props {
  loading: boolean
  isAlert: boolean
  alertType: AlertType
  alertText: string
  formData: CreateExperienceRequest | UpdateExperienceRequest
  producers: Producer[]
  experienceTypes: ExperienceType[]
  searchErrorMessage: string
  searchLoading: boolean
  videoUploading: boolean
  producerSearchKeyword?: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'update:files', files: FileList): void
  (e: 'update:video', files: FileList): void
  (
    e: 'update:form-data',
    formData: CreateExperienceRequest | UpdateExperienceRequest,
  ): void
  (e: 'click:search-address'): void
  (e: 'submit'): void
  (e: 'update:producer-search-keyword', v: string): void
}>()

const thumbnailIndex = computed<number>({
  get: (): number => props.formData.media.findIndex(item => item.isThumbnail),
  set: (index: number): void => {
    if (formDataValue.value.media.length <= index) {
      return
    }
    formDataValue.value.media = formDataValue.value.media.map((item, i) => {
      if (i === index) {
        return {
          ...item,
          isThumbnail: true,
        }
      }
      else {
        return {
          ...item,
          isThumbnail: false,
        }
      }
    })
  },
})

const formDataValue = computed({
  get: (): CreateExperienceRequest | UpdateExperienceRequest => props.formData,
  set: (formData: CreateExperienceRequest | UpdateExperienceRequest): void =>
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
const experienceStatusValue = computed<ExperienceStatus>(() => {
  if (!formDataValue.value._public) {
    return ExperienceStatus.ExperienceStatusPrivate
  }
  if (formDataValue.value.soldOut) {
    return ExperienceStatus.ExperienceStatusSoldOut
  }
  const now = dayjs().unix()
  if (formDataValue.value.startAt > now) {
    return ExperienceStatus.ExperienceStatusWaiting
  }
  if (formDataValue.value.endAt !== 0 && formDataValue.value.endAt < now) {
    return ExperienceStatus.ExperienceStatusFinished
  }
  return ExperienceStatus.ExperienceStatusAccepting
})

const producerSearchKeywordValue = computed({
  get: (): string => props.producerSearchKeyword || '',
  set: (v: string): void => emit('update:producer-search-keyword', v),
})

const formDataValidate = useVuelidate(
  CreateExperienceValidationRules,
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

const notSameTimeValidate = useVuelidate(
  () => NotSameTimeDataValidationRules(props.formData.startAt, '販売開始日時'),
  formDataValue,
)

const getStatus = (status: ExperienceStatus): string => {
  const value = experienceStatues.find(s => s.value === status)
  return value ? value.title : ''
}

const getStatusColor = (status: ExperienceStatus): string => {
  switch (status) {
    case ExperienceStatus.ExperienceStatusWaiting:
      return 'info'
    case ExperienceStatus.ExperienceStatusAccepting:
      return 'primary'
    case ExperienceStatus.ExperienceStatusSoldOut:
      return 'secondary'
    case ExperienceStatus.ExperienceStatusPrivate:
      return 'warning'
    case ExperienceStatus.ExperienceStatusFinished:
      return 'error'
    default:
      return ''
  }
}

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

const onClickThumbnail = (i: number): void => {
  thumbnailIndex.value = i
}

const onDeleteThumbnail = (i: number): void => {
  const targetItem = props.formData.media.find((_, index) => index === i)
  if (!targetItem) {
    return
  }

  const media = targetItem.isThumbnail
    ? props.formData.media
        .filter((_, index) => index !== i)
        .map((item, i) => {
          return i === 0 ? { ...item, isThumbnail: true } : item
        })
    : props.formData.media.filter((_, index) => index !== i)
  formDataValue.value.media = media
}

const onClickImageUpload = (files?: FileList): void => {
  if (!files) {
    return
  }

  emit('update:files', files)
}

const onChangeVideo = (files?: FileList): void => {
  if (!files) {
    return
  }
  emit('update:video', files)
}

const onClickSearchAddress = (): void => {
  emit('click:search-address')
}

const onSubmit = async (): Promise<void> => {
  const formDataValid = await formDataValidate.value.$validate()
  const startTimeDataValid = await startTimeDataValidate.value.$validate()
  const endTimeDataValid = await endTimeDataValidate.value.$validate()
  const notSameTimeValid = await notSameTimeValidate.value.$validate()
  if (
    !formDataValid
    || !startTimeDataValid
    || !endTimeDataValid
    || !notSameTimeValid
  ) {
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
        体験登録
      </h1>
      <p class="text-body-1 text-grey-darken-1">
        新しい体験の情報を登録します。各セクションを順番に入力してください。
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
              :icon="mdiCalendarCheck"
              size="24"
              class="mr-3 text-primary"
            />
            <span class="text-h6 font-weight-medium">基本情報</span>
          </v-card-title>
          <v-card-text class="pa-6">
            <v-autocomplete
              v-model="formDataValidate.producerId.$model"
              v-model:search="producerSearchKeywordValue"
              label="生産者名 *"
              :items="producers"
              item-title="username"
              item-value="id"
              variant="outlined"
              density="comfortable"
              :error-messages="getErrorMessage(formDataValidate.producerId.$errors)"
              class="mb-4"
            />
            <v-text-field
              v-model="formDataValidate.title.$model"
              label="体験名 *"
              variant="outlined"
              density="comfortable"
              :error-messages="getErrorMessage(formDataValidate.title.$errors)"
              class="mb-4"
            />
            <v-textarea
              v-model="formDataValidate.description.$model"
              label="体験説明 *"
              maxlength="2000"
              variant="outlined"
              density="comfortable"
              rows="4"
              counter
              :error-messages="getErrorMessage(formDataValidate.description.$errors)"
            />
          </v-card-text>
        </v-card>

        <!-- 体験画像管理セクション -->
        <v-card
          :loading="props.loading"
          class="form-section-card mb-6"
          elevation="2"
        >
          <v-card-title class="d-flex align-center section-header">
            <v-icon
              :icon="mdiImageMultiple"
              size="24"
              class="mr-3 text-primary"
            />
            <span class="text-h6 font-weight-medium">体験画像管理</span>
          </v-card-title>
          <v-card-text class="pa-6">
            <div class="mb-4">
              <atoms-file-upload-filed
                text="体験画像をアップロード"
                @update:files="onClickImageUpload"
              />
            </div>

            <v-radio-group
              v-if="formDataValue.media.length > 0"
              v-model="thumbnailIndex"
              :error-messages="getErrorMessage(formDataValidate.media.$errors)"
              class="image-gallery"
            >
              <div class="mb-3">
                <v-chip
                  color="primary"
                  variant="outlined"
                  size="small"
                >
                  サムネイルを選択してください
                </v-chip>
              </div>
              <v-row>
                <v-col
                  v-for="(img, i) in formDataValue.media"
                  :key="i"
                  cols="6"
                  sm="4"
                  md="3"
                >
                  <v-card
                    class="image-card"
                    :class="{ 'thumbnail-selected': img.isThumbnail }"
                    @click="onClickThumbnail(i)"
                  >
                    <v-img
                      :src="img.url"
                      aspect-ratio="1"
                      class="image-preview"
                    >
                      <div class="image-overlay">
                        <v-radio
                          :value="i"
                          color="primary"
                          class="thumbnail-radio"
                        />
                        <v-btn
                          :icon="mdiClose"
                          color="error"
                          variant="text"
                          size="small"
                          class="delete-btn"
                          @click.stop="onDeleteThumbnail(i)"
                        />
                      </div>
                    </v-img>
                    <v-card-text class="pa-2 text-center">
                      <v-chip
                        v-if="img.isThumbnail"
                        color="primary"
                        size="x-small"
                        variant="elevated"
                      >
                        サムネイル
                      </v-chip>
                      <span
                        v-else
                        class="text-caption text-grey"
                      >
                        画像 {{ i + 1 }}
                      </span>
                    </v-card-text>
                  </v-card>
                </v-col>
              </v-row>
            </v-radio-group>
          </v-card-text>
        </v-card>

        <!-- 紹介動画セクション -->
        <v-card
          class="form-section-card mb-6"
          elevation="2"
        >
          <v-card-title class="d-flex align-center section-header">
            <v-icon
              :icon="mdiVideo"
              size="24"
              class="mr-3 text-primary"
            />
            <span class="text-h6 font-weight-medium">紹介動画</span>
          </v-card-title>
          <v-card-text class="pa-6">
            <molecules-video-select-form
              label="紹介動画をアップロード"
              :loading="props.videoUploading"
              @update:file="onChangeVideo"
            />
            <template v-if="formDataValue.promotionVideoUrl">
              <v-responsive
                :aspect-ratio="16 / 9"
                class="mt-4 border rounded"
              >
                <video
                  class="w-100"
                  controls
                  :src="formDataValue.promotionVideoUrl"
                />
              </v-responsive>
            </template>
          </v-card-text>
        </v-card>

        <!-- おすすめポイント・詳細セクション -->
        <v-card
          class="form-section-card mb-6"
          elevation="2"
        >
          <v-card-title class="d-flex align-center section-header">
            <v-icon
              :icon="mdiStar"
              size="24"
              class="mr-3 text-primary"
            />
            <span class="text-h6 font-weight-medium">おすすめポイント・詳細</span>
          </v-card-title>
          <v-card-text class="pa-6">
            <div class="mb-4">
              <p class="text-subtitle-2 mb-3 text-grey-darken-1">
                おすすめポイント
              </p>
              <v-text-field
                v-model="formDataValidate.recommendedPoint1.$model"
                :error-messages="
                  getErrorMessage(formDataValidate.recommendedPoint1.$errors)
                "
                label="ポイント 1"
                variant="outlined"
                density="comfortable"
                class="mb-3"
              />
              <v-text-field
                v-model="formDataValidate.recommendedPoint2.$model"
                :error-messages="
                  getErrorMessage(formDataValidate.recommendedPoint2.$errors)
                "
                label="ポイント 2"
                variant="outlined"
                density="comfortable"
                class="mb-3"
              />
              <v-text-field
                v-model="formDataValidate.recommendedPoint3.$model"
                :error-messages="
                  getErrorMessage(formDataValidate.recommendedPoint3.$errors)
                "
                label="ポイント 3"
                variant="outlined"
                density="comfortable"
              />
            </div>

            <v-row>
              <v-col
                cols="12"
                sm="6"
              >
                <v-text-field
                  v-model.number="formDataValidate.duration.$model"
                  label="所要時間"
                  type="number"
                  :min="0"
                  :max="24"
                  variant="outlined"
                  density="comfortable"
                  suffix="時間"
                  :error-messages="getErrorMessage(formDataValidate.duration.$errors)"
                />
              </v-col>
            </v-row>

            <v-textarea
              v-model="formDataValidate.direction.$model"
              label="アクセス方法"
              maxlength="2000"
              variant="outlined"
              density="comfortable"
              rows="3"
              counter
              :error-messages="getErrorMessage(formDataValidate.direction.$errors)"
            />
          </v-card-text>
        </v-card>

        <!-- 営業時間・場所セクション -->
        <v-card
          class="form-section-card mb-6"
          elevation="2"
        >
          <v-card-title class="d-flex align-center section-header">
            <v-icon
              :icon="mdiMapMarker"
              size="24"
              class="mr-3 text-primary"
            />
            <span class="text-h6 font-weight-medium">営業時間・場所</span>
          </v-card-title>
          <v-card-text class="pa-6">
            <div class="mb-4">
              <p class="text-subtitle-2 mb-3 text-grey-darken-1">
                営業時間
              </p>
              <v-row>
                <v-col
                  cols="12"
                  sm="6"
                >
                  <v-text-field
                    v-model="formDataValidate.businessOpenTime.$model"
                    type="time"
                    label="開始時間"
                    variant="outlined"
                    density="comfortable"
                  />
                </v-col>
                <v-col
                  cols="12"
                  sm="6"
                >
                  <v-text-field
                    v-model="formDataValidate.businessCloseTime.$model"
                    type="time"
                    label="終了時間"
                    variant="outlined"
                    density="comfortable"
                  />
                </v-col>
              </v-row>
            </div>

            <molecules-address-form
              v-model:postal-code="formDataValue.hostPostalCode"
              v-model:prefecture="formDataValue.hostPrefectureCode"
              v-model:city="formDataValue.hostCity"
              v-model:address-line1="formDataValue.hostAddressLine1"
              v-model:address-line2="formDataValue.hostAddressLine2"
              :error-messages="props.searchErrorMessage"
              :loading="props.searchLoading"
              @click:search="onClickSearchAddress"
            />
          </v-card-text>
        </v-card>

        <!-- 価格設定セクション -->
        <v-card
          class="form-section-card mb-6"
          elevation="2"
        >
          <v-card-title class="d-flex align-center section-header">
            <v-icon
              :icon="mdiCurrencyJpy"
              size="24"
              class="mr-3 text-primary"
            />
            <span class="text-h6 font-weight-medium">価格設定</span>
          </v-card-title>
          <v-card-text class="pa-6">
            <v-row>
              <v-col
                cols="12"
                sm="6"
              >
                <v-text-field
                  v-model.number="formDataValidate.priceAdult.$model"
                  :error-messages="
                    getErrorMessage(formDataValidate.priceAdult.$errors)
                  "
                  label="大人(高校生以上）"
                  type="number"
                  min="0"
                  suffix="円"
                  variant="outlined"
                  density="comfortable"
                />
              </v-col>
              <v-col
                cols="12"
                sm="6"
              >
                <v-text-field
                  v-model.number="formDataValidate.priceSenior.$model"
                  :error-messages="
                    getErrorMessage(formDataValidate.priceSenior.$errors)
                  "
                  label="シニア (65歳〜）"
                  type="number"
                  min="0"
                  suffix="円"
                  variant="outlined"
                  density="comfortable"
                />
              </v-col>
              <v-col
                cols="12"
                sm="6"
              >
                <v-text-field
                  v-model.number="formDataValidate.priceJuniorHighSchool.$model"
                  :error-messages="
                    getErrorMessage(formDataValidate.priceJuniorHighSchool.$errors)
                  "
                  label="中学生"
                  type="number"
                  min="0"
                  suffix="円"
                  variant="outlined"
                  density="comfortable"
                />
              </v-col>
              <v-col
                cols="12"
                sm="6"
              >
                <v-text-field
                  v-model.number="formDataValidate.priceElementarySchool.$model"
                  :error-messages="
                    getErrorMessage(formDataValidate.priceElementarySchool.$errors)
                  "
                  label="小学生"
                  type="number"
                  min="0"
                  suffix="円"
                  variant="outlined"
                  density="comfortable"
                />
              </v-col>
              <v-col
                cols="12"
                sm="6"
              >
                <v-text-field
                  v-model.number="formDataValidate.pricePreschool.$model"
                  :error-messages="
                    getErrorMessage(formDataValidate.pricePreschool.$errors)
                  "
                  label="未就学児 (3歳〜）"
                  type="number"
                  min="0"
                  suffix="円"
                  variant="outlined"
                  density="comfortable"
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
        <!-- 販売設定セクション -->
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
            <span class="text-h6 font-weight-medium">販売設定</span>
          </v-card-title>
          <v-card-text class="pa-6">
            <v-alert
              :color="getStatusColor(experienceStatusValue)"
              variant="tonal"
              density="compact"
              class="mb-4"
            >
              現在の状況: {{ getStatus(experienceStatusValue) }}
            </v-alert>

            <v-select
              v-model="formDataValue._public"
              label="公開状況 *"
              :items="experiencePublicationStatuses"
              item-title="title"
              item-value="value"
              variant="outlined"
              density="comfortable"
              class="mb-4"
            />
            <v-select
              v-model="formDataValue.soldOut"
              label="販売状況 *"
              :items="experienceSoldStatus"
              item-title="title"
              item-value="value"
              variant="outlined"
              density="comfortable"
              class="mb-4"
            />

            <div class="date-time-section">
              <p class="text-subtitle-2 mb-3 text-grey-darken-1">
                販売開始日時 *
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
                販売終了日時 *
              </p>
              <v-row>
                <v-col
                  cols="12"
                  sm="6"
                >
                  <v-text-field
                    v-model="endTimeDataValidate.date.$model"
                    :error-messages="
                      getErrorMessage(endTimeDataValidate.date.$errors)
                    "
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
                    :error-messages="
                      getErrorMessage(notSameTimeValidate.endAt.$errors)
                    "
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

        <!-- 詳細分類セクション -->
        <v-card
          class="form-section-card mb-6"
          elevation="2"
        >
          <v-card-title class="d-flex align-center section-header">
            <v-icon
              :icon="mdiTagMultiple"
              size="24"
              class="mr-3 text-primary"
            />
            <span class="text-h6 font-weight-medium">詳細分類</span>
          </v-card-title>
          <v-card-text class="pa-6">
            <v-select
              v-model="formDataValidate.experienceTypeId.$model"
              :error-messages="
                getErrorMessage(formDataValidate.experienceTypeId.$errors)
              "
              label="カテゴリ *"
              :items="experienceTypes"
              item-title="name"
              item-value="id"
              variant="outlined"
              density="comfortable"
              no-data-text="カテゴリを選択してください。"
              clearable
            />
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
        :loading="loading"
        color="primary"
        variant="elevated"
        size="large"
        @click="onSubmit"
      >
        <v-icon
          :icon="mdiPlus"
          start
        />
        体験を登録
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

.image-gallery {
  margin-top: 16px;
}

.image-card {
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid transparent;
}

.image-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgb(0 0 0 / 10%);
}

.thumbnail-selected {
  border-color: rgb(33 150 243);
  background: rgb(33 150 243 / 5%);
}

.image-preview {
  position: relative;
}

.image-overlay {
  position: absolute;
  top: 4px;
  right: 4px;
  display: flex;
  gap: 4px;
}

.thumbnail-radio {
  background: rgb(255 255 255 / 90%);
  border-radius: 50%;
}

.delete-btn {
  background: rgb(255 255 255 / 90%) !important;
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

  .image-card {
    margin-bottom: 16px;
  }
}
</style>
