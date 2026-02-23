import { useApiClient } from '~/composables/useApiClient'
import { fileUpload } from './helper'
import { useCategoryStore } from './category'
import { useCoordinatorStore } from './coordinator'
import { useProductTypeStore } from './product-type'
import { useProductTagStore } from './product-tag'
import { useProducerStore } from './producer'
import { ProductApi, ProductReviewApi, UploadApi } from '~/types/api/v1'
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

export const useProductStore = defineStore('product', () => {
  const { create, errorHandler } = useApiClient()
  const productApi = () => create(ProductApi)
  const productReviewApi = () => create(ProductReviewApi)
  const uploadApi = () => create(UploadApi)

  const product = ref<Product>({} as Product)
  const products = ref<Product[]>([])
  const totalItems = ref<number>(0)

  async function fetchProducts(limit = 20, offset = 0): Promise<void> {
    try {
      const params: V1ProductsGetRequest = { limit, offset }
      const res = await productApi().v1ProductsGet(params)

      const coordinatorStore = useCoordinatorStore()
      const producerStore = useProducerStore()
      const categoryStore = useCategoryStore()
      const productTypeStore = useProductTypeStore()
      const productTagStore = useProductTagStore()
      products.value = res.products
      totalItems.value = res.total
      coordinatorStore.coordinators = res.coordinators
      producerStore.producers = res.producers
      categoryStore.categories = res.categories
      productTypeStore.productTypes = res.productTypes
      productTagStore.productTags = res.productTags
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function searchProducts(
    name = '',
    producerId = '',
    productIds: string[] = [],
  ): Promise<void> {
    try {
      const params: V1ProductsGetRequest = { name, producerId }
      const res = await productApi().v1ProductsGet(params)
      const merged: Product[] = []
      products.value.forEach((p: Product): void => {
        if (!productIds.includes(p.id)) {
          return
        }
        merged.push(p)
      })
      res.products.forEach((p: Product): void => {
        if (merged.find((v): boolean => v.id === p.id)) {
          return
        }
        merged.push(p)
      })
      products.value = merged
      totalItems.value = res.total
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function getProduct(productId: string): Promise<ProductResponse> {
    try {
      const params: V1ProductsProductIdGetRequest = { productId }
      const res = await productApi().v1ProductsProductIdGet(params)

      const coordinatorStore = useCoordinatorStore()
      const producerStore = useProducerStore()
      const categoryStore = useCategoryStore()
      const productTypeStore = useProductTypeStore()
      const productTagStore = useProductTagStore()
      product.value = res.product
      coordinatorStore.coordinators = [res.coordinator]
      producerStore.producers = [res.producer]
      categoryStore.categories = [res.category]
      productTypeStore.productTypes = [res.productType]
      productTagStore.productTags = res.productTags
      return res
    }
    catch (err: any) {
      return errorHandler(err, {
        403: '商品を閲覧する権限がありません',
        404: '対象の商品が存在しません',
      })
    }
  }

  async function uploadProductMedia(payload: File): Promise<string> {
    const contentType = payload.type
    try {
      const res = await getProductMediaUploadUrl(contentType)
      return await fileUpload(uploadApi(), payload, res.key, res.url)
    }
    catch (err) {
      return errorHandler(err, {
        400: 'このファイルはアップロードできません。',
      })
    }
  }

  async function getProductMediaUploadUrl(contentType: string): Promise<UploadURLResponse> {
    const body: GetUploadURLRequest = { fileType: contentType }
    if (contentType.includes('image/')) {
      try {
        const params: V1UploadProductsImagePostRequest = { getUploadURLRequest: body }
        return await uploadApi().v1UploadProductsImagePost(params)
      }
      catch (err) {
        return errorHandler(err, { 400: 'このファイルはアップロードできません。' })
      }
    }
    if (contentType.includes('video/')) {
      try {
        const params: V1UploadProductsVideoPostRequest = { getUploadURLRequest: body }
        return await uploadApi().v1UploadProductsVideoPost(params)
      }
      catch (err) {
        return errorHandler(err, { 400: 'このファイルはアップロードできません。' })
      }
    }
    throw new Error('不明なMINEタイプです。')
  }

  async function createProduct(payload: CreateProductRequest): Promise<void> {
    try {
      const params: V1ProductsPostRequest = {
        createProductRequest: { ...payload, inventory: Number(payload.inventory) },
      }
      await productApi().v1ProductsPost(params)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります',
        403: '商品を登録する権限がありません',
      })
    }
  }

  async function updateProduct(productId: string, payload: UpdateProductRequest) {
    try {
      const params: V1ProductsProductIdPatchRequest = {
        productId,
        updateProductRequest: { ...payload, inventory: Number(payload.inventory) },
      }
      await productApi().v1ProductsProductIdPatch(params)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります',
        403: '商品を更新する権限がありません',
        404: '対象の商品が存在しません',
      })
    }
  }

  async function deleteProduct(productId: string) {
    try {
      const params: V1ProductsProductIdDeleteRequest = { productId }
      await productApi().v1ProductsProductIdDelete(params)
      const index = products.value.findIndex(p => p.id === productId)
      products.value.splice(index, 1)
      totalItems.value--
    }
    catch (err) {
      return errorHandler(err, {
        403: '商品を削除する権限がありません',
        404: '対象の商品が存在しません',
      })
    }
  }

  async function createProductReview(productId: string, payload: CreateProductReviewRequest): Promise<void> {
    try {
      const params: V1ProductsProductIdReviewsPostRequest = {
        productId,
        createProductReviewRequest: payload,
      }
      await productReviewApi().v1ProductsProductIdReviewsPost(params)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります',
        403: 'レビューを投稿する権限がありません',
        404: '対象の商品が存在しません',
      })
    }
  }

  return {
    product,
    products,
    totalItems,
    fetchProducts,
    searchProducts,
    getProduct,
    uploadProductMedia,
    getProductMediaUploadUrl,
    createProduct,
    updateProduct,
    deleteProduct,
    createProductReview,
  }
})
