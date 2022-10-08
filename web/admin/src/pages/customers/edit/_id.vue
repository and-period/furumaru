<template>
  <div>
    <v-card-title>顧客管理</v-card-title>
    <v-tabs v-model="tab" grow color="dark">
      <v-tabs-slider color="accent"></v-tabs-slider>
      <v-tab
        v-for="tabItem in tabItems"
        :key="tabItem.value"
        :href="`#${tabItem.value}`"
      >
        {{ tabItem.name }}
      </v-tab>
    </v-tabs>

    <v-tabs-items v-model="tab">
      <v-tab-item value="customers">
        <v-card-text>
          <v-text-field
            name="name"
            label="名前"
            value="testName"
            readonly
          ></v-text-field>

          <span>アカウントの有無：</span>
          <v-chip color="primary">{{ accountValue }}</v-chip>

          <v-text-field
            name="email"
            label="連絡先:Email"
            value="and-period-assistant@gmail.com"
            readonly
          ></v-text-field>

          <v-text-field
            name="telePhone"
            label="連絡先:電話番号"
            value="09012345678"
            readonly
          ></v-text-field>

          <v-text-field
            name="prefecture"
            label="都道府県"
            value="広島県"
            readonly
          ></v-text-field>

          <v-text-field
            name="city"
            label="市区町村"
            value="豊田郡大崎上島町"
            readonly
          ></v-text-field>

          <v-text-field
            name="address"
            label="番地"
            value="2067番地1"
            readonly
          ></v-text-field>

          <v-text-field
            name="building"
            label="建物名・部屋番号"
            value="フルマルビルディング 201"
            readonly
          ></v-text-field>
        </v-card-text>
      </v-tab-item>

      <v-tab-item value="customerItems">
        <v-data-table :headers="headers" :items="items">
          <template #foot>
            <tfoot>
              <tr>
                <td :colspan="headers.length - 1">合計</td>
                <td>¥21,000</td>
              </tr>
            </tfoot>
          </template>
        </v-data-table>
      </v-tab-item>
    </v-tabs-items>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from '@nuxtjs/composition-api'

import { Customer } from '~/types/props/customer'
export default defineComponent({
  setup() {
    const tab = ref<string>('customers')
    const tabItems: Customer[] = [
      { name: '顧客情報', value: 'customers' },
      { name: '購入に関して', value: 'customerItems' },
    ]

    const accountValue = '有'

    const headers = [
      {
        text: '関連マルシェ',
        value: 'marche',
      },
      {
        text: '開催日時',
        value: 'date',
      },
      {
        text: '購入金額',
        value: 'price',
      },
    ]

    const items = [
      {
        marche: '大崎上島マルシェ',
        date: '2023/04/01 14:34',
        price: 1000,
      },
      {
        marche: '福岡マルシェ',
        date: '2023/04/02 14:34',
        price: 2000,
      },
      {
        marche: '大崎上島マルシェ',
        date: '2023/04/01 14:34',
        price: 1000,
      },
      {
        marche: '福岡マルシェ',
        date: '2023/04/02 14:34',
        price: 2000,
      },
      {
        marche: '大崎上島マルシェ',
        date: '2023/04/01 14:34',
        price: 1000,
      },
      {
        marche: '福岡マルシェ',
        date: '2023/04/02 14:34',
        price: 2000,
      },
      {
        marche: '大崎上島マルシェ',
        date: '2023/04/01 14:34',
        price: 1000,
      },
      {
        marche: '福岡マルシェ',
        date: '2023/04/02 14:34',
        price: 2000,
      },
    ]

    return {
      accountValue,
      headers,
      items,
      tabItems,
      tab,
    }
  },
})
</script>
