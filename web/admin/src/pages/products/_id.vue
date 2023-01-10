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

import { useProductStore } from '~/store/product'
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

    const { fetchState } = useFetch(async () => {
      try {
        await productStore.getProduct(id)
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
