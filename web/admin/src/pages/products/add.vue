<template>
  <div>
    <p class="text-h6">商品登録</p>
    <v-breadcrumbs :items="breadcrumbsItem" large class="pa-0 mb-6" />
    <div class="mb-4">
      <v-card
        v-for="(item, idx) in formData"
        :key="idx"
        elevation="0"
        class="mb-4"
      >
        <v-card-text>
          <div class="mb-2">
            <the-file-upload-filed text="商品画像" />
          </div>
          <v-text-field label="商品名" />
          <v-text-field label="商品概要" />
          <v-text-field label="重さ" />
          <v-text-field label="ステータス" />
          <v-text-field label="在庫" />
          <div class="d-flex">
            <v-text-field class="mr-4" label="ジャンル" />
            <v-text-field label="品目" />
          </div>
          <div class="d-flex">
            <v-text-field class="mr-4" label="原産地（都道府県）" />
            <v-text-field label="原産地（市町村）" />
          </div>
          <v-text-field label="店舗" />
        </v-card-text>
      </v-card>
    </div>
    <v-btn block outlined @click="addFormItem">
      <v-icon left>mdi-plus</v-icon>
      追加
    </v-btn>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from '@vue/composition-api'

export default defineComponent({
  setup() {
    const breadcrumbsItem = [
      {
        text: '商品管理',
        href: '/products',
        disabled: false,
      },
      {
        text: '商品登録',
        href: 'add',
        disabled: true,
      },
    ]

    const addFormItem = () => {
      console.log('not')
      formData.value.push(1)
    }

    const formData = ref<number[]>([1])

    return {
      breadcrumbsItem,
      addFormItem,
      formData,
    }
  },
})
</script>
