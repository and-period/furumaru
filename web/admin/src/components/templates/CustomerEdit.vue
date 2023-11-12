<script lang="ts" setup>
import { VDataTable } from 'vuetify/lib/labs/components.mjs'
import { type PrefecturesListItem, prefecturesList } from '~/constants'
import { convertI18nToJapanesePhoneNumber } from '~/lib/formatter'
import type { User } from '~/types/api'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  customer: {
    type: Object as PropType<User>,
    default: () => ({})
  }
})

const headers: VDataTable['headers'] = [
  {
    title: 'No.',
    key: 'id',
    sortable: false
  },
  {
    title: '関連マルシェ',
    key: 'marche',
    sortable: false
  },
  {
    title: '開催日時',
    key: 'date',
    sortable: false
  },
  {
    title: '購入金額',
    key: 'price',
    sortable: false
  }
]
// dummy data
const orders = [
  {
    id: '0000000001',
    marche: '大崎上島マルシェ',
    date: '2023/04/01 14:34',
    price: 1000
  },
  {
    id: '0000000002',
    marche: '福岡マルシェ',
    date: '2023/04/02 14:34',
    price: 2000
  },
  {
    id: '0000000003',
    marche: '大崎上島マルシェ',
    date: '2023/04/01 14:34',
    price: 1000
  },
  {
    id: '0000000004',
    marche: '福岡マルシェ',
    date: '2023/04/02 14:34',
    price: 2000
  },
  {
    id: '0000000005',
    marche: '大崎上島マルシェ',
    date: '2023/04/01 14:34',
    price: 1000
  },
  {
    id: '0000000006',
    marche: '福岡マルシェ',
    date: '2023/04/02 14:34',
    price: 2000
  },
  {
    id: '0000000007',
    marche: '大崎上島マルシェ',
    date: '2023/04/01 14:34',
    price: 1000
  },
  {
    id: '0000000008',
    marche: '福岡マルシェ',
    date: '2023/04/02 14:34',
    price: 2000
  }
]
const activities = [
  {
    eventType: 'notification',
    detail: '注文(#1000)の発想が完了しました。',
    createdAt: '2023/04/12 10:34'
  },
  {
    eventType: 'notification',
    detail: '注文(#1000)の発送済みメールを送りました。',
    createdAt: '2023/04/10 10:34'
  },
  {
    eventType: 'comment',
    username: 'ふるマル管理者',
    detail: '発送準備をコーディネータに依頼済み',
    createdAt: '2023/04/06 12:00'
  },
  {
    eventType: 'notification',
    detail: '注文(#1000)の支払い完了メールを送りました。',
    createdAt: '2023/04/05 10:34'
  }
]

const getUsername = (): string => {
  return `${props.customer.lastname} ${props.customer.firstname}`
}
const getUsernameKana = (): string => {
  return `${props.customer.lastnameKana} ${props.customer.firstnameKana}`
}
const getPhoneNumber = (): string => {
  return convertI18nToJapanesePhoneNumber(props.customer.phoneNumber)
}
const getAddressArea = (): string => {
  const prefecture = prefecturesList.find((prefecture: PrefecturesListItem): boolean => {
    return prefecture.value === props.customer.prefectureCode
  })
  return prefecture ? `${prefecture.text} ${props.customer.city}` : props.customer.city
}
const getStatus = (): string => {
  return props.customer.registered ? '登録済み' : '未登録'
}
const getStatusColor = (): string => {
  return props.customer.registered ? 'primary' : 'red'
}
</script>

<template>
  <v-row :loading="loading">
    <v-col sm="12" md="12" lg="4" order-lg="2">
      <v-card elevation="0">
        <v-card-text>
          <v-list>
            <v-list-item class="mb-4">
              <v-list-item-subtitle class="pb-2">
                氏名
              </v-list-item-subtitle>
              <div>{{ getUsername() }}</div>
              <div>{{ getUsernameKana() }}</div>
            </v-list-item>
            <v-list-item class="mb-4">
              <v-list-item-subtitle class="mb-2">
                登録状況
              </v-list-item-subtitle>
              <v-chip size="small" :color="getStatusColor()">
                {{ getStatus() }}
              </v-chip>
            </v-list-item>
            <v-list-item class="mb-4">
              <v-list-item-subtitle class="pb-2">
                連絡先情報
              </v-list-item-subtitle>
              <div>{{ props.customer.email }}</div>
              <div>{{ getPhoneNumber() }}</div>
            </v-list-item>
            <v-list-item>
              <v-list-item-subtitle class="pb-2">
                請求先情報
              </v-list-item-subtitle>
              <div>&#12306; {{ props.customer.postalCode }}</div>
              <div>{{ getAddressArea() }}</div>
              <div>{{ props.customer.addressLine1 }}</div>
              <div>{{ props.customer.addressLine2 }}</div>
            </v-list-item>
          </v-list>
        </v-card-text>
      </v-card>
    </v-col>

    <v-col sm="12" md="12" lg="8" order-lg="1">
      <v-card elevation="0" class="mb-4">
        <v-card-text>
          <v-row>
            <v-col>
              <v-card-subtitle class="pb-4">
                支払い金額
              </v-card-subtitle>
              <div class="px-4">
                &yen; 0
              </div>
            </v-col>
            <v-col>
              <v-card-subtitle class="pb-4">
                注文数
              </v-card-subtitle>
              <div class="px-4">
                0
              </div>
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
      <v-card elevation="0" class="mb-4">
        <v-card-text>
          <v-card-title class="pb-4">
            購入情報
          </v-card-title>
          <v-data-table :headers="headers" :items="orders" />
        </v-card-text>
      </v-card>

      <div class="pa-4">
        <h4 class="pb-2">
          タイムライン
        </h4>
        <v-divider />

        <v-timeline side="end" density="compact">
          <template v-for="(activity, i) in activities" :key="i">
            <v-timeline-item v-if="activity.eventType === 'notification'" class="mb-4" dot-color="grey" size="small" max-width="75vw">
              <div class="d-flex flex-column flex-lg-row justify-space-between flex-grow-1">
                <div>{{ activity.detail }}</div>
                <div class="flex-shrink-0 text-grey">
                  {{ activity.createdAt }}
                </div>
              </div>
            </v-timeline-item>
            <v-timeline-item v-if="activity.eventType === 'comment'" class="mb-4" dot-color="grey" size="small" max-width="75vw">
              <template #icon>
                <v-avatar image="https://i.pravatar.cc/64" />
              </template>
              <v-card class="elevation-0">
                <v-card-title class="d-lg-flex flex-lg-row align-center">
                  <div class="pr-2">
                    {{ activity.username }}
                  </div>
                  <div class="text-subtitle-2 text-grey">
                    {{ activity.createdAt }}
                  </div>
                </v-card-title>
                <v-card-text>
                  <div>{{ activity.detail }}</div>
                </v-card-text>
              </v-card>
            </v-timeline-item>
          </template>
        </v-timeline>
      </div>
    </v-col>
  </v-row>
</template>
