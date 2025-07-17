<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import type { PropType } from 'vue'
import { getErrorMessage } from '~/lib/validations'
import { Weekday } from '~/types/api'
import type { ProductType, Shop, UpdateShopRequest } from '~/types/api'
import { UpdateShopValidationRules } from '~/types/validations'

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
      businessDays: [],
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

const weekdays = [
  { title: '日曜日', value: Weekday.SUNDAY },
  { title: '月曜日', value: Weekday.MONDAY },
  { title: '火曜日', value: Weekday.TUESDAY },
  { title: '水曜日', value: Weekday.WEDNESDAY },
  { title: '木曜日', value: Weekday.THURSDAY },
  { title: '金曜日', value: Weekday.FRIDAY },
  { title: '土曜日', value: Weekday.SATURDAY },
]

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
  <v-card>
    <v-form @submit.prevent="onSubmit">
      <v-card-text>
        <v-text-field
          v-model="validate.name.$model"
          :error-messages="getErrorMessage(validate.name.$errors)"
          class="mr-4"
          label="マルシェ名"
        />
        <v-autocomplete
          v-model="formDataValue.productTypeIds"
          label="取り扱い品目"
          :items="productTypes"
          item-title="name"
          item-value="id"
          chips
          closable-chips
          multiple
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
        <v-select
          v-model="validate.businessDays.$model"
          label="営業日(発送可能日)"
          :error-messages="getErrorMessage(validate.businessDays.$errors)"
          :items="weekdays"
          item-title="title"
          item-value="value"
          chips
          closable-chips
          multiple
        />
      </v-card-text>

      <v-card-actions>
        <v-btn
          block
          :loading="loading"
          variant="outlined"
          color="primary"
          type="submit"
        >
          更新
        </v-btn>
      </v-card-actions>
    </v-form>
  </v-card>
</template>
