<script lang="ts" setup>
import { mdiClose, mdiPlus } from '@mdi/js'
import useVuelidate from '@vuelidate/core'
import { AlertType } from '~/lib/hooks'
import { CreateShippingRequest } from '~/types/api'
import { required, getErrorMessage, minValue } from '~/lib/validations'
import { PrefecturesListSelectItems, getSelectablePrefecturesList } from '~/lib/prefectures'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  isAlert: {
    type: Boolean,
    default: false
  },
  alertType: {
    type: String as PropType<AlertType>,
    default: undefined
  },
  alertText: {
    type: String,
    default: ''
  },
  formData: {
    type: Object as PropType<CreateShippingRequest>,
    default: (): CreateShippingRequest => ({
      cordinatorId: '',
      name: '',
      isDefault: false,
      box60Rates: [
        {
          name: '',
          price: 0,
          prefectures: []
        }
      ],
      box60Refrigerated: 0,
      box60Frozen: 0,
      box80Rates: [
        {
          name: '',
          price: 0,
          prefectures: []
        }
      ],
      box80Refrigerated: 0,
      box80Frozen: 0,
      box100Rates: [
        {
          name: '',
          price: 0,
          prefectures: []
        }
      ],
      box100Refrigerated: 0,
      box100Frozen: 0,
      hasFreeShipping: false,
      freeShippingRates: 0
    })
  }
})

const emit = defineEmits<{
  (e: 'update:form-data', formData: CreateShippingRequest): void
  (e: 'submit'): void
}>()

