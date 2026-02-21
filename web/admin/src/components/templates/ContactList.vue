<script lang="ts" setup>
import type { VDataTable } from 'vuetify/components'

import type { AlertType } from '~/lib/hooks'
import { ContactStatus } from '~/types/api/v1'
import type { ContactsResponse } from '~/types/api/v1'

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
  sortBy: {
    type: Array as PropType<VDataTable['sortBy']>,
    default: () => [],
  },
  contacts: {
    type: Array<ContactsResponse>,
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
  (e: 'click:update-page', page: number): void
  (e: 'click:update-items-per-page', page: number): void
  (e: 'click:row', contactId: string): void
  (e: 'update:sort-by', sortBy: VDataTable['sortBy']): void
}>()

const headers: VDataTable['headers'] = [
  {
    title: '件名',
    key: 'title',
  },
  {
    title: 'メールアドレス',
    key: 'email',
  },
  {
    title: 'ステータス',
    key: 'status',
  },
]

const getStatusColor = (status: ContactStatus): string => {
  switch (status) {
    case ContactStatus.ContactStatusWaiting:
      return 'error'
    case ContactStatus.ContactStatusInprogress:
      return 'secondary'
    case ContactStatus.ContactStatusDone:
      return 'primary'
    case ContactStatus.ContactStatusDiscard:
      return 'info'
    default:
      return 'unknown'
  }
}

const getStatus = (status: ContactStatus): string => {
  switch (status) {
    case ContactStatus.ContactStatusWaiting:
      return '未着手'
    case ContactStatus.ContactStatusInprogress:
      return '進行中'
    case ContactStatus.ContactStatusDone:
      return '対応完了'
    case ContactStatus.ContactStatusDiscard:
      return '対応不要'
    default:
      return '不明'
  }
}

const onClickUpdatePage = (page: number): void => {
  emit('click:update-page', page)
}

const onClickUpdateItemsPerPage = (page: number): void => {
  emit('click:update-items-per-page', page)
}

const onClickUpdateSortBy = (sortBy: VDataTable['sortBy']): void => {
  emit('update:sort-by', sortBy)
}

const onClickRow = (contactId: string): void => {
  emit('click:row', contactId)
}
</script>

<template>
  <v-alert
    v-show="props.isAlert"
    :type="props.alertType"
    v-text="props.alertText"
  />

  <v-card>
    <v-card-title>お問い合わせ管理</v-card-title>

    <v-card-text>
      <v-data-table-server
        :headers="headers"
        :loading="loading"
        :items="props.contacts"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
        :sort-by="props.sortBy"
        :multi-sort="true"
        hover
        @update:page="onClickUpdatePage"
        @update:items-per-page="onClickUpdateItemsPerPage"
        @update:sort-by="onClickUpdateSortBy"
        @click:row="(_:any, { item }: any) => onClickRow(item.id)"
      >
        <template #[`item.status`]="{ item }">
          <v-chip
            :color="getStatusColor(item.status)"
            size="small"
          >
            {{ getStatus(item.status) }}
          </v-chip>
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
