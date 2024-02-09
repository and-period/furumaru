import 'pinia'
import {
  AddressApi,
  AuthApi,
  CartApi,
  CheckoutApi,
  ProductApi,
  TopApi,
  ScheduleApi,
  CoordinatorApi,
  StatusApi,
  OrderApi,
  PromotionApi,
} from '../api'
import { CustomErrorMessage } from '~/plugins/api-error-handler'
import VueI18n from 'vue-i18n'

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
    i18n: VueI18
  }
}

declare module '#app' {
  interface NuxtApp {
    $md: MarkdownIt
  }
}
