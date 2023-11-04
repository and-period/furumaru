import { PiniaPluginContext } from 'pinia'
import ApiClientFactory from './helpter/factory'
import { AuthApi, CartApi, ProductApi } from '~/types/api'

function apiClientInjector({ store }: PiniaPluginContext) {
  const apiClientFactory = new ApiClientFactory()

  // 認証関連のAPIをstoreに定義
  const authApiClient = (token?: string) =>
    apiClientFactory.create(AuthApi, token)

  // 商品関連のAPIをstoreに定義
  const productApiClient = (token?: string): ProductApi =>
    apiClientFactory.create<ProductApi>(ProductApi, token)

  // カート関連のAPIをstoreに定義
  const cartApiClient = (token?: string): CartApi =>
    apiClientFactory.create<CartApi>(CartApi, token)

  store.authApiClient = authApiClient
  store.productApiClient = productApiClient
  store.cartApiClient = cartApiClient
}

/**
 * Pinia に ApiClient をinjection する Nuxt プラグイン
 */
const apiClientPlugin = defineNuxtPlugin(({ $pinia }) => {
  $pinia.use(apiClientInjector) // type: ignore
})

export default apiClientPlugin
