<template>
  <div>
    <v-card-title>コーディネーター登録</v-card-title>
    <v-card>
      <the-coordinator-create-form
        :form-data="formData"
        :thumbnail-upload-status="thumbnailUploadStatus"
        :header-upload-status="headerUploadStatus"
        :search-loading="searchLoading"
        :search-error-message="searchErrorMessage"
        @update:thumbnailFile="handleUpdateThumbnail"
        @update:headerFile="handleUpdateHeader"
        @submit="handleSubmit"
        @click:search="searchAddress"
      />
    </v-card>
  </div>
</template>

<script lang="ts">
import { reactive, useRouter, defineComponent } from '@nuxtjs/composition-api'

import { useSearchAddress } from '~/lib/hooks'
import { useCoordinatorStore } from '~/store/coordinator'
import { CreateCoordinatorRequest } from '~/types/api'
import { ImageUploadStatus } from '~/types/props'

export default defineComponent({
  setup() {
    const router = useRouter()
    const {
      createCoordinator,
      uploadCoordinatorThumbnail,
      uploadCoordinatorHeader,
    } = useCoordinatorStore()

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

    const handleUpdateHeader = (files: FileList) => {
      if (files.length > 0) {
        uploadCoordinatorHeader(files[0])
          .then((res) => {
            formData.headerUrl = res.url
          })
          .catch(() => {
            headerUploadStatus.error = true
            headerUploadStatus.message = 'アップロードに失敗しました。'
          })
      }
    }

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
    }
  },
})
</script>
