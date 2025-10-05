import { defineStore } from 'pinia';
import type { ProductApi } from '@/types/api/v1';
import type { ProductApi as FacilityProductApi, ProductsResponse } from '@/types/api/facility';

declare module 'pinia' {
  export interface PiniaCustomProperties {
    productApiClient: (token?: string) => ProductApi;
    facilityProductApiClient: (token?: string) => FacilityProductApi;
  }
}

export const useProductStore = defineStore('product', {
  state: () => ({
    products: [] as ProductsResponse['products'],
    isLoading: false,
    error: null as string | null,
  }),

  actions: {
    async fetchProducts(facilityId: string) {
      this.isLoading = true;
      this.error = null;

      try {
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
