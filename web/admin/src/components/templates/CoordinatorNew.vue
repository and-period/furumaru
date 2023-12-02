<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import { mdiFacebook, mdiInstagram } from '@mdi/js'
import type { AlertType } from '~/lib/hooks'
import { type CreateCoordinatorRequest, Prefecture, type ProductType, Weekday } from '~/types/api'
import { email, kana, getErrorMessage, maxLength, required, tel } from '~/lib/validations'
import type { ImageUploadStatus } from '~/types/props'

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
    type: Object as PropType<CreateCoordinatorRequest>,
    default: (): CreateCoordinatorRequest => ({
      lastname: '',
      lastnameKana: '',
      firstname: '',
      firstnameKana: '',
      marcheName: '',
      username: '',
      email: '',
      phoneNumber: '',
      postalCode: '',
      prefectureCode: Prefecture.UNKNOWN,
      city: '',
      addressLine1: '',
      addressLine2: '',
      profile: '',
      productTypeIds: [],
      thumbnailUrl: '',
      headerUrl: '',
      promotionVideoUrl: '',
      bonusVideoUrl: '',
      instagramId: '',
      facebookId: '',
      businessDays: []
    })
  },
  productTypes: {
    type: Array<ProductType>,
    default: () => []
  },
  thumbnailUploadStatus: {
    type: Object,
    default: (): ImageUploadStatus => ({
      error: false,
      message: ''
    })
  },
  headerUploadStatus: {
    type: Object,
    default: (): ImageUploadStatus => ({
      error: false,
      message: ''
    })
  },
  promotionVideoUploadStatus: {
    type: Object,
    default: (): ImageUploadStatus => ({
      error: false,
      message: ''
    })
  },
  bonusVideoUploadStatus: {
    type: Object,
    default: (): ImageUploadStatus => ({
      error: false,
      message: ''
    })
  },
  searchErrorMessage: {
    type: String,
    default: ''
  },
  searchLoading: {
    type: Boolean,
    default: false
  }
})

const weekdays = [
  { title: '日曜日', value: Weekday.SUNDAY },
  { title: '月曜日', value: Weekday.MONDAY },
  { title: '火曜日', value: Weekday.TUESDAY },
  { title: '水曜日', value: Weekday.WEDNESDAY },
  { title: '木曜日', value: Weekday.THURSDAY },
  { title: '金曜日', value: Weekday.FRIDAY },
  { title: '土曜日', value: Weekday.SATURDAY }
]

const emit = defineEmits<{
  (e: 'update:form-data', formData: CreateCoordinatorRequest): void
  (e: 'update:thumbnail-file', files: FileList): void
  (e: 'update:header-file', files: FileList): void
  (e: 'update:promotion-video', files: FileList): void
  (e: 'update:bonus-video', files: FileList): void
  (e: 'update:search-product-type', name: string): void
  (e: 'click:search-address'): void
  (e: 'submit'): void
}>()

const rules = computed(() => ({
  lastname: { required, maxLength: maxLength(16) },
  firstname: { required, maxLength: maxLength(16) },
  lastnameKana: { required, kana, maxLength: maxLength(32) },
  firstnameKana: { required, kana, maxLength: maxLength(32) },
  marcheName: { required, maxLength: maxLength(64) },
  username: { required, maxLength: maxLength(64) },
  email: { required, email },
  phoneNumber: { required, tel },
  profile: { maxLength: maxLength(2000) },
  instagramId: { maxLength: maxLength(30) },
  facebookId: { maxLength: maxLength(50) },
  businessDays: {}
}))
const formDataValue = computed({
  get: (): CreateCoordinatorRequest => props.formData,
  set: (formData: CreateCoordinatorRequest): void => emit('update:form-data', formData)
})

const validate = useVuelidate(rules, formDataValue)

const onChangeThumbnailFile = (files?: FileList) => {
  if (!files) {
    return
  }
  emit('update:thumbnail-file', files)
}

const onChangeHeaderFile = (files?: FileList) => {
  if (!files) {
    return
  }
  emit('update:header-file', files)
}

const onChangePromotionVideo = (files?: FileList) => {
  if (!files) {
    return
  }
  emit('update:promotion-video', files)
}

const onChangeBonusVideo = (files?: FileList) => {
  if (!files) {
    return
  }
  emit('update:bonus-video', files)
}

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

const onClickSearchAddress = (): void => {
  emit('click:search-address')
}
</script>

