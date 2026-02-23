<script lang="ts" setup>
import useVuelidate from '@vuelidate/core'
import {
  mdiFacebook,
  mdiInstagram,
  mdiAccount,
  mdiImage,
  mdiVideo,
  mdiEmail,
  mdiPhone,
  mdiMapMarker,
  mdiShareVariant,
  mdiArrowLeft,
  mdiPlus,
} from '@mdi/js'
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
  <v-container class="pa-6">
    <atoms-app-alert
      :show="props.isAlert"
      :type="props.alertType"
      :text="props.alertText"
      class="mb-6"
    />

    <div class="mb-6">
      <v-btn
        variant="text"
        :prepend-icon="mdiArrowLeft"
        @click="$router.back()"
      >
        戻る
      </v-btn>
      <h1 class="text-h4 font-weight-bold mt-2 mb-2">
        生産者登録
      </h1>
      <p class="text-body-1 text-grey-darken-1">
        新しい生産者の情報を登録します。必須項目をすべて入力してください。
      </p>
    </div>

    <v-form @submit.prevent="onSubmit">
      <!-- 基本情報セクション -->
      <v-card
        class="form-section-card mb-6"
        elevation="2"
      >
        <v-card-title class="d-flex align-center section-header">
          <v-icon
            :icon="mdiAccount"
            size="24"
            class="mr-3 text-primary"
          />
          <span class="text-h6 font-weight-medium">基本情報</span>
        </v-card-title>
        <v-card-text class="pa-6">
          <v-text-field
            v-model="validate.username.$model"
            :error-messages="getErrorMessage(validate.username.$errors)"
            label="生産者名 *"
            variant="outlined"
            density="comfortable"
            class="mb-4"
          />
          <v-textarea
            v-model="validate.profile.$model"
            :error-messages="getErrorMessage(validate.profile.$errors)"
            label="プロフィール"
            variant="outlined"
            density="comfortable"
            rows="4"
            maxlength="2000"
            counter
            class="mb-4"
          />
          <v-row>
            <v-col
              cols="12"
              sm="6"
            >
              <v-text-field
                v-model="validate.lastname.$model"
                :error-messages="getErrorMessage(validate.lastname.$errors)"
                label="姓 *"
                variant="outlined"
                density="comfortable"
              />
            </v-col>
            <v-col
              cols="12"
              sm="6"
            >
              <v-text-field
                v-model="validate.firstname.$model"
                :error-messages="getErrorMessage(validate.firstname.$errors)"
                label="名 *"
                variant="outlined"
                density="comfortable"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col
              cols="12"
              sm="6"
            >
              <v-text-field
                v-model="validate.lastnameKana.$model"
                :error-messages="getErrorMessage(validate.lastnameKana.$errors)"
                label="姓（ふりがな） *"
                variant="outlined"
                density="comfortable"
              />
            </v-col>
            <v-col
              cols="12"
              sm="6"
            >
              <v-text-field
                v-model="validate.firstnameKana.$model"
                :error-messages="
                  getErrorMessage(validate.firstnameKana.$errors)
                "
                label="名（ふりがな） *"
                variant="outlined"
                density="comfortable"
              />
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
      <!-- メディアセクション -->
      <v-card
        class="form-section-card mb-6"
        elevation="2"
      >
        <v-card-title class="d-flex align-center section-header">
          <v-icon
            :icon="mdiImage"
            size="24"
            class="mr-3 text-primary"
          />
          <span class="text-h6 font-weight-medium">画像・動画</span>
        </v-card-title>
        <v-card-text class="pa-6">
          <v-row>
            <v-col
              cols="12"
              sm="6"
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
          <v-row class="mt-4">
            <v-col
              cols="12"
              sm="6"
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
              sm="6"
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
        </v-card-text>
      </v-card>
      <!-- 連絡先セクション -->
      <v-card
        class="form-section-card mb-6"
        elevation="2"
      >
        <v-card-title class="d-flex align-center section-header">
          <v-icon
            :icon="mdiEmail"
            size="24"
            class="mr-3 text-primary"
          />
          <span class="text-h6 font-weight-medium">連絡先</span>
        </v-card-title>
        <v-card-text class="pa-6">
          <v-text-field
            v-model="validate.email.$model"
            :error-messages="getErrorMessage(validate.email.$errors)"
            label="メールアドレス *"
            type="email"
            variant="outlined"
            density="comfortable"
            class="mb-4"
          />
          <v-text-field
            v-model="validate.phoneNumber.$model"
            :error-messages="getErrorMessage(validate.phoneNumber.$errors)"
            label="電話番号 *"
            variant="outlined"
            density="comfortable"
            placeholder="090-1234-5678"
          />
        </v-card-text>
      </v-card>

      <!-- 住所セクション -->
      <v-card
        class="form-section-card mb-6"
        elevation="2"
      >
        <v-card-title class="d-flex align-center section-header">
          <v-icon
            :icon="mdiMapMarker"
            size="24"
            class="mr-3 text-primary"
          />
          <span class="text-h6 font-weight-medium">住所</span>
        </v-card-title>
        <v-card-text class="pa-6">
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
        </v-card-text>
      </v-card>
      <!-- SNS連携セクション -->
      <v-card
        class="form-section-card mb-6"
        elevation="2"
      >
        <v-card-title class="d-flex align-center section-header">
          <v-icon
            :icon="mdiShareVariant"
            size="24"
            class="mr-3 text-primary"
          />
          <span class="text-h6 font-weight-medium">SNS連携</span>
        </v-card-title>
        <v-card-text class="pa-6">
          <v-text-field
            v-model="validate.instagramId.$model"
            :error-messages="getErrorMessage(validate.instagramId.$errors)"
            :prepend-icon="mdiInstagram"
            prefix="https://www.instagram.com/"
            label="Instagram ID"
            variant="outlined"
            density="comfortable"
            class="mb-4"
          />
          <v-text-field
            v-model="validate.facebookId.$model"
            :error-messages="getErrorMessage(validate.facebookId.$errors)"
            :prepend-icon="mdiFacebook"
            prefix="https://www.facebook.com/"
            label="Facebook ID"
            variant="outlined"
            density="comfortable"
          />
        </v-card-text>
      </v-card>

      <v-footer
        app
        color="white"
        elevation="8"
        class="px-6 py-4 fixed-footer-actions"
      >
        <v-container
          fluid
          class="pa-0"
        >
          <div class="d-flex align-center justify-center flex-wrap ga-3">
            <v-btn
              variant="text"
              size="large"
              @click="$router.back()"
            >
              キャンセル
            </v-btn>
            <v-btn
              :loading="loading"
              color="primary"
              variant="elevated"
              size="large"
              type="submit"
            >
              <v-icon
                :icon="mdiPlus"
                start
              />
              生産者を登録
            </v-btn>
          </div>
        </v-container>
      </v-footer>
    </v-form>
  </v-container>
</template>

<style scoped>
.form-section-card {
  border-radius: 12px;
  max-width: 1200px;
  margin-left: auto;
  margin-right: auto;
}

.section-header {
  background: linear-gradient(90deg, rgb(33 150 243 / 5%) 0%, rgb(33 150 243 / 0%) 100%);
  border-bottom: 1px solid rgb(0 0 0 / 5%);
  padding: 20px 24px;
}

@media (width <= 600px) {
  .form-section-card {
    border-radius: 8px;
  }

  .section-header {
    padding: 16px 20px;
  }
}
</style>
