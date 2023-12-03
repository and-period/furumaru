import { defineStore } from 'pinia'

import { apiClient } from '~/plugins/api-client'
import type {
  Contact,
  ContactResponse,
  UpdateContactRequest
} from '~/types/api'

export const useContactStore = defineStore('contact', {
  state: () => ({
    contact: {} as Contact,
    contacts: [] as Contact[],
    total: 0
  }),

  actions: {
    /**
     * お問い合わせの一覧を取得する非同期関数
     * @param limit 最大取得件数
     * @param offset 取得開始位置
     * @param orders ソートキー
     */
    async fetchContacts (limit = 20, offset = 0, orders: string[] = []): Promise<void> {
      try {
        const res = await apiClient.contactApi().v1ListContacts(limit, offset)
        this.contacts = res.data.contacts
        this.total = res.data.total
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * お問い合わせの一覧を取得する非同期関数
     * @param contactId お問い合わせID
     */
    async getContact (contactId: string): Promise<ContactResponse> {
      try {
        const res = await apiClient.contactApi().v1GetContact(contactId)
        this.contact = res.data.contact
        return res.data
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    async updateContact (contactId: string, payload: UpdateContactRequest): Promise<void> {
      try {
        await apiClient.contactApi().v1UpdateContact(contactId, payload)
      } catch (err) {
        return this.errorHandler(err)
      }
    }
  }
})
