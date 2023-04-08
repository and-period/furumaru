<script lang="ts" setup>
import { useVuelidate } from '@vuelidate/core'

import { getSelectablePrefecturesList } from '~/lib/prefectures'
import { required, getErrorMessage } from '~/lib/validations'
import { CreateShippingRequest, UpdateShippingRequest } from '~/types/api'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  value: {
    type: Object,
    default: (): UpdateShippingRequest | CreateShippingRequest => ({
      name: '',
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
  (
    e: 'update:value',
    formData: CreateShippingRequest | UpdateShippingRequest
  ): void
  (e: 'click:addBox60RateItem'): void
  (e: 'click:addBox80RateItem'): void
  (e: 'click:addBox100RateItem'): void
  (e: 'click:removeItemButton', rate: '60' | '80' | '100', index: number): void
  (e: 'submit'): void
}>()

const formData = computed({
  get: (): UpdateShippingRequest | CreateShippingRequest =>
    props.value as CreateShippingRequest | UpdateShippingRequest,
  set: (val: UpdateShippingRequest | CreateShippingRequest) =>
    emit('update:value', val)
})

const rules = computed(() => {
  return {
    name: { required },
    hasFreeShipping: { required }
  }
})

const v$ = useVuelidate(rules, formData)

const box60RateItemsSize = computed(() => {
  return [...Array(formData.value.box60Rates.length).keys()]
})

const box80RateItemsSize = computed(() => {
  return [...Array(formData.value.box80Rates.length).keys()]
})

const box100RateItemsSize = computed(() => {
  return [...Array(formData.value.box100Rates.length).keys()]
})

const getSelectableBox60RatePrefecturesList = (i: number) => {
  return getSelectablePrefecturesList(formData.value.box60Rates, i)
}

const getSelectableBox80RatePrefecturesList = (i: number) => {
  return getSelectablePrefecturesList(formData.value.box80Rates, i)
}

const getSelectableBox100RatePrefecturesList = (i: number) => {
  return getSelectablePrefecturesList(formData.value.box100Rates, i)
}

const handleClickSelectAll = (rate: '60' | '80' | '100', i: number) => {
  switch (rate) {
    case '60':
      formData.value.box60Rates[i].prefectures =
        getSelectableBox60RatePrefecturesList(i)
          .filter(item => !item.disabled)
          .map(item => item.value)
      break
    case '80':
      formData.value.box80Rates[i].prefectures =
        getSelectableBox80RatePrefecturesList(i)
          .filter(item => !item.disabled)
          .map(item => item.value)
      break
    case '100':
      formData.value.box100Rates[i].prefectures =
        getSelectableBox100RatePrefecturesList(i)
          .filter(item => !item.disabled)
          .map(item => item.value)
      break
  }
}

const addBox60RateItem = () => {
  emit('click:addBox60RateItem')
}

const addBox80RateItem = () => {
  emit('click:addBox80RateItem')
}

const addBox100RateItem = () => {
  emit('click:addBox100RateItem')
}

const handleSubmit = async () => {
  const result = await v$.value.$validate()
  if (!result) {
    return
  }
  emit('submit')
}

const handleClickRemoveItemButton = (
  rate: '60' | '80' | '100',
  index: number
) => {
  emit('click:removeItemButton', rate, index)
}
</script>

<template>
  <v-card :loading="props.loading">
    <form @submit.prevent="handleSubmit">
      <v-card-text>
        <v-text-field
          v-model="v$.name.$model"
          label="名前"
          :error-messages="getErrorMessage(v$.name.$errors)"
        />
        <v-switch
          v-model="v$.hasFreeShipping.$model"
          :label="
            v$.hasFreeShipping.$model
              ? `無料配送オプション: 有り`
              : `無料配送オプション: 無し`
          "
        />
        <div class="my-4">
          <p class="text-h6">
            サイズ60配送オプション
          </p>
          <div class="d-flex">
            <v-text-field
              v-model.number="formData.box60Refrigerated"
              label="冷蔵配送価格"
              class="mr-4"
            />
            <v-text-field
              v-model.number="formData.box60Frozen"
              label="冷凍配送価格"
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
                :disabled="box60RateItemsSize.length === 1"
                @click="handleClickRemoveItemButton('60', i)"
              >
                <v-icon>mdi-close</v-icon>
              </v-btn>
            </div>
            <v-text-field v-model="formData.box60Rates[i].name" label="名前" />
            <v-text-field
              v-model.number="formData.box60Rates[i].price"
              label="価格"
              type="number"
            />
            <v-select
              v-model="formData.box60Rates[i].prefectures"
              label="都道府県"
              chips
              multiple
              :items="getSelectableBox60RatePrefecturesList(i)"
            >
              <template #prepend-item>
                <v-list-item
                  ripple
                  @click="handleClickSelectAll('60', i)"
                  @mousedown.prevent
                >
                  <v-list-item-content>
                    <v-list-item-title> すべて選択 </v-list-item-title>
                  </v-list-item-content>
                </v-list-item>
              </template>
            </v-select>
          </div>
          <v-btn color="primary" outlined block @click="addBox60RateItem">
            <v-icon>mdi-plus</v-icon>
            追加
          </v-btn>
        </div>
        <div class="my-4">
          <p class="text-h6">
            サイズ80配送オプション
          </p>
          <div class="d-flex">
            <v-text-field
              v-model.number="formData.box80Refrigerated"
              label="冷蔵配送価格"
              class="mr-4"
            />
            <v-text-field
              v-model.number="formData.box80Frozen"
              label="冷凍配送価格"
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
                :disabled="box80RateItemsSize.length === 1"
                @click="handleClickRemoveItemButton('80', i)"
              >
                <v-icon>mdi-close</v-icon>
              </v-btn>
            </div>
            <v-text-field v-model="formData.box80Rates[i].name" label="名前" />
            <v-text-field
              v-model.number="formData.box80Rates[i].price"
              label="価格"
              type="number"
            />
            <v-select
              v-model="formData.box80Rates[i].prefectures"
              label="都道府県"
              chips
              multiple
              :items="getSelectableBox80RatePrefecturesList(i)"
            >
              <template #prepend-item>
                <v-list-item
                  ripple
                  @click="handleClickSelectAll('80', i)"
                  @mousedown.prevent
                >
                  <v-list-item-content>
                    <v-list-item-title> すべて選択 </v-list-item-title>
                  </v-list-item-content>
                </v-list-item>
              </template>
            </v-select>
          </div>
          <v-btn color="primary" outlined block @click="addBox80RateItem">
            <v-icon>mdi-plus</v-icon>
            追加
          </v-btn>
        </div>

        <div class="my-4">
          <p class="text-h6">
            サイズ100配送オプション
          </p>
          <div class="d-flex">
            <v-text-field
              v-model.number="formData.box100Refrigerated"
              label="冷蔵配送価格"
              class="mr-4"
            />
            <v-text-field
              v-model.number="formData.box100Frozen"
              label="冷凍配送価格"
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
                :disabled="box100RateItemsSize.length === 1"
                @click="handleClickRemoveItemButton('100', i)"
              >
                <v-icon>mdi-close</v-icon>
              </v-btn>
            </div>
            <v-text-field v-model="formData.box100Rates[i].name" label="名前" />
            <v-text-field
              v-model.number="formData.box100Rates[i].price"
              label="価格"
              type="number"
            />
            <v-select
              v-model="formData.box100Rates[i].prefectures"
              label="都道府県"
              chips
              multiple
              :items="getSelectableBox100RatePrefecturesList(i)"
            >
              <template #prepend-item>
                <v-list-item
                  ripple
                  @click="handleClickSelectAll('100', i)"
                  @mousedown.prevent
                >
                  <v-list-item-content>
                    <v-list-item-title> すべて選択 </v-list-item-title>
                  </v-list-item-content>
                </v-list-item>
              </template>
            </v-select>
          </div>
          <v-btn color="primary" outlined block @click="addBox100RateItem">
            <v-icon>mdi-plus</v-icon>
            追加
          </v-btn>
        </div>
      </v-card-text>
      <v-card-actions>
        <v-btn type="submit" outlined color="primary" block>
          登録
        </v-btn>
      </v-card-actions>
    </form>
  </v-card>
</template>
