<template>
  <div>
    <v-card-title>コーディネータ編集</v-card-title>

    <v-skeleton-loader v-if="fetchState.pending" type="article" />

    <the-coordinator-edit-form
      v-else
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

import TheCoordinatorEditForm from '~/components/organisms/TheCoordinatorEditForm.vue'
import { useSearchAddress } from '~/lib/hooks'
import {
  kana,
  getErrorMessage,
  required,
  tel,
  maxLength,
} from '~/lib/validations'
import { useCoordinatorStore } from '~/store/coordinator'
import { UpdateCoordinatorRequest } from '~/types/api'
import { ImageUploadStatus } from '~/types/props'

export default defineComponent({
  components: { TheCoordinatorEditForm },
  setup() {
    const coordinatorStore = useCoordinatorStore()

    const route = useRoute()
    const id = route.value.params.id
    const router = useRouter()

    const { uploadCoordinatorThumbnail, uploadCoordinatorHeader } =
      useCoordinatorStore()

    const { getCoordinator } = useCoordinatorStore()

    const formData = reactive<UpdateCoordinatorRequest>({
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
      phoneNumber: '',
      postalCode: '',
      prefecture: '',
      city: '',
      addressLine1: '',
      addressLine2: '',
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
      formData.phoneNumber = coordinator.phoneNumber.replace('+81', '0')
      formData.postalCode = coordinator.postalCode
      formData.prefecture = coordinator.prefecture
      formData.city = coordinator.city
      formData.addressLine1 = coordinator.addressLine1
      formData.addressLine2 = coordinator.addressLine2
    })

    const rules = computed(() => ({
      storeName: { required, maxLength: maxLength(64) },
      companyName: { required, maxLength: maxLength(64) },
      firstname: { required, maxLength: maxLength(16) },
      lastname: { required, maxLength: maxLength(16) },
      firstnameKana: { required, kana },
      lastnameKana: { required, kana },
      phoneNumber: { required, tel },
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
        router.push('/coordinators')
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
