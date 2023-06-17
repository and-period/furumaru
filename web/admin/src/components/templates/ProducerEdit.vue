<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'

import { AlertType } from '~/lib/hooks'
import { getErrorMessage, maxLength, required, tel } from '~/lib/validations'
import { AdminStatus, ProducerResponse, UpdateProducerRequest } from '~/types/api'
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
  producer: {
    type: Object as PropType<ProducerResponse>,
    default: (): ProducerResponse => ({
      id: '',
      status: AdminStatus.UNKNOWN,
      coordinatorId: '',
      lastname: '',
      lastnameKana: '',
      firstname: '',
      firstnameKana: '',
      storeName: '',
      email: '',
      phoneNumber: '',
      postalCode: '',
      prefecture: '',
      city: '',
      addressLine1: '',
      addressLine2: '',
      thumbnailUrl: '',
      thumbnails: [],
      headerUrl: '',
      headers: [],
      createdAt: 0,
      updatedAt: 0
    })
  },
  formData: {
    type: Object as PropType<UpdateProducerRequest>,
    default: (): UpdateProducerRequest => ({
      lastname: '',
      lastnameKana: '',
      firstname: '',
      firstnameKana: '',
      storeName: '',
      thumbnailUrl: '',
      headerUrl: '',
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
  (e: 'update:form-data', formData: UpdateProducerRequest): void
  (e: 'update:producer', producer: ProducerResponse): void
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
  storeName: { required, maxLength: maxLength(64) },
  phoneNumber: { required, tel },
  postalCode: {},
  prefecture: {},
  city: {},
  addressLine1: {},
  addressLine2: {}
}))
const formDataValue = computed({
  get: (): UpdateProducerRequest => props.formData,
  set: (val: UpdateProducerRequest): void => emit('update:form-data', val)
})
const producerValue = computed({
  get: (): ProducerResponse => props.producer,
  set: (producer: ProducerResponse): void => emit('update:producer', producer)
})

const validate = useVuelidate(rules, formDataValue)

const updateThumbnailFileHandler = (files?: FileList): void => {
  if (!files) {
    return
  }

  emit('update:thumbnail-file', files)
}

const updateHeaderFileHandler = (files?: FileList): void => {
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
    <v-card-title>生産者編集</v-card-title>

    <v-skeleton-loader v-if="props.loading" type="article" />
    <v-form v-else @submit.prevent="onSubmit">
      <v-card-text>
        <v-text-field
          v-model="validate.storeName.$model"
          :error-messages="getErrorMessage(validate.storeName.$errors)"
          label="店舗名"
        />
        <div class="mb-2 d-flex">
          <molecules-profile-select-form
            class="mr-4 flex-grow-1 flex-shrink-1"
            :img-url="props.formData.thumbnailUrl"
            :error="props.thumbnailUploadStatus.error"
            :message="props.thumbnailUploadStatus.message"
            @update:file="updateThumbnailFileHandler"
          />
          <molecules-header-select-form
            class="flex-grow-1 flex-shrink-1"
            :img-url="props.formData.headerUrl"
            :error="props.headerUploadStatus.error"
            :message="props.headerUploadStatus.message"
            @update:file="updateHeaderFileHandler"
          />
        </div>
        <div class="d-flex">
          <v-text-field
            v-model="validate.lastname.$model"
            :error-messages="getErrorMessage(validate.lastname.$errors)"
            class="mr-4"
            label="生産者名:姓"
          />
          <v-text-field
            v-model="validate.firstname.$model"
            :error-messages="getErrorMessage(validate.firstname.$errors)"
            label="生産者名:名"
          />
        </div>
        <div class="d-flex">
          <v-text-field
            v-model="validate.lastnameKana.$model"
            :error-messages="getErrorMessage(validate.lastnameKana.$errors)"
            class="mr-4"
            label="生産者名:姓（ふりがな）"
          />
          <v-text-field
            v-model="validate.firstnameKana.$model"
            :error-messages="getErrorMessage(validate.firstnameKana.$errors)"
            label="生産者名:名（ふりがな）"
          />
        </div>
        <v-text-field
          v-model="producerValue.email"
          label="連絡先（Email）"
          type="email"
          readonly
        />
        <v-text-field
          v-model="validate.phoneNumber.$model"
          :error-messages="getErrorMessage(validate.phoneNumber.$errors)"
          label="連絡先（電話番号）"
        />

        <molecules-address-form
          v-model:postal-code="validate.postalCode.$model"
          v-model:prefecture="validate.prefecture.$model"
          v-model:city="validate.city.$model"
          v-model:address-line1="validate.addressLine1.$model"
          v-model:address-line2="validate.addressLine2.$model"
          :loading="loading"
          @click:search="onClickSearchAddress"
        />
      </v-card-text>
      <v-card-actions>
        <v-btn block :loading="loading" variant="outlined" color="primary" type="submit">
          更新
        </v-btn>
      </v-card-actions>
    </v-form>
  </v-card>
</template>
