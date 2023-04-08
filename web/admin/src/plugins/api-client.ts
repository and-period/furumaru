import { Plugin } from '@nuxt/types'
import { PiniaPluginContext } from 'pinia'

import ApiClientFactory from './factory'

import { OrderApi } from '~/types/api'

function apiClientInjector ({ store }: PiniaPluginContext) {
  const apiClientFactory = new ApiClientFactory()

  const orderApiClient = (token: string) =>
    apiClientFactory.create(OrderApi, token)

  store.orderApiClient = orderApiClient
}

const apiClientPlugin: Plugin = (ctx, _inject) => {
  ctx.$pinia.use(apiClientInjector)
}

export default apiClientPlugin
