<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useAlert } from '~/lib/hooks'
import {
  useCommonStore,
  useProducerStore,
  useProductStore,
  useProductTypeStore
} from '~/store'
import { UpdateProductRequest, CreateProductRequestMediaInner, UploadImageResponse } from '~/types/api'

const route = useRoute()
const router = useRouter()
const commonStore = useCommonStore()
const productStore = useProductStore()
const productTypeStore = useProductTypeStore()
const producerStore = useProducerStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const { product } = storeToRefs(productStore)
const { producers } = storeToRefs(producerStore)
const { productTypes } = storeToRefs(productTypeStore)

const productId = route.params.id as string

const loading = ref<boolean>(false)
const formData = ref<UpdateProductRequest>({
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

const fetchState = useAsyncData(async (): Promise<void> => {
  try {
    await Promise.all([
      productStore.getProduct(productId),
      productTypeStore.fetchProductTypes(),
      producerStore.fetchProducers()
    ])
    formData.value = { ...product.value }
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)

    window.scrollTo({
      top: 0,
      behavior: 'smooth'
    })
  }
})

const isLoading = (): boolean => {
  return fetchState?.pending?.value || loading.value
}

const handleImageUpload = async (files: FileList): Promise<void> => {
  loading.value = true
  for (const [index, file] of Array.from(files).entries()) {
    try {
      const uploadImage = await productStore.uploadProductImage(file)
      formData.value.media.push({
        ...uploadImage,
        isThumbnail: index === 0
      })
    } catch (err) {
      if (err instanceof Error) {
        show(err.message)
      }
      console.log(err)
    }
  }
  loading.value = false

  const thumbnailItem = formData.value.media.find(item => item.isThumbnail)
  if (thumbnailItem) {
    return
  }
  formData.value.media = formData.value.media.map((item, i): CreateProductRequestMediaInner => ({
    ...item,
    isThumbnail: i === 0
  }))
}

const handleSubmit = async (): Promise<void> => {
  try {
    loading.value = true
    await productStore.updateProduct(productId, formData.value)
    commonStore.addSnackbar({
      color: 'success',
      message: '商品を更新しました。'
    })
    router.push('/products')
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  } finally {
    loading.value = false
  }
}

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-product-edit
    v-model:form-data="formData"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :product="product"
    :producers="producers"
    :product-types="productTypes"
    @update:files="handleImageUpload"
    @submit="handleSubmit"
  />
</template>
