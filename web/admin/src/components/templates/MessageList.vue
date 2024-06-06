<script lang="ts" setup>
import { mdiEmailOpenOutline, mdiEmailOutline } from '@mdi/js'
import type { AlertType } from '~/lib/hooks'
import type { Message } from '~/types/api'

interface Props {
  loading: boolean
  isAlert: boolean
  alertType: AlertType
  alertText: string
  message: Message
  messages: Message[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'click:message', messageId: string): void
}>()

const onClickMessage = (messageId: string): void => {
  emit('click:message', messageId)
}

const parsedMessage = computed<string>(() => {
  if (props.message && props.message.body) {
    return props.message.body.replaceAll('\\n', '\n')
  }
  else {
    return ''
  }
})
</script>

<template>
  <v-alert
    v-show="isAlert"
    class="mb-4"
    :type="alertType"
    v-text="alertText"
  />

  <v-row no-gutters>
    <v-col cols="2">
      <v-card class="elevation-1 flex-shrink-0 mr-3">
        <v-list-item>
          <v-list-item-title>メッセージ一覧</v-list-item-title>
        </v-list-item>
        <v-divider />
        <v-list
          nav
          class="pa-2"
        >
          <v-list-item
            v-for="item in messages"
            :key="item.id"
            :prepend-icon="item.read ? mdiEmailOpenOutline : mdiEmailOutline"
            link-k
            @click="onClickMessage(item.id)"
          >
            <v-list-item-title>{{ item.title }}</v-list-item-title>
          </v-list-item>
          <v-list-item v-if="messages.length === 0">
            <v-list-item-title>メッセージなし</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-card>
    </v-col>

    <v-col cols="10">
      <v-card class="elevation-1 d-flex flex-grow-1 flex-column">
        <div v-if="message.title">
          <v-card-title>
            {{ `件名：${message.title}` }}
          </v-card-title>
          <v-divider />
          <v-card-text
            class="message-area"
            v-text="parsedMessage"
          />
        </div>
        <v-card-text
          v-else
          class="text-center"
        >
          メッセージを選択してください
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>

<style scoped>
.message-area {
  white-space: pre-wrap;
}
</style>
