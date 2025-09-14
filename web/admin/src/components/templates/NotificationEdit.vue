<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import dayjs from 'dayjs'
import type { AlertType } from '~/lib/hooks'

import { getErrorMessage } from '~/lib/validations'
import { AdminType, NotificationStatus, NotificationType, PromotionStatus, PromotionTargetType } from '~/types/api/v1'
import type { Notification, Promotion, UpdateNotificationRequest } from '~/types/api/v1'
import { TimeDataValidationRules } from '~/types/validations'
import { UpdateNotificationValidationRules } from '~/types/validations/notification'
import { NOTIFICATION_TYPES, NOTIFICATION_TARGETS } from '~/constants/notification'
import { useNotificationForm } from '~/composables/useNotificationForm'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
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
    type: Object as PropType<UpdateNotificationRequest>,
    default: (): UpdateNotificationRequest => ({
      targets: [],
      title: '',
      body: '',
      note: '',
      publishedAt: dayjs().unix(),
    }),
  },
  notification: {
    type: Object as PropType<Notification>,
    default: (): Notification => ({
      id: '',
      type: NotificationType.NotificationTypeUnknown,
      status: NotificationStatus.NotificationStatusUnknown,
      targets: [],
      title: '',
      body: '',
      note: '',
      publishedAt: dayjs().unix(),
      promotionId: '',
      createdBy: '',
      createdAt: 0,
      updatedBy: '',
      updatedAt: 0,
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
      targetType: PromotionTargetType.PromotionTargetTypeUnknown,
      status: PromotionStatus.PromotionStatusUnknown,
      discountType: DiscountType.DiscountTypeUnknown,
      discountRate: 0,
      code: '',
      startAt: 0,
      endAt: 0,
      createdAt: 0,
      updatedAt: 0,
      usedAmount: 0,
      usedCount: 0,
    }),
  },
})

const emit = defineEmits<{
  (e: 'update:form-data', formData: UpdateNotificationRequest): void
  (e: 'submit'): void
}>()

const formDataValue = computed({
  get: (): UpdateNotificationRequest => props.formData as UpdateNotificationRequest,
  set: (formData: UpdateNotificationRequest) => emit('update:form-data', formData),
})

const { timeDataValue, onChangePublishedAt } = useNotificationForm(formDataValue)

const notificationValue = computed((): Notification => {
  return props.notification
})

const formDataValidate = useVuelidate(UpdateNotificationValidationRules, formDataValue)
const timeDataValidate = useVuelidate(TimeDataValidationRules, timeDataValue)

const isEditable = (): boolean => {
  return props.adminType === AdminType.AdminTypeAdministrator
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
    <v-card-title>お知らせ編集</v-card-title>

    <v-form @submit.prevent="onSubmit">
      <v-card-text class="pa-6">
        <div class="d-flex flex-column ga-4">
          <v-select
            v-model="notificationValue.type"
            :items="NOTIFICATION_TYPES"
            label="お知らせ種別"
            item-title="title"
            item-value="value"
            variant="outlined"
            readonly
          />

          <!-- セール情報 -->
          <div v-if="notification.type === NotificationType.NotificationTypePromotion">
            <NotificationPromotionDisplay
              :promotion="promotion"
              show-title
            />
          </div>

          <!-- その他 -->
          <v-text-field
            v-else
            v-model="formDataValidate.title.$model"
            :error-messages="getErrorMessage(formDataValidate.title.$errors)"
            label="タイトル *"
            variant="outlined"
            maxlength="128"
          />

          <!-- 共通部分 -->
          <v-autocomplete
            v-model="formDataValidate.targets.$model"
            :error-messages="getErrorMessage(formDataValidate.targets.$errors)"
            :items="NOTIFICATION_TARGETS"
            label="公開範囲 *"
            multiple
            item-title="title"
            item-value="value"
            variant="outlined"
          />

          <div class="d-flex flex-column ga-2">
            <v-label class="text-body-2 font-weight-medium">
              投稿日時 *
            </v-label>
            <div class="d-flex align-center ga-2">
              <v-text-field
                v-model="timeDataValidate.date.$model"
                :error-messages="getErrorMessage(timeDataValidate.date.$errors)"
                type="date"
                variant="outlined"
                density="compact"
                hide-details="auto"
                @update:model-value="onChangePublishedAt"
              />
              <v-text-field
                v-model="timeDataValidate.time.$model"
                :error-messages="getErrorMessage(timeDataValidate.time.$errors)"
                type="time"
                variant="outlined"
                density="compact"
                hide-details="auto"
                @update:model-value="onChangePublishedAt"
              />
            </div>
          </div>

          <v-textarea
            v-model="formDataValidate.body.$model"
            :error-messages="getErrorMessage(formDataValidate.body.$errors)"
            label="本文 *"
            placeholder="ユーザーに公開される内容を記載してください"
            variant="outlined"
            maxlength="2000"
            rows="4"
            counter
          />

          <v-textarea
            v-model="formDataValidate.note.$model"
            :error-messages="getErrorMessage(formDataValidate.note.$errors)"
            label="備考"
            placeholder="ユーザーには非公開にしたいコメント等を記載してください"
            variant="outlined"
            maxlength="2000"
            rows="3"
            counter
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
