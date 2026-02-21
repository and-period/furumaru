<script lang="ts" setup>
import { mdiDelete, mdiPlus, mdiPlayCircle } from '@mdi/js'
import { unix } from 'dayjs'
import type { VDataTable } from 'vuetify/components'

import { getResizedImages } from '~/lib/helpers'
import type { AlertType } from '~/lib/hooks'
import { videoStatusToString, videoStatusToColor } from '~/lib/formatter'
import {
  AdminType,
} from '~/types/api/v1'
import type { Video } from '~/types/api/v1'

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
  videos: {
    type: Array<Video>,
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
  (e: 'click:row', videoId: string): void
  (e: 'click:add'): void
  (e: 'click:delete', videoId: string): void
}>()

const headers: VDataTable['headers'] = [
  {
    title: '',
    key: 'thumbnail',
    sortable: false,
  },
  {
    title: 'タイトル',
    key: 'title',
    sortable: false,
  },
  {
    title: '公開日時',
    key: 'publishedAt',
    sortable: false,
  },
  {
    title: 'ステータス',
    key: 'status',
    sortable: false,
  },
  {
    title: '',
    key: 'actions',
    sortable: false,
  },
]

const { dialogVisible, selectedItem, open: openDeleteDialog, close: closeDeleteDialog } = useDeleteDialog<Video>()

const isRegisterable = (): boolean => {
  return props.adminType === AdminType.AdminTypeCoordinator
}

const isDeletable = (): boolean => {
  const targets: AdminType[] = [AdminType.AdminTypeAdministrator, AdminType.AdminTypeCoordinator]
  return targets.includes(props.adminType)
}

const getThumbnail = (video: Video): string => {
  return video.thumbnailUrl || ''
}

const getResizedThumbnails = (video: Video): string => {
  if (!video.thumbnailUrl) {
    return ''
  }
  return getResizedImages(video.thumbnailUrl)
}

const getPublishedAt = (publishedAt: number): string => {
  if (publishedAt === 0) {
    return '未設定'
  }
  return unix(publishedAt).format('YYYY/MM/DD HH:mm')
}

const onClickUpdatePage = (page: number): void => {
  emit('click:update-page', page)
}

const onClickUpdateItemsPerPage = (page: number): void => {
  emit('click:update-items-per-page', page)
}

const onClickRow = (videoId: string): void => {
  emit('click:row', videoId)
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
    title="動画削除の確認"
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
          :icon="mdiPlayCircle"
          size="28"
          class="mr-3 text-primary"
        />
        <div>
          <h1 class="text-h6 text-sm-h5 font-weight-bold text-primary">
            動画管理
          </h1>
          <p class="text-caption text-sm-body-2 text-medium-emphasis ma-0">
            動画コンテンツの登録・編集・削除を行います
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
          動画登録
        </v-btn>
      </div>
    </v-card-title>

    <v-card-text>
      <v-data-table-server
        :headers="headers"
        :loading="loading"
        :items="props.videos"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
        hover
        no-data-text="登録されている動画がありません。"
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
            />
            <v-icon
              v-else
              :icon="mdiPlayCircle"
              color="grey"
            />
          </v-avatar>
        </template>
        <template #[`item.title`]="{ item }">
          <div class="text-body-1 font-weight-medium">
            {{ item.title }}
          </div>
          <div class="text-body-2 text-medium-emphasis">
            {{ item.description?.substring(0, 50) }}{{ item.description?.length > 50 ? '...' : '' }}
          </div>
        </template>
        <template #[`item.publishedAt`]="{ item }">
          <div class="text-body-2">
            {{ getPublishedAt(item.publishedAt) }}
          </div>
        </template>
        <template #[`item.status`]="{ item }">
          <v-chip
            :color="videoStatusToColor(item.status)"
            size="small"
          >
            {{ videoStatusToString(item.status) }}
          </v-chip>
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn
            v-show="isDeletable()"
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
