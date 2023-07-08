<script lang="ts" setup>
import { mdiDelete, mdiPencil } from '@mdi/js'
import { VDataTable } from 'vuetify/lib/labs/components.mjs'

import { useCategoryStore } from '~/store'
import {
  CategoriesResponseCategoriesInner,
  UpdateCategoryRequest
} from '~/types/api'

const props = defineProps({
  categories: {
    type: Array<CategoriesResponseCategoriesInner>,
    default: () => []
  },
  tableItemsPerPage: {
    type: Number,
    default: 20
  },
  tableItemsLength: {
    type: Number,
    default: 0
  },
  tableFooterOptions: {
    type: Object,
    default: () => {}
  }
})

const emit = defineEmits<{
  (e: 'update:items-per-page', page: number): void
  (e: 'update:page', page: number): void
}>()

const categoryStore = useCategoryStore()
const deleteDialog = ref<boolean>(false)
const editDialog = ref<boolean>(false)
const selectedItem = ref<string>('')
const selectedName = ref<string>('')
const categoryId = ref<string>('')
const categoryFormData = reactive<UpdateCategoryRequest>({
  name: ''
})

const categoryHeaders: VDataTable['headers'] = [
  {
    title: 'カテゴリー',
    key: 'name'
  },
  {
    title: 'Actions',
    key: 'actions',
    width: 200,
    align: 'end',
    sortable: false
  }
]

const handleUpdateItemsPerPage = (page: number) => {
  emit('update:items-per-page', page)
}

const handleUpdatePage = (page: number) => {
  emit('update:page', page)
}

const deleteCancel = (): void => {
  deleteDialog.value = false
}

const editCancel = (): void => {
  editDialog.value = false
}

const openEditDialog = (item: CategoriesResponseCategoriesInner): void => {
  categoryFormData.name = item.name
  categoryId.value = item.id
  editDialog.value = true
}

const openDeleteDialog = (item: CategoriesResponseCategoriesInner): void => {
  selectedItem.value = item.id
  selectedName.value = item.name
  deleteDialog.value = true
}

const handleEdit = async (): Promise<void> => {
  try {
    await categoryStore.updateCategory(categoryId.value, categoryFormData)
    editDialog.value = false
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
</script>

<template>
  <div>
    <v-data-table-server
      :headers="categoryHeaders"
      :items="categories"
      :items-per-page="tableItemsPerPage"
      :items-length="tableItemsLength"
      :footer-props="tableFooterOptions"
      @update:page="handleUpdatePage"
      @update:items-per-page="handleUpdateItemsPerPage"
    >
      <template #[`item.actions`]="{ item }">
        <v-btn class="mr-2" variant="outlined" color="primary" size="small" @click="openEditDialog(item.raw)">
          <v-icon size="small" :icon="mdiPencil" />
          編集
        </v-btn>
        <v-btn variant="outlined" color="primary" size="small" @click="openDeleteDialog(item.raw)">
          <v-icon size="small" :icon="mdiDelete" />
          削除
        </v-btn>
      </template>
    </v-data-table-server>

    <v-dialog v-model="deleteDialog" width="500">
      <v-card>
        <v-card-title class="text-h7">
          {{ selectedName }}を本当に削除しますか？
        </v-card-title>
        <v-card-actions>
          <v-spacer />
          <v-btn color="error" variant="text" @click="deleteCancel">
            キャンセル
          </v-btn>
          <v-btn color="primary" variant="outlined" @click="handleDelete">
            削除
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="editDialog" width="500">
      <v-card>
        <v-card-title class="primaryLight">
          カテゴリー編集
        </v-card-title>
        <v-card-text class="mt-6">
          <v-text-field
            v-model="categoryFormData.name"
            maxlength="32"
            label="カテゴリー"
          />
        </v-card-text>
        <v-divider />

        <v-card-actions>
          <v-spacer />
          <v-btn color="error" variant="text" @click="editCancel">
            キャンセル
          </v-btn>
          <v-btn color="primary" variant="outlined" @click="handleEdit">
            編集
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>
