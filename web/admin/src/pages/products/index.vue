<script lang="ts" setup>
import { useAlert, usePagination } from '~/lib/hooks'
import { useProductStore } from '~/store'

const router = useRouter()
const productStore = useProductStore()
const pagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const products = computed(() => {
  return productStore.products
})
const totalItems = computed(() => {
  return productStore.totalItems
})

watch(pagination.itemsPerPage, () => {
  fetchState.refresh()
})

const fetchState = useAsyncData(async () => {
  await fetchProducts()
})

const fetchProducts = async (): Promise<void> => {
  try {
    await productStore.fetchProducts(pagination.itemsPerPage.value, pagination.offset.value)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const isLoading = (): boolean => {
  return fetchState?.pending?.value || false
}

const handleUpdatePage = async (page: number): Promise<void> => {
  pagination.updateCurrentPage(page)
  await fetchState.refresh()
}

const handleClickShow = (productId: string): void => {
  router.push(`/products/${productId}`)
}

const handleClickNew = () => {
  router.push('/products/add')
}

const handleClickDelete = async (productId: string): Promise<void> => {
  try {
    await productStore.deleteProduct(productId)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-product-list
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :products="products"
    :table-items-total="totalItems"
    :table-items-per-page="pagination.itemsPerPage.value"
    @click:update-page="handleUpdatePage"
    @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
    @click:show="handleClickShow"
    @click:new="handleClickNew"
    @click:delete="handleClickDelete"
  />
</template>
