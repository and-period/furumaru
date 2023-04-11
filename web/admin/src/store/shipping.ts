import axios from 'axios'
import { defineStore } from 'pinia'
import { apiClient } from '~/plugins/api-client'

import {
  CreateShippingRequest,
  ShippingResponse,
  ShippingsResponseShippingsInner,
  UpdateShippingRequest
} from '~/types/api'
import {
  AuthError,
  ConnectionError,
  InternalServerError,
  NotFoundError,
  ValidationError
} from '~/types/exception'

export const useShippingStore = defineStore('shipping', {
  state: () => ({
    shippings: [] as ShippingsResponseShippingsInner[],
    totalItems: 0
  }),

  actions: {
    /**
     * 配送情報一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     * @returns
     */
    async fetchShippings (limit = 20, offset = 0): Promise<void> {
      try {
        const res = await apiClient.shippingApi().v1ListShippings(
          limit,
          offset
        )
        this.shippings = res.data.shippings
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
     * 指定したIDの配送設定情報を取得する非同期関数
     * @param id 配送設定情報ID
     * @returns 配送設定情報
     */
    async getShipping (id: string): Promise<ShippingResponse> {
      try {
        const res = await apiClient.shippingApi().v1GetShipping(id)
        return res.data
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
     * 配送情報を新規作成する非同期関数
     * @param payload
     * @returns
     */
    async createShipping (payload: CreateShippingRequest): Promise<void> {
      try {
        await apiClient.shippingApi().v1CreateShipping(payload)
      } catch (error) {
        if (axios.isAxiosError(error)) {
          if (!error.response) {
            return Promise.reject(new ConnectionError(error))
          }
          switch (error.response.status) {
            case 400:
              return Promise.reject(
                new ValidationError('入力内容に誤りがあります。', error)
              )
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
     * 指定したIDの配送情報を更新する非同期関数
     * @param id 配送情報ID
     * @param payload
     * @returns
     */
    async updateShipping (
      id: string,
      payload: UpdateShippingRequest
    ): Promise<void> {
      try {
        await apiClient.shippingApi().v1UpdateShipping(id, payload)
      } catch (error) {
        if (axios.isAxiosError(error)) {
          if (!error.response) {
            return Promise.reject(new ConnectionError(error))
          }
          switch (error.response.status) {
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
                new NotFoundError('この配送情報は存在しません。', error)
              )
            case 500:
            default:
              return Promise.reject(new InternalServerError(error))
          }
        }
        throw new InternalServerError(error)
      }
    }
  }
})
