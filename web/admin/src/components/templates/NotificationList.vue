<script lang="ts" setup>
import { mdiDelete, mdiPlus } from '@mdi/js'
import { unix } from 'dayjs'
import { VDataTable } from 'vuetify/lib/labs/components'
import { AlertType } from '~/lib/hooks'
import { NotificationsResponseNotificationsInner } from '~/types/api'

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
  notifications: {
    type: Array<NotificationsResponseNotificationsInner>,
    defualt: () => []
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
  (e: 'click:row', notificationId: string): void
  (e: 'click:add'): void
  (e: 'click:delete', notificationId: string): void
  (e: 'update:sort-by', sortBy: VDataTable['sortBy']): void
}>()

const headers: VDataTable['headers'] = [
  {
    title: 'タイトル',
    key: 'title'
  },
  {
    title: '公開状況',
    key: 'public'
  },
  {
    title: '投稿範囲',
    key: 'targets'
  },
  {
    title: '掲載開始時間',
    key: 'publishedAt'
  },
  {
    title: 'Actions',
    key: 'actions',
    sortable: false
  }
]

const dialog = ref<boolean>(false)
const selectedItem = ref<NotificationsResponseNotificationsInner>()

const getStatusColor = (status: boolean): string => {
  if (status) {
    return 'primary'
  } else {
    return 'error'
  }
}

const getPublic = (isPublic: boolean): string => {
  if (isPublic) {
    return '公開'
  } else {
    return '非公開'
  }
}

const getTarget = (targets: number[]): string => {
  const actors: string[] = targets?.map((target: number): string => {
    switch (target) {
      case 1:
        return 'ユーザー'
      case 2:
        return '生産者'
      case 3:
        return 'コーディネータ'
      default:
        return ''
    }
  }) || []
  return actors.join(', ')
}

const getDay = (unixTime: number): string => {
  return unix(unixTime).format('YYYY/MM/DD HH:mm')
}

const onClickOpen = (notification: NotificationsResponseNotificationsInner): void => {
  selectedItem.value = notification
  dialog.value = true
}

const onClickClose = (): void => {
  dialog.value = false
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
  emit('click:delete', selectedItem?.value?.id || '')
}
</script>

<template>
  <v-card-title class="d-flex flex-row">
    お知らせ管理
    <v-spacer />
    <v-btn color="primary" variant="outlined" @click="onClickAdd">
      <v-icon start :icon="mdiPlus" />
      お知らせ登録
    </v-btn>
  </v-card-title>

  <v-dialog v-model="dialog" width="500">
    <v-card>
      <v-card-title class="text-h7">
        {{ selectedItem?.title || '' }}を本当に削除しますか？
      </v-card-title>
      <v-card-actions>
        <v-spacer />
        <v-btn color="error" variant="text" @click="onClickClose">
          キャンセル
        </v-btn>
        <v-btn color="primary" variant="outlined" @click="onClickDelete">
          削除
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-card class="mt-4" flat :loading="props.loading">
    <v-card-text>
      <v-data-table-server
        v-model:sort-by="props.tableSortBy"
        :headers="headers"
        :items="notifications"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
        :multi-sort="true"
        hover
        no-data-text="登録されているお知らせ情報がありません"
        @update:page="onClickUpdatePage"
        @update:items-per-page="onClickUpdateItemsPerPage"
        @update:sort-by="onClickUpdateSortBy"
        @update:sort-desc="onClickUpdateSortBy"
        @click:row="(_: any, { item }: any) => onClickRow(item.raw.id)"
      >
        <template #[`item.public`]="{ item }">
          <v-chip size="small" :color="getStatusColor(item.raw.public)">
            {{ getPublic(item.raw.public) }}
          </v-chip>
        </template>
        <template #[`item.targets`]="{ item }">
          {{ getTarget(item.raw.targets) }}
        </template>
        <template #[`item.publishedAt`]="{ item }">
          {{ getDay(item.raw.publishedAt) }}
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn
            variant="outlined"
            color="primary"
            size="small"
            @click.stop="onClickOpen(item.raw)"
          >
            <v-icon size="small" :icon="mdiDelete" />
            削除
          </v-btn>
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
