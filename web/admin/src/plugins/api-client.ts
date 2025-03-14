import { client } from './axios'
import type { AxiosInstance } from 'axios'
import {
  AddressApi,
  AdministratorApi,
  AuthApi,
  BroadcastApi,
  CategoryApi,
  Configuration,
  ContactApi,
  CoordinatorApi,
  ExperienceApi,
  ExperienceTypeApi,
  GuestApi,
  LiveApi,
  MessageApi,
  NotificationApi,
  OrderApi,
  OtherApi,
  PaymentSystemApi,
  ProducerApi,
  ProductApi,
  ProductTagApi,
  ProductTypeApi,
  PromotionApi,
  ScheduleApi,
  ShippingApi,
  ShopApi,
  SpotTypeApi,
  TopApi,
  UserApi,
  VideoApi,
} from '~/types/api'

let apiClient: ApiClient

export class ApiClient {
  basePath: string
  config: Configuration
  instance: AxiosInstance

  constructor(basePath: string) {
    this.basePath = basePath
    this.config = new Configuration()
    this.instance = client
  }

  addressApi() {
    return new AddressApi(this.config, this.basePath, this.instance)
  }

  administratorApi() {
    return new AdministratorApi(this.config, this.basePath, this.instance)
  }

  authApi() {
    return new AuthApi(this.config, this.basePath, this.instance)
  }

  broadcastApi() {
    return new BroadcastApi(this.config, this.basePath, this.instance)
  }

  categoryApi() {
    return new CategoryApi(this.config, this.basePath, this.instance)
  }

  contactApi() {
    return new ContactApi(this.config, this.basePath, this.instance)
  }

  coordinatorApi() {
    return new CoordinatorApi(this.config, this.basePath, this.instance)
  }

  experienceApi() {
    return new ExperienceApi(this.config, this.basePath, this.instance)
  }

  experienceTypeApi() {
    return new ExperienceTypeApi(this.config, this.basePath, this.instance)
  }

  liveApi() {
    return new LiveApi(this.config, this.basePath, this.instance)
  }

  messageApi() {
    return new MessageApi(this.config, this.basePath, this.instance)
  }

  notificationApi() {
    return new NotificationApi(this.config, this.basePath, this.instance)
  }

  orderApi() {
    return new OrderApi(this.config, this.basePath, this.instance)
  }

  paymentSystemApi() {
    return new PaymentSystemApi(this.config, this.basePath, this.instance)
  }

  producerApi() {
    return new ProducerApi(this.config, this.basePath, this.instance)
  }

  productApi() {
    return new ProductApi(this.config, this.basePath, this.instance)
  }

  productTagApi() {
    return new ProductTagApi(this.config, this.basePath, this.instance)
  }

  productTypeApi() {
    return new ProductTypeApi(this.config, this.basePath, this.instance)
  }

  promotionApi() {
    return new PromotionApi(this.config, this.basePath, this.instance)
  }

  scheduleApi() {
    return new ScheduleApi(this.config, this.basePath, this.instance)
  }

  shippingApi() {
    return new ShippingApi(this.config, this.basePath, this.instance)
  }

  shopApi() {
    return new ShopApi(this.config, this.basePath, this.instance)
  }

  spotTypeApi() {
    return new SpotTypeApi(this.config, this.basePath, this.instance)
  }

  topApi() {
    return new TopApi(this.config, this.basePath, this.instance)
  }

  userApi() {
    return new UserApi(this.config, this.basePath, this.instance)
  }

  otherApi() {
    return new OtherApi(this.config, this.basePath, this.instance)
  }

  guestApi() {
    return new GuestApi(this.config, this.basePath, this.instance)
  }

  videoApi() {
    return new VideoApi(this.config, this.basePath, this.instance)
  }
}

export default defineNuxtPlugin(() => {
  const runtimeConfig = useRuntimeConfig()
  const baseUrl = runtimeConfig.public.API_BASE_URL

  apiClient = new ApiClient(baseUrl)

  return { provide: {} }
})

export { apiClient }
