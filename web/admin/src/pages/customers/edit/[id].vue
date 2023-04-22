<script lang="ts" setup>
import { VDataTable } from 'vuetify/lib/labs/components'
import { Customer } from '~/types/props/customer'

const tab = ref<string>('customers')

const tabItems: Customer[] = [
  { name: '顧客情報', value: 'customers' },
  { name: '購入に関して', value: 'customerItems' }
]

const accountValue = '有'

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

const items = [
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
</script>

<template>
  <div>
    <v-card-title>顧客管理</v-card-title>

    <v-card>
      <v-tabs v-model="tab" grow color="dark">
        <v-tabs-slider color="accent" />
        <v-tab
          v-for="tabItem in tabItems"
          :key="tabItem.value"
          :value="tabItem.value"
        >
          {{ tabItem.name }}
        </v-tab>
      </v-tabs>

      <v-card-text>
        <v-window v-model="tab">
          <v-window-item value="customers">
            <v-text-field name="name" label="名前" value="testName" readonly />

            <span>アカウントの有無：</span>
            <v-chip color="primary">
              {{ accountValue }}
            </v-chip>

            <v-text-field
              name="email"
              label="連絡先:Email"
              value="and-period-assistant@gmail.com"
              readonly
            />

            <v-text-field
              name="telePhone"
              label="連絡先:電話番号"
              value="09012345678"
              readonly
            />

            <v-text-field
              name="prefecture"
              label="都道府県"
              value="広島県"
              readonly
            />

            <v-text-field
              name="city"
              label="市区町村"
              value="豊田郡大崎上島町"
              readonly
            />

            <v-text-field
              name="address"
              label="番地"
              value="2067番地1"
              readonly
            />

            <v-text-field
              name="building"
              label="建物名・部屋番号"
              value="フルマルビルディング 201"
              readonly
            />
          </v-window-item>

          <v-window-item value="customerItems">
            <v-data-table :headers="headers" :items="items">
              <template #foot>
                <tfoot>
                  <tr>
                    <td :colspan="headers.length - 1">
                      合計
                    </td>
                    <td>¥21,000</td>
                  </tr>
                </tfoot>
              </template>
            </v-data-table>
          </v-window-item>
        </v-window>
      </v-card-text>
    </v-card>
  </div>
</template>
