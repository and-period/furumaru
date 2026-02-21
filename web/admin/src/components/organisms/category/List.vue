<script lang="ts" setup>
import { mdiDelete, mdiPencil, mdiPlus } from '@mdi/js'
import useVuelidate from '@vuelidate/core'
import type { VDataTable } from 'vuetify/components'

import { AdminType } from '~/types/api/v1'
import type { Category, CreateCategoryRequest, UpdateCategoryRequest } from '~/types/api/v1'
import { getErrorMessage } from '~/lib/validations'
import { CreateCategoryValidationRules, UpdateCategoryValidationRules } from '~/types/validations'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  adminType: {
    type: Number as PropType<AdminType>,
    default: AdminType.AdminTypeUnknown,
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
    type: Object as PropType<CreateCategoryRequest>,
    default: (): CreateCategoryRequest => ({
      name: '',
    }),
  },
  updateFormData: {
    type: Object as PropType<UpdateCategoryRequest>,
    default: (): UpdateCategoryRequest => ({
      name: '',
    }),
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
  (e: 'click:edit', categoryId: string): void
  (e: 'click:delete', categoryId: string): void
  (e: 'update:create-dialog', v: boolean): void
  (e: 'update:update-dialog', v: boolean): void
  (e: 'update:delete-dialog', v: boolean): void
  (e: 'update:create-form-data', formData: CreateCategoryRequest): void
  (e: 'update:update-form-data', formData: UpdateCategoryRequest): void
  (e: 'update:page', page: number): void
  (e: 'update:items-per-page', page: number): void
  (e: 'submit:create'): void
  (e: 'submit:update'): void
  (e: 'submit:delete'): void
}>()

const headers: VDataTable['headers'] = [
  {
    title: 'カテゴリー',
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
  get: (): CreateCategoryRequest => props.createFormData,
  set: (formData: CreateCategoryRequest): void => emit('update:create-form-data', formData),
})
const updateFormDataValue = computed({
  get: (): UpdateCategoryRequest => props.updateFormData,
  set: (formData: UpdateCategoryRequest): void => emit('update:update-form-data', formData),
})

const createFormDataValidate = useVuelidate(CreateCategoryValidationRules, createFormDataValue)
const updateFormDataValidate = useVuelidate(UpdateCategoryValidationRules, updateFormDataValue)

const isRegisterable = (): boolean => {
  return props.adminType === AdminType.AdminTypeAdministrator
}

const isEditable = (): boolean => {
  return props.adminType === AdminType.AdminTypeAdministrator
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

const onSubmitCreate = async (): Promise<void> => {
  const valid = await createFormDataValidate.value.$validate()
  if (!valid) {
    return
  }
  emit('submit:create')
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
        カテゴリー登録
      </v-card-title>
      <v-card-text class="mt-6">
        <v-text-field
          v-model="createFormDataValidate.name.$model"
          :error-messages="getErrorMessage(createFormDataValidate.name.$errors)"
          label="カテゴリー"
          maxlength="32"
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
        <v-text-field
          v-model="updateFormDataValidate.name.$model"
          :error-messages="getErrorMessage(updateFormDataValidate.name.$errors)"
          label="カテゴリー"
          maxlength="32"
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
        {{ props.category?.name || '' }}を本当に削除しますか？
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
        カテゴリ登録
      </v-btn>
    </v-card-title>

    <v-card-text>
      <v-data-table-server
        :headers="headers"
        :loading="loading"
        :items="categories"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
        no-data-text="登録されている商品種別情報がありません"
        @update:page="onClickUpdatePage"
        @update:items-per-page="onClickUpdateItemsPerPage"
      >
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
