<script lang="ts" setup>
import { mdiAccountGroup } from '@mdi/js'
import type { VDataTable } from 'vuetify/components'

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

  <v-card
    class="mt-6"
    elevation="0"
    rounded="lg"
  >
    <v-card-title class="d-flex align-center justify-space-between pa-6 pb-4">
      <div class="d-flex align-center">
        <v-icon
          :icon="mdiAccountGroup"
          size="28"
          class="mr-3 text-primary"
        />
        <div>
          <h1 class="text-h5 font-weight-bold text-primary">
            顧客管理
          </h1>
          <p class="text-body-2 text-medium-emphasis ma-0">
            顧客情報の確認・管理を行います
          </p>
        </div>
      </div>
    </v-card-title>

    <v-card-text>
      <v-data-table-server
        :headers="headers"
        :loading="loading"
        :items="props.customers"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
        :sort-by="props.tableSortBy"
        hover
        no-data-text="登録されている顧客がありません。"
        @update:page="onClickUpdatePage"
        @update:items-per-page="onClickUpdateItemsPerPage"
        @update:sort-by="onClickUpdateSortBy"
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
            :color="getStatusColor(item.status)"
          >
            {{ getStatus(item.status) }}
          </v-chip>
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
