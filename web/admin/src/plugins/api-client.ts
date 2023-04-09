import { PiniaPluginContext } from 'pinia'

import ApiClientFactory from '~/lib/api/factory'
import { AdministratorApi, AuthApi, CategoryApi, ContactApi, CoordinatorApi, MessageApi, NotificationApi, OrderApi, ProducerApi, ProductApi, ProductTypeApi, PromotionApi, ShippingApi } from '~/types/api'

function apiClientInjector ({ store }: PiniaPluginContext) {
  const factory = new ApiClientFactory()

  const administratorApiClient = (token?: string) => factory.create(AdministratorApi, token)
  const authApiClient = (token?: string) => factory.create(AuthApi, token)
  const categoryApiClient = (token?: string) => factory.create(CategoryApi, token)
  const contactApiCilent = (token?: string) => factory.create(ContactApi, token)
  const coordinatorApiClient = (token?: string) => factory.create(CoordinatorApi, token)
  const messageApiClient = (token?: string) => factory.create(MessageApi, token)
  const notificationApiClient = (token?: string) => factory.create(NotificationApi, token)
  const orderApiClient = (token?: string) => factory.create(OrderApi, token)
  const producerApiClient = (token?: string) => factory.create(ProducerApi, token)
  const productApiClient = (token?: string) => factory.create(ProductApi, token)
  const productTypeApiClient = (token?: string) => factory.create(ProductTypeApi, token)
  const promotionApiClient = (token?: string) => factory.create(PromotionApi, token)
  const shippingApiClient = (token?: string) => factory.create(ShippingApi, token)

  store.administratorApiClient = administratorApiClient
  store.authApiClient = authApiClient
  store.categoryApiClient = categoryApiClient
  store.contactApiClient = contactApiCilent
  store.coordinatorApiClient = coordinatorApiClient
  store.messageApiClient = messageApiClient
  store.notificationApiClient = notificationApiClient
  store.orderApiClient = orderApiClient
  store.producerApiClient = producerApiClient
  store.productApiClient = productApiClient
  store.productTypeApiClient = productTypeApiClient
  store.promotionApiClient = promotionApiClient
  store.shippingApiClient = shippingApiClient
}

export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.$pinia.use(apiClientInjector)
})
