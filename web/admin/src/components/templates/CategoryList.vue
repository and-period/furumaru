<script lang="ts" setup>
import { mdiPlus } from '@mdi/js'

import { AlertType } from '~/lib/hooks'
import { CreateCategoryRequest, CategoriesResponseCategoriesInner, ProductTypesResponseProductTypesInner, CreateProductTypeRequest } from '~/types/api'
import { Category } from '~/types/props'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  dialog: {
    type: Object,
    default: () => ({
      dialog: false,
      productType: false
    })
  },
  isAlert: {
    type: Boolean,
    default: false
  },
  alertType: {
    type: String as PropType<AlertType>,
    default: undefined
  },
  alertText: {
    type: String,
    default: ''
  },
  categories: {
    type: Array<CategoriesResponseCategoriesInner>,
    default: () => []
  },
  categoryFormData: {
    type: Object as PropType<CreateCategoryRequest>,
    default: () => ({
      name: ''
    })
  },
  categoryTableItemsPerPage: {
    type: Number,
    default: 20
  },
  categoryTableItemsTotal: {
    type: Number,
    default: 0
  },
  productTypes: {
    type: Array<ProductTypesResponseProductTypesInner>,
    default: () => []
  },
  productTypeFormData: {
    type: Object as PropType<CreateProductTypeRequest>,
    default: () => ({
      name: '',
      iconUrl: ''
    })
  },
  productTypeTableItemsPerPage: {
    type: Number,
    default: 20
  },
  productTypeTableItemsTotal: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits<{
  (e: 'click:category-update-page', page: number): void
  (e: 'click:category-update-items-per-page', page: number): void
  (e: 'click:category-more-page'): void
  (e: 'click:category-close-dialog'): void
  (e: 'click:product-type-update-page', page: number): void
  (e: 'click:product-type-update-items-per-page', page: number): void
  (e: 'click:product-type-close-dialog'): void
  (e: 'update:tab', key: string): void
  (e: 'update:product-type-upload-icon', file: FileList): void
  (e: 'submit:category'): void
  (e: 'submit:product-type', categoryId: string): void
}>()

const tabs: Category[] = [
  { name: 'カテゴリー', value: 'categories' },
  { name: '品目', value: 'productTypes' }
]

const selector = ref<string>('categories')
const categoryId = ref<string>('')
const productTypeIcon = ref<HTMLIFrameElement>()

watch(selector, () => {
  emit('update:tab', selector.value)
})

const onClickCategoryPage = (page: number): void => {
  emit('click:category-update-page', page)
}

const onClickCategoryItemsPerPage = (page: number): void => {
  emit('click:category-update-items-per-page', page)
}

const onClickCategoryMorePage = (): void => {
  emit('click:category-more-page')
}

const onClickCategoryCloseDialog = (): void => {
  emit('click:category-close-dialog')
}

const onClickProductTypePage = (page: number): void => {
  emit('click:product-type-update-page', page)
}

const onClickProductTypeItemsPerPage = (page: number): void => {
  emit('click:product-type-update-items-per-page', page)
}

const onClickProductTypeCloseDialog = (): void => {
  emit('click:product-type-close-dialog')
}

const onUploadProductTypeIcon = (event: Event): void => {
  const target = event.target as HTMLInputElement
  if (!target.files) {
    return
  }
  emit('update:product-type-upload-icon', target.files)
}

const onSubmitCategory = (): void => {
  emit('submit:category')
}

const onSubmitProductType = (): void => {
  emit('submit:product-type', categoryId.value)
}
</script>

