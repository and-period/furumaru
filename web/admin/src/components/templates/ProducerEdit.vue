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
  },
  formDataLoading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits<{
  (e: 'update:formData', formData: CreateProducerRequest): void
  (e: 'update:thumbnail-file', files?: FileList): void
  (e: 'update:header-file', files?: FileList): void
  (e: 'click:search'): void
  (e: 'submit'): void
}>()

const formDataValue = computed({
  get: (): CreateProducerRequest => props.formData as CreateProducerRequest,
  set: (val: CreateProducerRequest) => emit('update:formData', val)
})

const updateThumbnailFileHandler = (files?: FileList) => {
  emit('update:thumbnail-file', files)
}

const updateHeaderFileHandler = (files?: FileList) => {
  emit('update:header-file', files)
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
      生産者編集
    </p>
    <v-skeleton-loader v-if="props.formDataLoading" type="article" />
    <organisms-producer-form
      v-else
      form-type="edit"
      :form-data="formDataValue"
      :thumbnail-upload-status="props.thumbnailUploadStatus"
      :header-upload-status="props.headerUploadStatus"
      :search-loading="props.searchLoading"
      :search-error-message="props.searchErrorMessage"
      :form-data-loading="props.formDataLoading"
      @update:thumbnail-file="updateThumbnailFileHandler"
      @update:header-file="updateHeaderFileHandler"
      @submit="handleSubmit"
      @click:search="handleSearchClick"
    />
  </div>
</template>
