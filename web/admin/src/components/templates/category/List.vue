<script lang="ts" setup>
import type { VTabs } from 'vuetify/components'
import type { AlertType } from '~/lib/hooks'
import { AdminType } from '~/types/api/v1'
import type { CreateCategoryRequest, Category, ProductType, CreateProductTypeRequest, UpdateCategoryRequest, UpdateProductTypeRequest } from '~/types/api/v1'
import type { ImageUploadStatus } from '~/types/props'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  adminType: {
    type: Number as PropType<AdminType>,
    default: AdminType.AdminTypeUnknown,
  },
  isAlert: {
    type: Boolean,
    default: false,
  },
  alertType: {
    type: String as PropType<AlertType>,
    default: undefined,
  },
  alertText: {
    type: String,
    default: '',
  },
  selectedTabItem: {
    type: String,
    default: 'categories',
  },
  // Category
  createCategoryDialog: {
    type: Boolean,
    default: false,
  },
  updateCategoryDialog: {
    type: Boolean,
    default: false,
  },
  deleteCategoryDialog: {
    type: Boolean,
    default: false,
  },
  createCategoryFormData: {
    type: Object as PropType<CreateCategoryRequest>,
    default: (): CreateCategoryRequest => ({
      name: '',
    }),
  },
  updateCategoryFormData: {
    type: Object as PropType<UpdateCategoryRequest>,
    default: (): UpdateCategoryRequest => ({
      name: '',
    }),
  },
  categoryTableItemsPerPage: {
    type: Number,
    default: 20,
  },
  categoryTableItemsTotal: {
    type: Number,
    default: 0,
  },
  category: {
    type: Object as PropType<Category>,
    default: (): Category => ({
      id: '',
      name: '',
      createdAt: 0,
      updatedAt: 0,
    }),
  },
  categories: {
    type: Array<Category>,
    default: () => [],
  },
  // ProductType
  createProductTypeDialog: {
    type: Boolean,
    default: false,
  },
  updateProductTypeDialog: {
    type: Boolean,
    default: false,
  },
  deleteProductTypeDialog: {
    type: Boolean,
    default: false,
  },
  createProductTypeFormData: {
    type: Object as PropType<CreateProductTypeRequest>,
    default: () => ({
      name: '',
      iconUrl: '',
    }),
  },
  updateProductTypeFormData: {
    type: Object as PropType<UpdateProductTypeRequest>,
    default: () => ({
      name: '',
      iconUrl: '',
    }),
  },
  createProductTypeIconUploadStatus: {
    type: Object as PropType<ImageUploadStatus>,
    default: (): ImageUploadStatus => ({
      error: false,
      message: '',
    }),
  },
  updateProductTypeIconUploadStatus: {
    type: Object as PropType<ImageUploadStatus>,
    default: (): ImageUploadStatus => ({
      error: false,
      message: '',
    }),
  },
  productTypeTableItemsPerPage: {
    type: Number,
    default: 20,
  },
  productTypeTableItemsTotal: {
    type: Number,
    default: 0,
  },
  productType: {
    type: Object as PropType<ProductType>,
    default: (): ProductType => ({
      id: '',
      categoryId: '',
      name: '',
      iconUrl: '',
      createdAt: 0,
      updatedAt: 0,
    }),
  },
  productTypes: {
    type: Array<ProductType>,
    default: () => [],
  },
})

