<template>
  <div>
    <v-form flat :loading="fetchState.pending">
      <v-data-table :headers="productTypeHeaders" :items="productTypes">
        <template #[`item.category`]="{ item }">
          {{ `${item.categoryName}` }}
        </template>
        <template #[`item.productType`]="{ item }">
          {{ `${item.name}` }}
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn outlined color="primary" small @click="handleEdit(item)">
            <v-icon small>mdi-pencil</v-icon>
            編集
          </v-btn>
          <v-btn outlined color="primary" small @click="openDialog(item)">
            <v-icon small>mdi-delete</v-icon>
            削除
          </v-btn>
        </template>
      </v-data-table>
    </v-form>
    <v-dialog v-model="deleteDialog" width="500">
      <v-card>
        <v-card-title class="text-h7">
          {{ selectedName }}を本当に削除しますか？
        </v-card-title>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="accentDarken" text @click="deleteCancel">
            キャンセル
          </v-btn>
          <v-btn color="primary" outlined @click="handleDelete"> 削除 </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script lang="ts">
import { ref, useFetch } from '@nuxtjs/composition-api'
import { computed, defineComponent } from '@vue/composition-api'
import { DataTableHeader } from 'vuetify'

import { useProductTypeStore } from '~/store/product-type'
import { ProductTypesResponseProductTypesInner } from '~/types/api'

export default defineComponent({
  setup() {
    const productTypeStore = useProductTypeStore()
    const deleteDialog = ref<boolean>(false)
    const selectedCategoryId = ref<string>('')
    const selectedItemId = ref<string>('')
    const selectedName = ref<string>('')
    const productTypes = computed(() => {
      return productTypeStore.productTypes
    })
    const productTypeHeaders: DataTableHeader[] = [
      {
        text: 'カテゴリー',
        value: 'category',
      },
      {
        text: '品目',
        value: 'productType',
      },
      {
        text: 'Actions',
        value: 'actions',
        width: 200,
        align: 'end',
        sortable: false,
      },
    ]

    const openDialog = (item: ProductTypesResponseProductTypesInner): void => {
      selectedCategoryId.value = item.categoryId
      selectedItemId.value = item.id
      selectedName.value = item.name
      deleteDialog.value = true
    }

    const handleDelete = async (): Promise<void> => {
      try {
        await productTypeStore.deleteProductType(
          selectedCategoryId.value,
          selectedItemId.value
        )
      } catch (error) {
        console.log(error)
      }
      deleteDialog.value = false
    }

    const deleteCancel = (): void => {
      deleteDialog.value = false
    }

    const { fetchState } = useFetch(async () => {
      try {
        await productTypeStore.fetchProductTypes()
      } catch (err) {
        console.log(err)
      }
    })

    return {
      productTypeHeaders,
      fetchState,
      productTypes,
      deleteDialog,
      selectedName,
      openDialog,
      handleDelete,
      deleteCancel,
    }
  },
})
</script>
