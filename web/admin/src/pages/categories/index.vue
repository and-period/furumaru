<template>
  <div>
    <p class="text-h6">カテゴリー・品目設定</p>
    <v-tabs v-model="selector" grow color="dark">
      <v-tabs-slider color="accent"></v-tabs-slider>
      <v-tab v-for="item in items" :key="item.value" :href="`#tab-${item.value}`">
        {{ item.name }}
      </v-tab>
    </v-tabs>

    <v-tabs-items v-model="selector">
      <v-tab-item value="tab-categories">
        <div class="d-flex pt-3 pr-3">
          <v-spacer />
          <v-btn outlined @click="handleClickAddButton">
            <v-icon left>mdi-plus</v-icon>
            追加
          </v-btn>
        </div>
        <p> Category list will be displayed </p>
      </v-tab-item>
      <v-tab-item value="tab-categoryItems">
        <div class="d-flex pt-3 pr-3">
          <v-spacer />
          <v-btn outlined @click="handleClickAddButton">
            <v-icon left>mdi-plus</v-icon>
            追加
          </v-btn>
        </div>
        <p> Category Item list will be displayed </p>
      </v-tab-item>
    </v-tabs-items>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, useRouter } from '@nuxtjs/composition-api'

import { Category } from '~/types/props/category'

export default defineComponent({
  setup() {
    const router = useRouter()
    const selector = ref<string>('categories')
    const items: Category[] = [
      { name: 'カテゴリー', value: 'categories'},
      { name: '品目', value: 'categoryItems'}
    ]

    const handleClickAddButton = () => {
      router.push('/categories/add')
    }

    return {
      items,
      selector,
      handleClickAddButton,
    }
  },
})
</script>
