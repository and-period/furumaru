import { Configuration, ResponseError } from '~/types/api/v1'
import type { BaseAPI, Middleware, RequestContext } from '~/types/api/v1/runtime'
import { useAuthStore } from '~/store'
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

const STATUS_CODES = [400, 401, 403, 404, 409, 412, 429, 499] as const
type StatusCode = typeof STATUS_CODES[number]
export type CustomErrorMessage = {
  [key in StatusCode]?: string
}

const authMiddleware: Middleware = {
  async pre(ctx: RequestContext) {
    const store = useAuthStore()
    const headers = new Headers(ctx.init?.headers || {})
    const token: string | undefined = store.accessToken
    if (token) {
      headers.set('Authorization', `Bearer ${token}`)
    }
    return {
      url: ctx.url,
      init: { ...ctx.init, headers },
    }
  },
}

function createApiClient<T extends BaseAPI>(Client: new (config: Configuration) => T): T {
  const runtimeConfig = useRuntimeConfig()
  const baseUrl = runtimeConfig.public.API_BASE_URL
  const config = new Configuration({
    basePath: baseUrl,
    credentials: 'include',
    middleware: [authMiddleware],
  })
  return new Client(config)
}

function errorHandler(err: unknown, customObject?: CustomErrorMessage) {
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
          new ValidationError(customMessage || '入力内容に誤りがあります。', err),
        )
      case 401:
        return Promise.reject(
          new AuthError(customMessage || '認証エラー。再度ログインをしてください。', err),
        )
      case 403:
        return Promise.reject(
          new PermissionError(customMessage || 'この操作を実施する権限がありません。', err),
        )
      case 404:
        return Promise.reject(
          new NotFoundError(customMessage || '指定したリソースが見つかりません。', err),
        )
      case 409:
        return Promise.reject(
          new ConflictError(customMessage || '指定したリソースは競合しています。', err),
        )
      case 412:
        return Promise.reject(
          new PreconditionError(customMessage || '指定したリソースは条件を満たしていません。', err),
        )
      case 429:
        return Promise.reject(
          new TooManyRequestsError(customMessage || 'リクエスト数が上限を超えています。', err),
        )
      case 499:
        return Promise.reject(
          new CancelledError(customMessage || '正常に処理を完了できませんでした。', err),
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

export function useApiClient() {
  return {
    create: <T extends BaseAPI>(Client: new (config: Configuration) => T): T => createApiClient(Client),
    errorHandler,
  }
}
