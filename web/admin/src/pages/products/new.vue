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
  useProductTypeStore,
} from '~/store'
import { Prefecture } from '~/types'
import { DeliveryType, ProductScope, StorageMethodType } from '~/types/api/v1'
import type { CreateProductRequest, CreateProductMedia } from '~/types/api/v1'

const route = useRoute()
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

const copyFromProductId = computed(() => {
  const fromId = route.query.from as string
  if (fromId) {
    return fromId
  }
  else {
    return ''
  }
})

const loading = ref<boolean>(false)
const selectedCategoryId = ref<string>()
const formData = ref<CreateProductRequest>({
  name: '',
  description: '',
  scope: ProductScope.ProductScopePublic,
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
  deliveryType: DeliveryType.DeliveryTypeNormal,
  recommendedPoint1: '',
  recommendedPoint2: '',
  recommendedPoint3: '',
  expirationDate: 0,
  storageMethodType: StorageMethodType.StorageMethodTypeNormal,
  box60Rate: 0,
  box80Rate: 0,
  box100Rate: 0,
  originPrefectureCode: Prefecture.HOKKAIDO,
  originCity: '',
  startAt: dayjs().unix(),
  endAt: dayjs().unix(),
})

const fetchAndSetFormDataByCopyFromProduct = async (productId: string) => {
  try {
    const sourceProduct = await productStore.getProduct(
      copyFromProductId.value,
    )
    formData.value.name = `${sourceProduct.product.name}のコピー`
    formData.value.description = sourceProduct.product.description
    formData.value.scope = sourceProduct.product.scope
    formData.value.coordinatorId = sourceProduct.product.coordinatorId
    formData.value.producerId = sourceProduct.product.producerId
    formData.value.productTypeId = sourceProduct.product.productTypeId
    formData.value.productTagIds = sourceProduct.product.productTagIds
    formData.value.price = sourceProduct.product.price
    formData.value.cost = sourceProduct.product.cost
    formData.value.inventory = sourceProduct.product.inventory
    formData.value.weight = sourceProduct.product.weight
    formData.value.itemUnit = sourceProduct.product.itemUnit
    formData.value.itemDescription = sourceProduct.product.itemDescription
    formData.value.deliveryType = sourceProduct.product.deliveryType
    formData.value.recommendedPoint1 = sourceProduct.product.recommendedPoint1
    formData.value.recommendedPoint2 = sourceProduct.product.recommendedPoint2
    formData.value.recommendedPoint3 = sourceProduct.product.recommendedPoint3
    formData.value.expirationDate = sourceProduct.product.expirationDate
    formData.value.storageMethodType = sourceProduct.product.storageMethodType
    formData.value.box60Rate = sourceProduct.product.box60Rate
    formData.value.box80Rate = sourceProduct.product.box80Rate
    formData.value.box100Rate = sourceProduct.product.box100Rate
    formData.value.originPrefectureCode
      = sourceProduct.product.originPrefectureCode
    formData.value.originCity = sourceProduct.product.originCity
    formData.value.startAt = sourceProduct.product.startAt
    formData.value.endAt = sourceProduct.product.endAt
    formData.value.media = sourceProduct.product.media.map(
      (item): CreateProductMedia => ({
        url: item.url,
        isThumbnail: item.isThumbnail,
      }),
    )
    selectedCategoryId.value = sourceProduct.product.categoryId
  }
  catch (err) {
    if (err instanceof Error) {
      show(`複製元の商品を取得できませんでした。${err.message}`)
    }
    console.log(err)
  }
}

const fetchState = useAsyncData(async (): Promise<void> => {
  if (copyFromProductId.value) {
    fetchAndSetFormDataByCopyFromProduct(copyFromProductId.value)
  }
  await Promise.all([
    categoryStore.fetchCategories(),
    producerStore.fetchProducers(20, 0, ''),
    productTagStore.fetchProductTags(20, 0, []),
  ])
})

watch(selectedCategoryId, (newValue?: string, oldValue?: string): void => {
  productTypeStore.fetchProductTypesByCategoryId(
    selectedCategoryId.value || '',
  )
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
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleSearchCategory = async (name: string): Promise<void> => {
  try {
    const categoryIds: string[] = selectedCategoryId.value
      ? [selectedCategoryId.value]
      : []
    await categoryStore.searchCategories(name, categoryIds)
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleSearchProductType = async (name: string): Promise<void> => {
  try {
    const productTypeIds: string[] = formData.value.productTypeId
      ? [formData.value.productTypeId]
      : []
    await productTypeStore.searchProductTypes(
      name,
      selectedCategoryId.value,
      productTypeIds,
    )
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleSearchProductTag = async (name: string): Promise<void> => {
  try {
    await productTagStore.searchProductTags(name, formData.value.productTagIds)
  }
  catch (err) {
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
    }
    catch (err) {
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
  formData.value.media = formData.value.media.map(
    (item, i): CreateProductMedia => ({
      ...item,
      isThumbnail: i === 0,
    }),
  )
}

const handleSubmit = async (): Promise<void> => {
  const req = {
    ...formData.value,
    coordinatorId: auth.value?.adminId || '',
  }
  try {
    loading.value = true
    await productStore.createProduct(req)
    router.push('/products')
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)

    window.scrollTo({
      top: 0,
      behavior: 'smooth',
    })
  }
  finally {
    loading.value = false
  }
}

try {
  await fetchState.execute()
}
catch (err) {
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
