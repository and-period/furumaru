<script lang="ts" setup>
import { mdiClose, mdiPlus } from '@mdi/js'
import useVuelidate from '@vuelidate/core'
import type { AlertType } from '~/lib/hooks'
import type { Shipping, UpdateDefaultShippingRequest, UpsertShippingRequest } from '~/types/api'
import { getErrorMessage } from '~/lib/validations'
import { type PrefecturesListSelectItems, getSelectablePrefecturesList } from '~/lib/prefectures'
import { UpsertShippingValidationRules } from '~/types/validations'

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
    type: Object as PropType<UpdateDefaultShippingRequest | UpsertShippingRequest>,
    default: (): UpdateDefaultShippingRequest | UpsertShippingRequest => ({
      box60Rates: [
        {
          name: '',
          price: 0,
          prefectureCodes: [],
        },
      ],
      box60Frozen: 0,
      box80Rates: [
        {
          name: '',
          price: 0,
          prefectureCodes: [],
        },
      ],
      box80Frozen: 0,
      box100Rates: [
        {
          name: '',
          price: 0,
          prefectureCodes: [],
        },
      ],
      box100Frozen: 0,
      hasFreeShipping: false,
      freeShippingRates: 0,
    }),
  },
  shipping: {
    type: Object as PropType<Shipping>,
    default: (): Shipping => ({
      id: '',
      isDefault: false,
      box60Rates: [],
      box60Frozen: 0,
      box80Rates: [],
      box80Frozen: 0,
      box100Rates: [],
      box100Frozen: 0,
      hasFreeShipping: false,
      freeShippingRates: 0,
      createdAt: 0,
      updatedAt: 0,
    }),
  },
})

const emit = defineEmits<{
  (e: 'update:form-data', v: UpdateDefaultShippingRequest | UpsertShippingRequest): void
  (e: 'submit'): void
}>()

const formDataValue = computed({
  get: (): UpdateDefaultShippingRequest | UpsertShippingRequest => props.formData,
  set: (formData: UpdateDefaultShippingRequest | UpsertShippingRequest): void => emit('update:form-data', formData),
})
const box60RateItemsSize = computed(() => {
  return [...Array(formDataValue.value.box60Rates.length).keys()]
})
const box80RateItemsSize = computed(() => {
  return [...Array(formDataValue.value.box80Rates.length).keys()]
})
const box100RateItemsSize = computed(() => {
  return [...Array(formDataValue.value.box100Rates.length).keys()]
})

const validate = useVuelidate(UpsertShippingValidationRules, formDataValue)

const addBox60RateItem = () => {
  formDataValue.value.box60Rates.push({
    name: '',
    price: 0,
    prefectureCodes: [],
  })
}

const addBox80RateItem = () => {
  formDataValue.value.box80Rates.push({
    name: '',
    price: 0,
    prefectureCodes: [],
  })
}

const addBox100RateItem = () => {
  formDataValue.value.box100Rates.push({
    name: '',
    price: 0,
    prefectureCodes: [],
  })
}

const getSelectableBox60RatePrefecturesList = (i: number): PrefecturesListSelectItems[] => {
  return getSelectablePrefecturesList(formDataValue.value.box60Rates, i)
}

const getSelectableBox80RatePrefecturesList = (i: number): PrefecturesListSelectItems[] => {
  return getSelectablePrefecturesList(formDataValue.value.box80Rates, i)
}

const getSelectableBox100RatePrefecturesList = (i: number): PrefecturesListSelectItems[] => {
  return getSelectablePrefecturesList(formDataValue.value.box100Rates, i)
}

const onClickSelectAll = (rate: '60' | '80' | '100', i: number): void => {
  switch (rate) {
    case '60':
      formDataValue.value.box60Rates[i].prefectureCodes
        = getSelectableBox60RatePrefecturesList(i)
          .filter(item => !item.disabled)
          .map(item => item.value)
      break
    case '80':
      formDataValue.value.box80Rates[i].prefectureCodes
        = getSelectableBox80RatePrefecturesList(i)
          .filter(item => !item.disabled)
          .map(item => item.value)
      break
    case '100':
      formDataValue.value.box100Rates[i].prefectureCodes
        = getSelectableBox100RatePrefecturesList(i)
          .filter(item => !item.disabled)
          .map(item => item.value)
      break
  }
}

const onClickRemoveItem = (rate: '60' | '80' | '100', index: number): void => {
  switch (rate) {
    case '60':
      formDataValue.value.box60Rates.splice(index, 1)
      break
    case '80':
      formDataValue.value.box80Rates.splice(index, 1)
      break
    case '100':
      formDataValue.value.box100Rates.splice(index, 1)
      break
  }
}

const onSubmit = async (): Promise<void> => {
  const valid = await validate.value.$validate()
  if (!valid) {
    return
  }

  emit('submit')
}
</script>

