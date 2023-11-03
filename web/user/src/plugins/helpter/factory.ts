import { Configuration } from '~/types/api'
import { BaseAPI } from '~/types/api/runtime'

/**
 * API クライアントのインスタンスを生成するファクトリ
 */
export default class ApiClientFactory {
  create<T extends BaseAPI>(
    Client: new (config: Configuration) => T,
    token?: string,
  ): T {
    const { API_BASE_URL: basePath } = useRuntimeConfig()
    const config = new Configuration({ accessToken: token, basePath })
    return new Client(config)
  }
}
