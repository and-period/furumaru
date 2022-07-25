<template>
  <the-producer-edit-form-page
    :form-data="formData"
    :form-data-loading="fetchState.pending"
    :thumbnail-upload-status="thumbnailUploadStatus"
    :header-upload-status="headerUploadStatus"
    :search-loading="searchLoading"
    :search-error-message="searchErrorMessage"
    @update:thumbnailFile="handleUpdateThumbnail"
    @update:headerFile="handleUpdateHeader"
    @submit="handleSubmit"
    @click:search="searchAddress"
  />
</template>

<script lang="ts">
import {
  defineComponent,
  reactive,
  useFetch,
  useRoute,
} from '@nuxtjs/composition-api'

import { useSearchAddress } from '~/lib/hooks'
import { useProducerStore } from '~/store/producer'
import { ProducerResponse } from '~/types/api'
import { ImageUploadStatus } from '~/types/props'

export default defineComponent({
  setup() {
    const route = useRoute()
    const id = route.value.params.id

    const { getProducer } = useProducerStore()

    const { uploadProducerThumbnail, uploadProducerHeader } = useProducerStore()

    const thumbnailUploadStatus = reactive<ImageUploadStatus>({
      error: false,
      message: '',
    })

    const headerUploadStatus = reactive<ImageUploadStatus>({
      error: false,
      message: '',
    })

    const formData = reactive<ProducerResponse>({
      id,
      lastname: '',
      lastnameKana: '',
      firstname: '',
      firstnameKana: '',
      addressLine1: '',
      addressLine2: '',
      city: '',
      prefecture: '',
      phoneNumber: '',
      postalCode: '',
      storeName: '',
      headerUrl: '',
      createdAt: -1,
      updatedAt: -1,
      thumbnailUrl: '',
      email: '',
    })

    const { fetchState } = useFetch(async () => {
      const producer = await getProducer(id)
      formData.lastname = producer.lastname
      formData.lastnameKana = producer.lastnameKana
      formData.firstname = producer.firstname
      formData.firstnameKana = producer.firstnameKana
      formData.addressLine1 = producer.addressLine1
      formData.addressLine2 = producer.addressLine2
      formData.city = producer.city
      formData.prefecture = producer.prefecture
      formData.phoneNumber = producer.phoneNumber
      formData.postalCode = producer.postalCode
      formData.storeName = producer.storeName
      formData.headerUrl = producer.headerUrl
      formData.thumbnailUrl = producer.thumbnailUrl
      formData.email = producer.email
      formData.createdAt = producer.createdAt
      formData.updatedAt = producer.updatedAt
    })

    const handleUpdateThumbnail = (files: FileList) => {
      if (files.length > 0) {
        uploadProducerThumbnail(files[0])
          .then((res) => {
            formData.thumbnailUrl = res.url
          })
          .catch(() => {
            thumbnailUploadStatus.error = true
            thumbnailUploadStatus.message = 'アップロードに失敗しました。'
          })
      }
    }

    const handleUpdateHeader = async (files: FileList) => {
      if (files.length > 0) {
        await uploadProducerHeader(files[0])
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

    const handleSubmit = () => {
      console.log('未実装')
    }

    return {
      id,
      fetchState,
      formData,
      handleUpdateThumbnail,
      handleUpdateHeader,
      thumbnailUploadStatus,
      headerUploadStatus,
      searchAddress,
      searchLoading,
      searchErrorMessage,
      handleSubmit,
    }
  },
})
</script>
