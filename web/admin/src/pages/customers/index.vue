<template>
  <div>
    <v-card-title>顧客管理</v-card-title>
    <v-card>
      <v-card-text>
        <v-data-table :headers="headers" :items="items" :items-per-page="5">
          <template #[`item.account`]="{ item }">
            <v-chip :color="getAccountColor(item.account)" dark>
              {{ getAccount(item.account) }}
            </v-chip>
          </template>
          <template #[`item.action1`]>
            <v-btn outlined color="primary" small @click="handleEdit()">
              <v-icon small>mdi-pencil</v-icon>
              編集
            </v-btn>
          </template>
          <template #[`item.action2`]>
            <v-btn outlined color="primary" small>
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
import { defineComponent, useRouter } from '@nuxtjs/composition-api'

export default defineComponent({
  setup() {
    const router = useRouter()
    const headers = [
      {
        text: '名前',
        value: 'name',
      },
      {
        text: '住所',
        value: 'address',
      },
      {
        text: '購入数',
        value: 'quantity',
      },
      {
        text: '購入金額',
        value: 'price',
      },
      {
        text: 'アカウントの有無',
        value: 'account',
      },
      {
        text: 'Action1',
        value: 'action1',
      },
      {
        text: 'Action2',
        value: 'action2',
      },
    ]

    const items = [
      {
        name: 'namae',
        address: 'juu',
        quantity: 1,
        price: 1000,
        account: 0,
      },
      {
        name: 'namae2',
        address: 'juu2',
        quantity: 2,
        price: 2000,
        account: 1,
      },
    ]

    const getAccountColor = (account: any): string => {
      switch (account) {
        case 0:
          return 'red'
        case 1:
          return 'primary'
        default:
          return ''
      }
    }

    const getAccount = (account: any): string => {
      switch (account) {
        case 0:
          return '無'
        case 1:
          return '有'
        default:
          return '無'
      }
    }

    const handleEdit = () => {
      router.push(`/customers/edit/id`)
    }

    return {
      headers,
      items,
      getAccountColor,
      getAccount,
      handleEdit,
    }
  },
})
</script>

