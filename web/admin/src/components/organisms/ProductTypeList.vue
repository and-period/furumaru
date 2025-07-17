<script lang="ts" setup>
import { mdiAccount, mdiPencil, mdiDelete, mdiPlus } from '@mdi/js'
import useVuelidate from '@vuelidate/core'
import type { VDataTable } from 'vuetify/lib/components/index.mjs'

import { AdminType } from '~/types/api'
import type { Category, CreateProductTypeRequest, ProductType, UpdateProductTypeRequest } from '~/types/api'
import type { ImageUploadStatus } from '~/types/props'
import { getErrorMessage } from '~/lib/validations'
import { getResizedImages } from '~/lib/helpers'
import { CreateProductTypeValidationRules, UpdateProductTypeValidationRules } from '~/types/validations'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  adminType: {
    type: Number as PropType<AdminType>,
    default: AdminType.UNKNOWN,
  },
  createDialog: {
    type: Boolean,
    default: false,
  },
  updateDialog: {
    type: Boolean,
    default: false,
  },
  deleteDialog: {
    type: Boolean,
    default: false,
  },
  createFormData: {
    type: Object as PropType<CreateProductTypeRequest>,
    default: (): CreateProductTypeRequest => ({
      name: '',
      iconUrl: '',
    }),
  },
  updateFormData: {
    type: Object as PropType<UpdateProductTypeRequest>,
    default: (): UpdateProductTypeRequest => ({
      name: '',
      iconUrl: '',
    }),
  },
  createIconUploadStatus: {
    type: Object as PropType<ImageUploadStatus>,
    default: (): ImageUploadStatus => ({
      error: false,
      message: '',
    }),
  },
  updateIconUploadStatus: {
    type: Object as PropType<ImageUploadStatus>,
    default: (): ImageUploadStatus => ({
      error: false,
      message: '',
    }),
  },
  categories: {
    type: Array<Category>,
    default: () => [],
  },
  productType: {
    type: Object as PropType<ProductType>,
    default: (): ProductType => ({
      id: '',
      categoryId: '',
      name: '',
      iconUrl: '',
      icons: [],
      createdAt: 0,
      updatedAt: 0,
    }),
  },
  productTypes: {
    type: Array<ProductType>,
    default: () => [],
  },
  tableItemsPerPage: {
    type: Number,
    default: 20,
  },
  tableItemsTotal: {
    type: Number,
    default: 0,
  },
})

const emit = defineEmits<{
  (e: 'click:new'): void
  (e: 'click:edit', productTypeId: string): void
  (e: 'click:delete', productTypeId: string): void
  (e: 'update:create-dialog', v: boolean): void
  (e: 'update:update-dialog', v: boolean): void
  (e: 'update:delete-dialog', v: boolean): void
  (e: 'update:create-form-data', formData: CreateProductTypeRequest): void
  (e: 'update:update-form-data', formData: UpdateProductTypeRequest): void
  (e: 'update:product-type', productType: ProductType): void
  (e: 'update:page', page: number): void
  (e: 'update:items-per-page', page: number): void
  (e: 'update:create-icon', files: FileList): void
  (e: 'update:update-icon', files: FileList): void
  (e: 'search:category', name: string): void
  (e: 'submit:create', categoryId: string): void
  (e: 'submit:update'): void
  (e: 'submit:delete'): void
}>()

const headers: VDataTable['headers'] = [
  {
    title: 'アイコン',
    key: 'icon',
    sortable: false,
  },
  {
    title: 'カテゴリー',
    key: 'category',
    sortable: false,
  },
  {
    title: '品目',
    key: 'name',
    sortable: false,
  },
  {
    title: '',
    key: 'actions',
    align: 'end',
    sortable: false,
  },
]

const selectedCategoryId = ref<string>()

const productTypeValue = computed({
  get: (): ProductType => props.productType,
  set: (productType: ProductType): void => emit('update:product-type', productType),
})
const createDialogValue = computed({
  get: (): boolean => props.createDialog,
  set: (val: boolean): void => emit('update:create-dialog', val),
})
const updateDialogValue = computed({
  get: (): boolean => props.updateDialog,
  set: (val: boolean): void => emit('update:update-dialog', val),
})
const deleteDialogValue = computed({
  get: (): boolean => props.deleteDialog,
  set: (val: boolean): void => emit('update:delete-dialog', val),
})
const createFormDataValue = computed({
  get: (): CreateProductTypeRequest => props.createFormData,
  set: (formData: CreateProductTypeRequest): void => emit('update:create-form-data', formData),
})
const updateFormDataValue = computed({
  get: (): UpdateProductTypeRequest => props.updateFormData,
  set: (formData: UpdateProductTypeRequest): void => emit('update:update-form-data', formData),
})

const createFormDataValidate = useVuelidate(CreateProductTypeValidationRules, createFormDataValue)
const updateFormDataValidate = useVuelidate(UpdateProductTypeValidationRules, updateFormDataValue)

const isRegisterable = (): boolean => {
  return props.adminType === AdminType.ADMINISTRATOR
}

const isEditable = (): boolean => {
  return props.adminType === AdminType.ADMINISTRATOR
}

const getCategoryName = (categoryId: string): string => {
  const category = props.categories.find((category: Category): boolean => {
    return category.id === categoryId
  })
  return category ? category.name : ''
}

const getIcons = (productType: ProductType): string => {
  if (!productType.icons) {
    return ''
  }
  return getResizedImages(productType.iconUrl)
}

const onClickNew = (): void => {
  emit('click:new')
}

const onClickCloseCreateDialog = (): void => {
  createDialogValue.value = false
}

const onClickEdit = (categoryId: string): void => {
  emit('click:edit', categoryId)
}

