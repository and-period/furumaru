<script setup lang="ts">
import { mdiPlus } from '@mdi/js'
import { unix } from 'dayjs'
import { usePagination } from '~/lib/hooks'
import { useAuthStore, useVideoStore } from '~/store'
import { videoStatusToString, videoStatusToColor } from '~/lib/formatter'
import { AdminRole } from '~/types/api'

const videoStore = useVideoStore()
const { videoResponse } = storeToRefs(videoStore)

const pagination = usePagination()
const router = useRouter()

const authStore = useAuthStore()
const { role } = storeToRefs(authStore)

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
    key: '',
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
</script>

<template>
  <v-card>
    <v-card-title>動画管理</v-card-title>

    <v-card-text
      v-if="role === AdminRole.COORDINATOR"
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
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
