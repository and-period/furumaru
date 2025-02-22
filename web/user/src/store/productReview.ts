import { useAuthStore } from './auth'
import type { CreateProductReviewRequest, ProductReview } from '~/types/api'

export const useProductReviewStore = defineStore('product-review', {
  state: () => {
    return {
      reviewsFetchState: {
        isLoading: false,
        productId: '',
      },
      reviewsResponse: {
        reviews: [] as ProductReview[],
        nextToken: '',
      },
    }
  },

  getters: {
    reviews: (state) => {
      return state.reviewsResponse.reviews
    },
  },

  actions: {
    async fetchReviews(productId: string, nextToken?: string, limit: number = 20): Promise<void> {
      try {
        this.reviewsFetchState.isLoading = true
        this.reviewsFetchState.productId = productId
        const response = await this.productApiClient().v1ListProductReviews({
          productId,
          nextToken,
          limit,
        })
        this.reviewsResponse = response
      }
      catch (error) {
        return this.errorHandler(error)
      }
      finally {
        this.reviewsFetchState.isLoading = false
      }
    },

    /**
     * 商品レビューを投稿する関数
     * @param productId レビュー対象の商品ID
     * @param payload レビュー内容
     * @returns
     */
    async postReview(productId: string, payload: CreateProductReviewRequest): Promise<void> {
      const authStore = useAuthStore()
      const { accessToken } = authStore
      try {
        await this.productApiClient(accessToken).v1CreateProductReview({
          productId: productId,
          body: payload,
        })
      }
      catch (error) {
        return this.errorHandler(error)
      }
    },
  },
})
