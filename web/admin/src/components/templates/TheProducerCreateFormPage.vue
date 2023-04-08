<script lang="ts" setup>
import { CreateProducerRequest } from '~/types/api'
import { ImageUploadStatus } from '~/types/props'

const props = defineProps({
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
  (e: 'submit'): void
  (e: 'click:search'): void
}>()

const formDataValue = computed({
  get: (): CreateProducerRequest => props.formData as CreateProducerRequest,
  set: (val: CreateProducerRequest) => emit('update:formData', val)
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
  <div>
    <p class="text-h6">
      生産者登録
    </p>
    <the-producer-form
      :form-data="formDataValue"
      :thumbnail-upload-status="props.thumbnailUploadStatus"
      :header-upload-status="props.headerUploadStatus"
      :search-loading="props.searchLoading"
      :search-error-message="props.searchErrorMessage"
      @update:thumbnailFile="updateThumbnailFileHandler"
      @update:headerFile="updateHeaderFileHandler"
      @submit="handleSubmit"
      @click:search="handleSearchClick"
    />
  </div>
</template>
