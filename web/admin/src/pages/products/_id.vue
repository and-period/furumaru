<template>
  <div>
    <the-product-update-form-page
      :loading="fetchState.pending"
      :form-data="formData"
    />
  </div>
</template>

<script lang="ts">
import {
  defineComponent,
  reactive,
  useFetch,
  useRoute,
} from '@nuxtjs/composition-api'

import { useCategoryStore } from '~/store/category'
import { useProducerStore } from '~/store/producer'
import { useProductStore } from '~/store/product'
import { useProductTypeStore } from '~/store/product-type'
import { UpdateProductRequest } from '~/types/api'

export default defineComponent({
  setup() {
    const route = useRoute()
    const id = route.value.params.id

    const formData = reactive<UpdateProductRequest>({
      name: '',
      description: '',
      producerId: '',
      productTypeId: '',
      public: false,
      inventory: 0,
      weight: 0,
      itemUnit: '',
      itemDescription: '',
      media: [],
      price: 0,
      deliveryType: 0,
      box60Rate: 0,
      box80Rate: 0,
      box100Rate: 0,
      originPrefecture: '',
      originCity: '',
    })

    const productStore = useProductStore()
    const productTypeStore = useProductTypeStore()
    const categoryStore = useCategoryStore()
    const producerStore = useProducerStore()

    const { fetchState } = useFetch(async () => {
      try {
        Promise.all([
          productTypeStore.fetchProductTypes(),
          categoryStore.fetchCategories(),
          producerStore.fetchProducers(20, 0, ''),
        ])
        const data = await productStore.getProduct(id)
        Object.assign(formData, data)
      } catch (error) {
        console.log(error)
      }
    })

    return {
      fetchState,
      formData,
    }
  },
})
</script>
