<script lang="ts" setup>
import dayjs from 'dayjs'
import { storeToRefs } from 'pinia'

import { useAlert } from '~/lib/hooks'
import {
  useAuthStore,
  useCategoryStore,
  useProducerStore,
  useProductStore,
  useProductTagStore,
  useProductTypeStore
} from '~/store'
import { type CreateProductRequest, type CreateProductRequestMediaInner, DeliveryType, Prefecture, StorageMethodType } from '~/types/api'

const router = useRouter()
const authStore = useAuthStore()
const categoryStore = useCategoryStore()
const producerStore = useProducerStore()
const productStore = useProductStore()
const productTagStore = useProductTagStore()
const productTypeStore = useProductTypeStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const { auth } = storeToRefs(authStore)
const { categories } = storeToRefs(categoryStore)
const { producers } = storeToRefs(producerStore)
const { productTags } = storeToRefs(productTagStore)
const { productTypes } = storeToRefs(productTypeStore)

const loading = ref<boolean>(false)
const selectedCategoryId = ref<string>()
const formData = ref<CreateProductRequest>({
  name: '',
  description: '',
  public: false,
  coordinatorId: '',
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
  deliveryType: DeliveryType.NORMAL,
  recommendedPoint1: '',
  recommendedPoint2: '',
  recommendedPoint3: '',
  expirationDate: 0,
  storageMethodType: StorageMethodType.NORMAL,
  box60Rate: 0,
  box80Rate: 0,
  box100Rate: 0,
  originPrefectureCode: Prefecture.HOKKAIDO,
  originCity: '',
  startAt: dayjs().unix(),
  endAt: dayjs().unix()
})

const fetchState = useAsyncData(async (): Promise<void> => {
  await Promise.all([
    categoryStore.fetchCategories(),
    producerStore.fetchProducers(20, 0, ''),
    productTagStore.fetchProductTags(20, 0, [])
  ])
})

watch(selectedCategoryId, (newValue?: string, oldValue?: string): void => {
  productTypeStore.fetchProductTypesByCategoryId(selectedCategoryId.value || '')
  if (newValue === oldValue) {
    return
  }
  formData.value.productTypeId = ''
})

const isLoading = (): boolean => {
  return fetchState?.pending?.value || loading.value
}

const handleSearchProducer = async (name: string): Promise<void> => {
  try {
    await producerStore.searchProducers(name)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleSearchCategory = async (name: string): Promise<void> => {
  try {
    const categoryIds: string[] = selectedCategoryId.value ? [selectedCategoryId.value] : []
    await categoryStore.searchCategories(name, categoryIds)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleSearchProductType = async (name: string): Promise<void> => {
  try {
    const productTypeIds: string[] = formData.value.productTypeId ? [formData.value.productTypeId] : []
    await productTypeStore.searchProductTypes(name, selectedCategoryId.value, productTypeIds)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleSearchProductTag = async (name: string): Promise<void> => {
  try {
    await productTagStore.searchProductTags(name, formData.value.productTagIds)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleImageUpload = async (files: FileList): Promise<void> => {
  loading.value = true
  for (const [index, file] of Array.from(files).entries()) {
    try {
      const url: string = await productStore.uploadProductMedia(file)
      formData.value.media.push({ url, isThumbnail: index === 0 })
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
  const req = {
    ...formData.value,
    coordinatorId: auth.value?.adminId || ''
  }
  try {
    loading.value = true
    await productStore.createProduct(req)
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
    v-model:selected-category-id="selectedCategoryId"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :producers="producers"
    :categories="categories"
    :product-types="productTypes"
    :product-tags="productTags"
    @update:files="handleImageUpload"
    @update:search-producer="handleSearchProducer"
    @update:search-category="handleSearchCategory"
    @update:search-product-type="handleSearchProductType"
    @update:search-product-tag="handleSearchProductTag"
    @submit="handleSubmit"
  />
</template>
