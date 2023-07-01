<script lang="ts" setup>
import { mdiEmailOpenOutline, mdiEmailOutline } from '@mdi/js'
import { AlertType } from '~/lib/hooks'
import { MessageResponse, MessagesResponseMessagesInner } from '~/types/api'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  isAlert: {
    type: Boolean,
    default: false
  },
  alertType: {
    type: String as PropType<AlertType>,
    default: undefined
  },
  alertText: {
    type: String,
    default: ''
  },
  message: {
    type: Object as PropType<MessageResponse>,
    default: () => ({})
  },
  messages: {
    type: Array<MessagesResponseMessagesInner>,
    default: () => []
  }
})

const emit = defineEmits<{
  (e: 'click:message', messageId: string): void
}>()

const onClickMessage = (messageId: string): void => {
  emit('click:message', messageId)
}
</script>

<template>
  <v-alert v-show="props.isAlert" :type="props.alertType" v-text="props.alertText" />

  <div class="d-flex flex-row mt-2">
    <v-card :loading="loading" class="elevation-1 flex-shrink-0 mr-3">
      <v-list-item>
        <v-list-item-title>メッセージ一覧</v-list-item-title>
      </v-list-item>
      <v-divider />
      <v-list nav class="pa-2">
        <v-list-item
          v-for="item in messages"
          :key="item.id"
          :prepend-icon="item.read ? mdiEmailOpenOutline : mdiEmailOutline"
          link-k
          @click="onClickMessage(message.id)"
        >
          <v-list-item-title>{{ message.title }}</v-list-item-title>
        </v-list-item>
        <v-list-item v-if="messages.length === 0">
          <v-list-item-title>メッセージなし</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-card>

    <v-card class="elevation-1 d-flex flex-grow-1 flex-column">
      <v-card-title>
        {{ message.title ? `件名：${message.title}` : 'メッセージを選択してください' }}
      </v-card-title>
      <v-divider />
      <v-card-text>{{ message.body }}</v-card-text>
    </v-card>
  </div>
</template>
