<template>
  <div>
    <v-card-title>お問合せ管理</v-card-title>
    <v-card elevation="0">
      <v-col cols="4">
        <v-text-field
          name="name"
          label="名前"
          value="古丸太郎"
          readonly
        ></v-text-field>
      </v-col>

      <v-col cols="4">
        <v-text-field
          name="subject"
          label="件名"
          value="商品が届かない件について"
          readonly
        ></v-text-field>
      </v-col>

      <v-col cols="6">
        <v-textarea
          name="contact"
          label="お問合せ内容"
          value="商品が届いていないです。どうなっていますでしょうか。ああああああああああああああああああああ"
          readonly
        ></v-textarea>
      </v-col>

      <v-row>
        <v-col cols="2">
          <v-subheader>優先度：</v-subheader>
        </v-col>
        <v-col cols="5">
          <v-autocomplete
            :items="priority"
            chips
            :item-color="red"
            label="優先度"
          ></v-autocomplete>
        </v-col>
      </v-row>

      <v-row>
        <v-col cols="2">
          <v-subheader>ステータス：</v-subheader>
        </v-col>
        <v-col cols="5">
          <v-autocomplete
            :items="status"
            chips
            label="ステータス"
          ></v-autocomplete>
        </v-col>
      </v-row>

      <v-col cols="4">
        <v-text-field
          name="mailAddress"
          label="メールアドレス"
          value="and-period@gmail.com"
          readonly
        ></v-text-field>
      </v-col>

      <v-col cols="4">
        <v-text-field
          name="telephoneNumber"
          label="電話番号"
          value="09012345678"
          readonly
        ></v-text-field>
      </v-col>

      <v-col cols="6">
        <v-textarea name="memo" label="メモ"></v-textarea>
      </v-col>
    </v-card>
  </div>
</template>

<script lang="ts">
export default {
  setup() {
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
    return {
      //        selected: ['Trevor Handsen'],
      priority: ['High', 'Middle', 'Low', 'Unknown'],
      status: ['未着手', '進行中', '完了', '不明'],
      }
  },
}
</script>
