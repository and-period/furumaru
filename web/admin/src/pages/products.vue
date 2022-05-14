<template>
  <div>
    <p>商品一覧</p>
    <v-btn block outlined class="mb-4" @click="handleClickAddBtn">追加</v-btn>
    <v-data-table
      v-model="selectedProducts"
      :headers="headers"
      :items="products"
      show-select
    />
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from '@vue/composition-api'

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
    const handleClickAddBtn = () => {
      console.log('NOT IMPLEMENTED')
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
      handleClickAddBtn,
      products,
      selectedProducts,
    }
  },
})
</script>
