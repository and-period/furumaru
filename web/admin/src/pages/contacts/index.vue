<template>
  <div>
    <v-card-title>お問い合わせ管理</v-card-title>
    <v-card>
      <v-card-text>
        <v-data-table :headers="headers" :items="contacts" :items-per-page="5">
          <template #[`item.priority`]="{ item }">
            <v-chip :color="getPriorityColor(item.priority)" dark>
              {{ getPriority(item.priority) }}
            </v-chip>
          </template>
          <template #[`item.status`]="{ item }">
            <v-chip :color="getStatusColor(item.status)" dark>
              {{ getStatus(item.status) }}
            </v-chip>
          </template>
          <template #[`item.actions`]="{ item }">
            <v-btn outlined color="primary" small @click="handleEdit(item)">
              <v-icon small>mdi-pencil</v-icon>
              編集
            </v-btn>
          </template>
        </v-data-table>
      </v-card-text>
    </v-card>
  </div>
</template>

<script lang="ts">
import { defineComponent, useRouter } from '@nuxtjs/composition-api'
import { DataTableHeader } from 'vuetify'

export default defineComponent({
  setup() {
    const router = useRouter()
    const headers: DataTableHeader[] = [
      {
        text: '件名',
        value: 'subject',
      },
      {
        text: 'メールアドレス',
        value: 'mailAddress',
      },
      {
        text: '優先度',
        value: 'priority',
      },
      {
        text: 'ステータス',
        value: 'status',
      },
      {
        text: 'メモ',
        value: 'memo',
      },
      {
        text: 'Actions',
        value: 'actions',
        sortable: false,
      },
    ]
    const contacts = [
      {
        subject: '商品が届かない件について',
        mailAddress: 'and-period@gmail.com',
        priority: 1,
        status: 1,
        memo: '明日配送します',
        actions: 'あくしょん',
      },
      {
        subject: '商品が届かない件について',
        mailAddress: 'and-period@gmail.com',
        priority: 2,
        status: 2,
        memo: '明日配送します',
        actions: 'あくしょん',
      },
      {
        subject: '商品が届かない件について',
        mailAddress: 'and-period@gmail.com',
        priority: 3,
        status: 3,
        memo: '明日配送します',
        actions: 'あくしょん',
      },
      {
        subject: '商品が届かない件について',
        mailAddress: 'and-period@gmail.com',
        priority: 4,
        status: 4,
        memo: '明日配送します',
        actions: 'あくしょん',
      },
    ]

    const getPriorityColor = (priority: any): string => {
      switch (priority) {
        case 1:
          return 'red'
        case 2:
          return 'orange'
        case 3:
          return 'blue'
        default:
          return ''
      }
    }

    const getPriority = (priority: any): string => {
      switch (priority) {
        case 1:
          return 'High'
        case 2:
          return 'Middle'
        case 3:
          return 'Low'
        default:
          return 'Unknown'
      }
    }

    const getStatusColor = (status: any): string => {
      switch (status) {
        case 1:
          return 'red'
        case 2:
          return 'orange'
        case 3:
          return 'blue'
        default:
          return ''
      }
    }

    const getStatus = (status: any): string => {
      switch (status) {
        case 1:
          return '未着手'
        case 2:
          return '進行中'
        case 3:
          return '完了'
        default:
          return '不明'
      }
    }

    const handleEdit = () => {
      router.push('/contacts/edit')
    }

    return {
      headers,
      contacts,
      getPriority,
      getPriorityColor,
      getStatus,
      getStatusColor,
      handleEdit,
    }
  },
})
</script>
