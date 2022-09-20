import axios from 'axios'
import { defineStore } from 'pinia'

import { useAuthStore } from './auth'
import { useCommonStore } from './common'

import ApiClientFactory from '~/plugins/factory'
import {
  ContactApi,
  ContactResponse,
  ContactsResponse,
  UpdateContactRequest,
} from '~/types/api'
import {
  ConnectionError,
  NotFoundError,
  ValidationError,
} from '~/types/exception'

export const useContactStore = defineStore('Contact', {
  state: () => ({
    contacts: [] as ContactsResponse['contacts'],
  }),
  actions: {
    async fetchContacts(limit: number = 20, offset: number = 0): Promise<void> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(new Error('認証エラー'))
        }
        const factory = new ApiClientFactory()
        const contactsApiClient = factory.create(ContactApi, accessToken)
        const res = await contactsApiClient.v1ListContacts(limit, offset)
        this.contacts = res.data.contacts
      } catch (error) {
        throw new Error('Internal Server Error')
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
  },
})
