<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import dayjs from 'dayjs'

import type { AlertType } from '~/lib/hooks'
import { getErrorMessage } from '~/lib/validations'
import { NotificationType } from '~/types/api/v1'
import type { CreateNotificationRequest, Promotion } from '~/types/api/v1'
import { TimeDataValidationRules } from '~/types/validations'
import { CreateNotificationValidationRules } from '~/types/validations/notification'
import { NOTIFICATION_TYPES, NOTIFICATION_TARGETS } from '~/constants/notification'
import { useNotificationForm } from '~/composables/useNotificationForm'

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
    type: Object,
    default: (): CreateNotificationRequest => ({
      type: NotificationType.NotificationTypeUnknown,
      targets: [],
      title: '',
      body: '',
      note: '',
      publishedAt: dayjs().unix(),
      promotionId: '',
    }),
  },
  promotions: {
    type: Array<Promotion>,
    default: (): Promotion[] => [],
  },
})

const emit = defineEmits<{
  (e: 'update:form-data', formData: CreateNotificationRequest): void
  (e: 'update:notification-type', type: NotificationType): void
  (e: 'update:search-promotion', name: string): void
  (e: 'submit'): void
}>()

const selectedPromotion = ref<Promotion>()

const formDataValue = computed({
  get: (): CreateNotificationRequest => props.formData as CreateNotificationRequest,
  set: (formData: CreateNotificationRequest) => emit('update:form-data', formData),
})

const { timeDataValue, onChangePublishedAt } = useNotificationForm(formDataValue)

const formDataValidate = useVuelidate(CreateNotificationValidationRules, formDataValue)
const timeDataValidate = useVuelidate(TimeDataValidationRules, timeDataValue)

const onChangeType = (type: NotificationType): void => {
  emit('update:notification-type', type)
}

const onChangeSearchPromotion = (name: string): void => {
  emit('update:search-promotion', name)
}

const onChangePromotion = (promotionId: string): void => {
  const promotion = props.promotions.find((promotion): boolean => promotion.id === promotionId)
  selectedPromotion.value = promotion
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
  <v-alert
    v-show="props.isAlert"
    :type="props.alertType"
    v-text="props.alertText"
  />

  <v-card>
    <v-card-title>お知らせ登録</v-card-title>

    <v-form @submit.prevent="onSubmit">
      <v-card-text class="pa-6">
        <div class="d-flex flex-column ga-4">
          <v-select
            v-model="formDataValidate.type.$model"
            :error-messages="getErrorMessage(formDataValidate.type.$errors)"
            :items="NOTIFICATION_TYPES"
            :loading="loading"
            :disabled="loading"
            label="お知らせ種別 *"
            item-title="title"
            item-value="value"
            variant="outlined"
            @update:model-value="onChangeType"
          />

          <!-- セール情報 -->
          <div
            v-if="formDataValue.type === NotificationType.NotificationTypePromotion"
            class="d-flex flex-column ga-3"
          >
            <v-autocomplete
              v-model="formDataValidate.promotionId.$model"
              :error-messages="getErrorMessage(formDataValidate.promotionId.$errors)"
              :items="promotions"
              label="セール情報 *"
              item-title="title"
              item-value="id"
              variant="outlined"
              @update:model-value="onChangePromotion"
              @update:search="onChangeSearchPromotion"
            />
            <molecules-notification-promotion-display :promotion="selectedPromotion" />
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
          block
          :loading="loading"
          variant="outlined"
          color="primary"
          type="submit"
        >
          登録
        </v-btn>
      </v-card-actions>
    </v-form>
  </v-card>
</template>
