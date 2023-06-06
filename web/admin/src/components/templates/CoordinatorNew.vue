<script lang="ts" setup>
import { AlertType } from '~/lib/hooks'
import { CreateCoordinatorRequest } from '~/types/api'
import { email, getErrorMessage, maxLength, postalCode, required, tel } from '~/lib/validations'
import { ImageUploadStatus } from '~/types/props'
import useVuelidate from '@vuelidate/core'

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
      thumbnailUrl: '',
      headerUrl: '',
      twitterAccount: '',
      instagramAccount: '',
      facebookAccount: '',
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
  (e: 'update:form-data', v: CreateCoordinatorRequest): void
  (e: 'update:thumbnail-file', files?: FileList): void
  (e: 'update:header-file', files?: FileList): void
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
  thumbnailUrl: {},
  headerUrl: {},
  twitterAccount: {},
  instagramAccount: {},
  facebookAccount: {},
  email: { required, email },
  phoneNumber: { required, tel },
  postalCode: { maxLength: maxLength(7) },
  prefecture: { maxLength: maxLength(32) },
  city: { maxLength: maxLength(32) },
  addressLine1: { maxLength: maxLength(64) },
  addressLine2: { maxLength: maxLength(64) },
}))
const formDataValue = computed({
  get: (): CreateCoordinatorRequest => props.formData,
  set: (v: CreateCoordinatorRequest): void => emit('update:form-data', v)
})

const validate = useVuelidate(rules, formDataValue)

const updateThumbnailFileHandler = (files?: FileList) => {
  emit('update:thumbnail-file', files)
}

const updateHeaderFileHandler = (files?: FileList) => {
  emit('update:header-file', files)
}

const onSubmit = (): void => {
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

    <organisms-coordinator-create-form
      :form-data="formDataValue"
      :thumbnail-upload-status="props.thumbnailUploadStatus"
      :header-upload-status="props.headerUploadStatus"
      :search-loading="props.searchLoading"
      :search-error-message="props.searchErrorMessage"
      @update:thumbnail-file="updateThumbnailFileHandler"
      @update:header-file="updateHeaderFileHandler"
      @submit="onSubmit"
      @click:search="onClickSearchAddress"
    />
  </v-card>
</template>
