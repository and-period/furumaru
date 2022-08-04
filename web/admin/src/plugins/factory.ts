import { Configuration } from '~/types/api'
import { BaseAPI } from '~/types/api/base'

const basePath = process.env.API_BASE_URL || 'http://localhost:18010'

export default class ApiClientFactory {
  create<T extends BaseAPI>(
    Client: new (config: Configuration) => T,
    token?: string
  ): T {
    const config = new Configuration({ accessToken: token, basePath })
    return new Client(config)
  }
}