const rules = computed(() => ({
  name: { required },
  hasFreeShipping: { required },
  box60Refrigerated: { required, minValue: minValue(0) },
  box60Frozen: { required, minValue: minValue(0) },
  box80Refrigerated: { required, minValue: minValue(0) },
  box80Frozen: { required, minValue: minValue(0) },
  box100Refrigerated: { required, minValue: minValue(0) },
  box100Frozen: { required, minValue: minValue(0) }
}))
const formDataValue = computed({
  get: (): CreateShippingRequest => props.formData,
  set: (formData: CreateShippingRequest): void => emit('update:form-data', formData)
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

const validate = useVuelidate(rules, formDataValue)

const addBox60RateItem = (): void => {
  formDataValue.value.box60Rates.push({
    name: '',
    price: 0,
    prefectures: []
  })
}

const addBox80RateItem = (): void => {
  formDataValue.value.box80Rates.push({
    name: '',
    price: 0,
    prefectures: []
  })
}

const addBox100RateItem = (): void => {
  formDataValue.value.box100Rates.push({
    name: '',
    price: 0,
    prefectures: []
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
      formDataValue.value.box60Rates[i].prefectures =
        getSelectableBox60RatePrefecturesList(i)
          .filter(item => !item.disabled)
          .map(item => item.value)
      break
    case '80':
      formDataValue.value.box80Rates[i].prefectures =
        getSelectableBox80RatePrefecturesList(i)
          .filter(item => !item.disabled)
          .map(item => item.value)
      break
    case '100':
      formDataValue.value.box100Rates[i].prefectures =
        getSelectableBox100RatePrefecturesList(i)
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
  <v-alert v-show="props.isAlert" :type="props.alertType" v-text="props.alertText" />

  <v-card>
    <v-card-title>配送情報登録</v-card-title>

    <v-form @submit.prevent="onSubmit">
      <v-card-text>
        <v-text-field
          v-model="validate.name.$model"
          label="名前"
          :error-messages="getErrorMessage(validate.name.$errors)"
        />
        <v-switch
          v-model="validate.isDefault.$model"
          :label="デフォルト設定"
        />
        <v-switch
          v-model="validate.hasFreeShipping.$model"
          :label="無料配送オプション"
        />

        <div class="my-6">
          <p class="text-h6">
            サイズ60配送オプション
          </p>
          <div class="d-flex">
            <v-text-field
              v-model.number="validate.box60Refrigerated.$model"
              :error-messages="getErrorMessage(validate.box60Refrigerated.$errors)"
              label="冷蔵配送価格"
              type="number"
              class="mr-4"
            />
            <v-text-field
              v-model.number="validate.box60Frozen.$model"
              :error-messages="getErrorMessage(validate.box60Frozen.$errors)"
              label="冷凍配送価格"
              type="number"
            />
          </div>
          <div v-for="i in box60RateItemsSize" :key="i">
            <div class="d-flex align-center">
              <p class="mb-0">
                オプション{{ i + 1 }}
              </p>
              <v-spacer />
              <v-btn
                icon
                color="error"
                variant="text"
                size="small"
                :disabled="box60RateItemsSize.length === 1"
                @click="onClickRemoveItem('60', i)"
              >
                <v-icon :icon="mdiClose" />
              </v-btn>
            </div>
            <molecules-shipping-rate-form
              v-model="formDataValue.box60Rates[i]"
              :selectable-prefecture-list="getSelectableBox60RatePrefecturesList(i)"
              @click:select-all="onClickSelectAll('60', i)"
            />
          </div>
          <v-btn color="primary" variant="outlined" block @click="addBox60RateItem">
            <v-icon :icon="mdiPlus" />
            追加
          </v-btn>
        </div>

        <div class="my-6">
          <p class="text-h6">
            サイズ80配送オプション
          </p>
          <div class="d-flex">
            <v-text-field
              v-model.number="validate.box80Refrigerated.$model"
              :error-messages="getErrorMessage(validate.box80Refrigerated.$errors)"
              label="冷蔵配送価格"
              type="number"
              class="mr-4"
            />
            <v-text-field
              v-model.number="validate.box80Frozen.$model"
              :error-messages="getErrorMessage(validate.box80Frozen.$errors)"
              label="冷凍配送価格"
              type="number"
            />
          </div>
          <div v-for="i in box80RateItemsSize" :key="i">
            <div class="d-flex align-center">
              <p class="mb-0">
                オプション{{ i + 1 }}
              </p>
              <v-spacer />
              <v-btn
                icon
                color="error"
                variant="text"
                size="small"
                :disabled="box80RateItemsSize.length === 1"
                @click="onClickRemoveItem('80', i)"
              >
                <v-icon :icon="mdiClose" />
              </v-btn>
            </div>

            <molecules-shipping-rate-form
              v-model="formDataValue.box80Rates[i]"
              :selectable-prefecture-list="getSelectableBox80RatePrefecturesList(i)"
              @click:select-all="onClickSelectAll('80', i)"
            />
          </div>
          <v-btn color="primary" variant="outlined" block @click="addBox80RateItem">
            <v-icon :icon="mdiPlus" />
            追加
          </v-btn>
        </div>

        <div class="my-6">
          <p class="text-h6">
            サイズ100配送オプション
          </p>
          <div class="d-flex">
            <v-text-field
              v-model.number="validate.box100Refrigerated.$model"
              :error-messages="getErrorMessage(validate.box100Refrigerated.$errors)"
              label="冷蔵配送価格"
              class="mr-4"
              type="number"
            />
            <v-text-field
              v-model.number="validate.box100Frozen.$model"
              :error-messages="getErrorMessage(validate.box100Frozen.$errors)"
              label="冷凍配送価格"
              type="number"
            />
          </div>
          <div v-for="i in box100RateItemsSize" :key="i">
            <div class="d-flex align-center">
              <p class="mb-0">
                オプション{{ i + 1 }}
              </p>
              <v-spacer />
              <v-btn
                icon
                color="error"
                variant="text"
                size="small"
                :disabled="box100RateItemsSize.length === 1"
                @click="onClickRemoveItem('100', i)"
              >
                <v-icon :icon="mdiClose" />
              </v-btn>
            </div>

            <molecules-shipping-rate-form
              v-model="formDataValue.box100Rates[i]"
              :selectable-prefecture-list="getSelectableBox100RatePrefecturesList(i)"
              @click:select-all="onClickSelectAll('100', i)"
            />
          </div>
          <v-btn color="primary" variant="outlined" block @click="addBox100RateItem">
            <v-icon :icon="mdiPlus" />
            追加
          </v-btn>
        </div>
      </v-card-text>

      <v-card-actions>
        <v-btn :loading="loading" type="submit" variant="outlined" color="primary" block>
          登録
        </v-btn>
      </v-card-actions>
    </v-form>
  </v-card>
</template>
