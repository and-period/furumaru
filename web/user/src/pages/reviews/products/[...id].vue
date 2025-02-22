<script setup lang="ts">
import { useProductStore } from '~/store/product'
import type { CreateProductReviewRequest } from '~/types/api'
import type { I18n } from '~/types/locales'

const i18n = useI18n()

const productStore = useProductStore()
const { fetchProduct } = productStore

const lt = (str: keyof I18n['reviews']) => {
  return i18n.t(`reviews.${str}`)
}

const route = useRoute()

const productId = computed<string>(() => {
  const id = route.params.id
  if (id instanceof Array) {
    return id[0]
  }
  return route.params.id as string
})

const formData = ref<CreateProductReviewRequest>({
  rate: 0,
  title: '',
  comment: '',
})

/**
 * レビュー対象の商品情報取得
 */
const { data: product, status, error } = useAsyncData('target-product', () => {
  return fetchProduct(productId.value)
})

const handleSubmit = () => {
  // TODO: レビューを投稿する処理を実装する
  console.log('submit')
}

useSeoMeta({
  title: lt('postReviewTitle'),
})
</script>

<template>
  <div
    class="flex flex-col bg-white px-[15px] py-[48px] text-main md:px-[36px]"
  >
    <div class="container mx-auto p-4 xl:p-0">
      <p
        class="text-center text-[14px] font-bold tracking-[2px] md:text-[20px]"
      >
        {{ lt('postReviewTitle') }}
      </p>
      <hr class="my-[40px]">

      <!-- エラー表示 -->
      <template v-if="status === 'error'">
        <the-alert>
          {{ error }}
        </the-alert>
      </template>

      <template v-if="status === 'success'">
        <div class="flex flex-col gap-4">
          <template v-if="product">
            <the-review-target-product
              :product="product?.product"
            />
          </template>
          <the-review-form
            v-model="formData"
            @submit="handleSubmit"
          />
        </div>
      </template>
    </div>
  </div>
</template>
