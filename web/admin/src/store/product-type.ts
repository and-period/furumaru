import { useApiClient } from '~/composables/useApiClient'
import { fileUpload } from './helper'
import { useCategoryStore } from './category'
import { ProductTypeApi, UploadApi } from '~/types/api/v1'
import type {
  CreateProductTypeRequest,
  ProductType,
  UpdateProductTypeRequest,
  V1CategoriesCategoryIdProductTypesGetRequest,
  V1CategoriesCategoryIdProductTypesPostRequest,
  V1CategoriesCategoryIdProductTypesProductTypeIdDeleteRequest,
  V1CategoriesCategoryIdProductTypesProductTypeIdPatchRequest,
  V1UploadProductTypesIconPostRequest,
} from '~/types/api/v1'

export const useProductTypeStore = defineStore('productType', () => {
  const { create, errorHandler } = useApiClient()
  const productTypeApi = () => create(ProductTypeApi)
  const uploadApi = () => create(UploadApi)

  const productTypes = ref<ProductType[]>([])
  const totalItems = ref<number>(0)

  async function fetchProductTypes(limit = 20, offset = 0, orders = []): Promise<void> {
    try {
      const params: V1CategoriesCategoryIdProductTypesGetRequest = {
        categoryId: '-',
        limit,
        offset,
        orders: orders.join(','),
      }
      const res = await productTypeApi().v1CategoriesCategoryIdProductTypesGet(params)

      const categoryStore = useCategoryStore()
      productTypes.value = res.productTypes
      totalItems.value = res.total
      categoryStore.categories = res.categories || []
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function fetchProductTypesByCategoryId(categoryId: string, limit = 20, offset = 0): Promise<void> {
    if (categoryId === '') {
      productTypes.value = []
      totalItems.value = 0
      return
    }

    try {
      const params: V1CategoriesCategoryIdProductTypesGetRequest = { categoryId, limit, offset }
      const res = await productTypeApi().v1CategoriesCategoryIdProductTypesGet(params)
      productTypes.value = res.productTypes
      totalItems.value = res.total
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function searchProductTypes(name = '', categoryId = '', productTypeIds: string[] = []): Promise<void> {
    try {
      const params: V1CategoriesCategoryIdProductTypesGetRequest = {
        categoryId: categoryId || '-',
        name,
      }
      const res = await productTypeApi().v1CategoriesCategoryIdProductTypesGet(params)
      const merged: ProductType[] = []
      productTypes.value.forEach((productType: ProductType): void => {
        if (!productTypeIds.includes(productType.id)) {
          return
        }
        merged.push(productType)
      })
      res.productTypes.forEach((productType: ProductType): void => {
        if (merged.find((v): boolean => v.id === productType.id)) {
          return
        }
        merged.push(productType)
      })
      productTypes.value = merged
      totalItems.value = res.total
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function createProductType(categoryId: string, payload: CreateProductTypeRequest): Promise<void> {
    try {
      const params: V1CategoriesCategoryIdProductTypesPostRequest = {
        categoryId,
        createProductTypeRequest: payload,
      }
      const res = await productTypeApi().v1CategoriesCategoryIdProductTypesPost(params)
      productTypes.value.unshift(res.productType)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります。',
        409: '対象の品目名はすでに登録されているため、登録できません。',
      })
    }
  }

  async function updateProductType(categoryId: string, productTypeId: string, payload: UpdateProductTypeRequest) {
    try {
      const params: V1CategoriesCategoryIdProductTypesProductTypeIdPatchRequest = {
        categoryId,
        productTypeId,
        updateProductTypeRequest: payload,
      }
      await productTypeApi().v1CategoriesCategoryIdProductTypesProductTypeIdPatch(params)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります。',
        404: '対象の商品種別または品目が存在しません。',
        409: '対象の品目名はすでに登録されているため、登録できません。',
      })
    }
  }

  async function deleteProductType(categoryId: string, productTypeId: string): Promise<void> {
    try {
      const params: V1CategoriesCategoryIdProductTypesProductTypeIdDeleteRequest = { categoryId, productTypeId }
      await productTypeApi().v1CategoriesCategoryIdProductTypesProductTypeIdDelete(params)
      fetchProductTypes()
    }
    catch (err) {
      return errorHandler(err, { 404: '対象の商品種別または品目が存在しません。' })
    }
  }

  async function uploadProductTypeIcon(payload: File): Promise<string> {
    const contentType = payload.type
    try {
      const params: V1UploadProductTypesIconPostRequest = {
        getUploadURLRequest: { fileType: contentType },
      }
      const res = await uploadApi().v1UploadProductTypesIconPost(params)
      return await fileUpload(uploadApi(), payload, res.key, res.url)
    }
    catch (err) {
      return errorHandler(err, { 400: 'このファイルはアップロードできません。' })
    }
  }

  return {
    productTypes,
    totalItems,
    fetchProductTypes,
    fetchProductTypesByCategoryId,
    searchProductTypes,
    createProductType,
    updateProductType,
    deleteProductType,
    uploadProductTypeIcon,
  }
})
