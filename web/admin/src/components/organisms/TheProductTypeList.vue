<template>
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

import { useProductTypeStore } from '~/store/product-type'

export default defineComponent({
  setup() {
    const productTypeStore = useProductTypeStore()
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
    }
  },
})
</script>
