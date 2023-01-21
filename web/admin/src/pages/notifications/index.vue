<template>
  <div>
    <v-card-title>
      お知らせ管理
      <v-spacer />
      <v-btn outlined color="primary" @click="handleClickAddButton">
        <v-icon left>mdi-plus</v-icon>
        お知らせ登録
      </v-btn>
    </v-card-title>

    <v-dialog v-model="deleteDialog" width="500">
      <v-card>
        <v-card-title class="text-h7">
          {{ selectedName }}を本当に削除しますか？
        </v-card-title>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="error" text @click="hideDeleteDialog">
            キャンセル
          </v-btn>
          <v-btn color="primary" outlined @click="handleDelete"> 削除 </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-card class="mt-4" flat>
      <v-card-text>
        <v-data-table
          :headers="headers"
          :items="notifications"
          no-data-text="登録されているお知らせ情報がありません"
        >
          <template #[`item.public`]="{ item }">
            <v-chip small :color="getStatusColor(item.public)">
              {{ getPublic(item.public) }}
            </v-chip>
          </template>
          <template #[`item.targets`]="{ item }">
            {{ getTarget(item.targets) }}
          </template>
          <template #[`item.publishedAt`]="{ item }">
            {{ getDay(item.publishedAt) }}
          </template>
          <template #[`item.actions`]="{ item }">
            <v-btn outlined color="primary" small @click="handleEdit(item)">
              <v-icon small>mdi-pencil</v-icon>
              編集
            </v-btn>
            <v-btn
              outlined
              color="primary"
              small
              @click="openDeleteDialog(item)"
            >
              <v-icon small>mdi-delete</v-icon>
              削除
            </v-btn>
          </template>
        </v-data-table>
      </v-card-text>
    </v-card>
  </div>
</template>

<script lang="ts">
import { computed, useFetch, useRouter } from '@nuxtjs/composition-api'
import { defineComponent, ref } from '@vue/composition-api'
import dayjs from 'dayjs'
import { DataTableHeader } from 'vuetify'

import { useNotificationStore } from '~/store/notification'
import { NotificationsResponseNotificationsInner } from '~/types/api'

export default defineComponent({
  setup() {
    const router = useRouter()
    const notificationStore = useNotificationStore()

    const deleteDialog = ref<boolean>(false)
    const selectedId = ref<string>('')
    const selectedName = ref<string>('')
    const notifications = computed(() => {
      return notificationStore.notifications
    })

    const headers: DataTableHeader[] = [
      {
        text: 'タイトル',
        value: 'title',
      },
      {
        text: '公開状況',
        value: 'public',
      },
      {
        text: '投稿範囲',
        value: 'targets',
      },
      {
        text: '掲載開始時間',
        value: 'publishedAt',
      },
      {
        text: 'Actions',
        value: 'actions',
        sortable: false,
      },
    ]

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
      const actors: string[] = targets.map((target: number): string => {
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
      })
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

    const { fetchState } = useFetch(async () => {
      try {
        await notificationStore.fetchNotifications()
      } catch (err) {
        console.log(err)
      }
    })

    return {
      headers,
      notifications,
      fetchState,
      selectedName,
      deleteDialog,
      openDeleteDialog,
      hideDeleteDialog,
      getStatusColor,
      handleClickAddButton,
      handleEdit,
      handleDelete,
      getPublic,
      getTarget,
      getDay,
    }
  },
})
</script>
