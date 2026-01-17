import { fileUpload } from './helper'
import { useCategoryStore } from './category'
import { useCoordinatorStore } from './coordinator'
import { useProductTypeStore } from './product-type'
import { useProductTagStore } from './product-tag'
import { useProducerStore } from './producer'
import type {
  CreateProductRequest,
  CreateProductReviewRequest,
  ProductResponse,
  Product,
  UpdateProductRequest,
  V1ProductsGetRequest,
  V1ProductsProductIdGetRequest,
  UploadURLResponse,
  GetUploadURLRequest,
  V1UploadProductsImagePostRequest,
  V1UploadProductsVideoPostRequest,
  V1ProductsPostRequest,
  V1ProductsProductIdPatchRequest,
  V1ProductsProductIdDeleteRequest,
  V1ProductsProductIdReviewsPostRequest,
} from '~/types/api/v1'

export const useProductStore = defineStore('product', {
  state: () => ({
    product: {} as Product,
    products: [] as Product[],
    totalItems: 0,
  }),

  actions: {
    /**
     * 商品一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     * @returns
     */
    async fetchProducts(limit = 20, offset = 0): Promise<void> {
      try {
        const params: V1ProductsGetRequest = {
          limit,
          offset,
        }
        const res = await this.productApi().v1ProductsGet(params)

        const coordinatorStore = useCoordinatorStore()
        const producerStore = useProducerStore()
        const categoryStore = useCategoryStore()
        const productTypeStore = useProductTypeStore()
        const productTagStore = useProductTagStore()
        this.products = res.products
        this.totalItems = res.total
        coordinatorStore.coordinators = res.coordinators
        producerStore.producers = res.producers
        categoryStore.categories = res.categories
        productTypeStore.productTypes = res.productTypes
        productTagStore.productTags = res.productTags
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 商品情報を検索する非同期関数
     * @param name 商品名(あいまい検索)
     * @param producerId 生産者ID
     * @param productIds stateの更新時に残しておく必要がある商品情報
     */
    async searchProducts(
      name = '',
      producerId = '',
      productIds: string[] = [],
    ): Promise<void> {
      try {
        const params: V1ProductsGetRequest = {
          name,
          producerId,
        }
        const res = await this.productApi().v1ProductsGet(params)
        const products: Product[] = []
        this.products.forEach((product: Product): void => {
          if (!productIds.includes(product.id)) {
            return
          }
          products.push(product)
        })
        res.products.forEach((product: Product): void => {
          if (products.find((v): boolean => v.id === product.id)) {
            return
          }
          products.push(product)
        })
        this.products = products
        this.totalItems = res.total
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 商品詳細を取得する非同期関数
     * @param productId
     * @returns
     */
    async getProduct(productId: string): Promise<ProductResponse> {
      try {
        const params: V1ProductsProductIdGetRequest = {
          productId,
        }
        const res = await this.productApi().v1ProductsProductIdGet(params)

        const coordinatorStore = useCoordinatorStore()
        const producerStore = useProducerStore()
        const categoryStore = useCategoryStore()
        const productTypeStore = useProductTypeStore()
        const productTagStore = useProductTagStore()
        this.product = res.product
        coordinatorStore.coordinators = [res.coordinator]
        producerStore.producers = [res.producer]
        categoryStore.categories = [res.category]
        productTypeStore.productTypes = [res.productType]
        productTagStore.productTags = res.productTags
        return res
      }
      catch (err: any) {
        return this.errorHandler(err, {
          403: '商品を閲覧する権限がありません',
          404: '対象の商品が存在しません',
        })
      }
    },

    /**
     * 商品メディアファイルをアップロードする非同期関数
     * @param payload
     * @returns
     */
    async uploadProductMedia(payload: File): Promise<string> {
      const contentType = payload.type
      try {
        const res = await this.getProductMediaUploadUrl(contentType)

        return await fileUpload(this.uploadApi(), payload, res.key, res.url)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: 'このファイルはアップロードできません。',
        })
      }
    },

    async getProductMediaUploadUrl(contentType: string): Promise<UploadURLResponse> {
      const body: GetUploadURLRequest = {
        fileType: contentType,
      }
      if (contentType.includes('image/')) {
        try {
          const params: V1UploadProductsImagePostRequest = {
            getUploadURLRequest: body,
          }
          const res = await this.uploadApi().v1UploadProductsImagePost(params)
          return res
        }
        catch (err) {
          return this.errorHandler(err, {
            400: 'このファイルはアップロードできません。',
          })
        }
      }
      if (contentType.includes('video/')) {
        try {
          const params: V1UploadProductsVideoPostRequest = {
            getUploadURLRequest: body,
          }
          const res = await this.uploadApi().v1UploadProductsVideoPost(params)
          return res
        }
        catch (err) {
          return this.errorHandler(err, {
            400: 'このファイルはアップロードできません。',
          })
        }
      }
      throw new Error('不明なMINEタイプです。')
    },

    /**
     * 商品を作成する非同期関数
     */
    async createProduct(payload: CreateProductRequest): Promise<void> {
      try {
        const params: V1ProductsPostRequest = {
          createProductRequest: {
            ...payload,
            inventory: Number(payload.inventory),
          },
        }
        await this.productApi().v1ProductsPost(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります',
          403: '商品を登録する権限がありません',
        })
      }
    },

    /**
     * 商品を更新する関数
     * @param productId
     * @param payload
     */
    async updateProduct(productId: string, payload: UpdateProductRequest) {
      try {
        const params: V1ProductsProductIdPatchRequest = {
          productId,
          updateProductRequest: {
            ...payload,
            inventory: Number(payload.inventory),
          },
        }
        await this.productApi().v1ProductsProductIdPatch(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります',
          403: '商品を更新する権限がありません',
          404: '対象の商品が存在しません',
        })
      }
    },

    /**
     * 商品を削除する関数
     * @param productId
     * @returns
     */
    async deleteProduct(productId: string) {
      try {
        const params: V1ProductsProductIdDeleteRequest = {
          productId,
        }
        await this.productApi().v1ProductsProductIdDelete(params)
        const index = this.products.findIndex(
          product => product.id === productId,
        )
        this.products.splice(index, 1)
        this.totalItems--
      }
      catch (err) {
        return this.errorHandler(err, {
          403: '商品を削除する権限がありません',
          404: '対象の商品が存在しません',
        })
      }
    },

    /**
     * ダミー商品レビューを作成する関数
     * @param productId 商品ID
     * @param payload レビュー情報
     */
    async createProductReview(productId: string, payload: CreateProductReviewRequest): Promise<void> {
      try {
        const params: V1ProductsProductIdReviewsPostRequest = {
          productId,
          createProductReviewRequest: payload,
        }
        await this.productReviewApi().v1ProductsProductIdReviewsPost(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります',
          403: 'レビューを投稿する権限がありません',
          404: '対象の商品が存在しません',
        })
      }
    },
  },
})
