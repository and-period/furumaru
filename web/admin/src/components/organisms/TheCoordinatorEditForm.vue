<script lang="ts" setup>
import { useVuelidate } from '@vuelidate/core'

import {
  kana,
  getErrorMessage,
  required,
  tel,
  maxLength,
} from '~/lib/validations'
import { UpdateCoordinatorRequest } from '~/types/api'
import { ImageUploadStatus } from '~/types/props'

const props = defineProps({
  formData: {
    type: Object,
    default: (): UpdateCoordinatorRequest => {
      return {
        lastname: '',
        lastnameKana: '',
        firstname: '',
        firstnameKana: '',
        companyName: '',
        storeName: '',
        thumbnailUrl: '',
        headerUrl: '',
        phoneNumber: '',
        postalCode: '',
        prefecture: '',
        city: '',
        addressLine1: '',
        addressLine2: '',
      }
    },
  },
  thumbnailUploadStatus: {
    type: Object,
    default: (): ImageUploadStatus => {
      return {
        error: false,
        message: '',
      }
    },
  },
  headerUploadStatus: {
    type: Object,
    default: (): ImageUploadStatus => {
      return {
        error: false,
        message: '',
      }
    },
  },
  searchErrorMessage: {
    type: String,
    default: '',
  },
  searchLoading: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits<{
  (e: 'update:formData', formData: UpdateCoordinatorRequest): void
  (e: 'update:thumbnailFile', files?: FileList): void
  (e: 'update:headerFile', files?: FileList): void
  (e: 'click:search'): void
  (e: 'submit'): void
}>()

const formDataValue = computed({
  get: (): UpdateCoordinatorRequest =>
    props.formData as UpdateCoordinatorRequest,
  set: (val: UpdateCoordinatorRequest) => emit('update:formData', val),
})

const rules = computed(() => ({
  storeName: { required, maxLength: maxLength(64) },
  companyName: { required, maxLength: maxLength(64) },
  firstname: { required, maxLength: maxLength(16) },
  lastname: { required, maxLength: maxLength(16) },
  firstnameKana: { required, kana },
  lastnameKana: { required, kana },
  phoneNumber: { required, tel },
}))

const v$ = useVuelidate(rules, formDataValue)

const updateThumbnailFileHandler = (files?: FileList) => {
  emit('update:thumbnailFile', files)
}

const updateHeaderFileHandler = (files?: FileList) => {
  emit('update:headerFile', files)
}

const handleSearchClick = () => {
  emit('click:search')
}

const handleSubmit = async () => {
  const result = await v$.value.$validate()
  if (!result) {
    return
  }

  emit('submit')
}
</script>

<template>
  <form @submit.prevent="handleSubmit">
    <v-card>
      <v-card-text>
        <v-text-field
          v-model="v$.companyName.$model"
          :error-messages="getErrorMessage(v$.companyName.$errors)"
          label="会社名"
        />
        <v-text-field
          v-model="v$.storeName.$model"
          :error-messages="getErrorMessage(v$.storeName.$errors)"
          label="店舗名"
        />
        <div class="mb-2 d-flex">
          <the-profile-select-form
            class="mr-4 flex-grow-1 flex-shrink-1"
            :img-url="props.formData.thumbnailUrl"
            :error="props.thumbnailUploadStatus.error"
            :message="props.thumbnailUploadStatus.message"
            @update:file="updateThumbnailFileHandler"
          />
          <the-header-select-form
            class="flex-grow-1 flex-shrink-1"
            :img-url="props.formData.headerUrl"
            :error="props.headerUploadStatus.error"
            :message="props.headerUploadStatus.message"
            @update:file="updateHeaderFileHandler"
          />
        </div>

        <div class="d-flex">
          <v-text-field
            v-model="v$.lastname.$model"
            :error-messages="getErrorMessage(v$.lastname.$errors)"
            class="mr-4"
            label="コーディネータ:姓"
          />
          <v-text-field
            v-model="v$.firstname.$model"
            :error-messages="getErrorMessage(v$.firstname.$errors)"
            label="コーディネータ:名"
          />
        </div>
        <div class="d-flex">
          <v-text-field
            v-model="v$.lastnameKana.$model"
            :error-messages="getErrorMessage(v$.lastnameKana.$errors)"
            class="mr-4"
            label="コーディネータ:姓（ふりがな）"
          />
          <v-text-field
            v-model="v$.firstnameKana.$model"
            :error-messages="getErrorMessage(v$.firstnameKana.$errors)"
            label="コーディネータ:名（ふりがな）"
          />
        </div>

        <v-text-field
          v-model="v$.phoneNumber.$model"
          :error-messages="getErrorMessage(v$.phoneNumber.$errors)"
          type="tel"
          label="連絡先（電話番号）"
        />

        <the-address-form
          :postal-code.sync="props.formData.postalCode"
          :prefecture.sync="props.formData.prefecture"
          :city.sync="props.formData.city"
          :address-line1.sync="props.formData.addressLine1"
          :address-line2.sync="props.formData.addressLine2"
          :error-message="props.searchErrorMessage"
          :loading="props.searchLoading"
          @click:search="handleSearchClick"
        />
      </v-card-text>
      <v-card-actions>
        <v-btn block outlined color="primary" type="submit">更新</v-btn>
      </v-card-actions>
    </v-card>
  </form>
</template>
