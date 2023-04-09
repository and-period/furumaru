<script lang="ts" setup>
import { useAlert } from '~/lib/hooks'
import {
  useCommonStore,
  useProducerStore,
  useProductStore,
  useProductTypeStore
} from '~/store'
import { UpdateProductRequest, UploadImageResponse } from '~/types/api'
import { ApiBaseError } from '~/types/exception'

const route = useRoute()
const id = route.params.id as string

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
  originCity: ''
})

const productStore = useProductStore()
const productTypeStore = useProductTypeStore()
const producerStore = useProducerStore()

const fetchState = useAsyncData(async () => {
  try {
    Promise.all([
      productTypeStore.fetchProductTypes(),
      producerStore.fetchProducers()
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

const handleImageUpload = async (files?: FileList) => {
  if (!files) {
    return
  }

  for (const [, file] of Array.from(files).entries()) {
    try {
      const uploadImage: UploadImageResponse =
        await productStore.uploadProductImage(file)
      formData.media.push({
        ...uploadImage,
        isThumbnail: false
      })
    } catch (error) {
      console.log(error)
    }
  }
}

const commonStore = useCommonStore()

const isLoading = (): boolean => {
  return fetchState?.pending?.value || false
}

const handleSubmit = async () => {
  try {
    await productStore.updateProduct(id, formData)
    commonStore.addSnackbar({
      color: 'success',
      message: '商品を更新しました。'
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
</script>

<template>
  <div>
    <v-alert v-model="isShow" :type="alertType" v-text="alertText" />

    <the-product-update-form-page
      :loading="isLoading"
      :form-data="formData"
      :producers-items="producersItems"
      :product-types-items="productTypesItems"
      @update:files="handleImageUpload"
      @submit="handleSubmit"
    />
  </div>
</template>
