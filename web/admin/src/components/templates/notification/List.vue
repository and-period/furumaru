<script lang="ts" setup>
import { mdiDelete, mdiPlus, mdiBellOutline } from '@mdi/js'
import { unix } from 'dayjs'
import type { VDataTable } from 'vuetify/components'

import type { AlertType } from '~/lib/hooks'
import type { Admin, Notification, NotificationTarget } from '~/types/api/v1'
import { AdminType, NotificationStatus, NotificationType } from '~/types/api/v1'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  adminType: {
    type: Number as PropType<AdminType>,
    default: AdminType.AdminTypeUnknown,
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
  notifications: {
    type: Array<Notification>,
    default: () => [],
  },
  admins: {
    type: Array<Admin>,
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
  (e: 'click:row', notificationId: string): void
  (e: 'click:add'): void
  (e: 'click:delete', notificationId: string): void
  (e: 'update:sort-by', sortBy: VDataTable['sortBy']): void
}>()

const headers: VDataTable['headers'] = [
  {
    title: 'カテゴリ',
    key: 'type',
    sortable: false,
  },
  {
    title: 'タイトル',
    key: 'title',
    sortable: false,
  },
  {
    title: '状態',
    key: 'status',
    sortable: false,
  },
  {
    title: '投稿範囲',
    key: 'targets',
    sortable: false,
  },
  {
    title: '投稿日時',
    key: 'publishedAt',
    sortable: false,
  },
  {
    title: '作成者',
    key: 'createdBy',
    sortable: false,
  },
  {
    title: '',
    key: 'actions',
    sortable: false,
  },
]

const { dialogVisible, selectedItem, open: openDeleteDialog, close: closeDeleteDialog } = useDeleteDialog<Notification>()

const isRegisterable = (): boolean => {
  return props.adminType === AdminType.AdminTypeAdministrator
}

const isEditable = (): boolean => {
  return props.adminType === AdminType.AdminTypeAdministrator
}

const getAdminName = (adminId: string): string => {
  const admin = props.admins.find((admin: Admin): boolean => {
    return admin.id === adminId
  })
  return admin ? `${admin.lastname} ${admin.firstname}` : ''
}

const getType = (type: NotificationType): string => {
  switch (type) {
    case NotificationType.NotificationTypeSystem:
      return 'システム関連'
    case NotificationType.NotificationTypeLive:
      return 'ライブ関連'
    case NotificationType.NotificationTypePromotion:
      return 'セール関連'
    case NotificationType.NotificationTypeOther:
      return 'その他'
    default:
      return 'その他'
  }
}

const getStatus = (status: NotificationStatus): string => {
  switch (status) {
    case NotificationStatus.NotificationStatusWaiting:
      return '投稿前'
    case NotificationStatus.NotificationStatusNotified:
      return '投稿済み'
    default:
      return '不明'
  }
}

const getStatusColor = (status: NotificationStatus): string => {
  switch (status) {
    case NotificationStatus.NotificationStatusWaiting:
      return 'error'
    case NotificationStatus.NotificationStatusNotified:
      return 'primary'
    default:
      return ''
  }
}

const getTargets = (targets: NotificationTarget[]): string => {
  if (targets.length === 4) {
    return '全員'
  }
  const actors: string[] = targets?.map((target: number): string => {
    switch (target) {
      case 1:
        return 'ユーザー'
      case 2:
        return '生産者'
      case 3:
        return 'コーディネーター'
      case 4:
        return '管理者'
      default:
        return ''
    }
  }) || []
  return actors.join(',')
}

const getDay = (unixTime: number): string => {
  return unix(unixTime).format('YYYY/MM/DD HH:mm')
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

const onClickRow = (notificationId: string): void => {
  emit('click:row', notificationId)
}

const onClickAdd = (): void => {
  emit('click:add')
}

const onClickDelete = (): void => {
  emit('click:delete', selectedItem.value?.id || '')
  closeDeleteDialog()
}
</script>

<template>
  <atoms-app-alert
    :show="props.isAlert"
    :type="props.alertType"
    :text="props.alertText"
  />

  <atoms-app-confirm-dialog
    v-model="dialogVisible"
    title="お知らせ削除の確認"
    :message="`「${selectedItem?.title || ''}」を削除しますか？`"
    :loading="loading"
    @confirm="onClickDelete"
  />

  <v-card
    class="mt-6"
    elevation="0"
    rounded="lg"
  >
    <v-card-title class="d-flex flex-column flex-sm-row align-start align-sm-center justify-space-between pa-4 pa-sm-6 pb-4">
      <div class="d-flex align-center mb-3 mb-sm-0">
        <v-icon
          :icon="mdiBellOutline"
          size="28"
          class="mr-3 text-primary"
        />
        <div>
          <h1 class="text-h6 text-sm-h5 font-weight-bold text-primary">
            お知らせ管理
          </h1>
          <p class="text-caption text-sm-body-2 text-medium-emphasis ma-0">
            お知らせの登録・編集・削除を行います
          </p>
        </div>
      </div>
      <div class="d-flex ga-3 w-100 w-sm-auto">
        <v-btn
          v-show="isRegisterable()"
          variant="elevated"
          color="primary"
          :size="$vuetify.display.smAndDown ? 'default' : 'large'"
          class="w-100 w-sm-auto"
          @click="onClickAdd"
        >
          <v-icon
            start
            :icon="mdiPlus"
          />
          お知らせ登録
        </v-btn>
      </div>
    </v-card-title>

    <v-card-text>
      <v-skeleton-loader
        v-if="loading"
        type="table-heading, table-row-divider@5"
      />
      <v-data-table-server
        v-else
        :headers="headers"
        :items="notifications"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
        :sort-by="props.tableSortBy"
        hover
        no-data-text="登録されているお知らせがありません。"
        @update:page="onClickUpdatePage"
        @update:items-per-page="onClickUpdateItemsPerPage"
        @update:sort-by="onClickUpdateSortBy"
        @click:row="(_: any, { item }: any) => onClickRow(item.id)"
      >
        <template #[`item.type`]="{ item }">
          {{ getType(item.type) }}
        </template>
        <template #[`item.status`]="{ item }">
          <v-chip
            :color="getStatusColor(item.status)"
          >
            {{ getStatus(item.status) }}
          </v-chip>
        </template>
        <template #[`item.targets`]="{ item }">
          {{ getTargets(item.targets) }}
        </template>
        <template #[`item.publishedAt`]="{ item }">
          {{ getDay(item.publishedAt) }}
        </template>
        <template #[`item.createdBy`]="{ item }">
          {{ getAdminName(item.createdBy) }}
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn
            v-show="isEditable()"
            variant="outlined"
            color="error"
            size="small"
            :prepend-icon="mdiDelete"
            @click.stop="openDeleteDialog(item)"
          >
            削除
          </v-btn>
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
