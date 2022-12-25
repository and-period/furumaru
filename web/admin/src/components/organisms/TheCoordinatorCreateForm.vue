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
          v-model="v$.email.$model"
          label="連絡先（Email）"
          :error-messages="getErrorMessage(v$.email.$errors)"
        />
        <v-text-field
          v-model="v$.phoneNumber.$model"
          :error-messages="getErrorMessage(v$.phoneNumber.$errors)"
          type="tel"
          label="連絡先（電話番号）"
        />

        <the-address-form
          :postal-code.sync="formData.postalCode"
          :prefecture.sync="formData.prefecture"
          :city.sync="formData.city"
          :address-line1.sync="formData.addressLine1"
          :address-line2.sync="formData.addressLine2"
          :error-message="searchErrorMessage"
          :loading="searchLoading"
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
import { useVuelidate } from '@vuelidate/core'

import {
  kana,
  getErrorMessage,
  required,
  email,
  tel,
  maxLength,
} from '~/lib/validations'
import { CreateCoordinatorRequest } from '~/types/api'
import { ImageUploadStatus } from '~/types/props'

export default defineComponent({
  props: {
    formData: {
      type: Object as PropType<CreateCoordinatorRequest>,
      default: () => {
        return {
          lastname: '',
          lastnameKana: '',
          firstname: '',
          firstnameKana: '',
          companyName: '',
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
      get: (): CreateCoordinatorRequest => props.formData,
      set: (val: CreateCoordinatorRequest) => emit('update:formData', val),
    })

    const rules = computed(() => ({
      storeName: { required, maxLength: maxLength(64) },
      companyName: { required, maxLength: maxLength(64) },
      firstname: { required, maxLength: maxLength(16) },
      lastname: { required, maxLength: maxLength(16) },
      firstnameKana: { required, kana },
      lastnameKana: { required, kana },
      phoneNumber: { required, tel },
      email: { required, email },
    }))

    const v$ = useVuelidate(rules, formDataValue)

    const updateThumbnailFileHandler = (files: FileList) => {
      emit('update:thumbnailFile', files)
    }

    const updateHeaderFileHandler = (files: FileList) => {
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

    return {
      formDataValue,
      updateThumbnailFileHandler,
      updateHeaderFileHandler,
      handleSearchClick,
      handleSubmit,
      getErrorMessage,
      v$,
    }
  },
})
</script>
