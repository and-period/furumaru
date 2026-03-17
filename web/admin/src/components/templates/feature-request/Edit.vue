<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import dayjs from 'dayjs'
import type { AlertType } from '~/lib/hooks'
import { getErrorMessage, maxLength, required } from '~/lib/validations'
import { useFeatureRequestLabels } from '~/composables/useFeatureRequestLabels'
import { useDeleteDialog } from '~/composables/useDeleteDialog'
import {
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

const { getStatusLabel, getStatusColor, getCategoryLabel, getPriorityLabel } = useFeatureRequestLabels()
const { dialogVisible: deleteDialogVisible, open: openDeleteDialog, close: closeDeleteDialog } = useDeleteDialog<string>()

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
              @click="openDeleteDialog(featureRequest.id)"
            >
              削除
            </v-btn>
          </div>
        </v-form>
      </template>
    </v-card-text>
  </v-card>

  <v-dialog
    v-model="deleteDialogVisible"
    width="500"
  >
    <v-card>
      <v-card-title class="text-h7">
        この要望リクエストを本当に削除しますか？
      </v-card-title>
      <v-card-actions>
        <v-spacer />
        <v-btn
          color=""
          variant="text"
          @click="closeDeleteDialog"
        >
          キャンセル
        </v-btn>
        <v-btn
          :loading="loading"
          color="error"
          variant="text"
          @click="closeDeleteDialog(); emit('click:delete')"
        >
          削除
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
