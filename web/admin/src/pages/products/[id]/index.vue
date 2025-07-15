<script lang="ts" setup>
import dayjs from 'dayjs'
import { storeToRefs } from 'pinia'
import { useAlert } from '~/lib/hooks'
import {
  useAuthStore,
  useCategoryStore,
  useCommonStore,
  useProducerStore,
  useProductStore,
  useProductTagStore,
  useProductTypeStore,
} from '~/store'
import {
  type UpdateProductRequest,
  type CreateProductRequestMediaInner,
  DeliveryType,
  StorageMethodType,
  Prefecture,
} from '~/types/api'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const categoryStore = useCategoryStore()
const commonStore = useCommonStore()
const producerStore = useProducerStore()
const productStore = useProductStore()
const productTagStore = useProductTagStore()
const productTypeStore = useProductTypeStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const { adminType } = storeToRefs(authStore)
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
  endAt: dayjs().unix(),
})

const fetchState = useAsyncData(async (): Promise<void> => {
  try {
    await productStore.getProduct(productId)
    selectedCategoryId.value = product.value.categoryId
    formData.value = { ...product.value }
    if (categories.value.length === 0) {
      categoryStore.fetchCategories(20)
    }
    if (productTags.value.length === 0) {
      productTagStore.fetchProductTags(20)
    }
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
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
  for (const [_, file] of Array.from(files).entries()) {
    try {
      const url: string = await productStore.uploadProductMedia(file)
      formData.value.media.push({ url, isThumbnail: false })
    }
    catch (err) {
      if (err instanceof Error) {
        show(err.message)
      }
      console.log(err)
    }
  }
  loading.value = false

  // サムネイル画像が設定済みかをmediaの配列を走査して確認
  const thumbnailItem = formData.value.media.find(item => item.isThumbnail)
  if (thumbnailItem) {
    // 設定されていれば処理終了
    return
  }

  // 設定されていなければ、mediaの最初の要素をサムネイルに設定
  formData.value.media = formData.value.media.map(
    (item, i): CreateProductRequestMediaInner => ({
      ...item,
      isThumbnail: i === 0,
    }),
  )
}

const handleSubmit = async (): Promise<void> => {
  try {
    loading.value = true
    await productStore.updateProduct(productId, formData.value)
    commonStore.addSnackbar({
      color: 'success',
      message: '商品を更新しました。',
    })
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

const handleClickBackButton = () => {
  router.back()
}

const productEditRef = ref<UpdateProductRequest>(null)

const handleClickUpdateButton = async () => {
  // Call the onSubmit method of the ProductEdit component to trigger validation
  if (productEditRef.value?.onSubmit) {
    await productEditRef.value.onSubmit()
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
  <div>
    <templates-product-edit
      ref="productEditRef"
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
      :admin-type="adminType"
      class="mb-16"
      @update:files="handleImageUpload"
      @update:search-category="handleSearchCategory"
      @update:search-product-type="handleSearchProductType"
      @update:search-product-tag="handleSearchProductTag"
      @submit="handleSubmit"
    />
    <div
      class="position-fixed bottom-0 left-0 w-100 bg-white pa-4 text-right elevation-3"
    >
      <div class="d-inline-flex ga-4">
        <v-btn
          color="secondary"
          variant="outlined"
          @click="handleClickBackButton"
        >
          戻る
        </v-btn>
        <v-btn
          color="primary"
          variant="outlined"
          :loading="loading"
          @click="handleClickUpdateButton"
        >
          更新
        </v-btn>
      </div>
    </div>
  </div>
</template>
