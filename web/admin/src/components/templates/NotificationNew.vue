<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import dayjs, { unix } from 'dayjs'

import type { AlertType } from '~/lib/hooks'
import { getErrorMessage } from '~/lib/validations'
import { type CreateNotificationRequest, DiscountType, NotificationTarget, NotificationType, type Promotion } from '~/types/api'
import type { DateTimeInput } from '~/types/props'
import { TimeDataValidationRules } from '~/types/validations'
import { CreateNotificationValidationRules } from '~/types/validations/notification'

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
    type: Object,
    default: (): CreateNotificationRequest => ({
      type: NotificationType.UNKNOWN,
      targets: [],
      title: '',
      body: '',
      note: '',
      publishedAt: dayjs().unix(),
      promotionId: ''
    })
  },
  promotions: {
    type: Array<Promotion>,
    default: (): Promotion[] => []
  }
})

const emit = defineEmits<{
  (e: 'update:form-data', formData: CreateNotificationRequest): void
  (e: 'update:notification-type', type: NotificationType): void
  (e: 'update:search-promotion', name: string): void
  (e: 'submit'): void
}>()

const typeList = [
  { title: 'システム関連', value: NotificationType.SYSTEM },
  { title: 'ライブ関連', value: NotificationType.LIVE },
  { title: 'セール関連', value: NotificationType.PROMOTION },
  { title: 'その他', value: NotificationType.OTHER }
]
const targetList = [
  { title: 'ユーザー', value: NotificationTarget.USERS },
  { title: '生産者', value: NotificationTarget.PRODUCERS },
  { title: 'コーディネーター', value: NotificationTarget.COORDINATORS },
  { title: '管理者', value: NotificationTarget.ADMINISTRATORS }
]

const selectedPromotion = ref<Promotion>()

const formDataValue = computed({
  get: (): CreateNotificationRequest => props.formData as CreateNotificationRequest,
  set: (formData: CreateNotificationRequest) => emit('update:form-data', formData)
})
const timeDataValue = computed({
  get: (): DateTimeInput => ({
    date: unix(props.formData.publishedAt).format('YYYY-MM-DD'),
    time: unix(props.formData.publishedAt).format('HH:mm')
  }),
  set: (timeData: DateTimeInput): void => {
    const publishedAt = dayjs(`${timeData.date} ${timeData.time}`)
    formDataValue.value.publishedAt = publishedAt.unix()
  }
})

const formDataValidate = useVuelidate(CreateNotificationValidationRules, formDataValue)
const timeDataValidate = useVuelidate(TimeDataValidationRules, timeDataValue)

const onChangePublishedAt = (): void => {
  const publishedAt = dayjs(`${timeDataValue.value.date} ${timeDataValue.value.time}`)
  formDataValue.value.publishedAt = publishedAt.unix()
}

const getDateTime = (unixTime: number): string => {
  if (unixTime === 0) {
    return ''
  }
  return unix(unixTime).format('YYYY/MM/DD HH:mm')
}

const getPromotionTerm = (): string => {
  if (!selectedPromotion?.value) {
    return ''
  }

  const startAt = getDateTime(selectedPromotion.value.startAt)
  const endAt = getDateTime(selectedPromotion.value.endAt)
  return `${startAt} ${endAt}`
}

const getPromotionDiscount = (): string => {
  if (!selectedPromotion.value) {
    return ''
  }

  switch (selectedPromotion.value.discountType) {
    case DiscountType.AMOUNT:
      return '￥' + selectedPromotion.value.discountRate
    case DiscountType.RATE:
      return selectedPromotion.value.discountRate + '％'
    case DiscountType.FREE_SHIPPING:
      return '送料無料'
    default:
      return ''
  }
}

const onChangeType = (type: NotificationTarget): void => {
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
  <v-alert v-show="props.isAlert" :type="props.alertType" v-text="props.alertText" />

  <v-card>
    <v-card-title>お知らせ登録</v-card-title>

    <v-form @submit.prevent="onSubmit">
      <v-card-text>
        <v-select
          v-model="formDataValidate.type.$model"
          :error-messages="getErrorMessage(formDataValidate.type.$errors)"
          :items="typeList"
          :loading="loading"
          :disabled="loading"
          label="お知らせ種別"
          item-title="title"
          item-value="value"
          @update:model-value="onChangeType"
        />
        <!-- セール情報 -->
        <div v-if="formDataValue.type === NotificationType.PROMOTION">
          <v-autocomplete
            v-model="formDataValidate.promotionId.$model"
            :error-messages="getErrorMessage(formDataValidate.promotionId.$errors)"
            :items="promotions"
            label="セール情報"
            item-title="title"
            item-value="id"
            @update:model-value="onChangePromotion"
            @update:search="onChangeSearchPromotion"
          />
          <v-table>
            <tbody>
              <tr>
                <td>割引コード</td>
                <td>{{ selectedPromotion?.code || '' }}</td>
              </tr>
              <tr>
                <td>割引額</td>
                <td>{{ getPromotionDiscount() }}</td>
              </tr>
              <tr>
                <td>使用期間</td>
                <td>{{ getPromotionTerm() }}</td>
              </tr>
            </tbody>
          </v-table>
        </div>
        <!-- その他 -->
        <v-text-field
          v-else
          v-model="formDataValidate.title.$model"
          :error-messages="getErrorMessage(formDataValidate.title.$errors)"
          label="タイトル"
          required
          maxlength="128"
        />
        <!-- 共通部分 -->
        <v-autocomplete
          v-model="formDataValidate.targets.$model"
          :error-messages="getErrorMessage(formDataValidate.targets.$errors)"
          :items="targetList"
          label="公開範囲"
          multiple
          item-title="title"
          item-value="value"
        />
        <p class="text-subtitle-2 text-grey py-2">
          投稿日時
        </p>
        <div class="d-flex align-center">
          <v-text-field
            v-model="timeDataValidate.date.$model"
            :error-messages="getErrorMessage(timeDataValidate.date.$errors)"
            type="date"
            class="mr-2"
            variant="outlined"
            density="compact"
            @update:model-value="onChangePublishedAt"
          />
          <v-text-field
            v-model="timeDataValidate.time.$model"
            :error-messages="getErrorMessage(timeDataValidate.time.$errors)"
            type="time"
            variant="outlined"
            density="compact"
            @update:model-value="onChangePublishedAt"
          />
        </div>
        <v-textarea
          v-model="formDataValidate.body.$model"
          :error-messages="getErrorMessage(formDataValidate.body.$errors)"
          label="本文"
          placeholder="ユーザーに公開される内容を記載してください"
          maxlength="2000"
        />
        <v-textarea
          v-model="formDataValidate.note.$model"
          :error-messages="getErrorMessage(formDataValidate.note.$errors)"
          label="備考"
          placeholder="ユーザーには非公開にしたいコメント等を記載してください"
          maxlength="2000"
        />
      </v-card-text>

      <v-card-actions>
        <v-btn block :loading="loading" variant="outlined" color="primary" type="submit">
          登録
        </v-btn>
      </v-card-actions>
    </v-form>
  </v-card>
</template>
