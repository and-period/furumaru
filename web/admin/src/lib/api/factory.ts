import { Configuration } from '~/types/api'
import { BaseAPI } from '~/types/api/base'

export default class ApiClientFactory {
  create<T extends BaseAPI> (Client: new (config: Configuration) => T, token?: string): T {
    const runtimeConfig = useRuntimeConfig()
    const basePath = runtimeConfig.public.API_BASE_URL

    const config = new Configuration({ accessToken: token, basePath })
    return new Client(config)
  }
}
