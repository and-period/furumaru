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
            @update:file="handleUpdateThumbnail"
          />
          <the-header-select-form
            class="flex-grow-1 flex-shrink-1"
            :img-url="formData.headerUrl"
            :error="headerUploadStatus.error"
            :message="headerUploadStatus.message"
            @update:file="handleUpdateHeader"
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
        <v-text-field v-model="formData.postalCode" label="郵便番号" />
        <v-text-field v-model="formData.prefecture" label="都道府県" />
        <v-text-field v-model="formData.city" label="市区町村" />
        <v-text-field v-model="formData.addressLine1" label="番地" />
        <v-text-field
          v-model="formData.addressLine2"
          label="アパート名・部屋番号（任意）"
        />
      </v-card-text>
      <v-card-actions>
        <v-btn block outlined color="primary" type="submit">登録</v-btn>
      </v-card-actions>
    </v-card>
  </form>
</template>

<script lang="ts">
import { defineComponent, reactive, useRouter } from '@nuxtjs/composition-api'

import TheProfileSelectForm from '~/components/molecules/TheProfileSelectForm.vue'
import { useProducerStore } from '~/store/producer'
import { CreateProducerRequest } from '~/types/api'

interface ImageUploadStatus {
  error: boolean
  message: string
}

export default defineComponent({
  components: { TheProfileSelectForm },
  setup() {
    const router = useRouter()
    const { createProducer, uploadProducerThumbnail, uploadProducerHeader } =
      useProducerStore()

    const formData = reactive<CreateProducerRequest>({
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
        await createProducer({
          ...formData,
          phoneNumber: formData.phoneNumber.replace('0', '+81'),
        })
        router.push('/producers')
      } catch (error) {
        console.log(error)
      }
    }

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

    return {
      formData,
      handleSubmit,
      handleUpdateThumbnail,
      handleUpdateHeader,
      thumbnailUploadStatus,
      headerUploadStatus,
    }
  },
})
</script>
