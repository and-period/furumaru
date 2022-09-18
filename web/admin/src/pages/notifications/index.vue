<template>
  <div>
    <v-card-title> お知らせ管理 </v-card-title>
    <div class="d-flex">
      <v-spacer />
      <v-btn outlined color="primary" @click="handleClickAddButton">
        <v-icon left>mdi-plus</v-icon>
        お知らせ登録
      </v-btn>
    </div>
    <v-card class="mt-4" flat>
      <v-data-table
        :headers="headers"
        :items="notifications"
        no-data-text="登録されているお知らせ情報がありません"
      >
        <template #[`item.title`]="{ item }">
          {{ item.title }}
        </template>
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
      </v-data-table>
    </v-card>
  </div>
</template>

<script lang="ts">
import { computed, useFetch } from '@nuxtjs/composition-api'
import { defineComponent } from '@vue/composition-api'
import dayjs from 'dayjs'
import { DataTableHeader } from 'vuetify'

import { useNotificationStore } from '~/store/notification'

export default defineComponent({
  setup() {
    const notificationStore = useNotificationStore()
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
    ]

    const getStatusColor = (status: boolean): string => {
      if (status) {
        return 'primary'
      } else {
        return 'accentDarken'
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
            return 'コーディネーター'
          default:
            return ''
        }
      })
      return actors.join(', ')
    }

    const getDay = (unixTime: number): string => {
      return dayjs.unix(unixTime).format('YYYY/MM/DD hh:mm')
    }

    // TODO
    const handleClickAddButton = () => {}

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
      getStatusColor,
      handleClickAddButton,
      getPublic,
      getTarget,
      getDay,
    }
  },
})
</script>
