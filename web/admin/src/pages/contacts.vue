<template>
  <div>
    <v-card-title>お問い合わせ</v-card-title>
    <v-card>
      <v-card-text>
        <v-data-table :headers="headers" :items="inquiries" :items-per-page="5">
          <template #[`item.priority`]="{ item }">
            <v-chip :color="getPriorityColor(item.priority)" dark>
              {{ item.priority }}
            </v-chip>
          </template>
          <template #[`item.status`]="{ item }">
            <v-chip :color="getStatusColor(item.status)" dark>
              {{ item.status }}
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
import { defineComponent } from '@nuxtjs/composition-api'
import { DataTableHeader } from 'vuetify'
// import { DataTableItemProps } from 'vuetify'

export default defineComponent({
  setup() {
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
    const inquiries = [
      {
        subject: '商品が届かない件について',
        mailAddress: 'and-period@gmail.com',
        priority: 'High',
        status: '未着手',
        memo: '明日配送します',
        actions: 'あくしょん',
      },
      {
        subject: '商品が届かない件について',
        mailAddress: 'and-period@gmail.com',
        priority: 'Middle',
        status: '進行中',
        memo: '明日配送します',
        actions: 'あくしょん',
      },
      {
        subject: '商品が届かない件について',
        mailAddress: 'and-period@gmail.com',
        priority: 'Low',
        status: '完了',
        memo: '明日配送します',
        actions: 'あくしょん',
      },
      {
        subject: '商品が届かない件について',
        mailAddress: 'and-period@gmail.com',
        priority: 'unknown',
        status: '不明',
        memo: '明日配送します',
        actions: 'あくしょん',
      },
    ]
    return {
      headers,
      inquiries,
    }
  },
  methods: {
    getPriorityColor(priority: string) {
      if (priority === 'High') return 'red'
      else if (priority === 'Middle') return 'orange'
      else if (priority === 'Low') return 'blue'
      else return ''
    },
    getStatusColor(status: string) {
      if (status === '未着手') return 'red'
      else if (status === '進行中') return 'orange'
      else if (status === '完了') return 'blue'
      else return ''
    },
  },
})
</script>
