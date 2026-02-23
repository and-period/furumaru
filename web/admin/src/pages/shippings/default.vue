<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import useVuelidate from '@vuelidate/core'
import { useAlert } from '~/lib/hooks'
import { useCommonStore, useShippingStore } from '~/store'
import type { UpdateDefaultShippingRequest } from '~/types/api/v1'
import { UpdateDefaultShippingValidationRules } from '~/types/validations'
import { getSelectablePrefecturesList } from '~/lib/prefectures'
import type { PrefecturesListSelectItems } from '~/lib/prefectures'
import { mdiClose, mdiPlus, mdiPackageVariant, mdiSnowflake, mdiCurrencyJpy, mdiContentSave } from '@mdi/js'
import { getErrorMessage } from '~/lib/validations'

const commonStore = useCommonStore()
const shippingStore = useShippingStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const { shipping } = storeToRefs(shippingStore)

const loading = ref<boolean>(false)
const formData = ref<UpdateDefaultShippingRequest>({
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
})

/**
 * サイズ60の配送料オプションのインデックス配列を計算する
 * @returns サイズ60の配送料オプションのインデックス配列
 */
const box60RateItemsSize = computed(() => {
  return [...Array(formData.value.box60Rates.length).keys()]
})

/**
 * サイズ80の配送料オプションのインデックス配列を計算する
 * @returns サイズ80の配送料オプションのインデックス配列
 */
const box80RateItemsSize = computed(() => {
  return [...Array(formData.value.box80Rates.length).keys()]
})

/**
 * サイズ100の配送料オプションのインデックス配列を計算する
 * @returns サイズ100の配送料オプションのインデックス配列
 */
const box100RateItemsSize = computed(() => {
  return [...Array(formData.value.box100Rates.length).keys()]
})

const validate = useVuelidate(UpdateDefaultShippingValidationRules, formData)

const fetchState = useAsyncData('shipping-default', async (): Promise<void> => {
  try {
    await shippingStore.fetchDefaultShipping()
    formData.value = { ...shipping.value }
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
})

const isLoading = (): boolean => {
  return fetchState?.pending?.value || loading.value
}

/**
 * サイズ60の配送料オプションを追加する
 */
const addBox60RateItem = () => {
  formData.value.box60Rates.push({
    name: '',
    price: 0,
    prefectureCodes: [],
  })
}

/**
 * サイズ80の配送料オプションを追加する
 */
const addBox80RateItem = () => {
  formData.value.box80Rates.push({
    name: '',
    price: 0,
    prefectureCodes: [],
  })
}

/**
 * サイズ100の配送料オプションを追加する
 */
const addBox100RateItem = () => {
  formData.value.box100Rates.push({
    name: '',
    price: 0,
    prefectureCodes: [],
  })
}

/**
 * サイズ60の配送料オプションで選択可能な都道府県リストを取得する
 * @param i オプションのインデックス
 * @returns 選択可能な都道府県リスト
 */
const getSelectableBox60RatePrefecturesList = (i: number): PrefecturesListSelectItems[] => {
  return getSelectablePrefecturesList(formData.value.box60Rates, i)
}

/**
 * サイズ80の配送料オプションで選択可能な都道府県リストを取得する
 * @param i オプションのインデックス
 * @returns 選択可能な都道府県リスト
 */
const getSelectableBox80RatePrefecturesList = (i: number): PrefecturesListSelectItems[] => {
  return getSelectablePrefecturesList(formData.value.box80Rates, i)
}

/**
 * サイズ100の配送料オプションで選択可能な都道府県リストを取得する
 * @param i オプションのインデックス
 * @returns 選択可能な都道府県リスト
 */
const getSelectableBox100RatePrefecturesList = (i: number): PrefecturesListSelectItems[] => {
  return getSelectablePrefecturesList(formData.value.box100Rates, i)
}

/**
 * 「すべて選択」ボタンがクリックされた時、指定されたサイズと指定されたインデックスの配送料オプションで
 * 選択可能なすべての都道府県を選択する
 * @param rate サイズ（'60', '80', '100'）
 * @param i オプションのインデックス
 */
const onClickSelectAll = (rate: '60' | '80' | '100', i: number): void => {
  switch (rate) {
    case '60':
      formData.value.box60Rates[i].prefectureCodes
        = getSelectableBox60RatePrefecturesList(i)
          .filter(item => !item.disabled)
          .map(item => item.value)
      break
    case '80':
      formData.value.box80Rates[i].prefectureCodes
        = getSelectableBox80RatePrefecturesList(i)
          .filter(item => !item.disabled)
          .map(item => item.value)
      break
    case '100':
      formData.value.box100Rates[i].prefectureCodes
        = getSelectableBox100RatePrefecturesList(i)
          .filter(item => !item.disabled)
          .map(item => item.value)
      break
  }
}

/**
 * 配送料オプションを削除する
 * @param rate サイズ（'60', '80', '100'）
 * @param index 削除するオプションのインデックス
 */
const onClickRemoveItem = (rate: '60' | '80' | '100', index: number): void => {
  switch (rate) {
    case '60':
      formData.value.box60Rates.splice(index, 1)
      break
    case '80':
      formData.value.box80Rates.splice(index, 1)
      break
    case '100':
      formData.value.box100Rates.splice(index, 1)
      break
  }
}

const handleSubmit = async (): Promise<void> => {
  try {
    loading.value = true
    await shippingStore.updateDefaultShipping(formData.value)
    commonStore.addSnackbar({
      color: 'info',
      message: 'デフォルト配送設定を更新しました。',
    })
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)

    window.scrollTo({
      top: 0,
      behavior: 'smooth',
    })
  }
  finally {
    loading.value = false
  }
}

