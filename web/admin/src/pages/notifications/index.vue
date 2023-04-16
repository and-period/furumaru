<script lang="ts" setup>
import { mdiPlus, mdiPencil, mdiDelete } from '@mdi/js'
import * as dayjs from 'dayjs'
import { VDataTable } from 'vuetify/labs/components'

import { usePagination } from '~/lib/hooks'
import { useNotificationStore } from '~/store'
import { NotificationsResponseNotificationsInner } from '~/types/api'

const router = useRouter()
const notificationStore = useNotificationStore()
const {
  itemsPerPage,
  offset,
  options,
  updateCurrentPage,
  handleUpdateItemsPerPage
} = usePagination()

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

const fetchState = useAsyncData(async () => {
  await fetchNotifications()
})

const deleteDialog = ref<boolean>(false)
const selectedId = ref<string>('')
const selectedName = ref<string>('')
const sortBy = ref<VDataTable['sortBy']>([])

const notifications = computed(() => {
  return notificationStore.notifications
})
const total = computed(() => {
  return notificationStore.totalItems
})

watch(itemsPerPage, () => {
  fetchNotifications()
})

const handleUpdatePage = async (page: number) => {
  updateCurrentPage(page)
  await fetchNotifications()
}

const fetchNotifications = async () => {
  try {
    await notificationStore.fetchNotifications(
      itemsPerPage.value,
      offset.value
    )
  } catch (err) {
    console.log(err)
  }
}

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
  return dayjs.unix(unixTime).format('YYYY/MM/DD HH:mm')
}

const openDeleteDialog = (
  item: NotificationsResponseNotificationsInner
): void => {
  selectedId.value = item.id
  selectedName.value = item.title
  deleteDialog.value = true
}

const hideDeleteDialog = () => {
  deleteDialog.value = false
}

const handleClickAddButton = () => {
  router.push('/notifications/add')
}

const handleEdit = (item: NotificationsResponseNotificationsInner) => {
  router.push(`/notifications/edit/${item.id}`)
}

const handleDelete = async (): Promise<void> => {
  try {
    await notificationStore.deleteNotification(selectedId.value)
  } catch (err) {
    console.log(err)
  }
  deleteDialog.value = false
}

const isLoading = (): boolean => {
  return fetchState?.pending?.value || false
}

try {
  fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <div>
    <v-card-title class="d-flex flex-row">
      お知らせ管理
      <v-spacer />
      <v-btn color="primary" variant="outlined" @click="handleClickAddButton">
        <v-icon start :icon="mdiPlus" />
        お知らせ登録
      </v-btn>
    </v-card-title>

    <v-dialog v-model="deleteDialog" width="500">
      <v-card>
        <v-card-title class="text-h7">
          {{ selectedName }}を本当に削除しますか？
        </v-card-title>
        <v-card-actions>
          <v-spacer />
          <v-btn color="error" variant="text" @click="hideDeleteDialog">
            キャンセル
          </v-btn>
          <v-btn color="primary" variant="outlined" @click="handleDelete">
            削除
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-card class="mt-4" flat :loading="isLoading()">
      <v-card-text>
        <v-data-table-server
          v-model:sort-by="sortBy"
          :headers="headers"
          :items="notifications"
          :items-per-page="itemsPerPage"
          :items-length="total"
          :footer-props="options"
          :multi-sort="true"
          no-data-text="登録されているお知らせ情報がありません"
          @update:page="handleUpdatePage"
          @update:items-per-page="handleUpdateItemsPerPage"
          @update:sort-by="fetchState.refresh"
          @update:sort-desc="fetchState.refresh"
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
            <v-btn class="mr-2" variant="outlined" color="primary" size="small" @click="handleEdit(item.raw)">
              <v-icon size="small" :icon="mdiPencil" />
              編集
            </v-btn>
            <v-btn
              variant="outlined"
              color="primary"
              size="small"
              @click="openDeleteDialog(item.raw)"
            >
              <v-icon size="small" :icon="mdiDelete" />
              削除
            </v-btn>
          </template>
        </v-data-table-server>
      </v-card-text>
    </v-card>
  </div>
</template>
