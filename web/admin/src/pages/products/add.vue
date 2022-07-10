<template>
  <div>
    <v-card-title>商品登録</v-card-title>
    <v-breadcrumbs :items="breadcrumbsItem" large class="pa-0 mb-6" />
    <div class="mb-4">
      <v-card elevation="0" class="mb-4">
        <v-card-title>商品ステータス</v-card-title>
        <v-card-text>
          <v-select label="ステータス" />
        </v-card-text>
      </v-card>

      <v-card elevation="0" class="mb-4">
        <v-card-title>基本情報</v-card-title>
        <v-card-text>
          <v-text-field label="商品名" outlined />
          <client-only>
            <tiptap-editor
              label="商品詳細"
              :value="productRef"
              @update:value="handleUpdateProduct"
            />
          </client-only>
        </v-card-text>
      </v-card>

      <v-card elevation="0" class="mb-4">
        <v-card-title>在庫</v-card-title>
        <v-card-text>
          <div class="d-flex">
            <v-text-field label="在庫数" />
            <v-spacer />
          </div>

          <div class="d-flex">
            <v-select label="単位" />
            <v-spacer />
          </div>

          <div class="d-flex align-center">
            <v-text-field label="単位説明" />
            <p class="ml-12 mb-0">ex) 1kg → 5個入り</p>
            <v-spacer />
          </div>
        </v-card-text>
      </v-card>

      <v-card elevation="0" class="mb-4">
        <v-card-title>商品画像登録</v-card-title>
        <v-card-text>
          <div class="mb-2">
            <the-file-upload-filed text="商品画像" />
          </div>
          <p>※ check された商品画像がサムネイルになります</p>
        </v-card-text>
      </v-card>

      <v-card elevation="0" class="mb-4">
        <v-card-title>価格</v-card-title>
        <v-card-text>
          <v-text-field label="販売価格" />
        </v-card-text>
      </v-card>

      <v-card elevation="0" class="mb-4">
        <v-card-title>配送情報</v-card-title>
        <v-card-text>
          <div class="d-flex">
            <v-text-field label="重さ">
              <template #append>kg</template>
            </v-text-field>
            <v-spacer />
          </div>
          <div class="d-flex">
            <v-select label="配送種別" />
            <v-spacer />
          </div>

          <v-list>
            <v-list-item>
              <v-list-item-action>箱のサイズ</v-list-item-action>
              <v-list-item-content> 占有率 </v-list-item-content>
            </v-list-item>
            <v-list-item v-for="(size, i) in [60, 80, 100]" :key="i">
              <v-list-item-action>
                <p class="mb-0 mx-6 text-h6">{{ size }}</p>
              </v-list-item-action>
              <v-list-item-content>
                <v-text-field label="占有率">
                  <template #append>%</template>
                </v-text-field>
              </v-list-item-content>
            </v-list-item>
          </v-list>
        </v-card-text>
      </v-card>

      <v-card elevation="0" class="mb-4">
        <v-card-title>詳細情報</v-card-title>
        <v-card-text>
          <div class="d-flex">
            <v-select class="mr-4" label="カテゴリ" />
            <v-select label="品目" />
          </div>
          <div class="d-flex">
            <v-select class="mr-4" label="原産地（都道府県）" />
            <v-select label="原産地（市町村）" />
          </div>
          <v-select label="店舗名" />
        </v-card-text>
      </v-card>
    </div>
    <v-btn block outlined @click="addFormItem">
      <v-icon left>mdi-plus</v-icon>
      登録
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

    const productRef = ref<string>('')
    const handleUpdateProduct = (htmlString: string) => {
      productRef.value = htmlString
    }

    return {
      breadcrumbsItem,
      addFormItem,
      formData,
      productRef,
      handleUpdateProduct,
    }
  },
})
</script>
