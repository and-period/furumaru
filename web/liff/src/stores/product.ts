import { defineStore } from 'pinia';
import type { ProductApi, ProductResponse } from '@/types/api/v1';
import type {
  ProductApi as FacilityProductApi,
  ProductsResponse,
  ProductResponse as FacilityProductResponse,
} from '@/types/api/facility';

declare module 'pinia' {
  export interface PiniaCustomProperties {
    productApiClient: (token?: string) => ProductApi;
    facilityProductApiClient: (token?: string) => FacilityProductApi;
  }
}

type ProductDetailResponse = ProductResponse | FacilityProductResponse;

export const useProductStore = defineStore('product', {
  state: () => ({
    products: [] as ProductsResponse['products'],
    productDetails: {} as Record<string, ProductDetailResponse>,
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

    async fetchProductDetail(productId: string) {
      if (!productId) {
        throw new Error('Product ID is required to fetch product detail.');
      }

      this.isLoading = true;
      this.error = null;

      try {
        const api = this.productApiClient();
        const response = await api.productsProductIdGet({ productId });

        this.productDetails[productId] = response;
        this.upsertProduct(response.product);

        return response;
      }
      catch (error) {
        this.error = error instanceof Error ? error.message : 'Failed to fetch product detail';
        console.error('Failed to fetch product detail:', error);
        throw error;
      }
      finally {
        this.isLoading = false;
      }
    },

    async fetchFacilityProductDetail(facilityId: string, productId: string) {
      if (!facilityId) {
        throw new Error('Facility ID is required to fetch product detail.');
      }
      if (!productId) {
        throw new Error('Product ID is required to fetch product detail.');
      }

      this.isLoading = true;
      this.error = null;

      try {
        const api = this.facilityProductApiClient();
        const response = await api.facilitiesFacilityIdProductsProductIdGet({
          facilityId,
          productId,
        });

        this.productDetails[productId] = response;
        this.upsertProduct(response.product);

        return response;
      }
      catch (error) {
        this.error = error instanceof Error ? error.message : 'Failed to fetch facility product detail';
        console.error('Failed to fetch facility product detail:', error);
        throw error;
      }
      finally {
        this.isLoading = false;
      }
    },

    upsertProduct(product: ProductsResponse['products'][number]) {
      const index = this.products.findIndex(existing => existing.id === product.id);
      if (index >= 0) {
        this.products.splice(index, 1, product);
      }
      else {
        this.products.push(product);
      }
    },
  },
});
