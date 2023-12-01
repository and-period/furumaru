<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import dayjs, { unix } from 'dayjs'
import type { AlertType } from '~/lib/hooks'

import { getErrorMessage, maxLength, required } from '~/lib/validations'
import { AdminRole, DiscountType, type Notification, NotificationStatus, NotificationTarget, NotificationType, type Promotion, type UpdateNotificationRequest } from '~/types/api'
import type { NotificationTime } from '~/types/props'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  role: {
    type: Number as PropType<AdminRole>,
    default: AdminRole.UNKNOWN
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
    type: Object as PropType<UpdateNotificationRequest>,
    default: (): UpdateNotificationRequest => ({
      targets: [],
      title: '',
      body: '',
      note: '',
      publishedAt: dayjs().unix()
    })
  },
  notification: {
    type: Object as PropType<Notification>,
    default: (): Notification => ({
      id: '',
      type: NotificationType.UNKNOWN,
      status: NotificationStatus.UNKNOWN,
      targets: [],
      title: '',
      body: '',
      note: '',
      publishedAt: dayjs().unix(),
      promotionId: '',
      createdBy: '',
      createdAt: 0,
      updatedBy: '',
      updatedAt: 0
    })
  },
  promotion: {
    type: Object as PropType<Promotion>,
    default: (): Promotion => ({
      id: '',
      title: '',
      description: '',
      public: false,
      discountType: DiscountType.UNKNOWN,
      discountRate: 0,
      code: '',
      startAt: 0,
      endAt: 0,
      createdAt: 0,
      updatedAt: 0
    })
  }
})

const emit = defineEmits<{
  (e: 'update:form-data', formData: UpdateNotificationRequest): void
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

const formDataRules = computed(() => ({
  targets: {},
  title: { maxLength: maxLength(128) },
  body: { required, maxLength: maxLength(2000) },
  note: { required, maxLength: maxLength(2000) }
}))
const timeDataRules = computed(() => ({
  publishedDate: {},
  publishedTime: {}
}))
const formDataValue = computed({
  get: (): UpdateNotificationRequest => props.formData as UpdateNotificationRequest,
  set: (formData: UpdateNotificationRequest) => emit('update:form-data', formData)
})
const timeDataValue = computed({
  get: (): NotificationTime => ({
    publishedDate: unix(props.formData.publishedAt).format('YYYY-MM-DD'),
    publishedTime: unix(props.formData.publishedAt).format('HH:mm')
  }),
  set: (timeData: NotificationTime): void => {
    const publishedAt = dayjs(`${timeData.publishedDate} ${timeData.publishedTime}`)
    formDataValue.value.publishedAt = publishedAt.unix()
  }
})
const notificationValue = computed((): Notification => {
  return props.notification
})

const formDataValidate = useVuelidate(formDataRules, formDataValue)
const timeDataValidate = useVuelidate(timeDataRules, timeDataValue)

const isEditable = (): boolean => {
  return props.role === AdminRole.ADMINISTRATOR
}

const onChangePublishedAt = (): void => {
  const publishedAt = dayjs(`${timeDataValue.value.publishedDate} ${timeDataValue.value.publishedTime}`)
  formDataValue.value.publishedAt = publishedAt.unix()
}

const getDateTime = (unixTime: number): string => {
  if (unixTime === 0) {
    return ''
  }
  return unix(unixTime).format('YYYY/MM/DD HH:mm')
}

const getPromotionTerm = (): string => {
  if (!props.promotion) {
    return ''
  }

  const startAt = getDateTime(props.promotion.startAt)
  const endAt = getDateTime(props.promotion.endAt)
  return `${startAt} ${endAt}`
}

const getPromotionDiscount = (): string => {
  if (!props.promotion) {
    return ''
  }

  switch (props.promotion.discountType) {
    case DiscountType.AMOUNT:
      return '￥' + props.promotion.discountRate
    case DiscountType.RATE:
      return props.promotion.discountRate + '％'
    case DiscountType.FREE_SHIPPING:
      return '送料無料'
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
    <v-card-title>お知らせ編集</v-card-title>

    <v-form @submit.prevent="onSubmit">
      <v-card-text>
        <v-select
          v-model="notificationValue.type"
          :items="typeList"
          label="お知らせ種別"
          item-title="title"
          item-value="value"
          readonly
        />
        <!-- セール情報 -->
        <div v-if="notification.type === NotificationType.PROMOTION">
          <v-table>
            <tbody>
              <tr>
                <td>タイトル</td>
                <td>{{ promotion?.title || '' }}</td>
              </tr>
              <tr>
                <td>割引コード</td>
                <td>{{ promotion?.code || '' }}</td>
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
            v-model="timeDataValidate.publishedDate.$model"
            :error-messages="getErrorMessage(timeDataValidate.publishedDate.$errors)"
            type="date"
            class="mr-2"
            variant="outlined"
            density="compact"
            @update:model-value="onChangePublishedAt"
          />
          <v-text-field
            v-model="timeDataValidate.publishedTime.$model"
            :error-messages="getErrorMessage(timeDataValidate.publishedTime.$errors)"
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
