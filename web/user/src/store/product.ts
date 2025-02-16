import { defineStore } from 'pinia'
import type {
  Category,
  Coordinator,
  Producer,
  Product,
  ProductResponse,
  ProductTag,
  ProductType,
  ProductsResponse,
} from '~/types/api'

export const useProductStore = defineStore('product', {
  state: () => {
    return {
      productsFetchState: {
        isLoading: false,
      },
      productsResponse: {
        total: 0,
        products: [] as Product[],
        producers: [] as Producer[],
        coordinators: [] as Coordinator[],
        categories: [] as Category[],
        productTypes: [] as ProductType[],
        productTags: [] as ProductTag[],
      },
      productFetchState: {
        isLoading: false,
      },
      productResponse: {} as ProductResponse,
    }
  },

  actions: {
    async fetchProducts(limit = 20, offset = 0): Promise<void> {
      try {
        this.productsFetchState.isLoading = true
        const response: ProductsResponse
          = await this.productApiClient().v1ListProducts({
            limit,
            offset,
          })
        this.productsResponse = response
      }
      catch (error) {
        return this.errorHandler(error)
      }
      finally {
        this.productsFetchState.isLoading = false
      }
    },

    async fetchProduct(id: string): Promise<ProductResponse> {
      try {
        this.productFetchState.isLoading = true
        const response = await this.productApiClient().v1GetProduct({
          productId: id,
        })
        this.productResponse = response
        this.productFetchState.isLoading = false
        return response
      }
      catch (error) {
        return this.errorHandler(error)
      }
    },
  },

  getters: {
    totalProductsCount(state) {
      return state.productsResponse.total
    },

    products(state) {
      return state.productsResponse.products.map((product) => {
        const thumbnail = product.media.find(m => m.isThumbnail)
        return {
          ...product,
          // 在庫があるかのフラグ
          hasStock: product.inventory > 0,
          // サムネイル画像のマッピング
          thumbnail,
          // 生産者情報をマッピング
          producer: state.productsResponse.producers.find(
            producer => producer.id === product.producerId,
          ),
          // 商品タイプをマッピング
          productType: state.productsResponse.productTypes.find(
            productType => productType.id === product.productTypeId,
          ),
          // コーディネーター情報をマッピング
          coordinator: state.productsResponse.coordinators.find(
            coordinator => coordinator.id === product.coordinatorId,
          ),
          // カテゴリ情報をマッピング
          category: state.productsResponse.categories.find(
            category => category.id === product.categoryId,
          ),
          // 商品タグをマッピング
          productTags: product.productTagIds.map(id =>
            state.productsResponse.productTags.find(
              productTag => productTag.id === id,
            ),
          ),
          // サムネイルが動画かどうかのフラグ
          thumbnailIsVideo: thumbnail ? thumbnail.url.endsWith('.mp4') : false,
        }
      })
    },

    product(state) {
      const thumbnail = state.productResponse.product?.media.find(
        m => m.isThumbnail,
      )

      return {
        ...state.productResponse.product,
        // 在庫があるかのフラグ
        hasStock: state.productResponse.product?.inventory > 0,
        // サムネイル画像のマッピング
        thumbnail,
        // 生産者情報をマッピング
        producer: state.productResponse.producer,
        // コーディネーター情報をマッピング
        coordinator: state.productResponse.coordinator,
        // 商品タグをマッピング
        productTags: state.productResponse.product?.productTagIds.map(id =>
          state.productResponse.productTags.find(
            productTag => productTag.id === id,
          ),
        ),
        // サムネイルが動画かどうかのフラグ
        thumbnailIsVideo: thumbnail ? thumbnail.url.endsWith('.mp4') : false,
      }
    },
  },
})