const onClickCloseUpdateDialog = (): void => {
  updateDialogValue.value = false
}

const onClickDelete = (categoryId: string): void => {
  emit('click:delete', categoryId)
}

const onClickCloseDeleteDialog = (): void => {
  deleteDialogValue.value = false
}

const onClickUpdatePage = (page: number) => {
  emit('update:page', page)
}

const onClickUpdateItemsPerPage = (page: number) => {
  emit('update:items-per-page', page)
}

const onChangeCreateIconFile = (files?: FileList): void => {
  if (!files) {
    return
  }
  emit('update:create-icon', files)
}

const onChangeUpdateIconFile = (files?: FileList): void => {
  if (!files) {
    return
  }
  emit('update:update-icon', files)
}

const onSearchCategory = (name: string): void => {
  emit('search:category', name)
}

const onSubmitCreate = async (): Promise<void> => {
  const valid = await createFormDataValidate.value.$validate()
  if (!valid) {
    return
  }
  emit('submit:create', selectedCategoryId.value || '')
}

const onSubmitUpdate = async (): Promise<void> => {
  const valid = await updateFormDataValidate.value.$validate()
  if (!valid) {
    return
  }
  emit('submit:update')
}

const onSubmitDelete = (): void => {
  emit('submit:delete')
}
</script>

<template>
  <v-dialog
    v-model="createDialogValue"
    width="500"
  >
    <v-card>
      <v-card-title class="primaryLight">
        品目登録
      </v-card-title>
      <v-card-text class="mt-6">
        <v-autocomplete
          v-model="selectedCategoryId"
          label="カテゴリ"
          :items="categories"
          item-title="name"
          item-value="id"
          @update:search="onSearchCategory"
        />
        <v-text-field
          v-model="createFormDataValidate.name.$model"
          :error-messages="getErrorMessage(createFormDataValidate.name.$errors)"
          label="カテゴリー"
          maxlength="32"
        />
        <molecules-image-select-form
          label="アイコン"
          :loading="loading"
          :img-url="createFormDataValue.iconUrl"
          :error="props.createIconUploadStatus.error"
          :message="props.createIconUploadStatus.message"
          @update:file="onChangeCreateIconFile"
        />
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn
          color="error"
          variant="text"
          @click="onClickCloseCreateDialog"
        >
          キャンセル
        </v-btn>
        <v-btn
          :loading="loading"
          color="primary"
          variant="outlined"
          @click="onSubmitCreate"
        >
          登録
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog
    v-model="updateDialogValue"
    width="500"
  >
    <v-card>
      <v-card-title class="primaryLight">
        カテゴリー編集
      </v-card-title>
      <v-card-text class="mt-6">
        <v-autocomplete
          v-model="productTypeValue.categoryId"
          label="カテゴリ"
          :items="categories"
          item-title="name"
          item-value="id"
          readonly
        />
        <v-text-field
          v-model="updateFormDataValidate.name.$model"
          :error-messages="getErrorMessage(updateFormDataValidate.name.$errors)"
          label="カテゴリー"
          maxlength="32"
        />
        <molecules-image-select-form
          label="アイコン"
          :loading="loading"
          :img-url="updateFormDataValue.iconUrl"
          :error="props.updateIconUploadStatus.error"
          :message="props.updateIconUploadStatus.message"
          @update:file="onChangeUpdateIconFile"
        />
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn
          color="error"
          variant="text"
          @click="onClickCloseUpdateDialog"
        >
          キャンセル
        </v-btn>
        <v-btn
          :loading="loading"
          color="primary"
          variant="outlined"
          @click="onSubmitUpdate"
        >
          編集
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog
    v-model="deleteDialogValue"
    width="500"
  >
    <v-card>
      <v-card-title class="text-h7">
        {{ props.productType?.name || '' }}を本当に削除しますか？
      </v-card-title>
      <v-card-actions>
        <v-spacer />
        <v-btn
          color="error"
          variant="text"
          @click="onClickCloseDeleteDialog"
        >
          キャンセル
        </v-btn>
        <v-btn
          :loading="loading"
          color="primary"
          variant="outlined"
          @click="onSubmitDelete"
        >
          削除
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-card>
    <v-card-title class="d-flex flex-row">
      <v-spacer />
      <v-btn
        v-show="isRegisterable()"
        variant="outlined"
        color="primary"
        @click="onClickNew"
      >
        <v-icon :icon="mdiPlus" />
        品目登録
      </v-btn>
    </v-card-title>

    <v-card-text>
      <v-data-table-server
        :headers="headers"
        :loading="loading"
        :items="productTypes"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
        no-data-text="登録されている品目情報がありません"
        @update:page="onClickUpdatePage"
        @update:items-per-page="onClickUpdateItemsPerPage"
      >
        <template #[`item.icon`]="{ item }">
          <v-avatar>
            <v-img
              v-if="item.iconUrl !== ''"
              cover
              :src="item.iconUrl"
              :srcset="getIcons(item)"
            />
            <v-icon
              v-else
              :icon="mdiAccount"
            />
          </v-avatar>
        </template>
        <template #[`item.category`]="{ item }">
          {{ getCategoryName(item.categoryId) }}
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn
            v-show="isEditable()"
            class="mr-2"
            variant="outlined"
            color="primary"
            size="small"
            @click="onClickEdit(item.id)"
          >
            <v-icon
              size="small"
              :icon="mdiPencil"
            />
            編集
          </v-btn>
          <v-btn
            v-show="isEditable()"
            variant="outlined"
            color="primary"
            size="small"
            @click="onClickDelete(item.id)"
          >
            <v-icon
              size="small"
              :icon="mdiDelete"
            />
            削除
          </v-btn>
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
