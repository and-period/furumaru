<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useAlert, usePagination } from '~/lib/hooks'
import { useCategoryStore, useProductTypeStore } from '~/store'
import { CreateCategoryRequest, CreateProductTypeRequest } from '~/types/api'

const categoryStore = useCategoryStore()
const productTypeStore = useProductTypeStore()
const categoryPagination = usePagination()
const productTypePagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const { categories, total: categoryTotal } = storeToRefs(categoryStore)
const { productTypes, totalItems: productTypeTotal } = storeToRefs(productTypeStore)

const loading = ref<boolean>(false)
const selector = ref<string>('categories')
const selectedCategoryId = ref<string>('')
const categoryDialog = ref<boolean>(false)
const productTypeDialog = ref<boolean>(false)
const categoryFormData = ref<CreateCategoryRequest>({
  name: ''
})
const productTypeFormData = ref<CreateProductTypeRequest>({
  name: '',
  iconUrl: ''
})

const fetchState = useAsyncData(async (): Promise<void> => {
  await Promise.all([fetchCategories(), fetchProductTypes()])
})

watch(categoryPagination.itemsPerPage, (): void => {
  fetchCategories()
})
watch(productTypePagination.itemsPerPage, (): void => {
  fetchProductTypes()
})
watch(selector, async (): Promise<void> => {
  categoryPagination.updateCurrentPage(1)
  productTypePagination.updateCurrentPage(1)
  await fetchState.refresh()
})

/*
 * category methods
 */
const categoryState = useAsyncData(async (): Promise<void> => {
  await fetchCategories()
})

const fetchCategories = async (): Promise<void> => {
  try {
    await categoryStore.fetchCategories(categoryPagination.itemsPerPage.value, categoryPagination.offset.value)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleCategoryMorePage = async (): Promise<void> => {
  try {
    categoryPagination.updateCurrentPage(categoryPagination.currentPage.value + 1)
    await categoryStore.moreCategories(categoryPagination.itemsPerPage.value, categoryPagination.offset.value)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleCreateCategory = async (): Promise<void> => {
  try {
    loading.value = true
    await categoryStore.createCategory(categoryFormData.value)
    categoryFormData.value.name = ''
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  } finally {
    categoryDialog.value = false
    loading.value = false
  }
}

const handleUpdateCategoryPage = async (page: number): Promise<void> => {
  categoryPagination.updateCurrentPage(page)
  await fetchCategories()
}

/*
 * productType methods
 */
const productTypeState = useAsyncData(async (): Promise<void> => {
  await fetchProductTypes()
})

const fetchProductTypes = async (): Promise<void> => {
  try {
    await productTypeStore.fetchProductTypes(productTypePagination.itemsPerPage.value, productTypePagination.offset.value)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleCreateProductType = async (): Promise<void> => {
  try {
    loading.value = true
    await productTypeStore.createProductType(selectedCategoryId.value, productTypeFormData.value)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  } finally {
    productTypeDialog.value = false
    loading.value = false
  }
}

const handleUploadProductTypeIcon = async (files: FileList): Promise<void> => {
  if (files.length === 0) {
    return
  }

  try {
    const res = await productTypeStore.uploadProductTypeIcon(files[0])
    productTypeFormData.value.iconUrl = res.url
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleUpdateProductTypePage = async (page: number): Promise<void> => {
  productTypePagination.updateCurrentPage(page)
  await fetchProductTypes()
}

/**
 * common methods
 */
const isLoading = (): boolean => {
  return fetchState?.pending?.value || loading.value
}

const handleClickTabs = async (key: string) => {
  switch (key) {
    case 'categories':
      categoryPagination.updateCurrentPage(1)
      await categoryState.execute()
      break
    case 'productTypes':
      categoryPagination.updateCurrentPage(1)
      productTypePagination.updateCurrentPage(1)
      await fetchState.refresh()
      break
  }
}

try {
  await Promise.all([categoryState.execute(), productTypeState.execute()])
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-category-list
    v-model:category-dialog="categoryDialog"
    v-model:product-type-dialog="productTypeDialog"
    v-model:category-form-data="categoryFormData"
    v-model:product-type-form-data="productTypeFormData"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :categories="categories"
    :category-table-items-per-page="categoryPagination.itemsPerPage.value"
    :category-table-items-total="categoryTotal"
    :product-types="productTypes"
    :product-type-table-items-per-page="productTypePagination.itemsPerPage.value"
    :product-type-table-items-total="productTypeTotal"
    @click:category-update-page="handleUpdateCategoryPage"
    @click:category-update-items-per-page="categoryPagination.handleUpdateItemsPerPage"
    @click:category-more-page="handleCategoryMorePage"
    @click:product-type-update-page="handleUpdateProductTypePage"
    @click:product-type-update-items-per-page="productTypePagination.handleUpdateItemsPerPage"
    @update:tab="handleClickTabs"
    @update:product-type-upload-icon="handleUploadProductTypeIcon"
    @submit:category="handleCreateCategory"
    @submit:product-type="handleCreateProductType"
  />
</template>
