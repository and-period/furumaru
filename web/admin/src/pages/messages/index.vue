<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useAlert } from '~/lib/hooks'
import { useMessageStore } from '~/store'

const messageStore = useMessageStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const { message, messages } = storeToRefs(messageStore)

const loading = ref<boolean>(false)

const handleClickMessage = async (messageId: string) => {
  try {
    loading.value = true
    await messageStore.fetchMessage(messageId)
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    loading.value = false
  }
}

// 初期化処理
const fetchState = useAsyncData(async (): Promise<void> => {
  try {
    await messageStore.fetchMessages(50, 0) // 最大50件のメッセージを取得
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
})

try {
  await fetchState.execute()
}
catch (err) {
  console.log('failed to setup message list', err)
}
</script>

<template>
  <templates-message-list
    :loading="loading"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :message="message"
    :messages="messages"
    @click:message="handleClickMessage"
  />
</template>
