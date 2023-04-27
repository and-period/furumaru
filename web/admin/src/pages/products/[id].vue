<script lang="ts" setup>
import { mdiPlus } from '@mdi/js'
import useVuelidate from '@vuelidate/core'
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

const v$ = useVuelidate()

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
      formData.value.media.push({
        ...uploadImage,
        isThumbnail: false
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

const commonStore = useCommonStore()

const isLoading = (): boolean => {
  return fetchState?.pending?.value || false
}

const handleSubmit = async () => {
  const result = await v$.value.$validate()
  console.log(result, v$)
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
  <div>
    <v-alert v-model="isShow" :type="alertType" class="mb-2" v-text="alertText" />

    <organisms-product-form
      v-model:form-data="formData"
      :producers-items="producersItems"
      :product-types-items="productTypesItems"
      @update:files="handleImageUpload"
      @delete:thumbnail-image="handleDeleteThumbnailImageButton"
    />

    <v-btn block variant="outlined" @click="handleSubmit">
      <v-icon start :icon="mdiPlus" />
      更新
    </v-btn>
  </div>
</template>
