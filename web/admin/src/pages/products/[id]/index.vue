<script lang="ts" setup>
import dayjs from 'dayjs'
import { storeToRefs } from 'pinia'
import { mdiRobotOutline } from '@mdi/js'
import { useAlert } from '~/lib/hooks'
import { useProductAiAssistant } from '~/composables/useProductAiAssistant'
import {
  useAuthStore,
  useCategoryStore,
  useCommonStore,
  useProducerStore,
  useProductStore,
  useProductTagStore,
  useProductTypeStore,
} from '~/store'
import { useUnsavedChangesGuard } from '~/composables/useUnsavedChangesGuard'
import { Prefecture } from '~/types'
import { DeliveryType, ProductScope, StorageMethodType } from '~/types/api/v1'
import type { UpdateProductRequest, CreateProductMedia, CreateProductReviewRequest } from '~/types/api/v1'

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
const reviewLoading = ref<boolean>(false)
const currentTab = ref<number>(0)
const selectedCategoryId = ref<string>()
const reviewFormData = ref<CreateProductReviewRequest>({
  title: '',
  comment: '',
  rate: 5,
})
const formData = ref<UpdateProductRequest>({
  name: '',
  description: '',
  scope: ProductScope.ProductScopePublic,
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

const aiAssistant = useProductAiAssistant(formData as Ref<Record<string, unknown>>)

const { captureSnapshot, markAsSaved, showLeaveDialog, confirmLeave, cancelLeave }
  = useUnsavedChangesGuard(formData)

const fetchState = useAsyncData('product-detail', async (): Promise<void> => {
  try {
    await productStore.getProduct(productId)
    selectedCategoryId.value = product.value.categoryId
    formData.value = { ...product.value }
    captureSnapshot()
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
    (item, i): CreateProductMedia => ({
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
    markAsSaved()
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

const handleSubmitReview = async (): Promise<void> => {
  try {
    reviewLoading.value = true
    await productStore.createProductReview(productId, reviewFormData.value)
    commonStore.addSnackbar({
      color: 'success',
      message: 'ダミーレビューを投稿しました。',
    })
    // フォームをリセット
    reviewFormData.value = {
      title: '',
      comment: '',
      rate: 5,
    }
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
    reviewLoading.value = false
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
      v-model:form-data="formData"
      v-model:review-form-data="reviewFormData"
      v-model:selected-category-id="selectedCategoryId"
      v-model:current-tab="currentTab"
      :loading="isLoading()"
      :review-loading="reviewLoading"
      :is-alert="isShow"
      :alert-type="alertType"
      :alert-text="alertText"
      :product="product"
      :producers="producers"
      :categories="categories"
      :product-types="productTypes"
      :product-tags="productTags"
      :admin-type="adminType"
      class="mb-20"
      :style="aiAssistant.isPanelOpen.value ? 'margin-right: 400px' : ''"
      @update:files="handleImageUpload"
      @update:search-category="handleSearchCategory"
      @update:search-product-type="handleSearchProductType"
      @update:search-product-tag="handleSearchProductTag"
      @submit="handleSubmit"
      @submit:review="handleSubmitReview"
    />

    <!-- AI Assistant FAB -->
    <v-btn
      v-if="!aiAssistant.isPanelOpen.value"
      color="primary"
      icon
      size="large"
      class="ai-fab"
      elevation="8"
      @click="aiAssistant.togglePanel()"
    >
      <v-icon :icon="mdiRobotOutline" />
      <v-tooltip
        activator="parent"
        location="left"
      >
        AI アシスタント
      </v-tooltip>
    </v-btn>

    <!-- AI Chat Panel -->
    <organisms-ai-assistant-ai-chat-panel
      v-if="aiAssistant.isPanelOpen.value"
      :messages="aiAssistant.messages.value"
      :input="aiAssistant.input.value"
      :loading="aiAssistant.isChatLoading.value"
      :error="aiAssistant.chatError.value"
      :has-pending-approval="aiAssistant.hasPendingApproval.value"
      :pending-changes="aiAssistant.pendingChanges.value"
      :pending-tool-name="aiAssistant.pendingToolName.value"
      @update:input="aiAssistant.input.value = $event"
      @close="aiAssistant.togglePanel()"
      @submit="aiAssistant.sendMessage()"
      @approve="aiAssistant.applyUpdate()"
      @reject="aiAssistant.rejectUpdate()"
    />

    <atoms-app-confirm-dialog
      v-model="showLeaveDialog"
      title="未保存の変更があります"
      message="ページを離れると入力内容が失われます。よろしいですか？"
      confirm-text="破棄して離れる"
      confirm-color="warning"
      @confirm="confirmLeave"
      @cancel="cancelLeave"
    />
  </div>
</template>

<style scoped>
.ai-fab {
  position: fixed;
  bottom: 100px;
  right: 24px;
  z-index: 100;
}
</style>
