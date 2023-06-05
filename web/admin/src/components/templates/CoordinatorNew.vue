<script lang="ts" setup>
import { AlertType } from '~/lib/hooks'
import { CreateCoordinatorRequest } from '~/types/api'
import { ImageUploadStatus } from '~/types/props'

const props = defineProps({
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
    default: () => ({})
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
  (e: 'update:form-data', v: CreateCoordinatorRequest): void
  (e: 'update:thumbnail-file', files?: FileList): void
  (e: 'update:header-file', files?: FileList): void
  (e: 'click:search-address'): void
  (e: 'submit'): void
}>()

const formDataValue = computed({
  get: (): CreateCoordinatorRequest => props.formData,
  set: (v: CreateCoordinatorRequest): void => emit('update:form-data', v)
})

const updateThumbnailFileHandler = (files?: FileList) => {
  emit('update:thumbnail-file', files)
}

const updateHeaderFileHandler = (files?: FileList) => {
  emit('update:header-file', files)
}

const onSubmit = (): void => {
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
