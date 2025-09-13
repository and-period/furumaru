import type { CreateProductTagRequest, ProductTag, UpdateProductTagRequest, V1ProductTagsGetRequest, V1ProductTagsPostRequest, V1ProductTagsProductTagIdDeleteRequest, V1ProductTagsProductTagIdPatchRequest } from '~/types/api/v1'

export const useProductTagStore = defineStore('productTag', {
  state: () => ({
    productTag: {} as ProductTag,
    productTags: [] as ProductTag[],
    total: 0,
  }),

  actions: {
    /**
     * 商品タグ一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     * @param orders ソートキー
     */
    async fetchProductTags(limit = 20, offset = 0, orders: string[] = []): Promise<void> {
      try {
        const params: V1ProductTagsGetRequest = {
          limit,
          offset,
          orders: orders.join(','),
        }
        const res = await this.productTagApi().v1ProductTagsGet(params)
        this.productTags = res.productTags
        this.total = res.total
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 商品タグを検索する非同期関数
     * @param name 商品タグ名(あいまい検索)
     * @param productTagIds stateの更新時に残しておく必要がある商品タグ情報
     */
    async searchProductTags(name = '', productTagIds: string[] = []): Promise<void> {
      try {
        const params: V1ProductTagsGetRequest = {
          name,
        }
        const res = await this.productTagApi().v1ProductTagsGet(params)
        const productTags: ProductTag[] = []
        this.productTags.forEach((productTag: ProductTag): void => {
          if (!productTagIds.includes(productTag.id)) {
            return
          }
          productTags.push(productTag)
        })
        res.productTags.forEach((productTag: ProductTag): void => {
          if (productTags.find((v): boolean => v.id === productTag.id)) {
            return
          }
          productTags.push(productTag)
        })
        this.productTags = productTags
        this.total = res.total
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 商品タグを新規登録する非同期関数
     * @param payload
     */
    async createProductTag(payload: CreateProductTagRequest): Promise<void> {
      try {
        const params: V1ProductTagsPostRequest = {
          createProductTagRequest: payload,
        }
        const res = await this.productTagApi().v1ProductTagsPost(params)
        this.productTags.unshift(res.productTag)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります。',
          409: 'この商品タグ名はすでに登録されています。',
        })
      }
    },

    /**
     * 商品タグを更新する非同期関数
     * @param productTagId 商品タグID
     * @param payload
     */
    async updateProductTag(productTagId: string, payload: UpdateProductTagRequest): Promise<void> {
      try {
        const params: V1ProductTagsProductTagIdPatchRequest = {
          productTagId,
          updateProductTagRequest: payload,
        }
        await this.productTagApi().v1ProductTagsProductTagIdPatch(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります。',
          404: 'この商品タグは存在しません。',
          409: 'この商品タグ名はすでに登録されています。',
        })
      }
    },

    /**
     * 商品タグを削除する非同期関数
     * @param productTagId 商品タグID
     */
    async deleteProductTag(productTagId: string): Promise<void> {
      try {
        const params: V1ProductTagsProductTagIdDeleteRequest = {
          productTagId,
        }
        await this.productTagApi().v1ProductTagsProductTagIdDelete(params)
      }
      catch (err) {
        this.errorHandler(err, { 404: 'この商品タグは存在しません。' })
      }
    },
  },
})
