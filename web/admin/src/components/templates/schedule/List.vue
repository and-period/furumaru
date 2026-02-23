<script lang="ts" setup>
import { mdiDelete, mdiPencil, mdiPlus, mdiCalendarToday, mdiPlayCircle } from '@mdi/js'
import { unix } from 'dayjs'
import type { VDataTable } from 'vuetify/components'
import { scheduleStatuses } from '~/constants'

import { getResizedImages } from '~/lib/helpers'
import type { AlertType } from '~/lib/hooks'
import {

  ScheduleStatus,

  AdminType,
} from '~/types/api/v1'
import type { Coordinator, Schedule } from '~/types/api/v1'

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
  sortBy: {
    type: Array as PropType<VDataTable['sortBy']>,
    default: () => [],
  },
  coordinators: {
    type: Array<Coordinator>,
    default: () => [],
  },
  schedules: {
    type: Array<Schedule>,
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
  (e: 'click:row', scheduleId: string): void
  (e: 'click:add'): void
  (e: 'click:delete', scheduleId: string): void
  (e: 'click:approval', scheduleId: string): void
  (e: 'click:published', scheduleId: string): void
}>()

const headers: VDataTable['headers'] = [
  {
    title: '',
    key: 'thumbnail',
    sortable: false,
  },
  {
    title: 'マルシェ名',
    key: 'title',
    sortable: false,
  },
  {
    title: 'コーディネーター名',
    key: 'coordinatorName',
    sortable: false,
  },
  {
    title: 'ステータス',
    key: 'status',
    sortable: false,
  },
  {
    title: '開催期間',
    key: 'term',
    sortable: false,
  },
  {
    title: '',
    key: 'actions',
    sortable: false,
  },
]

const { dialogVisible, selectedItem, open: openDeleteDialog, close: closeDeleteDialog } = useDeleteDialog<Schedule>()

const isRegisterable = (): boolean => {
  return props.adminType === AdminType.AdminTypeCoordinator
}

const getCoordinatorName = (coordinatorId: string): string => {
  const coordinator = props.coordinators.find(
    (coordinator: Coordinator): boolean => {
      return coordinator.id === coordinatorId
    },
  )
  return coordinator ? coordinator.username : ''
}

const getThumbnail = (schedule: Schedule): string => {
  return schedule.thumbnailUrl || ''
}

const getResizedThumbnails = (schedule: Schedule): string => {
  if (!schedule.thumbnailUrl) {
    return ''
  }
  return getResizedImages(schedule.thumbnailUrl)
}

const getStatus = (status: ScheduleStatus): string => {
  const value = scheduleStatuses.find(s => s.value === status)
  return value?.title || '不明'
}

const getStatusColor = (status: ScheduleStatus): string => {
  switch (status) {
    case ScheduleStatus.ScheduleStatusPrivate:
      return 'error'
    case ScheduleStatus.ScheduleStatusInProgress:
      return 'warning'
    case ScheduleStatus.ScheduleStatusWaiting:
      return 'info'
    case ScheduleStatus.ScheduleStatusLive:
      return 'primary'
    case ScheduleStatus.ScheduleStatusClosed:
      return 'secondary'
    default:
      return ''
  }
}

const getDay = (unixTime: number): string => {
  return unix(unixTime).format('YYYY/MM/DD HH:mm')
}

const getTerm = (schedule: Schedule): string => {
  return `${getDay(schedule.startAt)} ~ ${getDay(schedule.endAt)}`
}

const getPublished = (schedule: Schedule): string => {
  return schedule._public ? '非公開にする' : '公開する'
}

const getPublishedColor = (schedule: Schedule): string => {
  return schedule._public ? 'secondary' : 'primary'
}

const onClickUpdatePage = (page: number): void => {
  emit('click:update-page', page)
}

const onClickUpdateItemsPerPage = (page: number): void => {
  emit('click:update-items-per-page', page)
}

const onClickRow = (scheduleId: string): void => {
  emit('click:row', scheduleId)
}

const onClickAdd = (): void => {
  emit('click:add')
}

const onClickDelete = (): void => {
  emit('click:delete', selectedItem.value?.id || '')
  closeDeleteDialog()
}

const onClickPublished = (scheduleId: string): void => {
  emit('click:published', scheduleId)
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
    title="ライブ配信削除の確認"
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
          :icon="mdiCalendarToday"
          size="28"
          class="mr-3 text-primary"
        />
        <div>
          <h1 class="text-h6 text-sm-h5 font-weight-bold text-primary">
            ライブ配信管理
          </h1>
          <p class="text-caption text-sm-body-2 text-medium-emphasis ma-0">
            ライブ配信の登録・編集・削除を行います
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
          ライブ配信登録
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
        :items="props.schedules"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
        hover
        no-data-text="登録されているスケジュールがありません。"
        @update:page="onClickUpdatePage"
        @update:items-per-page="onClickUpdateItemsPerPage"
        @click:row="(_: any, { item }: any) => onClickRow(item.id)"
      >
        <template #[`item.thumbnail`]="{ item }">
          <v-avatar size="40">
            <v-img
              v-if="getThumbnail(item) !== ''"
              cover
              :src="getThumbnail(item)"
              :srcset="getResizedThumbnails(item)"
              :alt="item.title || 'スケジュール画像'"
            />
            <v-icon
              v-else
              :icon="mdiPlayCircle"
              color="grey"
            />
          </v-avatar>
        </template>
        <template #[`item.status`]="{ item }">
          <v-chip :color="getStatusColor(item.status)">
            {{ getStatus(item.status) }}
          </v-chip>
        </template>
        <template #[`item.coordinatorName`]="{ item }">
          {{ getCoordinatorName(item.coordinatorId) }}
        </template>
        <template #[`item.term`]="{ item }">
          {{ getTerm(item) }}
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn
            variant="outlined"
            class="mr-2"
            size="small"
            :color="getPublishedColor(item)"
            @click.stop="onClickPublished(item.id)"
          >
            <v-icon
              size="small"
              :icon="mdiPencil"
            />
            {{ getPublished(item) }}
          </v-btn>
          <v-btn
            variant="outlined"
            color="error"
            size="small"
            @click.stop="openDeleteDialog(item)"
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
