<script lang="ts" setup>
import { mdiEmailOpenOutline, mdiEmailOutline } from '@mdi/js'
import { AlertType } from '~/lib/hooks'
import { Message } from '~/types/api'

interface Props {
  loading: boolean,
  isAlert: boolean,
  alertType: AlertType
  alertText: string,
  message: Message
  messages: Message[]
}

defineProps<Props>()

const emit = defineEmits<{
  (e: 'click:message', messageId: string): void
}>()

const onClickMessage = (messageId: string): void => {
  emit('click:message', messageId)
}
</script>

<template>
  <v-alert v-show="isAlert" class="mb-4" :type="alertType" v-text="alertText" />

  <v-row no-gutters>
    <v-col cols="2">
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
        <v-card-title>
          {{ message.title ? `件名：${message.title}` : 'メッセージを選択してください' }}
        </v-card-title>
        <v-divider />
        <v-card-text>{{ message.body }}</v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>
