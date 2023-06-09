<script lang="ts" setup>
import { mdiPlus } from '@mdi/js'
import useVuelidate from '@vuelidate/core'
import { AlertType } from '~/lib/hooks'
import { ProducersResponseProducersInner, UpdateCoordinatorRequest } from '~/types/api'
import { ImageUploadStatus } from '~/types/props'
import { email, getErrorMessage, kana, maxLength, required, tel } from '~/lib/validations'
import { Coordinator } from '~/types/props/coordinator'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  isAlert: {
    type: Boolean,
    default: false
  },
  alertType: {
    type: String as PropType<AlertType>,
    default: undefined
  },
  alertText: {
    type: String,
    default: ''
  },
  relatedProducersDialog: {
    type: Boolean,
    default: false
  },
  formData: {
    type: Object as PropType<UpdateCoordinatorRequest>,
    default: (): UpdateCoordinatorRequest => ({
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
      addressLine2: ''
    })
  },
  selectedProducerIds: {
    type: Array<string>,
    default: () => []
  },
  relatedProducers: {
    type: Array<ProducersResponseProducersInner>,
    default: () => []
  },
  unrelatedProducers: {
    type: Array<ProducersResponseProducersInner>,
    default: () => []
  },
  thumbnailUploadStatus: {
    type: Object,
    default: (): ImageUploadStatus => ({
      error: false,
      message: ''
    })
  },
  headerUploadStatus: {
    type: Object,
    default: (): ImageUploadStatus => ({
      error: false,
      message: ''
    })
  },
  relatedProducersTableItemsPerPage: {
    type: Number,
    default: 20
  },
  relatedProducersTableItemsTotal: {
    type: Number,
    default: 20
  }
})

const emit = defineEmits<{
  (e: 'update:form-data', v: UpdateCoordinatorRequest): void
  (e: 'update:thumbnail-file', files: FileList): void
  (e: 'update:header-file', files: FileList): void
  (e: 'update:selected-producer-ids', producerIds: string[]): void
  (e: 'update:related-producers-dialog', dialog: boolean): void
  (e: 'update:related-producers-table-page', page: number): void
  (e: 'update:related-producers-table-items-per-page', page: number): void
  (e: 'click:search-address'): void
  (e: 'submit:coordinator'): void
  (e: 'submit:related-producers'): void
}>()

const tabs: Coordinator[] = [
  { name: '基本情報', value: 'coordinators' },
  { name: '関連生産者', value: 'relationProducers' }
]

const selector = ref<string>('coordinators')

const rules = computed(() => ({
  lastname: { required, maxLength: maxLength(16) },
  lastnameKana: { required, kana },
  firstname: { required, maxLength: maxLength(16) },
  firstnameKana: { required, kana },
  storeName: { required, maxLength: maxLength(64) },
  companyName: { required, maxLength: maxLength(64) },
  phoneNumber: { required, tel }
}))
const formDataValue = computed({
  get: (): UpdateCoordinatorRequest => props.formData,
  set: (val: UpdateCoordinatorRequest): void => emit('update:form-data', val)
})
const relatedProducersDialogValue = computed({
  get: (): boolean => props.relatedProducersDialog,
  set: (val: boolean): void => emit('update:related-producers-dialog', val)
})
const selectedProducerIdsValue = computed({
  get: (): string[] => props.selectedProducerIds,
  set: (producerIds: string[]): void => emit('update:selected-producer-ids', producerIds)
})

const validate = useVuelidate(rules, formDataValue)

const onChangeThumbnailFile = (files?: FileList): void => {
  if (!files) {
    return
  }

  emit('update:thumbnail-file', files)
}

const onChangeHeaderFile = (files?: FileList): void => {
  if (!files) {
    return
  }

  emit('update:header-file', files)
}

const onClickRelatedProducersPage = (page: number): void => {
  emit('update:related-producers-table-page', page)
}

const onClickRelatedProducersItemsPerPage = (page: number): void => {
  emit('update:related-producers-table-items-per-page', page)
}

const onClickOpenRelatedProducersDialog = (): void => {
  relatedProducersDialogValue.value = true
}

const onClickCloseRelatedProducersDialog = (): void => {
  relatedProducersDialogValue.value = false
}

const onSubmitCoordinator = async (): Promise<void> => {
  const valid = await validate.value.$validate()
  if (!valid) {
    return
  }

  emit('submit:coordinator')
}

const onSubmitRelatedProducers = (): void => {
  emit('submit:related-producers')
}

const onCilckSearchAddress = (): void => {
  emit('click:search-address')
}
</script>

