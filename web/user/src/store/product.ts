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
      _products: [] as Product[],
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
        coordinators,
        categories,
        products,
        producers,
        productTags,
        productTypes,
      }: ProductsResponse = await this.productApiClient().v1ListProducts({
        limit,
        offset,
      })
      this.total = total
      this._products = products
      this.producers = producers
      this.productTags = productTags
      this.productTypes = productTypes
      this.coordinators = coordinators
      this.categories = categories
      this.isLoading = false
    },
  },

  getters: {
    products(state) {
      return state._products.map((product) => {
        return {
          ...product,
          // 在庫があるかのフラグ
          hasStock: product.inventory > 0,
          // サムネイル画像のマッピング
          thumbnail: product.media.find((m) => m.isThumbnail),
          // 生産者情報をマッピング
          producer: state.producers.find(
            (producer) => producer.id === product.producerId,
          ),
          // 商品タイプをマッピング
          productType: state.productTypes.find(
            (productType) => productType.id === product.productTypeId,
          ),
          // コーディネーター情報をマッピング
          coordinator: state.coordinators.find(
            (coordinator) => coordinator.id === product.coordinatorId,
          ),
          // カテゴリ情報をマッピング
          category: state.categories.find(
            (category) => category.id === product.categoryId,
          ),
          // 商品タグをマッピング
          productTags: product.productTagIds.map((id) =>
            state.productTags.find((productTag) => productTag.id === id),
          ),
        }
      })
    },
  },
})
