<script lang="ts" setup>
import { DataTableHeader } from 'vuetify'

import { useProductTypeStore } from '~/store'
import {
  ProductTypesResponseProductTypesInner,
  UpdateProductTypeRequest,
  UploadImageResponse,
} from '~/types/api'
import { ImageUploadStatus } from '~/types/props'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  tableFooterProps: {
    type: Object,
    default: () => {},
  },
  categories: {
    type: Array,
    default: () => [],
  },
})

const emit = defineEmits<{
  (e: 'update:items-per-page', page: number): void
  (e: 'update:page', page: number): void
  (e: 'click:more-item'): void
}>()

const productTypeStore = useProductTypeStore()
const inputRef = ref<HTMLInputElement | null>(null)
const deleteDialog = ref<boolean>(false)
const editDialog = ref<boolean>(false)
const selectedCategoryId = ref<string>('')
const selectedItemId = ref<string>('')
const selectedName = ref<string>('')

const editFormData = reactive<UpdateProductTypeRequest>({
  name: '',
  iconUrl: '',
})

const productTypes = computed(() => {
  return productTypeStore.productTypes
})

const totalItems = computed(() => {
  return productTypeStore.totalItems
})

const headerUploadStatus = reactive<ImageUploadStatus>({
  error: false,
  message: '',
})

const productTypeHeaders: DataTableHeader[] = [
  {
    text: 'アイコン',
    value: 'icon',
  },
  {
    text: 'カテゴリー',
    value: 'category',
  },
  {
    text: '品目',
    value: 'productType',
  },
  {
    text: 'Actions',
    value: 'actions',
    width: 200,
    align: 'end',
    sortable: false,
  },
]

const handleUpdateItemsPerPage = (page: number) => {
  emit('update:items-per-page', page)
}

const handleUpdatePage = (page: number) => {
  emit('update:page', page)
}

const handleMoreCategoryItems = () => {
  emit('click:more-item')
}

const openEditDialog = (item: ProductTypesResponseProductTypesInner) => {
  editDialog.value = true
  selectedCategoryId.value = item.categoryId
  selectedItemId.value = item.id
  editFormData.name = item.name
  editFormData.iconUrl = item.iconUrl
}

const hideEditDialog = () => {
  editDialog.value = false
}

const handleEdit = async () => {
  try {
    await productTypeStore.editProductType(
      selectedCategoryId.value,
      selectedItemId.value,
      editFormData
    )
    editDialog.value = false
  } catch (error) {
    console.log(error)
  }
}

const openDeleteDialog = (
  item: ProductTypesResponseProductTypesInner
): void => {
  selectedCategoryId.value = item.categoryId
  selectedItemId.value = item.id
  selectedName.value = item.name
  deleteDialog.value = true
}

const handleDelete = async (): Promise<void> => {
  try {
    await productTypeStore.deleteProductType(
      selectedCategoryId.value,
      selectedItemId.value
    )
  } catch (error) {
    console.log(error)
  }
  deleteDialog.value = false
}

const hideDeleteDialog = (): void => {
  deleteDialog.value = false
}

const handleClick = () => {
  if (inputRef.value !== null) {
    inputRef.value.click()
  }
}

const handleInputFileChange = () => {
  const files = inputRef.value?.files
  if (inputRef.value && inputRef.value.files) {
    if (files && files.length > 0) {
      productTypeStore
        .uploadProductTypeIcon(files[0])
        .then((res: UploadImageResponse) => {
          editFormData.iconUrl = res.url
        })
        .catch(() => {
          headerUploadStatus.error = true
          headerUploadStatus.message = 'アップロードに失敗しました。'
        })
    }
  }
}
</script>

<template>
  <div>
    <v-data-table
      :headers="productTypeHeaders"
      :items="productTypes"
      :loading="props.loading"
      :server-items-length="totalItems"
      :footer-props="props.tableFooterProps"
      @update:items-per-page="handleUpdateItemsPerPage"
      @update:page="handleUpdatePage"
    >
      <template #[`item.icon`]="{ item }">
        <v-avatar>
          <img
            v-if="item.iconlUrl !== ''"
            :src="item.iconUrl"
            :alt="`${item.categoryName}-profile`"
          />
          <v-icon v-else>mdi-account</v-icon>
        </v-avatar>
      </template>
      <template #[`item.category`]="{ item }">
        {{ `${item.categoryName}` }}
      </template>
      <template #[`item.productType`]="{ item }">
        {{ `${item.name}` }}
      </template>
      <template #[`item.actions`]="{ item }">
        <v-btn outlined color="primary" small @click="openEditDialog(item)">
          <v-icon small>mdi-pencil</v-icon>
          編集
        </v-btn>
        <v-btn outlined color="primary" small @click="openDeleteDialog(item)">
          <v-icon small>mdi-delete</v-icon>
          削除
        </v-btn>
      </template>
    </v-data-table>

    <v-dialog v-model="editDialog" width="500">
      <v-card>
        <v-card-title class="primaryLight">品目編集</v-card-title>
        <v-card-text class="mt-4">
          <v-autocomplete
            v-model="selectedCategoryId"
            :items="props.categories"
            item-text="name"
            item-value="id"
            label="カテゴリー"
          >
            <template #append-item>
              <div class="pa-2">
                <v-btn
                  outlined
                  block
                  color="primary"
                  @click="handleMoreCategoryItems"
                >
                  <v-icon>mdi-plus</v-icon>
                  さらに読み込む
                </v-btn>
              </div>
            </template>
          </v-autocomplete>
          <v-spacer />
          <v-text-field
            v-model="editFormData.name"
            maxlength="32"
            label="品目"
          />
        </v-card-text>
        <v-card class="text-center" role="button" flat @click="handleClick">
          <v-card-text>
            <v-avatar size="96">
              <v-icon v-if="editFormData.iconUrl === ''" x-large
                >mdi-plus</v-icon
              >
              <v-img
                v-else
                :src="editFormData.iconUrl"
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
              @change="handleInputFileChange"
            />
            <p class="ma-0">アイコン画像を選択</p>
          </v-card-text>
        </v-card>
        <p v-show="headerUploadStatus.error" class="red--text ma-0">
          {{ headerUploadStatus.message }}
        </p>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="error" text @click="hideEditDialog"> キャンセル </v-btn>
          <v-btn color="primary" outlined @click="handleEdit"> 編集 </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="deleteDialog" width="500">
      <v-card>
        <v-card-title class="text-h7">
          {{ selectedName }}を本当に削除しますか？
        </v-card-title>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="error" text @click="hideDeleteDialog">
            キャンセル
          </v-btn>
          <v-btn color="primary" outlined @click="handleDelete"> 削除 </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>
