import type { Message, MessagesResponse, V1MessagesGetRequest, V1MessagesMessageIdGetRequest } from '~/types/api/v1'

export const useMessageStore = defineStore('message', {
  state: () => ({
    message: {} as Message,
    messages: [] as Message[],
    total: 0,
    hasUnread: false,
  }),

  actions: {
    /**
     * メッセージの一覧を取得する非同期関数
     * @param limit 最大取得件数
     * @param offset 取得開始位置
     * @param orders ソートキー
     * @returns
     */
    async fetchMessages(
      limit = 20,
      offset = 0,
      orders: string[] = [],
    ): Promise<void> {
      try {
        if (orders.length === 0) {
          orders = ['-read', '-receivedAt'] // 優先順位: 未読 && 受信日時が新しい
        }
        const params: V1MessagesGetRequest = {
          limit,
          offset,
          orders: orders.join(','),
        }
        const res = await this.messageApi().v1MessagesGet(params)
        const { messages, total }: MessagesResponse = res

        this.messages = messages
        this.total = total
        this.hasUnread = messages.some((message): boolean => !message.read)
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * メッセージの一覧を取得する非同期関数
     * @param messageId メッセージID
     * @returns
     */
    async fetchMessage(messageId = ''): Promise<void> {
      try {
        const params: V1MessagesMessageIdGetRequest = {
          messageId,
        }
        const res = await this.messageApi().v1MessagesMessageIdGet(params)
        const message = res.message || {}

        this.message = message
        this.messages.forEach((v: Message, i: number) => {
          if (!this.messages[i]) {
            return
          }
          if (v.id === message.id) {
            this.messages[i].read = true
          }
        })
      }
      catch (err) {
        return this.errorHandler(err, { 404: '対象のメッセージが存在しません' })
      }
    },
  },
})
