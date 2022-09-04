<template>
  <div>
    <v-card-title>お問合せ管理</v-card-title>
    <v-card elevation="0">
      <v-card-text>
        <v-text-field
          name="name"
          label="名前"
          value="古丸太郎"
          readonly
        ></v-text-field>

        <v-text-field
          name="subject"
          label="件名"
          value="商品が届かない件について"
          readonly
        ></v-text-field>

        <v-textarea
          name="contact"
          label="お問合せ内容"
          value="商品が届いていないです。どうなっていますでしょうか。ああああああああああああああああああああ"
          readonly
        ></v-textarea>

        <v-autocomplete :items="priority" chips label="優先度"></v-autocomplete>

        <v-autocomplete
          :items="status"
          chips
          label="ステータス"
        ></v-autocomplete>

        <v-text-field
          name="mailAddress"
          label="メールアドレス"
          value="and-period@gmail.com"
          readonly
        ></v-text-field>

        <v-text-field
          name="telephoneNumber"
          label="電話番号"
          value="09012345678"
          readonly
        ></v-text-field>

        <v-textarea name="memo" label="メモ"></v-textarea>
      </v-card-text>
    </v-card>
  </div>
</template>

<script lang="ts">
export default {
  setup() {
    return {
      priority: ['High', 'Middle', 'Low', 'Unknown'],
      status: ['未着手', '進行中', '完了', '不明'],
    }
  },
}
</script>
