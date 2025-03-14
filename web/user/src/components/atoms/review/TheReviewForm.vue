<script setup lang="ts">
import type { CreateProductReviewRequest } from '~/types/api'
import type { I18n } from '~/types/locales'

interface Props {
  submitting: boolean
}

const props = defineProps<Props>()

interface Emits {
  (e: 'submit'): void
}

const emits = defineEmits<Emits>()

const model = defineModel<CreateProductReviewRequest>({ required: true })

const maxTitleLength = 64
const maxCommentLength = 2000

const i18n = useI18n()

const lt = (str: keyof I18n['reviews']) => {
  return i18n.t(`reviews.${str}`)
}

const ratingError = ref<string>('')

const handleSubmit = () => {
  if (props.submitting) {
    return
  }

  ratingError.value = ''

  // rateが1,2,3,4,5以外の場合はエラー
  if (![1, 2, 3, 4, 5].includes(model.value.rate)) {
    ratingError.value = '評価を選択してください'
    return
  }

  emits('submit')
}
</script>

<template>
  <form
    class="flex flex-col gap-4"
    @submit.prevent="handleSubmit"
  >
    <div class="flex flex-col gap-2">
      <label
        for="rating"
        class=" font-semibold text-[16px] tracking-[2px]"
      >
        {{ lt('ratingLabel') }}
      </label>
      <div class="flex flex-col gap-1">
        <the-rating-star-input v-model="model.rate" />
        <div class="text-orange text-[14px] tracking-[2px]">
          {{ ratingError }}
        </div>
      </div>
    </div>

    <div class="flex flex-col gap-2">
      <label
        for="title"
        class="font-semibold text-[16px] tracking-[2px]"
      >
        {{ lt('reviewTitleLabel') }}
      </label>
      <div class="w-full">
        <input
          id="title"
          v-model="model.title"
          :placeholder="lt('reviewTitlePlaceholder')"
          class="w-full border border-main py-1 rounded-md px-2"
          type="text"
          name="title"
          required
          :maxlength="maxTitleLength"
        >
        <div class="text-right text-[12px] text-gray-400">
          {{ model.title.length }} / {{ maxTitleLength }}
        </div>
      </div>
    </div>

    <div class="flex flex-col gap-2">
      <label
        for="comment"
        class=" font-semibold text-[16px] tracking-[2px]"
      >
        {{ lt('reviewCommentLabel') }}
      </label>
      <div class="w-full">
        <textarea
          id="comment"
          v-model="model.comment"
          :placeholder="lt('reviewCommentPlaceholder')"
          class="w-full border border-main py-1 rounded-md px-2"
          name="comment"
          rows="5"
          :maxlength="maxCommentLength"
          required
        />
        <div class="text-right text-[12px] text-gray-400">
          {{ model.comment.length }} / {{ maxCommentLength }}
        </div>
      </div>
    </div>
    <div class="text-center">
      <button
        :disabled="submitting"
        type="submit"
        class="bg-main text-white py-2 md:w-[400px] w-full disabled:cursor-wait"
      >
        <template v-if="submitting">
          <div class="flex items-center justify-center">
            <the-loading-icon />
          </div>
        </template>
        <template v-else>
          {{ lt('reviewSubmitButtonText') }}
        </template>
      </button>
    </div>
  </form>
</template>
