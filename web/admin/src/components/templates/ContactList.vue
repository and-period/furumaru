<script lang="ts" setup>
import { VDataTable } from 'vuetify/lib/labs/components'

import { AlertType } from '~/lib/hooks'
import { ContactPriority, ContactStatus, ContactsResponseContactsInner } from '~/types/api'

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
  sortBy: {
    type: Array as PropType<VDataTable['sortBy']>,
    default: () => []
  },
  contacts: {
    type: Array<ContactsResponseContactsInner>,
    default: () => []
  },
  tableItemsPerPage: {
    type: Number,
    default: 20
  },
  tableItemsTotal: {
    type: Number,
    default: 0
  }
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
    key: 'title'
  },
  {
    title: 'メールアドレス',
    key: 'email'
  },
  {
    title: '優先度',
    key: 'priority'
  },
  {
    title: 'ステータス',
    key: 'status'
  }
]

const getPriorityColor = (priority: ContactPriority): string => {
  switch (priority) {
    case ContactPriority.LOW:
      return 'primary'
    case ContactPriority.MIDDLE:
      return 'secondary'
    case ContactPriority.HIGH:
      return 'error'
    default:
      return 'unknown'
  }
}

const getPriority = (priority: ContactPriority): string => {
  switch (priority) {
    case ContactPriority.LOW:
      return '低'
    case ContactPriority.MIDDLE:
      return '中'
    case ContactPriority.HIGH:
      return '高'
    default:
      return '未設定'
  }
}

const getStatusColor = (status: ContactStatus): string => {
  switch (status) {
    case ContactStatus.TODO:
      return 'error'
    case ContactStatus.INPROGRESS:
      return 'secondary'
    case ContactStatus.DONE:
      return 'primary'
    case ContactStatus.DISCARD:
      return 'info'
    default:
      return 'unknown'
  }
}

const getStatus = (status: ContactStatus): string => {
  switch (status) {
    case ContactStatus.TODO:
      return '未着手'
    case ContactStatus.INPROGRESS:
      return '進行中'
    case ContactStatus.DONE:
      return '対応完了'
    case ContactStatus.DISCARD:
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

const onClickRow = (contact: ContactsResponseContactsInner): void => {
  emit('click:row', contact.id)
}
</script>

<template>
  <v-alert v-show="props.isAlert" :type="props.alertType" v-text="props.alertText" />

  <v-card>
    <v-card-title>お問い合わせ管理</v-card-title>

    <v-card-text>
      <v-data-table-server
        :headers="headers"
        :items="props.contacts"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
        :sort-by="props.sortBy"
        :multi-sort="true"
        hover
        @update:page="onClickUpdatePage"
        @update:items-per-page="onClickUpdateItemsPerPage"
        @update:sort-by="onClickUpdateSortBy"
        @click:row="(_:any, {item}: any) => onClickRow(item.raw)"
      >
        <template #[`item.priority`]="{ item }">
          <v-chip :color="getPriorityColor(item.raw.priority)" size="small">
            {{ getPriority(item.raw.priority) }}
          </v-chip>
        </template>
        <template #[`item.status`]="{ item }">
          <v-chip :color="getStatusColor(item.raw.status)" size="small">
            {{ getStatus(item.raw.status) }}
          </v-chip>
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
