<script lang="ts" setup>
import { VDataTable } from 'vuetify/lib/labs/components'
import { PrefecturesListItem, prefecturesList } from '~/constants'
import { convertI18nToJapanesePhoneNumber } from '~/lib/formatter'
import { UserResponse } from '~/types/api'
import { Customer } from '~/types/props/customer'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  customer: {
    type: Object as PropType<UserResponse>,
    default: () => ({})
  }
})

const emit = defineEmits<{
  (e: 'update:customer', customer: UserResponse): void
}>()

const tabs: Customer[] = [
  { name: '顧客情報', value: 'customers' },
  { name: '購入に関して', value: 'customerItems' }
]
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
    return prefecture.value === props.customer.prefecture
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
    <v-col sm="12" md="12" lg="8">
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
      <v-card elevation="0">
        <v-card-text>
          <v-card-title class="pb-4">
            購入情報
          </v-card-title>
          <v-data-table :headers="headers" :items="orders" />
        </v-card-text>
      </v-card>
    </v-col>
    <v-col sm="12" md="12" lg="4">
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
  </v-row>
</template>
