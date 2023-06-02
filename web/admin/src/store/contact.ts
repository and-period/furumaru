import { defineStore } from 'pinia'

import { useCommonStore } from './common'
import {
  ContactResponse,
  ContactsResponse,
  ContactsResponseContactsInner,
  UpdateContactRequest
} from '~/types/api'
import { apiClient } from '~/plugins/api-client'

export const useContactStore = defineStore('contact', {
  state: () => ({
    contact: {} as ContactResponse,
    contacts: [] as Array<ContactsResponseContactsInner>,
    total: 0
  }),

  actions: {
    /**
     * お問い合わせの一覧を取得する非同期関数
     * @param limit 最大取得件数
     * @param offset 取得開始位置
     * @param orders ソートキー
     * @returns
     */
    async fetchContacts (
      limit = 20,
      offset = 0,
      orders: string[] = []
    ): Promise<void> {
      try {
        const res = await apiClient.contactApi().v1ListContacts(
          limit,
          offset,
          orders.join(',')
        )
        const { contacts, total }: ContactsResponse = res.data

        this.contacts = contacts
        this.total = total
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    async getContact (id: string): Promise<ContactResponse> {
      try {
        const res = await apiClient.contactApi().v1GetContact(id)
        this.contact = res.data
        return res.data
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    async contactUpdate (
      payload: UpdateContactRequest,
      contactId: string
    ): Promise<void> {
      try {
        await apiClient.contactApi().v1UpdateContact(contactId, payload)
        const commonStore = useCommonStore()
        commonStore.addSnackbar({
          message: 'お問い合わせ情報が更新されました。',
          color: 'info'
        })
      } catch (err) {
        return this.errorHandler(err)
      }
    }
  }
})
