import { defineStore } from 'pinia'

import { useAuthStore } from './auth'

import ApiClientFactory from '~/plugins/factory'
import {
  MessageApi,
  MessageResponse,
  MessagesResponse,
  MessagesResponseMessagesInner,
} from '~/types/api'
import { AuthError } from '~/types/exception'

export const useMessageStore = defineStore('message', {
  state: () => {
    const apiClient = (token: string) => {
      const factory = new ApiClientFactory()
      return factory.create(MessageApi, token)
    }
    return {
      apiClient,
      message: {} as MessageResponse,
      messages: [] as Array<MessagesResponseMessagesInner>,
      total: 0,
      hasUnread: false,
    }
  },

  actions: {
    /**
     * メッセージの一覧を取得する非同期関数
     * @param limit 最大取得件数
     * @param offset 取得開始位置
     * @param orders ソートキー
     * @returns
     */
    async fetchMessages(
      limit: number = 20,
      offset: number = 0,
      orders: string[] = []
    ): Promise<void> {
      try {
        const accessToken = this.getAccessToken()
        if (orders.length === 0) {
          orders = ['-read', '-receivedAt'] // 優先順位: 未読 && 受信日時が新しい
        }
        const res = await this.apiClient(accessToken).v1ListMessages(
          limit,
          offset,
          orders.join(',')
        )
        const { messages, total }: MessagesResponse = res.data

        this.messages = messages
        this.total = total
        this.hasUnread = messages.some((message): boolean => !message.read)
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * メッセージの一覧を取得する非同期関数
     * @param limit 最大取得件数
     * @param offset 取得開始位置
     * @param orders ソートキー
     * @returns
     */
    async fetchMessage(messageId = ''): Promise<void> {
      try {
        const accessToken = this.getAccessToken()
        const res = await this.apiClient(accessToken).v1GetMessage(messageId)
        const message = res.data || {}

        this.message = message
        this.messages.forEach((v: MessagesResponseMessagesInner, i: number) => {
          if (v.id === message.id) this.messages[i].read = true
        })
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    getAccessToken(): string {
      const authStore = useAuthStore()
      const accessToken = authStore.accessToken
      if (!accessToken) {
        throw new AuthError('認証エラー。再度ログインをしてください。')
      }
      return accessToken
    },
  },
})
