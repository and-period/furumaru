import axios from 'axios'
import { defineStore } from 'pinia'

import { useAuthStore } from './auth'
import { useCommonStore } from './common'

import ApiClientFactory from '~/plugins/factory'
import {
  ContactApi,
  ContactResponse,
  ContactsResponse,
  ContactsResponseContactsInner,
  UpdateContactRequest,
} from '~/types/api'
import {
  AuthError,
  ConnectionError,
  NotFoundError,
  ValidationError,
} from '~/types/exception'

export const useContactStore = defineStore('contact', {
  state: () => {
    const apiClient = (token: string) => {
      const factory = new ApiClientFactory()
      return factory.create(ContactApi, token)
    }

    return {
      apiClient,
      contacts: [] as Array<ContactsResponseContactsInner>,
      total: 0,
    }
  },

  actions: {
    /**
     * お問い合わせの一覧を取得する非同期関数
     * @param limit 最大取得件数
     * @param offset 取得開始位置
     * @param orders ソートキー
     * @returns
     */
    async fetchContacts(limit: number = 20, offset: number = 0, orders: string[] = []): Promise<void> {
      try {
        const accessToken = this.getAccessToken()
        const res = await this.apiClient(accessToken).v1ListContacts(
          limit,
          offset,
          orders.join(','),
        )
        const { contacts, total }: ContactsResponse = res.data

        this.contacts = contacts
        this.total = total
      } catch (err) {
        return this.errorHandler(err)
      }
    },

    async getContact(id: string): Promise<ContactResponse> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) throw new Error('認証エラー')

        const factory = new ApiClientFactory()
        const contactsApiClient = factory.create(ContactApi, accessToken)
        const res = await contactsApiClient.v1GetContact(id)
        return res.data
      } catch (error) {
        if (axios.isAxiosError(error)) {
          if (!error.response) {
            return Promise.reject(new ConnectionError(error))
          }
          const statusCode = error.response.status
          switch (statusCode) {
            case 404:
              return Promise.reject(
                new NotFoundError(
                  '編集するお問い合わせが見つかりませんでした。',
                  error
                )
              )
          }
        }
        throw new Error('Internal Server Error')
      }
    },

    async contactUpdate(
      payload: UpdateContactRequest,
      contactId: string
    ): Promise<void> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(new Error('認証エラー'))
        }
        const factory = new ApiClientFactory()
        const contactsApiClient = factory.create(ContactApi, accessToken)
        await contactsApiClient.v1UpdateContact(contactId, payload)
        const commonStore = useCommonStore()
        commonStore.addSnackbar({
          message: 'お問い合わせ情報が更新されました。',
          color: 'info',
        })
      } catch (error) {
        if (axios.isAxiosError(error)) {
          if (!error.response) {
            return Promise.reject(new ConnectionError(error))
          }
          const statusCode = error.response.status
          switch (statusCode) {
            case 400:
              return Promise.reject(
                new ValidationError('入力された内容では更新できません。', error)
              )
            case 404:
              return Promise.reject(
                new NotFoundError(
                  '更新するお問い合わせが見つかりませんでした。',
                  error
                )
              )
          }
        }
        throw new Error('Internal Server Error')
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
