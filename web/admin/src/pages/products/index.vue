<script lang="ts" setup>
import { storeToRefs } from 'pinia'

import { useAlert, usePagination } from '~/lib/hooks'
import {
  useAuthStore,
  useCategoryStore,
  useProducerStore,
  useProductStore,
  useProductTagStore,
  useProductTypeStore,
} from '~/store'

const router = useRouter()
const authStore = useAuthStore()
const categoryStore = useCategoryStore()
const producerStore = useProducerStore()
const productStore = useProductStore()
const productTagStore = useProductTagStore()
const productTypeStore = useProductTypeStore()
const pagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const { adminType } = storeToRefs(authStore)
const { categories } = storeToRefs(categoryStore)
const { producers } = storeToRefs(producerStore)
const { products, totalItems } = storeToRefs(productStore)
const { productTags } = storeToRefs(productTagStore)
const { productTypes } = storeToRefs(productTypeStore)

const loading = ref<boolean>(false)
const selectedItemId = ref<string>('')

const fetchState = useAsyncData('products', async (): Promise<void> => {
  await fetchProducts()
})

watch(pagination.itemsPerPage, (): void => {
  fetchState.refresh()
})

const fetchProducts = async (): Promise<void> => {
  try {
    await productStore.fetchProducts(
      pagination.itemsPerPage.value,
      pagination.offset.value,
    )
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const isLoading = (): boolean => {
  return fetchState?.pending?.value || loading.value
}

const handleUpdatePage = async (page: number): Promise<void> => {
  pagination.updateCurrentPage(page)
  await fetchState.refresh()
}

const handleClickShow = (productId: string): void => {
  router.push(`/products/${productId}`)
}

const handleClickNew = (): void => {
  router.push('/products/new')
}

const handleClickCopyItem = (): void => {
  if (selectedItemId.value !== '') {
    router.push(`/products/new?from=${selectedItemId.value}`)
  }
}

const handleClickDelete = async (productId: string): Promise<void> => {
  try {
    loading.value = true
    await productStore.deleteProduct(productId)
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
  <templates-product-list
    v-model:selected-item-id="selectedItemId"
    :loading="isLoading()"
    :admin-type="adminType"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :categories="categories"
    :producers="producers"
    :products="products"
    :product-tags="productTags"
    :product-types="productTypes"
    :table-items-per-page="pagination.itemsPerPage.value"
    :table-items-total="totalItems"
    @click:show="handleClickShow"
    @click:new="handleClickNew"
    @click:delete="handleClickDelete"
    @click:update-page="handleUpdatePage"
    @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
    @click:copy-item="handleClickCopyItem"
  />
</template>
