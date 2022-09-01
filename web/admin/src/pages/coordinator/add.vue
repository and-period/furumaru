<template>
  <div>
    <v-card-title>コーディネーター登録</v-card-title>
    <v-card elevation="0">
      <v-card-text>
        <v-text-field v-model="formData.storeName" label="店舗名" />
        <div class="mb-2">
          <the-file-upload-filed text="コーディネーター画像" />
        </div>
        <v-text-field v-model="formData.companyName" label="会社名" />
        <div class="d-flex">
          <v-text-field
            v-model="formData.lastname"
            class="mr-4"
            label="コーディネーター:姓"
          />
          <v-text-field
            v-model="formData.firstname"
            label="コーディネーター:名"
          />
        </div>
        <div class="d-flex">
          <v-text-field
            v-model="formData.lastnameKana"
            class="mr-4"
            label="コーディネーター:姓（ふりがな）"
          />
          <v-text-field
            v-model="formData.firstnameKana"
            label="コーディネーター:名（ふりがな）"
          />
        </div>
        <v-text-field v-model="formData.email" label="連絡先（Email）" />
        <v-text-field
          v-model="formData.phoneNumber"
          type="tel"
          label="連絡先（電話番号）"
        />

        <the-address-form
          :postal-code.sync="formData.postalCode"
          :prefecture.sync="formData.prefecture"
          :city.sync="formData.city"
          :address-line1.sync="formData.addressLine1"
          :address-line2.sync="formData.addressLine2"
          :loading="searchLoading"
          :error-message="searchErrorMessage"
          @click:search="searchAddress"
        />
      </v-card-text>
      <v-card-actions>
        <v-btn block outlined color="primary" @click="handleSubmit">登録</v-btn>
      </v-card-actions>
    </v-card>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive } from '@nuxtjs/composition-api'

import { useSearchAddress } from '~/lib/hooks'
import { useCoordinatorStore } from '~/store/coordinator'
import { CreateCoordinatorRequest } from '~/types/api'

export default defineComponent({
  setup() {
    const formData = reactive<CreateCoordinatorRequest>({
      storeName: '',
      firstname: '',
      lastname: '',
      firstnameKana: '',
      lastnameKana: '',
      companyName: '',
      thumbnailUrl: '',
      headerUrl: '',
      twitterAccount: '',
      instagramAccount: '',
      facebookAccount: '',
      email: '',
      phoneNumber: '',
      postalCode: '',
      prefecture: '',
      city: '',
      addressLine1: '',
      addressLine2: '',
    })

    const { createCoordinator } = useCoordinatorStore()

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

    const handleSubmit = async () => {
      await createCoordinator({
        ...formData,
        phoneNumber: formData.phoneNumber.replace('0', '+81'),
      })
    }

    return {
      formData,
      searchLoading,
      searchErrorMessage,
      searchAddress,
      handleSubmit,
    }
  },
})
</script>
