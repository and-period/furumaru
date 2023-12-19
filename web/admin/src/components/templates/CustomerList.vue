<script lang="ts" setup>
import { mdiDelete } from '@mdi/js'
import { VDataTable } from 'vuetify/lib/labs/components.mjs'
import { type PrefecturesListItem, prefecturesList } from '~/constants'
import type { AlertType } from '~/lib/hooks'
import { UserStatus, type UserToList } from '~/types/api'

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
  customers: {
    type: Array<UserToList>,
    default: () => []
  },
  tableItemsPerPage: {
    type: Number,
    default: 20
  },
  tableItemsTotal: {
    type: Number,
    default: 0
  },
  tableSortBy: {
    type: Array as PropType<VDataTable['sortBy']>,
    default: () => []
  }
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
    sortable: false
  },
  {
    title: '住所',
    key: 'address',
    sortable: false
  },
  {
    title: '注文数',
    key: 'totalOrder',
    sortable: false
  },
  {
    title: '購入金額',
    key: 'totalAmount',
    sortable: false
  },
  {
    title: 'ステータス',
    key: 'status',
    sortable: false
  }
]

const getName = (item: UserToList): string => {
  if (item.lastname || item.firstname) {
    return `${item.lastname} ${item.firstname}`
  }
  return item.email
}

const getStatus = (status: UserStatus): string => {
  switch (status) {
    case UserStatus.GUEST:
      return 'ゲスト'
    case UserStatus.PROVISIONAL:
      return '仮登録'
    case UserStatus.VERIFIED:
      return '認証済み'
    case UserStatus.WITH_DRAWAL:
      return '退会済み'
    default:
      return '不明'
  }
}

const getStatusColor = (status: UserStatus): string => {
  switch (status) {
    case UserStatus.GUEST:
      return 'secondary'
    case UserStatus.PROVISIONAL:
      return 'warning'
    case UserStatus.VERIFIED:
      return 'primary'
    case UserStatus.WITH_DRAWAL:
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
  <v-alert v-show="props.isAlert" :type="props.alertType" v-text="props.alertText" />

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
        <template #[`item.totalAmount`]="{ item }">
          &yen; {{ `${item.totalAmount.toLocaleString()}` }}
        </template>
        <template #[`item.status`]="{ item }">
          <v-chip size="small" :color="getStatusColor(item.status)">
            {{ getStatus(item.status) }}
          </v-chip>
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
