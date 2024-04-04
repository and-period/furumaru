import type { PiniaPluginContext } from 'pinia'
import ApiClientFactory from './helpter/factory'
import {
  AddressApi,
  AuthApi,
  CartApi,
  CheckoutApi,
  CoordinatorApi,
  ProductApi,
  ScheduleApi,
  StatusApi,
  TopApi,
  OrderApi,
  PromotionApi,
  AuthUserApi,
  OtherApi,
} from '~/types/api'

function apiClientInjector({ store }: PiniaPluginContext) {
  const apiClientFactory = new ApiClientFactory()

  // 認証関連のAPIをstoreに定義
  const authApiClient = (token?: string) =>
    apiClientFactory.create(AuthApi, token)

  // ユーザー関連のAPIをstoreに定義
  const authUserApiClient = (token?: string) =>
    apiClientFactory.create(AuthUserApi, token)

  // トップページのAPIをstoreに定義
  const topPageApiClient = (token?: string) =>
    apiClientFactory.create<TopApi>(TopApi, token)

  // 商品関連のAPIをstoreに定義
  const productApiClient = (token?: string): ProductApi =>
    apiClientFactory.create<ProductApi>(ProductApi, token)

  // カート関連のAPIをstoreに定義
  const cartApiClient = (token?: string): CartApi =>
    apiClientFactory.create<CartApi>(CartApi, token)

  // 住所関連のAPIをStoreに定義
  const addressApiClient = (token?: string): AddressApi =>
    apiClientFactory.create<AddressApi>(AddressApi, token)

  // チェックアウト用のAPIをStoreに定義
  const checkoutApiClient = (token?: string): CheckoutApi =>
    apiClientFactory.create(CheckoutApi, token)

  // スケジュール
  const scheduleApiClient = (token?: string): ScheduleApi =>
    apiClientFactory.create(ScheduleApi, token)

  // コーディネーター関連のAPIをStoreに定義
  const coordinatorApiClient = (token?: string): CoordinatorApi =>
    apiClientFactory.create(CoordinatorApi, token)

  // 支払いシステム状態取得APIをStoreに定義
  const statusApiClient = (): StatusApi => apiClientFactory.create(StatusApi)

  // 注文履歴取得APIをStoreに定義
  const orderApiClient = (token?: string): OrderApi =>
    apiClientFactory.create(OrderApi, token)

  // クーポンコードのAPIをStoreに定義
  const promotionApiClient = (token?: string): PromotionApi =>
    apiClientFactory.create(PromotionApi, token)

  // その他のAPIをStoreに定義
  const otherApiClient = (token?: string) =>
    apiClientFactory.create(OtherApi, token)

  store.authApiClient = authApiClient
  store.authUserApiClient = authUserApiClient
  store.topPageApiClient = topPageApiClient
  store.productApiClient = productApiClient
  store.cartApiClient = cartApiClient
  store.addressApiClient = addressApiClient
  store.checkoutApiClient = checkoutApiClient
  store.scheduleApiClient = scheduleApiClient
  store.coordinatorApiClient = coordinatorApiClient
  store.statusApiClient = statusApiClient
  store.orderApiClient = orderApiClient
  store.promotionApiClient = promotionApiClient
  store.otherApiClient = otherApiClient
}

/**
 * Pinia に ApiClient をinjection する Nuxt プラグイン
 */
const apiClientPlugin = defineNuxtPlugin(({ $pinia }) => {
  $pinia.use(apiClientInjector) // type: ignore
})

export default apiClientPlugin
