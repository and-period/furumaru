<script lang="ts" setup>
import { mdiCheckCircle, mdiClose, mdiContentCopy, mdiContentSave, mdiCurrencyJpy, mdiDelete, mdiPackageVariant, mdiPlus, mdiSnowflake, mdiTruck } from '@mdi/js'
import useVuelidate from '@vuelidate/core'
import { unix } from 'dayjs'
import type { VDataTable } from 'vuetify/components'
import { getSelectablePrefecturesList } from '~/lib/prefectures'
import type { PrefecturesListSelectItems } from '~/lib/prefectures'
import { getErrorMessage } from '~/lib/validations'
import type { CreateShippingRequest, Shipping, UpdateShippingRequest } from '~/types/api/v1'
import { CreateShippingValidationRules, UpdateShippingValidationRules } from '~/types/validations'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  shipping: {
    type: Object as PropType<Shipping>,
    default: () => ({} as Shipping),
  },
  shippings: {
    type: Array<Shipping>,
    default: () => [],
  },
  tableItemsPerPage: {
    type: Number,
    default: 20,
  },
  tableItemsTotal: {
    type: Number,
    default: 0,
  },
  createFormData: {
    type: Object as PropType<CreateShippingRequest>,
    default: (): CreateShippingRequest => ({
      name: '',
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
  updateFormData: {
    type: Object as PropType<UpdateShippingRequest>,
    default: (): UpdateShippingRequest => ({
      name: '',
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
  createDialog: {
    type: Boolean,
    default: false,
  },
  updateDialog: {
    type: Boolean,
    default: false,
  },
  deleteDialog: {
    type: Boolean,
    default: false,
  },
  activeDialog: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits<{
  (e: 'update:create-form-data', formData: CreateShippingRequest): void
  (e: 'update:update-form-data', formData: UpdateShippingRequest): void
  (e: 'update:create-dialog', toggle: boolean): void
  (e: 'update:update-dialog', toggle: boolean): void
  (e: 'update:delete-dialog', toggle: boolean): void
  (e: 'update:active-dialog', toggle: boolean): void
  (e: 'click:update-page', page: number): void
  (e: 'click:update-items-per-page', page: number): void
  (e: 'click:create'): void
  (e: 'click:update', shippingId: string): void
  (e: 'click:delete', shippingId: string): void
  (e: 'click:active', shippingId: string): void
  (e: 'click:copy', shippingId: string): void
  (e: 'submit:create'): void
  (e: 'submit:update'): void
  (e: 'submit:delete'): void
  (e: 'submit:active'): void
}>()

const headers: VDataTable['headers'] = [
  {
    title: '配送設定名',
    key: 'name',
    sortable: false,
  },
  {
    title: 'デフォルト設定',
    key: 'isDefault',
    sortable: false,
  },
  {
    title: '登録日時',
    key: 'createdAt',
    sortable: false,
  },
  {
    title: '更新日時',
    key: 'updatedAt',
    sortable: false,
  },
  {
    title: '',
    key: 'actions',
    sortable: false,
  },
]

const createFormDataValue = computed({
  get: (): CreateShippingRequest => props.createFormData,
  set: (val: CreateShippingRequest): void => emit('update:create-form-data', val),
})
const updateFormDataValue = computed({
  get: (): UpdateShippingRequest => props.updateFormData,
  set: (val: UpdateShippingRequest): void => emit('update:update-form-data', val),
})
const createBox60RateItemsSize = computed(() => {
  return [...Array(createFormDataValue.value.box60Rates.length).keys()]
})
const createBox80RateItemsSize = computed(() => {
  return [...Array(createFormDataValue.value.box80Rates.length).keys()]
})
const createBox100RateItemsSize = computed(() => {
  return [...Array(createFormDataValue.value.box100Rates.length).keys()]
})
const updateBox60RateItemsSize = computed(() => {
  return [...Array(updateFormDataValue.value.box60Rates.length).keys()]
})
const updateBox80RateItemsSize = computed(() => {
  return [...Array(updateFormDataValue.value.box80Rates.length).keys()]
})
const updateBox100RateItemsSize = computed(() => {
  return [...Array(updateFormDataValue.value.box100Rates.length).keys()]
})
const createDialogValue = computed({
  get: (): boolean => props.createDialog,
  set: (val: boolean): void => emit('update:create-dialog', val),
})
const updateDialogValue = computed({
  get: (): boolean => props.updateDialog,
  set: (val: boolean): void => emit('update:update-dialog', val),
})
const deleteDialogValue = computed({
  get: (): boolean => props.deleteDialog,
  set: (val: boolean): void => emit('update:delete-dialog', val),
})
const activeDialogValue = computed({
  get: (): boolean => props.activeDialog,
  set: (val: boolean): void => emit('update:active-dialog', val),
})

const createFormDataValidate = useVuelidate(CreateShippingValidationRules, createFormDataValue)
const updateFormDataValidate = useVuelidate(UpdateShippingValidationRules, updateFormDataValue)

const addCreateBox60RateItem = () => {
  createFormDataValue.value.box60Rates.push({
    name: '',
    price: 0,
    prefectureCodes: [],
  })
}

const addCreateBox80RateItem = () => {
  createFormDataValue.value.box80Rates.push({
    name: '',
    price: 0,
    prefectureCodes: [],
  })
}

const addCreateBox100RateItem = () => {
  createFormDataValue.value.box100Rates.push({
    name: '',
    price: 0,
    prefectureCodes: [],
  })
}

const addUpdateBox60RateItem = () => {
  updateFormDataValue.value.box60Rates.push({
    name: '',
    price: 0,
    prefectureCodes: [],
  })
}

const addUpdateBox80RateItem = () => {
  updateFormDataValue.value.box80Rates.push({
    name: '',
    price: 0,
    prefectureCodes: [],
  })
}

const addUpdateBox100RateItem = () => {
  updateFormDataValue.value.box100Rates.push({
    name: '',
    price: 0,
    prefectureCodes: [],
  })
}

const removeCreateBoxItem = (rate: '60' | '80' | '100', index: number): void => {
  switch (rate) {
    case '60':
      createFormDataValue.value.box60Rates.splice(index, 1)
      break
    case '80':
      createFormDataValue.value.box80Rates.splice(index, 1)
      break
    case '100':
      createFormDataValue.value.box100Rates.splice(index, 1)
      break
  }
}

const removeUpdateBoxItem = (rate: '60' | '80' | '100', index: number): void => {
  switch (rate) {
    case '60':
      updateFormDataValue.value.box60Rates.splice(index, 1)
      break
    case '80':
      updateFormDataValue.value.box80Rates.splice(index, 1)
      break
    case '100':
      updateFormDataValue.value.box100Rates.splice(index, 1)
      break
  }
}

const getSelectableCreateBox60RatePrefecturesList = (i: number): PrefecturesListSelectItems[] => {
  return getSelectablePrefecturesList(createFormDataValue.value.box60Rates, i)
}

const getSelectableCreateBox80RatePrefecturesList = (i: number): PrefecturesListSelectItems[] => {
  return getSelectablePrefecturesList(createFormDataValue.value.box80Rates, i)
}

const getSelectableCreateBox100RatePrefecturesList = (i: number): PrefecturesListSelectItems[] => {
  return getSelectablePrefecturesList(createFormDataValue.value.box100Rates, i)
}

const getSelectableUpdateBox60RatePrefecturesList = (i: number): PrefecturesListSelectItems[] => {
  return getSelectablePrefecturesList(updateFormDataValue.value.box60Rates, i)
}

const getSelectableUpdateBox80RatePrefecturesList = (i: number): PrefecturesListSelectItems[] => {
  return getSelectablePrefecturesList(updateFormDataValue.value.box80Rates, i)
}

const getSelectableUpdateBox100RatePrefecturesList = (i: number): PrefecturesListSelectItems[] => {
  return getSelectablePrefecturesList(updateFormDataValue.value.box100Rates, i)
}

const selectCreateAllItem = (rate: '60' | '80' | '100', index: number): void => {
  switch (rate) {
    case '60':
      createFormDataValue.value.box60Rates[index].prefectureCodes
       = getSelectableCreateBox60RatePrefecturesList(index)
          .filter(item => !item.disabled)
          .map(item => item.value)
      break
    case '80':
      createFormDataValue.value.box80Rates[index].prefectureCodes
       = getSelectableCreateBox80RatePrefecturesList(index)
          .filter(item => !item.disabled)
          .map(item => item.value)
      break
    case '100':
      createFormDataValue.value.box100Rates[index].prefectureCodes
       = getSelectableCreateBox100RatePrefecturesList(index)
          .filter(item => !item.disabled)
          .map(item => item.value)
      break
  }
}

const selectUpdateAllItem = (rate: '60' | '80' | '100', index: number): void => {
  switch (rate) {
    case '60':
      updateFormDataValue.value.box60Rates[index].prefectureCodes
       = getSelectableUpdateBox60RatePrefecturesList(index)
          .filter(item => !item.disabled)
          .map(item => item.value)
      break
    case '80':
      updateFormDataValue.value.box80Rates[index].prefectureCodes
       = getSelectableUpdateBox80RatePrefecturesList(index)
          .filter(item => !item.disabled)
          .map(item => item.value)
      break
    case '100':
      updateFormDataValue.value.box100Rates[index].prefectureCodes
       = getSelectableUpdateBox100RatePrefecturesList(index)
          .filter(item => !item.disabled)
          .map(item => item.value)
      break
  }
}

const onClickCloseCreateDialog = (): void => {
  createDialogValue.value = false
}

const onClickCloseUpdateDialog = (): void => {
  updateDialogValue.value = false
}

const onClickCloseDeleteDialog = (): void => {
  deleteDialogValue.value = false
}

const onClickCloseActiveDialog = (): void => {
  activeDialogValue.value = false
}

const getShippingName = (shipping?: Shipping) => {
  if (!shipping) {
    return ''
  }
  return shipping.name
}

const getDateTimeString = (unixTime: number): string => {
  return unix(unixTime).format('YYYY/MM/DD HH:mm')
}

const getIsDefault = (isDefault: boolean): string => {
  return isDefault ? 'デフォルト' : '-'
}

const getIsDefaultColor = (isDefault: boolean): string => {
  return isDefault ? 'primary' : 'grey'
}

const isUpdatable = (): boolean => {
  return !props.shipping.isDefault
}

const isDeletable = (): boolean => {
  return props.shippings.length > 1
}

const onClickUpdatePage = (page: number): void => {
  emit('click:update-page', page)
}

const onClickUpdateItemsPerPage = (itemsPerPage: number): void => {
  emit('click:update-items-per-page', itemsPerPage)
}

const onClickCreate = (): void => {
  emit('click:create')
}

const onClickUpdate = (shippingId: string): void => {
  emit('click:update', shippingId)
}

const onClickDelete = (shippingId: string): void => {
  emit('click:delete', shippingId)
}

const onClickCopy = (shippingId: string): void => {
  emit('click:copy', shippingId)
}

const onClickActive = (shippingId: string): void => {
  emit('click:active', shippingId)
}

const onSubmitCreate = async (): Promise<void> => {
  const valid = await createFormDataValidate.value.$validate()
  if (!valid) {
    return
  }
  emit('submit:create')
}

const onSubmitUpdate = async (): Promise<void> => {
  const valid = await updateFormDataValidate.value.$validate()
  if (!valid) {
    return
  }
  emit('submit:update')
}

const onSubmitDelete = (): void => {
  emit('submit:delete')
}

const onSubmitActive = (): void => {
  emit('submit:active')
}
</script>

<template>
  <v-dialog
    v-model="createDialogValue"
    width="600"
    scrollable
  >
    <v-card>
      <v-card-title class="d-flex align-center pa-6 bg-blue-grey-lighten-5">
        <v-icon
          :icon="mdiPackageVariant"
          size="28"
          class="mr-3 text-primary"
        />
        <span class="text-h5 font-weight-medium">配送設定登録</span>
      </v-card-title>
      <v-card-text class="mt-6">
        <v-card
          class="form-section-card mb-6"
          elevation="2"
        >
          <v-card-title class="d-flex align-center section-header">
            <v-icon
              :icon="mdiTruck"
              size="24"
              class="mr-3 text-primary"
            />
            <span class="text-h6 font-weight-medium">配送設定名</span>
          </v-card-title>
          <v-card-text class="pa-6">
            <v-text-field
              v-model="createFormDataValidate.name.$model"
              :error-messages="getErrorMessage(createFormDataValidate.name.$errors)"
              label="配送設定名"
              placeholder="例）通常配送"
              variant="outlined"
              density="comfortable"
              autofocus
            />
          </v-card-text>
        </v-card>
        <v-card
          class="form-section-card mb-6"
          elevation="2"
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
                v-model.number="createFormDataValidate.box60Frozen.$model"
                :error-messages="getErrorMessage(createFormDataValidate.box60Frozen.$errors)"
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
              v-for="i in createBox60RateItemsSize"
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
                  v-show="createBox60RateItemsSize.length > 1"
                  :icon="mdiClose"
                  color="error"
                  variant="text"
                  size="small"
                  @click="removeCreateBoxItem('60', i)"
                />
              </div>
              <molecules-shipping-rate-form
                v-model="createFormDataValue.box60Rates[i]"
                :selectable-prefecture-list="getSelectableCreateBox60RatePrefecturesList(i)"
                @click:select-all="selectCreateAllItem('60', i)"
              />
            </div>
          </v-card-text>
          <v-card-actions class="pa-4">
            <v-btn
              color="primary"
              variant="tonal"
              block
              @click="addCreateBox60RateItem"
            >
              <v-icon :icon="mdiPlus" />
              オプション追加
            </v-btn>
          </v-card-actions>
        </v-card>

        <!-- サイズ80配送設定セクション -->
        <v-card
          class="form-section-card mb-6"
          elevation="2"
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
                v-model.number="createFormDataValidate.box80Frozen.$model"
                :error-messages="getErrorMessage(createFormDataValidate.box80Frozen.$errors)"
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
              v-for="i in createBox80RateItemsSize"
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
                  v-show="createBox80RateItemsSize.length > 1"
                  :icon="mdiClose"
                  color="error"
                  variant="text"
                  size="small"
                  @click="removeCreateBoxItem('80', i)"
                />
              </div>
              <molecules-shipping-rate-form
                v-model="createFormDataValue.box80Rates[i]"
                :selectable-prefecture-list="getSelectableCreateBox80RatePrefecturesList(i)"
                @click:select-all="selectCreateAllItem('80', i)"
              />
            </div>
          </v-card-text>
          <v-card-actions>
            <v-btn
              color="primary"
              variant="outlined"
              block
              @click="addCreateBox80RateItem"
            >
              <v-icon :icon="mdiPlus" />
              追加
            </v-btn>
          </v-card-actions>
        </v-card>

        <v-card
          class="form-section-card mb-6"
          elevation="2"
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
                v-model.number="createFormDataValidate.box100Frozen.$model"
                :error-messages="getErrorMessage(createFormDataValidate.box100Frozen.$errors)"
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
              v-for="i in createBox100RateItemsSize"
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
                  v-show="createBox100RateItemsSize.length > 1"
                  :icon="mdiClose"
                  color="error"
                  variant="text"
                  size="small"
                  @click="removeCreateBoxItem('100', i)"
                />
              </div>
              <molecules-shipping-rate-form
                v-model="createFormDataValue.box100Rates[i]"
                :selectable-prefecture-list="getSelectableCreateBox100RatePrefecturesList(i)"
                @click:select-all="selectCreateAllItem('100', i)"
              />
            </div>
          </v-card-text>
          <v-card-actions>
            <v-btn
              color="primary"
              variant="outlined"
              block
              @click="addCreateBox100RateItem"
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
              v-model="createFormDataValidate.hasFreeShipping.$model"
              label="送料無料オプションを有効にする"
              density="comfortable"
              class="mb-4"
            />
            <v-text-field
              v-if="createFormData.hasFreeShipping"
              v-model.number="createFormDataValidate.freeShippingRates.$model"
              label="送料無料になる購入金額"
              type="number"
              suffix="円以上"
              min="0"
              variant="outlined"
              density="comfortable"
            />
          </v-card-text>
        </v-card>
      </v-card-text>
      <v-divider />
      <v-card-actions class="d-flex justify-end pa-6 gap-3">
        <v-btn
          variant="text"
          size="large"
          @click="onClickCloseCreateDialog"
        >
          キャンセル
        </v-btn>
        <v-btn
          :loading="props.loading"
          color="primary"
          variant="elevated"
          size="large"
          :prepend-icon="mdiContentSave"
          @click="onSubmitCreate"
        >
          配送設定を登録
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog
    v-model="updateDialogValue"
    width="600"
    scrollable
  >
    <v-card>
      <v-card-title class="d-flex align-center pa-6 bg-blue-grey-lighten-5">
        <v-icon
          :icon="mdiPackageVariant"
          size="28"
          class="mr-3 text-primary"
        />
        <span class="text-h5 font-weight-medium">配送設定編集</span>
      </v-card-title>
      <v-card-text class="mt-6">
        <v-card
          class="form-section-card mb-6"
          elevation="2"
        >
          <v-card-title class="d-flex align-center section-header">
            <v-icon
              :icon="mdiTruck"
              size="24"
              class="mr-3 text-primary"
            />
            <span class="text-h6 font-weight-medium">配送設定名</span>
          </v-card-title>
          <v-card-text class="pa-6">
            <v-text-field
              v-model="updateFormDataValidate.name.$model"
              :error-messages="getErrorMessage(updateFormDataValidate.name.$errors)"
              label="配送設定名"
              placeholder="例）通常配送"
              variant="outlined"
              density="comfortable"
              autofocus
            />
          </v-card-text>
        </v-card>
        <v-card
          class="form-section-card mb-6"
          elevation="2"
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
                v-model.number="updateFormDataValidate.box60Frozen.$model"
                :error-messages="getErrorMessage(updateFormDataValidate.box60Frozen.$errors)"
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
              v-for="i in updateBox60RateItemsSize"
              :key="`60-${i}`"
              class="px-4 py-2 mb-2 border"
            >
              <div class="d-flex flex-row align-center">
                <p class="text-subtitle-2 text-grey">
                  オプション{{ i + 1 }}
                </p>
                <v-spacer />
                <v-btn
                  v-show="updateBox60RateItemsSize.length > 1"
                  :icon="mdiClose"
                  color="error"
                  variant="text"
                  size="small"
                  @click="removeCreateBoxItem('60', i)"
                />
              </div>
              <molecules-shipping-rate-form
                v-model="updateFormDataValue.box60Rates[i]"
                :selectable-prefecture-list="getSelectableUpdateBox60RatePrefecturesList(i)"
                @click:select-all="selectUpdateAllItem('60', i)"
              />
            </div>
          </v-card-text>
          <v-card-actions class="pa-4">
            <v-btn
              color="primary"
              variant="tonal"
              block
              @click="addUpdateBox60RateItem"
            >
              <v-icon :icon="mdiPlus" />
              オプション追加
            </v-btn>
          </v-card-actions>
        </v-card>

        <!-- サイズ80配送設定セクション -->
        <v-card
          class="form-section-card mb-6"
          elevation="2"
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
                v-model.number="updateFormDataValidate.box80Frozen.$model"
                :error-messages="getErrorMessage(updateFormDataValidate.box80Frozen.$errors)"
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
              v-for="i in updateBox80RateItemsSize"
              :key="`80-${i}`"
              class="px-4 py-2 mb-2 border"
            >
              <div class="d-flex flex-row align-center">
                <p class="text-subtitle-2 text-grey">
                  オプション{{ i + 1 }}
                </p>
                <v-spacer />
                <v-btn
                  v-show="updateBox80RateItemsSize.length > 1"
                  :icon="mdiClose"
                  color="error"
                  variant="text"
                  size="small"
                  @click="removeUpdateBoxItem('80', i)"
                />
              </div>
              <molecules-shipping-rate-form
                v-model="updateFormDataValue.box80Rates[i]"
                :selectable-prefecture-list="getSelectableUpdateBox80RatePrefecturesList(i)"
                @click:select-all="selectUpdateAllItem('80', i)"
              />
            </div>
          </v-card-text>
          <v-card-actions>
            <v-btn
              color="primary"
              variant="outlined"
              block
              @click="addUpdateBox80RateItem"
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
                v-model.number="updateFormDataValidate.box100Frozen.$model"
                :error-messages="getErrorMessage(updateFormDataValidate.box100Frozen.$errors)"
                label="冷凍配送価格"
                type="number"
                prefix="通常配送料＋"
                suffix="円"
                min="0"
              />
            </div>
            <div
              v-for="i in updateBox100RateItemsSize"
              :key="`100-${i}`"
              class="px-4 py-2 mb-2 border"
            >
              <div class="d-flex flex-row align-center">
                <p class="text-subtitle-2 text-grey">
                  オプション{{ i + 1 }}
                </p>
                <v-spacer />
                <v-btn
                  v-show="updateBox100RateItemsSize.length > 1"
                  :icon="mdiClose"
                  color="error"
                  variant="text"
                  size="small"
                  @click="removeUpdateBoxItem('100', i)"
                />
              </div>
              <molecules-shipping-rate-form
                v-model="updateFormDataValue.box100Rates[i]"
                :selectable-prefecture-list="getSelectableUpdateBox100RatePrefecturesList(i)"
                @click:select-all="selectUpdateAllItem('100', i)"
              />
            </div>
          </v-card-text>
          <v-card-actions class="pa-4">
            <v-btn
              color="primary"
              variant="tonal"
              block
              @click="addUpdateBox100RateItem"
            >
              <v-icon :icon="mdiPlus" />
              オプション追加
            </v-btn>
          </v-card-actions>
        </v-card>

        <!-- 送料無料設定セクション -->
        <v-card
          class="form-section-card mb-6"
          elevation="2"
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
              v-model="updateFormDataValidate.hasFreeShipping.$model"
              label="送料無料オプションを有効にする"
              density="comfortable"
              class="mb-4"
            />
            <v-text-field
              v-if="updateFormData.hasFreeShipping"
              v-model.number="updateFormDataValidate.freeShippingRates.$model"
              label="送料無料になる購入金額"
              type="number"
              suffix="円以上"
              min="0"
              variant="outlined"
              density="comfortable"
            />
          </v-card-text>
        </v-card>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn
          color="error"
          variant="text"
          @click="onClickCloseUpdateDialog"
        >
          {{ isUpdatable() ? 'キャンセル' : '閉じる' }}
        </v-btn>
        <v-btn
          v-if="isUpdatable()"
          :loading="props.loading"
          color="primary"
          variant="outlined"
          :prepend-icon="mdiContentSave"
          @click="onSubmitUpdate"
        >
          更新
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog
    v-model="deleteDialogValue"
    width="500"
  >
    <v-card>
      <v-card-title class="text-h6 py-4">
        配送設定削除の確認
      </v-card-title>
      <v-card-text class="pb-4">
        <div class="text-body-1">
          「{{ getShippingName(shipping) }}」を削除しますか？
        </div>
        <div class="text-body-2 text-medium-emphasis mt-2">
          この操作は取り消せません。
        </div>
      </v-card-text>
      <v-card-actions class="px-6 pb-4">
        <v-spacer />
        <v-btn
          color="medium-emphasis"
          variant="text"
          @click="onClickCloseDeleteDialog"
        >
          キャンセル
        </v-btn>
        <v-btn
          :loading="props.loading"
          color="error"
          variant="elevated"
          @click="onSubmitDelete"
        >
          削除
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog
    v-model="activeDialogValue"
    width="500"
  >
    <v-card>
      <v-card-title class="text-h6 py-4">
        配送設定有効化前確認
      </v-card-title>
      <v-card-text class="pb-4">
        <div class="text-body-1">
          「{{ getShippingName(shipping) }}」を有効化しますか？
        </div>
      </v-card-text>
      <v-card-actions class="px-6 pb-4">
        <v-spacer />
        <v-btn
          color="medium-emphasis"
          variant="text"
          @click="onClickCloseActiveDialog"
        >
          キャンセル
        </v-btn>
        <v-btn
          :loading="props.loading"
          color="primary"
          variant="outlined"
          @click="onSubmitActive"
        >
          有効化
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-card
    class="form-section-card"
    elevation="2"
  >
    <v-card-title class="d-flex align-center justify-space-between section-header">
      <div class="d-flex align-center">
        <v-icon
          :icon="mdiPackageVariant"
          size="24"
          class="mr-3 text-primary"
        />
        <span class="text-h6 font-weight-medium">配送設定一覧</span>
      </div>
      <v-btn
        variant="elevated"
        color="primary"
        size="default"
        :prepend-icon="mdiPlus"
        class="text-none"
        @click="onClickCreate"
      >
        新しい配送設定を追加
      </v-btn>
    </v-card-title>
    <v-card-text class="pa-6">
      <div
        v-if="props.shippings.length > 0"
        class="shipping-table-container"
      >
        <v-data-table-server
          :loading="props.loading"
          :headers="headers"
          :items="props.shippings"
          :items-per-page="props.tableItemsPerPage"
          :items-length="props.tableItemsTotal"
          hover
          class="elevation-1 rounded-lg"
          no-data-text="配送情報が登録されていません"
          @update:page="onClickUpdatePage"
          @update:items-per-page="onClickUpdateItemsPerPage"
          @click:row="(_:any, { item }: any) => onClickUpdate(item.id)"
        >
          <template #[`item.isDefault`]="{ item }">
            <v-chip
              :color="getIsDefaultColor(item.isDefault)"
              variant="elevated"
              size="small"
              class="font-weight-medium"
            >
              {{ getIsDefault(item.isDefault) }}
            </v-chip>
          </template>
          <template #[`item.createdAt`]="{ item }">
            <span class="text-body-2 text-grey-darken-1">
              {{ getDateTimeString(item.createdAt) }}
            </span>
          </template>
          <template #[`item.updatedAt`]="{ item }">
            <span class="text-body-2 text-grey-darken-1">
              {{ getDateTimeString(item.updatedAt) }}
            </span>
          </template>
          <template #[`item.actions`]="{ item }">
            <div class="d-flex gap-2">
              <v-btn
                v-show="!item.isDefault"
                variant="outlined"
                color="primary"
                size="small"
                :prepend-icon="mdiCheckCircle"
                class="text-none"
                @click.stop="onClickActive(item.id)"
              >
                有効化
              </v-btn>
              <v-btn
                variant="outlined"
                color="secondary"
                size="small"
                :prepend-icon="mdiContentCopy"
                class="text-none"
                @click.stop="onClickCopy(item.id)"
              >
                複製
              </v-btn>
              <v-btn
                v-show="isDeletable()"
                variant="outlined"
                color="error"
                size="small"
                :prepend-icon="mdiDelete"
                class="text-none"
                @click.stop="onClickDelete(item.id)"
              >
                削除
              </v-btn>
            </div>
          </template>
        </v-data-table-server>
      </div>
      <div
        v-else
        class="empty-state"
      >
        <v-icon
          :icon="mdiPackageVariant"
          size="64"
          class="text-grey-lighten-1 mb-4"
        />
        <h3 class="text-h6 text-grey-darken-1 mb-2">
          配送情報が登録されていません
        </h3>
        <p class="text-body-2 text-grey-darken-1 mb-4">
          配送料金や配送オプションを設定してください
        </p>
        <v-btn
          color="primary"
          variant="elevated"
          :prepend-icon="mdiPlus"
          @click="onClickCreate"
        >
          最初の配送設定を追加
        </v-btn>
      </div>
    </v-card-text>
  </v-card>
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

.shipping-table-container {
  background: rgb(250 250 250);
  border-radius: 8px;
  padding: 16px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: 64px 32px;
  background: rgb(250 250 250);
  border-radius: 12px;
  border: 2px dashed rgb(0 0 0 / 10%);
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

  .empty-state {
    padding: 48px 24px;
  }

  .shipping-rate-item {
    padding: 12px;
  }
}
</style>
