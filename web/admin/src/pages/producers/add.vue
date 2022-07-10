<template>
  <form @submit.prevent="handleSubmit">
    <p class="text-h6">生産者登録</p>
    <v-card elevation="0">
      <v-card-text>
        <v-text-field label="店舗名" required maxlength="64" />
        <div class="mb-2">
          <the-file-upload-filed text="生産者画像" />
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
          type="number"
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
import { defineComponent, reactive } from '@nuxtjs/composition-api'

import { useProducerStore } from '~/store/producer'
import { CreateProducerRequest } from '~/types/api'

export default defineComponent({
  setup() {
    const { createProducer } = useProducerStore()

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

    const handleSubmit = async () => {
      try {
        await createProducer(formData)
      } catch (error) {
        console.log(error)
      }
    }

    return { formData, handleSubmit }
  },
})
</script>
