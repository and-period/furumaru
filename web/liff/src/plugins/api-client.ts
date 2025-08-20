import type { PiniaPluginContext } from 'pinia';
import ApiClientFactory from './helper/factory';
import {
  ProductApi,
} from '@/types/api';

function apiClientInjector({ store }: PiniaPluginContext) {
  const apiClientFactory = new ApiClientFactory();

  // 商品関連のAPIをstoreに定義
  const productApiClient = (token?: string): ProductApi =>
    apiClientFactory.create<ProductApi>(ProductApi, token);

  store.productApiClient = productApiClient;
}

export default defineNuxtPlugin(() => {
  const { $pinia } = useNuxtApp();

  $pinia.use(apiClientInjector);
});
