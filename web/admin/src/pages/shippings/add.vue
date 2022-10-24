<template>
  <div>
    <v-card-title>配送情報登録</v-card-title>

    <v-alert v-model="isShow" :type="alertType" v-text="alertText" />

    <v-card>
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
            <p class="text-h6">サイズ60配送オプション</p>
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
                <p class="mb-0">オプション{{ i + 1 }}</p>
                <v-spacer />
                <v-btn
                  icon
                  :disabled="box60RateItemsSize.length === 1"
                  @click="handleClickCloseButton('60', i)"
                >
                  <v-icon>mdi-close</v-icon>
                </v-btn>
              </div>
              <v-text-field
                v-model="formData.box60Rates[i].name"
                label="名前"
              />
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
            <p class="text-h6">サイズ80配送オプション</p>
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
                <p class="mb-0">オプション{{ i + 1 }}</p>
                <v-spacer />
                <v-btn
                  icon
                  :disabled="box80RateItemsSize.length === 1"
                  @click="handleClickCloseButton('80', i)"
                  ><v-icon>mdi-close</v-icon></v-btn
                >
              </div>
              <v-text-field
                v-model="formData.box80Rates[i].name"
                label="名前"
              />
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
            <p class="text-h6">サイズ100配送オプション</p>
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
                <p class="mb-0">オプション{{ i + 1 }}</p>
                <v-spacer />
                <v-btn
                  icon
                  :disabled="box100RateItemsSize.length === 1"
                  @click="handleClickCloseButton('100', i)"
                  ><v-icon>mdi-close</v-icon></v-btn
                >
              </div>
              <v-text-field
                v-model="formData.box100Rates[i].name"
                label="名前"
              />
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
          <v-btn type="submit" outlined color="primary" block>登録</v-btn>
        </v-card-actions>
      </form>
    </v-card>
  </div>
</template>

<script lang="ts">
import {
  computed,
  defineComponent,
  reactive,
  useRouter,
} from '@nuxtjs/composition-api'
import { useVuelidate } from '@vuelidate/core'

import { prefecturesList } from '~/constants'
import { useAlert } from '~/lib/hooks'
import { getSelectablePrefecturesList } from '~/lib/prefectures'
import { required, getErrorMessage } from '~/lib/validations'
import { useCommonStore } from '~/store/common'
import { useShippingStore } from '~/store/shippings'
import { CreateShippingRequest } from '~/types/api'
import { ApiBaseError } from '~/types/exception'

export default defineComponent({
  setup() {
    const router = useRouter()

    const formData = reactive<CreateShippingRequest>({
      name: '',
      box60Rates: [
        {
          name: '',
          price: 0,
          prefectures: [],
        },
      ],
      box60Refrigerated: 0,
      box60Frozen: 0,
      box80Rates: [
        {
          name: '',
          price: 0,
          prefectures: [],
        },
      ],
      box80Refrigerated: 0,
      box80Frozen: 0,
      box100Rates: [
        {
          name: '',
          price: 0,
          prefectures: [],
        },
      ],
      box100Refrigerated: 0,
      box100Frozen: 0,
      hasFreeShipping: false,
      freeShippingRates: 0,
    })

    const rules = computed(() => {
      return {
        name: { required },
        hasFreeShipping: { required },
      }
    })

    const box60RateItemsSize = computed(() => {
      return [...Array(formData.box60Rates.length).keys()]
    })

    const box80RateItemsSize = computed(() => {
      return [...Array(formData.box80Rates.length).keys()]
    })

    const box100RateItemsSize = computed(() => {
      return [...Array(formData.box100Rates.length).keys()]
    })

    const addBox60RateItem = () => {
      formData.box60Rates.push({
        name: '',
        price: 0,
        prefectures: [],
      })
    }

    const addBox80RateItem = () => {
      formData.box80Rates.push({
        name: '',
        price: 0,
        prefectures: [],
      })
    }

    const addBox100RateItem = () => {
      formData.box100Rates.push({
        name: '',
        price: 0,
        prefectures: [],
      })
    }

    const getSelectableBox60RatePrefecturesList = (i: number) => {
      return getSelectablePrefecturesList(formData.box60Rates, i)
    }

    const getSelectableBox80RatePrefecturesList = (i: number) => {
      return getSelectablePrefecturesList(formData.box80Rates, i)
    }

    const getSelectableBox100RatePrefecturesList = (i: number) => {
      return getSelectablePrefecturesList(formData.box100Rates, i)
    }

    const handleClickSelectAll = (rate: '60' | '80' | '100', i: number) => {
      switch (rate) {
        case '60':
          formData.box60Rates[i].prefectures =
            getSelectableBox60RatePrefecturesList(i)
              .filter((item) => !item.disabled)
              .map((item) => item.value)
          break
        case '80':
          formData.box80Rates[i].prefectures =
            getSelectableBox80RatePrefecturesList(i)
              .filter((item) => !item.disabled)
              .map((item) => item.value)
          break
        case '100':
          formData.box100Rates[i].prefectures =
            getSelectableBox100RatePrefecturesList(i)
              .filter((item) => !item.disabled)
              .map((item) => item.value)
          break
      }
    }

    const v$ = useVuelidate(rules, formData)

    const { createShipping } = useShippingStore()

    const { alertType, isShow, alertText, show } = useAlert('error')

    const { addSnackbar } = useCommonStore()

    const handleSubmit = async (): Promise<void> => {
      const result = await v$.value.$validate()
      if (!result) {
        return
      }
      try {
        await createShipping(formData)
        addSnackbar({
          color: 'info',
          message: `${formData.name}を登録しました。`,
        })
        router.push('/shippings')
      } catch (error) {
        if (error instanceof ApiBaseError) {
          show(error.message)
          window.scrollTo({
            top: 0,
            behavior: 'smooth',
          })
        }
      }
    }

    const handleClickCloseButton = (
      rate: '60' | '80' | '100',
      index: number
    ) => {
      switch (rate) {
        case '60':
          formData.box60Rates.splice(index, 1)
          break
        case '80':
          formData.box80Rates.splice(index, 1)
          break
        case '100':
          formData.box100Rates.splice(index, 1)
          break
      }
    }

    return {
      formData,
      v$,
      box60RateItemsSize,
      box80RateItemsSize,
      box100RateItemsSize,
      prefecturesList,
      alertType,
      isShow,
      alertText,
      // 関数
      getErrorMessage,
      addBox60RateItem,
      addBox80RateItem,
      addBox100RateItem,
      handleSubmit,
      handleClickCloseButton,
      getSelectableBox60RatePrefecturesList,
      getSelectableBox80RatePrefecturesList,
      getSelectableBox100RatePrefecturesList,
      handleClickSelectAll,
    }
  },
})
</script>
