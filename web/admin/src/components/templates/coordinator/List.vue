<script lang="ts" setup>
import { mdiPlus, mdiAccount, mdiDelete } from '@mdi/js'
import type { VDataTable } from 'vuetify/components'
import { convertI18nToJapanesePhoneNumber } from '~/lib/formatter'

import { getResizedImages } from '~/lib/helpers'
import type { AlertType } from '~/lib/hooks'
import { AdminStatus } from '~/types/api/v1'
import type { Shop, Coordinator, ProductType } from '~/types/api/v1'

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
  coordinators: {
    type: Array<Coordinator>,
    default: () => [],
  },
  shops: {
    type: Array<Shop>,
    default: () => [],
  },
  productTypes: {
    type: Array<ProductType>,
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
  (e: 'click:row', coordinatorId: string): void
  (e: 'click:add'): void
  (e: 'click:delete', coordinatorId: string): void
}>()

const headers: VDataTable['headers'] = [
  {
    title: '',
    key: 'thumbnail',
    sortable: false,
  },
  {
    title: 'マルシェ名',
    key: 'marcheName',
    sortable: false,
  },
  {
    title: 'コーディネーター名',
    key: 'username',
    sortable: false,
  },
  {
    title: '生産者数',
    key: 'producerTotal',
    sortable: false,
  },
  {
    title: 'メールアドレス',
    key: 'email',
    sortable: false,
  },
  {
    title: '電話番号',
    key: 'phoneNumber',
    sortable: false,
  },
  {
    title: 'ステータス',
    key: 'status',
    sortable: false,
  },
  {
    title: '',
    key: 'actions',
    sortable: false,
  },
]

const { dialogVisible, selectedItem, open: openDeleteDialog, close: closeDeleteDialog } = useDeleteDialog<Coordinator>()

const getStatus = (status: AdminStatus): string => {
  switch (status) {
    case AdminStatus.AdminStatusActivated:
      return '有効'
    case AdminStatus.AdminStatusInvited:
      return '招待中'
    case AdminStatus.AdminStatusDeactivated:
      return '無効'
    default:
      return '不明'
  }
}

const getStatusColor = (status: AdminStatus): string => {
  switch (status) {
    case AdminStatus.AdminStatusActivated:
      return 'primary'
    case AdminStatus.AdminStatusInvited:
      return 'secondary'
    case AdminStatus.AdminStatusDeactivated:
      return 'error'
    default:
      return 'unknown'
  }
}

const coordinatorName = (coordinator?: Coordinator): string => {
  if (!coordinator) {
    return ''
  }
  return `${coordinator.lastname} ${coordinator.firstname}`
}

const getMarcheName = (coordinator?: Coordinator): string => {
  if (!coordinator) {
    return ''
  }
  const shop = props.shops.find((shop: Shop) => shop.coordinatorId === coordinator.id)
  if (!shop) {
    return ''
  }
  return shop.name
}

const getImages = (coordinator: Coordinator): string => {
  if (!coordinator.thumbnailUrl) {
    return ''
  }
  return getResizedImages(coordinator.thumbnailUrl)
}

const onClickUpdatePage = (page: number): void => {
  emit('click:update-page', page)
}

const onClickUpdateItemsPerPage = (page: number): void => {
  emit('click:update-items-per-page', page)
}

const onClickRow = (coordinatorId: string): void => {
  emit('click:row', coordinatorId)
}

const onClickAdd = (): void => {
  emit('click:add')
}

const onClickDelete = (): void => {
  emit('click:delete', selectedItem.value?.id || '')
  closeDeleteDialog()
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
    title="コーディネーター削除の確認"
    :message="`「${coordinatorName(selectedItem)}」を削除しますか？`"
    :loading="loading"
    @confirm="onClickDelete"
  />

  <v-card
    class="mt-4"
    flat
  >
    <v-card-title class="d-flex flex-row">
      コーディネーター管理
      <v-spacer />
      <v-btn
        variant="outlined"
        color="primary"
        @click="onClickAdd"
      >
        <v-icon
          start
          :icon="mdiPlus"
        />
        コーディネーター登録
      </v-btn>
    </v-card-title>

    <v-card-text>
      <v-data-table-server
        :headers="headers"
        :loading="loading"
        :items="coordinators"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
        hover
        no-data-text="登録されているコーディネーターがいません。"
        @update:page="onClickUpdatePage"
        @update:items-per-page="onClickUpdateItemsPerPage"
        @click:row="(_:any, { item }: any) => onClickRow(item.id)"
      >
        <template #[`item.thumbnail`]="{ item }">
          <v-avatar>
            <v-img
              v-if="item.thumbnailUrl !== ''"
              cover
              :src="item.thumbnailUrl"
              :srcset="getImages(item)"
              :alt="item.username || 'コーディネーター画像'"
            />
            <v-icon
              v-else
              :icon="mdiAccount"
            />
          </v-avatar>
        </template>
        <template #[`item.marcheName`]="{ item }">
          {{ getMarcheName(item) }}
        </template>
        <template #[`item.phoneNumber`]="{ item }">
          {{ convertI18nToJapanesePhoneNumber(item.phoneNumber) }}
        </template>
        <template #[`item.status`]="{ item }">
          <v-chip
            size="small"
            :color="getStatusColor(item.status)"
          >
            {{ getStatus(item.status) }}
          </v-chip>
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn
            variant="outlined"
            color="primary"
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
