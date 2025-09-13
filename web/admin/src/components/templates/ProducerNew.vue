<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import { mdiFacebook, mdiInstagram } from '@mdi/js'
import type { AlertType } from '~/lib/hooks'
import { Prefecture } from '~/types'
import type { CreateProducerRequest } from '~/types/api/v1'
import { getErrorMessage } from '~/lib/validations'
import type { ImageUploadStatus } from '~/types/props'
import { CreateProducerValidationRules } from '~/types/validations'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  isAlert: {
    type: Boolean,
    default: false,
  },
  alertType: {
    type: String as PropType<AlertType>,
    default: undefined,
  },
  alertText: {
    type: String,
    default: '',
  },
  formData: {
    type: Object as PropType<CreateProducerRequest>,
    default: (): CreateProducerRequest => ({
      coordinatorId: '',
      lastname: '',
      lastnameKana: '',
      firstname: '',
      firstnameKana: '',
      username: '',
      email: '',
      phoneNumber: '',
      postalCode: '',
      prefectureCode: Prefecture.UNKNOWN,
      city: '',
      addressLine1: '',
      addressLine2: '',
      profile: '',
      thumbnailUrl: '',
      headerUrl: '',
      promotionVideoUrl: '',
      bonusVideoUrl: '',
      instagramId: '',
      facebookId: '',
    }),
  },
  thumbnailUploadStatus: {
    type: Object,
    default: (): ImageUploadStatus => ({
      error: false,
      message: '',
    }),
  },
  headerUploadStatus: {
    type: Object,
    default: (): ImageUploadStatus => ({
      error: false,
      message: '',
    }),
  },
  promotionVideoUploadStatus: {
    type: Object,
    default: (): ImageUploadStatus => ({
      error: false,
      message: '',
    }),
  },
  bonusVideoUploadStatus: {
    type: Object,
    default: (): ImageUploadStatus => ({
      error: false,
      message: '',
    }),
  },
  searchErrorMessage: {
    type: String,
    default: '',
  },
  searchLoading: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits<{
  (e: 'update:form-data', formData: CreateProducerRequest): void
  (e: 'update:thumbnail-file', files: FileList): void
  (e: 'update:header-file', files: FileList): void
  (e: 'update:promotion-video', files: FileList): void
  (e: 'update:bonus-video', files: FileList): void
  (e: 'click:search-address'): void
  (e: 'submit'): void
}>()

const formDataValue = computed({
  get: (): CreateProducerRequest => props.formData,
  set: (formData: CreateProducerRequest) => emit('update:form-data', formData),
})

const validate = useVuelidate(CreateProducerValidationRules, formDataValue)

const onChangeThumbnailFile = (files?: FileList) => {
  if (!files) {
    return
  }
  emit('update:thumbnail-file', files)
}

const onChangeHeaderFile = (files?: FileList) => {
  if (!files) {
    return
  }
  emit('update:header-file', files)
}

const onChangePromotionVideo = (files?: FileList) => {
  if (!files) {
    return
  }
  emit('update:promotion-video', files)
}

const onChangeBonusVideo = (files?: FileList) => {
  if (!files) {
    return
  }
  emit('update:bonus-video', files)
}

const onSubmit = async (): Promise<void> => {
  const valid = await validate.value.$validate()
  if (!valid) {
    return
  }

  emit('submit')
}

const onClickSearchAddress = (): void => {
  emit('click:search-address')
}
</script>

