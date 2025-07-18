<script lang="ts" setup>
import { mdiDelete, mdiPlus, mdiContentCopy } from '@mdi/js'
import type { VDataTable } from 'vuetify/lib/components/index.mjs'
import { prefecturesList } from '~/constants'
import { getResizedImages } from '~/lib/helpers'
import type { AlertType } from '~/lib/hooks'
import {

  ExperienceStatus,
  AdminType,
} from '~/types/api'
import type { Experience, ExperienceType, ExperienceMediaInner, Prefecture, ExperiencesResponse, Producer } from '~/types/api'

const props = defineProps({
  selectedItemId: {
    type: String,
    default: '',
  },
  loading: {
    type: Boolean,
    default: false,
  },
  adminType: {
    type: Number as PropType<AdminType>,
    default: AdminType.UNKNOWN,
  },
  deleteDialog: {
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
  producers: {
    type: Array<Producer>,
    defaut: () => [],
  },
  experiencesResponse: {
    type: Array<ExperiencesResponse>,
    default: () => [],
  },
  experiences: {
    type: Array<Experience>,
    default: () => [],
  },
  experienceTypes: {
    type: Array<ExperienceType>,
    default: () => [],
  },
  tableItemsPerPage: {
    type: Number,
    default: 20,
  },
  tableItemsTotal: {
    type: Number,
    default: 0,
  },
})

const emit = defineEmits<{
  (e: 'click:update-page', page: number): void
  (e: 'click:update-items-per-page', page: number): void
  (e: 'click:show', productId: string): void
  (e: 'click:new'): void
  (e: 'click:copyItem'): void
  (e: 'click:delete', productId: string): void
  (e: 'update:delete-dialog', v: boolean): void
  (e: 'update:selectedItemId', v: string): void
}>()

const headers: VDataTable['headers'] = [
  {
    title: '',
    key: 'media',
    width: 80,
    sortable: false,
  },
  {
    title: '商品名',
    key: 'title',
    sortable: false,
  },
  {
    title: 'ステータス',
    key: 'status',
    sortable: false,
  },
  {
    title: 'カテゴリ',
    key: 'experienceType',
    sortable: false,
  },
  {
    title: '場所',
    key: 'place',
    sortable: false,
  },
  {
    title: '生産者名',
    key: 'producerName',
    sortable: false,
  },
  {
    title: '',
    key: 'actions',
    sortable: false,
  },
]

const handleUpdateSelectItemId = (itemIds: string[]): void => {
  if (itemIds.length === 0) {
    emit('update:selectedItemId', '')
  }
  else {
    emit('update:selectedItemId', itemIds[0])
  }
}
console.log(props.experiencesResponse)

const selectedItem = ref<Experience>()

const deleteDialogValue = computed({
  get: (): boolean => props.deleteDialog,
  set: (val: boolean): void => emit('update:delete-dialog', val),
})

const getThumbnail = (media: ExperienceMediaInner[]): string => {
  const thumbnail = media.find((media: ExperienceMediaInner) => {
    return media.isThumbnail
  })
  return thumbnail?.url || ''
}

const getResizedThumbnails = (media: ExperienceMediaInner[]): string => {
  const thumbnail = media.find((media: ExperienceMediaInner) => {
    return media.isThumbnail
  })
  if (!thumbnail) {
    return ''
  }
  return getResizedImages(thumbnail.url)
}

const getExperienceType = (experienceTypeId: string): string => {
  const experienceType = props.experienceTypes.find((experienceType: ExperienceType): boolean => {
    return experienceType.id === experienceTypeId
  })
  return experienceType ? experienceType.name : ''
}

const isRegisterable = (): boolean => {
  return props.adminType === AdminType.COORDINATOR
}

const isDeletable = (): boolean => {
  const targets: AdminType[] = [AdminType.ADMINISTRATOR, AdminType.COORDINATOR]
  return targets.includes(props.adminType)
}

const toggleDeleteDialog = (experience?: Experience): void => {
  if (experience) {
    selectedItem.value = experience
  }
  deleteDialogValue.value = !deleteDialogValue.value
}

const onUpdatePage = (page: number): void => {
  emit('click:update-page', page)
}

const onUpdateItemsPerPage = (page: number): void => {
  emit('click:update-items-per-page', page)
}

const onClickShow = (productId: string): void => {
  emit('click:show', productId)
}

const onClickNew = (): void => {
  emit('click:new')
}

const onClickDelete = (): void => {
  emit('click:delete', selectedItem?.value?.id || '')
}

const getProducerName = (producerId: string): string => {
  const producer = props.producers?.find((producer: Producer): boolean => {
    return producer.id === producerId
  })
  return producer ? producer.username : ''
}

const getStatus = (status: ExperienceStatus): string => {
  switch (status) {
    case ExperienceStatus.WAITING:
      return '販売開始前'
    case ExperienceStatus.ACCEPTING:
      return '体験受付中'
    case ExperienceStatus.SOLD_OUT:
      return '体験受付終了'
    case ExperienceStatus.PRIVATE:
      return '非公開'
    case ExperienceStatus.FINISHED:
      return '販売終了'
    default:
      return ''
  }
}

const getStatusColor = (status: ExperienceStatus): string => {
  switch (status) {
    case ExperienceStatus.WAITING:
      return 'info'
    case ExperienceStatus.ACCEPTING:
      return 'primary'
    case ExperienceStatus.SOLD_OUT:
      return 'secondary'
    case ExperienceStatus.PRIVATE:
      return 'warning'
    case ExperienceStatus.FINISHED:
      return 'error'
    default:
      return ''
  }
}

const getPrefecture = (hostPrefectureCode: Prefecture): string => {
  const item = prefecturesList.find(prefecture => prefecture.value === hostPrefectureCode)
  return item ? item.text : ''
}
</script>

<template>
  <v-alert
    v-show="props.isAlert"
    :type="props.alertType"
    v-text="props.alertText"
  />

  <v-dialog
    v-model="deleteDialogValue"
    width="500"
  >
    <v-card>
      <v-card-text class="text-h7">
        {{ selectedItem?.title }}を本当に削除しますか？
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn
          color="error"
          variant="text"
          @click="toggleDeleteDialog"
        >
          キャンセル
        </v-btn>
        <v-btn
          :loading="loading"
          color="primary"
          variant="outlined"
          @click="onClickDelete"
        >
          削除
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-card
    class="mt-4"
    flat
  >
    <v-card-title class="d-flex flex-row">
      体験管理
    </v-card-title>

    <div class="d-flex w-100 px-6 ga-2">
      <v-spacer />
      <v-btn
        v-show="isRegisterable()"
        variant="outlined"
        color="primary"
        @click="onClickNew"
      >
        <v-icon
          start
          :icon="mdiPlus"
        />
        体験登録
      </v-btn>
    </div>

    <v-card-text>
      <v-data-table-server
        :headers="headers"
        :loading="loading"
        :items="props.experiencesResponse.experiences"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
        hover
        select-strategy="single"
        show-select
        no-data-text="登録されている商品がありません。"
        @update:model-value="handleUpdateSelectItemId"
        @update:page="onUpdatePage"
        @update:items-per-page="onUpdateItemsPerPage"
        @click:row="(_: any, { item }: any) => onClickShow(item.id)"
      >
        <template #[`item.media`]="{ item }">
          <v-img
            aspect-ratio="1/1"
            :max-height="56"
            :max-width="80"
            :src="getThumbnail(item.media)"
            :srcset="getResizedThumbnails(item.media)"
          />
        </template>
        <template #[`item.status`]="{ item }">
          <v-chip :color="getStatusColor(item.status)">
            {{ getStatus(item.status) }}
          </v-chip>
        </template>
        <template #[`item.experienceType`]="{ item }">
          {{ getExperienceType(item.experienceTypeId) }}
        </template>
        <template #[`item.place`]="{ item }">
          {{ getPrefecture(item.hostPrefectureCode) }}
        </template>
        <template #[`item.producerName`]="{ item }">
          {{ getProducerName(item.producerId) }}
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn
            v-show="isDeletable()"
            variant="outlined"
            color="primary"
            size="small"
            @click.stop="toggleDeleteDialog(item)"
          >
            <v-icon
              size="small"
              :icon="mdiDelete"
            />
            削除
          </v-btn>
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
