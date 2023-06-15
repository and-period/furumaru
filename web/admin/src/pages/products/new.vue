<script lang="ts" setup>
import { storeToRefs } from 'pinia'

import { useAlert } from '~/lib/hooks'
import {
  useCategoryStore,
  useProducerStore,
  useProductStore,
  useProductTypeStore
} from '~/store'
import { CreateProductRequest, CreateProductRequestMediaInner } from '~/types/api'

const router = useRouter()
const productStore = useProductStore()
const productTypeStore = useProductTypeStore()
const categoryStore = useCategoryStore()
const producerStore = useProducerStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const { producers } = storeToRefs(producerStore)
const { productTypes } = storeToRefs(productTypeStore)

const fetchState = useAsyncData(async (): Promise<void> => {
  await Promise.all([
    productTypeStore.fetchProductTypes(),
    categoryStore.fetchCategories(),
    producerStore.fetchProducers(20, 0, '')
  ])
})

const loading = ref<boolean>(false)
const formData = ref<CreateProductRequest>({
  name: '',
  description: '',
  producerId: '',
  productTypeId: '',
  public: true,
  inventory: 0,
  weight: 0,
  itemUnit: '',
  itemDescription: '',
  media: [],
  price: 0,
  deliveryType: 1,
  box60Rate: 0,
  box80Rate: 0,
  box100Rate: 0,
  originPrefecture: '',
  originCity: ''
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
    await productStore.createProduct(formData.value)
    router.push('/products')
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)

    window.scrollTo({
      top: 0,
      behavior: 'smooth'
    })
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
  <templates-product-new
    v-model:form-data="formData"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :producers="producers"
    :product-types="productTypes"
    @update:files="handleImageUpload"
    @submit="handleSubmit"
  />
</template>
