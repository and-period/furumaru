<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import dayjs from 'dayjs'

import type { AlertType } from '~/lib/hooks'
import { getErrorMessage } from '~/lib/validations'
import { DiscountType } from '~/types/api/v1'
import type { CreatePromotionRequest } from '~/types/api/v1'
import { CreatePromotionValidationRules } from '~/types/validations'
import { DISCOUNT_METHODS } from '~/constants/promotion'
import { usePromotionForm } from '~/composables/usePromotionForm'
import { usePromotionValidation } from '~/composables/usePromotionValidation'

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
    type: Object as PropType<CreatePromotionRequest>,
    default: (): CreatePromotionRequest => ({
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
})

const emit = defineEmits<{
  (e: 'update:form-data', formData: CreatePromotionRequest): void
  (e: 'submit'): void
}>()

const formDataValue = computed({
  get: (): CreatePromotionRequest => props.formData,
  set: (val: CreatePromotionRequest) => emit('update:form-data', val),
})

const { startTimeDataValue, endTimeDataValue, onChangeStartAt, onChangeEndAt, generateRandomCode } = usePromotionForm(formDataValue)
const { getDiscountErrorMessage } = usePromotionValidation(formDataValue)

const formDataValidate = useVuelidate(CreatePromotionValidationRules, formDataValue)

const onClickGenerateCode = (): void => {
  formDataValue.value.code = generateRandomCode()
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
  <atoms-app-alert
    :show="props.isAlert"
    :type="props.alertType"
    :text="props.alertText"
  />

  <v-card>
    <v-card-title>セール情報登録</v-card-title>

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

          <v-textarea
            v-model="formDataValidate.description.$model"
            :error-messages="getErrorMessage(formDataValidate.description.$errors)"
            label="説明"
            variant="outlined"
            rows="3"
            maxlength="500"
            counter
          />

          <div class="d-flex align-center ga-2">
            <v-text-field
              v-model="formDataValidate.code.$model"
              :error-messages="getErrorMessage(formDataValidate.code.$errors)"
              label="割引コード (8文字) *"
              variant="outlined"
              maxlength="8"
            />
            <v-btn
              variant="tonal"
              color="primary"
              @click="onClickGenerateCode"
            >
              自動生成
            </v-btn>
          </div>

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
    </v-form>
  </v-card>

  <v-footer
    app
    color="white"
    elevation="8"
    class="px-6 py-4 fixed-footer-actions"
  >
    <v-container
      fluid
      class="pa-0"
    >
      <div class="d-flex align-center justify-center flex-wrap ga-3">
        <v-btn
          :loading="loading"
          color="primary"
          variant="elevated"
          size="large"
          @click="onSubmit"
        >
          登録
        </v-btn>
      </div>
    </v-container>
  </v-footer>
</template>
