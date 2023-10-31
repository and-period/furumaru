import { defineStore } from 'pinia'
import {
  Category,
  Coordinator,
  Producer,
  Product,
  ProductTag,
  ProductType,
  ProductsResponse,
} from '~/types/api'

export const useProductStore = defineStore('product', {
  state: () => {
    return {
      isLoading: false,
      total: 0,
      products: [] as Product[],
      producers: [] as Producer[],
      coordinators: [] as Coordinator[],
      categories: [] as Category[],
      productTypes: [] as ProductType[],
      productTags: [] as ProductTag[],
    }
  },

  actions: {
    async fetchProducts(limit = 20, offset = 0): Promise<void> {
      this.isLoading = true
      const {
        total,
        products,
        producers,
        productTags,
        productTypes,
      }: ProductsResponse = await this.productApiClient().v1ListProducts({
        limit,
        offset,
      })
      this.total = total
      this.products = products
      this.producers = producers
      this.productTags = productTags
      this.productTypes = productTypes
      this.isLoading = false
    },
  },
})
