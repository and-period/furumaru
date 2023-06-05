<script lang="ts" setup>
import { useAlert, usePagination } from '~/lib/hooks'
import { useCategoryStore, useProductTypeStore } from '~/store'
import { CreateCategoryRequest, CreateProductTypeRequest } from '~/types/api'

const categoryStore = useCategoryStore()
const productTypeStore = useProductTypeStore()
const categoryPagination = usePagination()
const productTypePagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const selector = ref<string>('categories')
const selectedCategoryId = ref<string>('')

const dialog = reactive({
  category: false,
  productType: false
})
const categoryFormData = reactive<CreateCategoryRequest>({
  name: ''
})
const productTypeFormData = reactive<CreateProductTypeRequest>({
  name: '',
  iconUrl: ''
})

const categories = computed(() => {
  return categoryStore.categories
})
const categoryTotal = computed(() => {
  return categoryStore.total
})
const productTypes = computed(() => {
  return productTypeStore.productTypes
})
const productTypeTotal = computed(() => {
  return productTypeStore.totalItems
})

watch(categoryPagination.itemsPerPage, () => {
  fetchCategories()
})
watch(productTypePagination.itemsPerPage, () => {
  fetchProductTypes()
})
watch(selector, async () => {
  categoryPagination.updateCurrentPage(1)
  productTypePagination.updateCurrentPage(1)
  await Promise.all([categoryState.execute(), productTypeState.execute()])
})

/*
 * category methods
 */
const categoryState = useAsyncData(async () => {
  await fetchCategories()
})

const fetchCategories = async () => {
  try {
    await categoryStore.fetchCategories(categoryPagination.itemsPerPage.value, categoryPagination.offset.value)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleCategoryMorePage = async () => {
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
    await categoryStore.createCategory(categoryFormData)
    categoryFormData.name = ''
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  } finally {
    dialog.category = false
  }
}

const handleUpdateCategoryPage = async (page: number) => {
  categoryPagination.updateCurrentPage(page)
  await fetchCategories()
}

/*
 * productType methods
 */
const productTypeState = useAsyncData(async () => {
  await fetchProductTypes()
})

const fetchProductTypes = async () => {
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
    await productTypeStore.createProductType(selectedCategoryId.value, productTypeFormData)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  } finally {
    dialog.productType = false
  }
}

const handleUploadProductTypeIcon = async (files: FileList) => {
  if (files.length === 0) {
    return
  }
  try {
    const res = await productTypeStore.uploadProductTypeIcon(files[0])
    productTypeFormData.iconUrl = res.url
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleUpdateProductTypePage = async (page: number) => {
  productTypePagination.updateCurrentPage(page)
  await fetchProductTypes()
}

/**
 * common methods
 */
const isLoading = (): boolean => {
  return categoryState.pending.value && productTypeState.pending.value
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
      await Promise.all([categoryState.refresh(), productTypeState.refresh()])
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
    v-model:category-dialog="dialog.category"
    v-model:product-type-dialog="dialog.productType"
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
