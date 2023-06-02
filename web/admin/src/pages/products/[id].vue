<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import { storeToRefs } from 'pinia'
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
const router = useRouter()
const commonStore = useCommonStore()
const productStore = useProductStore()
const productTypeStore = useProductTypeStore()
const producerStore = useProducerStore()
const { alertType, isShow, alertText, show } = useAlert('error')
const v$ = useVuelidate()

const { producers } = storeToRefs(producerStore)
const { productTypes } = storeToRefs(productTypeStore)

const id = route.params.id as string

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

const fetchState = useAsyncData(async () => {
  try {
    Promise.all([
      productTypeStore.fetchProductTypes(),
      producerStore.fetchProducers()
    ])
    const data = await productStore.getProduct(id)
    formData.value = data
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
})

const handleImageUpload = async (files?: FileList) => {
  if (!files) {
    return
  }

  for (const [, file] of Array.from(files).entries()) {
    try {
      const uploadImage: UploadImageResponse =
        await productStore.uploadProductImage(file)
      formData.value.media.push({
        ...uploadImage,
        isThumbnail: false
      })
    } catch (error) {
      console.log(error)
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
    await productStore.updateProduct(id, formData.value)
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

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-product-edit
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
