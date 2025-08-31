import { Configuration } from '@/types/api';
import { Configuration as FacilityConfiguration } from '@/types/api/facility';

/**
 * API クライアントのインスタンスを生成するファクトリ
 */
export default class ApiClientFactory {
  // 一般（購入者向け）API クライアント
  create<T>(Client: new (config: Configuration) => T, token?: string): T {
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

  // 施設向け API クライアント
  // TODO: createと統合したい
  createFacility<T>(
    Client: new (config: FacilityConfiguration) => T,
    token?: string,
  ): T {
    const runtimeConfig = useRuntimeConfig();
    const config = new FacilityConfiguration({
      headers: {
        Authorization: `Bearer ${token}`,
      },
      basePath: runtimeConfig.public.API_BASE_URL,
      credentials: 'include',
    });
    return new Client(config);
  }
}
