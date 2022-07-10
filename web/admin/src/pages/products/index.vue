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
      <v-btn outlined class="mb-4" @click="handleClickAddBtn">
        <v-icon left>mdi-plus</v-icon>
        商品登録
      </v-btn>
    </div>

    <v-data-table
      v-model="selectedProducts"
      :headers="headers"
      :items="products"
      show-select
    />
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, useRouter } from '@nuxtjs/composition-api'

interface IProduct {
  id: string
  name: string
  description: string
  public: 0 | 1
  type: string
  price: number
}

interface DataTableHeader {
  text: string
  value: string
  sortable?: boolean
}

export default defineComponent({
  setup() {
    const router = useRouter()

    const searchWord = ref<string>('')

    const handleClickAddBtn = () => {
      router.push('/products/add')
    }

    const headers: DataTableHeader[] = [
      {
        text: 'id',
        value: 'id',
      },
      {
        text: '商品名',
        value: 'name',
      },
      {
        text: '公開',
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
    ]

    const products = ref<IProduct[]>([
      {
        id: '1',
        name: 'みかん',
        description: '',
        public: 1,
        type: '果物',
        price: 1000,
      },
    ])

    const selectedProducts = ref<IProduct[]>([])

    return {
      headers,
      searchWord,
      handleClickAddBtn,
      products,
      selectedProducts,
    }
  },
})
</script>
