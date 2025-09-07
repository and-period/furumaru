import type { PiniaPluginContext } from 'pinia';
import ApiClientFactory from './helper/factory';
import { ProductApi } from '@/types/api/';
import { AuthApi, ProductApi as FacilityProductApi } from '@/types/api/facility';

function apiClientInjector({ store }: PiniaPluginContext) {
  const apiClientFactory = new ApiClientFactory();

  // 商品関連のAPIをstoreに定義
  const productApiClient = (token?: string): ProductApi =>
    apiClientFactory.create<ProductApi>(ProductApi, token);

  store.productApiClient = productApiClient;

  // 施設向け 認証APIクライアントをstoreに定義
  const facilityAuthApiClient = (token?: string): AuthApi =>
    apiClientFactory.createFacility<AuthApi>(AuthApi, token);

  store.facilityAuthApiClient = facilityAuthApiClient;

  // 施設向け 商品APIクライアントをstoreに定義
  const facilityProductApiClient = (token?: string): FacilityProductApi =>
    apiClientFactory.createFacility<FacilityProductApi>(FacilityProductApi, token);

  // @ts-expect-error 型拡張は各storeで宣言
  store.facilityProductApiClient = facilityProductApiClient;
}

export default defineNuxtPlugin(() => {
  const { $pinia } = useNuxtApp();

  $pinia.use(apiClientInjector);
});
