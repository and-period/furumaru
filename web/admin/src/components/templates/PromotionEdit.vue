<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import dayjs, { unix } from 'dayjs'

import type { AlertType } from '~/lib/hooks'
import { getErrorMessage } from '~/lib/validations'
import { AdminType, DiscountType, PromotionStatus, PromotionTargetType } from '~/types/api/v1'
import type { Promotion, Shop, UpdatePromotionRequest } from '~/types/api/v1'
import type { DateTimeInput } from '~/types/props'
import { TimeDataValidationRules, UpdatePromotionValidationRules } from '~/types/validations'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  shopIds: {
    type: Array as PropType<string[]>,
    default: () => [],
  },
  adminType: {
    type: Number as PropType<AdminType>,
    default: AdminType.AdminTypeUnknown,
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
    type: Object as PropType<UpdatePromotionRequest>,
    default: (): UpdatePromotionRequest => ({
      title: '',
      description: '',
      _public: false,
      discountType: DiscountType.DiscountTypeAmount,
      discountRate: 0,
      code: '',
      startAt: dayjs().unix(),
      endAt: dayjs().unix(),
    }),
  },
  promotion: {
    type: Object as PropType<Promotion>,
    default: (): Promotion => ({
      id: '',
      shopId: '',
      title: '',
      description: '',
      _public: false,
      status: PromotionStatus.PromotionStatusUnknown,
      targetType: PromotionTargetType.PromotionTargetTypeUnknown,
      discountType: DiscountType.DiscountTypeAmount,
      discountRate: 0,
      code: '',
      usedCount: 0,
      usedAmount: 0,
      startAt: dayjs().unix(),
      endAt: dayjs().unix(),
      createdAt: 0,
      updatedAt: 0,
    }),
  },
  shop: {
    type: Object as PropType<Shop>,
    default: (): Shop => ({
      id: '',
      coordinatorId: '',
      producerIds: [],
      productTypeIds: [],
      name: '',
      businessDays: [],
      createdAt: 0,
      updatedAt: 0,
    }),
  },
})

const emit = defineEmits<{
  (e: 'update:form-data', formData: UpdatePromotionRequest): void
  (e: 'update:promotion', promotion: Promotion): void
  (e: 'submit'): void
}>()

const discountMethodList = [
  { method: '円', value: DiscountType.DiscountTypeAmount },
  { method: '%', value: DiscountType.DiscountTypeRate },
  { method: '送料無料', value: DiscountType.DiscountTypeFreeShipping },
]

const formDataValue = computed({
  get: (): UpdatePromotionRequest => props.formData,
  set: (formData: UpdatePromotionRequest) => emit('update:form-data', formData),
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
const promotionValue = computed({
  get: (): Promotion => props.promotion,
  set: (promotion: Promotion): void => emit('update:promotion', promotion),
})
const getTarget = computed(() => {
  if (!props.promotion) {
    return ''
  }
  switch (props.promotion.targetType) {
    case PromotionTargetType.PromotionTargetTypeAllShop:
      return '全て'
    case PromotionTargetType.PromotionTargetTypeSpecificShop:
      return props.shop.name
    default:
      return ''
  }
})

const formDataValidate = useVuelidate(UpdatePromotionValidationRules, formDataValue)
const startTimeDataValidate = useVuelidate(TimeDataValidationRules, startTimeDataValue)
const endTimeDataValidate = useVuelidate(TimeDataValidationRules, endTimeDataValue)

const isEditable = (): boolean => {
  switch (props.adminType) {
    case AdminType.AdminTypeAdministrator:
      return true
    case AdminType.AdminTypeCoordinator:
      return props.shopIds.includes(props.promotion.shopId)
    default:
      return false
  }
}

const onChangeStartAt = (): void => {
  const startAt = dayjs(`${startTimeDataValue.value.date} ${startTimeDataValue.value.time}`)
  formDataValue.value.startAt = startAt.unix()
}

const onChangeEndAt = (): void => {
  const endAt = dayjs(`${endTimeDataValue.value.date} ${endTimeDataValue.value.time}`)
  formDataValue.value.endAt = endAt.unix()
}

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
  const startTimeDataValid = await startTimeDataValidate.value.$validate()
  const endTimeDataValid = await endTimeDataValidate.value.$validate()
  if (!formDataValid || !startTimeDataValid || !endTimeDataValid) {
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
        <v-text-field
          v-model="getTarget"
          label="対象マルシェ"
          readonly
        />
        <v-textarea
          v-model="formDataValidate.description.$model"
          :error-messages="getErrorMessage(formDataValidate.description.$errors)"
          label="説明"
        />
        <div class="d-flex align-center">
          <v-text-field
            v-model="promotionValue.code"
            class="mr-2"
            label="割引コード(8文字)"
            readonly
          />
          <v-spacer />
        </div>
        <div class="d-flex align-center">
          <v-select
            v-model="formDataValidate.discountType.$model"
            :error-messages="getErrorMessage(formDataValidate.discountType.$errors)"
            label="割引方法"
            :items="discountMethodList"
            item-title="method"
            item-value="value"
            class="mr-2"
          />
          <v-text-field
            v-if="formDataValue.discountType != 3"
            v-model.number="formDataValidate.discountRate.$model"
            :error-messages="getDiscountErrorMessage()"
            label="割引値"
            type="number"
          />
        </div>
        <p class="text-subtitle-2 text-grey py-2">
          使用期間
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
          <p class="text-subtitle-2 mx-4 pt-md-3 mb-4 pb-md-6">
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
        <v-switch
          v-model="formDataValue._public"
          label="クーポンを有効にする"
          color="primary"
        />
      </v-card-text>

      <v-card-actions>
        <v-btn
          v-show="isEditable()"
          block
          :loading="loading"
          variant="outlined"
          color="primary"
          type="submit"
        >
          更新
        </v-btn>
      </v-card-actions>
    </v-form>
  </v-card>
</template>
