<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useAlert, usePagination } from '~/lib/hooks'
import { useAuthStore, useCategoryStore, useCommonStore, useProductTypeStore } from '~/store'
import type { Category, CreateCategoryRequest, CreateProductTypeRequest, ProductType, UpdateCategoryRequest } from '~/types/api/v1'
import type { ImageUploadStatus } from '~/types/props'

const commonStore = useCommonStore()
const authStore = useAuthStore()
const categoryStore = useCategoryStore()
const productTypeStore = useProductTypeStore()
const categoryPagination = usePagination()
const productTypePagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const { adminType } = storeToRefs(authStore)
const { categories, total: categoryTotal } = storeToRefs(categoryStore)
const { productTypes, totalItems: productTypeTotal } = storeToRefs(productTypeStore)

const initialProductType: ProductType = {
  id: '',
  categoryId: '',
  name: '',
  iconUrl: '',
  createdAt: 0,
  updatedAt: 0,
}

const loading = ref<boolean>(false)
const selector = ref<string>('categories')
const selectedCategory = ref<Category>()
const selectedProductType = ref<ProductType>({ ...initialProductType })
const createCategoryDialog = ref<boolean>(false)
const updateCategoryDialog = ref<boolean>(false)
const deleteCategoryDialog = ref<boolean>(false)
const createProductTypeDialog = ref<boolean>(false)
const updateProductTypeDialog = ref<boolean>(false)
const deleteProductTypeDialog = ref<boolean>(false)
const createCategoryFormData = ref<CreateCategoryRequest>({
  name: '',
})
const updateCategoryFormData = ref<UpdateCategoryRequest>({
  name: '',
})
const createProductTypeFormData = ref<CreateProductTypeRequest>({
  name: '',
  iconUrl: '',
})
const updateProductTypeFormData = ref<CreateProductTypeRequest>({
  name: '',
  iconUrl: '',
})
const createProductTypeIconUploadStatus = ref<ImageUploadStatus>({
  error: false,
  message: '',
})
const updateProductTypeIconUploadStatus = ref<ImageUploadStatus>({
  error: false,
  message: '',
})

watch(updateProductTypeDialog, (): void => {
  if (updateProductTypeDialog) {
    return
  }
  selectedProductType.value = { ...initialProductType }
})
watch(deleteProductTypeDialog, (): void => {
  if (deleteProductTypeDialog) {
    return
  }
  selectedProductType.value = { ...initialProductType }
})
watch(categoryPagination.itemsPerPage, (): void => {
  fetchCategories()
})
watch(productTypePagination.itemsPerPage, (): void => {
  fetchProductTypes()
})
watch(selector, (): void => {
  fetchState.refresh()
})

const fetchState = useAsyncData(async (): Promise<void> => {
  switch (selector.value) {
    case 'categories':
      await handleUpdateCategoryPage(1)
      break
    case 'productTypes':
      await handleUpdateProductTypePage(1)
      break
  }
})

const isLoading = (): boolean => {
  return fetchState?.pending?.value || loading.value
}

