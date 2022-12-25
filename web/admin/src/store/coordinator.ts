import axios from 'axios'
import { defineStore } from 'pinia'

import { useAuthStore } from './auth'
import { useCommonStore } from './common'

import ApiClientFactory from '~/plugins/factory'
import {
  CoordinatorApi,
  CoordinatorResponse,
  CoordinatorsResponse,
  CreateCoordinatorRequest,
  RelateProducersRequest,
  UpdateCoordinatorRequest,
  UploadImageResponse,
} from '~/types/api'
import {
  AuthError,
  ConnectionError,
  InternalServerError,
  NotFoundError,
  PreconditionError,
  ValidationError,
} from '~/types/exception'

export const useCoordinatorStore = defineStore('Coordinator', {
  state: () => {
    const apiClient = (token: string) => {
      const factory = new ApiClientFactory()
      return factory.create(CoordinatorApi, token)
    }
    return {
      coordinators: [] as CoordinatorsResponse['coordinators'],
      totalItems: 0,
      apiClient,
    }
  },
  actions: {
    /**
     * コーディネータの一覧を取得する非同期関数
     * @param limit 最大取得件数
     * @param offset 取得開始位置
     * @returns
     */
    async fetchCoordinators(
      limit: number = 20,
      offset: number = 0
    ): Promise<void> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(
            new AuthError('認証エラー。再度ログインをしてください。')
          )
        }

        const factory = new ApiClientFactory()
        const coordinatorsApiClient = factory.create(
          CoordinatorApi,
          accessToken
        )
        const res = await coordinatorsApiClient.v1ListCoordinators(
          limit,
          offset
        )
        this.coordinators = res.data.coordinators
        this.totalItems = res.data.total
      } catch (error) {
        return this.errorHandler(error)
      }
    },

    /**
     * コーディネータを登録する非同期関数
     * @param payload
     * @returns
     */
    async createCoordinator(payload: CreateCoordinatorRequest) {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(
            new AuthError('認証エラー。再度ログインをしてください。')
          )
        }
        const factory = new ApiClientFactory()
        const coordinatorsApiClient = factory.create(
          CoordinatorApi,
          accessToken
        )
        const res = await coordinatorsApiClient.v1CreateCoordinator(payload)
        return res.data
      } catch (error) {
        console.log(error)
        if (axios.isAxiosError(error)) {
          return this.errorHandler(error)
        }
      }
    },

    /**
     * コーディネータの詳細情報を取得する非同期関数
     * @param id 対象のコーディネータのID
     * @returns
     */
    async getCoordinator(id: string): Promise<CoordinatorResponse> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(
            new AuthError('認証エラー。再度ログインをしてください')
          )
        }

        const factory = new ApiClientFactory()
        const coordinatorsApiClient = factory.create(
          CoordinatorApi,
          accessToken
        )
        const res = await coordinatorsApiClient.v1GetCoordinator(id)
        return res.data
      } catch (error) {
        return this.errorHandler(error, {
          404: '該当するコーディーネータが見つかりませんでした。',
        })
      }
    },

    /**
     * コーディネータの情報を更新する非同期関数
     * @param payload
     * @param coordinatorId 更新するコーディネータのID
     * @returns
     */
    async updateCoordinator(
      payload: UpdateCoordinatorRequest,
      coordinatorId: string
    ): Promise<void> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(new Error('認証エラー'))
        }
        const factory = new ApiClientFactory()
        const contactsApiClient = factory.create(CoordinatorApi, accessToken)
        await contactsApiClient.v1UpdateCoordinator(coordinatorId, payload)
        const commonStore = useCommonStore()
        commonStore.addSnackbar({
          message: 'コーディネータ情報が更新されました。',
          color: 'info',
        })
      } catch (error) {
        return this.errorHandler(error, {
          404: '該当するコーディーネータが見つかりませんでした。',
        })
      }
    },

    /**
     * コーディネータのサムネイル画像をアップロードする非同期関数
     * @param payload サムネイル画像
     * @returns アップロードされた画像のURI
     */
    async uploadCoordinatorThumbnail(
      payload: File
    ): Promise<UploadImageResponse> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(
            new AuthError('認証エラー。再度ログインをしてください。')
          )
        }

        const factory = new ApiClientFactory()
        const coordinatorsApiClient = factory.create(
          CoordinatorApi,
          accessToken
        )
        const res = await coordinatorsApiClient.v1UploadCoordinatorThumbnail(
          payload,
          {
            headers: {
              'Content-Type': 'multipart/form-data',
            },
          }
        )
        return res.data
      } catch (error) {
        return this.errorHandler(error, {
          400: 'このファイルはアップロードできません。',
        })
      }
    },

    /**
     * コーディネータのヘッダー画像をアップロードする非同期関数
     * @param payload ヘッダー画像
     * @returns アップロードされた画像のURI
     */
    async uploadCoordinatorHeader(payload: File): Promise<UploadImageResponse> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(
            new AuthError('認証エラー。再度ログインをしてください。')
          )
        }

        const factory = new ApiClientFactory()
        const coordinatorsApiClient = factory.create(
          CoordinatorApi,
          accessToken
        )
        const res = await coordinatorsApiClient.v1UploadCoordinatorHeader(
          payload,
          {
            headers: {
              'Content-Type': 'multipart/form-data',
            },
          }
        )
        return res.data
      } catch (error) {
        return this.errorHandler(error, {
          400: 'このファイルはアップロードできません。',
        })
      }
    },

    /**
     * コーディーネータを削除する非同期関数
     * @param id 削除するコーディネータのID
     * @returns
     */
    async deleteCoordinator(id: string) {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(new Error('認証エラー'))
        }

        const factory = new ApiClientFactory()
        const coordinatorsApiClient = factory.create(
          CoordinatorApi,
          accessToken
        )
        await coordinatorsApiClient.v1DeleteCoordinator(id)
        const commonStore = useCommonStore()
        commonStore.addSnackbar({
          message: 'コーディネーターの削除が完了しました',
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
                new ValidationError(
                  '削除できませんでした。管理者にお問い合わせしてください。',
                  error
                )
              )
            case 401:
              return Promise.reject(
                new AuthError('認証エラー。再度ログインをしてください。', error)
              )
            case 404:
              return Promise.reject(
                new NotFoundError(
                  '削除するコーディネーターが見つかりませんでした。',
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
      this.fetchCoordinators()
    },

    /**
     * コーディーネータに生産者を紐づける非同期関数
     * @param id 生産者を紐づけるコーディネータのID
     * @param payload コーディネーターに紐づく生産者
     * @returns
     */
    async relateProducers(
      id: string,
      payload: RelateProducersRequest
    ): Promise<void> {
      try {
        const authStore = useAuthStore()
        const accessToken = authStore.accessToken
        if (!accessToken) {
          return Promise.reject(new Error('認証エラー'))
        }
        const factory = new ApiClientFactory()
        const contactsApiClient = factory.create(CoordinatorApi, accessToken)
        await contactsApiClient.v1RelateProducers(id, payload)
        const commonStore = useCommonStore()
        commonStore.addSnackbar({
          message: 'コーディネーターと生産者の紐付けが完了しました',
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
                new ValidationError('入力内容に誤りがあります。', error)
              )
            case 401:
              return Promise.reject(
                new AuthError('認証エラー。再度ログインをしてください。', error)
              )
            case 404:
              return Promise.reject(
                new NotFoundError(
                  'コーディネーターが見つかりませんでした。',
                  error
                )
              )
            case 412:
              return Promise.reject(
                new PreconditionError(
                  '既に関連づけられている生産者または、存在しない生産者が指定されています。',
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
  },
})
