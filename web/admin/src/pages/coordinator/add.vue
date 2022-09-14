<template>
  <div>
    <v-card-title>コーディネーター登録</v-card-title>
    <v-card elevation="0">
      <v-card-text>
        <v-text-field
          v-model="v$.storeName.$model"
          :error-messages="getErrorMessage(v$.storeName.$errors)"
          label="店舗名"
        />
        <div class="mb-2">
          <the-file-upload-filed text="コーディネーター画像" />
        </div>
        <v-text-field
          v-model="v$.companyName.$model"
          :error-messages="getErrorMessage(v$.companyName.$errors)"
          label="会社名"
        />
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
import { computed, defineComponent, reactive } from '@nuxtjs/composition-api'
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
      const result = await v$.value.$validate()
      if (!result) {
        return
      }
      await createCoordinator({
        ...formData,
        phoneNumber: formData.phoneNumber.replace('0', '+81'),
      })
    }

    return {
      formData,
      v$,
      getErrorMessage,
      searchLoading,
      searchErrorMessage,
      searchAddress,
      handleSubmit,
    }
  },
})
</script>
