<script setup lang="ts">
import { useAuthStore } from '~/store/auth'
import { useProductStore } from '~/store/product'
import { useProductReviewStore } from '~/store/productReview'
import type { CreateProductReviewRequest } from '~/types/api'
import { ApiBaseError } from '~/types/exception'
import type { I18n } from '~/types/locales'

const i18n = useI18n()

const authStore = useAuthStore()
const { isAuthenticated } = storeToRefs(authStore)

const productStore = useProductStore()
const { fetchProduct } = productStore

const productReviewStore = useProductReviewStore()
const { postReview } = productReviewStore

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

const submitting = ref<boolean>(false)
const submitErrorMessage = ref<string>('')

const handleSubmit = async () => {
  submitting.value = true
  submitErrorMessage.value = ''
  try {
    await postReview(productId.value, formData.value)
  }
  catch (error) {
    if (error instanceof ApiBaseError) {
      submitErrorMessage.value = error.message
    }
    else {
      submitErrorMessage.value = ''
    }
  }
  finally {
    submitting.value = false
  }
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

      <template v-if="!isAuthenticated">
        <div class="flex flex-col md:gap-8 gap-4">
          <the-alert>
            {{ lt('requiredAuthMessage') }}
          </the-alert>

          <div class="text-center">
            <nuxt-link
              :to="`/signin?review_target_id=${productId}`"
              class=" bg-main text-white py-2 md:w-[400px] inline-block w-full"
            >
              {{ lt('loginButtonText') }}
            </nuxt-link>
          </div>
        </div>
      </template>

      <template v-if="isAuthenticated">
        <!-- エラー表示 -->
        <template v-if="status === 'error'">
          <the-alert>
            {{ error }}
          </the-alert>
        </template>

        <template v-if="submitErrorMessage">
          <the-alert class="mb-4">
            <div class=" font-semibold">
              {{ lt('reviewSubmitErrorMessage') }}
            </div>
            {{ submitErrorMessage }}
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
              :submitting="submitting"
              @submit="handleSubmit"
            />
          </div>
        </template>
      </template>
    </div>
  </div>
</template>
