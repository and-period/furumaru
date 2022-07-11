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
            <v-text-field v-model="categoryFormData.name" class="mx-4" maxlength="32" label="カテゴリー" />
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
        <p>Category list will be displayed</p>
      </v-tab-item>

      <v-tab-item value="tab-categoryItems">
        <v-dialog v-model="itemDialog" width="500">
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
              <v-select class="mx-4" label="カテゴリー" />
              <v-spacer />
            </div>
            <v-text-field v-model="itemFormData.name" class="mx-4" maxlength="32" label="品目" />
            <v-divider></v-divider>

            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="accentDarken" text @click="itemCancel">
                キャンセル
              </v-btn>
              <v-btn color="primary" outlined @click="itemRegister">
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
import { defineComponent, reactive, ref } from '@nuxtjs/composition-api'

import { useCategoryStore } from '~/store/category'
import { CreateCategoryRequest, CreateProductTypeRequest } from '~/types/api'
import { Category } from '~/types/props/category'

export default defineComponent({
  setup() {
    const { createCategory } = useCategoryStore()

    const selector = ref<string>('categories')
    const categoryDialog = ref<boolean>(false)
    const itemDialog = ref<boolean>(false)
    const items: Category[] = [
      { name: 'カテゴリー', value: 'categories' },
      { name: '品目', value: 'categoryItems' },
    ]

    const categoryFormData = reactive<CreateCategoryRequest>({
      name: '',
    })

    const itemFormData = reactive<CreateProductTypeRequest>({
      name: '',
    })

    const categoryCancel = (): void => {
      categoryDialog.value = false
    }

    const itemCancel = (): void => {
      itemDialog.value = false
    }

    const categoryRegister = async (): Promise<void> => {
      try {
        await createCategory(categoryFormData)
      } catch (error) {
        console.log(error)
      }
    }

    const itemRegister = async (): Promise<void> => {
      // TODO: categoryが実装できた後に実装する
    }

    return {
      items,
      selector,
      categoryDialog,
      categoryFormData,
      itemFormData,
      itemDialog,
      categoryCancel,
      itemCancel,
      categoryRegister,
      itemRegister,
    }
  },
})
</script>
