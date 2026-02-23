import { useApiClient } from '~/composables/useApiClient'
import { MessageApi } from '~/types/api/v1'
import type { Message, V1MessagesGetRequest, V1MessagesMessageIdGetRequest } from '~/types/api/v1'

export const useMessageStore = defineStore('message', () => {
  const { create, errorHandler } = useApiClient()
  const messageApi = () => create(MessageApi)

  const message = ref<Message>({} as Message)
  const messages = ref<Message[]>([])
  const total = ref<number>(0)
  const hasUnread = ref<boolean>(false)

  async function fetchMessages(limit = 20, offset = 0, orders: string[] = []): Promise<void> {
    try {
      if (orders.length === 0) {
        orders = ['-read', '-receivedAt']
      }
      const params: V1MessagesGetRequest = {
        limit,
        offset,
        orders: orders.join(','),
      }
      const res = await messageApi().v1MessagesGet(params)
      messages.value = res.messages
      total.value = res.total
      hasUnread.value = res.messages.some((m): boolean => !m.read)
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function fetchMessage(messageId = ''): Promise<void> {
    try {
      const params: V1MessagesMessageIdGetRequest = { messageId }
      const res = await messageApi().v1MessagesMessageIdGet(params)
      const msg = res.message || {}

      message.value = msg
      messages.value.forEach((v: Message, i: number) => {
        if (!messages.value[i]) {
          return
        }
        if (v.id === msg.id) {
          messages.value[i].read = true
        }
      })
    }
    catch (err) {
      return errorHandler(err, { 404: '対象のメッセージが存在しません' })
    }
  }

  return {
    message,
    messages,
    total,
    hasUnread,
    fetchMessages,
    fetchMessage,
  }
})
