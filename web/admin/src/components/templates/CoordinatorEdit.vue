<script lang="ts" setup>
import { mdiPlus } from '@mdi/js'
import { AlertType } from '~/lib/hooks'
import { ProducersResponseProducersInner, UpdateCoordinatorRequest } from '~/types/api'
import { ImageUploadStatus } from '~/types/props'
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
    default: () => ({})
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
  searchErrorMessage: {
    type: String,
    default: ''
  },
  searchLoading: {
    type: Boolean,
    default: false
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

const onChangeThumbnailFile = (files?: FileList): void => {
  if (!files || files.length === 0) {
    return
  }
  emit('update:thumbnail-file', files)
}

const onChangeHeaderFile = (files?: FileList): void => {
  if (!files || files.length === 0) {
    return
  }
  emit('update:header-file', files)
}

const onCilckSearch = (): void => {
  emit('click:search-address')
}

const onClickRelatedProducersPage = (page: number): void => {
  emit('update:related-producers-table-page', page)
}

const onClickRelatedProducersItemsPerPage = (page: number): void => {
  emit('update:related-producers-table-items-per-page', page)
}

const onCloseRelatedProducersDialog = (): void => {
  relatedProducersDialogValue.value = false
}

const onSubmitCoordinator = (): void => {
  emit('submit:coordinator')
}

const onSubmitRelatedProducers = (): void => {
  emit('submit:related-producers')
}
</script>

<template>
  <v-card-title>コーディネーター編集</v-card-title>

  <v-card>
    <v-tabs v-model="selector" grow color="dark">
      <v-tab v-for="item in tabs" :key="item.value" :value="item.value">
        {{ item.name }}
      </v-tab>
    </v-tabs>

    <v-window v-model="selector">
      <v-window-item value="coordinators">
        <v-skeleton-loader v-if="loading" type="article" />

        <organisms-coordinator-edit-form
          v-else
          v-model:form-data="formDataValue"
          :thumbnail-upload-status="thumbnailUploadStatus"
          :header-upload-status="headerUploadStatus"
          :search-loading="searchLoading"
          :search-error-message="searchErrorMessage"
          @update:thumbnail-file="onChangeThumbnailFile"
          @update:header-file="onChangeHeaderFile"
          @click:search="onCilckSearch"
          @submit="onSubmitCoordinator"
        />
      </v-window-item>

      <v-window-item value="relationProducers">
        <v-dialog v-model="relatedProducersDialogValue" width="500">
          <template #activator="{ props: dialog }">
            <div class="d-flex pt-3 pr-3">
              <v-spacer />
              <v-btn variant="outlined" color="primary" v-bind="dialog">
                <v-icon start :icon="mdiPlus" />
                生産者紐付け
              </v-btn>
            </div>
          </template>
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
              <v-btn color="error" variant="text" @click="onCloseRelatedProducersDialog">
                キャンセル
              </v-btn>
              <v-btn color="primary" variant="outlined" @click="onSubmitRelatedProducers">
                更新
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>

        <v-card-text>
          <organisms-related-producer-list
            :producers="relatedProducers"
            :table-items-per-page="relatedProducersTableItemsPerPage"
            :table-items-total="relatedProducersTableItemsTotal"
            @update:page="onClickRelatedProducersPage"
            @update:items-per-page="onClickRelatedProducersItemsPerPage"
          />
        </v-card-text>
      </v-window-item>
    </v-window>
  </v-card>
</template>
