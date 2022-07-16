<template>
  <form @submit.prevent="handleSubmit">
    <p class="text-h6">生産者登録</p>
    <v-card elevation="0">
      <v-card-text>
        <v-text-field
          v-model="formData.storeName"
          label="店舗名"
          required
          maxlength="64"
        />
        <div class="mb-2 d-flex">
          <the-profile-select-form
            class="mr-4 flex-grow-1 flex-shrink-1"
            :img-url="formData.thumbnailUrl"
            :error="thumbnailUploadStatus.error"
            :message="thumbnailUploadStatus.message"
            @update:file="updateThumbnailFileHandler"
          />
          <the-header-select-form
            class="flex-grow-1 flex-shrink-1"
            :img-url="formData.headerUrl"
            :error="headerUploadStatus.error"
            :message="headerUploadStatus.message"
            @update:file="updateHeaderFileHandler"
          />
        </div>
        <div class="d-flex">
          <v-text-field
            v-model="formData.lastname"
            class="mr-4"
            label="生産者名:姓"
            maxlength="16"
            required
          />
          <v-text-field
            v-model="formData.firstname"
            label="生産者名:名"
            maxlength="16"
            required
          />
        </div>
        <div class="d-flex">
          <v-text-field
            v-model="formData.lastnameKana"
            class="mr-4"
            label="生産者名:姓（ふりがな）"
            maxlength="32"
            required
          />
          <v-text-field
            v-model="formData.firstnameKana"
            label="生産者名:名（ふりがな）"
            maxlength="32"
            required
          />
        </div>
        <v-text-field
          v-model="formData.email"
          label="連絡先（Email）"
          type="email"
          required
        />
        <v-text-field
          v-model="formData.phoneNumber"
          label="連絡先（電話番号）"
          required
        />

        <the-address-form
          :postal-code.sync="formData.postalCode"
          :prefecture.sync="formData.prefecture"
          :city.sync="formData.city"
          :address-line1.sync="formData.addressLine1"
          :address-line2.sync="formData.addressLine2"
          :loading="searchLoading"
          :error-message="searchErrorMessage"
          @click:search="handleSearchClick"
        />
      </v-card-text>
      <v-card-actions>
        <v-btn block outlined color="primary" type="submit">登録</v-btn>
      </v-card-actions>
    </v-card>
  </form>
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
