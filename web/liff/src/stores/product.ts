import { defineStore } from 'pinia';
import type { ProductApi } from '@/types/api';
import type { ProductApi as FacilityProductApi, TypesProductsResponse } from '@/types/api/facility';

declare module 'pinia' {
  export interface PiniaCustomProperties {
    productApiClient: (token?: string) => ProductApi;
    facilityProductApiClient: (token?: string) => FacilityProductApi;
  }
}

export const useProductStore = defineStore('product', {
  state: () => ({
    products: [] as TypesProductsResponse['products'],
    isLoading: false,
    error: null as string | null,
  }),

  actions: {
    async fetchProducts() {
      this.isLoading = true;
      this.error = null;

      try {
        const route = useRoute();
        const facilityId = String(route.params.facilityId ?? '');
        if (!facilityId) {
          throw new Error('facilityId is not specified in params.');
        }

        const api = this.facilityProductApiClient();
        const response = await api.facilitiesFacilityIdProductsGet({
          facilityId,
          limit: 20,
          offset: 0,
        });

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
      return this.products?.find(product => product.id === id);
    },
  },
});
