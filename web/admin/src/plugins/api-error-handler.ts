import type { PiniaPluginContext } from 'pinia'
import { ResponseError } from '~/types/api/v1'

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
  const errorHandler = (err: unknown, customObject?: CustomErrorMessage) => {
    if (err instanceof ResponseError) {
      if (!err.response) {
        return Promise.reject(
          new AuthError('認証エラー。再度ログインをしてください。'),
        )
      }

      const statusCode = err.response.status

      let customMessage: string | undefined
      if (customObject && statusCode in customObject) {
        customMessage = customObject[statusCode as keyof CustomErrorMessage]
      }

      switch (statusCode) {
        case 400:
          return Promise.reject(
            new ValidationError(
              customMessage || '入力内容に誤りがあります。',
              err,
            ),
          )
        case 401:
          return Promise.reject(
            new AuthError(
              customMessage || '認証エラー。再度ログインをしてください。',
              err,
            ),
          )
        case 403:
          return Promise.reject(
            new PermissionError(
              customMessage || 'この操作を実施する権限がありません。',
              err,
            ),
          )
        case 404:
          return Promise.reject(
            new NotFoundError(
              customMessage || '指定したリソースが見つかりません。',
              err,
            ),
          )
        case 409:
          return Promise.reject(
            new ConflictError(
              customMessage || '指定したリソースは競合しています。',
              err,
            ),
          )
        case 412:
          return Promise.reject(
            new PreconditionError(
              customMessage || '指定したリソースは条件を満たしていません。',
              err,
            ),
          )
        case 429:
          return Promise.reject(
            new TooManyRequestsError(
              customMessage || 'リクエスト数が上限を超えています。',
              err,
            ),
          )
        case 499:
          return Promise.reject(
            new CancelledError(
              customMessage || '正常に処理を完了できませんでした。',
              err,
            ),
          )
        case 501:
          return Promise.reject(new NotImplementedError(err))
        case 503:
          return Promise.reject(new ServiceUnavailableError(err))
        case 500:
        default:
          return Promise.reject(new InternalServerError(err))
      }
    }
    return Promise.reject(new InternalServerError(err))
  }

  store.errorHandler = markRaw(errorHandler)
}

/**
 * piniaに共通エラーハンドラーを注入するプラグイン
 */
export default defineNuxtPlugin((nuxtApp) => {
  const pinia = nuxtApp.$pinia as any
  pinia.use(apiErrorHandler)
})
