<template>
  <v-form flat :loading="fetchState.pending">
    <v-data-table :headers="categoryHeaders" :items="categories">
      <template #[`item.category`]="{ item }">
        {{ `${item.name}` }}
      </template>
      <template #[`item.actions`]="{ item }">
        <v-btn outlined color="primary" small @click="handleEdit(item)">
          <v-icon small>mdi-pencil</v-icon>
          編集
        </v-btn>
        <v-btn outlined color="primary" small @click="handleDelete(item)">
          <v-icon small>mdi-delete</v-icon>
          削除
        </v-btn>
      </template>
    </v-data-table>
  </v-form>
</template>

<script lang="ts">
import { useFetch } from '@nuxtjs/composition-api'
import { computed, defineComponent } from '@vue/composition-api'
import { DataTableHeader } from 'vuetify'

import { useCategoryStore } from '~/store/category'
import { CategoriesResponseCategoriesInner } from '~/types/api'

export default defineComponent({
  setup() {
    const categoryStore = useCategoryStore()
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

    const handleDelete = async (
      item: CategoriesResponseCategoriesInner
    ): Promise<void> => {
      try {
        await categoryStore.deleteCategory(item.id)
      } catch (error) {
        console.log(error)
      }
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
      handleDelete,
    }
  },
})
</script>
