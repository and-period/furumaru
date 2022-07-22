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
            <v-divider></v-divider>

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
        <the-category-list />
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
            <v-card-title class="text-h6 primaryLight"> 品目登録 </v-card-title>
            <div>
              <v-select
                v-model="selectedCategoryId"
                class="mx-4"
                :items="categories"
                label="カテゴリー"
              />
              <v-spacer />
            </div>
            <v-text-field
              v-model="productTypeFormData.name"
              class="mx-4"
              maxlength="32"
              label="品目"
            />
            <v-divider></v-divider>

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
        <p>Category Item list will be displayed</p>
      </v-tab-item>
    </v-tabs-items>
  </div>
</template>

<script lang="ts">
import {
  computed,
  defineComponent,
  reactive,
  ref,
} from '@nuxtjs/composition-api'

import TheCategoryList from '~/components/organisms/TheCategoryList.vue'
import { useCategoryStore } from '~/store/category'
import { useCommonStore } from '~/store/common'
import { useProductTypeStore } from '~/store/product-type'
import { CreateCategoryRequest, CreateProductTypeRequest } from '~/types/api'
import { Category } from '~/types/props/category'

export default defineComponent({
  components: {
    TheCategoryList,
  },

  setup() {
    const categoryStore = useCategoryStore()
    const productTypeStore = useProductTypeStore()

    const categories = computed(() => {
      return categoryStore.categories.map((item) => {
        return {
          text: item.name,
          value: item.id,
        }
      })
    })

    const commonStore = useCommonStore()
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
        commonStore.addSnackbar({
          message: `カテゴリーを追加しました。`,
          color: 'info',
        })
      } catch (error) {
        console.log(error)
      }
    }

    const productTypeRegister = async (): Promise<void> => {
      try {
        console.log(selectedCategoryId)
        await productTypeStore.createProductType(
          selectedCategoryId.value,
          productTypeFormData
        )
        productTypeDialog.value = false
        commonStore.addSnackbar({
          message: `品目を追加しました。`,
          color: 'info',
        })
      } catch (error) {
        console.log(error)
      }
    }

    return {
      categories,
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
    }
  },
})
</script>