try {
  await fetchState.execute()
}
catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <v-container class="pa-6">
    <v-alert
      v-show="isShow"
      :type="alertType"
      class="mb-6"
      v-text="alertText"
    />

    <div class="mb-6">
      <h1 class="text-h4 font-weight-bold mb-2">
        デフォルト配送設定
      </h1>
      <p class="text-body-1 text-grey-darken-1">
        システム全体のデフォルト配送料金を設定します。各コーディネーターの配送設定が無い場合にこの設定が適用されます。
      </p>
    </div>

    <v-card
      class="form-section-card mb-6"
      elevation="2"
      :loading="loading"
    >
      <v-card-title class="d-flex align-center section-header">
        <v-icon
          :icon="mdiPackageVariant"
          size="24"
          class="mr-3 text-primary"
        />
        <span class="text-h6 font-weight-medium">配送オプション：サイズ60</span>
      </v-card-title>
      <v-card-text class="pa-6">
        <div class="d-flex flex-column flex-md-row justify-center">
          <v-text-field
            v-model.number="validate.box60Frozen.$model"
            :error-messages="getErrorMessage(validate.box60Frozen.$errors)"
            label="冷凍配送価格"
            type="number"
            prefix="通常配送料＋"
            suffix="円"
            min="0"
            variant="outlined"
            density="comfortable"
            :prepend-inner-icon="mdiSnowflake"
          />
        </div>
        <div
          v-for="i in box60RateItemsSize"
          :key="`60-${i}`"
          class="shipping-rate-item"
        >
          <div class="d-flex flex-row align-center mb-3">
            <v-chip
              size="small"
              color="primary"
              variant="outlined"
            >
              オプション{{ i + 1 }}
            </v-chip>
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
            v-model="formData.box60Rates[i]"
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

    <v-card
      class="form-section-card mb-6"
      elevation="2"
      :loading="loading"
    >
      <v-card-title class="d-flex align-center section-header">
        <v-icon
          :icon="mdiPackageVariant"
          size="24"
          class="mr-3 text-primary"
        />
        <span class="text-h6 font-weight-medium">配送オプション：サイズ80</span>
      </v-card-title>
      <v-card-text class="pa-6">
        <div class="d-flex flex-column flex-md-row justify-center">
          <v-text-field
            v-model.number="validate.box80Frozen.$model"
            :error-messages="getErrorMessage(validate.box80Frozen.$errors)"
            label="冷凍配送価格"
            type="number"
            prefix="通常配送料＋"
            suffix="円"
            min="0"
            variant="outlined"
            density="comfortable"
            :prepend-inner-icon="mdiSnowflake"
          />
        </div>
        <div
          v-for="i in box80RateItemsSize"
          :key="`80-${i}`"
          class="shipping-rate-item"
        >
          <div class="d-flex flex-row align-center mb-3">
            <v-chip
              size="small"
              color="primary"
              variant="outlined"
            >
              オプション{{ i + 1 }}
            </v-chip>
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
            v-model="formData.box80Rates[i]"
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

    <v-card
      class="form-section-card mb-6"
      elevation="2"
      :loading="loading"
    >
      <v-card-title class="d-flex align-center section-header">
        <v-icon
          :icon="mdiPackageVariant"
          size="24"
          class="mr-3 text-primary"
        />
        <span class="text-h6 font-weight-medium">配送オプション：サイズ100</span>
      </v-card-title>
      <v-card-text class="pa-6">
        <div class="d-flex flex-column flex-md-row justify-center">
          <v-text-field
            v-model.number="validate.box100Frozen.$model"
            :error-messages="getErrorMessage(validate.box100Frozen.$errors)"
            label="冷凍配送価格"
            type="number"
            prefix="通常配送料＋"
            suffix="円"
            min="0"
            variant="outlined"
            density="comfortable"
            :prepend-inner-icon="mdiSnowflake"
          />
        </div>
        <div
          v-for="i in box100RateItemsSize"
          :key="`100-${i}`"
          class="shipping-rate-item"
        >
          <div class="d-flex flex-row align-center mb-3">
            <v-chip
              size="small"
              color="primary"
              variant="outlined"
            >
              オプション{{ i + 1 }}
            </v-chip>
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
            v-model="formData.box100Rates[i]"
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

    <!-- 送料無料設定セクション -->
    <v-card
      class="form-section-card mb-6"
      elevation="2"
      :loading="loading"
    >
      <v-card-title class="d-flex align-center section-header">
        <v-icon
          :icon="mdiCurrencyJpy"
          size="24"
          class="mr-3 text-primary"
        />
        <span class="text-h6 font-weight-medium">送料無料設定</span>
      </v-card-title>
      <v-card-text class="pa-6">
        <v-checkbox
          v-model="validate.hasFreeShipping.$model"
          label="送料無料オプションを有効にする"
          density="comfortable"
          class="mb-4"
        />
        <v-text-field
          v-if="formData.hasFreeShipping"
          v-model.number="validate.freeShippingRates.$model"
          label="送料無料になる購入金額"
          type="number"
          suffix="円以上"
          min="0"
          variant="outlined"
          density="comfortable"
        />
      </v-card-text>
    </v-card>

    <v-footer
      app
      color="white"
      elevation="8"
      class="px-6 py-4 fixed-footer-actions"
    >
      <v-container
        fluid
        class="pa-0"
      >
        <div class="d-flex align-center justify-center flex-wrap ga-3">
          <v-btn
            variant="text"
            size="large"
            @click="$router.back()"
          >
            キャンセル
          </v-btn>
          <v-btn
            color="primary"
            variant="elevated"
            size="large"
            :loading="isLoading()"
            @click="handleSubmit"
          >
            <v-icon
              :icon="mdiContentSave"
              start
            />
            変更を保存
          </v-btn>
        </div>
      </v-container>
    </v-footer>
  </v-container>
</template>

<style scoped>
.form-section-card {
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.2s ease;
  border: 1px solid rgb(0 0 0 / 5%);
}

.section-header {
  background: linear-gradient(90deg, rgb(33 150 243 / 5%) 0%, rgb(33 150 243 / 0%) 100%);
  border-bottom: 1px solid rgb(0 0 0 / 5%);
  padding: 20px 24px;
}

.shipping-rate-item {
  background: rgb(248 249 250);
  border: 1px solid rgb(0 0 0 / 6%);
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 16px;
  transition: all 0.2s ease;
}

.shipping-rate-item:hover {
  border-color: rgb(33 150 243 / 20%);
  background: rgb(248 249 250 / 80%);
}

@media (width <= 600px) {
  .form-section-card {
    border-radius: 8px;
  }

  .section-header {
    padding: 16px 20px;
  }

  .shipping-rate-item {
    padding: 12px;
  }
}
</style>
