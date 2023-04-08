<script lang="ts" setup>
import { DataTableHeader } from 'vuetify'

import { usePagination } from '~/lib/hooks/'
import { useProductStore } from '~/store'
import { ProductsResponseProductsInner } from '~/types/api'

const router = useRouter()
const productStore = useProductStore()
const products = computed(() => productStore.products)
const totalItems = computed(() => productStore.totalItems)

const {
  updateCurrentPage,
  itemsPerPage,
  handleUpdateItemsPerPage,
  options,
  offset
} = usePagination()

watch(itemsPerPage, () => {
  productStore.fetchProducts(itemsPerPage.value, 0)
})

const handleUpdatePage = async (page: number) => {
  updateCurrentPage(page)
  await productStore.fetchProducts(itemsPerPage.value, offset.value)
}

const fetchState = useAsyncData(async () => {
  try {
    await productStore.fetchProducts(itemsPerPage.value, offset.value)
  } catch (error) {
    console.log(error)
  }
})

const searchWord = ref<string>('')

const handleRowClick = (
  _: any,
  { item }: { item: ProductsResponseProductsInner }
): void => {
  router.push(`/products/${item.id}`)
}

const handleClickAddBtn = () => {
  router.push('/products/add')
}

const headers: DataTableHeader[] = [
  {
    text: '',
    value: 'media'
  },
  {
    text: '商品名',
    value: 'name'
  },
  {
    text: 'ステータス',
    value: 'public'
  },
  {
    text: '種類',
    value: 'type'
  },
  {
    text: '価格',
    value: 'price'
  },
  {
    text: '在庫',
    value: 'inventory'
  },
  {
    text: 'ジャンル',
    value: 'categoryName'
  },
  {
    text: '品目',
    value: 'productTypeName'
  },
  {
    text: '農園名',
    value: 'storeName'
  }
]
</script>

<template>
  <div>
    <v-card-title>
      商品管理
      <v-spacer />
      <v-btn outlined color="primary" @click="handleClickAddBtn">
        <v-icon left>
          mdi-plus
        </v-icon>
        商品登録
      </v-btn>
    </v-card-title>

    <v-card :loading="fetchState.pending">
      <v-card-text>
        <div class="d-flex align-center mb-4">
          <v-spacer />
          <v-text-field
            v-model="searchWord"
            append-icon="mdi-magnify"
            label="商品名"
            hide-details
            single-line
          />
        </div>

        <v-data-table
          v-model:items-per-page="itemsPerPage"
          :headers="headers"
          :items="products"
          no-data-text="登録されている商品がありません。"
          :server-items-length="totalItems"
          :footer-props="options"
          @update:items-per-page="handleUpdateItemsPerPage"
          @update:page="handleUpdatePage"
          @click:row="handleRowClick"
        >
          <template #[`item.media`]="{ item }">
            <v-avatar tile>
              <v-img contain :src="item.media.find((m) => m.isThumbnail).url" />
            </v-avatar>
          </template>
          <template #[`item.public`]="{ item }">
            <v-chip :color="item.public ? 'primary' : 'warning'">
              {{ item.public ? '公開' : '非公開' }}
            </v-chip>
          </template>
        </v-data-table>
      </v-card-text>
    </v-card>
  </div>
</template>
