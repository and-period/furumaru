import type { ProductReview } from '~/types/api'

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
  },
})
