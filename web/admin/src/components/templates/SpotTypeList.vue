<script lang="ts" setup>
import { mdiDelete, mdiPencil, mdiPlus } from '@mdi/js'
import type { VDataTable } from 'vuetify/components'
import useVuelidate from '@vuelidate/core'

import type { AlertType } from '~/lib/hooks'
import { AdminType } from '~/types/api/v1'
import type { CreateSpotTypeRequest, SpotType, UpdateSpotTypeRequest } from '~/types/api/v1'
import { getErrorMessage } from '~/lib/validations'
import { CreateSpotTypeValidationRules, UpdateSpotTypeValidationRules } from '~/types/validations'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  adminType: {
    type: Number as PropType<AdminType>,
    default: AdminType.AdminTypeUnknown,
  },
  newDialog: {
    type: Boolean,
    default: false,
  },
  editDialog: {
    type: Boolean,
    default: false,
  },
  deleteDialog: {
    type: Boolean,
    default: false,
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
  tableItemsPerPage: {
    type: Number,
    default: 20,
  },
  tableItemsTotal: {
    type: Number,
    default: 0,
  },
  spotTypes: {
    type: Array<SpotType>,
    default: () => [],
  },
  newFormData: {
    type: Object as PropType<CreateSpotTypeRequest>,
    default: (): CreateSpotTypeRequest => ({
      name: '',
    }),
  },
  editFormData: {
    type: Object as PropType<UpdateSpotTypeRequest>,
    default: (): UpdateSpotTypeRequest => ({
      name: '',
    }),
  },
})

const emit = defineEmits<{
  (e: 'click:update-page', page: number): void
  (e: 'click:update-items-per-page', page: number): void
  (e: 'update:new-dialog', toggle: boolean): void
  (e: 'update:edit-dialog', toggle: boolean): void
  (e: 'update:delete-dialog', toggle: boolean): void
  (e: 'update:new-form-data', formData: CreateSpotTypeRequest): void
  (e: 'update:edit-form-data', formData: UpdateSpotTypeRequest): void
  (e: 'submit:create'): void
  (e: 'submit:update', spotTypeId: string): void
  (e: 'submit:delete', spotTypeId: string): void
}>()

const headers: VDataTable['headers'] = [
  {
    title: 'スポット種別',
    key: 'name',
    sortable: false,
  },
  {
    title: '',
    key: 'actions',
    width: 200,
    align: 'end',
    sortable: false,
  },
]

const selectedItem = ref<SpotType>()

const newDialogValue = computed({
  get: (): boolean => props.newDialog,
  set: (toggle: boolean): void => emit('update:new-dialog', toggle),
})
const editDialogValue = computed({
  get: (): boolean => props.editDialog,
  set: (toggle: boolean): void => emit('update:edit-dialog', toggle),
})
const deleteDialogValue = computed({
  get: (): boolean => props.deleteDialog,
  set: (toggle: boolean): void => emit('update:delete-dialog', toggle),
})
const newFormDataValue = computed({
  get: (): CreateSpotTypeRequest => props.newFormData,
  set: (formData: CreateSpotTypeRequest): void => emit('update:new-form-data', formData),
})
const editFormDataValue = computed({
  get: (): UpdateSpotTypeRequest => props.editFormData,
  set: (formData: UpdateSpotTypeRequest): void => emit('update:edit-form-data', formData),
})

const newValidate = useVuelidate<CreateSpotTypeRequest>(CreateSpotTypeValidationRules, newFormDataValue)
const editValidate = useVuelidate<UpdateSpotTypeRequest>(UpdateSpotTypeValidationRules, editFormDataValue)

const isRegisterable = (): boolean => {
  return props.adminType === AdminType.AdminTypeAdministrator
}

const isEditable = (): boolean => {
  return props.adminType === AdminType.AdminTypeAdministrator
}

const onClickUpdatePage = (page: number): void => {
  emit('click:update-page', page)
}

const onClickUpdateItemsPerPage = (page: number): void => {
  emit('click:update-items-per-page', page)
}

const onClickCloseNewDialog = (): void => {
  newDialogValue.value = false
}

const onClickCloseEditDialog = (): void => {
  editDialogValue.value = false
}

const onClickCloseDeleteDialog = (): void => {
  deleteDialogValue.value = false
}

const onClickAdd = (): void => {
  newDialogValue.value = true
}

const submitCreate = (): void => {
  const valid = newValidate.value.$validate()
  if (!valid) {
    return
  }

  emit('submit:create')
}

const onClickEdit = (item: SpotType): void => {
  selectedItem.value = item
  editDialogValue.value = true
}

const submitUpdate = (): void => {
  const valid = editValidate.value.$validate()
  if (!valid) {
    return
  }

  emit('submit:update', selectedItem.value?.id || '')
}

const onClickDelete = (item: SpotType): void => {
  selectedItem.value = item
  deleteDialogValue.value = true
}

const submitDelete = (): void => {
  emit('submit:delete', selectedItem.value?.id || '')
}
</script>

<template>
  <v-alert
    v-show="props.isAlert"
    :type="props.alertType"
    v-text="props.alertText"
  />

  <v-dialog
    v-model="newDialogValue"
    width="500"
  >
    <v-card>
      <v-card-title class="text-h6 primaryLight">
        スポット種別登録
      </v-card-title>
      <v-card-text>
        <v-text-field
          v-model="newValidate.name.$model"
          :error-messages="getErrorMessage(newValidate.name.$errors)"
          class="mx-4"
          label="スポット種別名"
        />
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn
          color="error"
          variant="text"
          @click="onClickCloseNewDialog"
        >
          キャンセル
        </v-btn>
        <v-btn
          :loading="loading"
          color="primary"
          variant="outlined"
          @click="submitCreate"
        >
          登録
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog
    v-model="editDialogValue"
    width="500"
  >
    <v-card>
      <v-card-title class="text-h6 primaryLight">
        スポット種別編集
      </v-card-title>
      <v-card-text>
        <v-text-field
          v-model="editValidate.name.$model"
          :error-messages="getErrorMessage(editValidate.name.$errors)"
          class="mx-4"
          label="スポット種別名"
        />
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn
          color="error"
          variant="text"
          @click="onClickCloseEditDialog"
        >
          キャンセル
        </v-btn>
        <v-btn
          :loading="loading"
          color="primary"
          variant="outlined"
          @click="submitUpdate"
        >
          更新
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog
    v-model="deleteDialogValue"
    width="500"
  >
    <v-card>
      <v-card-title>
        {{ selectedItem?.name || '' }}を本当に削除しますか？
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
          @click="submitDelete"
        >
          削除
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-card
    class="mt-4"
    flat
  >
    <v-card-title class="d-flex flex-row">
      スポット種別管理
      <v-spacer />
      <v-btn
        v-show="isRegisterable()"
        variant="outlined"
        color="primary"
        @click="onClickAdd"
      >
        <v-icon
          start
          :icon="mdiPlus"
        />
        スポット種別登録
      </v-btn>
    </v-card-title>

    <v-card-text>
      <v-data-table-server
        :headers="headers"
        :loading="loading"
        :items="props.spotTypes"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
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
            @click="onClickEdit(item)"
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
            @click="onClickDelete(item)"
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
