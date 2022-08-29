import axios from 'axios'
import { defineStore } from 'pinia'

import { useAuthStore } from './auth'

import ApiClientFactory from '~/plugins/factory'
import { CoordinatorApi, CreateCoordinatorRequest } from '~/types/api'
import {
  AuthError,
  ConnectionError,
  InternalServerError,
  ValidationError,
} from '~/types/exception'

export const useCoordinatorStore = defineStore('Coordinator', {
  state: () => {
    const apiClient = (token: string) => {
      const factory = new ApiClientFactory()
      return factory.create(CoordinatorApi, token)
    }
    return {
      apiClient,
      coordinators: [],
    }
  },

  actions: {
    /**
     * 仲介者を登録する非同期関数
     * @param payload
     * @returns
     */
    async createCoordinator(payload: CreateCoordinatorRequest) {
      try {
        const { accessToken } = useAuthStore()
        if (!accessToken) {
          return Promise.reject(
            new AuthError('認証エラー。再度ログインをしてください。')
          )
        }
        await this.apiClient(accessToken).v1CreateCoordinator(payload)
      } catch (error) {
        console.log(error)
        if (axios.isAxiosError(error)) {
          if (!error.response) {
            return Promise.reject(new ConnectionError(error))
          }
          switch (error.response.status) {
            case 401:
              return Promise.reject(
                new AuthError('認証エラー。再度ログインをしてください。', error)
              )
            case 400:
              return Promise.reject(
                new ValidationError('入力値に誤りがあります。', error)
              )
            case 500:
            default:
              return Promise.reject(new InternalServerError(error))
          }
        }
        throw new InternalServerError(error)
      }
    },
  },
})
