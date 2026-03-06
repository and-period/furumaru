<script lang="ts" setup>
import type { VDataTable } from 'vuetify/components'
import type { AlertType } from '~/lib/hooks'
import {
  FeatureRequestCategory,
  FeatureRequestPriority,
  FeatureRequestStatus,
} from '~/types/feature-request'
import type { FeatureRequest } from '~/types/feature-request'

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
  featureRequests: {
    type: Array as PropType<FeatureRequest[]>,
    default: () => [],
  },
  tableItemsPerPage: {
    type: Number,
    default: 20,
  },
  tableItemsTotal: {
    type: Number,
    default: 0,
  },
})

const emit = defineEmits<{
  (e: 'click:row', id: string): void
  (e: 'click:new'): void
  (e: 'click:update-page', page: number): void
  (e: 'click:update-items-per-page', perPage: number): void
}>()

const headers: VDataTable['headers'] = [
  { title: 'タイトル', key: 'title' },
  { title: 'カテゴリ', key: 'category' },
  { title: '優先度', key: 'priority' },
  { title: 'ステータス', key: 'status' },
  { title: '提出者', key: 'submitterName' },
]

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

const getPriorityColor = (priority: FeatureRequestPriority): string => {
  switch (priority) {
    case FeatureRequestPriority.Low: return 'default'
    case FeatureRequestPriority.Medium: return 'warning'
    case FeatureRequestPriority.High: return 'error'
    default: return 'default'
  }
}
</script>

<template>
  <atoms-app-alert
    :show="props.isAlert"
    :type="props.alertType"
    :text="props.alertText"
  />

  <v-card>
    <v-card-title class="d-flex align-center justify-space-between">
      要望リクエスト管理
      <v-btn
        color="primary"
        variant="elevated"
        @click="emit('click:new')"
      >
        新規要望を提出
      </v-btn>
    </v-card-title>

    <v-card-text>
      <v-skeleton-loader
        v-if="loading"
        type="table-heading, table-row-divider@5"
      />
      <v-data-table-server
        v-else
        :headers="headers"
        :items="props.featureRequests"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
        hover
        @update:page="emit('click:update-page', $event)"
        @update:items-per-page="emit('click:update-items-per-page', $event)"
        @click:row="(_: any, { item }: any) => emit('click:row', item.id)"
      >
        <template #[`item.category`]="{ item }">
          {{ getCategoryLabel(item.category) }}
        </template>

        <template #[`item.priority`]="{ item }">
          <v-chip
            :color="getPriorityColor(item.priority)"
            size="small"
            variant="tonal"
          >
            {{ getPriorityLabel(item.priority) }}
          </v-chip>
        </template>

        <template #[`item.status`]="{ item }">
          <v-chip
            :color="getStatusColor(item.status)"
            size="small"
          >
            {{ getStatusLabel(item.status) }}
          </v-chip>
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
