<template>
  <div>
    <v-data-table
      :headers="productTypeHeaders"
      :items="productTypes"
      :loading="loading"
      :server-items-length="totalItems"
      :footer-props="tableFooterProps"
      @update:items-per-page="handleUpdateItemsPerPage"
      @update:page="handleUpdatePage"
    >
      <template #[`item.category`]="{ item }">
        {{ `${item.categoryName}` }}
      </template>
      <template #[`item.productType`]="{ item }">
        {{ `${item.name}` }}
      </template>
      <template #[`item.actions`]="{ item }">
        <v-btn outlined color="primary" small @click="openEditDialog(item)">
          <v-icon small>mdi-pencil</v-icon>
          編集
        </v-btn>
        <v-btn outlined color="primary" small @click="openDeleteDialog(item)">
          <v-icon small>mdi-delete</v-icon>
          削除
        </v-btn>
      </template>
    </v-data-table>

    <v-dialog v-model="editDialog" width="500">
      <v-card>
        <v-card-title class="primaryLight">品目編集</v-card-title>
        <v-card-text class="mt-4">
          <v-autocomplete
            v-model="selectedCategoryId"
            :items="categories"
            item-text="name"
            item-value="id"
            label="カテゴリー"
          >
            <template #append-item>
              <div class="pa-2">
                <v-btn
                  outlined
                  block
                  color="primary"
                  @click="handleMoreCategoryItems"
                >
                  <v-icon>mdi-plus</v-icon>
                  さらに読み込む
                </v-btn>
              </div>
            </template>
          </v-autocomplete>
          <v-spacer />
          <v-text-field
            v-model="editFormData.name"
            maxlength="32"
            label="品目"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="accentDarken" text @click="hideEditDialog">
            キャンセル
          </v-btn>
          <v-btn color="primary" outlined @click="handleEdit"> 編集 </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="deleteDialog" width="500">
      <v-card>
        <v-card-title class="text-h7">
          {{ selectedName }}を本当に削除しますか？
        </v-card-title>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="accentDarken" text @click="hideDeleteDialog">
            キャンセル
          </v-btn>
          <v-btn color="primary" outlined @click="handleDelete"> 削除 </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script lang="ts">
import { ref, computed, defineComponent, reactive } from '@vue/composition-api'
import { DataTableHeader } from 'vuetify'

import { useProductTypeStore } from '~/store/product-type'
import {
  ProductTypesResponseProductTypesInner,
  UpdateProductTypeRequest,
} from '~/types/api'

export default defineComponent({
  props: {
    loading: {
      type: Boolean,
      default: false,
    },
    tableFooterProps: {
      type: Object,
      default: () => {},
    },
    categories: {
      type: Array,
      default: () => [],
    },
  },
  setup(_, { emit }) {
    const productTypeStore = useProductTypeStore()
    const deleteDialog = ref<boolean>(false)
    const editDialog = ref<boolean>(false)
    const selectedCategoryId = ref<string>('')
    const selectedItemId = ref<string>('')
    const selectedName = ref<string>('')

    const editFormData = reactive<UpdateProductTypeRequest>({
      name: '',
      iconUrl: '',
    })

    const productTypes = computed(() => {
      return productTypeStore.productTypes
    })

    const totalItems = computed(() => {
      return productTypeStore.totalItems
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

    const handleUpdateItemsPerPage = (page: number) => {
      emit('update:items-per-page', page)
    }

    const handleUpdatePage = (page: number) => {
      emit('update:page', page)
    }

    const handleMoreCategoryItems = () => {
      emit('click:more-item')
    }

    const openEditDialog = (item: ProductTypesResponseProductTypesInner) => {
      editDialog.value = true
      selectedCategoryId.value = item.categoryId
      selectedItemId.value = item.id
      editFormData.name = item.name
      editFormData.iconUrl = item.iconUrl
    }

    const hideEditDialog = () => {
      editDialog.value = false
    }

    const handleEdit = async () => {
      try {
        await productTypeStore.editProductType(
          selectedCategoryId.value,
          selectedItemId.value,
          editFormData
        )
        editDialog.value = false
      } catch (error) {
        console.log(error)
      }
    }

    const openDeleteDialog = (
      item: ProductTypesResponseProductTypesInner
    ): void => {
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

    const hideDeleteDialog = (): void => {
      deleteDialog.value = false
    }

    return {
      productTypeHeaders,
      productTypes,
      selectedName,
      selectedCategoryId,
      deleteDialog,
      editFormData,
      editDialog,
      openEditDialog,
      hideEditDialog,
      openDeleteDialog,
      hideDeleteDialog,
      handleEdit,
      handleDelete,
      totalItems,
      handleUpdateItemsPerPage,
      handleUpdatePage,
      handleMoreCategoryItems,
    }
  },
})
</script>