<template>
  <v-alert v-show="props.isAlert" :type="props.alertType" v-text="props.alertText" />
  <v-card-title>カテゴリー・品目設定</v-card-title>
  <v-tabs v-model="selector" grow color="dark">
    <v-tab v-for="item in tabs" :key="item.value" :value="item.value">
      {{ item.name }}
    </v-tab>
  </v-tabs>

  <v-window v-model="selector">
    <v-window-item value="categories">
      <v-dialog v-model="props.dialog.category" width="500">
        <template #activator="{ props }">
          <div class="d-flex pt-3 pr-3">
            <v-spacer />
            <v-btn variant="outlined" color="primary" v-bind="props">
              <v-icon start :icon="mdiPlus" />
              追加
            </v-btn>
          </div>
        </template>

        <v-card :loading="props.loading">
          <v-card-title class="text-h6 primaryLight">
            カテゴリー登録
          </v-card-title>
          <v-text-field
            v-model="props.categoryFormData.name"
            class="mx-4"
            maxlength="32"
            label="カテゴリー"
          />
          <v-divider />

          <v-card-actions>
            <v-spacer />
            <v-btn color="error" variant="text" @click="onClickCategoryCloseDialog">
              キャンセル
            </v-btn>
            <v-btn color="primary" variant="outlined" @click="onSubmitCategory">
              登録
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      <organisms-category-list
        :categories="props.categories"
        :table-items-per-page="props.categoryTableItemsPerPage"
        :table-items-length="props.categoryTableItemsTotal"
        @update:page="onClickCategoryPage"
        @update:items-per-page="onClickCategoryItemsPerPage"
      />
    </v-window-item>

    <v-window-item value="productTypes">
      <v-dialog v-model="props.dialog.productType" width="500">
        <template #activator="{ props }">
          <div class="d-flex pt-3 pr-3">
            <v-spacer />
            <v-btn variant="outlined" color="primary" v-bind="props">
              <v-icon start :icon="mdiPlus" />
              追加
            </v-btn>
          </div>
        </template>
        <v-card :loading="props.loading">
          <v-card-title class="primaryLight">
            品目登録
          </v-card-title>
          <v-card-text class="mt-4">
            <v-autocomplete
              v-model="categoryId"
              :items="categories"
              item-title="name"
              item-value="id"
              label="カテゴリー"
            >
              <template #append-item>
                <div class="pa-2">
                  <v-btn block color="primary" variant="outlined" @click="onClickCategoryMorePage">
                    <v-icon :icon="mdiPlus" />
                    さらに読み込む
                  </v-btn>
                </div>
              </template>
            </v-autocomplete>
            <v-spacer />
            <v-text-field
              v-model="props.productTypeFormData.name"
              maxlength="32"
              label="品目"
            />
          </v-card-text>
          <v-card class="text-center" role="button" flat @click="onSubmitProductType">
            <v-card-text>
              <v-avatar size="96">
                <v-icon v-if="props.productTypeFormData.iconUrl === ''" size="x-large" :icon="mdiPlus" />
                <v-img
                  v-else
                  :src="props.productTypeFormData.iconUrl"
                  aspect-ratio="1"
                  max-height="150"
                  cover
                />
              </v-avatar>
              <input
                ref="productTypeIcon"
                type="file"
                class="d-none"
                accept="image/png, image/jpeg"
                @change="onUploadProductTypeIcon"
              >
              <p class="ma-0">
                アイコン画像を選択
              </p>
            </v-card-text>
          </v-card>
          <v-divider />

          <v-card-actions>
            <v-spacer />
            <v-btn color="error" variant="text" @click="onClickProductTypeCloseDialog">
              キャンセル
            </v-btn>
            <v-btn color="primary" variant="outlined" @click="onSubmitProductType">
              登録
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
      <organisms-product-type-list
        :product-types="props.productTypes"
        :categories="props.categories"
        :table-items-per-page="props.productTypeTableItemsPerPage"
        :table-items-length="props.productTypeTableItemsTotal"
        @update:page="onClickProductTypePage"
        @update:items-per-page="onClickProductTypeItemsPerPage"
        @click:more-item="onClickCategoryMorePage"
      />
    </v-window-item>
  </v-window>
</template>
