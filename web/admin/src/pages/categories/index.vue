<template>
  <div>
    <v-card-title>カテゴリー・品目設定</v-card-title>
    <v-tabs v-model="selector" grow color="dark">
      <v-tabs-slider color="accent"></v-tabs-slider>
      <v-tab
        v-for="item in items"
        :key="item.value"
        :href="`#tab-${item.value}`"
      >
        {{ item.name }}
      </v-tab>
    </v-tabs>

    <v-tabs-items v-model="selector">
      <v-tab-item value="tab-categories">
        <v-dialog v-model="categoryDialog" width="500">
          <template #activator="{ on, attrs }">
            <div class="d-flex pt-3 pr-3">
              <v-spacer />
              <v-btn outlined color="primary" v-bind="attrs" v-on="on">
                <v-icon left>mdi-plus</v-icon>
                追加
              </v-btn>
            </div>
          </template>

          <v-card>
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
              <v-spacer></v-spacer>
              <v-btn color="accentDarken" text @click="categoryCancel">
                キャンセル
              </v-btn>
              <v-btn color="primary" outlined @click="categoryRegister">
                登録
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
        <the-category-list
          :loading="fetchState.pending"
          :table-footer-props="categoriesOptions"
          @update:items-per-page="handleUpdateCategoriesItemsPerPage"
          @update:page="handleUpdateCategoriesPage"
        />
      </v-tab-item>

      <v-tab-item value="tab-categoryItems">
        <v-dialog v-model="productTypeDialog" width="500">
          <template #activator="{ on, attrs }">
            <div class="d-flex pt-3 pr-3">
              <v-spacer />
              <v-btn outlined color="primary" v-bind="attrs" v-on="on">
                <v-icon left>mdi-plus</v-icon>
                追加
              </v-btn>
            </div>
          </template>
          <v-card>
            <v-card-title class="primaryLight"> 品目登録 </v-card-title>
            <v-card-text class="mt-4">
              <v-autocomplete
                v-model="selectedCategoryId"
                :items="categoriesItems.categories"
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
                v-model="productTypeFormData.name"
                maxlength="32"
                label="品目"
              />
            </v-card-text>
            <v-divider />

            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="accentDarken" text @click="productTypeCancel">
                キャンセル
              </v-btn>
              <v-btn color="primary" outlined @click="productTypeRegister">
                登録
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
        <the-product-type-list
          :loading="fetchState.pending"
          :table-footer-props="productTypesOptions"
          @update:items-per-page="handleUpdateProductTypesItemsPerPage"
          @update:page="handleUpdateProductTypesPage"
        />
      </v-tab-item>
    </v-tabs-items>
  </div>
</template>

<script lang="ts">
import {
  defineComponent,
  reactive,
  ref,
  useFetch,
  watch,
} from '@nuxtjs/composition-api'

import TheCategoryList from '~/components/organisms/TheCategoryList.vue'
import TheProductTypeList from '~/components/organisms/TheProductTypeList.vue'
import { usePagination } from '~/lib/hooks'
import { useAuthStore } from '~/store/auth'
import { useCategoryStore } from '~/store/category'
import { useProductTypeStore } from '~/store/product-type'
import {
  CategoriesResponseCategoriesInner,
  CreateCategoryRequest,
  CreateProductTypeRequest,
} from '~/types/api'
import { Category } from '~/types/props/category'

export default defineComponent({
  components: {
    TheCategoryList,
    TheProductTypeList,
  },

  setup() {
    const categoryStore = useCategoryStore()
    const productTypeStore = useProductTypeStore()
    const { accessToken } = useAuthStore()

    const categoriesItems = reactive<{
      offset: number
      categories: CategoriesResponseCategoriesInner[]
    }>({ offset: 0, categories: [] })

    const selector = ref<string>('categories')
    const categoryDialog = ref<boolean>(false)
    const productTypeDialog = ref<boolean>(false)
    const selectedCategoryId = ref<string>('')
    const items: Category[] = [
      { name: 'カテゴリー', value: 'categories' },
      { name: '品目', value: 'categoryItems' },
    ]

    const categoryFormData = reactive<CreateCategoryRequest>({
      name: '',
    })

    const productTypeFormData = reactive<CreateProductTypeRequest>({
      name: '',
      iconUrl: '',
    })

    const categoryCancel = (): void => {
      categoryDialog.value = false
    }

    const productTypeCancel = (): void => {
      productTypeDialog.value = false
    }

    const categoryRegister = async (): Promise<void> => {
      try {
        await categoryStore.createCategory(categoryFormData)
        categoryDialog.value = false
      } catch (error) {
        console.log(error)
      }
    }

    const productTypeRegister = async (): Promise<void> => {
      try {
        await productTypeStore.createProductType(
          selectedCategoryId.value,
          productTypeFormData
        )
        productTypeDialog.value = false
      } catch (error) {
        console.log(error)
      }
    }

    const {
      itemsPerPage: categoriesItemsPerPage,
      offset: categoriesOffset,
      options: categoriesOptions,
      handleUpdateItemsPerPage: handleUpdateCategoriesItemsPerPage,
      updateCurrentPage: _handleUpdateCategoriesPage,
    } = usePagination()

    watch(categoriesItemsPerPage, () => {
      categoryStore.fetchCategories(categoriesItemsPerPage.value, 0)
    })

    const handleUpdateCategoriesPage = async (page: number) => {
      _handleUpdateCategoriesPage(page)
      await categoryStore.fetchCategories(
        categoriesItemsPerPage.value,
        categoriesOffset.value
      )
    }

    const {
      itemsPerPage: productTypesItemsPerPage,
      offset: productTypesOffset,
      options: productTypesOptions,
      handleUpdateItemsPerPage: handleUpdateProductTypesItemsPerPage,
      updateCurrentPage: _handleUpdateProductTypesPage,
    } = usePagination()

    const { fetchState } = useFetch(async () => {
      try {
        await Promise.all([
          categoryStore.fetchCategories(categoriesItemsPerPage.value),
          productTypeStore.fetchProductTypes(productTypesItemsPerPage.value),
        ])
        categoriesItems.categories = categoryStore.categories
      } catch (err) {
        console.log(err)
      }
    })

    watch(productTypesItemsPerPage, () => {
      productTypeStore.fetchProductTypes(productTypesItemsPerPage.value)
    })

    const handleUpdateProductTypesPage = async (page: number) => {
      _handleUpdateProductTypesPage(page)
      await productTypeStore.fetchProductTypes(
        productTypesItemsPerPage.value,
        productTypesOffset.value
      )
    }

    const handleMoreCategoryItems = async () => {
      if (accessToken) {
        const limit = 20
        categoriesItems.offset = categoriesItems.offset + limit + 1
        const res = await categoryStore
          .apiClient(accessToken)
          .v1ListCategories(limit, categoriesItems.offset)
        categoriesItems.categories.push(...res.data.categories)
      }
    }

    return {
      fetchState,
      categoriesOptions,
      handleUpdateCategoriesItemsPerPage,
      handleUpdateCategoriesPage,
      productTypesOptions,
      handleUpdateProductTypesItemsPerPage,
      handleUpdateProductTypesPage,
      categoriesItems,
      items,
      selector,
      categoryDialog,
      categoryFormData,
      productTypeFormData,
      productTypeDialog,
      selectedCategoryId,
      categoryCancel,
      productTypeCancel,
      categoryRegister,
      productTypeRegister,
      handleMoreCategoryItems,
    }
  },
})
</script>
