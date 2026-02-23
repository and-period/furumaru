<script lang="ts" setup>
import { mdiDelete, mdiPlus, mdiContentCopy, mdiCalendarCheck, mdiAccount, mdiTent } from '@mdi/js'
import type { VDataTable } from 'vuetify/components'
import { experienceStatues, prefecturesList } from '~/constants'
import { getResizedImages } from '~/lib/helpers'
import type { AlertType } from '~/lib/hooks'
import {

  ExperienceStatus,
  AdminType,
} from '~/types/api/v1'
import type { Prefecture } from '~/types'
import type { Experience, ExperienceType, ExperienceMedia, ExperiencesResponse, Producer } from '~/types/api/v1'

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
    default: AdminType.AdminTypeUnknown,
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
  (e: 'click:copy-item'): void
  (e: 'click:delete', productId: string): void
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

const getStatus = (status: ExperienceStatus): string => {
  const value = experienceStatues.find(s => s.value === status)
  return value ? value.title : ''
}

const getStatusColor = (status: ExperienceStatus): string => {
  switch (status) {
    case ExperienceStatus.ExperienceStatusWaiting:
      return 'info'
    case ExperienceStatus.ExperienceStatusAccepting:
      return 'primary'
    case ExperienceStatus.ExperienceStatusSoldOut:
      return 'secondary'
    case ExperienceStatus.ExperienceStatusPrivate:
      return 'warning'
    case ExperienceStatus.ExperienceStatusFinished:
      return 'error'
    default:
      return ''
  }
}

const handleUpdateSelectItemId = (itemIds: string[]): void => {
  if (itemIds.length === 0) {
    emit('update:selectedItemId', '')
  }
  else {
    emit('update:selectedItemId', itemIds[0])
  }
}

const { dialogVisible, selectedItem, open: openDeleteDialog, close: closeDeleteDialog } = useDeleteDialog<Experience>()

const getThumbnail = (media: ExperienceMedia[]): string => {
  const thumbnail = media.find((media: ExperienceMedia) => {
    return media.isThumbnail
  })
  return thumbnail?.url || ''
}

const getResizedThumbnails = (media: ExperienceMedia[]): string => {
  const thumbnail = media.find((media: ExperienceMedia) => {
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
  return props.adminType === AdminType.AdminTypeCoordinator
}

const isDeletable = (): boolean => {
  const targets: AdminType[] = [AdminType.AdminTypeAdministrator, AdminType.AdminTypeCoordinator]
  return targets.includes(props.adminType)
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

const onClickCopyItem = (): void => {
  emit('click:copy-item')
}

const onClickDelete = (): void => {
  emit('click:delete', selectedItem.value?.id || '')
  closeDeleteDialog()
}

const getProducerName = (producerId: string): string => {
  const producer = props.producers?.find((producer: Producer): boolean => {
    return producer.id === producerId
  })
  return producer ? producer.username : ''
}

const getPrefecture = (hostPrefectureCode: Prefecture): string => {
  const item = prefecturesList.find(prefecture => prefecture.value === hostPrefectureCode)
  return item ? item.text : ''
}
</script>

<template>
  <atoms-app-alert
    :show="props.isAlert"
    :type="props.alertType"
    :text="props.alertText"
  />

  <atoms-app-confirm-dialog
    v-model="dialogVisible"
    title="体験削除の確認"
    :message="`「${selectedItem?.title}」を削除しますか？`"
    :loading="loading"
    @confirm="onClickDelete"
  />

  <v-card
    class="mt-6"
    elevation="0"
    rounded="lg"
  >
    <v-card-title class="d-flex flex-column flex-sm-row align-start align-sm-center justify-space-between pa-4 pa-sm-6 pb-4">
      <div class="d-flex align-center mb-3 mb-sm-0">
        <v-icon
          :icon="mdiCalendarCheck"
          size="28"
          class="mr-3 text-primary"
        />
        <div>
          <h1 class="text-h6 text-sm-h5 font-weight-bold text-primary">
            体験管理
          </h1>
          <p class="text-caption text-sm-body-2 text-medium-emphasis ma-0">
            体験の登録・編集・削除を行います
          </p>
        </div>
      </div>
      <div class="d-flex flex-column flex-sm-row ga-2 ga-sm-3 w-100 w-sm-auto">
        <v-btn
          v-show="isRegisterable()"
          variant="outlined"
          color="secondary"
          :size="$vuetify.display.smAndDown ? 'default' : 'large'"
          class="w-100 w-sm-auto"
          :disabled="selectedItemId === ''"
          @click="onClickCopyItem"
        >
          <v-icon
            start
            :icon="mdiContentCopy"
          />
          体験複製
          <v-tooltip
            v-if="selectedItemId === ''"
            activator="parent"
            location="bottom"
          >
            コピー元の体験をチェックする必要があります。
          </v-tooltip>
        </v-btn>
        <v-btn
          v-show="isRegisterable()"
          variant="elevated"
          color="primary"
          :size="$vuetify.display.smAndDown ? 'default' : 'large'"
          class="w-100 w-sm-auto"
          @click="onClickNew"
        >
          <v-icon
            start
            :icon="mdiPlus"
          />
          体験登録
        </v-btn>
      </div>
    </v-card-title>

    <v-card-text>
      <v-data-table-server
        :headers="headers"
        :loading="loading"
        :items="props.experiences"
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
          <v-avatar size="40">
            <v-img
              v-if="getThumbnail(item.media) !== ''"
              cover
              :src="getThumbnail(item.media)"
              :srcset="getResizedThumbnails(item.media)"
              :alt="item.title || '体験画像'"
            />
            <v-icon
              v-else
              :icon="mdiTent"
              color="grey"
            />
          </v-avatar>
        </template>
        <template #[`item.title`]="{ item }">
          {{ item.title.length > 24 ? item.title.slice(0, 24) + '...' : item.title }}
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
            color="error"
            size="small"
            @click.stop="openDeleteDialog(item)"
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