<template>
  <v-dialog v-model="relatedProducersDialogValue" width="500">
    <v-card>
      <v-card-title class="primaryLight">
        生産者紐付け
      </v-card-title>

      <v-autocomplete
        v-model="selectedProducerIdsValue"
        :items="unrelatedProducers"
        chips
        closable-chips
        item-title="firstname"
        item-value="id"
        label="関連生産者"
        multiple
      >
        <template #chip="{ props: chip, item }">
          <v-chip v-bind="chip" :prepend-avatar="item.raw.thumbnailUrl" :text="item.raw.firstname" />
        </template>
        <template #item="{ props: list, item }">
          <v-list-item
            v-bind="list"
            :prepend-avatar="item.raw.thumbnailUrl"
            :title="item.raw.firstname"
            :subtitle="item.raw.storeName"
          />
        </template>
      </v-autocomplete>

      <v-card-actions>
        <v-spacer />
        <v-btn color="error" variant="text" @click="onClickCloseRelatedProducersDialog">
          キャンセル
        </v-btn>
        <v-btn color="primary" variant="outlined" @click="onSubmitRelatedProducers">
          更新
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-card>
    <v-card-title class="d-flex flex-row">
      コーディネーター編集
      <v-spacer />
      <v-btn variant="outlined" color="primary" @click="onClickOpenRelatedProducersDialog">
        <v-icon start :icon="mdiPlus" />
        生産者紐付け
      </v-btn>
    </v-card-title>

    <v-tabs v-model="selector" grow color="dark">
      <v-tab v-for="item in tabs" :key="item.value" :value="item.value">
        {{ item.name }}
      </v-tab>
    </v-tabs>

    <v-window v-model="selector">
      <v-window-item value="coordinators">
        <v-skeleton-loader v-if="loading" type="article" />

        <v-form @submit.prevent="onSubmitCoordinator">
          <v-card-text>
            <v-text-field
              v-model="validate.companyName.$model"
              :error-messages="getErrorMessage(validate.companyName.$errors)"
              label="会社名"
            />
            <v-text-field
              v-model="validate.storeName.$model"
              :error-messages="getErrorMessage(validate.storeName.$errors)"
              label="店舗名"
            />
            <div class="mb-2 d-flex">
              <molecules-profile-select-form
                class="mr-4 flex-grow-1 flex-shrink-1"
                :img-url="formDataValue.thumbnailUrl"
                :error="props.thumbnailUploadStatus.error"
                :message="props.thumbnailUploadStatus.message"
                @update:file="onChangeThumbnailFile"
              />
              <molecules-header-select-form
                class="flex-grow-1 flex-shrink-1"
                :img-url="formDataValue.headerUrl"
                :error="props.headerUploadStatus.error"
                :message="props.headerUploadStatus.message"
                @update:file="onChangeHeaderFile"
              />
            </div>
            <div class="d-flex">
              <v-text-field
                v-model="validate.lastname.$model"
                :error-messages="getErrorMessage(validate.lastname.$errors)"
                class="mr-4"
                label="コーディネータ:姓"
              />
              <v-text-field
                v-model="validate.firstname.$model"
                :error-messages="getErrorMessage(validate.firstname.$errors)"
                label="コーディネータ:名"
              />
            </div>
            <div class="d-flex">
              <v-text-field
                v-model="validate.lastnameKana.$model"
                :error-messages="getErrorMessage(validate.lastnameKana.$errors)"
                class="mr-4"
                label="コーディネータ:姓（ふりがな）"
              />
              <v-text-field
                v-model="validate.firstnameKana.$model"
                :error-messages="getErrorMessage(validate.firstnameKana.$errors)"
                label="コーディネータ:名（ふりがな）"
              />
            </div>
            <v-text-field
              v-model="validate.phoneNumber.$model"
              :error-messages="getErrorMessage(validate.phoneNumber.$errors)"
              type="tel"
              label="連絡先（電話番号）"
            />
            <molecules-address-form
              v-model:postal-code="formDataValue.postalCode"
              v-model:prefecture="formDataValue.prefecture"
              v-model:city="formDataValue.city"
              v-model:address-line1="formDataValue.addressLine1"
              v-model:address-line2="formDataValue.addressLine2"
              @click:search="onCilckSearchAddress"
            />
          </v-card-text>
          <v-card-actions>
            <v-btn block variant="outlined" color="primary" type="submit">
              更新
            </v-btn>
          </v-card-actions>
        </v-form>
      </v-window-item>

      <v-window-item value="relationProducers">
        <organisms-related-producer-list
          :producers="relatedProducers"
          :table-items-per-page="relatedProducersTableItemsPerPage"
          :table-items-total="relatedProducersTableItemsTotal"
          @update:page="onClickRelatedProducersPage"
          @update:items-per-page="onClickRelatedProducersItemsPerPage"
        />
      </v-window-item>
    </v-window>
  </v-card>
</template>