const emits = defineEmits<{
  (e: 'click:new-category'): void
  (e: 'click:edit-category', categoryId: string): void
  (e: 'click:delete-category', categoryId: string): void
  (e: 'click:new-product-type'): void
  (e: 'click:edit-product-type', productTypeId: string): void
  (e: 'click:delete-product-type', productTypeId: string): void
  (e: 'update:selected-tab-item', item: string): void
  (e: 'update:create-category-dialog', v: boolean): void
  (e: 'update:update-category-dialog', v: boolean): void
  (e: 'update:delete-category-dialog', v: boolean): void
  (e: 'update:create-category-form-data', formData: CreateCategoryRequest): void
  (e: 'update:update-category-form-data', formData: UpdateCategoryRequest): void
  (e: 'update:create-product-type-dialog', v: boolean): void
  (e: 'update:update-product-type-dialog', v: boolean): void
  (e: 'update:delete-product-type-dialog', v: boolean): void
  (e: 'update:create-product-type-form-data', formData: CreateProductTypeRequest): void
  (e: 'update:update-product-type-form-data', formData: UpdateProductTypeRequest): void
  (e: 'update:create-product-type-icon', files: FileList): void
  (e: 'update:update-product-type-icon', files: FileList): void
  (e: 'update:category-page', page: number): void
  (e: 'update:category-items-per-page', page: number): void
  (e: 'update:product-type-page', page: number): void
  (e: 'update:product-type-items-per-page', page: number): void
  (e: 'search:category', name: string): void
  (e: 'submit:create-category'): void
  (e: 'submit:update-category'): void
  (e: 'submit:delete-category'): void
  (e: 'submit:create-product-type', categoryId: string): void
  (e: 'submit:update-product-type'): void
  (e: 'submit:delete-product-type'): void
}>()

const tabs: VTabs[] = [
  { name: 'カテゴリー', value: 'categories' },
  { name: '品目', value: 'productTypes' },
]

const selectedTabItemValue = computed({
  get: (): string => props.selectedTabItem,
  set: (item: string): void => emits('update:selected-tab-item', item),
})
const createCategoryDialogValue = computed({
  get: (): boolean => props.createCategoryDialog,
  set: (v: boolean): void => emits('update:create-category-dialog', v),
})
const updateCategoryDialogValue = computed({
  get: (): boolean => props.updateCategoryDialog,
  set: (v: boolean): void => emits('update:update-category-dialog', v),
})
const deleteCategoryDialogValue = computed({
  get: (): boolean => props.deleteCategoryDialog,
  set: (v: boolean): void => emits('update:delete-category-dialog', v),
})
const createCategoryFormDataValue = computed({
  get: (): CreateCategoryRequest => props.createCategoryFormData,
  set: (formData: CreateCategoryRequest): void => emits('update:create-category-form-data', formData),
})
const updateCategoryFormDataValue = computed({
  get: (): UpdateCategoryRequest => props.updateCategoryFormData,
  set: (formData: UpdateCategoryRequest): void => emits('update:update-category-form-data', formData),
})
const createProductTypeDialogValue = computed({
  get: (): boolean => props.createProductTypeDialog,
  set: (v: boolean): void => emits('update:create-product-type-dialog', v),
})
const updateProductTypeDialogValue = computed({
  get: (): boolean => props.updateProductTypeDialog,
  set: (v: boolean): void => emits('update:update-product-type-dialog', v),
})
const deleteProductTypeDialogValue = computed({
  get: (): boolean => props.deleteProductTypeDialog,
  set: (v: boolean): void => emits('update:delete-product-type-dialog', v),
})
const createProductTypeFormDataValue = computed({
  get: (): CreateProductTypeRequest => props.createProductTypeFormData,
  set: (formData: CreateProductTypeRequest): void => emits('update:create-product-type-form-data', formData),
})
const updateProductTypeFormDataValue = computed({
  get: (): UpdateProductTypeRequest => props.updateProductTypeFormData,
  set: (formData: UpdateProductTypeRequest): void => emits('update:update-product-type-form-data', formData),
})

const onClickNewCategory = (): void => {
  emits('click:new-category')
}

const onClickEditCategory = (categoryId: string): void => {
  emits('click:edit-category', categoryId)
}

const onClickDeleteCategory = (categoryId: string): void => {
  emits('click:delete-category', categoryId)
}

const onClickNewProductType = (): void => {
  emits('click:new-product-type')
}

const onClickEditProductType = (productTypeId: string): void => {
  emits('click:edit-product-type', productTypeId)
}

const onClickDeleteProductType = (productTypeId: string): void => {
  emits('click:delete-product-type', productTypeId)
}

const onChangeCreateProductTypeIcon = (files: FileList): void => {
  emits('update:create-product-type-icon', files)
}

const onChangeUpdateProductTypeIcon = (files: FileList): void => {
  emits('update:create-product-type-icon', files)
}

