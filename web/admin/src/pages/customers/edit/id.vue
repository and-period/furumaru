<template>
  <div>
    <v-card-title>顧客管理</v-card-title>
    <v-card elevation="0">
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

        <v-card-title>購入に関して</v-card-title>
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
      </v-card-text>
    </v-card>
  </div>
</template>

<script lang="ts">
import { defineComponent } from '@nuxtjs/composition-api'
export default defineComponent({
  setup() {
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
    }
  },
})
</script>
