<script lang="ts" setup>
import { mdiDelete, mdiPlus } from '@mdi/js'
import { unix } from 'dayjs'
import type { VDataTable } from 'vuetify/lib/components/index.mjs'

import type { AlertType } from '~/lib/hooks'
import type { Admin, Notification, NotificationTarget } from '~/types/api'
import { AdminRole, NotificationStatus, NotificationType } from '~/types/api'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  role: {
    type: Number as PropType<AdminRole>,
    default: AdminRole.UNKNOWN,
  },
  deleteDialog: {
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
  (e: 'update:delete-dialog', v: boolean): void
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

const selectedItem = ref<Notification>()

const deleteDialogValue = computed({
  get: (): boolean => props.deleteDialog,
  set: (val: boolean): void => emit('update:delete-dialog', val),
})

const isRegisterable = (): boolean => {
  return props.role === AdminRole.ADMINISTRATOR
}

const isEditable = (): boolean => {
  return props.role === AdminRole.ADMINISTRATOR
}

const getAdminName = (adminId: string): string => {
  const admin = props.admins.find((admin: Admin): boolean => {
    return admin.id === adminId
  })
  return admin ? `${admin.lastname} ${admin.firstname}` : ''
}

const getType = (type: NotificationType): string => {
  switch (type) {
    case NotificationType.SYSTEM:
      return 'システム関連'
    case NotificationType.LIVE:
      return 'ライブ関連'
    case NotificationType.PROMOTION:
      return 'セール関連'
    case NotificationType.OTHER:
      return 'その他'
    default:
      return 'その他'
  }
}

const getStatus = (status: NotificationStatus): string => {
  switch (status) {
    case NotificationStatus.WAITING:
      return '投稿前'
    case NotificationStatus.NOTIFIED:
      return '投稿済み'
    default:
      return '不明'
  }
}

const getStatusColor = (status: NotificationStatus): string => {
  switch (status) {
    case NotificationStatus.WAITING:
      return 'error'
    case NotificationStatus.NOTIFIED:
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

const onClickOpenDeleteDialog = (notification: Notification): void => {
  selectedItem.value = notification
  deleteDialogValue.value = true
}

const onClickCloseDeleteDialog = (): void => {
  deleteDialogValue.value = false
}

const onClickAdd = (): void => {
  emit('click:add')
}

const onClickDelete = (): void => {
  emit('click:delete', selectedItem?.value?.id || '')
}
</script>

<template>
  <v-dialog
    v-model="deleteDialogValue"
    width="500"
  >
    <v-card>
      <v-card-title class="text-h7">
        {{ selectedItem?.title || '' }}を本当に削除しますか？
      </v-card-title>
      <v-card-actions>
        <v-spacer />
        <v-btn
          color="error"
          variant="text"
          @click="onClickCloseDeleteDialog"
        >
          キャンセル
        </v-btn>
        <v-btn
          color="primary"
          variant="outlined"
          :loading="loading"
          @click="onClickDelete"
        >
          削除
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-card
    class="mt-4"
    flat
  >
    <v-card-title class="d-flex flex-row">
      お知らせ管理
      <v-spacer />
      <v-btn
        v-show="isRegisterable()"
        color="primary"
        variant="outlined"
        @click="onClickAdd"
      >
        <v-icon
          start
          :icon="mdiPlus"
        />
        お知らせ登録
      </v-btn>
    </v-card-title>

    <v-card-text>
      <v-data-table-server
        :headers="headers"
        :loading="loading"
        :items="notifications"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
        :sort-by="props.tableSortBy"
        :multi-sort="true"
        hover
        no-data-text="登録されているお知らせ情報がありません"
        @update:page="onClickUpdatePage"
        @update:items-per-page="onClickUpdateItemsPerPage"
        @update:sort-by="onClickUpdateSortBy"
        @update:sort-desc="onClickUpdateSortBy"
        @click:row="(_: any, { item }: any) => onClickRow(item.id)"
      >
        <template #[`item.type`]="{ item }">
          {{ getType(item.type) }}
        </template>
        <template #[`item.status`]="{ item }">
          <v-chip
            size="small"
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
            color="primary"
            size="small"
            @click.stop="onClickOpenDeleteDialog(item)"
          >
            <v-icon
              size="small"
              :icon="mdiDelete"
            />
            削除
          </v-btn>
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
