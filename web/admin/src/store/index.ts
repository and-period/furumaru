import { Configuration } from '~/types/api'
import { BaseAPI } from '~/types/api/base'

const BASE_API = process.env.API_BASE_URL || 'http://localhost:18010'

export class ApiClientFactory {
  create<T extends BaseAPI>(
    Client: new (config: Configuration) => T,
    token?: string
  ): T {
    const config = new Configuration({ accessToken: token, basePath: BASE_API })
    return new Client(config)
  }
}
