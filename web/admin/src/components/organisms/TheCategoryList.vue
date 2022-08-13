<template>
  <div>
    <v-form flat :loading="fetchState.pending">
      <v-data-table :headers="categoryHeaders" :items="categories">
        <template #[`item.category`]="{ item }">
          {{ `${item.name}` }}
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn outlined color="primary" small @click="openUpdateDialog(item)">
            <v-icon small>mdi-pencil</v-icon>
            編集
          </v-btn>
          <v-btn outlined color="primary" small @click="openDeleteDialog(item)">
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
    <v-dialog v-model="updateDialog" width="500">
      <v-card>
        <v-card-title class="text-h6 primaryLight">
          カテゴリー編集
        </v-card-title>
        <v-text-field
          v-model="categoryFormData.name"
          class="mx-4"
          maxlength="32"
          label="カテゴリー"
        />
        <v-divider></v-divider>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="accentDarken" text @click="updateCancel">
            キャンセル
          </v-btn>
          <v-btn color="primary" outlined @click="handleUpdate"> 編集 </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script lang="ts">
import { reactive, ref, useFetch } from '@nuxtjs/composition-api'
import { computed, defineComponent } from '@vue/composition-api'
import { DataTableHeader } from 'vuetify'

import { useCategoryStore } from '~/store/category'
import {
  CategoriesResponseCategoriesInner,
  UpdateCategoryRequest,
} from '~/types/api'

export default defineComponent({
  setup() {
    const categoryStore = useCategoryStore()
    const deleteDialog = ref<boolean>(false)
    const updateDialog = ref<boolean>(false)
    const selectedItem = ref<string>('')
    const selectedName = ref<string>('')
    const categoryId = ref<string>('')
    const categoryFormData = reactive<UpdateCategoryRequest>({
      name: '',
    })

    const categories = computed(() => {
      return categoryStore.categories
    })
    const categoryHeaders: DataTableHeader[] = [
      {
        text: 'カテゴリー',
        value: 'category',
      },
      {
        text: 'Actions',
        value: 'actions',
        width: 200,
        align: 'end',
        sortable: false,
      },
    ]

    const deleteCancel = (): void => {
      deleteDialog.value = false
    }

    const updateCancel = (): void => {
      updateDialog.value = false
    }

    const openUpdateDialog = (
      item: CategoriesResponseCategoriesInner
    ): void => {
      categoryFormData.name = item.name
      categoryId.value = item.id
      updateDialog.value = true
    }

    const openDeleteDialog = (
      item: CategoriesResponseCategoriesInner
    ): void => {
      selectedItem.value = item.id
      selectedName.value = item.name
      deleteDialog.value = true
    }

    const handleUpdate = async (): Promise<void> => {
      try {
        await categoryStore.updateCategory(categoryId.value, categoryFormData)
        updateDialog.value = false
      } catch (error) {
        console.log(error)
      }
    }

    const handleDelete = async (): Promise<void> => {
      try {
        await categoryStore.deleteCategory(selectedItem.value)
      } catch (error) {
        console.log(error)
      }
      deleteDialog.value = false
    }

    const { fetchState } = useFetch(async () => {
      try {
        await categoryStore.fetchCategories()
      } catch (err) {
        console.log(err)
      }
    })

    return {
      categoryHeaders,
      fetchState,
      categories,
      deleteDialog,
      updateDialog,
      selectedName,
      categoryFormData,
      openUpdateDialog,
      openDeleteDialog,
      updateCancel,
      deleteCancel,
      handleUpdate,
      handleDelete,
    }
  },
})
</script>
