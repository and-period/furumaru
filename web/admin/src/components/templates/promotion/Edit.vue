<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import dayjs from 'dayjs'

import type { AlertType } from '~/lib/hooks'
import { getErrorMessage } from '~/lib/validations'
import { AdminType, DiscountType, PromotionStatus, PromotionTargetType } from '~/types/api/v1'
import type { Promotion, Shop, UpdatePromotionRequest } from '~/types/api/v1'
import { UpdatePromotionValidationRules } from '~/types/validations'
import { DISCOUNT_METHODS } from '~/constants/promotion'
import { usePromotionForm } from '~/composables/usePromotionForm'
import { usePromotionValidation } from '~/composables/usePromotionValidation'

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

const formDataValue = computed({
  get: (): UpdatePromotionRequest => props.formData,
  set: (formData: UpdatePromotionRequest) => emit('update:form-data', formData),
})

const { startTimeDataValue, endTimeDataValue, onChangeStartAt, onChangeEndAt } = usePromotionForm(formDataValue)
const { getDiscountErrorMessage } = usePromotionValidation(formDataValue)

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

const onSubmit = async (): Promise<void> => {
  const formDataValid = await formDataValidate.value.$validate()
  if (!formDataValid) {
    return
  }

  emit('submit')
}
</script>

<template>
  <v-card>
    <v-card-title>セール情報詳細</v-card-title>

    <v-form @submit.prevent="onSubmit">
      <v-card-text class="pa-6">
        <div class="d-flex flex-column ga-4">
          <v-text-field
            v-model="formDataValidate.title.$model"
            :error-messages="getErrorMessage(formDataValidate.title.$errors)"
            label="タイトル *"
            variant="outlined"
            maxlength="128"
          />

          <v-text-field
            v-model="getTarget"
            label="対象マルシェ"
            variant="outlined"
            readonly
          />

          <v-textarea
            v-model="formDataValidate.description.$model"
            :error-messages="getErrorMessage(formDataValidate.description.$errors)"
            label="説明"
            variant="outlined"
            rows="3"
            maxlength="500"
            counter
          />

          <v-text-field
            v-model="promotionValue.code"
            label="割引コード (8文字)"
            variant="outlined"
            readonly
          />

          <div class="d-flex align-center ga-2">
            <v-select
              v-model="formDataValidate.discountType.$model"
              :error-messages="getErrorMessage(formDataValidate.discountType.$errors)"
              label="割引方法 *"
              :items="DISCOUNT_METHODS"
              item-title="method"
              item-value="value"
              variant="outlined"
            />
            <v-text-field
              v-if="formDataValue.discountType !== DiscountType.DiscountTypeFreeShipping"
              v-model.number="formDataValidate.discountRate.$model"
              :error-messages="getDiscountErrorMessage()"
              label="割引値 *"
              type="number"
              variant="outlined"
              min="0"
              :max="formDataValue.discountType === DiscountType.DiscountTypeRate ? 100 : undefined"
            />
          </div>

          <molecules-promotion-period-input
            v-model:start-time="startTimeDataValue"
            v-model:end-time="endTimeDataValue"
            @change:start-at="onChangeStartAt"
            @change:end-at="onChangeEndAt"
          />

          <v-switch
            v-model="formDataValue._public"
            label="クーポンを有効にする"
            color="primary"
            inset
          />
        </div>
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