const fetchCategories = async (): Promise<void> => {
  try {
    await categoryStore.fetchCategories(categoryPagination.itemsPerPage.value, categoryPagination.offset.value)
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const fetchProductTypes = async (): Promise<void> => {
  try {
    await productTypeStore.fetchProductTypes(productTypePagination.itemsPerPage.value, productTypePagination.offset.value)
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleUpdateCategoryPage = async (page: number): Promise<void> => {
  categoryPagination.updateCurrentPage(page)
  await fetchCategories()
}

const handleUpdateProductTypePage = async (page: number): Promise<void> => {
  productTypePagination.updateCurrentPage(page)
  await fetchProductTypes()
}

const handleSearchCategory = async (name: string): Promise<void> => {
  try {
    const categoryIds = productTypes.value.map((productType: ProductType): string => productType.categoryId)
    await categoryStore.searchCategories(name, categoryIds)
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleClickNewCategory = (): void => {
  createCategoryDialog.value = true
}

const handleClickEditCategory = (categoryId: string): void => {
  const category = categories.value.find((category: Category): boolean => {
    return category.id === categoryId
  })
  if (!category) {
    return
  }
  selectedCategory.value = category
  updateCategoryFormData.value = { ...category }
  updateCategoryDialog.value = true
}

const handleClickDeleteCategory = (categoryId: string): void => {
  const category = categories.value.find((category: Category): boolean => {
    return category.id === categoryId
  })
  if (!category) {
    return
  }
  selectedCategory.value = category
  deleteCategoryDialog.value = true
}

const handleClickNewProductType = (): void => {
  handleSearchCategory('')

  createProductTypeDialog.value = true
}

const handleClickEditProductType = (productTypeId: string): void => {
  const productType = productTypes.value.find((productType: ProductType): boolean => {
    return productType.id === productTypeId
  })
  if (!productType) {
    return
  }
  handleSearchCategory('')

  selectedProductType.value = productType
  updateProductTypeFormData.value = { ...productType }
  updateProductTypeDialog.value = true
}

const handleClickDeleteProductType = (productTypeId: string): void => {
  const productType = productTypes.value.find((productType: ProductType): boolean => {
    return productType.id === productTypeId
  })
  if (!productType) {
    return
  }

  selectedProductType.value = productType
  deleteProductTypeDialog.value = true
}

const handleUploadCreateProductTypeIcon = (files: FileList): void => {
  if (files.length === 0) {
    return
  }

  loading.value = true
  productTypeStore.uploadProductTypeIcon(files[0])
    .then((url: string) => {
      createProductTypeFormData.value.iconUrl = url
    })
    .catch(() => {
      createProductTypeIconUploadStatus.value.error = true
      createProductTypeIconUploadStatus.value.message = 'アップロードに失敗しました。'
    })
    .finally(() => {
      loading.value = false
    })
}

const handleUploadUpdateProductTypeIcon = (files: FileList): void => {
  if (files.length === 0) {
    return
  }

  loading.value = true
  productTypeStore.uploadProductTypeIcon(files[0])
    .then((url: string) => {
      updateProductTypeFormData.value.iconUrl = url
    })
    .catch(() => {
      updateProductTypeIconUploadStatus.value.error = true
      updateProductTypeIconUploadStatus.value.message = 'アップロードに失敗しました。'
    })
    .finally(() => {
      loading.value = false
    })
}

const handleSubmitCreateCategory = async (): Promise<void> => {
  try {
    loading.value = true
    await categoryStore.createCategory(createCategoryFormData.value)
    commonStore.addSnackbar({
      message: 'カテゴリーを追加しました。',
      color: 'info',
    })
    createCategoryDialog.value = false
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    loading.value = false
  }
}

const handleSubmitUpdateCategory = async (): Promise<void> => {
  try {
    loading.value = true
    await categoryStore.updateCategory(selectedCategory.value?.id || '', updateCategoryFormData.value)
    commonStore.addSnackbar({
      message: '変更しました。',
      color: 'info',
    })
    updateCategoryDialog.value = false
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    loading.value = false
  }
}

const handleSubmitDeleteCategory = async (): Promise<void> => {
  try {
    loading.value = true
    await categoryStore.deleteCategory(selectedCategory.value?.id || '')
    commonStore.addSnackbar({
      message: 'カテゴリー削除が完了しました',
      color: 'info',
    })
    deleteCategoryDialog.value = false
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    loading.value = false
  }
}

const handleSubmitCreateProductType = async (categoryId: string): Promise<void> => {
  try {
    loading.value = true
    await productTypeStore.createProductType(categoryId, createProductTypeFormData.value)
    commonStore.addSnackbar({
      message: '品目を追加しました。',
      color: 'info',
    })
    createProductTypeDialog.value = false
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    loading.value = false
  }
}

const handleSubmitUpdateProductType = async (): Promise<void> => {
  try {
    loading.value = true
    await productTypeStore.updateProductType(
      selectedProductType.value.categoryId,
      selectedProductType.value.id,
      createProductTypeFormData.value,
    )
    commonStore.addSnackbar({
      message: '品目の更新が完了しました',
      color: 'info',
    })
    updateProductTypeDialog.value = false
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    loading.value = false
  }
}

const handleSubmitDeleteProductType = async (): Promise<void> => {
  try {
    loading.value = true
    await productTypeStore.deleteProductType(
      selectedProductType.value.categoryId,
      selectedProductType.value.id,
    )
    commonStore.addSnackbar({
      message: '品目の削除が完了しました',
      color: 'info',
    })
    deleteProductTypeDialog.value = false
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    loading.value = false
  }
}

try {
  await fetchState.execute()
}
catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-category-list
    v-model:selected-tab-item="selector"
    v-model:create-category-dialog="createCategoryDialog"
    v-model:update-category-dialog="updateCategoryDialog"
    v-model:delete-category-dialog="deleteCategoryDialog"
    v-model:create-category-form-data="createCategoryFormData"
    v-model:update-category-form-data="updateCategoryFormData"
    v-model:create-product-type-dialog="createProductTypeDialog"
    v-model:update-product-type-dialog="updateProductTypeDialog"
    v-model:delete-product-type-dialog="deleteProductTypeDialog"
    v-model:create-product-type-form-data="createProductTypeFormData"
    v-model:update-product-type-form-data="updateProductTypeFormData"
    :loading="isLoading()"
    :admin-type="adminType"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :category="selectedCategory"
    :categories="categories"
    :category-table-items-per-page="categoryPagination.itemsPerPage.value"
    :category-table-items-total="categoryTotal"
    :product-type="selectedProductType"
    :product-types="productTypes"
    :product-type-table-items-per-page="productTypePagination.itemsPerPage.value"
    :product-type-table-items-total="productTypeTotal"
    :create-product-type-icon-upload-status="createProductTypeIconUploadStatus"
    :update-product-type-icon-upload-status="updateProductTypeIconUploadStatus"
    @click:new-category="handleClickNewCategory"
    @click:edit-category="handleClickEditCategory"
    @click:delete-category="handleClickDeleteCategory"
    @click:new-product-type="handleClickNewProductType"
    @click:edit-product-type="handleClickEditProductType"
    @click:delete-product-type="handleClickDeleteProductType"
    @update:category-page="handleUpdateCategoryPage"
    @update:category-items-per-page="categoryPagination.handleUpdateItemsPerPage"
    @update:product-type-page="handleUpdateProductTypePage"
    @update:product-type-items-per-page="productTypePagination.handleUpdateItemsPerPage"
    @update:create-product-type-icon="handleUploadCreateProductTypeIcon"
    @update:update-product-type-icon="handleUploadUpdateProductTypeIcon"
    @search:category="handleSearchCategory"
    @submit:create-category="handleSubmitCreateCategory"
    @submit:update-category="handleSubmitUpdateCategory"
    @submit:delete-category="handleSubmitDeleteCategory"
    @submit:create-product-type="handleSubmitCreateProductType"
    @submit:update-product-type="handleSubmitUpdateProductType"
    @submit:delete-product-type="handleSubmitDeleteProductType"
  />
</template>
