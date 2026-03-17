<script lang="ts" setup>
import type { VDataTable } from 'vuetify/components'
import type { AlertType } from '~/lib/hooks'
import { useFeatureRequestLabels } from '~/composables/useFeatureRequestLabels'
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

const { getStatusLabel, getStatusColor, getCategoryLabel, getPriorityLabel, getPriorityColor } = useFeatureRequestLabels()

const headers: VDataTable['headers'] = [
  { title: 'タイトル', key: 'title' },
  { title: 'カテゴリ', key: 'category' },
  { title: '優先度', key: 'priority' },
  { title: 'ステータス', key: 'status' },
  { title: '提出者', key: 'submitterName' },
]
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
        @click:row="(_: Event, { item }: { item: FeatureRequest }) => emit('click:row', item.id)"
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
