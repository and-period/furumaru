<template>
  <div>
    <v-card-title>商品登録</v-card-title>
    <v-breadcrumbs :items="breadcrumbsItem" large class="pa-0 mb-6" />

    <div class="mb-4">
      <v-card elevation="0" class="mb-4">
        <v-card-title>商品ステータス</v-card-title>
        <v-card-text>
          <v-select
            v-model="formData.public"
            label="ステータス"
            :items="statusItems"
          />
        </v-card-text>
      </v-card>

      <v-card elevation="0" class="mb-4">
        <v-card-title>基本情報</v-card-title>
        <v-card-text>
          <v-text-field v-model="formData.name" label="商品名" outlined />
          <client-only>
            <tiptap-editor
              label="商品詳細"
              :value="formData.description"
              @update:value="handleUpdateFormDataDescription"
            />
          </client-only>
        </v-card-text>
      </v-card>

      <v-card elevation="0" class="mb-4">
        <v-card-title>在庫</v-card-title>
        <v-card-text>
          <div class="d-flex">
            <v-text-field v-model="formData.inventory" label="在庫数" />
            <v-spacer />
          </div>

          <div class="d-flex">
            <v-select v-model="formData.itemUnit" label="単位" />
            <v-spacer />
          </div>

          <div class="d-flex align-center">
            <v-text-field v-model="formData.itemDescription" label="単位説明" />
            <p class="ml-12 mb-0">ex) 1kg → 5個入り</p>
            <v-spacer />
          </div>
        </v-card-text>
      </v-card>

      <v-card elevation="0" class="mb-4">
        <v-card-title>商品画像登録</v-card-title>
        <v-card-text>
          <div class="mb-2">
            <the-file-upload-filed
              text="商品画像"
              @update:files="handleImageUpload"
            />
          </div>
          <p>※ check された商品画像がサムネイルになります</p>
        </v-card-text>
      </v-card>

      <v-card elevation="0" class="mb-4">
        <v-card-title>価格</v-card-title>
        <v-card-text>
          <v-text-field v-model="formData.price" label="販売価格" />
        </v-card-text>
      </v-card>

      <v-card elevation="0" class="mb-4">
        <v-card-title>配送情報</v-card-title>
        <v-card-text>
          <div class="d-flex">
            <v-text-field v-model="formData.weight" label="重さ">
              <template #append>kg</template>
            </v-text-field>
            <v-spacer />
          </div>
          <div class="d-flex">
            <v-select
              v-model="formData.deliveryType"
              :items="deliveryTypeItems"
              label="配送種別"
            />
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
                <v-text-field
                  v-model="formData[`box${size}Rate`]"
                  type="number"
                  min="0"
                  max="100"
                  label="占有率"
                >
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
            <v-select
              v-model="formData.categoryId"
              class="mr-4"
              label="カテゴリ"
            />
            <v-select v-model="formData.productTypeId" label="品目" />
          </div>
          <div class="d-flex">
            <v-select
              v-model="formData.originPrefecture"
              class="mr-4"
              label="原産地（都道府県）"
            />
            <v-select v-model="formData.originCity" label="原産地（市町村）" />
          </div>
          <v-select v-model="formData.producerId" label="店舗名" />
        </v-card-text>
      </v-card>
    </div>
    <v-btn block outlined @click="handleFormSubmit">
      <v-icon left>mdi-plus</v-icon>
      登録
    </v-btn>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive, ref } from '@vue/composition-api'

import { useProductStore } from '~/store/product'
import { CreateProductRequest, UploadImageResponse } from '~/types/api'

export default defineComponent({
  setup() {
    const { uploadProductImage, createProduct } = useProductStore()
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

    const statusItems = [
      { text: '公開', value: true },
      { text: '非公開', value: false },
    ]
    const deliveryTypeItems = [
      { text: '通常便', value: 1 },
      { text: '冷蔵便', value: 2 },
      { text: '冷凍便', value: 3 },
    ]

    const formData = reactive<CreateProductRequest>({
      name: '',
      description: '',
      producerId: '',
      categoryId: '',
      productTypeId: '',
      public: true,
      inventory: 0,
      weight: 0,
      itemUnit: '',
      itemDescription: '',
      media: [],
      price: 0,
      deliveryType: 1,
      box60Rate: 0,
      box80Rate: 0,
      box100Rate: 0,
      originPrefecture: '',
      originCity: '',
    })

    const uploadFiles = ref<FileList | null>(null)

    const productRef = ref<string>('')
    const handleUpdateFormDataDescription = (htmlString: string) => {
      formData.description = htmlString
    }

    const handleImageUpload = async (files: FileList) => {
      for (const [index, file] of Array.from(files).entries()) {
        try {
          const uploadImage: UploadImageResponse = await uploadProductImage(
            file
          )
          formData.media.push({
            ...uploadImage,
            isThumbnail: index === 1,
          })
        } catch (error) {
          console.log(error)
        }
      }
    }

    const handleFormSubmit = async () => {
      console.log('submit')
      try {
        await createProduct(formData)
      } catch (error) {
        console.log(error)
      }
    }

    return {
      breadcrumbsItem,
      statusItems,
      deliveryTypeItems,
      formData,
      uploadFiles,
      productRef,
      handleUpdateFormDataDescription,
      handleImageUpload,
      handleFormSubmit,
    }
  },
})
</script>
