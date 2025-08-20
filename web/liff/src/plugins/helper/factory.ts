import { Configuration } from '@/types/api';
import type { BaseAPI } from '@/types/api/runtime';

/**
 * API クライアントのインスタンスを生成するファクトリ
 */
export default class ApiClientFactory {
  create<T extends BaseAPI>(
    Client: new (config: Configuration) => T,
    token?: string,
  ): T {
    const runtimeConfig = useRuntimeConfig();
    const config = new Configuration({
      headers: {
        Authorization: `Bearer ${token}`,
      },
      basePath: runtimeConfig.public.API_BASE_URL,
      credentials: 'include',
    });
    return new Client(config);
  }
}
