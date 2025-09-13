import 'pinia'

import type { CustomErrorMessage } from '~/plugins/api-error-handler'
import type { AdministratorApi, AuthApi, BroadcastApi, CategoryApi, ContactApi, CoordinatorApi, ExperienceApi, ExperienceTypeApi, GuestApi, LiveApi, MessageApi, NotificationApi, OrderApi, PaymentSystemApi, PostalCodeApi, ProducerApi, ProductApi, ProductTagApi, ProductTypeApi, PromotionApi, ScheduleApi, ShippingApi, ShopApi, SpotTypeApi, TopApi, UploadApi, UserApi, VideoApi } from '../api/v1'

declare module 'vue/types/vue' {
  interface Vue {
  }
}

declare module '@nuxt/types' {
  interface Context {
  }

  interface NuxtAppOptions {
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
    administratorApi: (token?: string | undefined) => AdministratorApi
    authApi: (token?: string | undefined) => AuthApi
    broadcastApi: (token?: string | undefined) => BroadcastApi
    categoryApi: (token?: string | undefined) => CategoryApi
    contactApi: (token?: string | undefined) => ContactApi
    coordinatorApi: (token?: string | undefined) => CoordinatorApi
    experienceApi: (token?: string | undefined) => ExperienceApi
    experienceTypeApi: (token?: string | undefined) => ExperienceTypeApi
    guestApi: (token?: string | undefined) => GuestApi
    liveApi: (token?: string | undefined) => LiveApi
    messageApi: (token?: string | undefined) => MessageApi
    notificationApi: (token?: string | undefined) => NotificationApi
    orderApi: (token?: string | undefined) => OrderApi
    paymentSystemApi: (token?: string | undefined) => PaymentSystemApi
    postalCodeApi: (token?: string | undefined) => PostalCodeApi
    producerApi: (token?: string | undefined) => ProducerApi
    productApi: (token?: string | undefined) => ProductApi
    productTagApi: (token?: string | undefined) => ProductTagApi
    productTypeApi: (token?: string | undefined) => ProductTypeApi
    promotionApi: (token?: string | undefined) => PromotionApi
    scheduleApi: (token?: string | undefined) => ScheduleApi
    shippingApi: (token?: string | undefined) => ShippingApi
    shopApi: (token?: string | undefined) => ShopApi
    spotTypeApi: (token?: string | undefined) => SpotTypeApi
    topApi: (token?: string | undefined) => TopApi
    uploadApi: (token?: string | undefined) => UploadApi
    userApi: (token?: string | undefined) => UserApi
    guestApi: (token?: string | undefined) => GuestApi
    videoApi: (token?: string | undefined) => VideoApi
  }
}
