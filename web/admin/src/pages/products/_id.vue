<template>
  <div>
    <v-alert v-model="isShow" :type="alertType" v-text="alertText" />

    <the-product-update-form-page
      :loading="fetchState.pending"
      :form-data="formData"
      :producers-items="producersItems"
      :product-types-items="productTypesItems"
      @update:files="handleImageUpload"
      @submit="handleSubmit"
    />
  </div>
</template>

<script lang="ts">
import {
  computed,
  defineComponent,
  reactive,
  useFetch,
  useRoute,
  useRouter,
} from '@nuxtjs/composition-api'

import { useAlert } from '~/lib/hooks'
import { useCommonStore } from '~/store/common'
import { useProducerStore } from '~/store/producer'
import { useProductStore } from '~/store/product'
import { useProductTypeStore } from '~/store/product-type'
import { UpdateProductRequest, UploadImageResponse } from '~/types/api'
import { ApiBaseError } from '~/types/exception'

export default defineComponent({
  setup() {
    const route = useRoute()
    const id = route.value.params.id

    const router = useRouter()

    const { alertType, isShow, alertText, show } = useAlert('error')

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
    const producerStore = useProducerStore()

    const { fetchState } = useFetch(async () => {
      try {
        Promise.all([
          productTypeStore.fetchProductTypes(),
          producerStore.fetchProducers(),
        ])
        const data = await productStore.getProduct(id)
        Object.assign(formData, data)
      } catch (error) {
        if (error instanceof ApiBaseError) {
          show(error.message)
        } else {
          show('不明なエラーが発生しました。')
        }
      }
    })

    const productTypesItems = computed(() => {
      return productTypeStore.productTypes
    })

    const producersItems = computed(() => {
      return producerStore.producers
    })

    const handleImageUpload = async (files: FileList) => {
      for (const [, file] of Array.from(files).entries()) {
        try {
          const uploadImage: UploadImageResponse =
            await productStore.uploadProductImage(file)
          formData.media.push({
            ...uploadImage,
            isThumbnail: false,
          })
        } catch (error) {
          console.log(error)
        }
      }
    }

    const commonStore = useCommonStore()

    const handleSubmit = async () => {
      try {
        await productStore.updateProduct(id, formData)
        commonStore.addSnackbar({
          color: 'success',
          message: '商品を更新しました。',
        })
        router.push('/products')
      } catch (error) {
        if (error instanceof ApiBaseError) {
          show(error.message)
        } else {
          show('不明なエラーが発生しました。')
        }
      }
    }

    return {
      // リアクティブ変数
      fetchState,
      isShow,
      alertType,
      alertText,
      formData,
      productTypesItems,
      producersItems,
      // 関数
      handleImageUpload,
      handleSubmit,
    }
  },
})
</script>