<template>
  <v-alert
    v-show="props.isAlert"
    :type="props.alertType"
    v-text="props.alertText"
  />

  <v-card>
    <v-card-title>生産者登録</v-card-title>

    <v-form @submit.prevent="onSubmit">
      <v-card-text>
        <v-text-field
          v-model="validate.username.$model"
          :error-messages="getErrorMessage(validate.username.$errors)"
          label="生産者名"
        />
        <v-row>
          <v-col
            cols="12"
            ms="12"
            lg="6"
          >
            <molecules-video-select-form
              label="紹介動画"
              :loading="loading"
              :video-url="formDataValue.promotionVideoUrl"
              :error="props.promotionVideoUploadStatus.error"
              :message="props.promotionVideoUploadStatus.message"
              @update:file="onChangePromotionVideo"
            />
          </v-col>
          <v-col
            cols="12"
            sm="12"
            lg="6"
          >
            <molecules-video-select-form
              label="サンキュー動画"
              :loading="loading"
              :video-url="formDataValue.bonusVideoUrl"
              :error="props.bonusVideoUploadStatus.error"
              :message="props.bonusVideoUploadStatus.message"
              @update:file="onChangeBonusVideo"
            />
          </v-col>
        </v-row>
        <v-textarea
          v-model="validate.profile.$model"
          :error-messages="getErrorMessage(validate.profile.$errors)"
          label="プロフィール"
          maxlength="2000"
        />
        <v-row>
          <v-col>
            <v-text-field
              v-model="validate.lastname.$model"
              :error-messages="getErrorMessage(validate.lastname.$errors)"
              class="mr-4"
              label="生産者名:姓"
            />
          </v-col>
          <v-col>
            <v-text-field
              v-model="validate.firstname.$model"
              :error-messages="getErrorMessage(validate.firstname.$errors)"
              label="生産者名:名"
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-text-field
              v-model="validate.lastnameKana.$model"
              :error-messages="getErrorMessage(validate.lastnameKana.$errors)"
              class="mr-4"
              label="生産者名:姓（ふりがな）"
            />
          </v-col>
          <v-col>
            <v-text-field
              v-model="validate.firstnameKana.$model"
              :error-messages="
                getErrorMessage(validate.firstnameKana.$errors)
              "
              label="生産者名:名（ふりがな）"
            />
          </v-col>
        </v-row>
        <v-text-field
          v-model="validate.email.$model"
          :error-messages="getErrorMessage(validate.email.$errors)"
          label="連絡先（Email）"
          type="email"
        />
        <v-text-field
          v-model="validate.phoneNumber.$model"
          :error-messages="getErrorMessage(validate.phoneNumber.$errors)"
          label="連絡先（電話番号）"
        />
        <v-row>
          <v-col
            cols="12"
            sm="6"
            md="6"
          >
            <molecules-icon-select-form
              label="アイコン画像"
              :loading="loading"
              :img-url="formDataValue.thumbnailUrl"
              :error="props.thumbnailUploadStatus.error"
              :message="props.thumbnailUploadStatus.message"
              @update:file="onChangeThumbnailFile"
            />
          </v-col>
          <v-col
            cols="12"
            sm="6"
            md="6"
          >
            <molecules-image-select-form
              label="ヘッダー画像"
              :loading="loading"
              :img-url="formDataValue.headerUrl"
              :error="props.headerUploadStatus.error"
              :message="props.headerUploadStatus.message"
              @update:file="onChangeHeaderFile"
            />
          </v-col>
        </v-row>
        <molecules-address-form
          v-model:postal-code="formDataValue.postalCode"
          v-model:prefecture="formDataValue.prefectureCode"
          v-model:city="formDataValue.city"
          v-model:address-line1="formDataValue.addressLine1"
          v-model:address-line2="formDataValue.addressLine2"
          :error-messages="props.searchErrorMessage"
          :loading="props.searchLoading"
          @click:search="onClickSearchAddress"
        />
        <p>SNS連携</p>
        <v-text-field
          v-model="validate.instagramId.$model"
          :error-messages="getErrorMessage(validate.instagramId.$errors)"
          :prepend-icon="mdiInstagram"
          prefix="https://www.instagram.com/"
        />
        <v-text-field
          v-model="validate.facebookId.$model"
          :error-messages="getErrorMessage(validate.facebookId.$errors)"
          :prepend-icon="mdiFacebook"
          prefix="https://www.facebook.com/"
        />
      </v-card-text>

      <v-card-actions>
        <v-btn
          block
          :loading="loading"
          variant="outlined"
          color="primary"
          type="submit"
        >
          登録
        </v-btn>
      </v-card-actions>
    </v-form>
  </v-card>
</template>
