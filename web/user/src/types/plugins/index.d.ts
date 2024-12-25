import 'pinia'
import type {
  AddressApi,
  AuthApi,
  AuthUserApi,
  CartApi,
  CheckoutApi,
  ProductApi,
  TopApi,
  ScheduleApi,
  CoordinatorApi,
  StatusApi,
  OrderApi,
  PromotionApi,
  ExperienceApi,
  SpotApi,
} from '../api'
import type { CustomErrorMessage } from '~/plugins/api-error-handler'

declare module 'pinia' {
  export interface PiniaCustomProperties {
    /**
     * apiクライアントのエラーをハンドリングする関数
     * @param error 発生したエラーオブジェクト
     * @param customObject エラーメッセージをカスタマイズするオブジェクト
     * @returns Promise.rejectを返す。呼び出す側で再度returnすることを想定している。
     */
    errorHandler: (
      error: unknown,
      customObject?: CustomErrorMessage,
    ) => Promise<never>
    authApiClient: (token?: string | undefined) => AuthApi
    authUserApiClient: (token?: string | undefined) => AuthUserApi
    topPageApiClient: (token?: string | undefined) => TopApi
    productApiClient: (token?: string | undefined) => ProductApi
    cartApiClient: (token?: string) => CartApi
    addressApiClient: (token?: string) => AddressApi
    checkoutApiClient: (token?: string) => CheckoutApi
    scheduleApiClient: (token?: string) => ScheduleApi
    coordinatorApiClient: (token?: string) => CoordinatorApi
    statusApiClient: () => StatusApi
    orderApiClient: (token?: string) => OrderApi
    promotionApiClient: (token?: string) => PromotionApi
    otherApiClient: (token?: string) => OtherApi
    spotApiClient: (token?: string) => SpotApi
    experienceApiClient: (token?: string) => ExperienceApi
    i18n: VueI18
  }
}

declare module '#app' {
  // eslint-disable-next-line no-unused-vars
  interface NuxtApp {
    $md: MarkdownIt
  }
}
