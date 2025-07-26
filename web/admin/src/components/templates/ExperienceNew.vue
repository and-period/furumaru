<script lang="ts" setup>
import { mdiClose, mdiPlus } from '@mdi/js'
import useVuelidate from '@vuelidate/core'

import dayjs, { unix } from 'dayjs'
import type { AlertType } from '~/lib/hooks'
import type {
  CreateExperienceRequest,
  ExperienceType,
  Producer,
  UpdateExperienceRequest,
} from '~/types/api'
import type { DateTimeInput } from '~/types/props'
import {
  CreateExperienceValidationRules,
  NotSameTimeDataValidationRules,
  TimeDataValidationRules,
} from '~/types/validations'
import { getErrorMessage } from '~/lib/validations'

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

const publicStatus = [
  { title: '公開', value: true },
  { title: '非公開', value: false },
]

const soldStatus = [
  { title: '販売中', value: false },
  { title: '在庫なし', value: true },
]

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
  // 動画ファイルのemits
  emit('update:video', files)
}

const onClickSearchAddress = (): void => {
  emit('click:search-address')
}

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
  <v-alert
    v-show="props.isAlert"
    :type="props.alertType"
    class="mb-2"
    v-text="props.alertText"
  />

  <v-card class="mb-16">
    <v-card-title>体験登録</v-card-title>
    <v-card-text>
      <v-row>
    <v-col
      sm="12"
      md="12"
      lg="8"
    >
      <div class="mb-4">
        <v-card
          elevation="0"
          class="mb-4"
          :loading="loading"
        >
          <v-card-title>基本情報</v-card-title>
          <v-card-text>
            <v-autocomplete
              v-model="formDataValidate.producerId.$model"
              label="生産者名"
              :items="producers"
              item-title="username"
              item-value="id"
            />
            <v-text-field
              v-model="formDataValidate.title.$model"
              label="体験名"
              outlined
            />
            <v-textarea
              v-model="formDataValidate.description.$model"
              label="体験説明"
              maxlength="2000"
            />
            <v-number-input
              v-model="formDataValidate.duration.$model"
              :max="24"
              :min="0"
              :reverse="false"
              control-variant="default"
              label="所用時間"
              :hide-input="false"
              :inset="true"
            />
            <p class="text-subtitle-2 text-grey py-2">
              営業時間
            </p>
            <div class="d-flex flex-column flex-md-row justify-center">
              <v-text-field
                v-model="formDataValidate.businessOpenTime.$model"
                type="time"
                variant="outlined"
                density="compact"
              />
              <div class="pa-3">
                〜
              </div>
              <v-text-field
                v-model="formDataValidate.businessCloseTime.$model"
                type="time"
                variant="outlined"
                density="compact"
              />
            </div>
            <v-textarea
              v-model="formDataValidate.direction.$model"
              label="アクセス方法"
              maxlength="2000"
            />
          </v-card-text>
          <v-card-subtitle>商品画像登録</v-card-subtitle>
          <v-card-text>
            <v-radio-group
              v-model="thumbnailIndex"
              :error-messages="getErrorMessage(formDataValidate.media.$errors)"
            >
              <v-row>
                <v-col
                  v-for="(img, i) in formDataValue.media"
                  :key="i"
                  cols="4"
                  class="d-flex flex-row align-center"
                >
                  <v-card
                    rounded
                    variant="outlined"
                    width="100%"
                    :class="{ 'thumbnail-border': img.isThumbnail }"
                    @click="onClickThumbnail(i)"
                  >
                    <v-img
                      :src="img.url"
                      aspect-ratio="1"
                    >
                      <div class="d-flex col">
                        <v-radio
                          :value="i"
                          color="primary"
                        />
                        <v-btn
                          :icon="mdiClose"
                          color="error"
                          variant="text"
                          size="small"
                          @click="onDeleteThumbnail(i)"
                        />
                      </div>
                    </v-img>
                  </v-card>
                </v-col>
              </v-row>
            </v-radio-group>
            <p
              v-show="formDataValue.media.length > 0"
              class="mt-2"
            >
              ※ check された商品画像がサムネイルになります
            </p>
            <div class="mb-2">
              <atoms-file-upload-filed
                text="商品画像"
                @update:files="onClickImageUpload"
              />
            </div>
          </v-card-text>
          <div class="mx-4">
            <molecules-video-select-form
              label="紹介動画"
              :loading="videoUploading"
              @update:file="onChangeVideo"
            />
            <template v-if="formDataValue.promotionVideoUrl">
              <v-responsive
                :aspect-ratio="16 / 9"
                class="border pa-4"
              >
                <video
                  class="w-100"
                  controls
                  :src="formDataValue.promotionVideoUrl"
                />
              </v-responsive>
            </template>
          </div>
          <v-card-text>
            <v-text-field
              v-model="formDataValidate.recommendedPoint1.$model"
              :error-messages="
                getErrorMessage(formDataValidate.recommendedPoint1.$errors)
              "
              label="おすすめポイント1"
            />
            <v-text-field
              v-model="formDataValidate.recommendedPoint2.$model"
              :error-messages="
                getErrorMessage(formDataValidate.recommendedPoint2.$errors)
              "
              label="おすすめポイント2"
            />
            <v-text-field
              v-model="formDataValidate.recommendedPoint3.$model"
              :error-messages="
                getErrorMessage(formDataValidate.recommendedPoint3.$errors)
              "
              label="おすすめポイント3"
            />
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
          <v-card-title>価格設定</v-card-title>
          <v-card-text>
            <v-text-field
              v-model.number="formDataValidate.priceAdult.$model"
              :error-messages="
                getErrorMessage(formDataValidate.priceAdult.$errors)
              "
              label="大人(高校生以上）(〜64歳)"
              type="number"
              min="0"
              suffix="円"
            />
            <v-text-field
              v-model.number="formDataValidate.priceJuniorHighSchool.$model"
              :error-messages="
                getErrorMessage(formDataValidate.priceJuniorHighSchool.$errors)
              "
              label="中学生"
              type="number"
              min="0"
              suffix="円"
            />
            <v-text-field
              v-model.number="formDataValidate.priceElementarySchool.$model"
              :error-messages="
                getErrorMessage(formDataValidate.priceElementarySchool.$errors)
              "
              label="小学生"
              type="number"
              min="0"
              suffix="円"
            />
            <v-text-field
              v-model.number="formDataValidate.pricePreschool.$model"
              :error-messages="
                getErrorMessage(formDataValidate.pricePreschool.$errors)
              "
              label="未就学児 (3歳〜）"
              type="number"
              min="0"
              suffix="円"
            />
            <v-text-field
              v-model.number="formDataValidate.priceSenior.$model"
              :error-messages="
                getErrorMessage(formDataValidate.priceSenior.$errors)
              "
              label="シニア (65歳〜）"
              type="number"
              min="0"
              suffix="円"
            />
          </v-card-text>
        </v-card>
      </div>
    </v-col>
    <v-col
      sm="12"
      md="12"
      lg="4"
    >
      <v-card
        elevation="0"
        class="mb-4"
      >
        <v-card-title>販売設定</v-card-title>
        <v-card-text>
          <v-select
            v-model="formDataValue.public"
            label="販売状況"
            :items="publicStatus"
            item-title="title"
            item-value="value"
          />
          <v-select
            v-model="formDataValue.soldOut"
            label="公開状況"
            :items="soldStatus"
            item-title="title"
            item-value="value"
          />
          <p class="text-subtitle-2 text-grey py-2">
            販売開始日時
          </p>
          <div class="d-flex flex-column flex-md-row justify-center">
            <v-text-field
              v-model="startTimeDataValidate.date.$model"
              :error-messages="
                getErrorMessage(startTimeDataValidate.date.$errors)
              "
              type="date"
              variant="outlined"
              density="compact"
              class="mr-md-2"
              @update:model-value="onChangeStartAt"
            />
            <v-text-field
              v-model="startTimeDataValidate.time.$model"
              :error-messages="
                getErrorMessage(startTimeDataValidate.time.$errors)
              "
              type="time"
              variant="outlined"
              density="compact"
              @update:model-value="onChangeStartAt"
            />
          </div>
          <p class="text-subtitle-2 text-grey py-2">
            販売終了日時
          </p>
          <div class="d-flex flex-column flex-md-row justify-center">
            <v-text-field
              v-model="endTimeDataValidate.date.$model"
              :error-messages="
                getErrorMessage(endTimeDataValidate.date.$errors)
              "
              type="date"
              variant="outlined"
              density="compact"
              class="mr-md-2"
              @update:model-value="onChangeEndAt"
            />
            <v-text-field
              v-model="endTimeDataValidate.time.$model"
              :error-messages="
                getErrorMessage(notSameTimeValidate.endAt.$errors)
              "
              type="time"
              variant="outlined"
              density="compact"
              @update:model-value="onChangeEndAt"
            />
          </div>
        </v-card-text>
        <v-card-title>詳細情報</v-card-title>
        <v-card-text>
          <v-select
            v-model="formDataValidate.experienceTypeId.$model"
            :error-messages="
              getErrorMessage(formDataValidate.experienceTypeId.$errors)
            "
            label="カテゴリ"
            :items="experienceTypes"
            item-title="name"
            item-value="id"
            no-data-text="カテゴリを選択してください。"
            clearable
          />
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
    </v-card-text>
  </v-card>
</template>

<style lang="scss">
.thumbnail-border {
  border: 2px;
  border-style: solid;
  border-color: rgb(var(--v-theme-secondary));
}
</style>
