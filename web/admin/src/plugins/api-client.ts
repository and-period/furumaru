import type { PiniaPluginContext } from 'pinia'
import { AdministratorApi, AuthApi, BroadcastApi, CategoryApi, Configuration, ContactApi, CoordinatorApi, ExperienceApi, ExperienceTypeApi, GuestApi, LiveApi, MessageApi, NotificationApi, OrderApi, PaymentSystemApi, PostalCodeApi, ProducerApi, ProductApi, ProductTagApi, ProductTypeApi, PromotionApi, ScheduleApi, ShippingApi, ShopApi, SpotTypeApi, TopApi, UploadApi, UserApi, VideoApi } from '~/types/api/v1'
import type { BaseAPI } from '~/types/api/v1/runtime'

/**
 * API クライアントのインスタンスを生成するファクトリ
 */
class ApiClientFactory {
  create<T extends BaseAPI>(Client: new (config: Configuration) => T, token?: string): T {
    const runtimeConfig = useRuntimeConfig()
    const baseUrl = runtimeConfig.public.API_BASE_URL

    const config = new Configuration({
      headers: {
        Authorization: `Bearer ${token}`,
      },
      basePath: baseUrl,
      credentials: 'include',
    })

    return new Client(config)
  }
}

/**
 * Pinia の store に API クライアントを注入する
 */
function apiClientInjector({ store }: PiniaPluginContext) {
  const apiClientFactory = new ApiClientFactory()

  store.administratorApi = (token?: string): AdministratorApi => apiClientFactory.create<AdministratorApi>(AdministratorApi, token)

  store.authApi = (token?: string): AuthApi => apiClientFactory.create<AuthApi>(AuthApi, token)

  store.broadcastApi = (token?: string): BroadcastApi => apiClientFactory.create<BroadcastApi>(BroadcastApi, token)

  store.categoryApi = (token?: string): CategoryApi => apiClientFactory.create<CategoryApi>(CategoryApi, token)

  store.contactApi = (token?: string): ContactApi => apiClientFactory.create<ContactApi>(ContactApi, token)

  store.coordinatorApi = (token?: string): CoordinatorApi => apiClientFactory.create<CoordinatorApi>(CoordinatorApi, token)

  store.experienceApi = (token?: string): ExperienceApi => apiClientFactory.create<ExperienceApi>(ExperienceApi, token)

  store.experienceTypeApi = (token?: string): ExperienceTypeApi => apiClientFactory.create<ExperienceTypeApi>(ExperienceTypeApi, token)

  store.guestApi = (token?: string): GuestApi => apiClientFactory.create<GuestApi>(GuestApi, token)

  store.liveApi = (token?: string): LiveApi => apiClientFactory.create<LiveApi>(LiveApi, token)

  store.messageApi = (token?: string): MessageApi => apiClientFactory.create<MessageApi>(MessageApi, token)

  store.notificationApi = (token?: string): NotificationApi => apiClientFactory.create<NotificationApi>(NotificationApi, token)

  store.orderApi = (token?: string): OrderApi => apiClientFactory.create<OrderApi>(OrderApi, token)

  store.paymentSystemApi = (token?: string): PaymentSystemApi => apiClientFactory.create<PaymentSystemApi>(PaymentSystemApi, token)

  store.postalCodeApi = (token?: string): PostalCodeApi => apiClientFactory.create<PostalCodeApi>(PostalCodeApi, token)

  store.producerApi = (token?: string): ProducerApi => apiClientFactory.create<ProducerApi>(ProducerApi, token)

  store.productApi = (token?: string): ProductApi => apiClientFactory.create<ProductApi>(ProductApi, token)

  store.productTagApi = (token?: string): ProductTagApi => apiClientFactory.create<ProductTagApi>(ProductTagApi, token)

  store.productTypeApi = (token?: string): ProductTypeApi => apiClientFactory.create<ProductTypeApi>(ProductTypeApi, token)

  store.promotionApi = (token?: string): PromotionApi => apiClientFactory.create<PromotionApi>(PromotionApi, token)

  store.scheduleApi = (token?: string): ScheduleApi => apiClientFactory.create<ScheduleApi>(ScheduleApi, token)

  store.shippingApi = (token?: string): ShippingApi => apiClientFactory.create<ShippingApi>(ShippingApi, token)

  store.shopApi = (token?: string): ShopApi => apiClientFactory.create<ShopApi>(ShopApi, token)

  store.spotTypeApi = (token?: string): SpotTypeApi => apiClientFactory.create<SpotTypeApi>(SpotTypeApi, token)

  store.topApi = (token?: string): TopApi => apiClientFactory.create<TopApi>(TopApi, token)

  store.uploadApi = (token?: string): UploadApi => apiClientFactory.create<UploadApi>(UploadApi, token)

  store.userApi = (token?: string): UserApi => apiClientFactory.create<UserApi>(UserApi, token)

  store.videoApi = (token?: string): VideoApi => apiClientFactory.create<VideoApi>(VideoApi, token)
}

/**
 * Pinia に ApiClient をinjection する Nuxt プラグイン
 */
const apiClientPlugin = defineNuxtPlugin(({ $pinia }) => {
  $pinia.use(apiClientInjector) // type: ignore
})

export default apiClientPlugin