<template>
  <v-alert
    v-show="props.isAlert"
    :type="props.alertType"
    v-text="props.alertText"
  />

  <v-card-title>配送情報詳細</v-card-title>
  <v-card class="mb-4 py-2">
    <v-card-title>配送オプション：サイズ60</v-card-title>
    <v-card-text>
      <div class="d-flex flex-column flex-md-row justify-center">
        <v-text-field
          v-model.number="validate.box60Frozen.$model"
          :error-messages="getErrorMessage(validate.box60Frozen.$errors)"
          label="冷凍配送価格"
          type="number"
          prefix="通常配送料＋"
          suffix="円"
          min="0"
        />
      </div>
      <div
        v-for="i in box60RateItemsSize"
        :key="`60-${i}`"
        class="px-4 py-2 mb-2 border"
      >
        <div class="d-flex flex-row align-center">
          <p class="text-subtitle-2 text-grey">
            オプション{{ i + 1 }}
          </p>
          <v-spacer />
          <v-btn
            v-show="box60RateItemsSize.length > 1"
            :icon="mdiClose"
            color="error"
            variant="text"
            size="small"
            @click="onClickRemoveItem('60', i)"
          />
        </div>
        <molecules-shipping-rate-form
          v-model="formDataValue.box60Rates[i]"
          :selectable-prefecture-list="getSelectableBox60RatePrefecturesList(i)"
          @click:select-all="onClickSelectAll('60', i)"
        />
      </div>
    </v-card-text>
    <v-card-actions>
      <v-btn
        color="primary"
        variant="outlined"
        block
        @click="addBox60RateItem"
      >
        <v-icon :icon="mdiPlus" />
        追加
      </v-btn>
    </v-card-actions>
  </v-card>

  <v-card class="mb-4 py-2">
    <v-card-title>配送オプション：サイズ80</v-card-title>
    <v-card-text>
      <div class="d-flex flex-column flex-md-row justify-center">
        <v-text-field
          v-model.number="validate.box80Frozen.$model"
          :error-messages="getErrorMessage(validate.box80Frozen.$errors)"
          label="冷凍配送価格"
          type="number"
          prefix="通常配送料＋"
          suffix="円"
          min="0"
        />
      </div>
      <div
        v-for="i in box80RateItemsSize"
        :key="`80-${i}`"
        class="px-4 py-2 mb-2 border"
      >
        <div class="d-flex flex-row align-center">
          <p class="text-subtitle-2 text-grey">
            オプション{{ i + 1 }}
          </p>
          <v-spacer />
          <v-btn
            v-show="box80RateItemsSize.length > 1"
            :icon="mdiClose"
            color="error"
            variant="text"
            size="small"
            @click="onClickRemoveItem('80', i)"
          />
        </div>
        <molecules-shipping-rate-form
          v-model="formDataValue.box80Rates[i]"
          :selectable-prefecture-list="getSelectableBox80RatePrefecturesList(i)"
          @click:select-all="onClickSelectAll('80', i)"
        />
      </div>
    </v-card-text>
    <v-card-actions>
      <v-btn
        color="primary"
        variant="outlined"
        block
        @click="addBox80RateItem"
      >
        <v-icon :icon="mdiPlus" />
        追加
      </v-btn>
    </v-card-actions>
  </v-card>

  <v-card class="mb-4 py-2">
    <v-card-title>配送オプション：サイズ100</v-card-title>
    <v-card-text>
      <div class="d-flex flex-column flex-md-row justify-center">
        <v-text-field
          v-model.number="validate.box100Frozen.$model"
          :error-messages="getErrorMessage(validate.box100Frozen.$errors)"
          label="冷凍配送価格"
          type="number"
          prefix="通常配送料＋"
          suffix="円"
          min="0"
        />
      </div>
      <div
        v-for="i in box100RateItemsSize"
        :key="`100-${i}`"
        class="px-4 py-2 mb-2 border"
      >
        <div class="d-flex flex-row align-center">
          <p class="text-subtitle-2 text-grey">
            オプション{{ i + 1 }}
          </p>
          <v-spacer />
          <v-btn
            v-show="box100RateItemsSize.length > 1"
            :icon="mdiClose"
            color="error"
            variant="text"
            size="small"
            @click="onClickRemoveItem('100', i)"
          />
        </div>
        <molecules-shipping-rate-form
          v-model="formDataValue.box100Rates[i]"
          :selectable-prefecture-list="getSelectableBox100RatePrefecturesList(i)"
          @click:select-all="onClickSelectAll('100', i)"
        />
      </div>
    </v-card-text>
    <v-card-actions>
      <v-btn
        color="primary"
        variant="outlined"
        block
        @click="addBox100RateItem"
      >
        <v-icon :icon="mdiPlus" />
        追加
      </v-btn>
    </v-card-actions>
  </v-card>

  <v-btn
    :loading="loading"
    block
    variant="outlined"
    @click="onSubmit"
  >
    更新
  </v-btn>
</template>
