import { defineStore } from 'pinia'

import { apiClient } from '~/plugins/api-client'
import type { CreateProductTagRequest, ProductTag, UpdateProductTagRequest } from '~/types/api'

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
        const res = await apiClient.productTagApi().v1ListProductTags(limit, offset, '', orders.join(','))
        this.productTags = res.data.productTags
        this.total = res.data.total
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
        const res = await apiClient.productTagApi().v1ListProductTags(undefined, undefined, name)
        const productTags: ProductTag[] = []
        this.productTags.forEach((productTag: ProductTag): void => {
          if (!productTagIds.includes(productTag.id)) {
            return
          }
          productTags.push(productTag)
        })
        res.data.productTags.forEach((productTag: ProductTag): void => {
          if (productTags.find((v): boolean => v.id === productTag.id)) {
            return
          }
          productTags.push(productTag)
        })
        this.productTags = productTags
        this.total = res.data.total
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
        const res = await apiClient.productTagApi().v1CreateProductTag(payload)
        this.productTags.unshift(res.data.productTag)
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
        await apiClient.productTagApi().v1UpdateProductTag(productTagId, payload)
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
        await apiClient.productTagApi().v1DeleteProductTag(productTagId)
      }
      catch (err) {
        this.errorHandler(err, { 404: 'この商品タグは存在しません。' })
      }
    },
  },
})
