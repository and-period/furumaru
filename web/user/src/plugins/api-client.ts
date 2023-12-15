import type { PiniaPluginContext } from 'pinia'
import ApiClientFactory from './helpter/factory'
import {
  AddressApi,
  AuthApi,
  CartApi,
  CheckoutApi,
  ProductApi,
  ScheduleApi,
  TopApi,
} from '~/types/api'

function apiClientInjector({ store }: PiniaPluginContext) {
  const apiClientFactory = new ApiClientFactory()

  // 認証関連のAPIをstoreに定義
  const authApiClient = (token?: string) =>
    apiClientFactory.create(AuthApi, token)

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

  store.authApiClient = authApiClient
  store.topPageApiClient = topPageApiClient
  store.productApiClient = productApiClient
  store.cartApiClient = cartApiClient
  store.addressApiClient = addressApiClient
  store.checkoutApiClient = checkoutApiClient
  store.scheduleApiClient = scheduleApiClient
}

/**
 * Pinia に ApiClient をinjection する Nuxt プラグイン
 */
const apiClientPlugin = defineNuxtPlugin(({ $pinia }) => {
  $pinia.use(apiClientInjector) // type: ignore
})

export default apiClientPlugin
