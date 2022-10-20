<template>
  <div>
    <v-card-title>コーディネーター登録</v-card-title>
    <v-card elevation="0">
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
            label="コーディネーター:姓"
          />
          <v-text-field
            v-model="v$.firstname.$model"
            :error-messages="getErrorMessage(v$.firstname.$errors)"
            label="コーディネーター:名"
          />
        </div>
        <div class="d-flex">
          <v-text-field
            v-model="v$.lastnameKana.$model"
            :error-messages="getErrorMessage(v$.lastnameKana.$errors)"
            class="mr-4"
            label="コーディネーター:姓（ふりがな）"
          />
          <v-text-field
            v-model="v$.firstnameKana.$model"
            :error-messages="getErrorMessage(v$.firstnameKana.$errors)"
            label="コーディネーター:名（ふりがな）"
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
        />
      </v-card-text>
      <v-card-actions>
        <v-btn block outlined color="primary" @click="handleSubmit">登録</v-btn>
      </v-card-actions>
    </v-card>
  </div>
</template>

<script lang="ts">
import { reactive, useRouter } from '@nuxtjs/composition-api'
import { computed, defineComponent } from '@vue/composition-api'
import { useVuelidate } from '@vuelidate/core'

import { useSearchAddress } from '~/lib/hooks'
import {
  kana,
  getErrorMessage,
  required,
  email,
  tel,
  maxLength,
} from '~/lib/validations'
import { useCoordinatorStore } from '~/store/coordinator'
import { CreateCoordinatorRequest } from '~/types/api'
import { ImageUploadStatus } from '~/types/props'

export default defineComponent({
  setup() {
    const router = useRouter()
    const { createCoordinator, uploadCoordinatorThumbnail } =
      useCoordinatorStore()

    const formData = reactive<CreateCoordinatorRequest>({
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
    })

    const thumbnailUploadStatus = reactive<ImageUploadStatus>({
      error: false,
      message: '',
    })

    const headerUploadStatus = reactive<ImageUploadStatus>({
      error: false,
      message: '',
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

    const v$ = useVuelidate(rules, formData)

    const handleSubmit = async () => {
      try {
        await createCoordinator({
          ...formData,
          phoneNumber: formData.phoneNumber.replace('0', '+81'),
        })
        router.push('/coordinators')
      } catch (error) {
        console.log(error)
      }
    }

    const handleUpdateThumbnail = (files: FileList) => {
      if (files.length > 0) {
        uploadCoordinatorThumbnail(files[0])
          .then((res) => {
            formData.thumbnailUrl = res.url
          })
          .catch(() => {
            thumbnailUploadStatus.error = true
            thumbnailUploadStatus.message = 'アップロードに失敗しました。'
          })
      }
    }

    const handleUpdateHeader = async (_files: FileList) => {}

    const {
      loading: searchLoading,
      errorMessage: searchErrorMessage,
      searchAddressByPostalCode,
    } = useSearchAddress()

    const searchAddress = async () => {
      searchLoading.value = true
      searchErrorMessage.value = ''
      const res = await searchAddressByPostalCode(Number(formData.postalCode))
      if (res) {
        formData.prefecture = res.prefecture
        formData.city = res.city
        formData.addressLine1 = res.addressLine1
      }
    }

    return {
      formData,
      handleSubmit,
      handleUpdateThumbnail,
      handleUpdateHeader,
      thumbnailUploadStatus,
      headerUploadStatus,
      searchAddress,
      searchLoading,
      searchErrorMessage,
      getErrorMessage,
      v$,
    }
  },
})
</script>
