// Example of how to use the generated LIFF API client
import { defineStore } from 'pinia';
import { ProductApi, Configuration, ProductsResponse } from '@/types/api';

// Example configuration for API client
const apiConfig = new Configuration({
  basePath: process.env.NODE_ENV === 'production' ? 'https://api.furumaru.and-period.co.jp' : 'http://localhost:8080',
  // Add authentication headers if needed
  // headers: {
  //   Authorization: 'Bearer YOUR_TOKEN_HERE',
  // },
});

// Initialize the Product API client
const productApi = new ProductApi(apiConfig);

export const useApiExampleStore = defineStore('apiExample', {
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
        // Use the generated API client to fetch products
        const response = await productApi.v1ProductsGet({
          limit: 20,
          offset: 0,
        });
        
        this.products = response.products || [];
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Failed to fetch products';
        console.error('Failed to fetch products:', error);
      } finally {
        this.isLoading = false;
      }
    },

    async fetchProductById(productId: string) {
      this.isLoading = true;
      this.error = null;
      
      try {
        // Use the generated API client to fetch a specific product
        const response = await productApi.v1ProductsProductIdGet({
          productId,
        });
        
        return response.product;
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Failed to fetch product';
        console.error('Failed to fetch product:', error);
        throw error;
      } finally {
        this.isLoading = false;
      }
    },
  },
});