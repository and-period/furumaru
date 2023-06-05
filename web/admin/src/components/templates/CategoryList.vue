<script lang="ts" setup>
import { mdiPlus } from '@mdi/js'

import { useVuelidate } from '@vuelidate/core'
import { Category } from '~/types/props'
import { AlertType } from '~/lib/hooks'
import { CreateCategoryRequest, CategoriesResponseCategoriesInner, ProductTypesResponseProductTypesInner, CreateProductTypeRequest } from '~/types/api'

import {
  required,
  getErrorMessage,
  maxLength
} from '~/lib/validations'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  categoryDialog: {
    type: Boolean,
    default: false
  },
  productTypeDialog: {
    type: Boolean,
    default: false
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

const emits = defineEmits<{
  (e: 'click:category-update-page', page: number): void
  (e: 'click:category-update-items-per-page', page: number): void
  (e: 'click:category-more-page'): void
  (e: 'click:product-type-update-page', page: number): void
  (e: 'click:product-type-update-items-per-page', page: number): void
  (e: 'update:category-dialog', val: boolean): void
  (e: 'update:product-type-dialog', val: boolean): void
  (e: 'update:category-form-data', val: CreateCategoryRequest): void
  (e: 'update:product-type-form-data', val: CreateProductTypeRequest): void
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

const categoryDialogValue = computed({
  get: () => props.categoryDialog,
  set: (val: boolean) => emits('update:category-dialog', val)
})
const productTypeDialogValue = computed({
  get: () => props.productTypeDialog,
  set: (val: boolean) => emits('update:product-type-dialog', val)
})
const categoryFormDataValue = computed({
  get: () => props.categoryFormData,
  set: (val: CreateCategoryRequest) => emits('update:category-form-data', val)
})
const categoryFormDataRules = computed(() => {
  return {
    name: { required, maxlength: maxLength(32) }
  }
})

const cv$ = useVuelidate<CreateCategoryRequest>(categoryFormDataRules, categoryFormDataValue)

const productTypeFormDataValue = computed({
  get: () => props.productTypeFormData,
  set: (val: CreateProductTypeRequest) => emits('update:product-type-form-data', val)
})

const productTypeFormDataRules = computed(() => {
  return {
    name: { required, maxlength: maxLength(32) },
    iconUrl: { required }
  }
})

const pv$ = useVuelidate<CreateProductTypeRequest>(productTypeFormDataRules, productTypeFormDataValue)

watch(selector, () => {
  emits('update:tab', selector.value)
})

const onClickCategoryPage = (page: number): void => {
  emits('click:category-update-page', page)
}

const onClickCategoryItemsPerPage = (page: number): void => {
  emits('click:category-update-items-per-page', page)
}

const onClickCategoryMorePage = (): void => {
  emits('click:category-more-page')
}

const onClickCategoryOpenDialog = (): void => {
  categoryDialogValue.value = true
}

const onClickCategoryCloseDialog = (): void => {
  categoryDialogValue.value = false
}

const onClickProductTypePage = (page: number): void => {
  emits('click:product-type-update-page', page)
}

const onClickProductTypeItemsPerPage = (page: number): void => {
  emits('click:product-type-update-items-per-page', page)
}

const onClickProductTypeOpenDialog = (): void => {
  productTypeDialogValue.value = true
}

const onClickProductTypeCloseDialog = (): void => {
  productTypeDialogValue.value = false
}

const onUploadProductTypeIcon = (event: Event): void => {
  const target = event.target as HTMLInputElement
  if (!target.files) {
    return
  }
  emits('update:product-type-upload-icon', target.files)
}

const onSubmitCategory = async () => {
  const result = await cv$.value.$validate()
  if (!result) {
    return
  }
  emits('submit:category')
}

const onSubmitProductType = (): void => {
  emits('submit:product-type', categoryId.value)
}
</script>

<template>
  <v-alert v-show="props.isAlert" :type="props.alertType" v-text="props.alertText" />

  <v-dialog v-model="categoryDialogValue" width="500">
    <v-card :loading="loading">
      <v-card-title class="text-h6 primaryLight">
        カテゴリー登録
      </v-card-title>

      <v-card-text>
        <v-text-field
          v-model="cv$.name.$model"
          class="mx-4"
          label="カテゴリー名"
          :error-messages="getErrorMessage(cv$.name.$errors)"
        />
      </v-card-text>
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

  <v-dialog v-model="productTypeDialogValue" width="500">
    <v-card :loading="loading">
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
          v-model="pv$.name.$model"
          :error-messages="getErrorMessage(pv$.name.$errors)"
          label="品目"
        />
        <v-card class="text-center" role="button" flat>
          <v-card-text>
            <v-avatar size="96">
              <v-icon v-if="productTypeFormData.iconUrl === ''" size="x-large" :icon="mdiPlus" />
              <v-img
                v-else
                :src="productTypeFormData.iconUrl"
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
      </v-card-text>
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

  <v-card>
    <v-card-title class="d-flex flex-row">
      カテゴリー・品目設定
      <v-spacer />
      <v-btn v-show="selector === 'categories'" variant="outlined" color="primary" @click="onClickCategoryOpenDialog">
        <v-icon start :icon="mdiPlus" />
        カテゴリー登録
      </v-btn>
      <v-btn v-show="selector === 'productTypes'" variant="outlined" color="primary" @click="onClickProductTypeOpenDialog">
        <v-icon start :icon="mdiPlus" />
        品目登録
      </v-btn>
    </v-card-title>

    <v-tabs v-model="selector" grow color="dark">
      <v-tab v-for="item in tabs" :key="item.value" :value="item.value">
        {{ item.name }}
      </v-tab>
    </v-tabs>

    <v-card-text>
      <v-window v-model="selector">
        <v-window-item value="categories">
          <organisms-category-list
            :categories="categories"
            :table-items-per-page="categoryTableItemsPerPage"
            :table-items-length="categoryTableItemsTotal"
            @update:page="onClickCategoryPage"
            @update:items-per-page="onClickCategoryItemsPerPage"
          />
        </v-window-item>

        <v-window-item value="productTypes">
          <organisms-product-type-list
            :product-types="productTypes"
            :categories="categories"
            :table-items-per-page="productTypeTableItemsPerPage"
            :table-items-length="productTypeTableItemsTotal"
            @update:page="onClickProductTypePage"
            @update:items-per-page="onClickProductTypeItemsPerPage"
            @click:more-item="onClickCategoryMorePage"
          />
        </v-window-item>
      </v-window>
    </v-card-text>
  </v-card>
</template>
