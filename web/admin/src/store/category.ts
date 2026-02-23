import { useApiClient } from '~/composables/useApiClient'
import { CategoryApi } from '~/types/api/v1'
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

export const useCategoryStore = defineStore('category', () => {
  const { create, errorHandler } = useApiClient()
  const categoryApi = () => create(CategoryApi)

  const categories = ref<Category[]>([])
  const total = ref<number>(0)

  async function listCategories(limit = 20, offset = 0, name = '', orders: string[] = []): Promise<CategoriesResponse> {
    const params: V1CategoriesGetRequest = {
      limit,
      offset,
      name,
      orders: orders.join(','),
    }
    const res = await categoryApi().v1CategoriesGet(params)
    return { ...res }
  }

  async function fetchCategories(limit = 20, offset = 0, orders = []): Promise<void> {
    try {
      const res = await listCategories(limit, offset, '', orders)
      categories.value = res.categories
      total.value = res.total
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function searchCategories(name = '', categoryIds: string[] = []): Promise<void> {
    try {
      const res = await listCategories(undefined, undefined, name, [])
      const merged: Category[] = []
      categories.value.forEach((category: Category): void => {
        if (!categoryIds.includes(category.id)) {
          return
        }
        merged.push(category)
      })
      res.categories.forEach((category: Category): void => {
        if (merged.find((v): boolean => v.id === category.id)) {
          return
        }
        merged.push(category)
      })
      categories.value = merged
      total.value = res.total
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function moreCategories(limit = 20, offset = 0, orders = []): Promise<void> {
    try {
      const res = await listCategories(limit, offset, '', orders)
      categories.value.push(...res.categories)
      total.value = res.total
    }
    catch (err) {
      return errorHandler(err)
    }
  }

  async function createCategory(payload: CreateCategoryRequest): Promise<void> {
    try {
      const params: V1CategoriesPostRequest = { createCategoryRequest: payload }
      const res = await categoryApi().v1CategoriesPost(params)
      categories.value.unshift(res.category)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります',
        409: 'このカテゴリー名はすでに登録されています。',
      })
    }
  }

  async function updateCategory(categoryId: string, payload: UpdateCategoryRequest) {
    try {
      const params: V1CategoriesCategoryIdPatchRequest = {
        categoryId,
        updateCategoryRequest: payload,
      }
      await categoryApi().v1CategoriesCategoryIdPatch(params)
    }
    catch (err) {
      return errorHandler(err, {
        400: '必須項目が不足しているか、内容に誤りがあります',
        404: '対象のカテゴリーが存在しません',
        409: 'このカテゴリー名はすでに登録されています。',
      })
    }
    fetchCategories()
  }

  async function deleteCategory(categoryId: string): Promise<void> {
    try {
      const params: V1CategoriesCategoryIdDeleteRequest = { categoryId }
      await categoryApi().v1CategoriesCategoryIdDelete(params)
    }
    catch (err) {
      return errorHandler(err, {
        404: '対象のカテゴリーが存在しません',
        412: '品目と紐付いているため削除できません',
      })
    }
    fetchCategories()
  }

  return {
    categories,
    total,
    fetchCategories,
    searchCategories,
    moreCategories,
    createCategory,
    updateCategory,
    deleteCategory,
    listCategories,
  }
})
