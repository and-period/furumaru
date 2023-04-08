<script lang="ts" setup>
import { CreateProducerRequest } from '~/types/api'
import { ImageUploadStatus } from '~/types/props'

const props = defineProps({
  formType: {
    type: String,
    default: 'create',
    validator: (value: string) => {
      return ['create', 'edit'].includes(value)
    }
  },
  formData: {
    type: Object,
    default: (): CreateProducerRequest => ({
      lastname: '',
      lastnameKana: '',
      firstname: '',
      firstnameKana: '',
      storeName: '',
      thumbnailUrl: '',
      headerUrl: '',
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

const emit = defineEmits<{
  (e: 'update:formData', formData: CreateProducerRequest): void
  (e: 'update:thumbnailFile', files?: FileList): void
  (e: 'update:headerFile', files?: FileList): void
  (e: 'click:search'): void
  (e: 'submit'): void
}>()

const formDataValue = computed({
  get: (): CreateProducerRequest => props.formData as CreateProducerRequest,
  set: (val: CreateProducerRequest) => emit('update:formData', val)
})

const btnText = computed(() => {
  return props.formType === 'create' ? '登録' : '更新'
})

const updateThumbnailFileHandler = (files?: FileList) => {
  emit('update:thumbnailFile', files)
}

const updateHeaderFileHandler = (files?: FileList) => {
  emit('update:headerFile', files)
}

const handleSubmit = () => {
  emit('submit')
}

const handleSearchClick = () => {
  emit('click:search')
}
</script>

<template>
  <form @submit.prevent="handleSubmit">
    <v-card elevation="0">
      <v-card-text>
        <v-text-field
          v-model="formDataValue.storeName"
          label="店舗名"
          required
          maxlength="64"
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
            v-model="formDataValue.lastname"
            class="mr-4"
            label="生産者名:姓"
            maxlength="16"
            required
          />
          <v-text-field
            v-model="formDataValue.firstname"
            label="生産者名:名"
            maxlength="16"
            required
          />
        </div>
        <div class="d-flex">
          <v-text-field
            v-model="formDataValue.lastnameKana"
            class="mr-4"
            label="生産者名:姓（ふりがな）"
            maxlength="32"
            required
          />
          <v-text-field
            v-model="formDataValue.firstnameKana"
            label="生産者名:名（ふりがな）"
            maxlength="32"
            required
          />
        </div>
        <v-text-field
          v-model="formDataValue.email"
          label="連絡先（Email）"
          type="email"
          required
        />
        <v-text-field
          v-model="formDataValue.phoneNumber"
          label="連絡先（電話番号）"
          required
        />

        <the-address-form
          v-model:postal-code="props.formData.postalCode"
          v-model:prefecture="props.formData.prefecture"
          v-model:city="props.formData.city"
          v-model:address-line1="props.formData.addressLine1"
          v-model:address-line2="props.formData.addressLine2"
          :loading="props.searchLoading"
          :error-message="props.searchErrorMessage"
          @click:search="handleSearchClick"
        />
      </v-card-text>
      <v-card-actions>
        <v-btn block outlined color="primary" type="submit">
          {{ btnText }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </form>
</template>
