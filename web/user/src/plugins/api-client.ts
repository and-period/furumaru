import { PiniaPluginContext } from 'pinia'
import ApiClientFactory from './helpter/factory'
import { AuthApi } from '~/types/api'

function apiClientInjector ({ store }: PiniaPluginContext) {
  const apiClientFactory = new ApiClientFactory()

  // authのAPIをstoreに定義
  const authApiClient = (token?: string) => apiClientFactory.create(AuthApi, token)
  store.authApiClient = authApiClient
}

/**
 * Pinia に ApiClient をinjection する Nuxt プラグイン
 */
const apiClientPlugin = defineNuxtPlugin(({ $pinia }) => {
  $pinia.use(apiClientInjector)
})

export default apiClientPlugin
