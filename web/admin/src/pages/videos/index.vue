<script setup lang="ts">
import { mdiDelete, mdiPlus } from '@mdi/js'
import { unix } from 'dayjs'
import { useAlert, usePagination } from '~/lib/hooks'
import { useAuthStore, useCommonStore, useVideoStore } from '~/store'
import { videoStatusToString, videoStatusToColor } from '~/lib/formatter'
import { AdminType } from '~/types/api'
import type { Video, VideoResponse } from '~/types/api'

const videoStore = useVideoStore()
const { videoResponse } = storeToRefs(videoStore)

const pagination = usePagination()
const router = useRouter()

const authStore = useAuthStore()
const { adminType } = storeToRefs(authStore)

const selectedItem = ref<Video | null>(null)
const commonStore = useCommonStore()
const deleteDialogValue = ref<boolean>(false)

const { isShow, alertText, alertType, show } = useAlert('error')

const headers: VDataTable['headers'] = [
  {
    title: 'サムネイル',
    key: 'thumbnailUrl',
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

const fetchVideos = async () => {
  await videoStore.fetchVideos()
}

const handleClickNewVideoButton = () => {
  router.push('/videos/new')
}

const getPublishedAt = (publishedAt: number) => {
  if (publishedAt === 0) {
    return '-'
  }
  return unix(publishedAt).format('YYYY/MM/DD HH:mm')
}

const handleClickRow = (id: string) => {
  router.push(`/videos/${id}`)
}

const { status } = useAsyncData(async () => {
  await fetchVideos()
})

const isDeletable = (): boolean => {
  const targets: AdminType[] = [AdminType.ADMINISTRATOR, AdminType.COORDINATOR]
  return targets.includes(adminType.value)
}

const toggleDeleteDialog = (item: Video): void => {
  if (item) {
    selectedItem.value = item
  }
  deleteDialogValue.value = !deleteDialogValue.value
}

const onClickDelete = async () => {
  try {
    if (selectedItem.value) {
      await videoStore.deleteVideo(selectedItem.value.id)
      commonStore.addSnackbar({
        color: 'info',
        message: '動画を削除しました。',
      })
    }
    fetchVideos()
    deleteDialogValue.value = !deleteDialogValue.value
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}
</script>

<template>
  <v-card>
    <v-dialog
      v-model="deleteDialogValue"
      width="500"
    >
      <v-card>
        <v-card-text class="text-h7">
          {{ selectedItem?.title }}を本当に削除しますか？
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn
            color="error"
            variant="text"
            @click="toggleDeleteDialog"
          >
            キャンセル
          </v-btn>
          <v-btn
            color="primary"
            variant="outlined"
            @click="onClickDelete"
          >
            削除
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-card-title>動画管理</v-card-title>

    <v-card-text
      v-if="adminType === AdminType.COORDINATOR"
      class="text-right"
    >
      <v-btn
        color="primary"
        variant="outlined"
        @click="handleClickNewVideoButton"
      >
        <v-icon :icon="mdiPlus" />
        新規動画
      </v-btn>
    </v-card-text>

    <v-card-text>
      <v-data-table-server
        :headers="headers"
        :items="videoResponse?.videos"
        :loading="status === 'pending'"
        :items-length="videoResponse?.total"
        :items-per-page="pagination.itemsPerPage.value"
        hover
        @click:row="(_: any, { item }: any) => handleClickRow(item.id)"
      >
        <template #[`item.thumbnailUrl`]="{ item }">
          <v-img
            :src="item.thumbnailUrl"
            width="50"
            aspect-ratio="16/9"
          />
        </template>
        <template #[`item.publishedAt`]="{ item }">
          {{ getPublishedAt(item.publishedAt) }}
        </template>
        <template #[`item.status`]="{ item }">
          <v-chip :color="videoStatusToColor(item.status)">
            {{ videoStatusToString(item.status) }}
          </v-chip>
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn
            v-show="isDeletable()"
            variant="outlined"
            color="primary"
            size="small"
            @click.stop="toggleDeleteDialog(item)"
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
