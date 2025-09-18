<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import type { PropType } from 'vue'
import { mdiStore, mdiCalendarMonth, mdiTagMultiple } from '@mdi/js'
import { getErrorMessage } from '~/lib/validations'
import type { TimeWeekday, ProductType, Shop, UpdateShopRequest } from '~/types/api/v1'
import { UpdateShopValidationRules } from '~/types/validations'
import { weekdays } from '~/constants'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  formData: {
    type: Object as PropType<UpdateShopRequest>,
    default: (): UpdateShopRequest => ({
      name: '',
      productTypeIds: [],
      businessDays: new Set<TimeWeekday>(),
    }),
  },
  shop: {
    type: Object as PropType<Shop>,
    default: (): Shop => ({
      id: '',
      name: '',
      coordinatorId: '',
      producerIds: [],
      productTypeIds: [],
      businessDays: [],
      createdAt: 0,
      updatedAt: 0,
    }),
  },
  productTypes: {
    type: Array<ProductType>,
    default: () => [],
  },
})

const emit = defineEmits<{
  (e: 'update:form-data', formData: UpdateShopRequest): void
  (e: 'update:search-product-type', name: string): void
  (e: 'click:search-address'): void
  (e: 'submit'): void
}>()

const formDataValue = computed({
  get: (): UpdateShopRequest => props.formData,
  set: (value: UpdateShopRequest): void => emit('update:form-data', value),
})

const validate = useVuelidate(UpdateShopValidationRules, formDataValue)

const onSubmit = async (): Promise<void> => {
  const valid = await validate.value.$validate()
  if (!valid) {
    return
  }

  emit('submit')
}

const onChangeSearchProductType = (name: string): void => {
  emit('update:search-product-type', name)
}
</script>

<template>
  <v-form @submit.prevent="onSubmit">
    <!-- 店舗基本情報セクション -->
    <v-card
      class="form-section-card mb-6"
      elevation="2"
    >
      <v-card-title class="d-flex align-center section-header">
        <v-icon
          :icon="mdiStore"
          size="24"
          class="mr-3 text-primary"
        />
        <span class="text-h6 font-weight-medium">店舗基本情報</span>
      </v-card-title>
      <v-card-text class="pa-6">
        <v-text-field
          v-model="validate.name.$model"
          :error-messages="getErrorMessage(validate.name.$errors)"
          label="マルシェ名 *"
          variant="outlined"
          density="comfortable"
          class="mb-4"
        />
      </v-card-text>
    </v-card>

    <!-- 取り扱い品目セクション -->
    <v-card
      class="form-section-card mb-6"
      elevation="2"
    >
      <v-card-title class="d-flex align-center section-header">
        <v-icon
          :icon="mdiTagMultiple"
          size="24"
          class="mr-3 text-primary"
        />
        <span class="text-h6 font-weight-medium">取り扱い品目</span>
      </v-card-title>
      <v-card-text class="pa-6">
        <v-autocomplete
          v-model="validate.productTypeIds.$model"
          label="取り扱い品目を選択"
          :error-messages="getErrorMessage(validate.productTypeIds.$errors)"
          :items="productTypes"
          item-title="name"
          item-value="id"
          chips
          closable-chips
          multiple
          variant="outlined"
          density="comfortable"
          clearable
          @update:search="onChangeSearchProductType"
        >
          <template #chip="{ props: val, item }">
            <v-chip
              v-bind="val"
              :prepend-avatar="item.raw?.iconUrl"
              :text="item.raw?.name"
              rounded
              class="px-4"
              variant="outlined"
            />
          </template>
          <template #item="{ props: val, item }">
            <v-list-item
              v-bind="val"
              :prepend-avatar="item?.raw?.iconUrl"
              :title="item?.raw?.name"
            />
          </template>
        </v-autocomplete>
      </v-card-text>
    </v-card>

    <!-- 営業日設定セクション -->
    <v-card
      class="form-section-card mb-6"
      elevation="2"
    >
      <v-card-title class="d-flex align-center section-header">
        <v-icon
          :icon="mdiCalendarMonth"
          size="24"
          class="mr-3 text-primary"
        />
        <span class="text-h6 font-weight-medium">営業日設定</span>
      </v-card-title>
      <v-card-text class="pa-6">
        <v-select
          v-model="validate.businessDays.$model"
          label="営業日（発送可能日）"
          :error-messages="getErrorMessage(validate.businessDays.$errors)"
          :items="weekdays"
          chips
          closable-chips
          multiple
          variant="outlined"
          density="comfortable"
        />
      </v-card-text>
    </v-card>

    <!-- 更新ボタン -->
    <div class="d-flex justify-end mt-6">
      <v-btn
        :loading="loading"
        variant="elevated"
        color="primary"
        size="large"
        type="submit"
        class="px-8 py-2"
        elevation="2"
      >
        <v-icon
          start
          size="20"
        >
          mdi-content-save
        </v-icon>
        店舗情報を更新
      </v-btn>
    </div>
  </v-form>
</template>

<style scoped>
.form-section-card {
  border-radius: 16px;
  overflow: hidden;
  transition: all 0.3s ease;
  border: 1px solid rgb(0 0 0 / 6%);
}

.form-section-card:hover {
  box-shadow: 0 4px 12px rgb(0 0 0 / 8%);
  transform: translateY(-2px);
}

.section-header {
  background: linear-gradient(135deg, rgb(33 150 243 / 8%) 0%, rgb(33 150 243 / 0%) 100%);
  border-bottom: 1px solid rgb(33 150 243 / 12%);
  padding: 18px 24px;
  font-weight: 500;
}

@media (width <= 600px) {
  .form-section-card {
    border-radius: 8px;
  }

  .section-header {
    padding: 12px 16px;
  }
}
</style>
