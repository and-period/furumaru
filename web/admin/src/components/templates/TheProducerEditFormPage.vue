<template>
  <div>
    <p class="text-h6">生産者編集</p>
    <v-skeleton-loader v-if="formDataLoading" type="article" />
    <the-producer-form
      v-else
      form-type="edit"
      :form-data="formDataValue"
      :thumbnail-upload-status="thumbnailUploadStatus"
      :header-upload-status="headerUploadStatus"
      :search-loading="searchLoading"
      :search-error-message="searchErrorMessage"
      :form-data-loading="formDataLoading"
      @update:thumbnailFile="updateThumbnailFileHandler"
      @update:headerFile="updateHeaderFileHandler"
      @submit="handleSubmit"
      @click:search="handleSearchClick"
    />
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, PropType } from '@vue/composition-api'

import { CreateProducerRequest } from '~/types/api'
import { ImageUploadStatus } from '~/types/props'

export default defineComponent({
  props: {
    formData: {
      type: Object as PropType<CreateProducerRequest>,
      default: () => {
        return {
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
          addressLine2: '',
        }
      },
    },
    thumbnailUploadStatus: {
      type: Object as PropType<ImageUploadStatus>,
      default: () => {
        return {
          error: false,
          message: '',
        }
      },
    },
    headerUploadStatus: {
      type: Object as PropType<ImageUploadStatus>,
      default: () => {
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
    formDataLoading: {
      type: Boolean,
      default: false,
    },
  },

  setup(props, { emit }) {
    const formDataValue = computed({
      get: (): CreateProducerRequest => props.formData,
      set: (val: CreateProducerRequest) => emit('update:formData', val),
    })

    const updateThumbnailFileHandler = (files: FileList) => {
      emit('update:thumbnailFile', files)
    }

    const updateHeaderFileHandler = (files: FileList) => {
      emit('update:headerFile', files)
    }

    const handleSubmit = () => {
      emit('submit')
    }

    const handleSearchClick = () => {
      emit('click:search')
    }

    return {
      formDataValue,
      updateThumbnailFileHandler,
      updateHeaderFileHandler,
      handleSubmit,
      handleSearchClick,
    }
  },
})
</script>
