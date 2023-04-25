<script lang="ts" setup>
import { mdiPlus } from '@mdi/js'
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

const productTypeStore = useProductTypeStore()
const categoryStore = useCategoryStore()
const producerStore = useProducerStore()

const { producers } = storeToRefs(producerStore)
const { productTypes } = storeToRefs(productTypeStore)

const fetchState = useAsyncData(async () => {
  await Promise.all([
    productTypeStore.fetchProductTypes(),
    categoryStore.fetchCategories(),
    producerStore.fetchProducers(20, 0, '')
  ])
})

const router = useRouter()

const { uploadProductImage, createProduct } = useProductStore()
const breadcrumbsItem = [
  {
    text: '商品管理',
    href: '/products',
    disabled: false
  },
  {
    text: '商品登録',
    href: 'add',
    disabled: true
  }
]

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

const v$ = useVuelidate()

const handleImageUpload = async (files: FileList) => {
  for (const [index, file] of Array.from(files).entries()) {
    try {
      const uploadImage: UploadImageResponse = await uploadProductImage(file)
      formData.value.media.push({
        ...uploadImage,
        isThumbnail: index === 0
      })
    } catch (error) {
      console.log(error)
    }
  }
}

const handleDeleteThumbnailImageButton = (index: number) => {
  formData.value.media = formData.value.media.filter((_, i) => {
    return i !== index
  })
}

const { alertType, isShow, alertText, show } = useAlert('error')

const handleFormSubmit = async () => {
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
  } catch (error) {
    if (error instanceof ApiBaseError) {
      show(error.message)
    } else {
      show('不明なエラーが発生しました。')
    }
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
  <div>
    <v-card-title>商品登録</v-card-title>
    <v-breadcrumbs :items="breadcrumbsItem" large class="pa-0 mb-6" />

    <v-alert v-model="isShow" class="mb-4" :type="alertType" v-text="alertText" />

    <organisms-product-form
      v-model:form-data="formData"
      :producers-items="producers"
      :product-types-items="productTypes"
      @update:files="handleImageUpload"
      @delete:thumbnail-image="handleDeleteThumbnailImageButton"
    />

    <v-btn block variant="outlined" @click="handleFormSubmit">
      <v-icon start :icon="mdiPlus" />
      登録
    </v-btn>
  </div>
</template>
