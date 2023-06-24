<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useAlert } from '~/lib/hooks'
import {
  useCategoryStore,
  useCommonStore,
  useProducerStore,
  useProductStore,
  useProductTagStore,
  useProductTypeStore
} from '~/store'
import { UpdateProductRequest, CreateProductRequestMediaInner, DeliveryType, StorageMethodType, Prefecture } from '~/types/api'

const route = useRoute()
const router = useRouter()
const categoryStore = useCategoryStore()
const commonStore = useCommonStore()
const producerStore = useProducerStore()
const productStore = useProductStore()
const productTagStore = useProductTagStore()
const productTypeStore = useProductTypeStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const { categories } = storeToRefs(categoryStore)
const { producers } = storeToRefs(producerStore)
const { productTags } = storeToRefs(productTagStore)
const { productTypes } = storeToRefs(productTypeStore)
const { product } = storeToRefs(productStore)

const productId = route.params.id as string

const loading = ref<boolean>(false)
const selectedCategoryId = ref<string>()
const formData = ref<UpdateProductRequest>({
  name: '',
  description: '',
  public: false,
  producerId: '',
  productTypeId: '',
  productTagIds: [],
  media: [],
  price: 0,
  cost: 0,
  inventory: 0,
  weight: 0,
  itemUnit: '',
  itemDescription: '',
  deliveryType: DeliveryType.UNKNOWN,
  recommendedPoint1: '',
  recommendedPoint2: '',
  recommendedPoint3: '',
  expirationDate: 0,
  storageMethodType: StorageMethodType.UNKNOWN,
  box60Rate: 0,
  box80Rate: 0,
  box100Rate: 0,
  originPrefecture: Prefecture.HOKKAIDO,
  originCity: ''
})

const fetchState = useAsyncData(async (): Promise<void> => {
  try {
    await Promise.all([
      productStore.getProduct(productId),
      categoryStore.fetchCategories(20, 0),
      producerStore.fetchProducers(20, 0, ''),
      productTagStore.fetchProductTags(20, 0, [])
    ])
    selectedCategoryId.value = product.value.categoryId
    formData.value = { ...product.value }
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
})

watch(selectedCategoryId, (): void => {
  productTypeStore.fetchProductTypesByCategoryId(selectedCategoryId.value || '')
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
  <templates-product-edit
    v-model:form-data="formData"
    v-model:selected-category-id="selectedCategoryId"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :product="product"
    :producers="producers"
    :categories="categories"
    :product-types="productTypes"
    :product-tags="productTags"
    @update:files="handleImageUpload"
    @submit="handleSubmit"
  />
</template>
