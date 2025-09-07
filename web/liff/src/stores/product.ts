import { defineStore } from 'pinia';
import type { ProductsResponse, ProductApi } from '@/types/api';

declare module 'pinia' {
  export interface PiniaCustomProperties {
    productApiClient: (token?: string) => ProductApi;
  }
}

export const useProductStore = defineStore('product', {
  state: () => ({
    products: [] as ProductsResponse['products'],
    isLoading: false,
    error: null as string | null,
  }),

  actions: {
    async fetchProducts() {
      this.isLoading = true;
      this.error = null;

      try {
        const route = useRoute();
        const runtimeConfig = useRuntimeConfig();

        const facilityId = String(route.params.facilityId ?? '');
        if (!facilityId) {
          throw new Error('facilityId is not specified in params.');
        }

        const response = await $fetch<ProductsResponse>(
          `${runtimeConfig.public.API_BASE_URL}/facilities/${facilityId}/products`,
          {
            method: 'GET',
            credentials: 'include',
            query: {
              limit: 20,
              offset: 0,
            },
          },
        );

        this.products = response.products || [];
      }
      catch (error) {
        this.error = error instanceof Error ? error.message : 'Failed to fetch products';
        console.error('Failed to fetch products:', error);
      }
      finally {
        this.isLoading = false;
      }
    },

    getProductById(id: string) {
      return this.products.find(product => product.id === id);
    },
  },
});
