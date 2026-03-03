<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import dayjs from 'dayjs'
import type { AlertType } from '~/lib/hooks'
import { getErrorMessage, maxLength, required } from '~/lib/validations'
import {
  FeatureRequestCategory,
  FeatureRequestPriority,
  FeatureRequestStatus,
} from '~/types/feature-request'
import type { FeatureRequest, UpdateFeatureRequestInput } from '~/types/feature-request'
import { AdminType } from '~/types/api/v1'

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
  adminType: {
    type: Number as PropType<AdminType>,
    default: AdminType.AdminTypeUnknown,
  },
  featureRequest: {
    type: Object as PropType<FeatureRequest>,
    default: () => ({
      id: '',
      title: '',
      description: '',
      category: FeatureRequestCategory.Feature,
      priority: FeatureRequestPriority.Medium,
      status: FeatureRequestStatus.Waiting,
      note: '',
      submittedBy: '',
      submitterName: '',
      createdAt: 0,
      updatedAt: 0,
    }),
  },
  formData: {
    type: Object as PropType<UpdateFeatureRequestInput>,
    default: () => ({
      status: FeatureRequestStatus.Waiting,
      note: '',
    }),
  },
})

const emit = defineEmits<{
  (e: 'update:form-data', formData: UpdateFeatureRequestInput): void
  (e: 'submit'): void
  (e: 'click:delete'): void
}>()

const isAdmin = computed(
  () => props.adminType === AdminType.AdminTypeAdministrator,
)

const statuses = [
  { title: '受付中', value: FeatureRequestStatus.Waiting },
  { title: '検討中', value: FeatureRequestStatus.Reviewing },
  { title: '採用決定', value: FeatureRequestStatus.Adopted },
  { title: '開発中', value: FeatureRequestStatus.InProgress },
  { title: '完了', value: FeatureRequestStatus.Done },
  { title: '却下', value: FeatureRequestStatus.Rejected },
]

const getCategoryLabel = (category: FeatureRequestCategory): string => {
  switch (category) {
    case FeatureRequestCategory.UI: return 'UI/UX改善'
    case FeatureRequestCategory.Feature: return '新機能'
    case FeatureRequestCategory.Performance: return 'パフォーマンス改善'
    case FeatureRequestCategory.Other: return 'その他'
    default: return '不明'
  }
}

const getPriorityLabel = (priority: FeatureRequestPriority): string => {
  switch (priority) {
    case FeatureRequestPriority.Low: return '低'
    case FeatureRequestPriority.Medium: return '中'
    case FeatureRequestPriority.High: return '高'
    default: return '不明'
  }
}

const getStatusLabel = (status: FeatureRequestStatus): string => {
  switch (status) {
    case FeatureRequestStatus.Waiting: return '受付中'
    case FeatureRequestStatus.Reviewing: return '検討中'
    case FeatureRequestStatus.Adopted: return '採用決定'
    case FeatureRequestStatus.InProgress: return '開発中'
    case FeatureRequestStatus.Done: return '完了'
    case FeatureRequestStatus.Rejected: return '却下'
    default: return '不明'
  }
}

const getStatusColor = (status: FeatureRequestStatus): string => {
  switch (status) {
    case FeatureRequestStatus.Waiting: return 'warning'
    case FeatureRequestStatus.Reviewing: return 'info'
    case FeatureRequestStatus.Adopted: return 'secondary'
    case FeatureRequestStatus.InProgress: return 'primary'
    case FeatureRequestStatus.Done: return 'success'
    case FeatureRequestStatus.Rejected: return 'error'
    default: return 'default'
  }
}

const formatDate = (unixTime: number): string => {
  if (!unixTime)
    return '-'
  return dayjs.unix(unixTime).format('YYYY/MM/DD HH:mm')
}

const formDataValue = computed({
  get: (): UpdateFeatureRequestInput => props.formData,
  set: (v: UpdateFeatureRequestInput): void => emit('update:form-data', v),
})

const rules = computed(() => ({
  status: { required },
  note: { maxLength: maxLength(2000) },
}))

const validate = useVuelidate(rules, formDataValue)

