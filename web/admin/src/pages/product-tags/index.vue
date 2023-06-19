<script lang="ts" setup>
import { storeToRefs } from 'pinia';
import { VDataTable } from 'vuetify/lib/labs/components';
import { useAlert, usePagination } from '~/lib/hooks';
import { useProductTagStore } from '~/store';
import { CreateProductTagRequest, UpdateProductTagRequest } from '~/types/api';

const productTagStore = useProductTagStore()
const pagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const { productTags, total } = storeToRefs(productTagStore)

const loading = ref<boolean>(false)
const sortBy = ref<VDataTable['sortBy']>([])
const newDialog = ref<boolean>(false)
const editDialog = ref<boolean>(false)
const deleteDialog = ref<boolean>(false)
const newFormData = ref<CreateProductTagRequest>({
  name: '',
})
const editFormData = ref<UpdateProductTagRequest>({
  name: '',
})

const fetchState = useAsyncData(async (): Promise<void> => {
  await fetchProductTags()
})

watch(pagination.itemsPerPage, (): void => {
  fetchProductTags()
})
watch(sortBy, (): void => {
  fetchState.refresh()
})

const fetchProductTags = async (): Promise<void> => {
  try {
    const orders: string[] = sortBy.value.map((item) => {
      switch (item.order) {
        case 'asc':
          return item.key
        case 'desc':
          return `-${item.key}`
        default:
          return item.order ? item.key : `-${item.key}`
      }
    }) || []

    await productTagStore.fetchProductTags(pagination.itemsPerPage.value, pagination.offset.value, orders)
  } catch (err) {
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
  await fetchProductTags()
}

const handleCreate = async (): Promise<void> => {
  try {
    loading.value = true
    await productTagStore.createProductTag(newFormData.value)
    fetchState.refresh()
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  } finally {
    newDialog.value = false
    loading.value = false
  }
}

const handleUpdate = async (productTagId: string): Promise<void> => {
  try {
    loading.value = true
    await productTagStore.updateProductTag(productTagId, editFormData.value)
    fetchState.refresh()
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  } finally {
    editDialog.value = false
    loading.value = false
  }
}

const handleDelete = async (productTagId: string): Promise<void> => {
  try {
    loading.value = true
    await productTagStore.deleteProductTag(productTagId)
    fetchState.refresh()
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  } finally {
    deleteDialog.value = false
    loading.value = false
  }
}

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-product-tag-list
    v-model:new-form-data="newFormData"
    v-model:edit-form-data="editFormData"
    v-model:new-dialog="newDialog"
    v-model:edit-dialog="editDialog"
    v-model:delete-dialog="deleteDialog"
    v-model:sort-by="sortBy"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :product-tags="productTags"
    :table-items-per-page="pagination.itemsPerPage.value"
    :table-items-total="total"
    @click:update-page="handleUpdatePage"
    @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
    @update:sort-by="fetchState.refresh"
    @submit:create="handleCreate"
    @submit:update="handleUpdate"
    @submit:delete="handleDelete"
  />
</template>
