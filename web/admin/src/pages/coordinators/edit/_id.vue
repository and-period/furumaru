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
        <v-textarea v-model="formNote" name="note" label="メモ"></v-textarea>
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
            label="コーディネーター:姓（ふりがな）"
          />
          <v-text-field
            v-model="v$.firstnameKana.$model"
            :error-messages="getErrorMessage(v$.firstnameKana.$errors)"
            class="mr-4"
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
import {
  computed,
  defineComponent,
  reactive,
  useFetch,
  useRoute,
  useRouter,
} from '@nuxtjs/composition-api'
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
import { CoordinatorResponse } from '~/types/api'
import { ImageUploadStatus } from '~/types/props'

export default defineComponent({
  setup() {
    const coordinatorStore = useCoordinatorStore()

    const route = useRoute()
    const id = route.value.params.id
    const router = useRouter()

    const { uploadCoordinatorThumbnail, uploadCoordinatorHeader } =
      useCoordinatorStore()

    const { getCoordinator } = useCoordinatorStore()

    const formData = reactive<CoordinatorResponse>({
      id,
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
      createdAt: 0,
      updatedAt: 0,
    })

    const { fetchState } = useFetch(async () => {
      const coordinator = await getCoordinator(id)
      formData.storeName = coordinator.storeName
      formData.firstname = coordinator.firstname
      formData.lastname = coordinator.lastname
      formData.firstnameKana = coordinator.firstnameKana
      formData.lastnameKana = coordinator.lastnameKana
      formData.companyName = coordinator.companyName
      formData.thumbnailUrl = coordinator.thumbnailUrl
      formData.headerUrl = coordinator.headerUrl
      formData.twitterAccount = coordinator.twitterAccount
      formData.instagramAccount = coordinator.instagramAccount
      formData.facebookAccount = coordinator.facebookAccount
      formData.email = coordinator.email
      formData.phoneNumber = coordinator.phoneNumber.replace('+81', '0')
      formData.postalCode = coordinator.postalCode
      formData.prefecture = coordinator.prefecture
      formData.city = coordinator.city
      formData.addressLine1 = coordinator.addressLine1
      formData.addressLine2 = coordinator.addressLine2
      formData.createdAt = coordinator.createdAt
      formData.updatedAt = coordinator.updatedAt
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

    const thumbnailUploadStatus = reactive<ImageUploadStatus>({
      error: false,
      message: '',
    })

    const headerUploadStatus = reactive<ImageUploadStatus>({
      error: false,
      message: '',
    })

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

    const handleUpdateHeader = async (files: FileList) => {
      if (files.length > 0) {
        await uploadCoordinatorHeader(files[0])
          .then((res) => {
            formData.headerUrl = res.url
          })
          .catch(() => {
            headerUploadStatus.error = true
            headerUploadStatus.message = 'アップロードに失敗しました。'
          })
      }
    }

    const handleSubmit = async (): Promise<void> => {
      try {
        const result = await v$.value.$validate()
        if (!result) {
          return
        }
        await coordinatorStore.updateCoordinator(
          {
            ...formData,
            phoneNumber: formData.phoneNumber.replace('0', '+81'),
          },
          id
        )
        router.push('/')
      } catch (error) {
        console.log(error)
      }
    }

    return {
      id,
      fetchState,
      formData,
      v$,
      getErrorMessage,
      searchLoading,
      searchErrorMessage,
      searchAddress,
      handleSubmit,
      handleUpdateThumbnail,
      thumbnailUploadStatus,
      headerUploadStatus,
      handleUpdateHeader,
    }
  },
})
</script>
