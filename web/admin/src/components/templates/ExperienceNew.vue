<script lang="ts" setup>
import { mdiClose, mdiPlus } from '@mdi/js'
import useVuelidate from '@vuelidate/core'

import dayjs, { unix } from 'dayjs'
import type { AlertType } from '~/lib/hooks'
import type {
  CreateExperienceRequest,
  Producer,
} from '~/types/api'
import type { DateTimeInput } from '~/types/props'
import {
  CreateProductValidationRules,
  NotSameTimeDataValidationRules,
  TimeDataValidationRules,
} from '~/types/validations'
import { getErrorMessage } from '~/lib/validations'

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
    type: Object as PropType<CreateExperienceRequest>,
    default: (): CreateExperienceRequest => ({
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
    }),
  },
  producers: {
    type: Array<Producer>,
    default: () => [],
  },
})

const formDataValue = computed({
  get: (): CreateExperienceRequest => props.formData,
  set: (formData: CreateExperienceRequest): void =>
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
  CreateProductValidationRules,
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
</script>

<template>
  <v-alert
    v-show="props.isAlert"
    :type="props.alertType"
    class="mb-2"
    v-text="props.alertText"
  />

  <v-card-title>体験登録</v-card-title>

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
        >
          <v-card-title>基本情報</v-card-title>
          <v-card-text>
            <v-autocomplete
              label="生産者名"
              :items="producers"
              item-title="username"
              item-value="id"
            />
            <v-text-field
              label="体験名"
              outlined
            />
            <v-textarea
              label="体験説明"
              maxlength="2000"
            />
            <v-number-input
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
                type="time"
                variant="outlined"
                density="compact"
              />
              <div class="pa-3">〜</div>
              <v-text-field
                type="time"
                variant="outlined"
                density="compact"
              />
            </div>
            <v-textarea
              label="アクセス方法"
              maxlength="2000"
            />
          </v-card-text>
          <v-card-subtitle>商品画像登録</v-card-subtitle>
          <v-card-text>
            <v-radio-group
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
              />
            </div>
          </v-card-text>
          <v-card-text>
            <v-text-field
              :error-messages="
                getErrorMessage(formDataValidate.recommendedPoint1.$errors)
              "
              label="おすすめポイント1"
            />
            <v-text-field
              :error-messages="
                getErrorMessage(formDataValidate.recommendedPoint2.$errors)
              "
              label="おすすめポイント2"
            />
            <v-text-field
              :error-messages="
                getErrorMessage(formDataValidate.recommendedPoint3.$errors)
              "
              label="おすすめポイント3"
            />
            <molecules-address-form />
          </v-card-text>
          <v-card-title>価格設定</v-card-title>
          <v-card-text>
            <v-text-field
              :error-messages="getErrorMessage(formDataValidate.price.$errors)"
              label="大人(高校生以上）(〜64歳)"
              type="number"
              min="0"
              suffix="円"
            />
            <v-text-field
              :error-messages="getErrorMessage(formDataValidate.price.$errors)"
              label="中学生"
              type="number"
              min="0"
              suffix="円"
            />
            <v-text-field
              :error-messages="getErrorMessage(formDataValidate.price.$errors)"
              label="小学生"
              type="number"
              min="0"
              suffix="円"
            />
            <v-text-field
              :error-messages="getErrorMessage(formDataValidate.price.$errors)"
              label="未就学児 (3歳〜）"
              type="number"
              min="0"
              suffix="円"
            />
            <v-text-field
              :error-messages="getErrorMessage(formDataValidate.price.$errors)"
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
            label="販売状況"
            item-title="title"
            item-value="value"
            variant="plain"
            readonly
          />
          <v-select
            label="公開状況"
          />
          <p class="text-subtitle-2 text-grey py-2">
            販売開始日時
          </p>
          <div class="d-flex flex-column flex-md-row justify-center">
            <v-text-field
              type="date"
              variant="outlined"
              density="compact"
              class="mr-md-2"
            />
            <v-text-field
              type="time"
              variant="outlined"
              density="compact"
            />
          </div>
          <p class="text-subtitle-2 text-grey py-2">
            販売終了日時
          </p>
          <div class="d-flex flex-column flex-md-row justify-center">
            <v-text-field
              type="date"
              variant="outlined"
              density="compact"
              class="mr-md-2"
            />
            <v-text-field
              type="time"
              variant="outlined"
              density="compact"
            />
          </div>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>

  <v-btn
    :loading="loading"
    block
    variant="outlined"
    @click="onSubmit"
  >
    <v-icon
      start
      :icon="mdiPlus"
    />
    登録
  </v-btn>
</template>

<style lang="scss">
.thumbnail-border {
  border: 2px;
  border-style: solid;
  border-color: rgb(var(--v-theme-secondary));
}
</style>
