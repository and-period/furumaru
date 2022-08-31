<template>
  <div>
    <v-card-title>商品登録</v-card-title>
    <v-breadcrumbs :items="breadcrumbsItem" large class="pa-0 mb-6" />

    <v-alert v-model="isShow" :type="alertType" v-text="alertText" />

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
          <v-text-field
            v-model="v$.name.$model"
            label="商品名"
            outlined
            :error="v$.name.$error"
            :error-messages="getErrorMessage('name')"
          />
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
            <v-text-field
              v-model="v$.inventory.$model"
              :error-messages="getErrorMessage('inventory')"
              type="number"
              label="在庫数"
            />
            <v-spacer />
          </div>

          <div class="d-flex">
            <v-select
              v-model="formData.itemUnit"
              label="単位"
              :items="['個']"
            />
            <v-spacer />
          </div>

          <div class="d-flex align-center">
            <v-text-field
              v-model="v$.itemDescription.$model"
              label="単位説明"
              :error-messages="getErrorMessage('itemDescription')"
            />
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
          <v-radio-group>
            <div
              v-for="(img, i) in formData.media"
              :key="i"
              class="d-flex flex-row align-center"
            >
              <v-radio :value="i" />
              <img :src="img.url" width="200" class="mx-4" />
              <p class="mb-0">{{ img.url }}</p>
            </div>
          </v-radio-group>
          <p>※ check された商品画像がサムネイルになります</p>
        </v-card-text>
      </v-card>

      <v-card elevation="0" class="mb-4">
        <v-card-title>価格</v-card-title>
        <v-card-text>
          <v-text-field
            v-model="v$.price.$model"
            label="販売価格"
            :error-messages="getErrorMessage('price')"
          />
        </v-card-text>
      </v-card>

      <v-card elevation="0" class="mb-4">
        <v-card-title>配送情報</v-card-title>
        <v-card-text>
          <div class="d-flex">
            <v-text-field
              v-model="v$.weight.$model"
              label="重さ"
              :error-messages="getErrorMessage('weight')"
            >
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
              :items="categoriesItem"
              item-text="name"
              item-value="id"
            />
            <v-select
              v-model="formData.productTypeId"
              label="品目"
              :items="productTypesItem"
              item-text="name"
              item-value="id"
            />
          </div>
          <div class="d-flex">
            <v-select
              v-model="formData.originPrefecture"
              class="mr-4"
              label="原産地（都道府県）"
            />
            <v-select v-model="formData.originCity" label="原産地（市町村）" />
          </div>
          <v-select
            v-model="formData.producerId"
            label="店舗名"
            :items="producersItem"
            item-text="storeName"
            item-value="id"
          />
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
import { useFetch, useRouter } from '@nuxtjs/composition-api'
import {
  computed,
  defineComponent,
  reactive,
  Ref,
  ref,
} from '@vue/composition-api'
import useVuelidate from '@vuelidate/core'
import { required, minValue } from '@vuelidate/validators'

import { useAlert } from '~/lib/hooks'
import { useCategoryStore } from '~/store/category'
import { useProducerStore } from '~/store/producer'
import { useProductStore } from '~/store/product'
import { useProductTypeStore } from '~/store/product-type'
import { CreateProductRequest, UploadImageResponse } from '~/types/api'

export default defineComponent({
  setup() {
    const productTypeStore = useProductTypeStore()
    const categoryStore = useCategoryStore()
    const producerStore = useProducerStore()

    useFetch(async () => {
      await Promise.all([
        productTypeStore.fetchProductTypes(),
        categoryStore.fetchCategories(),
        producerStore.fetchProducers(),
      ])
    })

    const router = useRouter()

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

    const rules = computed(() => ({
      name: { required },
      inventory: { required, minValue: minValue(0) },
      price: { required, minValue: minValue(0) },
      weight: { required, minValue: minValue(0) },
      itemUnit: { required },
      itemDescription: { required },
    }))

    const v$ = useVuelidate(rules, formData)

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
            isThumbnail: index === 0,
          })
        } catch (error) {
          console.log(error)
        }
      }
    }

    const { alertType, isShow, alertText, show } = useAlert('error')

    const handleFormSubmit = async () => {
      const result = await v$.value.v$alidate()
      if (!result) {
        return
      }
      try {
        await createProduct(formData)
        router.push('/products')
      } catch (error) {
        show(error.message)
        window.scrollTo({
          top: 0,
          behavior: 'smooth',
        })
      }
    }

    const getErrorMessage = (key: string): string | Ref<string> => {
      const error = v$.value.$errors.find((e) => {
        return e.$property === key
      })
      return error ? error.$message : ''
    }

    return {
      alertType,
      isShow,
      alertText,
      productTypesItem: productTypeStore.productTypes,
      categoriesItem: categoryStore.categories,
      producersItem: producerStore.producers,
      breadcrumbsItem,
      statusItems,
      deliveryTypeItems,
      formData,
      v$,
      uploadFiles,
      productRef,
      handleUpdateFormDataDescription,
      handleImageUpload,
      handleFormSubmit,
      getErrorMessage,
    }
  },
})
</script>
