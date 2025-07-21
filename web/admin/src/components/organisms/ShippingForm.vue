<script setup lang="ts">
import { mdiClose, mdiPlus } from '@mdi/js'
import useVuelidate from '@vuelidate/core'
import type { CreateShippingRequest, UpdateShippingRequest } from '~/types/api'

import { getErrorMessage } from '~/lib/validations'
import { getSelectablePrefecturesList } from '~/lib/prefectures'
import type { PrefecturesListSelectItems } from '~/lib/prefectures'
import { UpsertShippingValidationRules } from '~/types/validations'

interface Props {
  formType: 'create' | 'update'
  submitting: boolean
  loading: boolean
}

defineProps<Props>()

interface Emits {
  (e: 'submit'): void
}

const emit = defineEmits<Emits>()

const formData = defineModel<CreateShippingRequest | UpdateShippingRequest>({ required: true })

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

const validate = useVuelidate(UpsertShippingValidationRules, formData)

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

/**
 * フォームの送信処理を行う
 * バリデーションが成功した場合のみsubmitイベントを発行する
 */
const handleSubmit = async () => {
  const valid = await validate.value.$validate()
  if (!valid) {
    return
  }

  emit('submit')
}
</script>

<template>
  <div class="mb-16">
    <v-card
      class="mb-4"
      :loading="loading"
    >
      <v-card-text>
        <v-text-field
          v-model="validate.name.$model"
          :error-messages="getErrorMessage(validate.name.$errors)"
          label="配送設定名"
          type="text"
        />
      </v-card-text>
    </v-card>
    <v-card
      class="mb-4 py-2"
      :loading="loading"
    >
      <v-card-title>配送オプション：サイズ60</v-card-title>
      <v-card-text>
        <div class="d-flex flex-column flex-md-row justify-center">
          <v-text-field
            v-model.number="validate.box60Frozen.$model"
            :error-messages="getErrorMessage(validate.box60Frozen.$errors)"
            label="クール配送価格"
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
      class="mb-4 py-2"
      :loading="loading"
    >
      <v-card-title>配送オプション：サイズ80</v-card-title>
      <v-card-text>
        <div class="d-flex flex-column flex-md-row justify-center">
          <v-text-field
            v-model.number="validate.box80Frozen.$model"
            :error-messages="getErrorMessage(validate.box80Frozen.$errors)"
            label="クール配送価格"
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
      class="mb-4 py-2"
      :loading="loading"
    >
      <v-card-title>配送オプション：サイズ100</v-card-title>
      <v-card-text>
        <div class="d-flex flex-column flex-md-row justify-center">
          <v-text-field
            v-model.number="validate.box100Frozen.$model"
            :error-messages="getErrorMessage(validate.box100Frozen.$errors)"
            label="クール配送価格"
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
  </div>
</template>
