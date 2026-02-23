import { useApiClient } from '~/composables/useApiClient'
import { ProductTagApi } from '~/types/api/v1'
import type { CreateProductTagRequest, ProductTag, UpdateProductTagRequest, V1ProductTagsGetRequest, V1ProductTagsPostRequest, V1ProductTagsProductTagIdDeleteRequest, V1ProductTagsProductTagIdPatchRequest } from '~/types/api/v1'

export const useProductTagStore = defineStore('productTag', () => {
  const { create, errorHandler } = useApiClient()
  const productTagApi = () => create(ProductTagApi)

  const productTag = ref<ProductTag>({} as ProductTag)
  const productTags = ref<ProductTag[]>([])
  const total = ref<number>(0)

  async function fetchProductTags(limit = 20, offset = 0, orders: string[] = []): Promise<void> {
    try {
      const params: V1ProductTagsGetRequest = {
        limit,
        offset,
        orders: orders.join(','),
      }
      const res = await productTagApi().v1ProductTagsGet(params)
      productTags.value = res.productTags
      total.value = res.total
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function searchProductTags(name = '', productTagIds: string[] = []): Promise<void> {
    try {
      const params: V1ProductTagsGetRequest = { name }
      const res = await productTagApi().v1ProductTagsGet(params)
      const merged: ProductTag[] = []
      productTags.value.forEach((pt: ProductTag): void => {
        if (!productTagIds.includes(pt.id)) {
          return
        }
        merged.push(pt)
      })
      res.productTags.forEach((pt: ProductTag): void => {
        if (merged.find((v): boolean => v.id === pt.id)) {
          return
        }
        merged.push(pt)
      })
      productTags.value = merged
      total.value = res.total
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function createProductTag(payload: CreateProductTagRequest): Promise<void> {
    try {
      const params: V1ProductTagsPostRequest = { createProductTagRequest: payload }
      const res = await productTagApi().v1ProductTagsPost(params)
      productTags.value.unshift(res.productTag)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります。',
        409: 'この商品タグ名はすでに登録されています。',
      })
    }
  }

  async function updateProductTag(productTagId: string, payload: UpdateProductTagRequest): Promise<void> {
    try {
      const params: V1ProductTagsProductTagIdPatchRequest = {
        productTagId,
        updateProductTagRequest: payload,
      }
      await productTagApi().v1ProductTagsProductTagIdPatch(params)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります。',
        404: 'この商品タグは存在しません。',
        409: 'この商品タグ名はすでに登録されています。',
      })
    }
  }

  async function deleteProductTag(productTagId: string): Promise<void> {
    try {
      const params: V1ProductTagsProductTagIdDeleteRequest = { productTagId }
      await productTagApi().v1ProductTagsProductTagIdDelete(params)
    }
    catch (err) {
      errorHandler(err, { 404: 'この商品タグは存在しません。' })
    }
  }

  return {
    productTag,
    productTags,
    total,
    fetchProductTags,
    searchProductTags,
    createProductTag,
    updateProductTag,
    deleteProductTag,
  }
})
