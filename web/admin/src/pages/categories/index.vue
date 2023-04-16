<script lang="ts" setup>
import { mdiPlus } from '@mdi/js'
import { usePagination } from '~/lib/hooks'
import { useCategoryStore, useProductTypeStore } from '~/store'
import {
  CreateCategoryRequest,
  CreateProductTypeRequest
} from '~/types/api'
import { ImageUploadStatus } from '~/types/props'
import { Category } from '~/types/props/category'

const categoryStore = useCategoryStore()
const productTypeStore = useProductTypeStore()
const categoryPagination = usePagination()
const productTypePagination = usePagination()

const tabItems: Category[] = [
  { name: 'カテゴリー', value: 'categories' },
  { name: '品目', value: 'categoryItems' }
]

const inputRef = ref<HTMLInputElement | null>(null)
const selector = ref<string>('categories')
const categoryDialog = ref<boolean>(false)
const productTypeDialog = ref<boolean>(false)
const selectedCategoryId = ref<string>('')

const headerUploadStatus = reactive<ImageUploadStatus>({
  error: false,
  message: ''
})
const categoryFormData = reactive<CreateCategoryRequest>({
  name: ''
})
const productTypeFormData = reactive<CreateProductTypeRequest>({
  name: '',
  iconUrl: ''
})

const categories = computed(() => {
  return categoryStore.categories
})
const categoryTotal = computed(() => {
  return categoryStore.totalCategoryItems
})
const productTypes = computed(() => {
  return productTypeStore.productTypes
})
const productTypeTotal = computed(() => {
  return productTypeStore.totalItems
})

watch(categoryPagination.itemsPerPage, () => {
  fetchCategories()
})
watch(productTypePagination.itemsPerPage, () => {
  fetchProductTypes()
})
watch(selector, async () => {
  categoryPagination.updateCurrentPage(1)
  productTypePagination.updateCurrentPage(1)
  await Promise.all([categoryState.execute(), productTypeState.execute()])
})

/*
 * category methods
 */
const categoryState = useAsyncData(async () => {
  await fetchCategories()
})

const fetchCategories = async () => {
  try {
    await categoryStore.fetchCategories(categoryPagination.itemsPerPage.value, categoryPagination.offset.value)
  } catch (err) {
    console.log(err)
  }
}

const moreCategories = async () => {
  try {
    categoryPagination.updateCurrentPage(categoryPagination.offset.value + 1)
    await categoryStore.moreCategories(categoryPagination.itemsPerPage.value, categoryPagination.offset.value)
  } catch (err) {
    console.log(err)
  }
}

const categoryRegister = async (): Promise<void> => {
  try {
    await categoryStore.createCategory(categoryFormData)
    categoryDialog.value = false
  } catch (err) {
    console.log(err)
  }
}

const categoryCancel = (): void => {
  categoryDialog.value = false
}

/*
 * productType methods
 */
const productTypeState = useAsyncData(async () => {
  await fetchProductTypes()
})

const handleUpdateCategoryPage = async (page: number) => {
  categoryPagination.updateCurrentPage(page)
  await fetchCategories()
}

const fetchProductTypes = async () => {
  try {
    await productTypeStore.fetchProductTypes(productTypePagination.itemsPerPage.value, productTypePagination.offset.value)
  } catch (err) {
    console.log(err)
  }
}

const productTypeRegister = async (): Promise<void> => {
  try {
    await productTypeStore.createProductType(selectedCategoryId.value, productTypeFormData)
    productTypeDialog.value = false
  } catch (err) {
    console.log(err)
  }
}

const handleUploadProductTypeIcon = () => {
  const files = inputRef.value?.files
  if (inputRef.value && inputRef.value.files) {
    if (files && files.length > 0) {
      productTypeStore
        .uploadProductTypeIcon(files[0])
        .then((res) => {
          productTypeFormData.iconUrl = res.url
        })
        .catch(() => {
          headerUploadStatus.error = true
          headerUploadStatus.message = 'アップロードに失敗しました。'
        })
    }
  }
}

const productTypeCancel = (): void => {
  productTypeDialog.value = false
}

/**
 * common methods
 */
const isLoading = (): boolean => {
  return categoryState.pending.value && productTypeState.pending.value
}

const handleClick = () => {
  if (!inputRef.value) {
    return
  }
  inputRef.value.click()
}

