<script lang="ts" setup>
import { mdiPlus, mdiAccount, mdiDelete } from '@mdi/js'
import { VDataTable } from 'vuetify/lib/components/index.mjs'
import { convertI18nToJapanesePhoneNumber } from '~/lib/formatter'

import { getResizedImages } from '~/lib/helpers'
import type { AlertType } from '~/lib/hooks'
import { AdminStatus, type Coordinator, type ProductType } from '~/types/api'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  deleteDialog: {
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
  coordinators: {
    type: Array<Coordinator>,
    default: () => []
  },
  productTypes: {
    type: Array<ProductType>,
    default: () => []
  },
  tableItemsPerPage: {
    type: Number,
    default: 20
  },
  tableItemsTotal: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits<{
  (e: 'click:update-page', page: number): void
  (e: 'click:update-items-per-page', page: number): void
  (e: 'click:row', coordinatorId: string): void
  (e: 'click:add'): void
  (e: 'click:delete', coordinatorId: string): void
  (e: 'update:delete-dialog', v: boolean): void
}>()

const headers: VDataTable['headers'] = [
  {
    title: '',
    key: 'thumbnail',
    sortable: false
  },
  {
    title: 'マルシェ名',
    key: 'marcheName',
    sortable: false
  },
  {
    title: 'コーディネーター名',
    key: 'username',
    sortable: false
  },
  {
    title: '生産者数',
    key: 'producerTotal',
    sortable: false
  },
  {
    title: 'メールアドレス',
    key: 'email',
    sortable: false
  },
  {
    title: '電話番号',
    key: 'phoneNumber',
    sortable: false
  },
  {
    title: 'ステータス',
    key: 'status',
    sortable: false
  },
  {
    title: '',
    key: 'actions',
    sortable: false
  }
]

const selectedItem = ref<Coordinator>()

const deleteDialogValue = computed({
  get: () => props.deleteDialog,
  set: (val: boolean) => emit('update:delete-dialog', val)
})

const getStatus = (status: AdminStatus): string => {
  switch (status) {
    case AdminStatus.ACTIVATED:
      return '有効'
    case AdminStatus.INVITED:
      return '招待中'
    case AdminStatus.DEACTIVATED:
      return '無効'
    default:
      return '不明'
  }
}

const getStatusColor = (status: AdminStatus): string => {
  switch (status) {
    case AdminStatus.ACTIVATED:
      return 'primary'
    case AdminStatus.INVITED:
      return 'secondary'
    case AdminStatus.DEACTIVATED:
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

const onClickOpen = (coordinator: Coordinator): void => {
  selectedItem.value = coordinator
  deleteDialogValue.value = true
}

const onClickClose = (): void => {
  deleteDialogValue.value = false
}

const onClickAdd = (): void => {
  emit('click:add')
}

const onClickDelete = (): void => {
  emit('click:delete', selectedItem?.value?.id || '')
}
</script>

<template>
  <v-alert v-show="props.isAlert" :type="props.alertType" v-text="props.alertText" />

  <v-dialog v-model="deleteDialogValue" width="500">
    <v-card>
      <v-card-title>
        {{ coordinatorName(selectedItem) }}を本当に削除しますか？
      </v-card-title>
      <v-card-actions>
        <v-spacer />
        <v-btn color="error" variant="text" @click="onClickClose">
          キャンセル
        </v-btn>
        <v-btn :loading="loading" color="primary" variant="outlined" @click="onClickDelete">
          削除
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-card class="mt-4" flat>
    <v-card-title class="d-flex flex-row">
      コーディネーター管理
      <v-spacer />
      <v-btn variant="outlined" color="primary" @click="onClickAdd">
        <v-icon start :icon="mdiPlus" />
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
        @click:row="(_:any, {item}: any) => onClickRow(item.id)"
      >
        <template #[`item.thumbnail`]="{ item }">
          <v-avatar>
            <v-img
              v-if="item.thumbnailUrl !== ''"
              cover
              :src="item.thumbnailUrl"
              :srcset="getImages(item)"
            />
            <v-icon v-else :icon="mdiAccount" />
          </v-avatar>
        </template>
        <template #[`item.phoneNumber`]="{ item }">
          {{ convertI18nToJapanesePhoneNumber(item.phoneNumber) }}
        </template>
        <template #[`item.status`]="{ item }">
          <v-chip size="small" :color="getStatusColor(item.status)">
            {{ getStatus(item.status) }}
          </v-chip>
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn
            variant="outlined"
            color="primary"
            size="small"
            @click.stop="onClickOpen(item)"
          >
            <v-icon size="small" :icon="mdiDelete" />
            削除
          </v-btn>
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
