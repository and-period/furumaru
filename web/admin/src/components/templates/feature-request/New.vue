<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import type { AlertType } from '~/lib/hooks'
import { getErrorMessage, maxLength, required } from '~/lib/validations'
import {
  FeatureRequestCategory,
  FeatureRequestPriority,
} from '~/types/feature-request'
import type { CreateFeatureRequestInput } from '~/types/feature-request'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  isAlert: {
    type: Boolean,
    default: false,
  },
  alertType: {
    type: String as PropType<AlertType>,
    default: undefined,
  },
  alertText: {
    type: String,
    default: '',
  },
  formData: {
    type: Object as PropType<CreateFeatureRequestInput>,
    default: () => ({
      title: '',
      description: '',
      category: FeatureRequestCategory.Feature,
      priority: FeatureRequestPriority.Medium,
      submittedBy: '',
      submitterName: '',
    }),
  },
})

const emit = defineEmits<{
  (e: 'update:form-data', formData: CreateFeatureRequestInput): void
  (e: 'submit'): void
}>()

const categories = [
  { title: 'UI/UX改善', value: FeatureRequestCategory.UI },
  { title: '新機能', value: FeatureRequestCategory.Feature },
  { title: 'パフォーマンス改善', value: FeatureRequestCategory.Performance },
  { title: 'その他', value: FeatureRequestCategory.Other },
]

const priorities = [
  { title: '低', value: FeatureRequestPriority.Low },
  { title: '中', value: FeatureRequestPriority.Medium },
  { title: '高', value: FeatureRequestPriority.High },
]

const formDataValue = computed({
  get: (): CreateFeatureRequestInput => props.formData,
  set: (v: CreateFeatureRequestInput): void => emit('update:form-data', v),
})

const rules = computed(() => ({
  title: { required, maxLength: maxLength(128) },
  description: { required, maxLength: maxLength(2000) },
  category: { required },
  priority: { required },
}))

const validate = useVuelidate(rules, formDataValue)

const onSubmit = async (): Promise<void> => {
  const valid = await validate.value.$validate()
  if (!valid)
    return
  emit('submit')
}
</script>

<template>
  <atoms-app-alert
    :show="props.isAlert"
    :type="props.alertType"
    :text="props.alertText"
  />

  <v-card elevation="0">
    <v-card-title>要望リクエストを提出</v-card-title>
    <v-card-subtitle class="mt-1">
      プラットフォームへの機能改善・新機能の要望を送信できます。
    </v-card-subtitle>

    <v-form @submit.prevent="onSubmit">
      <v-card-text class="d-flex flex-column ga-2">
        <v-text-field
          v-model="validate.title.$model"
          :error-messages="getErrorMessage(validate.title.$errors)"
          label="タイトル"
          placeholder="例: 商品CSVの一括インポート機能を追加してほしい"
          counter="128"
          maxlength="128"
          required
        />

        <v-row>
          <v-col
            cols="12"
            sm="6"
          >
            <v-select
              v-model="validate.category.$model"
              :error-messages="getErrorMessage(validate.category.$errors)"
              :items="categories"
              item-title="title"
              item-value="value"
              label="カテゴリ"
              required
            />
          </v-col>
          <v-col
            cols="12"
            sm="6"
          >
            <v-select
              v-model="validate.priority.$model"
              :error-messages="getErrorMessage(validate.priority.$errors)"
              :items="priorities"
              item-title="title"
              item-value="value"
              label="優先度（ご自身の感覚で）"
              required
            />
          </v-col>
        </v-row>

        <v-textarea
          v-model="validate.description.$model"
          :error-messages="getErrorMessage(validate.description.$errors)"
          label="要望内容"
          placeholder="どのような機能が欲しいか、なぜ必要なのかを詳しく教えてください。"
          rows="6"
          counter="2000"
          maxlength="2000"
          required
        />
      </v-card-text>

      <v-card-actions class="px-4 pb-4">
        <v-btn
          :loading="loading"
          type="submit"
          color="primary"
          variant="elevated"
          block
        >
          提出する
        </v-btn>
      </v-card-actions>
    </v-form>
  </v-card>
</template>
