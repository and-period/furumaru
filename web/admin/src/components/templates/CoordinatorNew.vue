<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import { AlertType } from '~/lib/hooks'
import { CreateCoordinatorRequest } from '~/types/api'
import { email, getErrorMessage, maxLength, required, tel } from '~/lib/validations'
import { ImageUploadStatus } from '~/types/props'

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
      companyName: '',
      storeName: '',
      email: '',
      phoneNumber: '',
      postalCode: '',
      prefecture: '',
      city: '',
      addressLine1: '',
      addressLine2: ''
    })
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
  }
})

const emit = defineEmits<{
  (e: 'update:form-data', formData: CreateCoordinatorRequest): void
  (e: 'update:thumbnail-file', files: FileList): void
  (e: 'update:header-file', files: FileList): void
  (e: 'click:search-address'): void
  (e: 'submit'): void
}>()

const rules = computed(() => ({
  lastname: { required, maxLength: maxLength(16) },
  firstname: { required, maxLength: maxLength(16) },
  lastnameKana: { required, maxLength: maxLength(32) },
  firstnameKana: { required, maxLength: maxLength(32) },
  companyName: { required, maxLength: maxLength(64) },
  storeName: { required, maxLength: maxLength(64) },
  email: { required, email },
  phoneNumber: { required, tel }
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

const onSubmit = async (): Promise<void> => {
  const valid = await validate.value.$validate()
  if (!valid) {
    return
  }

  emit('submit')
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
        <v-text-field
          v-model="validate.companyName.$model"
          :error-messages="getErrorMessage(validate.companyName.$errors)"
          label="会社名"
        />
        <v-text-field
          v-model="validate.storeName.$model"
          :error-messages="getErrorMessage(validate.storeName.$errors)"
          label="店舗名"
        />
        <div class="mb-2 d-flex">
          <molecules-profile-select-form
            class="mr-4 flex-grow-1 flex-shrink-1"
            :img-url="formDataValue.thumbnailUrl"
            :error="props.thumbnailUploadStatus.error"
            :message="props.thumbnailUploadStatus.message"
            @update:file="onChangeThumbnailFile"
          />
          <molecules-header-select-form
            class="flex-grow-1 flex-shrink-1"
            :img-url="formDataValue.headerUrl"
            :error="props.headerUploadStatus.error"
            :message="props.headerUploadStatus.message"
            @update:file="onChangeHeaderFile"
          />
        </div>
        <div class="d-flex">
          <v-text-field
            v-model="validate.lastname.$model"
            :error-messages="getErrorMessage(validate.lastname.$errors)"
            class="mr-4"
            label="コーディネータ:姓"
          />
          <v-text-field
            v-model="validate.firstname.$model"
            :error-messages="getErrorMessage(validate.firstname.$errors)"
            label="コーディネータ:名"
          />
        </div>
        <div class="d-flex">
          <v-text-field
            v-model="validate.lastnameKana.$model"
            :error-messages="getErrorMessage(validate.lastnameKana.$errors)"
            class="mr-4"
            label="コーディネータ:姓（ふりがな）"
          />
          <v-text-field
            v-model="validate.firstnameKana.$model"
            :error-messages="getErrorMessage(validate.firstnameKana.$errors)"
            label="コーディネータ:名（ふりがな）"
          />
        </div>
        <v-text-field
          v-model="validate.email.$model"
          label="連絡先（Email）"
          :error-messages="getErrorMessage(validate.email.$errors)"
        />
        <v-text-field
          v-model="validate.phoneNumber.$model"
          :error-messages="getErrorMessage(validate.phoneNumber.$errors)"
          type="tel"
          label="連絡先（電話番号）"
        />

        <molecules-address-form
          v-model:postal-code="formDataValue.postalCode"
          v-model:prefecture="formDataValue.prefecture"
          v-model:city="formDataValue.city"
          v-model:address-line1="formDataValue.addressLine1"
          v-model:address-line2="formDataValue.addressLine2"
          @click:search="onClickSearchAddress"
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
