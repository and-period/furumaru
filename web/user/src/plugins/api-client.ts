import { PiniaPluginContext } from 'pinia'
import ApiClientFactory from './helpter/factory'
import { AuthApi, ProductApi } from '~/types/api'

function apiClientInjector({ store }: PiniaPluginContext) {
  const apiClientFactory = new ApiClientFactory()

  // authのAPIをstoreに定義
  const authApiClient = (token?: string) =>
    apiClientFactory.create(AuthApi, token)

  const productApiClient = (token?: string): ProductApi =>
    apiClientFactory.create<ProductApi>(ProductApi, token)

  store.authApiClient = authApiClient
  store.productApiClient = productApiClient
}

/**
 * Pinia に ApiClient をinjection する Nuxt プラグイン
 */
const apiClientPlugin = defineNuxtPlugin(({ $pinia }) => {
  $pinia.use(apiClientInjector)
})

export default apiClientPlugin