<template>
  <v-alert v-show="props.isAlert" :type="props.alertType" v-text="props.alertText" />

  <v-card>
    <v-card-title>コーディネーター登録</v-card-title>

    <v-form @submit.prevent="onSubmit">
      <v-card-text>
        <v-row>
          <v-col>
            <v-text-field
              v-model="validate.marcheName.$model"
              :error-messages="getErrorMessage(validate.marcheName.$errors)"
              class="mr-4"
              label="マルシェ名"
            />
          </v-col>
          <v-col>
            <v-text-field
              v-model="validate.username.$model"
              :error-messages="getErrorMessage(validate.username.$errors)"
              class="mr-4"
              label="コーディネーター名"
            />
          </v-col>
        </v-row>
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
              :prepend-avatar="item.raw.iconUrl"
              :text="item.raw.name"
              rounded
              class="px-4"
              variant="outlined"
            />
          </template>
          <template #item="{ props: val, item }">
            <v-list-item
              v-bind="val"
              :prepend-avatar="item?.raw.iconUrl"
              :title="item?.raw.name"
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
        <v-row>
          <v-col cols="12" ms="12" lg="6">
            <molecules-video-select-form
              label="紹介動画"
              :video-url="formDataValue.promotionVideoUrl"
              :error="props.promotionVideoUploadStatus.error"
              :message="props.promotionVideoUploadStatus.message"
              @update:file="onChangePromotionVideo"
            />
          </v-col>
          <v-col cols="12" sm="12" lg="6">
            <molecules-video-select-form
              label="サンキュー動画"
              :video-url="formDataValue.bonusVideoUrl"
              :error="props.bonusVideoUploadStatus.error"
              :message="props.bonusVideoUploadStatus.message"
              @update:file="onChangeBonusVideo"
            />
          </v-col>
        </v-row>
        <v-textarea
          v-model="validate.profile.$model"
          :error-messages="getErrorMessage(validate.profile.$errors)"
          label="プロフィール"
          maxlength="2000"
        />
        <v-row>
          <v-col>
            <v-text-field
              v-model="validate.lastname.$model"
              :error-messages="getErrorMessage(validate.lastname.$errors)"
              class="mr-4"
              label="コーディネーター:姓"
            />
          </v-col>
          <v-col>
            <v-text-field
              v-model="validate.firstname.$model"
              :error-messages="getErrorMessage(validate.firstname.$errors)"
              label="コーディネーター:名"
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-text-field
              v-model="validate.lastnameKana.$model"
              :error-messages="getErrorMessage(validate.lastnameKana.$errors)"
              class="mr-4"
              label="コーディネーター:姓（ふりがな）"
            />
          </v-col>
          <v-col>
            <v-text-field
              v-model="validate.firstnameKana.$model"
              :error-messages="getErrorMessage(validate.firstnameKana.$errors)"
              label="コーディネーター:名（ふりがな）"
            />
          </v-col>
        </v-row>
        <v-text-field
          v-model="validate.email.$model"
          :error-messages="getErrorMessage(validate.email.$errors)"
          label="連絡先（Email）"
        />
        <v-text-field
          v-model="validate.phoneNumber.$model"
          :error-messages="getErrorMessage(validate.phoneNumber.$errors)"
          type="tel"
          label="連絡先（電話番号）"
        />
        <v-row>
          <v-col cols="12" sm="6" md="6">
            <molecules-icon-select-form
              label="アイコン画像"
              :img-url="formDataValue.thumbnailUrl"
              :error="props.thumbnailUploadStatus.error"
              :message="props.thumbnailUploadStatus.message"
              @update:file="onChangeThumbnailFile"
            />
          </v-col>
          <v-col cols="12" sm="6" md="6">
            <molecules-image-select-form
              label="ヘッダー画像"
              :img-url="formDataValue.headerUrl"
              :error="props.headerUploadStatus.error"
              :message="props.headerUploadStatus.message"
              @update:file="onChangeHeaderFile"
            />
          </v-col>
        </v-row>
        <molecules-address-form
          v-model:postal-code="formDataValue.postalCode"
          v-model:prefecture="formDataValue.prefectureCode"
          v-model:city="formDataValue.city"
          v-model:address-line1="formDataValue.addressLine1"
          v-model:address-line2="formDataValue.addressLine2"
          :error-message="props.searchErrorMessage"
          :loading="props.searchLoading"
          @click:search="onClickSearchAddress"
        />
        <p>SNS連携</p>
        <v-text-field
          v-model="validate.instagramId.$model"
          :error-messages="getErrorMessage(validate.instagramId.$errors)"
          :prepend-icon="mdiInstagram"
          prefix="https://www.instagram.com/"
        />
        <v-text-field
          v-model="validate.facebookId.$model"
          :error-messages="getErrorMessage(validate.facebookId.$errors)"
          :prepend-icon="mdiFacebook"
          prefix="https://www.facebook.com/"
        />
      </v-card-text>

      <v-card-actions>
        <v-btn block :loading="loading" variant="outlined" color="primary" type="submit">
          登録
        </v-btn>
      </v-card-actions>
    </v-form>
  </v-card>
</template>
