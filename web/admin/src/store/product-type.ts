import { fileUpload } from './helper'
import { useCategoryStore } from './category'
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

export const useProductTypeStore = defineStore('productType', {
  state: () => ({
    productTypes: [] as ProductType[],
    totalItems: 0,
  }),

  actions: {
    /**
     * 品目一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     * @param orders ソートキー
     */
    async fetchProductTypes(
      limit = 20,
      offset = 0,
      orders = [],
    ): Promise<void> {
      try {
        const params: V1CategoriesCategoryIdProductTypesGetRequest = {
          categoryId: '-', // 全カテゴリ対象
          limit,
          offset,
          orders: orders.join(','),
        }
        const res = await this.productTypeApi().v1CategoriesCategoryIdProductTypesGet(params)

        const categoryStore = useCategoryStore()
        this.productTypes = res.productTypes
        this.totalItems = res.total
        categoryStore.categories = res.categories || []
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * カテゴリに紐づく品目一覧を取得する非同期関数
     * @param categoryId カテゴリID
     * @param limit 取得上限数
     * @param offset 取得開始位置
     * @returns
     */
    async fetchProductTypesByCategoryId(
      categoryId: string,
      limit = 20,
      offset = 0,
    ): Promise<void> {
      if (categoryId === '') {
        this.productTypes = []
        this.totalItems = 0
        return
      }

      try {
        const params: V1CategoriesCategoryIdProductTypesGetRequest = {
          categoryId,
          limit,
          offset,
        }
        const res = await this.productTypeApi().v1CategoriesCategoryIdProductTypesGet(params)

        this.productTypes = res.productTypes
        this.totalItems = res.total
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 品目を検索する非同期関数
     * @param name 品目名(あいまい検索)
     * @param categoryId カテゴリ名
     * @param productTypeIds stateの更新時に残しておく必要がある品目情報
     */
    async searchProductTypes(
      name = '',
      categoryId = '',
      productTypeIds: string[] = [],
    ): Promise<void> {
      try {
        const params: V1CategoriesCategoryIdProductTypesGetRequest = {
          categoryId: categoryId || '-',
          name,
        }
        const res = await this.productTypeApi().v1CategoriesCategoryIdProductTypesGet(params)
        const productTypes: ProductType[] = []
        this.productTypes.forEach((productType: ProductType): void => {
          if (!productTypeIds.includes(productType.id)) {
            return
          }
          productTypes.push(productType)
        })
        res.productTypes.forEach((productType: ProductType): void => {
          if (productTypes.find((v): boolean => v.id === productType.id)) {
            return
          }
          productTypes.push(productType)
        })
        this.productTypes = productTypes
        this.totalItems = res.total
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * 品目を新規登録する非同期関数
     * @param categoryId 品目の親となるカテゴリのID
     * @param payload
     * @returns
     */
    async createProductType(
      categoryId: string,
      payload: CreateProductTypeRequest,
    ): Promise<void> {
      try {
        const params: V1CategoriesCategoryIdProductTypesPostRequest = {
          categoryId,
          createProductTypeRequest: payload,
        }
        const res = await this.productTypeApi().v1CategoriesCategoryIdProductTypesPost(params)
        this.productTypes.unshift(res.productType)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります。',
          409: '対象の品目名はすでに登録されているため、登録できません。',
        })
      }
    },

    /**
     * 登録済みの品目を更新する非同期関数
     * @param categoryId カテゴリID
     * @param productTypeId 品目ID
     * @param payload 品目情報
     * @returns
     */
    async updateProductType(
      categoryId: string,
      productTypeId: string,
      payload: UpdateProductTypeRequest,
    ) {
      try {
        const params: V1CategoriesCategoryIdProductTypesProductTypeIdPatchRequest = {
          categoryId,
          productTypeId,
          updateProductTypeRequest: payload,
        }
        await this.productTypeApi().v1CategoriesCategoryIdProductTypesProductTypeIdPatch(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります。',
          404: '対象の商品種別または品目が存在しません。',
          409: '対象の品目名はすでに登録されているため、登録できません。',
        })
      }
    },

    /**
     * 品目を削除する非同期関数
     * @param categoryId カテゴリID
     * @param productTypeId 品目ID
     * @returns
     */
    async deleteProductType(
      categoryId: string,
      productTypeId: string,
    ): Promise<void> {
      try {
        const params: V1CategoriesCategoryIdProductTypesProductTypeIdDeleteRequest = {
          categoryId,
          productTypeId,
        }
        await this.productTypeApi().v1CategoriesCategoryIdProductTypesProductTypeIdDelete(params)
        this.fetchProductTypes()
      }
      catch (err) {
        return this.errorHandler(err, {
          404: '対象の商品種別または品目が存在しません。',
        })
      }
    },

    /**
     * 品目画像をアップロードする関数
     * @param payload 品目画像のファイルオブジェクト
     * @returns アップロード後の品目画像のパスを含んだオブジェクト
     */
    async uploadProductTypeIcon(payload: File): Promise<string> {
      const contentType = payload.type
      try {
        const params: V1UploadProductTypesIconPostRequest = {
          getUploadURLRequest: {
            fileType: contentType,
          },
        }
        const res = await this.uploadApi().v1UploadProductTypesIconPost(params)

        return await fileUpload(this.uploadApi(), payload, res.key, res.url)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: 'このファイルはアップロードできません。',
        })
      }
    },
  },
})
