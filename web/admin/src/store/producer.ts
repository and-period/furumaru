import axios from 'axios'
import { defineStore } from 'pinia'

import ApiClientFactory from '../plugins/factory'

import { useAuthStore } from './auth'
import { useCommonStore } from './common'

import {
  CreateProducerRequest,
  ProducerApi,
  ProducerResponse,
  ProducersResponse,
  UpdateProducerRequest,
  UploadImageResponse
} from '~/types/api'
import {
  AuthError,
  ConflictError,
  ConnectionError,
  InternalServerError,
  NotFoundError,
  ValidationError
} from '~/types/exception'

export const useProducerStore = defineStore('Producer', {
  state: () => {
    const apiClient = (token: string) => {
      const factory = new ApiClientFactory()
      return factory.create(ProducerApi, token)
    }

    return {
      producers: [] as ProducersResponse['producers'],
      totalItems: 0,
      apiClient
    }
  },

  actions: {
    /**
     * 登録済みの生産者一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     */
    async fetchProducers (limit = 20, offset = 0, options = ''): Promise<void> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(
            new AuthError('認証エラー。再度ログインをしてください。')
          )
        }

        const res = await this.apiClient(accessToken).v1ListProducers(
          limit,
          offset,
          options
        )
        this.producers = res.data.producers
        this.totalItems = res.data.total
      } catch (error) {
        if (axios.isAxiosError(error)) {
          if (!error.response) {
            return Promise.reject(new ConnectionError(error))
          }
          switch (error.response.status) {
            case 401:
              return Promise.reject(
                new AuthError('認証エラー。再度ログインをしてください。', error)
              )
            case 500:
            default:
              return Promise.reject(new InternalServerError(error))
          }
        }
        throw new InternalServerError(error)
      }
    },

    /**
     * 生産者を新規登録する非同期関数
     * @param payload
     */
    async createProducer (payload: CreateProducerRequest): Promise<void> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(
            new AuthError('認証エラー。再度ログインをしてください。')
          )
        }

        await this.apiClient(accessToken).v1CreateProducer(payload)
        const commonStore = useCommonStore()
        commonStore.addSnackbar({
          message: `${payload.storeName}を作成しました。`,
          color: 'info'
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
                new ValidationError('入力内容に誤りがあります。', error)
              )
            case 401:
              return Promise.reject(
                new AuthError('認証エラー。再度ログインをしてください。', error)
              )
            case 409:
              return Promise.reject(
                new ConflictError(
                  'このメールアドレスはすでに登録されているため、登録できません。',
                  error
                )
              )
            case 500:
            default:
              return Promise.reject(new InternalServerError(error))
          }
        }

        throw new InternalServerError(error)
      }
    },

    /**
     * 生産者のサムネイル画像をアップロードする関数
     * @param payload サムネイル画像のファイルオブジェクト
     * @returns アップロード後のサムネイル画像のパスを含んだオブジェクト
     */
    async uploadProducerThumbnail (payload: File): Promise<UploadImageResponse> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(
            new AuthError('認証エラー。再度ログインをしてください。')
          )
        }

        const res = await this.apiClient(accessToken).v1UploadProducerThumbnail(
          payload,
          {
            headers: {
              'Content-Type': 'multipart/form-data'
            }
          }
        )
        return res.data
      } catch (error) {
        if (axios.isAxiosError(error)) {
          if (!error.response) {
            return Promise.reject(new ConnectionError(error))
          }
          const statusCode = error.response.status
          switch (statusCode) {
            case 401:
              return Promise.reject(
                new AuthError('認証エラー。再度ログインをしてください', error)
              )
            case 400:
              return Promise.reject(
                new ValidationError(
                  'このファイルはアップロードできません。',
                  error
                )
              )
            case 500:
            default:
              return Promise.reject(new InternalServerError(error))
          }
        }
        throw new InternalServerError(error)
      }
    },

    /**
     * 生産者のヘッダー画像をアップロードする関数
     * @param payload ヘッダー画像のファイルオブジェクト
     * @returns アップロード後のヘッダー画像のパスを含んだオブジェクト
     */
    async uploadProducerHeader (payload: File): Promise<UploadImageResponse> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(
            new AuthError('認証エラー。再度ログインをしてください。')
          )
        }

        const res = await this.apiClient(accessToken).v1UploadProducerHeader(
          payload,
          {
            headers: {
              'Content-Type': 'multipart/form-data'
            }
          }
        )
        return res.data
      } catch (error) {
        if (axios.isAxiosError(error)) {
          if (!error.response) {
            return Promise.reject(new ConnectionError(error))
          }
          const statusCode = error.response.status
          switch (statusCode) {
            case 401:
              return Promise.reject(
                new AuthError('認証エラー。再度ログインをしてください', error)
              )
            case 400:
              return Promise.reject(
                new ValidationError(
                  'このファイルはアップロードできません。',
                  error
                )
              )
            case 500:
            default:
              return Promise.reject(new InternalServerError(error))
          }
        }
        throw new InternalServerError(error)
      }
    },

    /**
     * 生産者IDから生産者の情報を取得する非同期関数
     * @param id 生産者ID
     * @returns 生産者の情報
     */
    async getProducer (id: string): Promise<ProducerResponse> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(
            new AuthError('認証エラー。再度ログインをしてください')
          )
        }

        const res = await this.apiClient(accessToken).v1GetProducer(id)
        return res.data
      } catch (error) {
        if (axios.isAxiosError(error)) {
          if (!error.response) {
            return Promise.reject(new ConnectionError(error))
          }
          const statusCode = error.response.status
          switch (statusCode) {
            case 401:
              return Promise.reject(
                new AuthError('認証エラー。再度ログインをしてください', error)
              )
            case 404:
              return Promise.reject(
                new NotFoundError(
                  '一致する生産者が見つかりませんでした。',
                  error
                )
              )
            case 500:
            default:
              return Promise.reject(new InternalServerError(error))
          }
        }
        throw new InternalServerError(error)
      }
    },

    /**
     * 生産者を更新する非同期関数
     * @param id 更新対象の生産者ID
     * @param payload
     * @returns
     */
    async updateProducer (id: string, payload: UpdateProducerRequest) {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(
            new AuthError('認証エラー。再度ログインをしてください')
          )
        }

        await this.apiClient(accessToken).v1UpdateProducer(id, payload)
      } catch (error) {
        if (axios.isAxiosError(error)) {
          if (!error.response) {
            return Promise.reject(new ConnectionError(error))
          }
          const statusCode = error.response.status
          switch (statusCode) {
            case 400:
              return Promise.reject(
                new ValidationError('入力内容に誤りがあります。', error)
              )
            case 401:
              return Promise.reject(
                new AuthError('認証エラー。再度ログインをしてください', error)
              )
            case 404:
              return Promise.reject(
                new NotFoundError(
                  '一致する生産者が見つかりませんでした。',
                  error
                )
              )
            case 500:
            default:
              return Promise.reject(new InternalServerError(error))
          }
        }
        throw new InternalServerError(error)
      }
    },

    /**
     * 生産者を削除する非同期関数
     * @param id 削除する生産者のID
     * @returns
     */
    async deleteProducer (id: string) {
      try {
        const { accessToken } = useAuthStore()
        if (!accessToken) {
          return Promise.reject(
            new AuthError('認証エラー。再度ログインをしてください。')
          )
        }
        await this.apiClient(accessToken).v1DeleteProducer(id)
      } catch (error) {
        return this.errorHandler(error)
      }
    }
  }
})