try {
  await Promise.all([categoryState.execute(), productTypeState.execute()])
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <div>
    <v-card-title>カテゴリー・品目設定</v-card-title>
    <v-tabs v-model="selector" grow color="dark">
      <v-tabs-slider color="accent" />
      <v-tab v-for="item in tabItems" :key="item.value" :value="item.value">
        {{ item.name }}
      </v-tab>
    </v-tabs>

    <v-window v-model="selector">
      <v-window-item value="categories">
        <v-dialog v-model="categoryDialog" width="500">
          <template #activator="{ on, attrs }">
            <div class="d-flex pt-3 pr-3">
              <v-spacer />
              <v-btn variant="outlined" color="primary" v-bind="attrs" v-on="on">
                <v-icon start :icon="mdiPlus" />
                追加
              </v-btn>
            </div>
          </template>

          <v-card :loading="isLoading()">
            <v-card-title class="text-h6 primaryLight">
              カテゴリー登録
            </v-card-title>
            <v-text-field
              v-model="categoryFormData.name"
              class="mx-4"
              maxlength="32"
              label="カテゴリー"
            />
            <v-divider />

            <v-card-actions>
              <v-spacer />
              <v-btn color="error" variant="text" @click="categoryCancel">
                キャンセル
              </v-btn>
              <v-btn color="primary" variant="outlined" @click="categoryRegister">
                登録
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
        <organisms-category-list
          :categories="categories"
          :table-items-per-page="categoryPagination.itemsPerPage.value"
          :table-items-length="categoryTotal"
          :table-footer-options="categoryPagination.options"
          @update:page="handleUpdateCategoryPage"
          @update:items-per-page="categoryPagination.handleUpdateItemsPerPage"
        />
      </v-window-item>

      <v-window-item value="categoryItems">
        <v-dialog v-model="productTypeDialog" width="500">
          <template #activator="{ on, attrs }">
            <div class="d-flex pt-3 pr-3">
              <v-spacer />
              <v-btn variant="outlined" color="primary" v-bind="attrs" v-on="on">
                <v-icon start :icon="mdiPlus" />
                追加
              </v-btn>
            </div>
          </template>
          <v-card :loading="isLoading()">
            <v-card-title class="primaryLight">
              品目登録
            </v-card-title>
            <v-card-text class="mt-4">
              <v-autocomplete
                v-model="selectedCategoryId"
                :items="categories"
                item-title="name"
                item-value="id"
                label="カテゴリー"
              >
                <template #append-item>
                  <div class="pa-2">
                    <v-btn
                      block
                      color="primary"
                      variant="outlined"
                      @click="moreCategories"
                    >
                      <v-icon :icon="mdiPlus" />
                      さらに読み込む
                    </v-btn>
                  </div>
                </template>
              </v-autocomplete>
              <v-spacer />
              <v-text-field
                v-model="productTypeFormData.name"
                maxlength="32"
                label="品目"
              />
            </v-card-text>
            <v-card class="text-center" role="button" flat @click="handleClick">
              <v-card-text>
                <v-avatar size="96">
                  <v-icon v-if="productTypeFormData.iconUrl === ''" size="x-large" :icon="mdiPlus" />
                  <v-img
                    v-else
                    :src="productTypeFormData.iconUrl"
                    aspect-ratio="1"
                    max-height="150"
                    contain
                  />
                </v-avatar>
                <input
                  ref="inputRef"
                  type="file"
                  class="d-none"
                  accept="image/png, image/jpeg"
                  @change="handleUploadProductTypeIcon"
                >
                <p class="ma-0">
                  アイコン画像を選択
                </p>
              </v-card-text>
            </v-card>
            <p v-show="headerUploadStatus.error" class="red--text ma-0">
              {{ headerUploadStatus.message }}
            </p>
            <v-divider />

            <v-card-actions>
              <v-spacer />
              <v-btn color="error" variant="text" @click="productTypeCancel">
                キャンセル
              </v-btn>
              <v-btn color="primary" variant="outlined" @click="productTypeRegister">
                登録
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
        <organisms-product-type-list
          :product-types="productTypes"
          :categories="categories"
          :table-items-per-page="productTypePagination.itemsPerPage.value"
          :table-items-length="productTypeTotal"
          :table-footer-options="productTypePagination.options"
          @update:page="handleUpdateCategoryPage"
          @update:items-per-page="productTypePagination.handleUpdateItemsPerPage"
          @click:more-item="moreCategories"
        />
      </v-window-item>
    </v-window>
  </div>
</template>
