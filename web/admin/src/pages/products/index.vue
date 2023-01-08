<template>
  <div>
    <div class="d-flex align-center mb-4">
      <v-card-title>商品管理</v-card-title>
      <v-spacer />
      <v-text-field
        v-model="searchWord"
        append-icon="mdi-magnify"
        label="商品名"
        hide-details
        single-line
      />
    </div>
    <div class="d-flex">
      <v-spacer />
      <v-btn outlined class="mb-4" color="primary" @click="handleClickAddBtn">
        <v-icon left>mdi-plus</v-icon>
        商品登録
      </v-btn>
    </div>

    <v-card :loading="fetchState.pending">
      <v-card-text>
        <v-data-table
          :headers="headers"
          :items="products"
          no-data-text="登録されている商品がありません。"
          :items-per-page.sync="itemsPerPage"
          :server-items-length="totalItems"
          :footer-props="options"
          @update:items-per-page="handleUpdateItemsPerPage"
          @update:page="handleUpdatePage"
        >
          <template #[`item.public`]="{ item }">
            <v-chip :color="item.public ? 'primary' : 'warning'">{{
              item.public ? '公開' : '非公開'
            }}</v-chip>
          </template>
        </v-data-table>
      </v-card-text>
    </v-card>
  </div>
</template>

<script lang="ts">
import {
  defineComponent,
  ref,
  useFetch,
  useRouter,
  computed,
  watch,
} from '@nuxtjs/composition-api'
import { DataTableHeader } from 'vuetify'

import { usePagination } from '~/lib/hooks/'
import { useProductStore } from '~/store/product'

export default defineComponent({
  setup() {
    const router = useRouter()
    const productStore = useProductStore()
    const products = computed(() => productStore.products)
    const totalItems = computed(() => productStore.totalItems)

    const {
      updateCurrentPage,
      itemsPerPage,
      handleUpdateItemsPerPage,
      options,
      offset,
    } = usePagination()

    watch(itemsPerPage, () => {
      productStore.fetchProducts(itemsPerPage.value, 0)
    })

    const handleUpdatePage = async (page: number) => {
      updateCurrentPage(page)
      await productStore.fetchProducts(itemsPerPage.value, offset.value)
    }

    const { fetchState } = useFetch(async () => {
      try {
        await productStore.fetchProducts(itemsPerPage.value, offset.value)
      } catch (error) {
        console.log(error)
      }
    })

    const searchWord = ref<string>('')

    const handleClickAddBtn = () => {
      router.push('/products/add')
    }

    const headers: DataTableHeader[] = [
      {
        text: '商品名',
        value: 'name',
      },
      {
        text: 'ステータス',
        value: 'public',
      },
      {
        text: '種類',
        value: 'type',
      },
      {
        text: '価格',
        value: 'price',
      },
      {
        text: '在庫',
        value: 'inventory',
      },
      {
        text: 'ジャンル',
        value: 'categoryName',
      },
      {
        text: '品目',
        value: 'productTypeName',
      },
      {
        text: '農園名',
        value: 'storeName',
      },
    ]

    return {
      fetchState,
      headers,
      searchWord,
      handleClickAddBtn,
      products,
      totalItems,
      itemsPerPage,
      handleUpdateItemsPerPage,
      handleUpdatePage,
      options,
    }
  },
})
</script>
