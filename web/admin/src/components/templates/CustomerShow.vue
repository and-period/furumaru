<script lang="ts" setup>
import { VDataTable } from 'vuetify/lib/labs/components'
import { convertI18nToJapanesePhoneNumber } from '~/lib/formatter'
import { UserResponse } from '~/types/api'
import { Customer } from '~/types/props/customer'

const props = defineProps({
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
    title: '関連マルシェ',
    key: 'marche'
  },
  {
    title: '開催日時',
    key: 'date'
  },
  {
    title: '購入金額',
    key: 'price'
  }
]
// dummy data
const orders = [
  {
    marche: '大崎上島マルシェ',
    date: '2023/04/01 14:34',
    price: 1000
  },
  {
    marche: '福岡マルシェ',
    date: '2023/04/02 14:34',
    price: 2000
  },
  {
    marche: '大崎上島マルシェ',
    date: '2023/04/01 14:34',
    price: 1000
  },
  {
    marche: '福岡マルシェ',
    date: '2023/04/02 14:34',
    price: 2000
  },
  {
    marche: '大崎上島マルシェ',
    date: '2023/04/01 14:34',
    price: 1000
  },
  {
    marche: '福岡マルシェ',
    date: '2023/04/02 14:34',
    price: 2000
  },
  {
    marche: '大崎上島マルシェ',
    date: '2023/04/01 14:34',
    price: 1000
  },
  {
    marche: '福岡マルシェ',
    date: '2023/04/02 14:34',
    price: 2000
  }
]

const selector = ref<string>('customers')

const customerName = computed((): string => {
  return `${props.customer.lastname} ${props.customer.firstname}`
})
const customerNameKana = computed((): string => {
  return `${props.customer.lastnameKana} ${props.customer.firstnameKana}`
})
const customerPhoneNumber = computed((): string => {
  return convertI18nToJapanesePhoneNumber(props.customer.phoneNumber)
})
const orderTotal = computed((): number => {
  let total = 0
  orders.forEach((order) => {
    total += order.price
  })
  return total
})
const customerValue = computed({
  get: (): UserResponse => props.customer,
  set: (customer: UserResponse): void => emit('update:customer', customer)
})
</script>

<template>
  <v-card>
    <v-card-title>顧客管理</v-card-title>

    <v-tabs v-model="selector" grow color="dark">
      <v-tab v-for="item in tabs" :key="item.value" :value="item.value">
        {{ item.name }}
      </v-tab>
    </v-tabs>

    <v-card-text>
      <v-window v-model="selector">
        <v-window-item value="customers">
          <v-text-field v-model="customerName" name="name" label="名前" readonly />

          <v-text-field v-model="customerNameKana" name="nameKana" label="名前（かな）" readonly />

          <span>アカウントの有無：</span>
          <v-chip color="primary">
            {{ customerValue.registered ? '有' : '無' }}
          </v-chip>

          <v-text-field
            v-model="customerValue.email"
            name="email"
            label="連絡先:Email"
            readonly
          />

          <v-text-field
            v-model="customerPhoneNumber"
            name="telePhone"
            label="連絡先:電話番号"
            readonly
          />

          <v-text-field
            v-model="customerValue.prefecture"
            name="prefecture"
            label="都道府県"
            readonly
          />

          <v-text-field
            v-model="customerValue.city"
            name="city"
            label="市区町村"
            readonly
          />

          <v-text-field
            v-model="customerValue.addressLine1"
            name="address"
            label="番地"
            readonly
          />

          <v-text-field
            v-model="customerValue.addressLine2"
            name="building"
            label="建物名・部屋番号"
            readonly
          />
        </v-window-item>

        <v-window-item value="customerItems">
          <div class="d-flex mb-4">
            <v-spacer />
            <div>
              <span>合計</span>
              <span>&yen;{{ orderTotal }}</span>
            </div>
          </div>
          <v-data-table :headers="headers" :items="orders" />
        </v-window-item>
      </v-window>
    </v-card-text>
  </v-card>
</template>
