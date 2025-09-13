<script lang="ts" setup>
import { mdiDelete } from '@mdi/js'
import type { VDataTable } from 'vuetify/lib/components/index.mjs'

import { prefecturesList } from '~/constants'
import type { PrefecturesListItem } from '~/constants'
import type { AlertType } from '~/lib/hooks'
import { UserStatus } from '~/types/api/v1'
import type { UserToList } from '~/types/api/v1'

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
  customers: {
    type: Array<UserToList>,
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
  tableSortBy: {
    type: Array as PropType<VDataTable['sortBy']>,
    default: () => [],
  },
})

const emit = defineEmits<{
  (e: 'click:update-page', page: number): void
  (e: 'click:update-items-per-page', page: number): void
  (e: 'click:row', customerId: string): void
  (e: 'update:sort-by', sortBy: VDataTable['sortBy']): void
}>()

const headers: VDataTable['headers'] = [
  {
    title: '氏名',
    key: 'name',
    sortable: false,
  },
  {
    title: '住所',
    key: 'address',
    sortable: false,
  },
  {
    title: '注文数',
    key: 'paymentTotalCount',
    sortable: false,
  },
  {
    title: 'ステータス',
    key: 'status',
    sortable: false,
  },
]

const getName = (item: UserToList): string => {
  if (item.lastname || item.firstname) {
    return `${item.lastname} ${item.firstname}`
  }
  return item.email
}

const getStatus = (status: UserStatus): string => {
  switch (status) {
    case UserStatus.UserStatusGuest:
      return 'ゲスト'
    case UserStatus.UserStatusProvisional:
      return '仮登録'
    case UserStatus.UserStatusVerified:
      return '認証済み'
    case UserStatus.UserStatusDeactivated:
      return '退会済み'
    default:
      return '不明'
  }
}

const getStatusColor = (status: UserStatus): string => {
  switch (status) {
    case UserStatus.UserStatusGuest:
      return 'secondary'
    case UserStatus.UserStatusProvisional:
      return 'warning'
    case UserStatus.UserStatusVerified:
      return 'primary'
    case UserStatus.UserStatusDeactivated:
      return 'error'
    default:
      return 'unknown'
  }
}

const getAddress = (customer: UserToList): string => {
  const prefecture = prefecturesList.find((prefecture: PrefecturesListItem): boolean => {
    return prefecture.value === customer.prefectureCode
  })
  return `${prefecture?.text || ''} ${customer.city}`
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

const onClickRow = (item: UserToList): void => {
  emit('click:row', item.id || '')
}
</script>

<template>
  <v-alert
    v-show="props.isAlert"
    :type="props.alertType"
    v-text="props.alertText"
  />

  <v-card flat>
    <v-card-title>顧客管理</v-card-title>

    <v-card-text>
      <v-data-table-server
        :headers="headers"
        :loading="loading"
        :items="props.customers"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
        :sort-by="props.tableSortBy"
        no-data-text="登録されている顧客情報がありません"
        hover
        @update:page="onClickUpdatePage"
        @update:items-per-page="onClickUpdateItemsPerPage"
        @update:sort-by="onClickUpdateSortBy"
        @update:sort-desc="onClickUpdateSortBy"
        @click:row="(_: any, { item }: any) => onClickRow(item)"
      >
        <template #[`item.name`]="{ item }">
          {{ getName(item) }}
        </template>
        <template #[`item.address`]="{ item }">
          {{ getAddress(item) }}
        </template>
        <template #[`item.status`]="{ item }">
          <v-chip
            size="small"
            :color="getStatusColor(item.status)"
          >
            {{ getStatus(item.status) }}
          </v-chip>
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
