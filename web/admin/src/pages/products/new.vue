<script lang="ts" setup>
import { useVuelidate } from '@vuelidate/core'
import { storeToRefs } from 'pinia'

import { useAlert } from '~/lib/hooks'
import {
  useCategoryStore,
  useProducerStore,
  useProductStore,
  useProductTypeStore
} from '~/store'
import { CreateProductRequest, UploadImageResponse } from '~/types/api'
import { ApiBaseError } from '~/types/exception'

const router = useRouter()
const productTypeStore = useProductTypeStore()
const categoryStore = useCategoryStore()
const producerStore = useProducerStore()
const { alertType, isShow, alertText, show } = useAlert('error')
const v$ = useVuelidate()

const { producers } = storeToRefs(producerStore)
const { productTypes } = storeToRefs(productTypeStore)
const { uploadProductImage, createProduct } = useProductStore()

const fetchState = useAsyncData(async () => {
  await Promise.all([
    productTypeStore.fetchProductTypes(),
    categoryStore.fetchCategories(),
    producerStore.fetchProducers(20, 0, '')
  ])
})

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

const handleImageUpload = async (files: FileList) => {
  for (const [index, file] of Array.from(files).entries()) {
    try {
      const uploadImage: UploadImageResponse = await uploadProductImage(file)
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

  const thumbnailItem = formData.value.media.find(item => item.isThumbnail)
  if (!thumbnailItem) {
    formData.value.media = formData.value.media.map((item, i) => {
      return {
        ...item,
        isThumbnail: i === 0
      }
    })
  }
}

const handleSubmit = async () => {
  const result = await v$.value.$validate()
  if (!result) {
    window.scrollTo({
      top: 0,
      behavior: 'smooth'
    })
    return
  }
  try {
    await createProduct(formData.value)
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
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :producers="producers"
    :product-types="productTypes"
    @update:files="handleImageUpload"
    @submit="handleSubmit"
  />
</template>
