import type {
  CategoriesResponse,
  Category,
  CreateCategoryRequest,
  UpdateCategoryRequest,
  V1CategoriesCategoryIdDeleteRequest,
  V1CategoriesCategoryIdPatchRequest,
  V1CategoriesGetRequest,
  V1CategoriesPostRequest,
} from '~/types/api/v1'

export const useCategoryStore = defineStore('category', {
  state: () => ({
    categories: [] as Category[],
    total: 0,
  }),

  actions: {
    /**
     * カテゴリ一覧を取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     * @param orders ソートキー
     */
    async fetchCategories(limit = 20, offset = 0, orders = []): Promise<void> {
      try {
        const res = await this.listCategories(limit, offset, '', orders)
        this.categories = res.categories
        this.total = res.total
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * カテゴリを検索をする非同期関数
     * @param name カテゴリ名(あいまい検索)
     * @param categoryIds stateの更新時に残しておく必要があるカテゴリ情報
     */
    async searchCategories(name = '', categoryIds: string[] = []): Promise<void> {
      try {
        const res = await this.listCategories(undefined, undefined, name, [])
        const categories: Category[] = []
        this.categories.forEach((category: Category): void => {
          if (!categoryIds.includes(category.id)) {
            return
          }
          categories.push(category)
        })
        res.categories.forEach((category: Category): void => {
          if (categories.find((v): boolean => v.id === category.id)) {
            return
          }
          categories.push(category)
        })
        this.categories = categories
        this.total = res.total
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * カテゴリを追加取得する非同期関数
     * @param limit 取得上限数
     * @param offset 取得開始位置
     * @param orders ソートキー
     */
    async moreCategories(limit = 20, offset = 0, orders = []): Promise<void> {
      try {
        const res = await this.listCategories(limit, offset, '', orders)
        this.categories.push(...res.categories)
        this.total = res.total
      }
      catch (err) {
        return this.errorHandler(err)
      }
    },

    /**
     * カテゴリを新規登録する非同期関数
     * @param payload
     */
    async createCategory(payload: CreateCategoryRequest): Promise<void> {
      try {
        const params: V1CategoriesPostRequest = {
          createCategoryRequest: payload,
        }
        const res = await this.categoryApi().v1CategoriesPost(params)
        this.categories.unshift(res.category)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります',
          409: 'このカテゴリー名はすでに登録されています。',
        })
      }
    },

    /**
     * カテゴリを編集する非同期関数
     * @param categoryId カテゴリID
     * @param payload
     */
    async updateCategory(categoryId: string, payload: UpdateCategoryRequest) {
      try {
        const params: V1CategoriesCategoryIdPatchRequest = {
          categoryId,
          updateCategoryRequest: payload,
        }
        await this.categoryApi().v1CategoriesCategoryIdPatch(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          400: '必須項目が不足しているか、内容に誤りがあります',
          404: '対象のカテゴリーが存在しません',
          409: 'このカテゴリー名はすでに登録されています。',
        })
      }
      this.fetchCategories()
    },

    /**
     * カテゴリを削除する非同期関数
     * @param categoryId カテゴリID
     */
    async deleteCategory(categoryId: string): Promise<void> {
      try {
        const params: V1CategoriesCategoryIdDeleteRequest = {
          categoryId,
        }
        await this.categoryApi().v1CategoriesCategoryIdDelete(params)
      }
      catch (err) {
        return this.errorHandler(err, {
          404: '対象のカテゴリーが存在しません',
          412: '品目と紐付いているため削除できません',
        })
      }
      this.fetchCategories()
    },

    async listCategories(limit = 20, offset = 0, name = '', orders: string[] = []): Promise<CategoriesResponse> {
      const params: V1CategoriesGetRequest = {
        limit,
        offset,
        name,
        orders: orders.join(','),
      }
      const res = await this.categoryApi().v1CategoriesGet(params)
      return { ...res }
    },
  },
})
