<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import dayjs, { unix } from 'dayjs'

import { AlertType } from '~/lib/hooks'
import { getErrorMessage, maxLength, minLength, minValue, required } from '~/lib/validations'
import { DiscountType, PromotionResponse, UpdatePromotionRequest } from '~/types/api'
import { PromotionTime } from '~/types/props'

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
    type: Object as PropType<UpdatePromotionRequest>,
    default: (): UpdatePromotionRequest => ({
      title: '',
      description: '',
      public: false,
      discountType: DiscountType.AMOUNT,
      discountRate: 0,
      code: '',
      startAt: dayjs().unix(),
      endAt: dayjs().unix()
    })
  },
  promotion: {
    type: Object as PropType<PromotionResponse>,
    default: (): PromotionResponse => ({
      id: '',
      title: '',
      description: '',
      public: false,
      discountType: DiscountType.AMOUNT,
      discountRate: 0,
      code: '',
      startAt: dayjs().unix(),
      endAt: dayjs().unix(),
      createdAt: 0,
      updatedAt: 0
    })
  }
})

const emit = defineEmits<{
  (e: 'update:form-data', formData: UpdatePromotionRequest): void
  (e: 'update:promotion', promotion: PromotionResponse): void
  (e: 'submit'): void
}>()

const statusList = [
  { status: '有効', value: true },
  { status: '無効', value: false }
]
const discountMethodList = [
  { method: '円', value: DiscountType.AMOUNT },
  { method: '%', value: DiscountType.RATE },
  { method: '送料無料', value: DiscountType.FREE_SHIPPING }
]

const formDataRules = computed(() => ({
  title: { required, maxLength: maxLength(200) },
  description: { required, maxLength: maxLength(2000) },
  discountType: {},
  discountRate: { minValue: minValue(0) },
  code: { required, minLength: minLength(8), maxLength: maxLength(8) }
}))
const timeDataRules = computed(() => ({
  startDate: { required },
  startTime: { required },
  endDate: { required },
  endTime: { required }
}))
const formDataValue = computed({
  get: (): UpdatePromotionRequest => props.formData,
  set: (formData: UpdatePromotionRequest) => emit('update:form-data', formData)
})
const timeDataValue = computed({
  get: (): PromotionTime => ({
    startDate: unix(props.formData.startAt).format('YYYY-MM-DD'),
    startTime: unix(props.formData.startAt).format('HH:mm'),
    endDate: unix(props.formData.endAt).format('YYYY-MM-DD'),
    endTime: unix(props.formData.endAt).format('HH:mm')
  }),
  set: (timeData: PromotionTime): void => {
    const startAt = dayjs(`${timeData.startDate} ${timeData.startTime}`)
    const endAt = dayjs(`${timeData.endDate} ${timeData.endTime}`)
    formDataValue.value.startAt = startAt.unix()
    formDataValue.value.endAt = endAt.unix()
  }
})
const promotionValue = computed({
  get: (): PromotionResponse => props.promotion,
  set: (promotion: PromotionResponse): void => emit('update:promotion', promotion)
})

const formDataValidate = useVuelidate(formDataRules, formDataValue)
const timeDataValidate = useVuelidate(timeDataRules, timeDataValue)

const getDiscountErrorMessage = (): string => {
  switch (formDataValue.value.discountType) {
    case 1:
      if (formDataValue.value.discountRate >= 0) {
        return ''
      }
      return '0以上の値を指定してください'
    case 2:
      if (formDataValue.value.discountRate >= 0 && formDataValue.value.discountRate <= 100) {
        return ''
      }
      return '0~100の値を指定してください'
    default:
      return ''
  }
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
  <v-card>
    <v-card-title>セール情報詳細</v-card-title>

    <v-form @submit.prevent="onSubmit">
      <v-card-text>
        <v-text-field
          v-model="formDataValidate.title.$model"
          :error-messages="getErrorMessage(formDataValidate.title.$errors)"
          label="タイトル"
        />
        <v-textarea
          v-model="formDataValidate.description.$model"
          :error-messages="getErrorMessage(formDataValidate.description.$errors)"
          label="説明"
        />
        <div class="d-flex align-center">
          <v-text-field
            v-model="promotionValue.code"
            class="mr-4"
            label="割引コード(8文字)"
            readonly
          />
          <v-spacer />
        </div>
        <div class="d-flex align-center">
          <v-select
            v-model="formDataValidate.discountType.$model"
            :error-messages="getErrorMessage(formDataValidate.discountType.$errors)"
            :items="discountMethodList"
            item-title="method"
            item-value="value"
            label="割引方法"
          />
          <v-text-field
            v-if="formDataValue.discountType != 3"
            v-model="formDataValidate.discountRate.$model"
            :error-messages="getDiscountErrorMessage()"
            class="ml-4"
            type="number"
            label="割引値"
          />
        </div>
        <p class="text-h6">
          使用期間
        </p>
        <div class="d-flex align-center">
          <v-text-field
            v-model="timeDataValidate.startDate.$model"
            :error-messages="getErrorMessage(timeDataValidate.startDate.$errors)"
            type="date"
            variant="outlined"
            class="mr-2"
          />
          <v-text-field
            v-model="timeDataValidate.startTime.$model"
            :error-messages="getErrorMessage(timeDataValidate.startTime.$errors)"
            type="time"
            variant="outlined"
          />
          <p class="text-h6 mx-4 mb-6">
            〜
          </p>
          <v-text-field
            v-model="timeDataValidate.endDate.$model"
            :error-messages="getErrorMessage(timeDataValidate.endDate.$errors)"
            type="date"
            variant="outlined"
            class="mr-2"
          />
          <v-text-field
            v-model="timeDataValidate.endTime.$model"
            :error-messages="getErrorMessage(timeDataValidate.endTime.$errors)"
            type="time"
            variant="outlined"
          />
        </div>
        <v-switch v-model="formDataValue.public" label="クーポンを有効にする" color="primary" />
      </v-card-text>

      <v-card-actions>
        <v-btn block :loading="loading" variant="outlined" color="primary" type="submit">
          更新
        </v-btn>
      </v-card-actions>
    </v-form>
  </v-card>
</template>