const onSubmit = async (): Promise<void> => {
  const valid = await validate.value.$validate()
  if (!valid)
    return
  emit('submit')
}
</script>

<template>
  <atoms-app-alert
    :show="props.isAlert"
    :type="props.alertType"
    :text="props.alertText"
  />

  <v-card elevation="0">
    <v-card-title class="d-flex align-center justify-space-between flex-wrap ga-2">
      要望リクエスト詳細
      <v-chip
        :color="getStatusColor(featureRequest.status)"
        variant="elevated"
      >
        {{ getStatusLabel(featureRequest.status) }}
      </v-chip>
    </v-card-title>

    <v-card-text>
      <!-- 要望情報（読み取り専用） -->
      <v-card
        variant="tonal"
        class="mb-4"
      >
        <v-card-text class="d-flex flex-column ga-3">
          <div>
            <div class="text-caption text-medium-emphasis mb-1">
              タイトル
            </div>
            <div class="text-body-1 font-weight-medium">
              {{ featureRequest.title }}
            </div>
          </div>

          <v-row>
            <v-col
              cols="6"
              sm="3"
            >
              <div class="text-caption text-medium-emphasis mb-1">
                カテゴリ
              </div>
              <div class="text-body-2">
                {{ getCategoryLabel(featureRequest.category) }}
              </div>
            </v-col>
            <v-col
              cols="6"
              sm="3"
            >
              <div class="text-caption text-medium-emphasis mb-1">
                優先度
              </div>
              <div class="text-body-2">
                {{ getPriorityLabel(featureRequest.priority) }}
              </div>
            </v-col>
            <v-col
              cols="6"
              sm="3"
            >
              <div class="text-caption text-medium-emphasis mb-1">
                提出者
              </div>
              <div class="text-body-2">
                {{ featureRequest.submitterName }}
              </div>
            </v-col>
            <v-col
              cols="6"
              sm="3"
            >
              <div class="text-caption text-medium-emphasis mb-1">
                提出日時
              </div>
              <div class="text-body-2">
                {{ formatDate(featureRequest.createdAt) }}
              </div>
            </v-col>
          </v-row>

          <div>
            <div class="text-caption text-medium-emphasis mb-1">
              要望内容
            </div>
            <div
              class="text-body-2"
              style="white-space: pre-wrap;"
            >
              {{ featureRequest.description }}
            </div>
          </div>
        </v-card-text>
      </v-card>

      <!-- 管理者コメント表示（読み取り専用 for coordinator） -->
      <template v-if="!isAdmin && featureRequest.note">
        <div class="text-subtitle-2 mb-2">
          管理者コメント
        </div>
        <v-card
          variant="outlined"
          class="mb-4"
        >
          <v-card-text>
            <div
              class="text-body-2"
              style="white-space: pre-wrap;"
            >
              {{ featureRequest.note }}
            </div>
          </v-card-text>
        </v-card>
      </template>

      <!-- 管理者向け編集フォーム -->
      <template v-if="isAdmin">
        <v-divider class="mb-4" />
        <div class="text-subtitle-1 font-weight-bold mb-3">
          ステータス管理
        </div>

        <v-form @submit.prevent="onSubmit">
          <v-select
            v-model="validate.status.$model"
            :error-messages="getErrorMessage(validate.status.$errors)"
            :items="statuses"
            item-title="title"
            item-value="value"
            label="ステータス"
            class="mb-2"
          />

          <v-textarea
            v-model="validate.note.$model"
            :error-messages="getErrorMessage(validate.note.$errors)"
            label="管理者コメント（提出者に表示されます）"
            rows="4"
            counter="2000"
            maxlength="2000"
            class="mb-4"
          />

          <div class="d-flex ga-2">
            <v-btn
              :loading="loading"
              type="submit"
              color="primary"
              variant="elevated"
              flex="1"
            >
              更新する
            </v-btn>
            <v-btn
              color="error"
              variant="outlined"
              @click="emit('click:delete')"
            >
              削除
            </v-btn>
          </div>
        </v-form>
      </template>
    </v-card-text>
  </v-card>
</template>
