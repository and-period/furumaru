import axios from 'axios'
import type { PiniaPluginContext } from 'pinia'

import {
  AuthError,
  CancelledError,
  ConflictError,
  InternalServerError,
  NotFoundError,
  NotImplementedError,
  PermissionError,
  PreconditionError,
  ServiceUnavailableError,
  TooManyRequestsError,
  ValidationError,
} from '~/types/exception'

// メッセージを変更できるステータスコードの一覧（500系のエラーはシステム固有のメッセージを利用する）
const STATUS_CODES = [400, 401, 403, 404, 409, 412, 429, 499] as const

type StatusCode = typeof STATUS_CODES[number]

export type CustomErrorMessage = {
  [key in StatusCode]?: string
}

/**
 * APIのエラーハンドリングを共通化するpiniaのプラグイン
 */
function apiErrorHandler({ store }: PiniaPluginContext) {
  /**
   * apiクライアントのエラーをハンドリングする関数
   * @param error 発生したエラーオブジェクト
   * @param customObject エラーメッセージをカスタマイズするオブジェクト
   * @returns Promise.rejectを返す。呼び出す側で再度returnすることを想定している。
   */
  const errorHandler = (error: unknown, customObject?: CustomErrorMessage) => {
    if (axios.isAxiosError(error)) {
      if (!error.response) {
        return Promise.reject(
          new AuthError('認証エラー。再度ログインをしてください。'),
        )
      }

      const statusCode = error.response.status

      let customMessage: string | undefined
      if (customObject && statusCode in customObject) {
        customMessage = customObject[statusCode as keyof CustomErrorMessage]
      }

      switch (statusCode) {
        case 400:
          return Promise.reject(
            new ValidationError(
              customMessage || '入力内容に誤りがあります。',
              error,
            ),
          )
        case 401:
          return Promise.reject(
            new AuthError(
              customMessage || '認証エラー。再度ログインをしてください。',
              error,
            ),
          )
        case 403:
          return Promise.reject(
            new PermissionError(
              customMessage || 'この操作を実施する権限がありません。',
              error,
            ),
          )
        case 404:
          return Promise.reject(
            new NotFoundError(
              customMessage || '指定したリソースが見つかりません。',
              error,
            ),
          )
        case 409:
          return Promise.reject(
            new ConflictError(
              customMessage || '指定したリソースは競合しています。',
              error,
            ),
          )
        case 412:
          return Promise.reject(
            new PreconditionError(
              customMessage || '指定したリソースは条件を満たしていません。',
              error,
            ),
          )
        case 429:
          return Promise.reject(
            new TooManyRequestsError(
              customMessage || 'リクエスト数が上限を超えています。',
              error,
            ),
          )
        case 499:
          return Promise.reject(
            new CancelledError(
              customMessage || '正常に処理を完了できませんでした。',
              error,
            ),
          )
        case 501:
          return Promise.reject(new NotImplementedError(error))
        case 503:
          return Promise.reject(new ServiceUnavailableError(error))
        case 500:
        default:
          return Promise.reject(new InternalServerError(error))
      }
    }
    return Promise.reject(new InternalServerError(error))
  }

  store.errorHandler = markRaw(errorHandler)
}

/**
 * piniaに共通エラーハンドラーを注入するプラグイン
 */
export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.$pinia.use(apiErrorHandler)
})
