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
      <v-card-actions>
        <v-btn block outlined color="primary">登録</v-btn>
      </v-card-actions>
    </v-card>
  </div>
</template>

<script lang="ts">
import { defineComponent, useFetch, useRoute } from '@nuxtjs/composition-api'
import { useContactStore } from '~/store/contact'

export default defineComponent ({
  setup() {
    const route = useRoute()
    const id = route.value.params.id
    const { getContact } = useContactStore()
    var title = ''
    var content = ''
    var username = ''
    var email = ''
    var phoneNumber = ''
    var _status = 0
    var _priority = 0
    var note = ''

const {fetchState} = useFetch(async () => {
    const contact = await getContact(id)
    title = contact.title
    content = contact.content
    username = contact.username
    email = contact.email
    phoneNumber = contact.phoneNumber
    _status = contact.status
    _priority = contact.priority
    note = contact.note
})

    return {
      priority: ['High', 'Middle', 'Low', 'Unknown'],
      status: ['未着手', '進行中', '完了', '不明'],
      id,
      fetchState,
      title,
      content,
      username,
      email,
      phoneNumber,
      _status,
      _priority,
      note,
    }
  },
})
</script>
