<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useAlert } from '~/lib/hooks'
import { useMessageStore } from '~/store'

const messageStore = useMessageStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const { message, messages } = storeToRefs(messageStore)

const handleClickMessage = async (messageId: string) => {
  try {
    await messageStore.fetchMessage(messageId)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}
</script>

<template>
  <templates-message-list
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :message="message"
    :messages="messages"
    @click:message="handleClickMessage"
  />
</template>