const onClickCategoryPage = (page: number): void => {
  emits('update:category-page', page)
}

const onClickCategoryItemsPerPage = (page: number): void => {
  emits('update:category-items-per-page', page)
}

const onClickProductTypePage = (page: number): void => {
  emits('update:product-type-page', page)
}

const onClickProductTypeItemsPerPage = (page: number): void => {
  emits('update:product-type-items-per-page', page)
}

const onSearchCategory = (name: string): void => {
  emits('search:category', name)
}

const onSubmitCreateCategory = (): void => {
  emits('submit:create-category')
}

const onSubmitUpdateCategory = (): void => {
  emits('submit:update-category')
}

const onSubmitDeleteCategory = (): void => {
  emits('submit:delete-category')
}

const onSubmitCreateProductType = (categoryId: string): void => {
  emits('submit:create-product-type', categoryId)
}

const onSubmitUpdateProductType = (): void => {
  emits('submit:update-product-type')
}

const onSubmitDeleteProductType = (): void => {
  emits('submit:delete-product-type')
}
</script>

<template>
  <atoms-app-alert
    :show="props.isAlert"
    :type="props.alertType"
    :text="props.alertText"
  />

  <v-card>
    <v-card-title class="d-flex flex-row">
      カテゴリー・品目設定
    </v-card-title>

    <v-card-text>
      <v-tabs
        v-model="selectedTabItemValue"
        grow
        color="dark"
      >
        <v-tab
          v-for="item in tabs"
          :key="item.value"
          :value="item.value"
        >
          {{ item.name }}
        </v-tab>
      </v-tabs>
    </v-card-text>

    <v-window v-model="selectedTabItemValue">
      <v-window-item value="categories">
        <organisms-category-list
          v-model:create-form-data="createCategoryFormDataValue"
          v-model:update-form-data="updateCategoryFormDataValue"
          v-model:create-dialog="createCategoryDialogValue"
          v-model:update-dialog="updateCategoryDialogValue"
          v-model:delete-dialog="deleteCategoryDialogValue"
          :loading="loading"
          :admin-type="adminType"
          :category="category"
          :categories="categories"
          :table-items-per-page="categoryTableItemsPerPage"
          :table-items-total="categoryTableItemsTotal"
          @click:new="onClickNewCategory"
          @click:edit="onClickEditCategory"
          @click:delete="onClickDeleteCategory"
          @update:page="onClickCategoryPage"
          @update:items-per-page="onClickCategoryItemsPerPage"
          @submit:create="onSubmitCreateCategory"
          @submit:update="onSubmitUpdateCategory"
          @submit:delete="onSubmitDeleteCategory"
        />
      </v-window-item>

      <v-window-item value="productTypes">
        <organisms-product-type-list
          v-model:create-form-data="createProductTypeFormDataValue"
          v-model:update-form-data="updateProductTypeFormDataValue"
          v-model:create-dialog="createProductTypeDialogValue"
          v-model:update-dialog="updateProductTypeDialogValue"
          v-model:delete-dialog="deleteProductTypeDialogValue"
          :loading="loading"
          :admin-type="adminType"
          :product-type="productType"
          :product-types="productTypes"
          :categories="categories"
          :table-items-per-page="productTypeTableItemsPerPage"
          :table-items-total="productTypeTableItemsTotal"
          :create-icon-upload-status="createProductTypeIconUploadStatus"
          :update-icon-upload-status="updateProductTypeIconUploadStatus"
          @click:new="onClickNewProductType"
          @click:edit="onClickEditProductType"
          @click:delete="onClickDeleteProductType"
          @update:page="onClickProductTypePage"
          @update:items-per-page="onClickProductTypeItemsPerPage"
          @update:create-icon="onChangeCreateProductTypeIcon"
          @update:update-icon="onChangeUpdateProductTypeIcon"
          @search:category="onSearchCategory"
          @submit:create="onSubmitCreateProductType"
          @submit:update="onSubmitUpdateProductType"
          @submit:delete="onSubmitDeleteProductType"
        />
      </v-window-item>
    </v-window>
  </v-card>
</template>
