<script lang="ts" setup>
import { useMessageStore } from '~/store/message'

const messageStore = useMessageStore()

const message = computed(() => {
  return messageStore.message
})
const messages = computed(() => {
  return messageStore.messages
})

const handleClickMessage = async (messageId: string) => {
  await messageStore.fetchMessage(messageId)
}
</script>

<template>
  <div class="d-flex flex-row mt-2">
    <v-card class="elevation-1 flex-shrink-0 mr-3">
      <v-list-item>
        <v-list-item-content>
          <v-list-item-title>メッセージ一覧</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
      <v-divider />
      <v-list nav class="pa-2">
        <v-list-item
          v-for="message in messages"
          :key="message.id"
          link
          @click="handleClickMessage(message.id)"
        >
          <v-list-item-icon>
            <v-icon>{{
              message.read ? 'mdi-email-open-outline' : 'mdi-email-outline'
            }}</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>{{ message.title }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item v-if="messages.length === 0">
          <v-list-item-content>
            <v-list-item-title>メッセージなし</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-card>
    <v-card class="elevation-1 d-flex flex-grow-1 flex-column">
      <v-card-title>{{
        message.title
          ? `件名：${message.title}`
          : 'メッセージを選択してください'
      }}</v-card-title>
      <v-divider />
      <v-card-text>{{ message.body }}</v-card-text>
    </v-card>
  </div>
</template>
