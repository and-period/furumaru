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
import { computed, defineComponent, useFetch } from '@nuxtjs/composition-api'
import { DataTableHeader } from 'vuetify'

import { useCategoryStore } from '~/store/category'

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
        sortable: false,
      },
    ]

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
    }
  },
})
</script>
