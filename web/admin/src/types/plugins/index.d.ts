import 'pinia'

import { CustomErrorMessage } from '~/plugins/api-error-handler'
import { useAuthStore } from '~/store'
import { AdministratorApi, AuthApi, CategoryApi, ContactApi, CoordinatorApi, MessageApi, NotificationApi, OrderApi, ProducerApi, ProductApi, ProductTypeApi, PromotionApi, ShippingApi } from '~/types/api'

declare module 'vue/types/vue' {
  interface Vue {
    $auth: ReturnType<typeof useAuthStore>
  }
}

declare module '@nuxt/types' {
  interface Context {
    $auth: ReturnType<typeof useAuthStore>
  }

  interface NuxtAppOptions {
    $auth: ReturnType<typeof useAuthStore>
  }
}

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
      customObject?: CustomErrorMessage
    ) => Promise<never>

    /**
     * apiクライアントの定義
     * @param token APIリクエスト時の認証用トークン
     * @returns OpenApi Generatorによって自動生成されたAPIクライアントを返す。
     */
    administratorApiClient: (token?: string) => AdministratorApi
    authApiClient: (token?: string) => AuthApi
    categoryApiClient: (token?: string) => CategoryApi
    contactApiClient: (token?: string) => ContactApi
    coordinatorApiClient: (token?: string) => CoordinatorApi
    messageApiClient: (token?: string) => MessageApi
    notificationApiClient: (token?: string) => NotificationApi
    orderApiClient: (token?: string) => OrderApi
    producerApiClient: (token?: string) => ProducerApi
    productApiClient: (token?: string) => ProductApi
    productTypeApiClient: (token?: string) => ProductTypeApi
    promotionApiClient: (token?: string) => PromotionApi
    shippingApiClient: (token?: string) => ShippingApi
  }
}
